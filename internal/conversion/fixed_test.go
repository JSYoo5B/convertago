package conversion

import (
	"errors"
	"reflect"
	"testing"
)

var fixedProfile = Profile{Platform: "fixed", Roles: map[string]Role{
	"banner": {DefaultSlot: "text", Slots: map[string]Slot{
		"text":  {Repeated: true, Required: true},
		"style": {Values: []string{"blue", "red"}, Rule: "banner.style.value"},
		"tone":  {Values: []string{"calm", "loud"}, Rule: "banner.tone.value", Required: true},
		"bold":  {Bool: true},
		"wide":  {Bool: true},
		"note":  {},
	}, Styles: []string{"bold", "italic"}},
}}

func init() { Register(fixedProfile) }

func TestParseFixedValues(t *testing.T) {
	tag, err := Parse(fixedProfile, "banner;style=blue;bold;wide=false;tone=calm")
	if err != nil {
		t.Fatal(err)
	}
	want := []Fixed{{"style", "blue"}, {"bold", "true"}, {"wide", "false"}, {"tone", "calm"}}
	if !reflect.DeepEqual(tag.Fixed, want) || tag.Style != nil {
		t.Fatalf("Fixed = %#v, Style = %#v", tag.Fixed, tag.Style)
	}

	tag, err = Parse(fixedProfile, "banner;style=bold,italic;tone=loud")
	if err != nil || !reflect.DeepEqual(tag.Style, []string{"bold", "italic"}) || len(tag.Fixed) != 1 {
		t.Fatalf("text styles = %#v, fixed = %#v, %v", tag.Style, tag.Fixed, err)
	}
}

func TestParseRejectsInvalidFixedValues(t *testing.T) {
	for raw, code := range map[string]string{
		"banner;style=green":          "banner.style.value",
		"banner;style=bold,blue":      "banner.style.value",
		"banner;tone":                 "banner.tone.value",
		"banner;bold=yes":             "invalid_tag",
		"banner;note=x":               "invalid_tag",
		"banner;colour=red":           "invalid_tag",
		"banner;style=":               "invalid_tag",
		"banner;bold;bold":            "invalid_tag",
		"banner;style=blue;style=red": "invalid_tag",
		"part;bold":                   "invalid_tag",
		"flatten;bold":                "invalid_tag",
	} {
		t.Run(raw, func(t *testing.T) {
			_, err := Parse(fixedProfile, raw)
			if err == nil {
				t.Fatal("expected an error")
			}
			if got := TagErrorCode(err); got != code {
				t.Fatalf("code = %q, want %q (%v)", got, code, err)
			}
		})
	}
}

func TestFixedValuesAndFieldOverrides(t *testing.T) {
	type banner struct {
		Text  string `fixed:"part"`
		Style string `fixed:"part;slot=style;omitempty"`
	}
	source := struct {
		Fixed    string `fixed:"banner;style=blue;bold;tone=calm"`
		Override banner `fixed:"banner;style=blue;tone=loud"`
	}{Fixed: "a", Override: banner{"b", "red"}}
	nodes, err := Prepare(source, "fixed", nil)
	if err != nil {
		t.Fatal(err)
	}
	first, second := nodes[0], nodes[1]
	if first.Text("style") != "blue" || first.Text("bold") != "true" || !first.Has("tone") || first.Has("wide") {
		t.Fatalf("fixed node = %#v", first)
	}
	if first.FieldPath("style") != "$.Fixed" {
		t.Fatalf("fixed path = %q", first.FieldPath("style"))
	}
	if second.Text("style") != "red" || second.FieldPath("style") != "$.Override.Style" {
		t.Fatalf("override = %q at %q", second.Text("style"), second.FieldPath("style"))
	}

	source.Override.Style = ""
	nodes, err = Prepare(source, "fixed", nil)
	if err != nil || nodes[1].Text("style") != "blue" {
		t.Fatalf("omitted override = %#v, %v", nodes, err)
	}
}

func TestFixedValuesSatisfyRequiredSlots(t *testing.T) {
	missing := struct {
		Text string `fixed:"banner"`
	}{"a"}
	_, err := Prepare(missing, "fixed", nil)
	var diagnostic Diagnostic
	if !errors.As(err, &diagnostic) || diagnostic.Code != "missing_input" {
		t.Fatalf("missing tone = %v", err)
	}
	fixed := struct {
		Text string `fixed:"banner;tone=calm"`
	}{"a"}
	if _, err := Prepare(fixed, "fixed", nil); err != nil {
		t.Fatal(err)
	}
}

func TestGroupMembersShareFixedValues(t *testing.T) {
	same := struct {
		A string `fixed:"banner;group=g;tone=calm"`
		B string `fixed:"banner;group=g;tone=calm"`
	}{"a", "b"}
	nodes, err := Prepare(same, "fixed", nil)
	if err != nil || len(nodes) != 1 || nodes[0].Text("text") != "ab" {
		t.Fatalf("group = %#v, %v", nodes, err)
	}
	mixed := struct {
		A string `fixed:"banner;group=g;tone=calm"`
		B string `fixed:"banner;group=g;tone=loud"`
	}{"a", "b"}
	_, err = Prepare(mixed, "fixed", nil)
	var diagnostic Diagnostic
	if !errors.As(err, &diagnostic) || diagnostic.Code != "group_conflict" {
		t.Fatalf("mixed group = %v", err)
	}
}

func TestInvalidFixedValueDiagnostic(t *testing.T) {
	source := struct {
		Text string `fixed:"banner;style=green;tone=calm"`
	}{"a"}
	_, err := Prepare(source, "fixed", nil)
	var diagnostic Diagnostic
	if !errors.As(err, &diagnostic) || diagnostic.Code != "banner.style.value" || diagnostic.Path != "$.Text" {
		t.Fatalf("diagnostic = %#v", err)
	}
}
