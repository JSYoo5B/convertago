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
			return Error(profile.Platform, fieldPath, TagErrorCode(field.Err), field.Err.Error())
		}
		if err := validateContext(profile, field.Tag, builder); err != nil {
			return Error(profile.Platform, fieldPath, TagErrorCode(err), err.Error())
		}
		if err := checkGroup(groups, field.Tag, builder, profile); err != nil {
			return Error(profile.Platform, fieldPath, "group_conflict", err.Error())
		}
		if err := validateNativeShape(profile, field.Tag, builder, field.Value, fieldPath); err != nil {
			return err
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
	if tag.Role == "" || tag.Role == "flatten" {
		return nil
	}
	if builder == "" {
		role := profile.Roles[tag.Role]
		if role.NestedOnly {
			return fmt.Errorf("%s requires a containing builder", tag.Role)
		}
		if tag.Slot != "" && !role.Unavailable {
			if _, exists := role.Slots[tag.Slot]; !exists {
				return fmt.Errorf("unknown slot %q for %s", tag.Slot, tag.Role)
			}
		}
		return nil
	}
	if profile.Roles[builder].Unavailable {
		return nil
	}
	_, err := inputSlot(profile, tag, builder)
	if err != nil {
		return err
	}
	return nil
}

func inputSlot(profile Profile, tag Tag, builder string) (string, error) {
	role := profile.Roles[builder]
	name := tag.Slot
	if tag.Role == "part" {
		if name == "" {
			name = role.DefaultSlot
		}
		slot, exists := role.Slots[name]
		if !exists {
			return "", fmt.Errorf("unknown slot %q for %s", name, builder)
		}
		if len(slot.Children) != 0 && !slot.Scalar {
			return "", fmt.Errorf("slot %q for %s requires a child builder", name, builder)
		}
		return name, nil
	}
	if name == "" && contains(role.Slots[role.DefaultChildSlot].Children, tag.Role) {
		name = role.DefaultChildSlot
	}
	if name == "" {
		for candidate, slot := range role.Slots {
			if contains(slot.Children, tag.Role) {
				if name != "" {
					return "", fmt.Errorf("%s requires slot in %s", tag.Role, builder)
				}
				name = candidate
			}
		}
	}
	slot, exists := role.Slots[name]
	if !exists || !contains(slot.Children, tag.Role) {
		return "", fmt.Errorf("%s does not accept %s in slot %q", builder, tag.Role, name)
	}
	return name, nil
}

func checkGroup(groups map[string]string, tag Tag, builder string, profile Profile) error {
	if tag.Group == "" {
		return nil
	}
	signature := tag.Role
	for _, fixed := range tag.Fixed {
		signature += ";" + fixed.Slot + "=" + fixed.Value
	}
	if builder != "" && tag.Role != "part" {
		slot, _ := inputSlot(profile, tag, builder)
		signature += ":" + slot
	}
	if role, exists := groups[tag.Group]; exists && role != signature {
		return fmt.Errorf("group %q mixes %s and %s; members must share their role and fixed values", tag.Group, role, signature)
	}
	groups[tag.Group] = signature
	return nil
}

func validateNativeShape(profile Profile, tag Tag, builder string, shape *Shape, path string) error {
	if tag.Role == "" || tag.Role == "part" || tag.Role == "flatten" || profile.Roles[tag.Role].Unavailable {
		return nil
	}
	for shape.Kind == "pointer" || shape.Kind == "list" {
		shape = shape.Elem
	}
	if shape.Kind == "scalar" {
		slot := profile.Roles[tag.Role].DefaultSlot
		if builder == "" && tag.Slot != "" {
			slot = tag.Slot
		}
		if slot == "" {
			return Error(profile.Platform, path, "invalid_source", tag.Role+" requires a tagged struct")
		}
	}
	return nil
}
