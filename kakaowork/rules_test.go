package kakaowork

import (
	"encoding/json"
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

func messageOf(blocks ...BubbleBlock) Message {
	return Message{Preview: "알림", Blocks: blocks}
}

func browser() ButtonAction { return OpenSystemBrowserAction{Value: "https://example.com"} }

func button() ButtonBlock { return ButtonBlock{Text: "열기", Action: browser()} }

func text(inlines ...Inline) TextBlock {
	block := TextBlock{Inlines: inlines}
	for _, inline := range inlines {
		block.Text += inline.String()
	}
	return block
}

type ruleCase struct {
	rule    validation.Rule
	path    string
	message Message
}

var ruleCases = []ruleCase{
	{ruleMessageTextRequired, "$.text", Message{Blocks: []BubbleBlock{TextBlock{Text: "본문"}}}},
	{ruleMessageTextLength, "$.text", Message{Preview: strings.Repeat("가", 10001)}},
	{ruleMessageBlockType, "$.blocks[0]", messageOf(InlineStyled{Text: "본문"})},
	{ruleMessageBlockType, "$.blocks[0]", messageOf(nil)},
	{ruleMessageButtonCount, "$.blocks", messageOf(ActionBlock{Elements: []ButtonBlock{button(), button(), button()}}, button())},
	{ruleMessageImageAfterButton, "$.blocks[1]", messageOf(button(), ImageBlock{Url: "https://example.com/a.png"})},
	{ruleMessageImageConsecutive, "$.blocks[1]", messageOf(ImageBlock{Url: "https://example.com/a.png"}, ImageBlock{Url: "https://example.com/b.png"})},
	{ruleMessageDividerEdge, "$.blocks[0]", messageOf(DividerBlock{}, TextBlock{Text: "본문"})},
	{ruleMessageDividerOnly, "$.blocks", messageOf(DividerBlock{})},

	{ruleHeaderPosition, "$.blocks[1]", messageOf(TextBlock{Text: "본문"}, HeaderBlock{Text: "알림"})},
	{ruleHeaderTextRequired, "$.blocks[0].text", messageOf(HeaderBlock{})},
	{ruleHeaderTextLength, "$.blocks[0].text", messageOf(HeaderBlock{Text: strings.Repeat("가", 21)})},
	{ruleHeaderTextLineBreak, "$.blocks[0].text", messageOf(HeaderBlock{Text: "두\n줄"})},
	{ruleHeaderStyleValue, "$.blocks[0].style", messageOf(HeaderBlock{Text: "알림", Style: "green"})},

	{ruleTextTextRequired, "$.blocks[0].text", messageOf(TextBlock{})},
	{ruleTextTextLength, "$.blocks[0].text", messageOf(TextBlock{Text: strings.Repeat("가", 501)})},
	{ruleTextInlinesMismatch, "$.blocks[0].inlines", messageOf(TextBlock{Text: "가", Inlines: []Inline{InlineStyled{Text: "나"}}})},
	{ruleTextInlineRequired, "$.blocks[0].inlines[0]", messageOf(TextBlock{Text: "가", Inlines: []Inline{nil}})},
	{ruleTextLinkCount, "$.blocks[0].inlines", messageOf(text(InlineLink{Text: "가", Url: "https://example.com"}, InlineLink{Text: "나", Url: "https://example.com"}))},
	{ruleStyledTextRequired, "$.blocks[0].inlines[1].text", messageOf(text(InlineStyled{Text: "가"}, InlineStyled{}))},
	{ruleStyledColorValue, "$.blocks[0].inlines[0].color", messageOf(text(InlineStyled{Text: "가", Color: "pink"}))},
	{ruleLinkTextRequired, "$.blocks[0].inlines[1].text", messageOf(text(InlineStyled{Text: "가"}, InlineLink{Url: "https://example.com"}))},
	{ruleLinkURLScheme, "$.blocks[0].inlines[0].url", messageOf(text(InlineLink{Text: "가", Url: "ftp://example.com"}))},
	{ruleMentionTextRequired, "$.blocks[0].inlines[1].text", messageOf(text(InlineStyled{Text: "가"}, InlineMention{UserId: 1}))},
	{ruleMentionUserID, "$.blocks[0].inlines[0].ref.value", messageOf(text(InlineMention{Text: "@ryan"}))},

	{ruleImageURLFormat, "$.blocks[0].url", messageOf(ImageBlock{Url: "/relative.png"})},

	{ruleButtonTextRequired, "$.blocks[0].text", messageOf(ButtonBlock{Action: browser()})},
	{ruleButtonTextLength, "$.blocks[0].text", messageOf(ButtonBlock{Text: strings.Repeat("가", 21), Action: browser()})},
	{ruleButtonStyleValue, "$.blocks[0].style", messageOf(ButtonBlock{Text: "열기", Style: "blue", Action: browser()})},
	{ruleButtonActionRequired, "$.blocks[0].action", messageOf(ButtonBlock{Text: "열기"})},
	{ruleActionElementsCount, "$.blocks[0].elements", messageOf(ActionBlock{Elements: []ButtonBlock{button()}})},

	{ruleDescriptionTermRequired, "$.blocks[0].term", messageOf(DescriptionBlock{Content: TextBlock{Text: "본문"}})},
	{ruleDescriptionTermLength, "$.blocks[0].term", messageOf(DescriptionBlock{Content: TextBlock{Text: "본문"}, Term: strings.Repeat("가", 11)})},

	{ruleSystemBrowserURLFormat, "$.blocks[0].action.value", messageOf(ButtonBlock{Text: "열기", Action: OpenSystemBrowserAction{Value: "example.com"}})},
	{ruleInAppBrowserURLFormat, "$.blocks[0].action.value", messageOf(ButtonBlock{Text: "열기", Action: OpenInAppBrowserAction{Value: "example.com"}})},
	{ruleInAppBrowserSizeStandalone, "$.blocks[0].action.standalone", messageOf(ButtonBlock{Text: "열기", Action: OpenInAppBrowserAction{Value: "https://example.com", Width: 800}})},
	{ruleInAppBrowserSizeValue, "$.blocks[0].action.height", messageOf(ButtonBlock{Text: "열기", Action: OpenInAppBrowserAction{Value: "https://example.com", Standalone: true, Height: -1}})},
	{ruleExternalAppValueFormat, "$.blocks[0].action.value", messageOf(ButtonBlock{Text: "열기", Action: OpenExternalAppAction{Value: "kakaomap://look"}})},
	{ruleSubmitNameRequired, "$.blocks[0].action.name", messageOf(ButtonBlock{Text: "열기", Action: SubmitAction{Value: "1"}})},
	{ruleSubmitValueRequired, "$.blocks[0].action.value", messageOf(ButtonBlock{Text: "열기", Action: SubmitAction{Name: "accept"}})},
	{ruleCallModalValueRequired, "$.blocks[0].action.value", messageOf(ButtonBlock{Text: "열기", Action: CallModalAction{}})},
	{ruleExclusiveDefaultRequired, "$.blocks[0].action.default", messageOf(ButtonBlock{Text: "열기", Action: ExclusiveAction{Pc: browser()}})},
	{ruleExclusiveActionType, "$.blocks[0].action.ios", messageOf(ButtonBlock{Text: "열기", Action: ExclusiveAction{Default: browser(), Ios: ExclusiveAction{Default: browser()}}})},
}

func TestValidMessageHasNoViolations(t *testing.T) {
	message := messageOf(
		HeaderBlock{Text: "알림", Style: HeaderStyleBlue},
		text(InlineStyled{Text: "안녕하세요 "}, InlineMention{Text: "@ryan", UserId: 7}, InlineLink{Text: " 문서", Url: "mailto:a@example.com"}),
		ImageBlock{Url: "https://example.com/a.png"},
		DividerBlock{},
		DescriptionBlock{Content: TextBlock{Text: "완료"}, Term: "상태"},
		SectionBlock{Content: TextBlock{Text: "본문"}, Accessory: &ImageBlock{Url: "https://example.com/a.png"}, Action: SubmitAction{Name: "accept", Value: "1"}},
		ContextBlock{Content: TextBlock{Text: "작성자"}, Image: ImageBlock{Url: "https://example.com/a.png"}},
		ActionBlock{Elements: []ButtonBlock{
			{Text: "열기", Action: OpenInAppBrowserAction{Value: "https://example.com", Standalone: true, Width: 800, Height: 600}},
			{Text: "앱", Style: ButtonStylePrimary, Action: ExclusiveAction{Default: CallModalAction{Value: "1"}, Ios: OpenExternalAppAction{Value: "ios=kakaomap%3A%2F%2Flook&aos=kakaomap%3A%2F%2Flook"}}},
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
			var diagnostic conversion.Diagnostic
			switch test.rule.Severity {
			case validation.Fatal:
				require.True(t, errors.As(err, &diagnostic))
				require.True(t, errors.As(warningAsError, &diagnostic))
			case validation.Warning:
				require.True(t, errors.As(warningAsError, &diagnostic))
				require.Contains(t, reported, conversion.Diagnostic{Platform: "kakaowork", Path: test.path, Code: test.rule.ID, Message: test.rule.Message, Severity: validation.Warning})
			case validation.Advisory:
				require.Contains(t, reported, conversion.Diagnostic{Platform: "kakaowork", Path: test.path, Code: test.rule.ID, Message: test.rule.Message, Severity: validation.Advisory})
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
	type mention struct {
		Text string `kakaowork:"part"`
		ID   int    `kakaowork:"part;slot=user_id"`
	}
	type action struct {
		Value string `kakaowork:"part"`
		Width int    `kakaowork:"part;slot=width"`
	}
	type external struct {
		Value string `kakaowork:"part"`
	}
	type content struct {
		Text string `kakaowork:"part"`
	}
	source := struct {
		Title string `kakaowork:"header"`
		Body  struct {
			Prefix string  `kakaowork:"part"`
			User   mention `kakaowork:"mention"`
		} `kakaowork:"text"`
		Image   string `kakaowork:"image_link"`
		Buttons struct {
			First struct {
				Text   string `kakaowork:"part"`
				Action action `kakaowork:"open_inapp_browser"`
			} `kakaowork:"button"`
			Second struct {
				Text   string   `kakaowork:"part"`
				Action external `kakaowork:"open_external_app"`
			} `kakaowork:"button"`
		} `kakaowork:"action"`
		Details struct {
			Body content `kakaowork:"text"`
			Term string  `kakaowork:"part;slot=term"`
		} `kakaowork:"description"`
		Divider struct{} `kakaowork:"divider"`
	}{Title: strings.Repeat("가", 21) + "\n", Image: "/relative.png"}
	source.Body.User = mention{"@ryan", 0}
	source.Buttons.First.Text = strings.Repeat("나", 21)
	source.Buttons.First.Action = action{"https://example.com", 800}
	source.Buttons.Second.Text = "앱"
	source.Buttons.Second.Action.Value = "kakaomap://look"
	source.Details.Body.Text = "완료"
	source.Details.Term = strings.Repeat("다", 11)

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
		ruleMessageTextRequired.ID, ruleHeaderTextLength.ID, ruleHeaderTextLineBreak.ID, ruleMentionUserID.ID,
		ruleImageURLFormat.ID, ruleButtonTextLength.ID, ruleInAppBrowserSizeStandalone.ID,
		ruleExternalAppValueFormat.ID, ruleDescriptionTermLength.ID, ruleMessageDividerEdge.ID,
	})

	_, err = ToMessage(source, conversion.WithWarningAsError())
	var diagnostic conversion.Diagnostic
	require.True(t, errors.As(err, &diagnostic))
	require.Equal(t, validation.Warning, diagnostic.Severity)
	require.Equal(t, "$.Title", diagnostic.Path)
}

func TestDefaults(t *testing.T) {
	data, err := json.Marshal(messageOf(HeaderBlock{Text: "알림"}, ButtonBlock{Text: "열기", Action: browser()}))
	require.NoError(t, err)
	require.JSONEq(t, `{"text":"알림","blocks":[
		{"type":"header","text":"알림","style":"white"},
		{"type":"button","text":"열기","style":"default","action":{"type":"open_system_browser","value":"https://example.com"}}
	]}`, string(data))

	source := struct {
		Preview string `kakaowork:"preview"`
		Title   string `kakaowork:"header"`
		Button  struct {
			Text   string `kakaowork:"part"`
			Action string `kakaowork:"open_system_browser"`
		} `kakaowork:"button"`
	}{Preview: "알림", Title: "알림"}
	source.Button.Text, source.Button.Action = "열기", "https://example.com"
	message, err := ToMessage(source)
	require.NoError(t, err)
	converted, err := json.Marshal(message)
	require.NoError(t, err)
	require.JSONEq(t, string(data), string(converted))
}
