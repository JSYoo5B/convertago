package slack

import "encoding/json"

// HeaderBlock displays a heading in a larger, bold font.
//
// Reference: https://docs.slack.dev/reference/block-kit/blocks/header-block/
type HeaderBlock struct {
	// Text is the heading's plain text, up to 150 characters.
	Text PlainTextObject `json:"text"`
	// BlockID identifies the block, up to 255 characters. Replace it when updating a message.
	BlockID string `json:"block_id,omitempty" validate:"max=255"`
	// Level selects a heading level from 1 to 4.
	Level int `json:"level,omitempty" validate:"omitempty,min=1,max=4"`
}

func (b HeaderBlock) Type() string   { return "header" }
func (b HeaderBlock) String() string { return b.Text.String() }
func (HeaderBlock) block()           {}
func (b HeaderBlock) MarshalJSON() ([]byte, error) {
	type Embed HeaderBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{b.Type(), Embed(b)})
}
