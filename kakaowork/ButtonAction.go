package kakaowork

import "encoding/json"

type ButtonAction interface {
	MessageBubbleBlock
	ActionType() string
}

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/#open_system_browser

type OpenSystemBrowserAction struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value"`
}

func (o OpenSystemBrowserAction) Type() string {
	return "action"
}

func (o OpenSystemBrowserAction) String() string {
	return o.Value
}

func (o OpenSystemBrowserAction) MarshalJSON() ([]byte, error) {
	type Embed OpenSystemBrowserAction
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  o.ActionType(),
		Embed: (Embed)(o),
	})
}

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/#open_inapp_browser

func (o OpenSystemBrowserAction) ActionType() string {
	return "open_system_browser"
}

type OpenInAppBrowserAction struct {
	Name       string `json:"name,omitempty"`
	Value      string `json:"value"`
	Standalone bool   `json:"standalone,omitempty"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
}

func (o OpenInAppBrowserAction) Type() string {
	return "action"
}

func (o OpenInAppBrowserAction) String() string {
	return o.Value
}

func (o OpenInAppBrowserAction) MarshalJSON() ([]byte, error) {
	if !o.Standalone {
		o.Width, o.Height = 0, 0
	}

	type Embed OpenInAppBrowserAction
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  o.ActionType(),
		Embed: (Embed)(o),
	})
}

func (o OpenInAppBrowserAction) ActionType() string {
	return "open_inapp_browser"
}

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/#open_external_app

type OpenExternalAppAction struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value"`
}

func (o OpenExternalAppAction) Type() string {
	return "action"
}

func (o OpenExternalAppAction) String() string {
	return o.Value
}

func (o OpenExternalAppAction) MarshalJSON() ([]byte, error) {
	type Embed OpenExternalAppAction
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  o.ActionType(),
		Embed: (Embed)(o),
	})
}

func (o OpenExternalAppAction) ActionType() string {
	return "open_external_app"
}

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/#submit_action

type SubmitAction struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value"`
}

func (s SubmitAction) Type() string {
	return "action"
}

func (s SubmitAction) String() string {
	return s.Value
}

func (s SubmitAction) MarshalJSON() ([]byte, error) {
	type Embed SubmitAction
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  s.ActionType(),
		Embed: (Embed)(s),
	})
}

func (s SubmitAction) ActionType() string {
	return "submit_action"
}

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/#call_modal

type CallModalAction struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value"`
}

func (c CallModalAction) Type() string {
	return "action"
}

func (c CallModalAction) String() string {
	return c.Value
}

func (c CallModalAction) MarshalJSON() ([]byte, error) {
	type Embed CallModalAction
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  c.ActionType(),
		Embed: (Embed)(c),
	})
}

func (c CallModalAction) ActionType() string {
	return "call_modal"
}

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/#exclusive-action

type ExclusiveAction struct {
	Default ButtonAction `json:"default"`
	Pc      ButtonAction `json:"pc,omitempty"`
	Mobile  ButtonAction `json:"mobile,omitempty"`
	Windows ButtonAction `json:"windows,omitempty"`
	MacOs   ButtonAction `json:"macos,omitempty"`
	Android ButtonAction `json:"android,omitempty"`
	Ios     ButtonAction `json:"ios,omitempty"`
}

func (e ExclusiveAction) Type() string {
	return "action"
}

func (e ExclusiveAction) String() string {
	return e.Default.String()
}

func (e ExclusiveAction) MarshalJSON() ([]byte, error) {
	type Embed ExclusiveAction
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  e.ActionType(),
		Embed: (Embed)(e),
	})
}

func (e ExclusiveAction) ActionType() string {
	return "exclusive"
}
