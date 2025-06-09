package kakaowork_test

import (
	"github.com/JSYoo5B/convertago"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalKakaoworkMessage_Preview(t *testing.T) {
	type TestCase struct {
		Input   any
		Preview string
	}

	testCases := map[string]TestCase{
		"single field": {
			Input: struct {
				Preview string `kakaowork:"Preview"`
			}{
				Preview: "hello world",
			},
			Preview: "hello world",
		},
		"multiple fields": {
			Input: struct {
				Hello string `kakaowork:"Preview"`
				World string `kakaowork:"Preview"`
			}{
				Hello: "hello ",
				World: "world",
			},
			Preview: "hello world",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			output, err := convertago.MarshalKakaoworkMessage(tc.Input)

			assert.NoError(t, err)
			assert.Equal(t, tc.Preview, output.Preview)
			assert.Len(t, output.Blocks, 0)
		})
	}
}
