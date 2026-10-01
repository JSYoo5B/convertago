package conversion

import (
	"fmt"
	"sort"
	"strings"
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
	Role  string
	Path  string
	Parts []Part
}

func (n Node) Text(slot string) string {
	var text strings.Builder
	for _, part := range n.Parts {
		if part.Slot == slot {
			text.WriteString(part.Text)
		}
	}
	return text.String()
}

func (n Node) Has(slot string) bool {
	for _, part := range n.Parts {
		if part.Slot == slot {
			return true
		}
	}
	return false
}

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
			return nil, a.fail(fieldPath, "invalid_tag", err.Error())
		}
		if err := checkGroup(groups, field.Tag); err != nil {
			return nil, a.fail(fieldPath, "group_conflict", err.Error())
		}
		if field.Value.shape != nil {
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
		if tag.Role == "part" {
			tag.Role = builder
			tag.Style = mergeStyles(inherited.Style, tag.Style)
			if tag.Slot == "" {
				tag.Slot = inherited.Slot
			}
			if tag.Format == "" {
				tag.Format = inherited.Format
			}
		}
		// Reserve the first declared member's position before checking omissions.
		anchor := -1
		if tag.Group != "" {
			var exists bool
			anchor, exists = anchors[tag.Group]
			if !exists {
				anchor = len(nodes)
				anchors[tag.Group] = anchor
				nodes = append(nodes, Node{Role: tag.Role, Path: fieldPath})
			}
		}
		unavailable, err := CheckTag(a.profile, tag)
		if err != nil {
			return nil, a.fail(fieldPath, "invalid_tag", err.Error())
		}
		if field.Value.Err != nil {
			return nil, a.fail(fieldPath, "invalid_source", field.Value.Err.Error())
		}
		if field.Value.Nil || (tag.OmitEmpty && field.Value.Empty) {
			continue
		}
		if unavailable != "" {
			diagnostic := Diagnostic{a.profile.Platform, fieldPath, "unsupported_feature", unavailable + " is unavailable in this converter"}
			if !tag.Optional || a.options.Strict {
				return nil, diagnostic
			}
			if a.options.Diagnostic != nil {
				a.options.Diagnostic(diagnostic)
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
			parts, err := a.parts(field.Value, tag, fieldPath)
			if err != nil {
				return nil, err
			}
			nodes[anchor].Parts = append(nodes[anchor].Parts, parts...)
			continue
		}
		if builder == "" && field.Value.Kind == "list" {
			for i, item := range field.Value.Items {
				itemPath := fmt.Sprintf("%s[%d]", fieldPath, i)
				parts, err := a.parts(item, tag, itemPath)
				if err != nil {
					return nil, err
				}
				nodes = append(nodes, Node{Role: tag.Role, Path: itemPath, Parts: parts})
			}
		} else {
			parts, err := a.parts(field.Value, tag, fieldPath)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, Node{Role: tag.Role, Path: fieldPath, Parts: parts})
		}
	}
	result := nodes[:0]
	for _, node := range nodes {
		if len(node.Parts) != 0 {
			result = append(result, node)
		}
	}
	return result, nil
}

func (a assembler) parts(value Value, tag Tag, path string) ([]Part, error) {
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
		return []Part{{slot, value.Text, tag.Style, tag.Format, path}}, nil
	case "object":
		nodes, err := a.fields(value.Fields, tag.Role, tag, path)
		if err != nil {
			return nil, err
		}
		var parts []Part
		for _, node := range nodes {
			parts = append(parts, node.Parts...)
		}
		// A present builder with no inputs is an error, rather than an omission.
		if len(parts) == 0 {
			return nil, a.fail(path, "missing_input", "builder has no inputs")
		}
		return parts, nil
	case "list":
		var parts []Part
		for i, item := range value.Items {
			child, err := a.parts(item, tag, fmt.Sprintf("%s[%d]", path, i))
			if err != nil {
				return nil, err
			}
			parts = append(parts, child...)
		}
		return parts, nil
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
	for _, part := range node.Parts {
		if err := ValidateText(a.profile.Platform, part.Path, part.Text, 0, 0); err != nil {
			return err
		}
		slot, exists := role.Slots[part.Slot]
		if !exists {
			return a.fail(part.Path, "invalid_tag", fmt.Sprintf("unknown slot %q for %s", part.Slot, node.Role))
		}
		counts[part.Slot]++
		if counts[part.Slot] > 1 && !slot.Repeated {
			return a.fail(part.Path, "duplicate_input", "duplicate input for slot "+part.Slot)
		}
	}
	var names []string
	for name := range role.Slots {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if role.Slots[name].Required && counts[name] == 0 {
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
