package googlechat

import (
	"math"
	"testing"
)

func TestGoogleActionParameters(t *testing.T) {
	type param struct {
		Key   string `googlechat:"part;slot=key"`
		Value string `googlechat:"part;slot=value"`
	}
	type action struct {
		Function   string  `googlechat:"part"`
		Parameters []param `googlechat:"actionParameter"`
	}
	type button struct {
		Text   string `googlechat:"part"`
		Action action `googlechat:"action"`
	}
	source := struct {
		Buttons struct {
			Button button `googlechat:"button"`
		} `googlechat:"buttonList"`
	}{}
	source.Buttons.Button = button{"Send", action{"send", []param{{"id", "42"}}}}
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	native := message.CardsV2[0].Card.Sections[0].Widgets[0].Content.(ButtonList).Buttons[0]
	if native.OnClick.Action.Parameters[0].Value != "42" {
		t.Fatal("parameter lost")
	}
	source.Buttons.Button.Action.Parameters = append(source.Buttons.Button.Action.Parameters, param{"id", "99"})
	if _, err := ToMessage(source); err == nil {
		t.Fatal("duplicate parameter key accepted")
	}
}

func TestGoogleColorAndIconValidation(t *testing.T) {
	type color struct {
		Red float64 `googlechat:"part;slot=red"`
	}
	type button struct {
		Text  string `googlechat:"part"`
		Link  string `googlechat:"openLink"`
		Color color  `googlechat:"color"`
	}
	source := struct {
		Buttons struct {
			Button button `googlechat:"button"`
		} `googlechat:"buttonList"`
	}{}
	source.Buttons.Button = button{"Open", "https://example.com", color{1}}
	if _, err := ToMessage(source); err != nil {
		t.Fatal(err)
	}
	for _, value := range []float64{1.1, -0.1, math.NaN(), math.Inf(1)} {
		source.Buttons.Button.Color.Red = value
		if _, err := ToMessage(source); err == nil {
			t.Fatalf("accepted %v color", value)
		}
	}
}

func TestOnClickAndTrailingControlConflicts(t *testing.T) {
	type click struct {
		Action string `googlechat:"action"`
		URL    string `googlechat:"openLink"`
	}
	source := struct {
		Image struct {
			URL   string `googlechat:"part;slot=url"`
			Click click  `googlechat:"onClick"`
		} `googlechat:"image"`
	}{}
	source.Image.URL = "https://example.com/image.png"
	source.Image.Click = click{"action", "https://example.com"}
	if _, err := ToMessage(source); err == nil {
		t.Fatal("ambiguous onClick accepted")
	}
	type icon struct {
		Name string `googlechat:"part;slot=knownIcon"`
	}
	type toggle struct {
		Name string `googlechat:"part"`
	}
	decorated := struct {
		Text struct {
			Text   string `googlechat:"part"`
			End    icon   `googlechat:"icon;slot=endIcon"`
			Switch toggle `googlechat:"switchControl"`
		} `googlechat:"decoratedText"`
	}{}
	decorated.Text.Text = "text"
	decorated.Text.End.Name = "STAR"
	decorated.Text.Switch.Name = "selected"
	if _, err := ToMessage(decorated); err == nil {
		t.Fatal("multiple trailing controls accepted")
	}
}

func TestMaterialIconValues(t *testing.T) {
	type material struct {
		Name   string `googlechat:"part"`
		Weight int    `googlechat:"part;slot=weight"`
		Grade  int    `googlechat:"part;slot=grade"`
	}
	type icon struct {
		Material material `googlechat:"materialIcon"`
	}
	source := struct {
		Text struct {
			Text      string `googlechat:"part"`
			Icon      icon   `googlechat:"icon;slot=startIcon"`
			Alignment string `googlechat:"part;slot=startIconVerticalAlignment"`
		} `googlechat:"decoratedText"`
	}{}
	source.Text.Text = "Status"
	source.Text.Icon.Material = material{"check", 400, 200}
	source.Text.Alignment = "MIDDLE"
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	native := message.CardsV2[0].Card.Sections[0].Widgets[0].Content.(DecoratedText)
	if native.StartIcon.MaterialIcon.Grade != 200 || native.StartIconVerticalAlignment != VerticalAlignmentMiddle {
		t.Fatalf("native=%#v", native)
	}
	source.Text.Icon.Material.Weight = 250
	if _, err := ToMessage(source); err == nil {
		t.Fatal("unsupported material weight accepted")
	}
	source.Text.Icon.Material = material{"check", 400, 10}
	if _, err := ToMessage(source); err == nil {
		t.Fatal("unsupported material grade accepted")
	}
}

func TestOverflowMenuConversion(t *testing.T) {
	type item struct {
		Text     string `googlechat:"part"`
		URL      string `googlechat:"openLink"`
		Disabled bool   `googlechat:"part;slot=disabled"`
	}
	type menu struct {
		Items []item `googlechat:"overflowMenuItem"`
	}
	type button struct {
		Text string `googlechat:"part"`
		Menu menu   `googlechat:"overflowMenu"`
	}
	source := struct {
		Buttons struct {
			Button button `googlechat:"button"`
		} `googlechat:"buttonList"`
	}{}
	source.Buttons.Button = button{"More", menu{[]item{{"Docs", "https://example.com", false}, {"Archive", "https://example.com/archive", true}}}}
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	native := message.CardsV2[0].Card.Sections[0].Widgets[0].Content.(ButtonList).Buttons[0].OnClick.OverflowMenu
	if len(native.Items) != 2 || !native.Items[1].Disabled || native.Items[0].OnClick.OpenLink.URL != "https://example.com" {
		t.Fatalf("menu=%#v", native)
	}
}

func TestDecoratedParagraphAndSwitchAction(t *testing.T) {
	type paragraph struct {
		Text     string `googlechat:"part"`
		MaxLines int    `googlechat:"part;slot=maxLines"`
	}
	type toggle struct {
		Name     string `googlechat:"part"`
		Selected bool   `googlechat:"part;slot=selected"`
		Action   string `googlechat:"action"`
	}
	source := struct {
		Text struct {
			Text   string    `googlechat:"part"`
			Top    paragraph `googlechat:"textParagraph;slot=topLabelText"`
			Switch toggle    `googlechat:"switchControl"`
		} `googlechat:"decoratedText"`
	}{}
	source.Text.Text = "Status"
	source.Text.Top = paragraph{"<label>", 0}
	source.Text.Switch = toggle{"selected", true, "change"}
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	native := message.CardsV2[0].Card.Sections[0].Widgets[0].Content.(DecoratedText)
	if native.TopLabelText.Text != "&lt;label&gt;" || !native.SwitchControl.Selected || native.SwitchControl.OnChangeAction.Function != "change" {
		t.Fatalf("native=%#v", native)
	}
}
