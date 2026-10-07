package conversion

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
)

// Tag is parsed source metadata. Style and Fixed may reference shared, read-only storage.
type Tag struct {
	Role      string
	Group     string
	Slot      string
	Style     []string
	Format    string
	OmitEmpty bool
	Optional  bool
	// Fixed supplies tag values for enum and boolean slots of the role, in declaration order.
	Fixed []Fixed
}

// Fixed is a slot value written in a tag, such as "color=red" or a bare "bold".
type Fixed struct {
	Slot  string
	Value string
}

// textStyles is the shared text style vocabulary of the style option.
var textStyles = []string{"bold", "italic", "strike", "code", "underline", "highlight", "client_highlight", "unlink"}

// reservedOptions name tag options that slots cannot use for fixed values.
var reservedOptions = []string{"group", "slot", "format", "omitempty", "optional"}

// valueError is an invalid fixed value, reported with the code of the slot's rule.
type valueError struct {
	code string
	err  error
}

func (e valueError) Error() string { return e.err.Error() }

// TagErrorCode returns the diagnostic code of a tag error.
func TagErrorCode(err error) string {
	var value valueError
	if errors.As(err, &value) && value.code != "" {
		return value.code
	}
	return "invalid_tag"
}

// Slot defines the inputs accepted by a platform builder.
type Slot struct {
	Repeated bool
	Required bool
	Children []string
	Scalar   bool
	Styles   []string
	Formats  []string
	// Values lists the values a tag may fix for an enum slot; Rule is the rule ID they enforce.
	Values []string
	Rule   string
	// Bool allows a tag to fix a boolean slot, as a bare name for true.
	Bool bool
}

// Fixable reports whether a tag may fix the slot's value.
func (s Slot) Fixable() bool { return s.Bool || len(s.Values) != 0 }

// Role describes a native builder or a recognized, unavailable feature.
type Role struct {
	DefaultSlot      string
	DefaultChildSlot string
	EmptyAllowed     bool
	NestedOnly       bool
	Slots            map[string]Slot
	Styles           []string
	FormatStyles     map[string][]string
	Formats          []string
	Unavailable      bool
}

// Profile is owned and registered by the corresponding messenger package.
type Profile struct {
	Platform string
	Roles    map[string]Role
}

var profiles sync.Map

// Register adds a messenger profile. Each platform registers once.
// It panics when slot names would make tag options ambiguous.
func Register(profile Profile) {
	for name, role := range profile.Roles {
		for slotName, slot := range role.Slots {
			if contains(reservedOptions, slotName) {
				panic("convertago: " + profile.Platform + " " + name + " slot " + slotName + " collides with a tag option")
			}
			for _, value := range slot.Values {
				if contains(textStyles, value) {
					panic("convertago: " + profile.Platform + " " + name + " slot " + slotName + " value " + value + " collides with a text style")
				}
			}
			if len(slot.Values) != 0 && slot.Rule == "" {
				panic("convertago: " + profile.Platform + " " + name + " slot " + slotName + " lists values without a rule")
			}
		}
	}
	if _, loaded := profiles.LoadOrStore(profile.Platform, profile); loaded {
		panic("convertago: duplicate platform " + profile.Platform)
	}
}

// Lookup returns the registered profile of a platform.
func Lookup(platform string) (Profile, error) {
	profile, ok := profiles.Load(platform)
	if !ok {
		return Profile{}, fmt.Errorf("convertago: unknown platform %q", platform)
	}
	return profile.(Profile), nil
}

// Parse uses exact, case-sensitive names. An absent, empty, or '-' tag is ignored.
// Options resolve in order: reserved options, then style as a text style list when every
// item is a text style, then a fixed value for an enum or boolean slot of the role.
func Parse(profile Profile, raw string) (Tag, error) {
	if raw == "" || raw == "-" {
		return Tag{}, nil
	}
	pieces := strings.Split(raw, ";")
	tag := Tag{Role: pieces[0]}
	if tag.Role == "" {
		return Tag{}, fmt.Errorf("tag requires a role")
	}
	seen := make(map[string]bool)
	for _, piece := range pieces[1:] {
		name, value, assigned := strings.Cut(piece, "=")
		if name == "" {
			return Tag{}, fmt.Errorf("tag option requires a name")
		}
		if seen[name] {
			return Tag{}, fmt.Errorf("duplicate option %q", name)
		}
		seen[name] = true
		switch name {
		case "group", "slot", "format":
			if !assigned || value == "" {
				return Tag{}, fmt.Errorf("option %q requires a value", name)
			}
			switch name {
			case "group":
				tag.Group = value
			case "slot":
				tag.Slot = value
			case "format":
				tag.Format = value
			}
		case "omitempty", "optional":
			if assigned {
				return Tag{}, fmt.Errorf("option %q is a flag", name)
			}
			if name == "omitempty" {
				tag.OmitEmpty = true
			} else {
				tag.Optional = true
			}
		default:
			if assigned && value == "" {
				return Tag{}, fmt.Errorf("option %q requires a value", name)
			}
			if name == "style" && assigned && allTextStyles(strings.Split(value, ",")) {
				tag.Style = strings.Split(value, ",")
				continue
			}
			if !assigned {
				value = "true"
			}
			tag.Fixed = append(tag.Fixed, Fixed{Slot: name, Value: value})
		}
	}
	if _, err := CheckTag(profile, tag); err != nil {
		return Tag{}, err
	}
	return tag, nil
}

