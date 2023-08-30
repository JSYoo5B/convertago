package kakaowork

import (
	"encoding/json"
	"strings"
)

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/

type TextBlock struct {
	Text    string
	Inlines []Inline
}

func (t TextBlock) Type() string {
	return "text"
}

func (t TextBlock) String() string {
	if t.Inlines == nil {
		return t.Text
	} else {
		var inlineTexts []string
		for _, inline := range t.Inlines {
			inlineTexts = append(inlineTexts, inline.String())
		}
		return strings.Join(inlineTexts, "")
	}
}

func (t TextBlock) MarshalJSON() ([]byte, error) {
	type _TextBlock struct {
		Type    string   `json:"type"`
		Text    string   `json:"text"`
		Inlines []Inline `json:"inlines,omitempty"`
	}

	return json.Marshal(_TextBlock{
		Type:    t.Type(),
		Text:    t.Text,
		Inlines: t.Inlines,
	})
}

type Inline interface {
	MessageBubbleBlock
	InlineType() string
}

type InlineStyled struct {
	Text   string
	Bold   bool
	Italic bool
	Strike bool
	Color  InlineColor
}

func (i InlineStyled) Type() string {
	return "inline"
}

func (i InlineStyled) String() string {
	return i.Text
}

func (i InlineStyled) MarshalJSON() ([]byte, error) {
	type _InlineStyled struct {
		Type   string      `json:"type"`
		Text   string      `json:"text"`
		Bold   bool        `json:"bold,omitempty"`
		Italic bool        `json:"italic,omitempty"`
		Strike bool        `json:"strike,omitempty"`
		Color  InlineColor `json:"color,omitempty"`
	}

	return json.Marshal(_InlineStyled{
		Type:   i.InlineType(),
		Text:   i.Text,
		Bold:   i.Bold,
		Italic: i.Italic,
		Strike: i.Strike,
		Color:  i.Color,
	})
}

func (i InlineStyled) InlineType() string {
	return "styled"
}

type InlineColor string

const (
	InlineColorDefault = InlineColor("default")
	InlineColorRed     = InlineColor("red")
	InlineColorBlue    = InlineColor("blue")
	InlineColorGrey    = InlineColor("grey")
)

type InlineLink struct {
	Text string
	Url  string
}

func (i InlineLink) Type() string {
	return "inline"
}
func (i InlineLink) String() string {
	return i.Text
}

func (i InlineLink) MarshalJSON() ([]byte, error) {
	type _InlineLink struct {
		Type string `json:"type"`
		Text string `json:"text"`
		Url  string `json:"url"`
	}

	return json.Marshal(_InlineLink{
		Type: i.InlineType(),
		Text: i.Text,
		Url:  i.Url,
	})
}

func (i InlineLink) InlineType() string {
	return "link"
}

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
	type _InlineMention struct {
		Type string `json:"type"`
		Text string `json:"text"`
		Ref  struct {
			Type  string `json:"type"`
			Value int    `json:"value"`
		} `json:"ref"`
	}

	return json.Marshal(_InlineMention{
		Type: i.InlineType(),
		Text: i.Text,
		Ref: struct {
			Type  string `json:"type"`
			Value int    `json:"value"`
		}{
			Type:  "kw",
			Value: i.UserId,
		},
	})
}

func (i InlineMention) InlineType() string {
	return "mention"
}
