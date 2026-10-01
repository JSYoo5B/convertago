package googlechat

import (
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/internal/conversion"
)

func TestParagraphKeepsPlainTextLiteral(t *testing.T) {
	source := struct {
		Plain string `googlechat:"textParagraph;style=bold"`
		HTML  string `googlechat:"textParagraph;format=html"`
		MD    string `googlechat:"textParagraph;format=markdown"`
	}{"<b>&\n*literal*", "<b>bold</b>", "**bold**"}
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	widgets := message.CardsV2[0].Card.Sections[0].Widgets
	if got := widgets[0].Content.(TextParagraph).Text; got != "<b>&lt;b&gt;&amp;<br>*literal*</b>" {
		t.Fatalf("escaped text = %q", got)
	}
	if got := widgets[1].Content.(TextParagraph).Text; got != "<b>bold</b>" {
		t.Fatalf("HTML text = %q", got)
	}
	if got := widgets[2].Content.(TextParagraph); got.Text != "**bold**" || got.TextSyntax != TextSyntaxMarkdown {
		t.Fatalf("Markdown paragraph = %#v", got)
	}
}

func TestToMessageValidatesGoogleChatConstraints(t *testing.T) {
	for _, source := range []any{
		struct {
			Text  string `googlechat:"textParagraph"`
			Title string `googlechat:"header"`
		}{"body", "late"},
		struct {
			URL string `googlechat:"image"`
		}{"http://example.com/image.png"},
		struct {
			Text []string `googlechat:"textParagraph"`
		}{strings.Fields(strings.Repeat("a ", 101))},
		struct {
			Text string `googlechat:"textParagraph"`
		}{strings.Repeat("a", 32*1024)},
		struct {
			A string `googlechat:"textParagraph;group=x"`
			B string `googlechat:"textParagraph;group=x;format=markdown"`
		}{"literal", "**bold**"},
	} {
		if _, err := ToMessage(source); err == nil {
			t.Fatalf("expected an error for %#v", source)
		}
	}
}

func TestOptionalMarkdownUnderline(t *testing.T) {
	source := struct {
		Text string `googlechat:"textParagraph;format=markdown;style=underline;optional"`
	}{"text"}
	var diagnostics []conversion.Diagnostic
	message, err := ToMessage(source, conversion.WithDiagnostics(func(d conversion.Diagnostic) { diagnostics = append(diagnostics, d) }))
	if err != nil || len(message.CardsV2) != 0 || len(diagnostics) != 1 {
		t.Fatalf("message = %#v, diagnostics = %#v, error = %v", message, diagnostics, err)
	}
	if _, err := ToMessage(source, conversion.WithStrict()); err == nil {
		t.Fatal("strict must reject the unsupported formatting")
	}
}
