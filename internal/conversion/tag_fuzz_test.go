package conversion

import (
	"reflect"
	"strings"
	"testing"
)

func FuzzParse(f *testing.F) {
	for _, raw := range []string{
		"", "-", "flatten;omitempty", "part;slot=text;style=bold;format=plain",
		"text;group=body;slot=text;style=bold,italic;format=plain;omitempty;optional",
		"header;style=bold;optional", "button;optional", "text;group=a=b",
		"typo;optional", "text;style=bold,bold", "text;slot=url", "text;",
	} {
		f.Add(raw)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		tag, err := Parse(testProfile, raw)
		if err != nil {
			return
		}
		canonical := tag.Role
		for _, option := range []struct{ key, value string }{
			{"group", tag.Group}, {"slot", tag.Slot},
			{"style", strings.Join(tag.Style, ",")}, {"format", tag.Format},
		} {
			if option.value != "" {
				canonical += ";" + option.key + "=" + option.value
			}
		}
		if tag.OmitEmpty {
			canonical += ";omitempty"
		}
		if tag.Optional {
			canonical += ";optional"
		}
		roundTrip, err := Parse(testProfile, canonical)
		if err != nil || !reflect.DeepEqual(tag, roundTrip) {
			t.Fatalf("canonical tag %q changed metadata: %#v -> %#v (%v)", canonical, tag, roundTrip, err)
		}
	})
}
