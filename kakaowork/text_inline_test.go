package kakaowork_test

import (
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

var (
	_ kakaowork.Inline = kakaowork.InlineStyled{}
	_ kakaowork.Inline = (*kakaowork.InlineStyled)(nil)
	_ kakaowork.Inline = kakaowork.InlineLink{}
	_ kakaowork.Inline = (*kakaowork.InlineLink)(nil)
	_ kakaowork.Inline = kakaowork.InlineMention{}
	_ kakaowork.Inline = (*kakaowork.InlineMention)(nil)
)

type externalInline struct {
	kakaowork.BubbleBlock
}

func (externalInline) InlineType() string { return "styled" }
func (externalInline) inline()            {}

func TestInlineRejectsExternalImplementation(t *testing.T) {
	value := externalInline{BubbleBlock: kakaowork.TextBlock{Text: "custom"}}
	for _, candidate := range []any{value, &value} {
		if _, ok := candidate.(kakaowork.Inline); ok {
			t.Errorf("external type %T implements Inline", candidate)
		}
	}
}

func TestInlineLink_Validate(t *testing.T) {
	v := validator.New()
	kakaowork.RegisterValidation(v)
	for _, test := range []struct {
		url   string
		valid bool
	}{
		{"https://example.com/path#fragment", true},
		{"http://example.com", true},
		{"mailto:jane@example.com", true},
		{"tel:+821012345678", true},
		{"ftp://example.com", false},
		{"/relative", false},
		{"https:#fragment", false},
		{"mailto:", false},
	} {
		t.Run(test.url, func(t *testing.T) {
			err := v.Struct(kakaowork.InlineLink{Text: "링크", Url: test.url})
			if test.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestInlineStyled_Validate(t *testing.T) {
	v := validator.New()
	inline := kakaowork.InlineStyled{Text: strings.Repeat("가", 500), Color: kakaowork.InlineColorGrey}
	require.NoError(t, v.Struct(inline))
	inline.Text += "나"
	require.Error(t, v.Struct(inline))
	inline.Text = "텍스트"
	inline.Color = "purple"
	require.Error(t, v.Struct(inline))
}

func TestInlineMention_Validate(t *testing.T) {
	v := validator.New()
	require.NoError(t, v.Struct(kakaowork.InlineMention{Text: "사용자", UserId: 1}))
	require.Error(t, v.Struct(kakaowork.InlineMention{Text: "사용자", UserId: 0}))
}
