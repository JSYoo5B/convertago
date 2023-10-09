package kakaowork

import "encoding/json"

type HeaderBlock struct {
	Text  string
	Style HeaderStyle
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
	type _HeaderBlock struct {
		Type  string      `json:"type"`
		Text  string      `json:"text"`
		Style HeaderStyle `json:"style"`
	}

	if _, exists := headerStyleConstants[h.Style]; !exists {
		h.Style = HeaderStyleWhite
	}
	return json.Marshal(_HeaderBlock{
		Type:  h.Type(),
		Text:  h.Text,
		Style: h.Style,
	})
}
