package conversion

import "testing"

var testProfile = Profile{Platform: "test", Roles: map[string]Role{
	"text":   {DefaultSlot: "text", Slots: map[string]Slot{"text": {Repeated: true, Required: true}}, Styles: []string{"bold", "italic"}, Formats: []string{"plain"}},
	"button": {Unavailable: true},
}}

func TestParse(t *testing.T) {
	tag, err := Parse(testProfile, "text;group=body;slot=text;style=bold,italic;format=plain;omitempty;optional")
	if err != nil || tag.Role != "text" || tag.Group != "body" || tag.Slot != "text" || len(tag.Style) != 2 || tag.Format != "plain" || !tag.OmitEmpty || !tag.Optional {
		t.Fatalf("Parse = %#v, %v", tag, err)
	}
	for _, raw := range []string{"", "-"} {
		tag, err := Parse(testProfile, raw)
		if err != nil || tag.Role != "" {
			t.Fatalf("Parse(%q) = %#v, %v", raw, tag, err)
		}
	}
}

func TestParseRejectsMistakesEvenWhenOptional(t *testing.T) {
	for _, raw := range []string{
		"Text", ";optional", "typo;optional", "text;", "text;unknown=1;optional",
		"text;group", "text;group=", "text;group=a;group=b",
		"text;omitempty=true", "text;style=bold,bold", "text;style=blod;optional",
		"text;slot=url", "text;format=typo;optional", "flatten;group=x",
	} {
		t.Run(raw, func(t *testing.T) {
			if _, err := Parse(testProfile, raw); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func TestRecognizedUnavailableFeatures(t *testing.T) {
	for _, raw := range []string{"button;optional", "text;style=underline;optional", "text;format=markdown;optional"} {
		tag, err := Parse(testProfile, raw)
		if err != nil {
			t.Fatal(err)
		}
		unavailable, err := CheckTag(testProfile, tag)
		if err != nil || unavailable == "" {
			t.Fatalf("CheckTag = %q, %v", unavailable, err)
		}
	}
}
