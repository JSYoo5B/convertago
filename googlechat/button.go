package googlechat

import "strings"

// ButtonList displays a collection of buttons.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#ButtonList
type ButtonList struct {
	// Buttons contains the buttons in display order.
	Buttons []Button `json:"buttons" validate:"dive"`
}

func (ButtonList) WidgetType() string { return "buttonList" }
func (b ButtonList) String() string {
	parts := make([]string, len(b.Buttons))
	for i, button := range b.Buttons {
		parts[i] = button.Text
	}
	return strings.Join(parts, " ")
}
func (ButtonList) widgetContent()       {}
func (ButtonList) columnWidgetContent() {}
func (ButtonList) nestedWidgetContent() {}

// Button displays text, an icon, or both, and performs a click action.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Button
type Button struct {
	// Text labels the button.
	Text string `json:"text,omitempty"`
	// Icon appears before Text when both are supplied.
	Icon *Icon `json:"icon,omitempty"`
	// Color sets RGB components from 0 to 1 and forces a filled button.
	Color *Color `json:"color,omitempty"`
	// OnClick is the required click behavior.
	OnClick OnClick `json:"onClick"`
	// Disabled prevents user interaction.
	Disabled bool `json:"disabled,omitempty"`
	// AltText describes the button's purpose for accessibility.
	AltText string `json:"altText,omitempty"`
	// Type selects the visual style. Omission uses outlined; Color overrides it with filled.
	Type ButtonType `json:"type,omitempty"`
}

// ButtonType selects the visual emphasis of a button.
type ButtonType string

const (
	ButtonTypeOutlined    ButtonType = "OUTLINED"
	ButtonTypeFilled      ButtonType = "FILLED"
	ButtonTypeFilledTonal ButtonType = "FILLED_TONAL"
	ButtonTypeBorderless  ButtonType = "BORDERLESS"
)

// Color expresses an RGB color with components in the range 0 to 1.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Color
type Color struct {
	// Red is the red component, from 0 to 1.
	Red float64 `json:"red"`
	// Green is the green component, from 0 to 1.
	Green float64 `json:"green"`
	// Blue is the blue component, from 0 to 1.
	Blue float64 `json:"blue"`
}
