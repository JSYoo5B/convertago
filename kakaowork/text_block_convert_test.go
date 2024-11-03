package kakaowork

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTextBlock_ConvertText(t *testing.T) {
	type testCase struct {
		input       any
		output      TextBlock
		expectError bool
	}

	testCases := map[string]testCase{
		`Simple struct with kakaowork:"TextBlock" tag on string`: {
			input: struct {
				Text string `kakaowork:"TextBlock"`
			}{Text: "Test Message"},
			output: TextBlock{Text: "Test Message"},
		},
		`Simple struct with kakaowork:"text" tag on string`: {
			input: struct {
				Text string `kakaowork:"text"`
			}{Text: "Test Message"},
			output: TextBlock{Text: "Test Message"},
		},
		`Simple struct with kakaowork:"text" on int`: {
			input: struct {
				Value int `kakaowork:"text"`
			}{Value: 1234},
			output: TextBlock{Text: "1234"},
		},
		`Simple struct without any kakaowork tag`: {
			input: struct {
				Text string `not_kakaowork:"text"`
			}{Text: "Test Message"},
			expectError: false, // Not error, empty result
			output:      TextBlock{},
		},
		`Simple struct with kakaowork:"notText" tag`: {
			input: struct {
				Text string `kakaowork:"notText"`
			}{Text: "Test Message"},
			expectError: false,
			output:      TextBlock{},
		},
		`Simple struct with kakaowork:"text" tag, but private`: {
			input: struct {
				text string `kakaowork:"text"`
			}{text: "Test Message"},
			expectError: false, // Not error, empty result
			output:      TextBlock{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actual, err := ConvertTextBlock(tc.input)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.output, actual)
			}
		})
	}
}
