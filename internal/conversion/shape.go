package conversion

import "fmt"

// ValidateShape is used by both the reflection compiler and the source generator.
func ValidateShape(profile Profile, shape *Shape) error {
	return validateShape(profile, shape, "", "$", make(map[shapeVisit]bool))
}

type shapeVisit struct {
	shape   *Shape
	builder string
}

func validateShape(profile Profile, shape *Shape, builder, path string, active map[shapeVisit]bool) error {
	key := shapeVisit{shape, builder}
	if active[key] {
		return nil
	}
	active[key] = true
	defer delete(active, key)
	if shape.Err != nil {
		return Error(profile.Platform, path, "invalid_source", shape.Err.Error())
	}
	if shape.Kind == "pointer" || shape.Kind == "list" {
		if shape.Kind == "list" {
			path += "[]"
		}
		return validateShape(profile, shape.Elem, builder, path, active)
	}
	if shape.Kind != "object" {
		return nil
	}
	groups := make(map[string]string)
	for _, field := range shape.Fields {
		fieldPath := path + "." + field.Name
		if field.Err != nil {
			return Error(profile.Platform, fieldPath, "invalid_tag", field.Err.Error())
		}
		if err := validateContext(profile, field.Tag, builder); err != nil {
			return Error(profile.Platform, fieldPath, "invalid_tag", err.Error())
		}
		if err := checkGroup(groups, field.Tag); err != nil {
			return Error(profile.Platform, fieldPath, "group_conflict", err.Error())
		}
		if field.Tag.Role == "flatten" {
			kind := field.Value
			for kind.Kind == "pointer" || kind.Kind == "list" {
				kind = kind.Elem
			}
			if kind.Kind != "object" && kind.Kind != "dynamic" {
				return Error(profile.Platform, fieldPath, "invalid_source", "flatten requires a struct or a list of structs")
			}
		}
		next := builder
		if field.Tag.Role != "flatten" && field.Tag.Role != "part" {
			next = field.Tag.Role
		}
		if err := validateShape(profile, field.Value, next, fieldPath, active); err != nil {
			return err
		}
	}
	return nil
}

func validateContext(profile Profile, tag Tag, builder string) error {
	if _, err := CheckTag(profile, tag); err != nil {
		return err
	}
	if builder == "" && tag.Role == "part" {
		return fmt.Errorf("part requires a containing builder")
	}
	if builder != "" && tag.Role != "part" && tag.Role != "flatten" {
		return fmt.Errorf("%s inputs must use part or flatten", builder)
	}
	if tag.Role == "part" && tag.Slot != "" {
		if role := profile.Roles[builder]; !role.Unavailable {
			if _, exists := role.Slots[tag.Slot]; !exists {
				return fmt.Errorf("unknown slot %q for %s", tag.Slot, builder)
			}
		}
	}
	return nil
}

func checkGroup(groups map[string]string, tag Tag) error {
	if tag.Group == "" {
		return nil
	}
	if role, exists := groups[tag.Group]; exists && role != tag.Role {
		return fmt.Errorf("group %q mixes %s and %s", tag.Group, role, tag.Role)
	}
	groups[tag.Group] = tag.Role
	return nil
}
