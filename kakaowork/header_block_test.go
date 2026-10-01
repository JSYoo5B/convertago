package kakaowork_test

import (
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestHeaderBlock_Validate(t *testing.T) {
	v := validator.New()
	for _, test := range []struct {
		name  string
		block kakaowork.HeaderBlock
		valid bool
	}{
		{"unicode boundary", kakaowork.HeaderBlock{Text: strings.Repeat("가", 20)}, true},
		{"too long", kakaowork.HeaderBlock{Text: strings.Repeat("가", 21)}, false},
		{"line feed", kakaowork.HeaderBlock{Text: "첫째\n둘째"}, false},
		{"carriage return", kakaowork.HeaderBlock{Text: "첫째\r둘째"}, false},
		{"background", kakaowork.HeaderBlock{Text: "알림", Style: kakaowork.HeaderStyleBlue}, true},
		{"unknown background", kakaowork.HeaderBlock{Text: "알림", Style: "purple"}, false},
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
