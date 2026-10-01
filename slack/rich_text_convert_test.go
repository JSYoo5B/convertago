package slack

import "testing"

func TestRichTextMixedHierarchy(t *testing.T) {
	source := struct {
		Body struct {
			Before string `slack:"part"`
			Quote  string `slack:"rich_text_quote"`
			After  string `slack:"part;style=highlight"`
		} `slack:"rich_text"`
	}{}
	source.Body.Before = "before"
	source.Body.Quote = "quote"
	source.Body.After = "after"
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	block := message.Blocks[0].(RichTextBlock)
	if len(block.Elements) != 3 || block.Elements[1].(RichTextQuote).Elements[0].String() != "quote" || !block.Elements[2].(RichTextSection).Elements[0].(TextInline).Style.Highlight {
		t.Fatalf("block=%#v", block)
	}
}

func TestPreformattedRejectsUnsupportedInlineEvenIfNil(t *testing.T) {
	source := struct {
		Body struct {
			Code struct {
				User *string `slack:"user;optional"`
			} `slack:"rich_text_preformatted"`
		} `slack:"rich_text"`
	}{}
	if _, err := ToMessage(source); err == nil {
		t.Fatal("preformatted accepted user")
	}
}

func TestOptionalRichTextStyle(t *testing.T) {
	source := struct {
		Body struct {
			Link string `slack:"link;style=code;optional"`
			Text string `slack:"part"`
		} `slack:"rich_text"`
	}{}
	source.Body.Link = "https://example.com"
	source.Body.Text = "keep"
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	if len(message.Blocks[0].(RichTextBlock).Elements[0].(RichTextSection).Elements) != 1 {
		t.Fatal("unsupported code on link must skip")
	}
}

func TestExplicitRichLinkLabelStyle(t *testing.T) {
	type link struct {
		URL  string `slack:"part;slot=url"`
		Text string `slack:"part;slot=text;style=italic"`
	}
	source := struct {
		Body struct {
			Link link `slack:"link"`
		} `slack:"rich_text"`
	}{}
	source.Body.Link = link{"https://example.com", "label"}
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	inline := message.Blocks[0].(RichTextBlock).Elements[0].(RichTextSection).Elements[0].(LinkInline)
	if inline.Style == nil || !inline.Style.Italic || inline.Text != "label" {
		t.Fatalf("inline=%#v", inline)
	}
}
