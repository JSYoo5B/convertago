package slack

import (
	"encoding/json"
	"strings"
)

// RichTextElement is a section, list, quote, or preformatted region in a RichTextBlock.
//
// Reference: https://docs.slack.dev/reference/block-kit/blocks/rich-text-block/
type RichTextElement interface {
	Type() string
	String() string
	MarshalJSON() ([]byte, error)
	richTextElement()
}

// RichTextBlock represents formatted content from Slack's message composer.
// Its elements can contain nested sections and styled inline content.
//
// Reference: https://docs.slack.dev/reference/block-kit/blocks/rich-text-block/
type RichTextBlock struct {
	// Elements contains rich text sections, lists, preformatted regions, or quotes.
	Elements []RichTextElement `json:"elements"`
	// BlockID identifies the block, up to 255 characters. Replace it when updating a message.
	BlockID string `json:"block_id,omitempty"`
}

func (b RichTextBlock) Type() string { return "rich_text" }
func (RichTextBlock) block()         {}
func (b RichTextBlock) String() string {
	var parts []string
	for _, element := range b.Elements {
		if element != nil {
			parts = append(parts, element.String())
		}
	}
	return strings.Join(parts, "\n")
}
func (b RichTextBlock) MarshalJSON() ([]byte, error) {
	type Embed RichTextBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{b.Type(), Embed(b)})
}

// RichTextSection combines rich text inline elements in display order.
//
// Reference: https://docs.slack.dev/reference/block-kit/block-elements/rich-text-section-element/
type RichTextSection struct {
	// Elements contains inline text, links, emoji, or mentions.
	Elements []RichTextInline `json:"elements"`
}

func (e RichTextSection) Type() string   { return "rich_text_section" }
func (e RichTextSection) String() string { return inlineString(e.Elements) }
func (RichTextSection) richTextElement() {}
func (e RichTextSection) MarshalJSON() ([]byte, error) {
	type Embed RichTextSection
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{e.Type(), Embed(e)})
}

// RichTextList displays rich text sections as a bulleted or numbered list.
//
// Reference: https://docs.slack.dev/reference/block-kit/block-elements/rich-text-list-element/
type RichTextList struct {
	// Style selects bullets or an ordered list.
	Style RichTextListStyle `json:"style"`
	// Elements contains the sections displayed as list items.
	Elements []RichTextSection `json:"elements"`
	// Indent sets the sub-list indentation level.
	Indent int `json:"indent,omitempty"`
	// Offset shifts the first number of an ordered list; 4 starts at 5.
	Offset int `json:"offset,omitempty"`
	// Border enables or disables the border. Nil leaves it unspecified.
	Border *int `json:"border,omitempty"`
}

// RichTextListStyle selects the list marker format.
type RichTextListStyle string

const (
	RichTextListStyleBullet  RichTextListStyle = "bullet"
	RichTextListStyleOrdered RichTextListStyle = "ordered"
)

func (e RichTextList) Type() string   { return "rich_text_list" }
func (RichTextList) richTextElement() {}
func (e RichTextList) String() string {
	parts := make([]string, len(e.Elements))
	for i, section := range e.Elements {
		parts[i] = section.String()
	}
	return strings.Join(parts, "\n")
}
func (e RichTextList) MarshalJSON() ([]byte, error) {
	type Embed RichTextList
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{e.Type(), Embed(e)})
}

// RichTextPreformatted displays text and links in a preformatted region.
//
// Reference: https://docs.slack.dev/reference/block-kit/block-elements/rich-text-preformatted-element/
type RichTextPreformatted struct {
	// Elements contains text or link elements.
	Elements []PreformattedInline `json:"elements"`
	// Border enables or disables the border. Nil leaves it unspecified.
	Border *int `json:"border,omitempty"`
	// Language selects a language for code syntax highlighting.
	Language string `json:"language,omitempty"`
}

// PreformattedInline is text or a link accepted by RichTextPreformatted.
//
// Reference: https://docs.slack.dev/reference/block-kit/block-elements/rich-text-preformatted-element/
type PreformattedInline interface {
	RichTextInline
	preformattedInline()
}

func (e RichTextPreformatted) Type() string   { return "rich_text_preformatted" }
func (RichTextPreformatted) richTextElement() {}
func (e RichTextPreformatted) String() string {
	var text strings.Builder
	for _, inline := range e.Elements {
		if inline != nil {
			text.WriteString(inline.String())
		}
	}
	return text.String()
}
func (e RichTextPreformatted) MarshalJSON() ([]byte, error) {
	type Embed RichTextPreformatted
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{e.Type(), Embed(e)})
}

// RichTextQuote displays rich text inline content as a quotation.
//
// Reference: https://docs.slack.dev/reference/block-kit/block-elements/rich-text-quote-element/
type RichTextQuote struct {
	// Elements contains the quotation's inline content.
	Elements []RichTextInline `json:"elements"`
	// Border enables or disables the border. Nil leaves it unspecified.
	Border *int `json:"border,omitempty"`
}

func (e RichTextQuote) Type() string   { return "rich_text_quote" }
func (e RichTextQuote) String() string { return inlineString(e.Elements) }
func (RichTextQuote) richTextElement() {}
func (e RichTextQuote) MarshalJSON() ([]byte, error) {
	type Embed RichTextQuote
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{e.Type(), Embed(e)})
}

func inlineString(elements []RichTextInline) string {
	var text strings.Builder
	for _, element := range elements {
		if element != nil {
			text.WriteString(element.String())
		}
	}
	return text.String()
}
