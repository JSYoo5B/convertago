package kakaowork

import (
	"encoding/json"
)

// HeaderBlock 은 말풍선의 최상단에만 지정할 수 있는 블록으로, 말풍선의 헤더를 색상으로 구분하여 표시할 수 있습니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/headerblock/
type HeaderBlock struct {
	// Text 에 내용을 입력.
	// 최대 20자까지 입력 가능,
	// 말풍선 사이즈에 따라 말줄임 처리,
	// 기본 Bold 로 처리,
	// 줄바꿈은 미지원,
	Text string `json:"text"`
	// Style 에 색상 설정, 기본값은 흰색.
	// 설정 가능한 값: HeaderStyleWhite, HeaderStyleBlue, HeaderStyleRed, HeaderStyleYellow
	Style HeaderStyle `json:"style"`
}

type HeaderStyle string

const (
	HeaderStyleWhite  = HeaderStyle("white")
	HeaderStyleBlue   = HeaderStyle("blue")
	HeaderStyleRed    = HeaderStyle("red")
	HeaderStyleYellow = HeaderStyle("yellow")
	// HeaderStylePlain 은 HeaderStyleWhite 와 동일
	HeaderStylePlain = HeaderStyle("plain")
	// HeaderStyleSuccess 는 HeaderStyleBlue 와 동일
	HeaderStyleSuccess = HeaderStyle("success")
	// HeaderStyleError 는 HeaderStyleRed 와 동일
	HeaderStyleError = HeaderStyle("error")
	// HeaderStyleWarning 은 HeaderStyleYellow 와 동일
	HeaderStyleWarning = HeaderStyle("warning")
)

func (HeaderBlock) BubbleType() string { return "header" }
func (h HeaderBlock) MarshalJSON() ([]byte, error) {
	h.Style = HeaderStyles[h.Style]

	type Embed HeaderBlock
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  h.BubbleType(),
		Embed: (Embed)(h),
	})
}

var HeaderStyles = map[HeaderStyle]HeaderStyle{
	HeaderStyle(""):   HeaderStyle(""),
	HeaderStyleWhite:  HeaderStyleWhite,
	HeaderStyleBlue:   HeaderStyleBlue,
	HeaderStyleRed:    HeaderStyleRed,
	HeaderStyleYellow: HeaderStyleYellow,

	HeaderStylePlain:   HeaderStyleWhite,
	HeaderStyleSuccess: HeaderStyleBlue,
	HeaderStyleError:   HeaderStyleRed,
	HeaderStyleWarning: HeaderStyleYellow,
}
