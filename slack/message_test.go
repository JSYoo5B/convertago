package slack_test

import (
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestMessage_Validate(t *testing.T) {
	v := validator.New()
	slack.RegisterValidation(v)
	blocks := make([]slack.Block, 51)
	for i := range blocks {
		blocks[i] = slack.DividerBlock{}
	}
	require.NoError(t, v.Struct(slack.Message{Blocks: blocks[:50]}))
	require.Error(t, v.Struct(slack.Message{Blocks: blocks}))
	var absent *slack.HeaderBlock
	require.Error(t, v.Struct(slack.Message{Blocks: []slack.Block{absent}}))

	message := slack.Message{Blocks: []slack.Block{
		slack.MarkdownBlock{Text: strings.Repeat("가", 6000)},
		&slack.MarkdownBlock{Text: strings.Repeat("나", 6000)},
	}}
	require.NoError(t, v.Struct(message))
	message.Blocks[1] = slack.MarkdownBlock{Text: strings.Repeat("나", 6001)}
	err := v.Struct(message)
	require.ErrorAs(t, err, new(validator.ValidationErrors))
	require.Equal(t, "markdown_max", err.(validator.ValidationErrors)[0].Tag())
}
