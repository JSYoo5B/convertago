package slack

import (
	"encoding/json"
	"strings"
)

// Element is an accessory or interactive component contained in a layout block.
//
// Reference: https://docs.slack.dev/reference/block-kit/block-elements/
type Element interface {
	Type() string
	String() string
	MarshalJSON() ([]byte, error)
	element()
}

// ActionsBlock groups interactive elements such as buttons and menus.
//
// Reference: https://docs.slack.dev/reference/block-kit/blocks/actions-block/
type ActionsBlock struct {
	// Elements contains up to 25 interactive elements. ImageElement belongs in a section or context.
	Elements []ActionElement `json:"elements"`
	// BlockID identifies the block, up to 255 characters. Replace it when updating a message.
	BlockID string `json:"block_id,omitempty"`
}

// ActionElement is an interactive element that can be placed in an ActionsBlock.
//
// Reference: https://docs.slack.dev/reference/block-kit/blocks/actions-block/
type ActionElement interface {
	Element
	actionElement()
}

func (b ActionsBlock) Type() string { return "actions" }
func (ActionsBlock) block()         {}
func (b ActionsBlock) String() string {
	var parts []string
	for _, element := range b.Elements {
		if element != nil {
			parts = append(parts, element.String())
		}
	}
	return strings.Join(parts, " ")
}
func (b ActionsBlock) MarshalJSON() ([]byte, error) {
	type Embed ActionsBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{b.Type(), Embed(b)})
}
