package conversion

import (
	"strconv"
	"strings"

	"github.com/JSYoo5B/convertago/internal/validation"
)

// Checker applies native rules to values built during one conversion.
type Checker struct {
	Platform string
	Options  Options
	// Slots maps native field names to differently named input slots, such as "alt_text" to "alt".
	Slots map[string]string
}

// Check runs the rules of a native value built from node.
// Children are not visited, because each child is checked when it is built.
func (c Checker) Check(node Node, check func(*validation.Check)) error {
	return c.CheckAt(func(field string) string {
		if slot, ok := c.Slots[field]; ok {
			field = slot
		}
		return node.FieldPath(field)
	}, check)
}

// CheckAt runs the rules of a native value, resolving each field to a source path.
func (c Checker) CheckAt(resolve func(field string) string, check func(*validation.Check)) error {
	var err error
	check(&validation.Check{Report: func(v validation.Violation) {
		if err == nil {
			err = c.Options.Handle(violationDiagnostic(c.Platform, resolve(v.Field), v))
		}
	}})
	return err
}

// Validate runs the rules of a native message and its children at native JSON paths.
func Validate(platform string, check func(*validation.Check), options []Option) error {
	configured := Configure(options)
	var err error
	validation.Walk("$", check, func(path string, v validation.Violation) {
		if err == nil {
			err = configured.Handle(violationDiagnostic(platform, path, v))
		}
	})
	return err
}

func violationDiagnostic(platform, path string, v validation.Violation) Diagnostic {
	return Diagnostic{Platform: platform, Path: path, Code: v.Rule.ID, Message: v.Rule.Message, Severity: v.Rule.Severity}
}

// FieldPath resolves a native field such as "text" or "elements[1]" to the source
// path of the matching input slot, falling back to the node path.
func (n Node) FieldPath(field string) string {
	field, _, _ = strings.Cut(field, ".")
	name, index := field, 0
	if open := strings.IndexByte(field, '['); open >= 0 && strings.HasSuffix(field, "]") {
		parsed, err := strconv.Atoi(field[open+1 : len(field)-1])
		if err == nil {
			name, index = field[:open], parsed
		}
	}
	for _, input := range n.Inputs {
		if input.Slot != name {
			continue
		}
		if index > 0 {
			index--
			continue
		}
		if input.Part != nil {
			return input.Part.Path
		}
		return input.Child.Path
	}
	return n.Path
}
