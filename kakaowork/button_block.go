package kakaowork

import (
	"encoding/json"
)

// ButtonBlock 은 말풍선에서 버튼을 표현하는 블록으로, 레이아웃 블록을 구성하는 엘리먼트의 속성으로 사용되기도 합니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/
type ButtonBlock struct {
	// Text 에 버튼이 표현할 텍스트를 입력
	Text string `json:"text"`
	// Style 은 버튼의 색상을 설정, 기본값은 회색
	// - 기본, 회색: ButtonStyleDefault 혹은 ButtonStyleGray
	// - 강조, 파랑: ButtonStylePrimary 혹은 ButtonStyleBlue
	// - 위험, 빨강: ButtonStyleDanger 혹은 ButtonStyleRed
	Style ButtonStyle `json:"style,omitempty"`
	// Action 에 버튼을 클릭했을 때 수행할 동작을 설정,
	// 적용 가능한 ButtonAction: OpenSystemBrowserAction, OpenInAppBrowserAction, OpenExternalAppAction, SubmitAction, CallModalAction, ExclusiveAction 참고
	Action ButtonAction `json:"action"`
}

type ButtonStyle string

const (
	ButtonStyleDefault = ButtonStyle("default")
	ButtonStylePrimary = ButtonStyle("primary")
	ButtonStyleDanger  = ButtonStyle("danger")
	// ButtonStyleGray 는 ButtonStyleDefault 와 동일
	ButtonStyleGray = ButtonStyleDefault
	// ButtonStyleBlue 는 ButtonStylePrimary 와 동일
	ButtonStyleBlue = ButtonStylePrimary
	// ButtonStyleRed 는 ButtonStyleDanger 와 동일
	ButtonStyleRed = ButtonStyleDanger
)

func (b ButtonBlock) Type() string   { return "button" }
func (b ButtonBlock) String() string { return b.Text + ": " + b.Action.String() }
func (b ButtonBlock) MarshalJSON() ([]byte, error) {
	if _, exists := buttonStyleConstants[b.Style]; !exists {
		b.Style = ButtonStyleDefault
	}

	type Embed ButtonBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  b.Type(),
		Embed: (Embed)(b),
	})
}

var buttonStyleConstants = map[ButtonStyle]bool{
	ButtonStyle(""): true,
	ButtonStyleGray: true,
	ButtonStyleBlue: true,
	ButtonStyleRed:  true,
}
