package kakaowork_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/kakaowork"
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
