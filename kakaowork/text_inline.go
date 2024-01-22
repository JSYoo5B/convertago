package kakaowork

import "encoding/json"

// Inline is type casting interface for TextBlock.Inlines.
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/#inlines
type Inline interface {
	InlineType() string
	BubbleBlock
}

// InlineStyled gives styles to text.
// Style are represented by each field, and it can be applied at the same time.
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/#styled
type InlineStyled struct {
	// Text to apply styles
	Text string `json:"text"`
	// Bold style, disabled by default
	Bold bool `json:"bold,omitempty"`
	// Italic style, disabled by default
	Italic bool `json:"italic,omitempty"`
	// Strike style, disabled by default
	Strike bool `json:"strike,omitempty"`
	// Color supports InlineColorRed, InlineColorBlue, InlineColorGrey.
	// InlineColorDefault of empty applies default plain text style
	Color InlineColor `json:"color,omitempty"`
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

// InlineLink links text to scheme (http/https, mailto, tel).
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/#link
type InlineLink struct {
	// Text to apply link
	Text string `json:"text"`
	// Url starts with "http:", "https:" will link to web browser.
	// Url starts with "tel:" will link to call.
	// Url starts with "mailto:" will link to mail client.
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

// InlineMention can mention the other user.
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/#mention
type InlineMention struct {
	// Text to apply mention
	Text string
	// UserId to mention
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
