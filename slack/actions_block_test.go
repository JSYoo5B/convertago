package slack_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestActionsBlock_Validate(t *testing.T) {
	v := validator.New()
	elements := make([]slack.ActionElement, 26)
	for i := range elements {
		elements[i] = slack.ButtonElement{Text: slack.PlainTextObject{Text: "Open"}}
	}
	require.NoError(t, v.Struct(slack.ActionsBlock{Elements: elements[:25]}))
	require.Error(t, v.Struct(slack.ActionsBlock{Elements: elements}))
	require.Error(t, v.Struct(slack.ActionsBlock{}))
	var absent *slack.ButtonElement
	require.Error(t, v.Struct(slack.ActionsBlock{Elements: []slack.ActionElement{absent}}))
}
