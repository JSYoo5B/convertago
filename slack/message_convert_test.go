package slack

import (
	"strings"
	"testing"
)

func TestToMessageValidatesSlackConstraints(t *testing.T) {
	for _, source := range []any{
		struct {
			Text string `slack:"header"`
		}{strings.Repeat("가", 151)},
		struct {
			Text string `slack:"section"`
		}{""},
		struct {
			Text string `slack:"section"`
		}{strings.Repeat("a", 3001)},
		struct {
			URL string `slack:"image"`
		}{"https://example.com/image.png"},
		struct {
			Texts []string `slack:"section"`
		}{strings.Fields(strings.Repeat("a ", 51))},
		struct {
			Plain string `slack:"section;group=x"`
			Raw   string `slack:"section;group=x;format=mrkdwn"`
		}{"*literal*", "*bold*"},
	} {
		if _, err := ToMessage(source); err == nil {
			t.Fatalf("expected an error for %#v", source)
		}
	}
}

func TestSectionMarkupIsExplicit(t *testing.T) {
	source := struct {
		Plain string `slack:"section"`
		Raw   string `slack:"section;format=mrkdwn"`
	}{"*literal*", "*bold*"}
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := message.Blocks[0].(SectionBlock).Text.(PlainTextObject); !ok {
		t.Fatal("default text must stay literal")
	}
	if _, ok := message.Blocks[1].(SectionBlock).Text.(MrkdwnTextObject); !ok {
		t.Fatal("explicit mrkdwn must retain its syntax")
	}
}
