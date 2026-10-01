package conversion

import (
	"errors"
	"testing"
)

func init() {
	Register(Profile{Platform: "tree", Roles: map[string]Role{
		"row":     {DefaultChildSlot: "items", Slots: map[string]Slot{"items": {Children: []string{"label", "divider"}, Repeated: true, Required: true}, "accessory": {Children: []string{"link"}}}},
		"label":   {DefaultSlot: "text", DefaultChildSlot: "text", Slots: map[string]Slot{"text": {Scalar: true, Children: []string{"link"}, Repeated: true, Required: true}}},
		"link":    {NestedOnly: true, DefaultSlot: "text", Slots: map[string]Slot{"text": {Required: true}, "url": {Required: true}}},
		"divider": {EmptyAllowed: true},
	}})
}

func TestNestedBuilderOrderAndEmptyNodes(t *testing.T) {
	type link struct {
		Text string `tree:"part"`
		URL  string `tree:"part;slot=url"`
	}
	type label struct {
		Before string `tree:"part"`
		Link   link   `tree:"link"`
		After  string `tree:"part"`
	}
	source := struct {
		Row struct {
			Label   label    `tree:"label"`
			Divider struct{} `tree:"divider"`
		} `tree:"row"`
	}{}
	source.Row.Label = label{"before ", link{"link", "https://example.com"}, " after"}
	nodes, err := Prepare(source, "tree", nil)
	if err != nil {
		t.Fatal(err)
	}
	children := nodes[0].Children("items")
	if len(children) != 2 || len(children[0].Inputs) != 3 || children[1].Role != "divider" {
		t.Fatalf("nodes = %#v", nodes)
	}
	if children[0].Inputs[0].Part.Text != "before " || children[0].Inputs[1].Child.Text("text") != "link" || children[0].Inputs[2].Part.Text != " after" {
		t.Fatal("source order lost")
	}
}

func TestNestedSlotsValidateBeforeOmission(t *testing.T) {
	source := struct {
		Row *struct {
			Invalid string `tree:"row;optional;omitempty"`
		} `tree:"row;omitempty"`
	}{}
	_, err := Prepare(source, "tree", nil)
	var diagnostic Diagnostic
	if !errors.As(err, &diagnostic) || diagnostic.Path != "$.Row.Invalid" || diagnostic.Code != "invalid_tag" {
		t.Fatalf("error=%v", err)
	}
}

func TestNestedListsAndGroupSlots(t *testing.T) {
	source := struct {
		Row struct {
			Labels  []string `tree:"label"`
			Divider struct{} `tree:"divider"`
		} `tree:"row"`
	}{}
	source.Row.Labels = []string{"one", "two"}
	nodes, err := Prepare(source, "tree", nil)
	if err != nil {
		t.Fatal(err)
	}
	children := nodes[0].Children("items")
	if len(children) != 3 || children[0].Text("text") != "one" || children[1].Text("text") != "two" {
		t.Fatalf("children=%#v", children)
	}
}

func TestInheritedStylesExcludeMetadata(t *testing.T) {

	source := struct {
		Label struct {
			Text    string `styledtree:"part"`
			Enabled bool   `styledtree:"part;slot=enabled"`
		} `styledtree:"label;style=bold"`
	}{}
	source.Label.Text = "label"
	source.Label.Enabled = true
	nodes, err := Prepare(source, "styledtree", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes[0].Parts[0].Style) != 1 || len(nodes[0].Parts[1].Style) != 0 {
		t.Fatalf("parts=%#v", nodes[0].Parts)
	}
}

func init() {
	Register(Profile{Platform: "styledtree", Roles: map[string]Role{"label": {DefaultSlot: "text", Slots: map[string]Slot{"text": {Required: true}, "enabled": {}}, Styles: []string{"bold"}}}})
}

func TestUniqueChildSlotOverridesUnrelatedDefault(t *testing.T) {
	type link struct {
		Text string `tree:"part"`
		URL  string `tree:"part;slot=url"`
	}
	source := struct {
		Row struct {
			Label string `tree:"label"`
			Link  link   `tree:"link"`
		} `tree:"row"`
	}{}
	source.Row.Label = "label"
	source.Row.Link = link{"link", "https://example.com"}
	nodes, err := Prepare(source, "tree", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes[0].Children("accessory")) != 1 {
		t.Fatal("default items slot hid unique accessory slot")
	}
}

func TestContainerScalarMismatchIsNotHiddenByNil(t *testing.T) {
	source := struct {
		Row *int `tree:"row;optional;omitempty"`
	}{}
	_, err := Prepare(source, "tree", nil)
	var diagnostic Diagnostic
	if !errors.As(err, &diagnostic) || diagnostic.Code != "invalid_source" || diagnostic.Path != "$.Row" {
		t.Fatalf("error=%v", err)
	}
}
