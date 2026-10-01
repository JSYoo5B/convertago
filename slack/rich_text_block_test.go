package slack_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestRichTextList_Validate(t *testing.T) {
	v := validator.New()
	border := 0
	list := slack.RichTextList{Style: slack.RichTextListStyleOrdered, Elements: []slack.RichTextSection{
		{Elements: []slack.RichTextInline{slack.TextInline{Text: "Item"}}},
	}, Offset: 4, Border: &border}
	require.NoError(t, v.Struct(list))
	list.Style = slack.RichTextListStyleBullet
	require.Error(t, v.Struct(list))
	list.Offset = 0
	require.NoError(t, v.Struct(list))
	border = 2
	require.Error(t, v.Struct(list))
	border = 1
	list.Indent = -1
	require.Error(t, v.Struct(list))
	list.Indent = 0
	list.Elements = nil
	require.Error(t, v.Struct(list))
}

func TestRichTextBlock_ValidateNestedContent(t *testing.T) {
	v := validator.New()
	slack.RegisterValidation(v)
	block := slack.RichTextBlock{Elements: []slack.RichTextElement{
		&slack.RichTextSection{Elements: []slack.RichTextInline{
			&slack.LinkInline{URL: "https://example.com", Style: &slack.RichTextStyle{Code: true}},
		}},
	}}
	err := v.Struct(slack.Message{Blocks: []slack.Block{&block}})
	require.ErrorAs(t, err, new(validator.ValidationErrors))
	require.Equal(t, "Message.Blocks[0].Elements[0].Elements[0].Style.Code", err.(validator.ValidationErrors)[0].Namespace())
	block.Elements = []slack.RichTextElement{slack.RichTextSection{Elements: []slack.RichTextInline{
		slack.TextInline{Text: "Code", Style: &slack.RichTextStyle{Code: true}},
	}}}
	require.NoError(t, v.Struct(block))
}

func TestRichTextPreformatted_Validate(t *testing.T) {
	v := validator.New()
	slack.RegisterValidation(v)
	border := 0
	region := slack.RichTextPreformatted{Elements: []slack.PreformattedInline{
		slack.TextInline{Text: "Code"},
		&slack.LinkInline{URL: "https://example.com"},
	}, Border: &border}
	require.NoError(t, v.Struct(region))
	border = -1
	require.Error(t, v.Struct(region))
	var absent *slack.TextInline
	region.Border = nil
	region.Elements = []slack.PreformattedInline{absent}
	require.Error(t, v.Struct(region))
}

func TestRichTextQuote_Validate(t *testing.T) {
	v := validator.New()
	border := 1
	region := slack.RichTextQuote{Elements: []slack.RichTextInline{slack.TextInline{Text: "Quote"}}, Border: &border}
	require.NoError(t, v.Struct(region))
	border = 2
	require.Error(t, v.Struct(region))
	region.Border = nil
	region.Elements = nil
	require.Error(t, v.Struct(region))
}
