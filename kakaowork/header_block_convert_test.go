package kakaowork_test

import (
	"github.com/JSYoo5B/convertago"
	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalKakaoworkMessage_HeaderBlock(t *testing.T) {
	type TestCase struct {
		Input  any
		Blocks []kakaowork.BubbleBlock
	}

	testCases := map[string]TestCase{
		"simple field": {
			Input: struct {
				Header string `kakaowork:"Header"`
			}{Header: "hello world"},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.HeaderBlock{
					Text: "hello world",
				},
			},
		},
		"with style option": {
			Input: struct {
				Header string `kakaowork:"Header;blue"`
			}{Header: "hello world"},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.HeaderBlock{
					Text:  "hello world",
					Style: kakaowork.HeaderStyleBlue,
				},
			},
		},
		"case insensitive": {
			Input: struct {
				Header string `kakaowork:"Header;BLUE"`
			}{Header: "hello world"},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.HeaderBlock{
					Text:  "hello world",
					Style: kakaowork.HeaderStyleBlue,
				},
			},
		},
		"unknown option ignored": {
			Input: struct {
				Header string `kakaowork:"Header;blue_and_yellow"`
			}{Header: "hello world"},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.HeaderBlock{
					Text:  "hello world",
					Style: kakaowork.HeaderStyle(""),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			output, err := convertago.MarshalKakaoworkMessage(tc.Input)

			assert.NoError(t, err)
			assert.Empty(t, output.Preview)
			if len(tc.Blocks) > 0 {
				assert.Equal(t, tc.Blocks, output.Blocks)
			} else {
				assert.Empty(t, output.Blocks)
			}
		})
	}
}
