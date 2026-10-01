package slack_test

import (
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestHeaderBlock_Validate(t *testing.T) {
	v := validator.New()
	slack.RegisterValidation(v)
	for _, test := range []struct {
		name  string
		block slack.HeaderBlock
		valid bool
	}{
		{"unicode boundary", slack.HeaderBlock{Text: slack.PlainTextObject{Text: strings.Repeat("가", 150)}, BlockID: strings.Repeat("a", 255), Level: 4}, true},
		{"long heading", slack.HeaderBlock{Text: slack.PlainTextObject{Text: strings.Repeat("가", 151)}}, false},
		{"long block ID", slack.HeaderBlock{Text: slack.PlainTextObject{Text: "Heading"}, BlockID: strings.Repeat("a", 256)}, false},
		{"level default", slack.HeaderBlock{Text: slack.PlainTextObject{Text: "Heading"}}, true},
		{"level range", slack.HeaderBlock{Text: slack.PlainTextObject{Text: "Heading"}, Level: 5}, false},
		{"empty heading", slack.HeaderBlock{}, false},
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
