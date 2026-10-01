// Package benchmarksource defines identical inputs for reflection and generated
// accessor benchmarks. It is used only by the repository's tests.
package benchmarksource

//go:generate go run ../../cmd/convertago -type FlatMessage,NestedMessage,DynamicMessage -output zz_convertago.gen.go

type FlatMessage struct {
	Title  string `kakaowork:"header" slack:"header" googlechat:"header"`
	Prefix string `kakaowork:"text;group=body" slack:"rich_text;group=body" googlechat:"textParagraph;group=body"`
	Name   string `kakaowork:"text;group=body;style=bold" slack:"rich_text;group=body;style=bold" googlechat:"textParagraph;group=body;style=bold"`
	Count  int    `kakaowork:"text" slack:"rich_text" googlechat:"textParagraph"`
}

type Picture struct {
	URL string `kakaowork:"part;slot=url" slack:"part;slot=url" googlechat:"part;slot=url"`
	Alt string `slack:"part;slot=alt" googlechat:"part;slot=alt"`
}

type Line struct {
	Prefix string `kakaowork:"text;group=line" slack:"rich_text;group=line" googlechat:"textParagraph;group=line"`
	Name   string `kakaowork:"text;group=line;style=bold" slack:"rich_text;group=line;style=bold" googlechat:"textParagraph;group=line;style=bold"`
}

type NestedMessage struct {
	Title   string   `kakaowork:"header" slack:"header" googlechat:"header"`
	Picture *Picture `kakaowork:"image_link" slack:"image" googlechat:"image"`
	Items   []Line   `kakaowork:"flatten" slack:"flatten" googlechat:"flatten"`
	Footer  *string  `kakaowork:"text;omitempty" slack:"rich_text;omitempty" googlechat:"textParagraph;omitempty"`
}

type DynamicMessage struct {
	Content any `kakaowork:"flatten" slack:"flatten" googlechat:"flatten"`
}
