package googlechat

import "strings"

// ChipList displays chips horizontally, wrapping or scrolling when space runs out.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#ChipList
type ChipList struct {
	// Layout selects wrapping or horizontal scrolling. Omission uses wrapping.
	Layout ChipListLayout `json:"layout,omitempty" validate:"omitempty,oneof=WRAPPED HORIZONTAL_SCROLLABLE"`
	// Chips contains the chips in display order.
	Chips []Chip `json:"chips" validate:"min=1,dive"`
}

func (ChipList) WidgetType() string   { return "chipList" }
func (ChipList) widgetContent()       {}
func (ChipList) columnWidgetContent() {}
func (c ChipList) String() string {
	parts := make([]string, len(c.Chips))
	for i, chip := range c.Chips {
		parts[i] = chip.Label
	}
	return strings.Join(parts, " ")
}

// Chip displays an icon, text, or both, with an optional click action.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Chip
type Chip struct {
	// Icon appears before Label when both are supplied.
	Icon *Icon `json:"icon,omitempty" validate:"required_without=Label,omitempty"`
	// Label is the displayed text.
	Label string `json:"label,omitempty" validate:"required_without=Icon"`
	// OnClick runs when the chip is clicked.
	OnClick *OnClick `json:"onClick,omitempty" validate:"omitempty"`
	// Disabled prevents the chip from responding to user actions.
	Disabled bool `json:"disabled,omitempty"`
	// AltText describes the chip's purpose for accessibility.
	AltText string `json:"altText,omitempty"`
}

// ChipListLayout selects wrapping or horizontal scrolling.
type ChipListLayout string

const (
	ChipListLayoutWrapped              ChipListLayout = "WRAPPED"
	ChipListLayoutHorizontalScrollable ChipListLayout = "HORIZONTAL_SCROLLABLE"
)
