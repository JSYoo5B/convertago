package sample

import (
	"fmt"
	"time"
)

//go:generate go run github.com/JSYoo5B/convertago/cmd/convertago -type Notice,VerboseNotice -output zz_convertago.gen.go

type Picture struct {
	URL string `kakaowork:"part;slot=url" slack:"part;slot=url" googlechat:"part;slot=url"`
	Alt string `slack:"part;slot=alt" googlechat:"part;slot=alt"`
}

type Item struct {
	Prefix string `kakaowork:"text;group=line" slack:"rich_text;group=line" googlechat:"textParagraph;group=line"`
	Name   string `kakaowork:"text;group=line;style=italic" slack:"rich_text;group=line;style=italic" googlechat:"textParagraph;group=line;style=italic"`
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
	Optional string       `kakaowork:"divider;optional" slack:"divider;optional" googlechat:"divider;optional"`
	Dynamic  any          `kakaowork:"flatten" slack:"flatten" googlechat:"flatten"`
	Next     *Notice      `kakaowork:"flatten" slack:"flatten" googlechat:"flatten"`
	Ignored  map[string]string
}

type VerboseNotice struct {
	Text string `kakaowork:"text" slack:"rich_text" googlechat:"textParagraph"`
}

func (VerboseNotice) String() string { return "root text methods do not replace tags" }
