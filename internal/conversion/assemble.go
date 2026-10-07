package conversion

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/JSYoo5B/convertago/internal/validation"
)

// Part is one source contribution to a platform-specific input slot.
type Part struct {
	Slot   string
	Text   string
	Style  []string
	Format string
	Path   string
}

// Node represents an assembled platform builder before native rendering.
type Node struct {
	Role   string
	Path   string
	Parts  []Part
	Inputs []Input
	// Fixed holds tag values for slots; an input to the same slot replaces its value.
	Fixed   []Part
	Slot    string
	present bool
	partial bool
}

// Input preserves the source order of scalar parts and native child builders.
type Input struct {
	Slot  string
	Part  *Part
	Child *Node
}

// Children returns the native builders assigned to a containing input slot.
func (n Node) Children(slot string) []Node {
	var children []Node
	for _, input := range n.Inputs {
		if input.Slot == slot && input.Child != nil {
			children = append(children, *input.Child)
		}
	}
	return children
}

func (n *Node) appendInputs(inputs []Input) {
	n.Inputs = append(n.Inputs, inputs...)
	for _, input := range inputs {
		if input.Part != nil {
			n.Parts = append(n.Parts, *input.Part)
		}
	}
}

// Text concatenates the scalar parts supplied to a slot, or returns its fixed tag value.
func (n Node) Text(slot string) string {
	var text strings.Builder
	supplied := false
	for _, part := range n.Parts {
		if part.Slot == slot {
			text.WriteString(part.Text)
			supplied = true
		}
	}
	if fixed, ok := n.fixed(slot); ok && !supplied {
		return fixed.Text
	}
	return text.String()
}

// Has reports whether any part, child, or fixed tag value was supplied to a slot.
func (n Node) Has(slot string) bool {
	for _, input := range n.Inputs {
		if input.Slot == slot {
			return true
		}
	}
	_, ok := n.fixed(slot)
	return ok
}

func (n Node) fixed(slot string) (Part, bool) {
	for _, part := range n.Fixed {
		if part.Slot == slot {
			return part, true
		}
	}
	return Part{}, false
}

// fixedParts records a tag's fixed values as parts at the tag's path.
func fixedParts(tag Tag, path string) []Part {
	if len(tag.Fixed) == 0 {
		return nil
	}
	parts := make([]Part, len(tag.Fixed))
	for i, fixed := range tag.Fixed {
		parts[i] = Part{Slot: fixed.Slot, Text: fixed.Value, Path: path}
	}
	return parts
}

