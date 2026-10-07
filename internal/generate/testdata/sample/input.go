package sample

import (
	"fmt"
	"time"
)

//go:generate go run github.com/JSYoo5B/convertago/cmd/convertago -type Notice,VerboseNotice,LayoutNotice -output zz_convertago.gen.go

type Picture struct {
	URL string `kakaowork:"part;slot=url" slack:"part;slot=url" googlechat:"part;slot=url"`
	Alt string `slack:"part;slot=alt" googlechat:"part;slot=alt"`
}

type Item struct {
	Prefix string `kakaowork:"text;group=line" slack:"rich_text;group=line" googlechat:"textParagraph;group=line"`
	Name   string `kakaowork:"text;group=line;style=italic" slack:"rich_text;group=line;style=italic" googlechat:"textParagraph;group=line;style=italic"`
}

type Alert struct {
	Word string `kakaowork:"styled;color=red;bold;italic" slack:"text;style=bold" googlechat:"part"`
}

type Label string

func (label Label) MarshalText() ([]byte, error) {
	if label == "bad" {
		return nil, fmt.Errorf("bad label")
	}
	return []byte(label), nil
}

type Caption string

func (caption Caption) String() string { return string(caption) }

type Notice struct {
	Title    string       `kakaowork:"header" slack:"header" googlechat:"header"`
	Prefix   string       `kakaowork:"text;group=body;omitempty" slack:"rich_text;group=body;omitempty" googlechat:"textParagraph;group=body;omitempty"`
	Between  string       `kakaowork:"text" slack:"section" googlechat:"textParagraph"`
	Name     string       `kakaowork:"text;group=body;style=bold" slack:"rich_text;group=body;style=bold" googlechat:"textParagraph;group=body;style=bold"`
	Image    *Picture     `kakaowork:"image_link" slack:"image" googlechat:"image"`
	Items    []Item       `kakaowork:"flatten" slack:"flatten" googlechat:"flatten"`
	Numbers  [2]int       `kakaowork:"text" slack:"rich_text" googlechat:"textParagraph"`
	Flags    []bool       `kakaowork:"text" slack:"rich_text" googlechat:"textParagraph"`
	Fraction float32      `kakaowork:"text;omitempty" slack:"rich_text;omitempty" googlechat:"textParagraph;omitempty"`
	Pointer  *int         `kakaowork:"text;omitempty" slack:"rich_text;omitempty" googlechat:"textParagraph;omitempty"`
	Label    Label        `kakaowork:"text;omitempty" slack:"rich_text;omitempty" googlechat:"textParagraph;omitempty"`
	Caption  fmt.Stringer `kakaowork:"text" slack:"rich_text" googlechat:"textParagraph"`
	Time     time.Time    `kakaowork:"text;omitempty" slack:"rich_text;omitempty" googlechat:"textParagraph;omitempty"`
	Optional string       `kakaowork:"header;style=bold;optional" slack:"header;style=bold;optional" googlechat:"header;style=bold;optional"`
	Alert    Alert        `kakaowork:"text" slack:"rich_text" googlechat:"textParagraph;horizontalAlignment=CENTER"`
	Wide     string       `kakaowork:"text;omitempty" slack:"section;expand;omitempty" googlechat:"decoratedText;wrapText;omitempty"`
	Dynamic  any          `kakaowork:"flatten" slack:"flatten" googlechat:"flatten"`
	Next     *Notice      `kakaowork:"flatten" slack:"flatten" googlechat:"flatten"`
	Ignored  map[string]string
}

type VerboseNotice struct {
	Text string `kakaowork:"text" slack:"rich_text" googlechat:"textParagraph"`
}

func (VerboseNotice) String() string { return "root text methods do not replace tags" }

type ActionInput struct {
	Name  string `kakaowork:"part;slot=name" slack:"part;slot=action_id" googlechat:"part;slot=function"`
	Value string `kakaowork:"part;slot=value" slack:"part;slot=value"`
}

type ButtonInput struct {
	Label  string      `kakaowork:"part" slack:"part" googlechat:"part"`
	Action ActionInput `kakaowork:"submit_action" slack:"flatten" googlechat:"action"`
}

type ButtonRow struct {
	Buttons []ButtonInput `kakaowork:"button" slack:"button" googlechat:"button"`
}

type LinkInput struct {
	Label string `kakaowork:"part" slack:"part;slot=text"`
	URL   string `kakaowork:"part;slot=url" slack:"part;slot=url"`
}

type RichInput struct {
	Before string    `kakaowork:"part" slack:"part"`
	Link   LinkInput `kakaowork:"link" slack:"link"`
	After  string    `kakaowork:"part;style=bold" slack:"part;style=bold"`
}

type SlackList struct {
	Style  string   `slack:"part;slot=style"`
	Border *int     `slack:"part;slot=border"`
	Items  []string `slack:"rich_text_section"`
}

type SlackRich struct {
	Before string    `slack:"part"`
	List   SlackList `slack:"rich_text_list"`
	Quote  string    `slack:"rich_text_quote"`
	After  string    `slack:"part;style=underline"`
}

type GoogleColumn struct {
	Text string `googlechat:"textParagraph"`
}

type GoogleColumns struct {
	Columns []GoogleColumn `googlechat:"column"`
}

type GoogleSection struct {
	Header  string        `googlechat:"part;slot=header"`
	Columns GoogleColumns `googlechat:"columns"`
}

type GoogleCard struct {
	Header   string          `googlechat:"header"`
	Sections []GoogleSection `googlechat:"section"`
}

type GoogleWrappedCard struct {
	ID   string     `googlechat:"part;slot=cardId"`
	Card GoogleCard `googlechat:"card"`
}

type LayoutNotice struct {
	Header  string              `kakaowork:"header" slack:"header"`
	Body    RichInput           `kakaowork:"text" slack:"rich_text"`
	Divider struct{}            `kakaowork:"divider" slack:"divider"`
	Buttons ButtonRow           `kakaowork:"action" slack:"actions"`
	Rich    SlackRich           `slack:"rich_text"`
	Cards   []GoogleWrappedCard `googlechat:"cardWithId"`
}
