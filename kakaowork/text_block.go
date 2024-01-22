package kakaowork

import (
	"encoding/json"
	"strings"
)

// TextBlock is a BubbleBlock to represent simple texts.
// When Text and Inlines differ, Inlines applied prior.
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/
type TextBlock struct {
	// Text represents plain text without formatting
	Text string `json:"text"`
	// Inlines are partial blocks for text with formatting
	// See also InlineStyled, InlineLink, InlineMention and Inline
	Inlines []Inline `json:"inlines,omitempty"`
}

func (t TextBlock) Type() string { return "text" }
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
	type Embed TextBlock
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  t.Type(),
		Embed: (Embed)(t),
	})
}