// Prepare reads a source and assembles its top-level builder nodes in declaration order.
func Prepare(input any, platform string, options []Option) ([]Node, error) {
	fields, err := Fields(input, platform)
	if err != nil {
		return nil, err
	}
	profile, err := Lookup(platform)
	if err != nil {
		return nil, err
	}
	a := assembler{profile: profile, options: Configure(options)}
	nodes, err := a.fields(fields, "", Tag{}, "$")
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		if err := a.checkSlots(node); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

type assembler struct {
	profile Profile
	options Options
}

func (a assembler) fail(path, code, message string) error {
	return Error(a.profile.Platform, path, code, message)
}

func (a assembler) fields(fields []Field, builder string, inherited Tag, path string) ([]Node, error) {
	groups := make(map[string]string)
	for _, field := range fields {
		fieldPath := path + "." + field.Name
		if err := validateContext(a.profile, field.Tag, builder); err != nil {
			return nil, a.fail(fieldPath, TagErrorCode(err), err.Error())
		}
		if err := checkGroup(groups, field.Tag, builder, a.profile); err != nil {
			return nil, a.fail(fieldPath, "group_conflict", err.Error())
		}
		if field.Value.shape != nil {
			if err := validateNativeShape(a.profile, field.Tag, builder, field.Value.shape, fieldPath); err != nil {
				return nil, err
			}
			next := builder
			if field.Tag.Role != "part" && field.Tag.Role != "flatten" {
				next = field.Tag.Role
			}
			if err := validateShape(a.profile, field.Value.shape, next, fieldPath, make(map[shapeVisit]bool)); err != nil {
				return nil, err
			}
		}
	}
	var nodes []Node
	anchors := make(map[string]int)
	for _, field := range fields {
		fieldPath := path + "." + field.Name
		tag := field.Tag
		if tag.Role == "" {
			continue
		}
		partial := tag.Role == "part"
		target := ""
		if builder != "" && tag.Role != "flatten" {
			target, _ = inputSlot(a.profile, tag, builder)
		}
		if partial {
			tag.Role = builder
			if tag.Slot == "" {
				tag.Slot = inherited.Slot
			}
			if tag.Slot == "" || tag.Slot == a.profile.Roles[builder].DefaultSlot {
				tag.Style = mergeStyles(inherited.Style, tag.Style)
				if tag.Format == "" {
					tag.Format = inherited.Format
				}
			}
		} else if builder != "" {
			tag.Slot = ""
		}
		// Reserve the first declared member's position before checking omissions.
		anchor := -1
		if tag.Group != "" {
			var exists bool
			anchor, exists = anchors[tag.Group]
			if !exists {
				anchor = len(nodes)
				anchors[tag.Group] = anchor
				nodes = append(nodes, Node{Role: tag.Role, Path: fieldPath, Fixed: fixedParts(tag, fieldPath), Slot: target, partial: partial})
			}
		}
		unavailable, err := CheckTag(a.profile, tag)
		if err != nil {
			return nil, a.fail(fieldPath, TagErrorCode(err), err.Error())
		}
		if field.Value.Err != nil {
			return nil, a.fail(fieldPath, "invalid_source", field.Value.Err.Error())
		}
		if field.Value.Nil || (tag.OmitEmpty && field.Value.Empty) {
			continue
		}
		if unavailable != "" {
			diagnostic := Diagnostic{a.profile.Platform, fieldPath, "unsupported_feature", unavailable + " is unavailable in this converter", validation.Fatal}
			if tag.Optional {
				diagnostic.Severity = validation.Warning
			}
			if err := a.options.Handle(diagnostic); err != nil {
				return nil, err
			}
			continue
		}
		if tag.Role == "flatten" {
			flattened, err := a.flatten(field.Value, builder, inherited, fieldPath)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, flattened...)
			continue
		}
		if anchor >= 0 {
			if !partial && nodes[anchor].Slot != target {
				return nil, a.fail(fieldPath, "group_conflict", "group mixes containing input slots")
			}
			inputs, err := a.inputs(field.Value, tag, fieldPath)
			if err != nil {
				return nil, err
			}
			nodes[anchor].appendInputs(inputs)
			nodes[anchor].present = nodes[anchor].present || len(inputs) != 0 || field.Value.Kind != "list"
			continue
		}
		values := []Value{field.Value}
		if !partial && field.Value.Kind == "list" {
			values = field.Value.Items
		}
		for i, value := range values {
			itemPath := fieldPath
			if !partial && field.Value.Kind == "list" {
				itemPath = fmt.Sprintf("%s[%d]", fieldPath, i)
			}
			if value.Nil {
				continue
			}
			inputs, err := a.inputs(value, tag, itemPath)
			if err != nil {
				return nil, err
			}
			node := Node{Role: tag.Role, Path: itemPath, Fixed: fixedParts(tag, itemPath), Slot: target, partial: partial, present: true}
			node.appendInputs(inputs)
			nodes = append(nodes, node)
		}
	}
	result := nodes[:0]
	for _, node := range nodes {
		if node.present {
			result = append(result, node)
		}
	}
	return result, nil
}

