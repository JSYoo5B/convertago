package googlechat

import (
	"encoding/json"
	"fmt"
)

// DecoratedText displays text with optional labels, a leading icon, and a trailing control.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#DecoratedText
type DecoratedText struct {
	// StartIcon appears before the text.
	StartIcon *Icon `json:"startIcon,omitempty" validate:"omitempty"`
	// StartIconVerticalAlignment positions StartIcon vertically; omission centers it.
	StartIconVerticalAlignment VerticalAlignment `json:"startIconVerticalAlignment,omitempty" validate:"omitempty,oneof=TOP MIDDLE BOTTOM"`
	// TopLabel appears above Text and truncates when necessary.
	TopLabel string `json:"topLabel,omitempty"`
	// TopLabelText supplies a TextParagraph in place of the simple top label.
	TopLabelText *TextParagraph `json:"topLabelText,omitempty" validate:"omitempty"`
	// Text is the required primary text and supports simple formatting.
	Text string `json:"text"`
	// ContentText supplies a TextParagraph representation of the primary text.
	ContentText *TextParagraph `json:"contentText,omitempty" validate:"omitempty"`
	// WrapText wraps Text onto multiple lines instead of truncating it.
	WrapText bool `json:"wrapText,omitempty"`
	// BottomLabel appears below Text and wraps when necessary.
	BottomLabel string `json:"bottomLabel,omitempty"`
	// BottomLabelText supplies a TextParagraph in place of the simple bottom label.
	BottomLabelText *TextParagraph `json:"bottomLabelText,omitempty" validate:"omitempty"`
	// OnClick runs when the top or bottom label is clicked.
	OnClick *OnClick `json:"onClick,omitempty" validate:"omitempty"`
	// Button is a trailing button. Supply at most one trailing control.
	Button *Button `json:"button,omitempty" validate:"excluded_with=SwitchControl EndIcon,omitempty"`
	// SwitchControl is a trailing switch or checkbox.
	SwitchControl *SwitchControl `json:"switchControl,omitempty" validate:"excluded_with=Button EndIcon,omitempty"`
	// EndIcon is a trailing icon.
	EndIcon *Icon `json:"endIcon,omitempty" validate:"excluded_with=Button SwitchControl,omitempty"`
}

func (DecoratedText) WidgetType() string   { return "decoratedText" }
func (t DecoratedText) String() string     { return t.Text }
func (DecoratedText) widgetContent()       {}
func (DecoratedText) columnWidgetContent() {}
func (t DecoratedText) MarshalJSON() ([]byte, error) {
	count := 0
	if t.Button != nil {
		count++
	}
	if t.SwitchControl != nil {
		count++
	}
	if t.EndIcon != nil {
		count++
	}
	if count > 1 {
		return nil, fmt.Errorf("googlechat: DecoratedText allows only one trailing control")
	}
	type Embed DecoratedText
	return json.Marshal(Embed(t))
}

// SwitchControl displays a toggle switch or checkbox inside DecoratedText.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#SwitchControl
type SwitchControl struct {
	// Name identifies the control in form input events.
	Name string `json:"name" validate:"required"`
	// Value is the value returned with the form input event.
	Value string `json:"value,omitempty"`
	// Selected controls the initial checked state.
	Selected bool `json:"selected,omitempty"`
	// OnChangeAction runs when the selection changes.
	OnChangeAction *Action `json:"onChangeAction,omitempty" validate:"omitempty"`
	// ControlType selects a switch or checkbox. Omission uses a switch.
	ControlType SwitchControlType `json:"controlType,omitempty" validate:"omitempty,oneof=SWITCH CHECK_BOX"`
}

// SwitchControlType selects the control's appearance.
type SwitchControlType string

const (
	SwitchControlTypeSwitch   SwitchControlType = "SWITCH"
	SwitchControlTypeCheckBox SwitchControlType = "CHECK_BOX"
)
