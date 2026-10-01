package slack

import "encoding/json"

// TextObject contains plain text or Slack mrkdwn for a block or element.
//
// Reference: https://docs.slack.dev/reference/block-kit/composition-objects/text-object/
type TextObject interface {
	Type() string
	String() string
	MarshalJSON() ([]byte, error)
	textObject()
}

// PlainTextObject contains unformatted text.
//
// Reference: https://docs.slack.dev/reference/block-kit/composition-objects/text-object/
type PlainTextObject struct {
	// Text contains 1 to 3000 characters, subject to the containing object's limit.
	Text string `json:"text"`
	// Emoji controls conversion of emoji to colon notation. Nil uses Slack's default.
	Emoji *bool `json:"emoji,omitempty"`
}

func (t PlainTextObject) Type() string   { return "plain_text" }
func (t PlainTextObject) String() string { return t.Text }
func (PlainTextObject) textObject()      {}
func (PlainTextObject) contextElement()  {}
func (t PlainTextObject) MarshalJSON() ([]byte, error) {
	type Embed PlainTextObject
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{t.Type(), Embed(t)})
}

// MrkdwnTextObject contains text formatted with Slack's mrkdwn syntax.
//
// Reference: https://docs.slack.dev/reference/block-kit/composition-objects/text-object/
type MrkdwnTextObject struct {
	// Text contains 1 to 3000 characters, subject to the containing object's limit.
	Text string `json:"text"`
	// Verbatim disables automatic links and mentions while retaining mrkdwn parsing.
	Verbatim bool `json:"verbatim,omitempty"`
}

func (t MrkdwnTextObject) Type() string   { return "mrkdwn" }
func (t MrkdwnTextObject) String() string { return t.Text }
func (MrkdwnTextObject) textObject()      {}
func (MrkdwnTextObject) contextElement()  {}
func (t MrkdwnTextObject) MarshalJSON() ([]byte, error) {
	type Embed MrkdwnTextObject
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{t.Type(), Embed(t)})
}
