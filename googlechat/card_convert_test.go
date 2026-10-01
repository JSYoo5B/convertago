package googlechat

import "testing"

func TestGoogleExplicitCardIDs(t *testing.T) {
	type card struct {
		Text string `googlechat:"textParagraph"`
	}
	type wrapper struct {
		ID   string `googlechat:"part;slot=cardId"`
		Card card   `googlechat:"card"`
	}
	source := struct {
		Cards []wrapper `googlechat:"cardWithId"`
	}{[]wrapper{{"one", card{"one"}}, {"two", card{"two"}}}}
	if _, err := ToMessage(source); err != nil {
		t.Fatal(err)
	}
	source.Cards[1].ID = "one"
	if _, err := ToMessage(source); err == nil {
		t.Fatal("duplicate cardId accepted")
	}
	source.Cards[1].ID = ""
	if _, err := ToMessage(source); err == nil {
		t.Fatal("missing cardId in multiple cards accepted")
	}
}

func TestGoogleCardWidgetCountAcrossSections(t *testing.T) {
	type section struct {
		Text []string `googlechat:"textParagraph"`
	}
	source := struct {
		Sections []section `googlechat:"section"`
	}{[]section{{make([]string, 50)}, {make([]string, 50)}}}
	if _, err := ToMessage(source); err != nil {
		t.Fatal(err)
	}
	source.Sections[1].Text = append(source.Sections[1].Text, "extra")
	if _, err := ToMessage(source); err == nil {
		t.Fatal("over 100 widgets across sections accepted")
	}
}

func TestCollapsePropertiesRequireCollapsible(t *testing.T) {
	source := struct {
		Section struct {
			Text        string `googlechat:"textParagraph"`
			Count       int    `googlechat:"part;slot=uncollapsibleWidgetsCount"`
			Collapsible bool   `googlechat:"part;slot=collapsible"`
		} `googlechat:"section"`
	}{}
	source.Section.Text = "one"
	source.Section.Count = 1
	if _, err := ToMessage(source); err == nil {
		t.Fatal("collapse count without collapsible accepted")
	}
	source.Section.Collapsible = true
	if _, err := ToMessage(source); err != nil {
		t.Fatal(err)
	}
	source.Section.Count = 2
	if _, err := ToMessage(source); err == nil {
		t.Fatal("collapse count exceeds widgets")
	}
}

func TestCustomCollapseControls(t *testing.T) {
	type button struct {
		Text   string `googlechat:"part"`
		Action string `googlechat:"action"`
	}
	type control struct {
		Expand   button `googlechat:"button;slot=expandButton"`
		Collapse button `googlechat:"button;slot=collapseButton"`
	}
	source := struct {
		Section struct {
			Text        string  `googlechat:"textParagraph"`
			Collapsible bool    `googlechat:"part;slot=collapsible"`
			Control     control `googlechat:"collapseControl"`
		} `googlechat:"section"`
	}{}
	source.Section.Text = "detail"
	source.Section.Collapsible = true
	source.Section.Control = control{button{"More", "expand"}, button{"Less", "collapse"}}
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	native := message.CardsV2[0].Card.Sections[0].CollapseControl
	if native.ExpandButton.Text != "More" || native.CollapseButton.OnClick.Action.Function != "collapse" {
		t.Fatalf("control=%#v", native)
	}
}
