package kakaowork_test

import (
	"github.com/JSYoo5B/convertago"
	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalKakaoworkMessage_ImageBlock(t *testing.T) {
	type TestCase struct {
		Input  any
		Blocks []kakaowork.BubbleBlock
	}

	testCases := map[string]TestCase{
		"simple field": {
			Input: struct {
				Image string `kakaowork:"image_link"`
			}{Image: "https://picsum.photos/200/300"},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.ImageBlock{Url: "https://picsum.photos/200/300"},
			},
		},
		"unknown option": {
			Input: struct {
				Image string `kakaowork:"image;unknown_option"`
			}{Image: "https://picsum.photos/200/300"},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.ImageBlock{Url: "https://picsum.photos/200/300"},
			},
		},
		"prefix field": {
			Input: struct {
				Height string `kakaowork:"image;prefix=https://picsum.photos/200/"`
			}{Height: "300"},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.ImageBlock{Url: "https://picsum.photos/200/300"},
			},
		},
		"format field": {
			Input: struct {
				Width int `kakaowork:"image;format=https://picsum.photos/%d/300"`
			}{Width: 200},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.ImageBlock{Url: "https://picsum.photos/200/300"},
			},
		},
		"format with slice": {
			Input: struct {
				Size []string `kakaowork:"image;format=https://picsum.photos/{0}/{1}"`
			}{Size: []string{"200", "300"}},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.ImageBlock{Url: "https://picsum.photos/200/300"},
			},
		},
		"format with array": {
			Input: struct {
				Size [2]int `kakaowork:"image;format=https://picsum.photos/{0}/{1}"`
			}{Size: [2]int{200, 300}},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.ImageBlock{Url: "https://picsum.photos/200/300"},
			},
		},
		"format with map": {
			Input: struct {
				Params map[string]string `kakaowork:"image;format=https://picsum.photos/{w}/{h}"`
			}{
				Params: map[string]string{
					"w": "200",
					"h": "300",
				},
			},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.ImageBlock{Url: "https://picsum.photos/200/300"},
			},
		},
		"format with struct": {
			Input: struct {
				Params struct{ w, h int } `kakaowork:"image;format=https://picsum.photos/{w}/{h}"`
			}{
				Params: struct{ w, h int }{
					w: 200,
					h: 300,
				},
			},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.ImageBlock{Url: "https://picsum.photos/200/300"},
			},
		},
		"format with slice pointer": {
			Input: struct {
				Size *[]string `kakaowork:"image;format=https://picsum.photos/{0}/{1}"`
			}{Size: &[]string{"200", "300"}},
			Blocks: []kakaowork.BubbleBlock{kakaowork.ImageBlock{
				Url: "https://picsum.photos/200/300",
			}},
		},
		"format with map pointer": {
			Input: struct {
				Params *map[string]int `kakaowork:"image;fmt=https://picsum.photos/{w}/{h}"`
			}{
				Params: &map[string]int{
					"w": 200,
					"h": 300,
				},
			},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.ImageBlock{Url: "https://picsum.photos/200/300"},
			},
		},
		"format with struct pointer": {
			Input: struct {
				Params *struct{ w, h int } `kakaowork:"image;format=https://picsum.photos/{w}/{h}"`
			}{
				Params: &struct{ w, h int }{
					w: 200,
					h: 300,
				},
			},
			Blocks: []kakaowork.BubbleBlock{
				kakaowork.ImageBlock{Url: "https://picsum.photos/200/300"},
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
