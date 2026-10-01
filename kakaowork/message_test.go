package kakaowork_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/kakaowork"
)

var (
	_ kakaowork.BubbleBlock = kakaowork.ButtonBlock{}
	_ kakaowork.BubbleBlock = (*kakaowork.ButtonBlock)(nil)
	_ kakaowork.BubbleBlock = kakaowork.ContextBlock{}
	_ kakaowork.BubbleBlock = (*kakaowork.ContextBlock)(nil)
	_ kakaowork.BubbleBlock = kakaowork.DescriptionBlock{}
	_ kakaowork.BubbleBlock = (*kakaowork.DescriptionBlock)(nil)
	_ kakaowork.BubbleBlock = kakaowork.DividerBlock{}
	_ kakaowork.BubbleBlock = (*kakaowork.DividerBlock)(nil)
	_ kakaowork.BubbleBlock = kakaowork.HeaderBlock{}
	_ kakaowork.BubbleBlock = (*kakaowork.HeaderBlock)(nil)
	_ kakaowork.BubbleBlock = kakaowork.ImageBlock{}
	_ kakaowork.BubbleBlock = (*kakaowork.ImageBlock)(nil)
	_ kakaowork.BubbleBlock = kakaowork.SectionBlock{}
	_ kakaowork.BubbleBlock = (*kakaowork.SectionBlock)(nil)
	_ kakaowork.BubbleBlock = kakaowork.TextBlock{}
	_ kakaowork.BubbleBlock = (*kakaowork.TextBlock)(nil)
)

type externalBubbleBlock struct{}

func (externalBubbleBlock) Type() string                 { return "text" }
func (externalBubbleBlock) String() string               { return "custom" }
func (externalBubbleBlock) MarshalJSON() ([]byte, error) { return []byte(`{"type":"text"}`), nil }
func (externalBubbleBlock) bubbleBlock()                 {}

func TestBubbleBlockRejectsExternalImplementation(t *testing.T) {
	for _, value := range []any{externalBubbleBlock{}, &externalBubbleBlock{}} {
		if _, ok := value.(kakaowork.BubbleBlock); ok {
			t.Errorf("external type %T implements BubbleBlock", value)
		}
	}
}
