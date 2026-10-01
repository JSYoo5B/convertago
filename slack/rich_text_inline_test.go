package slack_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestRichTextInline_Validate(t *testing.T) {
	v := validator.New()
	slack.RegisterValidation(v)
	for _, test := range []struct {
		name   string
		inline any
		valid  bool
	}{
		{"mail link", slack.LinkInline{URL: "mailto:jane@example.com"}, true},
		{"application link", slack.LinkInline{URL: "slack://channel?id=C123"}, true},
		{"relative link", slack.LinkInline{URL: "/channel"}, false},
		{"hostless link", slack.LinkInline{URL: "https:#channel"}, false},
		{"code link", slack.LinkInline{URL: "https://example.com", Style: &slack.RichTextStyle{Code: true}}, false},
		{"styled user", slack.UserInline{UserID: "U123", Style: &slack.RichTextStyle{Bold: true}}, true},
		{"code user", slack.UserInline{UserID: "U123", Style: &slack.RichTextStyle{Code: true}}, false},
		{"missing user ID", slack.UserInline{}, false},
		{"named emoji", slack.EmojiInline{Name: "wave"}, true},
		{"missing emoji name", slack.EmojiInline{}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := v.Struct(test.inline)
			if test.valid {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, new(validator.ValidationErrors))
			}
		})
	}
}
