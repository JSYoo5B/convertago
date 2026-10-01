package convertago_test

import (
	"runtime"
	"testing"

	"github.com/JSYoo5B/convertago"
)

type benchmarkFlatMessageSource struct {
	Title  string `kakaowork:"header" slack:"header" googlechat:"header"`
	Prefix string `kakaowork:"text;group=body" slack:"rich_text;group=body" googlechat:"textParagraph;group=body"`
	Name   string `kakaowork:"text;group=body;style=bold" slack:"rich_text;group=body;style=bold" googlechat:"textParagraph;group=body;style=bold"`
	Count  int    `kakaowork:"text" slack:"rich_text" googlechat:"textParagraph"`
}

type benchmarkPictureSource struct {
	URL string `kakaowork:"part;slot=url" slack:"part;slot=url" googlechat:"part;slot=url"`
	Alt string `slack:"part;slot=alt" googlechat:"part;slot=alt"`
}

type benchmarkLineSource struct {
	Prefix string `kakaowork:"text;group=line" slack:"rich_text;group=line" googlechat:"textParagraph;group=line"`
	Name   string `kakaowork:"text;group=line;style=bold" slack:"rich_text;group=line;style=bold" googlechat:"textParagraph;group=line;style=bold"`
}

type benchmarkNestedMessageSource struct {
	Title   string                  `kakaowork:"header" slack:"header" googlechat:"header"`
	Picture *benchmarkPictureSource `kakaowork:"image_link" slack:"image" googlechat:"image"`
	Items   []benchmarkLineSource   `kakaowork:"flatten" slack:"flatten" googlechat:"flatten"`
	Footer  *string                 `kakaowork:"text;omitempty" slack:"rich_text;omitempty" googlechat:"textParagraph;omitempty"`
}

type benchmarkDynamicMessageSource struct {
	Content any `kakaowork:"flatten" slack:"flatten" googlechat:"flatten"`
}

func BenchmarkToKakaoworkMessage(b *testing.B) {
	benchmarkToMessage(b, convertago.ToKakaoworkMessage)
}

func BenchmarkToSlackMessage(b *testing.B) {
	benchmarkToMessage(b, convertago.ToSlackMessage)
}

func BenchmarkToGoogleChatMessage(b *testing.B) {
	benchmarkToMessage(b, convertago.ToGoogleChatMessage)
}

// benchmarkToMessage uses identical source values across the three messengers.
// The initial call warms the reflection caches outside the timer. Timed calls
// include source reading, assembly, validation, and native message rendering.
// No fixture implements generated accessors; serialization by the caller is
// excluded, while serialization performed inside a converter remains included.
func benchmarkToMessage[T any](b *testing.B, convert func(any, ...convertago.Option) (T, error)) {
	footer := "Complete"
	nested := &benchmarkNestedMessageSource{
		Title: "Report",
		Picture: &benchmarkPictureSource{
			URL: "https://example.com/image.png", Alt: "Report image",
		},
		Items: []benchmarkLineSource{
			{Prefix: "Hello ", Name: "Jane"},
			{Prefix: "Hello ", Name: "John"},
			{Prefix: "Hello ", Name: "Alex"},
		},
		Footer: &footer,
	}
	for _, tc := range []struct {
		name  string
		input any
	}{
		{
			name: "Flat",
			input: &benchmarkFlatMessageSource{
				Title: "Notice", Prefix: "Hello ", Name: "Jane", Count: 123,
			},
		},
		{name: "Nested", input: nested},
		{name: "Dynamic", input: &benchmarkDynamicMessageSource{Content: nested}},
	} {
		b.Run(tc.name, func(b *testing.B) {
			message, err := convert(tc.input)
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				message, err = convert(tc.input)
				if err != nil {
					b.Fatal(err)
				}
			}
			runtime.KeepAlive(message)
		})
	}
}
