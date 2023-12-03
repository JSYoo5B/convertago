package kakaowork

import (
	"encoding/json"
)

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/buttonblock/

type ButtonBlock struct {
	Text   string       `json:"text"`
	Style  ButtonStyle  `json:"style,omitempty"`
	Action ButtonAction `json:"action"`
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