func (a assembler) inputs(value Value, tag Tag, path string) ([]Input, error) {
	if value.Err != nil {
		return nil, a.fail(path, "invalid_source", value.Err.Error())
	}
	if value.Nil {
		return nil, nil
	}
	if value.shape != nil {
		if err := validateShape(a.profile, value.shape, tag.Role, path, make(map[shapeVisit]bool)); err != nil {
			return nil, err
		}
	}
	switch value.Kind {
	case "scalar":
		slot := tag.Slot
		if slot == "" {
			slot = a.profile.Roles[tag.Role].DefaultSlot
		}
		if slot == "" {
			return nil, a.fail(path, "invalid_source", tag.Role+" requires a tagged struct")
		}
		part := &Part{slot, value.Text, tag.Style, tag.Format, path}
		return []Input{{Slot: slot, Part: part}}, nil
	case "object":
		nodes, err := a.fields(value.Fields, tag.Role, tag, path)
		if err != nil {
			return nil, err
		}
		var inputs []Input
		for _, node := range nodes {
			if node.partial {
				inputs = append(inputs, node.Inputs...)
			} else {
				child := node
				inputs = append(inputs, Input{Slot: node.Slot, Child: &child})
			}
		}
		if len(inputs) == 0 && !a.profile.Roles[tag.Role].EmptyAllowed {
			return nil, a.fail(path, "missing_input", "builder has no inputs")
		}
		return inputs, nil
	case "list":
		var inputs []Input
		for i, item := range value.Items {
			child, err := a.inputs(item, tag, fmt.Sprintf("%s[%d]", path, i))
			if err != nil {
				return nil, err
			}
			inputs = append(inputs, child...)
		}
		return inputs, nil
	default:
		return nil, a.fail(path, "invalid_source", "expected text, a tagged struct, or a list")
	}
}

func (a assembler) flatten(value Value, builder string, inherited Tag, path string) ([]Node, error) {
	if value.shape != nil {
		if err := validateShape(a.profile, value.shape, builder, path, make(map[shapeVisit]bool)); err != nil {
			return nil, err
		}
	}
	if value.Err != nil {
		return nil, a.fail(path, "invalid_source", value.Err.Error())
	}
	if value.Nil {
		return nil, nil
	}
	if value.Kind == "object" {
		return a.fields(value.Fields, builder, inherited, path)
	}
	if value.Kind == "list" {
		var nodes []Node
		for i, item := range value.Items {
			child, err := a.flatten(item, builder, inherited, fmt.Sprintf("%s[%d]", path, i))
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, child...)
		}
		return nodes, nil
	}
	return nil, a.fail(path, "invalid_source", "flatten requires a struct or a list of structs")
}

func (a assembler) checkSlots(node Node) error {
	role := a.profile.Roles[node.Role]
	counts := make(map[string]int)
	for _, input := range node.Inputs {
		slot, exists := role.Slots[input.Slot]
		path := node.Path
		if input.Part != nil {
			path = input.Part.Path
		} else {
			path = input.Child.Path
		}
		if !exists {
			return a.fail(path, "invalid_tag", fmt.Sprintf("unknown slot %q for %s", input.Slot, node.Role))
		}
		if input.Part != nil {
			if len(slot.Children) != 0 && !slot.Scalar {
				return a.fail(path, "invalid_source", "slot "+input.Slot+" requires a child builder")
			}
			if !utf8.ValidString(input.Part.Text) {
				return a.fail(path, "invalid_value", "text must be valid UTF-8")
			}
		} else {
			if !contains(slot.Children, input.Child.Role) {
				return a.fail(path, "invalid_tag", "invalid child for slot "+input.Slot)
			}
			if err := a.checkSlots(*input.Child); err != nil {
				return err
			}
		}
		counts[input.Slot]++
		if counts[input.Slot] > 1 && !slot.Repeated {
			return a.fail(path, "duplicate_input", "duplicate input for slot "+input.Slot)
		}
	}
	var names []string
	for name := range role.Slots {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if _, fixed := node.fixed(name); role.Slots[name].Required && counts[name] == 0 && !fixed {
			return a.fail(node.Path, "missing_input", "required slot "+name+" is missing")
		}
	}
	return nil
}

func mergeStyles(parent, child []string) []string {
	styles := append([]string(nil), parent...)
	for _, style := range child {
		if !contains(styles, style) {
			styles = append(styles, style)
		}
	}
	return styles
}
