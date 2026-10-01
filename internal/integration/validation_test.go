package integration_test

import (
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestSharedValidatorConfiguration(t *testing.T) {
	for _, requiredStructs := range []bool{false, true} {
		name := "default"
		if requiredStructs {
			name = "required structs"
		}
		t.Run(name, func(t *testing.T) {
			v := validator.New()
			if requiredStructs {
				v = validator.New(validator.WithRequiredStructEnabled())
			}
			kakaowork.RegisterValidation(v)
			slack.RegisterValidation(v)
			googlechat.RegisterValidation(v)

			type messages struct {
				Kakaowork  kakaowork.Message
				Slack      slack.Message
				GoogleChat googlechat.Message
			}
			value := messages{
				Kakaowork: kakaowork.Message{Blocks: []kakaowork.BubbleBlock{kakaowork.HeaderBlock{Text: "알림"}}},
				Slack:     slack.Message{Blocks: []slack.Block{slack.HeaderBlock{Text: slack.PlainTextObject{Text: "Notice"}}}},
				GoogleChat: googlechat.Message{CardsV2: []googlechat.CardWithID{{Card: googlechat.Card{
					Header: &googlechat.CardHeader{Title: "Notice"},
					Sections: []googlechat.Section{{Widgets: []googlechat.Widget{{Content: googlechat.ButtonList{Buttons: []googlechat.Button{{
						Text: "Save", Color: &googlechat.Color{}, OnClick: googlechat.OnClick{Action: &googlechat.Action{Function: "save"}},
					}}}}}}},
				}}}},
			}
			require.NoError(t, v.Struct(value))

			value.Slack.Blocks[0] = slack.HeaderBlock{Text: slack.PlainTextObject{Text: strings.Repeat("가", 151)}}
			err := v.Struct(value)
			require.ErrorAs(t, err, new(validator.ValidationErrors))
			failures := err.(validator.ValidationErrors)
			require.Equal(t, "messages.Slack.Blocks[0].Text.Text", failures[0].Namespace())
			require.Equal(t, "150", failures[0].Param())
		})
	}
}
