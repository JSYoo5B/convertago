package googlechat_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
)

func TestWidgetRejectsMissingContent(t *testing.T) {
	var paragraph *googlechat.TextParagraph
	for _, content := range []googlechat.WidgetContent{nil, paragraph} {
		if _, err := json.Marshal(googlechat.Widget{Content: content}); err == nil {
			t.Fatal("expected missing widget content to fail")
		}
	}
}

func TestOnClickUnion(t *testing.T) {
	link := &googlechat.OpenLink{URL: "https://example.com"}
	action := &googlechat.Action{Function: "approve"}
	menu := &googlechat.OverflowMenu{Items: []googlechat.OverflowMenuItem{{Text: "Open", OnClick: googlechat.OnClick{OpenLink: link}}}}
	for _, click := range []googlechat.OnClick{{OpenLink: link}, {Action: action}, {OverflowMenu: menu}} {
		if _, err := json.Marshal(click); err != nil {
			t.Fatal(err)
		}
	}
	for _, click := range []googlechat.OnClick{{}, {OpenLink: link, Action: action}, {Action: action, OverflowMenu: menu}, {OpenLink: link, OverflowMenu: menu}} {
		if _, err := json.Marshal(click); err == nil {
			t.Fatal("expected ambiguous click action to fail")
		}
	}
}

func TestIconUnion(t *testing.T) {
	for _, icon := range []googlechat.Icon{{KnownIcon: "EMAIL"}, {IconURL: "https://example.com/icon.png"}, {MaterialIcon: &googlechat.MaterialIcon{Name: "mail"}}} {
		if _, err := json.Marshal(icon); err != nil {
			t.Fatal(err)
		}
	}
	for _, icon := range []googlechat.Icon{{}, {KnownIcon: "EMAIL", IconURL: "https://example.com/icon.png"}, {KnownIcon: "EMAIL", MaterialIcon: &googlechat.MaterialIcon{Name: "mail"}}} {
		if _, err := json.Marshal(icon); err == nil {
			t.Fatal("expected ambiguous icon source to fail")
		}
	}
}

func TestDecoratedTextControlUnion(t *testing.T) {
	text := googlechat.DecoratedText{
		Text:          "Recipient",
		EndIcon:       &googlechat.Icon{KnownIcon: "EMAIL"},
		SwitchControl: &googlechat.SwitchControl{Name: "notify"},
	}
	if _, err := json.Marshal(googlechat.Widget{Content: text}); err == nil {
		t.Fatal("expected multiple trailing controls to fail")
	}
}

func TestTextParagraphIsRawInsideDecoratedText(t *testing.T) {
	widget := googlechat.Widget{Content: googlechat.DecoratedText{
		Text:        "Passed",
		ContentText: &googlechat.TextParagraph{Text: "**Passed**", TextSyntax: googlechat.TextSyntaxMarkdown},
	}}
	payload, err := json.Marshal(widget)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(payload), `"contentText":{"text":"**Passed**","textSyntax":"MARKDOWN"}`) {
		t.Fatalf("TextParagraph must be a raw object inside contentText: %s", payload)
	}
}
