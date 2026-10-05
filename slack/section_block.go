package slack

import (
	"encoding/json"
	"strings"
)

// SectionBlock displays text, compact fields, or text beside an accessory.
//
// Reference: https://docs.slack.dev/reference/block-kit/blocks/section-block/
type SectionBlock struct {
	// Text contains 1 to 3000 characters. It can be omitted when Fields is supplied.
	Text TextObject `json:"text,omitempty"`
	// BlockID identifies this block. Use a new ID for each message update, up to 255 characters.
	BlockID string `json:"block_id,omitempty"`
	// Fields contains up to 10 text objects, each up to 2000 characters, in two columns.
	Fields []TextObject `json:"fields,omitempty" validate:"dive"`
	// Accessory is an element compatible with sections, such as a button or image.
	Accessory Element `json:"accessory,omitempty"`
	// Expand displays the full text without requiring the reader to expand it.
	Expand bool `json:"expand,omitempty"`
}

func (b SectionBlock) Type() string { return "section" }
func (SectionBlock) block()         {}
func (b SectionBlock) String() string {
	var parts []string
	if b.Text != nil {
		parts = append(parts, b.Text.String())
	}
	for _, field := range b.Fields {
		if field != nil {
			parts = append(parts, field.String())
		}
	}
	return strings.Join(parts, "\n")
}
func (b SectionBlock) MarshalJSON() ([]byte, error) {
	type Embed SectionBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{b.Type(), Embed(b)})
}
