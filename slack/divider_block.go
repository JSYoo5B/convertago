package slack

import "encoding/json"

// DividerBlock separates content with a horizontal rule.
//
// Reference: https://docs.slack.dev/reference/block-kit/blocks/divider-block/
type DividerBlock struct {
	// BlockID identifies the block, up to 255 characters. Replace it when updating a message.
	BlockID string `json:"block_id,omitempty" validate:"max=255"`
}

func (b DividerBlock) Type() string { return "divider" }
func (DividerBlock) String() string { return "" }
func (DividerBlock) block()         {}
func (b DividerBlock) MarshalJSON() ([]byte, error) {
	type Embed DividerBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{b.Type(), Embed(b)})
}
