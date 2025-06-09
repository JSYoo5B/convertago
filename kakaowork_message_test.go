package convertago_test

import (
	"github.com/JSYoo5B/convertago"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalKakaoworkMessage_InvalidInputs(t *testing.T) {
	type TestCase struct {
		Input   any
		IsValid bool
	}

	testCases := map[string]TestCase{
		"bool": {
			Input:   true,
			IsValid: false,
		},
		"int": {
			Input:   32,
			IsValid: false,
		},
		"uint": {
			Input:   uint(32),
			IsValid: false,
		},
		"uintptr": {
			Input:   uintptr(32),
			IsValid: false,
		},
		"float": {
			Input:   32.0,
			IsValid: false,
		},
		"complex": {
			Input:   1.2 + 3.4i,
			IsValid: false,
		},
		"array": {
			Input:   [3]int{1, 2, 3},
			IsValid: false,
		},
		"channel": {
			Input:   make(chan int),
			IsValid: false,
		},
		"function": {
			Input:   func() {},
			IsValid: false,
		},
		"interface": {
			Input:   new(error),
			IsValid: false,
		},
		"map": {
			Input:   map[int]int{1: 1, 2: 2, 3: 3},
			IsValid: false,
		},
		"slice": {
			Input:   []int{1, 2, 3},
			IsValid: false,
		},
		"string": {
			Input:   "hello world",
			IsValid: false,
		},
		"struct": {
			Input:   struct{}{},
			IsValid: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := convertago.MarshalKakaoworkMessage(tc.Input)

			if tc.IsValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})

		t.Run("pointer of "+name, func(t *testing.T) {
			_, err := convertago.MarshalKakaoworkMessage(&tc.Input)

			if tc.IsValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestMarshalKakaoworkMessage_InvalidTags(t *testing.T) {
	type TestCase any

	testCases := map[string]TestCase{
		"no kakaowork field": struct {
			Text string `not_kakaowork:"Divider"`
		}{
			Text: "hello world",
		},
		"kakaowork tag but field is private": struct {
			text string `kakaowork:"Divider"`
		}{
			text: "hello world",
		},
		"unknown kakaowork tag": struct {
			Text string `kakaowork:"Unknown"`
		}{
			Text: "hello world",
		},
	}

	for name, input := range testCases {
		t.Run(name, func(t *testing.T) {
			output, err := convertago.MarshalKakaoworkMessage(input)

			assert.NoError(t, err)
			assert.Empty(t, output.Preview)
			assert.Len(t, output.Blocks, 0)
		})
	}
}
