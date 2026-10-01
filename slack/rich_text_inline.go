package slack

import "encoding/json"

// RichTextInline is an inline component of a rich text section or quotation.
//
// Reference: https://docs.slack.dev/reference/block-kit/block-elements/rich-text-section-element/
type RichTextInline interface {
	Type() string
	String() string
	MarshalJSON() ([]byte, error)
	richTextInline()
}

// RichTextStyle applies optional formatting flags to rich text.
// Code is supported by TextInline; the other flags are also supported by links and users.
//
// Reference: https://docs.slack.dev/reference/block-kit/block-elements/text-element/
type RichTextStyle struct {
	// Bold renders the content in bold.
	Bold bool `json:"bold,omitempty"`
	// Code applies code formatting to TextInline.
	Code bool `json:"code,omitempty"`
	// Italic renders the content in italics.
	Italic bool `json:"italic,omitempty"`
	// Strike applies strikethrough formatting.
	Strike bool `json:"strike,omitempty"`
	// Highlight requests highlighted content.
	Highlight bool `json:"highlight,omitempty"`
	// ClientHighlight records client highlighting.
	ClientHighlight bool `json:"client_highlight,omitempty"`
	// Underline underlines the content.
	Underline bool `json:"underline,omitempty"`
	// Unlink removes link formatting.
	Unlink bool `json:"unlink,omitempty"`
}

// TextInline displays text with optional rich text formatting.
//
// Reference: https://docs.slack.dev/reference/block-kit/block-elements/text-element/
type TextInline struct {
	// Text is the content displayed to the user.
	Text string `json:"text"`
	// Style applies optional formatting to the text.
	Style *RichTextStyle `json:"style,omitempty"`
}

func (e TextInline) Type() string      { return "text" }
func (e TextInline) String() string    { return e.Text }
func (TextInline) richTextInline()     {}
func (TextInline) preformattedInline() {}
func (e TextInline) MarshalJSON() ([]byte, error) {
	type Embed TextInline
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{e.Type(), Embed(e)})
}

// LinkInline displays a link with an optional label and formatting.
//
// Reference: https://docs.slack.dev/reference/block-kit/block-elements/link-element/
type LinkInline struct {
	// URL is the link's destination.
	URL string `json:"url"`
	// Text replaces the displayed URL when supplied.
	Text string `json:"text,omitempty"`
	// Unsafe identifies a potentially unsafe link.
	Unsafe bool `json:"unsafe,omitempty"`
	// FromLLM indicates that an LLM generated the link.
	FromLLM bool `json:"from_llm,omitempty"`
	// IsSlackURL indicates a Slack destination.
	IsSlackURL bool `json:"is_slack_url,omitempty"`
	// Truncated indicates that the displayed link has been shortened.
	Truncated bool `json:"truncated,omitempty"`
	// Style applies formatting other than Code.
	Style *RichTextStyle `json:"style,omitempty"`
}

func (e LinkInline) Type() string { return "link" }
func (e LinkInline) String() string {
	if e.Text != "" {
		return e.Text
	}
	return e.URL
}
func (LinkInline) richTextInline()     {}
func (LinkInline) preformattedInline() {}
func (e LinkInline) MarshalJSON() ([]byte, error) {
	type Embed LinkInline
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{e.Type(), Embed(e)})
}

// UserInline mentions a Slack user in rich text.
//
// Reference: https://docs.slack.dev/reference/block-kit/block-elements/user-element/
type UserInline struct {
	// UserID identifies the user being mentioned.
	UserID string `json:"user_id"`
	// Style applies formatting other than Code.
	Style *RichTextStyle `json:"style,omitempty"`
	// FromLLM indicates that an LLM generated the mention.
	FromLLM bool `json:"from_llm,omitempty"`
}

func (e UserInline) Type() string   { return "user" }
func (e UserInline) String() string { return "<@" + e.UserID + ">" }
func (UserInline) richTextInline()  {}
func (e UserInline) MarshalJSON() ([]byte, error) {
	type Embed UserInline
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{e.Type(), Embed(e)})
}

// EmojiInline displays an emoji identified by its Slack name.
//
// Reference: https://docs.slack.dev/reference/block-kit/block-elements/emoji-element/
type EmojiInline struct {
	// Name identifies the emoji, including an optional skin-tone suffix.
	Name string `json:"name"`
	// Unicode is the emoji's Unicode code point when applicable.
	Unicode string `json:"unicode,omitempty"`
}

func (e EmojiInline) Type() string   { return "emoji" }
func (e EmojiInline) String() string { return ":" + e.Name + ":" }
func (EmojiInline) richTextInline()  {}
func (e EmojiInline) MarshalJSON() ([]byte, error) {
	type Embed EmojiInline
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{e.Type(), Embed(e)})
}
