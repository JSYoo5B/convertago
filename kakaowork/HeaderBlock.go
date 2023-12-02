package kakaowork

import "encoding/json"

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/headerblock/

type HeaderBlock struct {
	Text  string      `json:"text"`
	Style HeaderStyle `json:"style"`
}

type HeaderStyle string

const (
	HeaderStyleEmpty  = HeaderStyle("")
	HeaderStyleWhite  = HeaderStyle("white")
	HeaderStyleBlue   = HeaderStyle("blue")
	HeaderStyleRed    = HeaderStyle("red")
	HeaderStyleYellow = HeaderStyle("yellow")
)

var headerStyleConstants = map[HeaderStyle]bool{
	HeaderStyleEmpty:  true,
	HeaderStyleWhite:  true,
	HeaderStyleBlue:   true,
	HeaderStyleRed:    true,
	HeaderStyleYellow: true,
}

func (h HeaderBlock) Type() string {
	return "header"
}

func (h HeaderBlock) String() string {
	return h.Text
}

func (h HeaderBlock) MarshalJSON() ([]byte, error) {
	if _, exists := headerStyleConstants[h.Style]; !exists {
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
