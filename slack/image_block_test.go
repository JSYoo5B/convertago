package slack_test

import (
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestImageBlock_Validate(t *testing.T) {
	v := validator.New()
	slack.RegisterValidation(v)
	for _, test := range []struct {
		name  string
		block slack.ImageBlock
		valid bool
	}{
		{"public image", slack.ImageBlock{ImageURL: "https://example.com/image.png", AltText: "Image"}, true},
		{"Slack ID", slack.ImageBlock{SlackFile: &slack.SlackFileObject{ID: "F123"}, AltText: "Image"}, true},
		{"Slack URL", slack.ImageBlock{SlackFile: &slack.SlackFileObject{URL: "https://example.com/file"}, AltText: "Image"}, true},
		{"missing source", slack.ImageBlock{AltText: "Image"}, false},
		{"two sources", slack.ImageBlock{ImageURL: "https://example.com/image.png", SlackFile: &slack.SlackFileObject{ID: "F123"}, AltText: "Image"}, false},
		{"empty file", slack.ImageBlock{SlackFile: &slack.SlackFileObject{}, AltText: "Image"}, false},
		{"two file sources", slack.ImageBlock{SlackFile: &slack.SlackFileObject{ID: "F123", URL: "https://example.com/file"}, AltText: "Image"}, false},
		{"unsupported scheme", slack.ImageBlock{ImageURL: "ftp://example.com/image.png", AltText: "Image"}, false},
		{"long alt", slack.ImageBlock{ImageURL: "https://example.com/image.png", AltText: strings.Repeat("가", 2001)}, false},
		{"long title", slack.ImageBlock{ImageURL: "https://example.com/image.png", AltText: "Image", Title: &slack.PlainTextObject{Text: strings.Repeat("가", 2001)}}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := v.Struct(test.block)
			if test.valid {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, new(validator.ValidationErrors))
			}
		})
	}
}
