package kakaowork_test

import (
	"github.com/JSYoo5B/convertago"
	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalKakaoworkMessage_DividerBlock(t *testing.T) {
	type TestCase struct {
		Input  any
		Blocks []kakaowork.BubbleBlock
	}

	testCases := map[string]TestCase{
		"simple field": {
			Input: struct {
				Divider int `kakaowork:"divider"`
			}{Divider: 1},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.DividerBlock{},
			},
		},
		"case insensitive field": {
			Input: struct {
				Divider int `kakaowork:"DIVIDER"`
			}{Divider: 1},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.DividerBlock{},
			},
		},
		"underbar field": {
			Input: struct {
				_ int `kakaowork:"Divider"`
			}{},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.DividerBlock{},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			output, err := convertago.MarshalKakaoworkMessage(tc.Input)

			assert.NoError(t, err)
			assert.Empty(t, output.Preview)
			if len(output.Blocks) > 0 {
				assert.Equal(t, tc.Blocks, output.Blocks)
			} else {
				assert.Empty(t, output.Blocks)
			}
		})
	}
}
