package conversion

import (
	"fmt"
	"strings"
	"sync"
)

// Tag is the parsed metadata emitted by the source generator.
type Tag struct {
	Role      string
	Group     string
	Slot      string
	Style     []string
	Format    string
	OmitEmpty bool
	Optional  bool
}

// Slot defines the inputs accepted by a platform builder.
type Slot struct {
	Repeated bool
	Required bool
}

// Role describes a native builder or a recognized, unavailable feature.
type Role struct {
	DefaultSlot string
	Slots       map[string]Slot
	Styles      []string
	Formats     []string
	Unavailable bool
}

// Profile is owned and registered by the corresponding messenger package.
type Profile struct {
	Platform string
	Roles    map[string]Role
}

var profiles sync.Map

func Register(profile Profile) {
	if _, loaded := profiles.LoadOrStore(profile.Platform, profile); loaded {
		panic("convertago: duplicate platform " + profile.Platform)
	}
}

func Lookup(platform string) (Profile, error) {
	profile, ok := profiles.Load(platform)
	if !ok {
		return Profile{}, fmt.Errorf("convertago: unknown platform %q", platform)
	}
	return profile.(Profile), nil
}

// Parse uses exact, case-sensitive names. An absent, empty, or '-' tag is ignored.
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
		if seen[name] {
			return Tag{}, fmt.Errorf("duplicate option %q", name)
		}
		seen[name] = true
		switch name {
		case "group", "slot", "style", "format":
			if !assigned || value == "" {
				return Tag{}, fmt.Errorf("option %q requires a value", name)
			}
			switch name {
			case "group":
				tag.Group = value
			case "slot":
				tag.Slot = value
			case "style":
				tag.Style = strings.Split(value, ",")
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
			return Tag{}, fmt.Errorf("unknown option %q", name)
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
		if tag.Group != "" || tag.Slot != "" || len(tag.Style) != 0 || tag.Format != "" {
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
	for _, style := range tag.Style {
		if !contains([]string{"bold", "italic", "strike", "code", "underline"}, style) {
			return "", fmt.Errorf("unknown style %q", style)
		}
		if seen[style] {
			return "", fmt.Errorf("duplicate style %q", style)
		}
		seen[style] = true
		if tag.Role != "part" && !contains(role.Styles, style) {
			unavailable = append(unavailable, "style "+style)
		}
	}
	if tag.Format != "" {
		if !contains([]string{"plain", "html", "markdown", "mrkdwn"}, tag.Format) {
			return "", fmt.Errorf("unknown format %q", tag.Format)
		}
		if tag.Role != "part" && !contains(role.Formats, tag.Format) {
			unavailable = append(unavailable, "format "+tag.Format)
		}
	}
	if tag.Role != "part" && !role.Unavailable && tag.Slot != "" {
		if _, exists := role.Slots[tag.Slot]; !exists {
			return "", fmt.Errorf("unknown slot %q for %s", tag.Slot, tag.Role)
		}
	}
	return strings.Join(unavailable, ", "), nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
