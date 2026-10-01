package kakaowork_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
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

func TestMessage_Validate(t *testing.T) {
	v := validator.New()
	kakaowork.RegisterValidation(v)
	var absent *kakaowork.HeaderBlock
	for _, test := range []struct {
		name   string
		blocks []kakaowork.BubbleBlock
		valid  bool
	}{
		{"header first", []kakaowork.BubbleBlock{&kakaowork.HeaderBlock{Text: "제목"}, kakaowork.TextBlock{Text: "본문"}}, true},
		{"header later", []kakaowork.BubbleBlock{kakaowork.TextBlock{Text: "본문"}, &kakaowork.HeaderBlock{Text: "제목"}}, false},
		{"duplicate headers", []kakaowork.BubbleBlock{kakaowork.HeaderBlock{Text: "제목"}, kakaowork.HeaderBlock{Text: "다음 제목"}}, false},
		{"nil block", []kakaowork.BubbleBlock{nil}, false},
		{"typed nil block", []kakaowork.BubbleBlock{absent}, false},
		{"inline at root", []kakaowork.BubbleBlock{kakaowork.InlineStyled{Text: "본문"}}, false},
		{"invalid nested image", []kakaowork.BubbleBlock{kakaowork.SectionBlock{Accessory: &kakaowork.ImageBlock{Url: "/image"}}}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := v.Struct(kakaowork.Message{Blocks: test.blocks})
			if test.valid {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, new(validator.ValidationErrors))
			}
		})
	}
}
