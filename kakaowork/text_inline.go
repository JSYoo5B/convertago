package kakaowork

import "encoding/json"

// Inline 은 텍스트에 다양한 추가 서식을 적용할 때 사용하는 TextBlock 포맷입니다.
// 포맷은 InlineStyled, InlineLink, InlineMention 속성으로 구성됩니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/#inlines
//
// TextBlock.Inlines 내 캐스팅을 위해 제공되는 interface 입니다.
type Inline interface {
	InlineType() string
	BubbleBlock
}

// InlineStyled 는 TextBlock 에 Bold, Italic, Strike, Color 스타일을 지정하여 텍스트를 꾸미는 속성입니다.
// 속성은 굵은 글씨체, 기울어진 글씨체, 취소선, 글자 색상의 스타일을 중복하여 지정할 수 있습니다.
// 값을 명시하지 않은 경우에는 속성별 기본값이 적용됩니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/#styled
type InlineStyled struct {
	// Text 에 스타일을 적용할 텍스트를 입력합니다.
	Text string `json:"text"`
	// Bold 는 굵은 글씨체를 적용하며, 기본값은 false 입니다.
	Bold bool `json:"bold,omitempty"`
	// Italic 은 기울어진 글씨체를 적용하며, 기본값은 false 입니다.
	Italic bool `json:"italic,omitempty"`
	// Strike 는 취소선을 적용하며, 기본값은 false 입니다.
	Strike bool `json:"strike,omitempty"`
	// Color 는 글자 색상을 적용하며, 기본값은 검은색입니다.
	// InlineColorRed, InlineColorBlue, InlineColorGrey 로 설정할 수 있습니다.
	Color InlineColor `json:"color,omitempty"`
}

type InlineColor string

const (
	InlineColorDefault = InlineColor("default")
	InlineColorRed     = InlineColor("red")
	InlineColorBlue    = InlineColor("blue")
	InlineColorGrey    = InlineColor("grey")
)

func (i InlineStyled) InlineType() string { return "styled" }
func (i InlineStyled) Type() string       { return "inline" }
func (i InlineStyled) String() string     { return i.Text }
func (i InlineStyled) MarshalJSON() ([]byte, error) {
	if _, exists := inlineColorConstants[i.Color]; !exists {
		i.Color = InlineColorDefault
	}

	type Embed InlineStyled
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  i.InlineType(),
		Embed: (Embed)(i),
	})
}

var inlineColorConstants = map[InlineColor]bool{
	InlineColor(""):    true, // Empty allowed. (will be applied as default)
	InlineColorDefault: true,
	InlineColorRed:     true,
	InlineColorBlue:    true,
	InlineColorGrey:    true,
}

// InlineLink 는 TextBlock 에 다음과 같은 HTTP 또는 HTTPS 스킴을 적용할 수 있는 속성입니다.
// 메일 주소(mailto:)와 연락처(tel:)가 입력될 경우, 각 클라이언트 특성에 맞게 링크가 자동으로 연결합니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/#link
type InlineLink struct {
	// Text 에 링크를 적용할 텍스트를 입력합니다.
	Text string `json:"text"`
	// Url 에 연결할 링크 주소를 입력합니다.
	// "http:", "https:" 로 시작한다면 브라우저로 연결합니다.
	// "tel:" 로 시작한다면 통화 또는 브라우저로 연결합니다.
	// "mailto:" 로 시작한다면 메일 작성으로 연결합니다.
	Url string `json:"url"`
}

func (i InlineLink) InlineType() string { return "link" }
func (i InlineLink) Type() string       { return "inline" }
func (i InlineLink) String() string     { return i.Text }
func (i InlineLink) MarshalJSON() ([]byte, error) {
	type Embed InlineLink
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  i.InlineType(),
		Embed: (Embed)(i),
	})
}

// InlineMention 은 TextBlock 에서 사용자를 멘션할 수 있는 속성이며, 채팅방에 멘션된 참여자가 존재할 경우에만 동작합니다.
// 멘션된 사용자를 클릭하면, 해당 사용자의 프로필 카드가 표시됩니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/#mention
type InlineMention struct {
	// Text 에 멘션을 적용할 텍스트를 입력합니다.
	Text string
	// UserId 에 멘션할 멤버의 고유 정보를 입력합니다.
	UserId int
}

func (i InlineMention) InlineType() string { return "mention" }
func (i InlineMention) Type() string       { return "inline" }
func (i InlineMention) String() string     { return i.Text }
func (i InlineMention) MarshalJSON() ([]byte, error) {
	type Ref struct {
		Type  string `json:"type"`
		Value int    `json:"value"`
	}

	return json.Marshal(&struct {
		Type string `json:"type"`
		Text string `json:"text"`
		Ref  `json:"ref"`
	}{
		Type: i.InlineType(),
		Text: i.Text,
		Ref: Ref{
			Type:  "kw",
			Value: i.UserId,
		},
	})
}
