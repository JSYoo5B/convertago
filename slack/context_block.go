package slack

import (
	"github.com/JSYoo5B/convertago/internal/validation"

	"encoding/json"
	"strings"
)

// ContextElement is a text object or image element displayed in a ContextBlock.
//
// Reference: https://docs.slack.dev/reference/block-kit/blocks/context-block/
type ContextElement interface {
	Type() string
	String() string
	MarshalJSON() ([]byte, error)
	contextElement()
	check(*validation.Check)
}

// ContextBlock displays contextual information with text and images.
//
// Reference: https://docs.slack.dev/reference/block-kit/blocks/context-block/
type ContextBlock struct {
	// Elements contains up to 10 text objects or image elements in display order.
	Elements []ContextElement `json:"elements" validate:"dive"`
	// BlockID identifies the block, up to 255 characters. Replace it when updating a message.
	BlockID string `json:"block_id,omitempty"`
}

func (b ContextBlock) Type() string { return "context" }
func (ContextBlock) block()         {}
func (b ContextBlock) String() string {
	var parts []string
	for _, element := range b.Elements {
		if element != nil {
			parts = append(parts, element.String())
		}
	}
	return strings.Join(parts, " ")
}
func (b ContextBlock) MarshalJSON() ([]byte, error) {
	type Embed ContextBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{b.Type(), Embed(b)})
}
