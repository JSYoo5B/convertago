package slack

import (
	"errors"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/internal/conversion"
	"github.com/JSYoo5B/convertago/internal/validation"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

type violation struct {
	path string
	rule validation.Rule
}

// violations returns every rule violation of message without stopping at Fatal ones.
func violations(message Message) []violation {
	var result []violation
	validation.Walk("$", message.check, func(path string, v validation.Violation) {
		result = append(result, violation{path, v.Rule})
	})
	return result
}

func plain(text string) PlainTextObject { return PlainTextObject{Text: text} }

func messageOf(blocks ...Block) Message { return Message{Text: "Notice", Blocks: blocks} }

func button() ButtonElement { return ButtonElement{Text: plain("Open")} }

func buttonWith(edit func(*ButtonElement)) Message {
	b := button()
	edit(&b)
	return messageOf(ActionsBlock{Elements: []ActionElement{b}})
}

func dialog() *ConfirmationDialogObject {
	return &ConfirmationDialogObject{Title: plain("Sure?"), Text: plain("Delete it"), Confirm: plain("Yes"), Deny: plain("No")}
}

func video() VideoBlock {
	return VideoBlock{Title: plain("Video"), VideoURL: "https://example.com/video", AltText: "video", ThumbnailURL: "https://example.com/thumb.png"}
}

func videoWith(edit func(*VideoBlock)) Message {
	v := video()
	edit(&v)
	return messageOf(v)
}

func image() ImageBlock { return ImageBlock{ImageURL: "https://example.com/a.png", AltText: "image"} }

func imageWith(edit func(*ImageBlock)) Message {
	i := image()
	edit(&i)
	return messageOf(i)
}

func rich(inlines ...RichTextInline) Message {
	return messageOf(RichTextBlock{Elements: []RichTextElement{RichTextSection{Elements: inlines}}})
}

func repeat[T any](value T, count int) []T {
	result := make([]T, count)
	for i := range result {
		result[i] = value
	}
	return result
}

func one() *int {
	value := 1
	return &value
}

type ruleCase struct {
	rule    validation.Rule
	path    string
	message Message
}

var ruleCases = []ruleCase{
	{ruleMessageContentRequired, "$", Message{}},
	{ruleMessageBlocksCount, "$.blocks", messageOf(repeat[Block](DividerBlock{}, 51)...)},
	{ruleMessageTextFallback, "$.text", Message{Blocks: []Block{DividerBlock{}}}},
	{ruleMessageTextLength, "$.text", Message{Text: strings.Repeat("x", 40001)}},
	{ruleMarkdownTextTotal, "$.blocks", messageOf(MarkdownBlock{Text: strings.Repeat("x", 6000)}, MarkdownBlock{Text: strings.Repeat("x", 6001)})},
	{ruleMarkdownTextRequired, "$.blocks[0].text", messageOf(MarkdownBlock{})},
	{ruleElementRequired, "$.blocks[0]", messageOf(nil)},
	{ruleBlockIDLength, "$.blocks[0].block_id", messageOf(DividerBlock{BlockID: strings.Repeat("x", 256)})},

	{ruleTextRequired, "$.blocks[0].text.text", messageOf(HeaderBlock{Text: plain("")})},
	{ruleTextLength, "$.blocks[0].text.text", messageOf(SectionBlock{Text: MrkdwnTextObject{Text: strings.Repeat("x", 3001)}})},

	{ruleHeaderTextLength, "$.blocks[0].text", messageOf(HeaderBlock{Text: plain(strings.Repeat("x", 151))})},
	{ruleHeaderLevelValue, "$.blocks[0].level", messageOf(HeaderBlock{Text: plain("Notice"), Level: 5})},
	{ruleSectionContent, "$.blocks[0]", messageOf(SectionBlock{})},
	{ruleSectionFieldCount, "$.blocks[0].fields", messageOf(SectionBlock{Fields: repeat[TextObject](plain("x"), 11)})},
	{ruleSectionFieldLen, "$.blocks[0].fields[0]", messageOf(SectionBlock{Fields: []TextObject{plain(strings.Repeat("x", 2001))}})},

	{ruleImageSource, "$.blocks[0]", imageWith(func(i *ImageBlock) { i.ImageURL = "" })},
	{ruleImageURLLength, "$.blocks[0].image_url", imageWith(func(i *ImageBlock) { i.ImageURL = "https://example.com/" + strings.Repeat("x", 3000) })},
	{ruleImageURLFormat, "$.blocks[0].image_url", imageWith(func(i *ImageBlock) { i.ImageURL = "/a.png" })},
	{ruleImageAltRequired, "$.blocks[0].alt_text", imageWith(func(i *ImageBlock) { i.AltText = "" })},
	{ruleImageAltLength, "$.blocks[0].alt_text", imageWith(func(i *ImageBlock) { i.AltText = strings.Repeat("x", 2001) })},
	{ruleImageTitleLength, "$.blocks[0].title", imageWith(func(i *ImageBlock) { title := plain(strings.Repeat("x", 2001)); i.Title = &title })},
	{ruleSlackFileSource, "$.blocks[0].slack_file", imageWith(func(i *ImageBlock) { i.ImageURL, i.SlackFile = "", &SlackFileObject{} })},
	{ruleSlackFileURL, "$.blocks[0].slack_file.url", imageWith(func(i *ImageBlock) { i.ImageURL, i.SlackFile = "", &SlackFileObject{URL: "file"} })},

	{ruleActionsCount, "$.blocks[0].elements", messageOf(ActionsBlock{Elements: repeat[ActionElement](button(), 26)})},
	{ruleActionsRequired, "$.blocks[0].elements", messageOf(ActionsBlock{})},
	{ruleContextCount, "$.blocks[0].elements", messageOf(ContextBlock{Elements: repeat[ContextElement](plain("x"), 11)})},
	{ruleContextRequired, "$.blocks[0].elements", messageOf(ContextBlock{})},

	{ruleButtonTextLength, "$.blocks[0].elements[0].text", buttonWith(func(b *ButtonElement) { b.Text = plain(strings.Repeat("x", 76)) })},
	{ruleButtonTextTruncation, "$.blocks[0].elements[0].text", buttonWith(func(b *ButtonElement) { b.Text = plain(strings.Repeat("x", 31)) })},
	{ruleButtonActionIDLength, "$.blocks[0].elements[0].action_id", buttonWith(func(b *ButtonElement) { b.ActionID = strings.Repeat("x", 256) })},
	{ruleButtonURLLength, "$.blocks[0].elements[0].url", buttonWith(func(b *ButtonElement) { b.URL = "https://example.com/" + strings.Repeat("x", 3000) })},
	{ruleButtonURLFormat, "$.blocks[0].elements[0].url", buttonWith(func(b *ButtonElement) { b.URL = "relative" })},
	{ruleButtonValueLength, "$.blocks[0].elements[0].value", buttonWith(func(b *ButtonElement) { b.Value = strings.Repeat("x", 2001) })},
	{ruleButtonStyleValue, "$.blocks[0].elements[0].style", buttonWith(func(b *ButtonElement) { b.Style = "green" })},
	{ruleButtonLabelLength, "$.blocks[0].elements[0].accessibility_label", buttonWith(func(b *ButtonElement) { b.AccessibilityLabel = strings.Repeat("x", 76) })},
	{ruleButtonPromptLength, "$.blocks[0].elements[0].agent_prompt", buttonWith(func(b *ButtonElement) { b.AgentPrompt = strings.Repeat("x", 4001) })},
	{ruleButtonPromptDisplay, "$.blocks[0].elements[0].agent_prompt_display", buttonWith(func(b *ButtonElement) { b.AgentPromptDisplay = "Ask" })},
	{ruleButtonVisibleRequired, "$.blocks[0].elements[0].visible_to_user_ids[0]", buttonWith(func(b *ButtonElement) { b.VisibleToUserIDs = []string{""} })},

	{ruleConfirmTitleLength, "$.blocks[0].elements[0].confirm.title", buttonWith(func(b *ButtonElement) { b.Confirm = dialog(); b.Confirm.Title = plain(strings.Repeat("x", 101)) })},
	{ruleConfirmTextRequired, "$.blocks[0].elements[0].confirm.text", buttonWith(func(b *ButtonElement) { b.Confirm = dialog(); b.Confirm.Text = nil })},
	{ruleConfirmTextLength, "$.blocks[0].elements[0].confirm.text", buttonWith(func(b *ButtonElement) { b.Confirm = dialog(); b.Confirm.Text = plain(strings.Repeat("x", 301)) })},
	{ruleConfirmConfirmLength, "$.blocks[0].elements[0].confirm.confirm", buttonWith(func(b *ButtonElement) { b.Confirm = dialog(); b.Confirm.Confirm = plain(strings.Repeat("x", 31)) })},
	{ruleConfirmDenyLength, "$.blocks[0].elements[0].confirm.deny", buttonWith(func(b *ButtonElement) { b.Confirm = dialog(); b.Confirm.Deny = plain(strings.Repeat("x", 31)) })},
	{ruleConfirmStyleValue, "$.blocks[0].elements[0].confirm.style", buttonWith(func(b *ButtonElement) { b.Confirm = dialog(); b.Confirm.Style = "green" })},

	{ruleVideoTitleLength, "$.blocks[0].title", videoWith(func(v *VideoBlock) { v.Title = plain(strings.Repeat("x", 200)) })},
	{ruleVideoDescriptionLength, "$.blocks[0].description", videoWith(func(v *VideoBlock) { d := plain(strings.Repeat("x", 200)); v.Description = &d })},
	{ruleVideoAuthorLength, "$.blocks[0].author_name", videoWith(func(v *VideoBlock) { v.AuthorName = strings.Repeat("x", 50) })},
	{ruleVideoURLHTTPS, "$.blocks[0].video_url", videoWith(func(v *VideoBlock) { v.VideoURL = "http://example.com/video" })},
	{ruleVideoTitleURLHTTPS, "$.blocks[0].title_url", videoWith(func(v *VideoBlock) { v.TitleURL = "http://example.com" })},
	{ruleVideoAltRequired, "$.blocks[0].alt_text", videoWith(func(v *VideoBlock) { v.AltText = "" })},
	{ruleVideoThumbnailFormat, "$.blocks[0].thumbnail_url", videoWith(func(v *VideoBlock) { v.ThumbnailURL = "" })},
	{ruleVideoProviderIcon, "$.blocks[0].provider_icon_url", videoWith(func(v *VideoBlock) { v.ProviderIconURL = "icon.png" })},

	{ruleRichTextRequired, "$.blocks[0].elements", messageOf(RichTextBlock{})},
	{ruleListStyleValue, "$.blocks[0].elements[0].style", messageOf(RichTextBlock{Elements: []RichTextElement{RichTextList{Style: "dash", Elements: []RichTextSection{{Elements: []RichTextInline{TextInline{Text: "x"}}}}}}})},
	{ruleListOffsetOrdered, "$.blocks[0].elements[0].offset", messageOf(RichTextBlock{Elements: []RichTextElement{RichTextList{Style: RichTextListStyleBullet, Offset: 1, Elements: []RichTextSection{{Elements: []RichTextInline{TextInline{Text: "x"}}}}}}})},
	{ruleListNumberValue, "$.blocks[0].elements[0].indent", messageOf(RichTextBlock{Elements: []RichTextElement{RichTextList{Style: RichTextListStyleBullet, Indent: -1, Elements: []RichTextSection{{Elements: []RichTextInline{TextInline{Text: "x"}}}}}}})},
	{ruleBorderValue, "$.blocks[0].elements[0].border", messageOf(RichTextBlock{Elements: []RichTextElement{RichTextQuote{Border: func() *int { v := 2; return &v }(), Elements: []RichTextInline{TextInline{Text: "x"}}}}})},
	{ruleRichSectionRequired, "$.blocks[0].elements[0].elements", messageOf(RichTextBlock{Elements: []RichTextElement{RichTextSection{}}})},
	{ruleInlineTextRequired, "$.blocks[0].elements[0].elements[0].text", rich(TextInline{})},
	{ruleLinkURLFormat, "$.blocks[0].elements[0].elements[0].url", rich(LinkInline{URL: "relative"})},
	{ruleLinkStyleCode, "$.blocks[0].elements[0].elements[0].style.code", rich(LinkInline{URL: "https://example.com", Style: &RichTextStyle{Code: true}})},
	{ruleUserIDRequired, "$.blocks[0].elements[0].elements[0].user_id", rich(UserInline{})},
	{ruleUserStyleCode, "$.blocks[0].elements[0].elements[0].style.code", rich(UserInline{UserID: "U1", Style: &RichTextStyle{Code: true}})},
	{ruleEmojiNameRequired, "$.blocks[0].elements[0].elements[0].name", rich(EmojiInline{})},
}

func TestValidMessageHasNoViolations(t *testing.T) {
	b := button()
	b.URL, b.Style, b.Confirm = "https://example.com", ButtonStyleDanger, dialog()
	b.AgentPrompt, b.AgentPromptDisplay, b.VisibleToUserIDs = "Summarize", "Summary", []string{"U1"}
	title := plain("Title")
	message := messageOf(
		HeaderBlock{Text: plain("Notice"), Level: 2, BlockID: "header"},
		SectionBlock{Text: MrkdwnTextObject{Text: "*Body*"}, Fields: []TextObject{plain("A"), plain("B")}, Accessory: ImageElement{ImageURL: "https://example.com/a.png", AltText: "a"}},
		ImageBlock{SlackFile: &SlackFileObject{ID: "F1"}, AltText: "file", Title: &title},
		ActionsBlock{Elements: []ActionElement{b}},
		ContextBlock{Elements: []ContextElement{plain("by"), ImageElement{ImageURL: "https://example.com/a.png", AltText: "a"}}},
		DividerBlock{},
		MarkdownBlock{Text: "**done**"},
		video(),
		RichTextBlock{Elements: []RichTextElement{
			RichTextSection{Elements: []RichTextInline{TextInline{Text: "Hi "}, UserInline{UserID: "U1"}, EmojiInline{Name: "wave"}}},
			RichTextList{Style: RichTextListStyleOrdered, Offset: 2, Border: one(), Elements: []RichTextSection{{Elements: []RichTextInline{LinkInline{URL: "https://example.com"}}}}},
			RichTextPreformatted{Elements: []PreformattedInline{TextInline{Text: "code"}}},
			RichTextQuote{Elements: []RichTextInline{TextInline{Text: "quote"}}},
		}},
	)
	require.Empty(t, violations(message))
	require.NoError(t, Validate(message, conversion.WithWarningAsError()))
}

func TestRules(t *testing.T) {
	for _, test := range ruleCases {
		t.Run(test.rule.ID, func(t *testing.T) {
			require.Contains(t, violations(test.message), violation{test.path, test.rule})
		})
	}
}

func TestEveryRuleIsTestedAndDocumented(t *testing.T) {
	readme, err := os.ReadFile("README.md")
	require.NoError(t, err)
	ids := map[string]bool{}
	for _, rule := range rules {
		require.False(t, ids[rule.ID], "duplicate rule %s", rule.ID)
		ids[rule.ID] = true
		require.True(t, slices.ContainsFunc(ruleCases, func(c ruleCase) bool { return c.rule == rule }), "untested rule %s", rule.ID)
		require.Contains(t, string(readme), "`"+rule.ID+"`", "undocumented rule")
		require.NotEmpty(t, rule.Reference)
	}
}

func TestValidateSeverities(t *testing.T) {
	for _, test := range ruleCases {
		t.Run(test.rule.ID, func(t *testing.T) {
			var reported []conversion.Diagnostic
			err := Validate(test.message, conversion.WithDiagnostics(func(d conversion.Diagnostic) { reported = append(reported, d) }))
			warningAsError := Validate(test.message, conversion.WithWarningAsError())
			want := conversion.Diagnostic{Platform: "slack", Path: test.path, Code: test.rule.ID, Message: test.rule.Message, Severity: test.rule.Severity}
			var diagnostic conversion.Diagnostic
			switch test.rule.Severity {
			case validation.Fatal:
				require.True(t, errors.As(err, &diagnostic))
			case validation.Warning:
				require.True(t, errors.As(warningAsError, &diagnostic))
				require.Contains(t, reported, want)
			case validation.Advisory:
				require.Contains(t, reported, want)
			}
		})
	}
}

func TestRegisterValidationReportsFatalRules(t *testing.T) {
	v := validator.New(validator.WithRequiredStructEnabled())
	RegisterValidation(v)
	for _, test := range ruleCases {
		t.Run(test.rule.ID, func(t *testing.T) {
			var want []string
			for _, found := range violations(test.message) {
				if found.rule.Severity == validation.Fatal {
					want = append(want, found.rule.ID)
				}
			}
			var got []string
			var failures validator.ValidationErrors
			if errors.As(v.Struct(test.message), &failures) {
				for _, failure := range failures {
					got = append(got, failure.Tag())
				}
			}
			require.ElementsMatch(t, want, got)
		})
	}
}

func TestConversionMatchesValidate(t *testing.T) {
	type user struct {
		ID string `slack:"part"`
	}
	type item struct {
		Text string `slack:"part"`
	}
	source := struct {
		Title string `slack:"header"`
		Image struct {
			URL string `slack:"part;slot=url"`
			Alt string `slack:"part;slot=alt"`
		} `slack:"image"`
		Actions struct {
			Button struct {
				Text string `slack:"part"`
				URL  string `slack:"part;slot=url"`
			} `slack:"button"`
		} `slack:"actions"`
		Body struct {
			Prefix string `slack:"part"`
			User   user   `slack:"user"`
			List   struct {
				Style  string `slack:"part;slot=style"`
				Offset int    `slack:"part;slot=offset"`
				Items  []struct {
					Item item `slack:"text"`
				} `slack:"rich_text_section"`
			} `slack:"rich_text_list"`
		} `slack:"rich_text"`
		Markdown string `slack:"markdown"`
	}{Title: "Notice", Markdown: ""}
	source.Image.URL, source.Image.Alt = "/a.png", ""
	source.Actions.Button.Text, source.Actions.Button.URL = strings.Repeat("x", 31), "relative"
	source.Body.List.Style, source.Body.List.Offset = "bullet", 3
	source.Body.List.Items = make([]struct {
		Item item `slack:"text"`
	}, 1)
	source.Body.List.Items[0].Item.Text = "first"

	var converted []string
	message, err := ToMessage(source, conversion.WithDiagnostics(func(d conversion.Diagnostic) {
		converted = append(converted, d.Code)
	}))
	require.NoError(t, err)
	var native []string
	for _, found := range violations(message) {
		native = append(native, found.rule.ID)
	}
	require.ElementsMatch(t, native, converted)
	require.Subset(t, converted, []string{
		ruleImageURLFormat.ID, ruleImageAltRequired.ID, ruleButtonTextTruncation.ID, ruleButtonURLFormat.ID,
		ruleInlineTextRequired.ID, ruleUserIDRequired.ID, ruleListOffsetOrdered.ID, ruleMarkdownTextRequired.ID,
		ruleMessageTextFallback.ID,
	})

	_, err = ToMessage(source, conversion.WithWarningAsError())
	var diagnostic conversion.Diagnostic
	require.True(t, errors.As(err, &diagnostic))
	require.Equal(t, validation.Warning, diagnostic.Severity)
	require.Equal(t, "$.Image.URL", diagnostic.Path)
}
