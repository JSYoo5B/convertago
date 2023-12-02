package kakaowork

import "encoding/json"

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/#inlines

type Inline interface {
	MessageBubbleBlock
	InlineType() string
}

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/#styled

type InlineStyled struct {
	Text   string      `json:"text"`
	Bold   bool        `json:"bold,omitempty"`
	Italic bool        `json:"italic,omitempty"`
	Strike bool        `json:"strike,omitempty"`
	Color  InlineColor `json:"color,omitempty"`
}

type InlineColor string

const (
	InlineColorEmpty   = InlineColor("")
	InlineColorDefault = InlineColor("default")
	InlineColorRed     = InlineColor("red")
	InlineColorBlue    = InlineColor("blue")
	InlineColorGrey    = InlineColor("grey")
)

var inlineColorConstants = map[InlineColor]bool{
	InlineColorEmpty:   true,
	InlineColorDefault: true,
	InlineColorRed:     true,
	InlineColorBlue:    true,
	InlineColorGrey:    true,
}

func (i InlineStyled) Type() string {
	return "inline"
}

func (i InlineStyled) String() string {
	return i.Text
}

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

func (i InlineStyled) InlineType() string {
	return "styled"
}

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/#link

type InlineLink struct {
	Text string `json:"text"`
	Url  string `json:"url"`
}

func (i InlineLink) Type() string {
	return "inline"
}
func (i InlineLink) String() string {
	return i.Text
}

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

func (i InlineLink) InlineType() string {
	return "link"
}

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/#mention

type InlineMention struct {
	Text   string
	UserId int
}

func (i InlineMention) Type() string {
	return "inline"
}

func (i InlineMention) String() string {
	return i.Text
}

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

func (i InlineMention) InlineType() string {
	return "mention"
}
