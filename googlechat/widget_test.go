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
