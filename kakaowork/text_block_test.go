package kakaowork_test

import (
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestTextBlock_Validate(t *testing.T) {
	v := validator.New()
	kakaowork.RegisterValidation(v)
	var absent *kakaowork.InlineStyled
	for _, test := range []struct {
		name  string
		block kakaowork.TextBlock
		valid bool
	}{
		{"unicode boundary", kakaowork.TextBlock{Text: strings.Repeat("가", 500)}, true},
		{"plain overflow", kakaowork.TextBlock{Text: strings.Repeat("가", 501)}, false},
		{"inline boundary", kakaowork.TextBlock{Inlines: []kakaowork.Inline{
			kakaowork.InlineStyled{Text: strings.Repeat("가", 250)},
			&kakaowork.InlineStyled{Text: strings.Repeat("나", 250)},
		}}, true},
		{"inline overflow", kakaowork.TextBlock{Inlines: []kakaowork.Inline{
			kakaowork.InlineStyled{Text: strings.Repeat("가", 250)},
			&kakaowork.InlineStyled{Text: strings.Repeat("나", 251)},
		}}, false},
		{"typed nil inline", kakaowork.TextBlock{Inlines: []kakaowork.Inline{absent}}, false},
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
