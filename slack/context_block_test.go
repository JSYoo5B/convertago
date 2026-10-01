package slack_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestContextBlock_Validate(t *testing.T) {
	v := validator.New()
	elements := make([]slack.ContextElement, 11)
	for i := range elements {
		elements[i] = slack.PlainTextObject{Text: "Context"}
	}
	require.NoError(t, v.Struct(slack.ContextBlock{Elements: elements[:10]}))
	require.Error(t, v.Struct(slack.ContextBlock{Elements: elements}))
	require.Error(t, v.Struct(slack.ContextBlock{}))
	require.Error(t, v.Struct(slack.ContextBlock{Elements: []slack.ContextElement{nil}}))
}
