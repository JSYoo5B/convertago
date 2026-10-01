package slack_test

import (
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestButtonElement_Validate(t *testing.T) {
	v := validator.New()
	slack.RegisterValidation(v)
	for _, test := range []struct {
		name   string
		change func(*slack.ButtonElement)
		valid  bool
	}{
		{"boundary", func(b *slack.ButtonElement) {
			b.Text.Text = strings.Repeat("가", 75)
			b.ActionID = strings.Repeat("a", 255)
			b.Value = strings.Repeat("a", 2000)
		}, true},
		{"label overflow", func(b *slack.ButtonElement) { b.Text.Text = strings.Repeat("가", 76) }, false},
		{"action ID overflow", func(b *slack.ButtonElement) { b.ActionID = strings.Repeat("a", 256) }, false},
		{"value overflow", func(b *slack.ButtonElement) { b.Value = strings.Repeat("a", 2001) }, false},
		{"accessibility overflow", func(b *slack.ButtonElement) { b.AccessibilityLabel = strings.Repeat("가", 76) }, false},
		{"prompt overflow", func(b *slack.ButtonElement) { b.AgentPrompt = strings.Repeat("a", 4001) }, false},
		{"application URL", func(b *slack.ButtonElement) { b.URL = "slack://channel?id=C123" }, true},
		{"relative URL", func(b *slack.ButtonElement) { b.URL = "/channel" }, false},
		{"hostless URL", func(b *slack.ButtonElement) { b.URL = "https:#channel" }, false},
		{"unknown style", func(b *slack.ButtonElement) { b.Style = "purple" }, false},
		{"nil audience member", func(b *slack.ButtonElement) { b.VisibleToUserIDs = []string{""} }, false},
		{"invalid confirmation", func(b *slack.ButtonElement) { b.Confirm = &slack.ConfirmationDialogObject{} }, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			button := slack.ButtonElement{Text: slack.PlainTextObject{Text: "Open"}}
			test.change(&button)
			err := v.Struct(button)
			if test.valid {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, new(validator.ValidationErrors))
			}
		})
	}
}

func TestConfirmationDialogObject_Validate(t *testing.T) {
	v := validator.New()
	slack.RegisterValidation(v)
	dialog := slack.ConfirmationDialogObject{
		Title:   slack.PlainTextObject{Text: strings.Repeat("가", 100)},
		Text:    &slack.MrkdwnTextObject{Text: strings.Repeat("가", 300)},
		Confirm: slack.PlainTextObject{Text: strings.Repeat("가", 30)},
		Deny:    slack.PlainTextObject{Text: strings.Repeat("가", 30)},
	}
	require.NoError(t, v.Struct(dialog))
	for _, test := range []struct {
		name   string
		change func(*slack.ConfirmationDialogObject)
	}{
		{"title", func(d *slack.ConfirmationDialogObject) { d.Title.Text += "나" }},
		{"body", func(d *slack.ConfirmationDialogObject) {
			d.Text = slack.PlainTextObject{Text: strings.Repeat("가", 301)}
		}},
		{"confirm", func(d *slack.ConfirmationDialogObject) { d.Confirm.Text += "나" }},
		{"deny", func(d *slack.ConfirmationDialogObject) { d.Deny.Text += "나" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			invalid := dialog
			test.change(&invalid)
			require.Error(t, v.Struct(invalid))
		})
	}
}