// CheckTag returns recognized features that the current builder cannot express.
// Unknown names always produce an error, even on an optional field.
func CheckTag(profile Profile, tag Tag) (string, error) {
	if tag.Role == "" {
		return "", nil
	}
	if tag.Role == "flatten" {
		if tag.Group != "" || tag.Slot != "" || len(tag.Style) != 0 || tag.Format != "" || len(tag.Fixed) != 0 {
			return "", fmt.Errorf("flatten only accepts omitempty and optional")
		}
		return "", nil
	}
	role, exists := profile.Roles[tag.Role]
	if !exists && tag.Role != "part" {
		return "", fmt.Errorf("unknown role %q", tag.Role)
	}
	var unavailable []string
	if role.Unavailable {
		unavailable = append(unavailable, "role "+tag.Role)
	}
	seen := make(map[string]bool)
	styles := role.Styles
	formats := role.Formats
	if tag.Role != "part" && tag.Slot != "" && tag.Slot != role.DefaultSlot {
		if slot, own := role.Slots[tag.Slot]; own {
			styles = slot.Styles
			formats = slot.Formats
		}
	}
	if allowed, exists := role.FormatStyles[tag.Format]; exists && (tag.Slot == "" || tag.Slot == role.DefaultSlot) {
		styles = allowed
	}
	if err := checkFixed(role, tag); err != nil {
		return "", err
	}
	for _, style := range tag.Style {
		if !contains(textStyles, style) {
			return "", fmt.Errorf("unknown style %q", style)
		}
		if seen[style] {
			return "", fmt.Errorf("duplicate style %q", style)
		}
		seen[style] = true
		if tag.Role != "part" && !contains(styles, style) {
			unavailable = append(unavailable, "style "+style)
		}
	}
	if tag.Format != "" {
		if !contains([]string{"plain", "html", "markdown", "mrkdwn"}, tag.Format) {
			return "", fmt.Errorf("unknown format %q", tag.Format)
		}
		if tag.Role != "part" && !contains(formats, tag.Format) {
			unavailable = append(unavailable, "format "+tag.Format)
		}
	}
	if tag.Role != "part" && !role.Unavailable && tag.Slot != "" {
		_, known := role.Slots[tag.Slot]
		for _, parent := range profile.Roles {
			if contains(parent.Slots[tag.Slot].Children, tag.Role) {
				known = true
			}
		}
		if !known {
			return "", fmt.Errorf("unknown slot %q for %s", tag.Slot, tag.Role)
		}
	}
	return strings.Join(unavailable, ", "), nil
}

func allTextStyles(values []string) bool {
	for _, value := range values {
		if !contains(textStyles, value) {
			return false
		}
	}
	return true
}

// checkFixed validates fixed slot values against the role's fixable slots.
func checkFixed(role Role, tag Tag) error {
	for _, fixed := range tag.Fixed {
		if tag.Role == "part" {
			return fmt.Errorf("part does not accept fixed slot values")
		}
		if role.Unavailable {
			continue
		}
		slot, exists := role.Slots[fixed.Slot]
		if !exists || !slot.Fixable() {
			if fixed.Slot == "style" {
				return fmt.Errorf("unknown style %q", fixed.Value)
			}
			return fmt.Errorf("unknown option %q", fixed.Slot)
		}
		if slot.Bool {
			if fixed.Value != "true" && fixed.Value != "false" {
				return fmt.Errorf("option %q requires true or false", fixed.Slot)
			}
			continue
		}
		if !slices.Contains(slot.Values, fixed.Value) {
			message := fmt.Errorf("unknown value %q for %s; allowed: %s", fixed.Value, fixed.Slot, strings.Join(slot.Values, ", "))
			if fixed.Slot == "style" {
				message = fmt.Errorf("unknown value %q for style; allowed: %s, or text styles", fixed.Value, strings.Join(slot.Values, ", "))
			}
			return valueError{slot.Rule, message}
		}
	}
	return nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
