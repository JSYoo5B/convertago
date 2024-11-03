package kakaowork

import "encoding/json"

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
)

func (HeaderBlock) Type() string     { return "header" }
func (h HeaderBlock) String() string { return h.Text }
func (h HeaderBlock) MarshalJSON() ([]byte, error) {
	if exists := headerStyleConstants[h.Style]; !exists {
		h.Style = HeaderStyleWhite
	}

	type Embed HeaderBlock
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  h.Type(),
		Embed: (Embed)(h),
	})
}

var headerStyleConstants = map[HeaderStyle]bool{
	HeaderStyle(""):   true,
	HeaderStyleWhite:  true,
	HeaderStyleBlue:   true,
	HeaderStyleRed:    true,
	HeaderStyleYellow: true,
}
