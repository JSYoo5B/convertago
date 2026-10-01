package slack_test

import (
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestSectionBlock_Validate(t *testing.T) {
	v := validator.New()
	slack.RegisterValidation(v)
	var absent *slack.PlainTextObject
	fields := make([]slack.TextObject, 11)
	for i := range fields {
		fields[i] = slack.PlainTextObject{Text: "Field"}
	}
	for _, test := range []struct {
		name  string
		block slack.SectionBlock
		valid bool
	}{
		{"text boundary", slack.SectionBlock{Text: &slack.PlainTextObject{Text: strings.Repeat("가", 3000)}}, true},
		{"text overflow", slack.SectionBlock{Text: slack.MrkdwnTextObject{Text: strings.Repeat("가", 3001)}}, false},
		{"field boundary", slack.SectionBlock{Fields: []slack.TextObject{&slack.PlainTextObject{Text: strings.Repeat("가", 2000)}}}, true},
		{"field overflow", slack.SectionBlock{Fields: []slack.TextObject{slack.MrkdwnTextObject{Text: strings.Repeat("가", 2001)}}}, false},
		{"10 fields", slack.SectionBlock{Fields: fields[:10]}, true},
		{"11 fields", slack.SectionBlock{Fields: fields}, false},
		{"missing text and fields", slack.SectionBlock{}, false},
		{"empty fields", slack.SectionBlock{Fields: []slack.TextObject{}}, false},
		{"typed nil field", slack.SectionBlock{Fields: []slack.TextObject{absent}}, false},
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
