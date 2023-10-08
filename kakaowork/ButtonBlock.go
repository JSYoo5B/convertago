package kakaowork

import (
	"encoding/json"
)

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/

type ButtonBlock struct {
	Text   string
	Style  ButtonStyle
	Action ButtonAction
}

type ButtonStyle string

const (
	ButtonStyleEmpty   = ButtonStyle("")
	ButtonStyleDefault = ButtonStyle("default")
	ButtonStyleGray    = ButtonStyleDefault
	ButtonStylePrimary = ButtonStyle("primary")
	ButtonStyleBlue    = ButtonStylePrimary
	ButtonStyleDanger  = ButtonStyle("danger")
	ButtonStyleRed     = ButtonStyleDanger
)

var buttonStyleConstants = map[ButtonStyle]bool{
	ButtonStyleEmpty: true,
	ButtonStyleGray:  true,
	ButtonStyleBlue:  true,
	ButtonStyleRed:   true,
}

func (b ButtonBlock) Type() string {
	return "button"
}

func (b ButtonBlock) String() string {
	return b.Text + ": " + b.Action.String()
}

func (b ButtonBlock) MarshalJSON() ([]byte, error) {
	type _ButtonBlock struct {
		Type   string       `json:"type"`
		Text   string       `json:"text"`
		Style  ButtonStyle  `json:"style,omitempty"`
		Action ButtonAction `json:"action"`
	}

	if _, exists := buttonStyleConstants[b.Style]; !exists {
		b.Style = ButtonStyleDefault
	}
	return json.Marshal(_ButtonBlock{
		Type:   b.Type(),
		Text:   b.Text,
		Style:  b.Style,
		Action: b.Action,
	})
}

type ButtonAction interface {
	MessageBubbleBlock
	ActionType() string
}

type OpenSystemBrowserAction struct {
	Name  string
	Value string
}

func (o OpenSystemBrowserAction) Type() string {
	return "action"
}

func (o OpenSystemBrowserAction) String() string {
	return o.Value
}

func (o OpenSystemBrowserAction) MarshalJSON() ([]byte, error) {
	type _OpenSystemBrowserAction struct {
		Type  string `json:"type"`
		Name  string `json:"name,omitempty"`
		Value string `json:"value"`
	}

	return json.Marshal(_OpenSystemBrowserAction{
		Type:  o.ActionType(),
		Name:  o.Name,
		Value: o.Value,
	})
}

func (o OpenSystemBrowserAction) ActionType() string {
	return "open_system_browser"
}

type OpenInAppBrowserAction struct {
	Name       string
	Value      string
	Standalone bool
	Width      int
	Height     int
}

func (o OpenInAppBrowserAction) Type() string {
	return "action"
}

func (o OpenInAppBrowserAction) String() string {
	return o.Value
}

func (o OpenInAppBrowserAction) MarshalJSON() ([]byte, error) {
	type _OpenInAppBrowserAction struct {
		Type       string `json:"type"`
		Name       string `json:"name,omitempty"`
		Value      string `json:"value"`
		Standalone bool   `json:"standalone,omitempty"`
		Width      int    `json:"width,omitempty"`
		Height     int    `json:"height,omitempty"`
	}

	if !o.Standalone {
		o.Width, o.Height = 0, 0
	}

	return json.Marshal(_OpenInAppBrowserAction{
		Type:       o.ActionType(),
		Name:       o.Name,
		Value:      o.Value,
		Standalone: o.Standalone,
		Width:      o.Width,
		Height:     o.Height,
	})
}

func (o OpenInAppBrowserAction) ActionType() string {
	return "open_inapp_browser"
}

type OpenExternalAppAction struct {
	Name  string
	Value string
}

func (o OpenExternalAppAction) Type() string {
	return "action"
}

func (o OpenExternalAppAction) String() string {
	return o.Value
}

func (o OpenExternalAppAction) MarshalJSON() ([]byte, error) {
	type _OpenExternalAppAction struct {
		Type  string `json:"type"`
		Name  string `json:"name,omitempty"`
		Value string `json:"value"`
	}

	return json.Marshal(_OpenExternalAppAction{
		Type:  o.ActionType(),
		Name:  o.Name,
		Value: o.Value,
	})
}

func (o OpenExternalAppAction) ActionType() string {
	return "open_external_app"
}

type SubmitAction struct {
	Name  string
	Value string
}

func (s SubmitAction) Type() string {
	return "action"
}

func (s SubmitAction) String() string {
	return s.Value
}

func (s SubmitAction) MarshalJSON() ([]byte, error) {
	type _SubmitAction struct {
		Type  string `json:"type"`
		Name  string `json:"name,omitempty"`
		Value string `json:"value"`
	}

	return json.Marshal(_SubmitAction{
		Type:  s.ActionType(),
		Name:  s.Name,
		Value: s.Value,
	})
}

func (s SubmitAction) ActionType() string {
	return "submit_action"
}

type CallModalAction struct {
	Name  string
	Value string
}

func (c CallModalAction) Type() string {
	return "action"
}

func (c CallModalAction) String() string {
	return c.Value
}

func (c CallModalAction) MarshalJSON() ([]byte, error) {
	type _CallModalAction struct {
		Type  string `json:"type"`
		Name  string `json:"name,omitempty"`
		Value string `json:"value"`
	}

	return json.Marshal(_CallModalAction{
		Type:  c.ActionType(),
		Name:  c.Name,
		Value: c.Value,
	})
}

func (c CallModalAction) ActionType() string {
	return "call_modal"
}

type ExclusiveAction struct {
	Default ButtonAction
	Pc      ButtonAction
	Mobile  ButtonAction
	Windows ButtonAction
	MacOs   ButtonAction
	Android ButtonAction
	Ios     ButtonAction
}

func (e ExclusiveAction) Type() string {
	return "action"
}

func (e ExclusiveAction) String() string {
	return e.Default.String()
}

func (e ExclusiveAction) MarshalJSON() ([]byte, error) {
	type _ExclusiveAction struct {
		Type    string       `json:"type"`
		Default ButtonAction `json:"default"`
		Pc      ButtonAction `json:"pc,omitempty"`
		Mobile  ButtonAction `json:"mobile,omitempty"`
		Windows ButtonAction `json:"windows,omitempty"`
		MacOs   ButtonAction `json:"macos,omitempty"`
		Android ButtonAction `json:"android,omitempty"`
		Ios     ButtonAction `json:"ios,omitempty"`
	}

	return json.Marshal(_ExclusiveAction{
		Type:    e.ActionType(),
		Default: e.Default,
		Pc:      e.Pc,
		Mobile:  e.Mobile,
		Windows: e.Windows,
		MacOs:   e.MacOs,
		Android: e.Android,
		Ios:     e.Ios,
	})
}

func (e ExclusiveAction) ActionType() string {
	return "exclusive"
}
