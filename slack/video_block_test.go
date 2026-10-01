package slack_test

import (
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestVideoBlock_Validate(t *testing.T) {
	v := validator.New()
	slack.RegisterValidation(v)
	for _, test := range []struct {
		name   string
		change func(*slack.VideoBlock)
		valid  bool
	}{
		{"boundary", func(b *slack.VideoBlock) {
			b.Title.Text = strings.Repeat("가", 199)
			b.AuthorName = strings.Repeat("a", 49)
		}, true},
		{"title overflow", func(b *slack.VideoBlock) { b.Title.Text = strings.Repeat("가", 200) }, false},
		{"description overflow", func(b *slack.VideoBlock) { b.Description = &slack.PlainTextObject{Text: strings.Repeat("가", 200)} }, false},
		{"author overflow", func(b *slack.VideoBlock) { b.AuthorName = strings.Repeat("a", 50) }, false},
		{"HTTP video", func(b *slack.VideoBlock) { b.VideoURL = "http://example.com/video" }, false},
		{"HTTP title", func(b *slack.VideoBlock) { b.TitleURL = "http://example.com/video" }, false},
		{"mixed-case HTTPS", func(b *slack.VideoBlock) { b.VideoURL = "HtTpS://example.com/video" }, true},
		{"invalid thumbnail", func(b *slack.VideoBlock) { b.ThumbnailURL = "/thumbnail" }, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			block := slack.VideoBlock{Title: slack.PlainTextObject{Text: "Video"}, VideoURL: "https://example.com/video", ThumbnailURL: "https://example.com/thumbnail.png", AltText: "Video"}
			test.change(&block)
			err := v.Struct(block)
			if test.valid {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, new(validator.ValidationErrors))
			}
		})
	}
}
