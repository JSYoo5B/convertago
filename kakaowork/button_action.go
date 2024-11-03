package kakaowork

import "encoding/json"

type ButtonAction interface {
	ActionType() string
	BubbleBlock
}

// OpenSystemBrowserAction 은 Value 속성 값의 URL 을 시스템 브라우저로 연결합니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/#open_system_browser
type OpenSystemBrowserAction struct {
	// Name 은 어떤 ButtonBlock 을 클릭했는지 구분하기 위해 사용,
	// (주로 API 호출 시 사용)
	Name string `json:"name,omitempty"`
	// Value 는 시스템 브라우저에서 연결할 URL 을 설정
	Value string `json:"value"`
}

func (OpenSystemBrowserAction) ActionType() string { return "open_system_browser" }
func (OpenSystemBrowserAction) Type() string       { return "action" }
func (o OpenSystemBrowserAction) String() string   { return o.Value }
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

// OpenInAppBrowserAction 은 URL 을 카카오워크 내부 팝업 브라우저로 호출할 수 있습니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/#open_inapp_browser
type OpenInAppBrowserAction struct {
	// Name 은 어떤 ButtonBlock 을 클릭했는지 구분하기 위해 사용,
	// (주로 API 호출 시 사용)
	Name string `json:"name,omitempty"`
	// Value 는 카카오워크 내부 팝업에서 연결할 URL 을 설정
	Value string `json:"value"`
	// Standalone 은 PC 인앱 브라우저를 별도 창으로 띄우려 할 때 true 로 설정
	Standalone bool `json:"standalone,omitempty"`
	// Width 는 Standalone 사용 시 가로 사이즈를 설정 (기본 사이즈: 980)
	Width int `json:"width,omitempty"`
	// Height 는 Standalone 사용 시 세로 사이즈를 설정 (기본 사이즈: 720)
	Height int `json:"height,omitempty"`
}

func (OpenInAppBrowserAction) ActionType() string { return "open_inapp_browser" }
func (OpenInAppBrowserAction) Type() string       { return "action" }
func (o OpenInAppBrowserAction) String() string   { return o.Value }
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

// OpenExternalAppAction 은 Custom App Scheme 을 통해 외부 앱 실행을 할 수 있습니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/#open_external_app
type OpenExternalAppAction struct {
	// Name 은 어떤 ButtonBlock 을 클릭했는지 구분하기 위해 사용,
	// (주로 API 호출 시 사용)
	Name string `json:"name,omitempty"`
	// Value 는 사용자의 디바이스에 연결할 Custom App Scheme 을 입력
	Value string `json:"value"`
}

func (OpenExternalAppAction) ActionType() string { return "open_external_app" }
func (OpenExternalAppAction) Type() string       { return "action" }
func (o OpenExternalAppAction) String() string   { return o.Value }
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

// SubmitAction 은 반응형 메시지 봇을 개발 할 때 사용자의 입력 선택을 제공할 수 있습니다.
// Name 과 Value 를 고객 서버에 콜백으로 전달할 수 있습니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/#submit_action
type SubmitAction struct {
	// Name 은 어떤 ButtonBlock 을 클릭했는지 구분하기 위해 사용,
	// (SubmitAction 은 해당 항목 필수 설정)
	Name string `json:"name"`
	// Value 는 사용자가 전송하게 될 값을 입력
	Value string `json:"value"`
}

func (SubmitAction) ActionType() string { return "submit_action" }
func (SubmitAction) Type() string       { return "action" }
func (s SubmitAction) String() string   { return s.Value }
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

// CallModalAction 은 버튼 클릭 시 Modal 을 화면에 띄우고 사전에 등록한 블록 데이터를 요청합니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/#call_modal
type CallModalAction struct {
	// Name 은 어떤 ButtonBlock 을 클릭했는지 구분하기 위해 사용,
	// (주로 API 호출 시 사용)
	Name string `json:"name,omitempty"`
	// Value 에 사전 등록한 Modal 을 찾기 위한 초기 설정값을 입력
	Value string `json:"value"`
}

func (CallModalAction) ActionType() string { return "call_modal" }
func (CallModalAction) Type() string       { return "action" }
func (c CallModalAction) String() string   { return c.Value }
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

// ExclusiveAction 은 특정 플랙폼 또는 특정 OS 별로 다른 ButtonAction 을 설정할 수 있습니다.
// 설정은 특정 OS, 특정 플랫폼, Default 순으로 적용됩니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/#exclusive-action
type ExclusiveAction struct {
	// Default 는 고객의 OS, 플랫폼에 설정된 값이 없을 때 기본 실행할 동작을 설정,
	// (다른 항목 설정여부와 관계 없이 해당 항목 필수 설정)
	Default ButtonAction `json:"default"`
	// Pc 는 PC 환경에서 MacOs, Windows 등 OS 값 설정이 없는 경우 실행할 동작을 설정
	Pc ButtonAction `json:"pc,omitempty"`
	// Mobile 은 모바일 환경에서 Android, Ios 등 OS 값 설정이 없는 경우 실행할 동작을 설정
	Mobile ButtonAction `json:"mobile,omitempty"`
	// Windows 에서 수행할 동작을 설정
	Windows ButtonAction `json:"windows,omitempty"`
	// MacOs 에서 수행할 동작을 설정
	MacOs ButtonAction `json:"macos,omitempty"`
	// Android 에서 수행할 동작을 설정
	Android ButtonAction `json:"android,omitempty"`
	// Ios 에서 수행할 동작을 설정
	Ios ButtonAction `json:"ios,omitempty"`
}

func (ExclusiveAction) ActionType() string { return "exclusive" }
func (ExclusiveAction) Type() string       { return "action" }
func (e ExclusiveAction) String() string   { return e.Default.String() }
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
