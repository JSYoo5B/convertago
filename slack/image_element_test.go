package slack_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestImageElement_Validate(t *testing.T) {
	v := validator.New()
	require.NoError(t, v.Struct(slack.ImageElement{SlackFile: &slack.SlackFileObject{ID: "F123"}, AltText: "Image"}))
	require.Error(t, v.Struct(slack.ImageElement{AltText: "Image"}))
	require.Error(t, v.Struct(slack.ImageElement{ImageURL: "https://example.com/image", SlackFile: &slack.SlackFileObject{ID: "F123"}, AltText: "Image"}))
}
