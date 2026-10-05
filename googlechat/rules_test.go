package googlechat

import (
	"errors"
	"math"
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

const (
	card0   = "$.cardsV2[0].card"
	widget0 = card0 + ".sections[0].widgets[0]"
	button0 = widget0 + ".buttonList.buttons[0]"
)

func cardOf(card Card) Message { return Message{CardsV2: []CardWithID{{Card: card}}} }

func widgets(contents ...WidgetContent) Message {
	section := Section{}
	for _, content := range contents {
		section.Widgets = append(section.Widgets, Widget{Content: content})
	}
	return cardOf(Card{Sections: []Section{section}})
}

func link() OnClick { return OnClick{OpenLink: &OpenLink{URL: "https://example.com"}} }

func buttonWith(edit func(*Button)) Message {
	b := Button{Text: "Open", OnClick: link()}
	edit(&b)
	return widgets(ButtonList{Buttons: []Button{b}})
}

func actionWith(edit func(*Action)) Message {
	return buttonWith(func(b *Button) {
		action := Action{Function: "save"}
		edit(&action)
		b.OnClick = OnClick{Action: &action}
	})
}

func gridItem(item GridItem) Message { return widgets(Grid{Items: []GridItem{item}}) }

func imageItem(image ImageComponent) Message {
	return gridItem(GridItem{Title: "Item", Image: &image})
}

type ruleCase struct {
	rule    validation.Rule
	path    string
	message Message
}

var ruleCases = []ruleCase{
	{ruleMessageCardsSize, "$.cardsV2", widgets(TextParagraph{Text: strings.Repeat("x", 33000)})},
	{ruleCardIDRequired, "$.cardsV2[0].cardId", Message{CardsV2: []CardWithID{{Card: Card{Header: &CardHeader{Title: "A"}}}, {CardID: "b", Card: Card{Header: &CardHeader{Title: "B"}}}}}},
	{ruleCardIDUnique, "$.cardsV2[1].cardId", Message{CardsV2: []CardWithID{{CardID: "a", Card: Card{Header: &CardHeader{Title: "A"}}}, {CardID: "a", Card: Card{Header: &CardHeader{Title: "B"}}}}}},
	{ruleCardContent, card0, cardOf(Card{})},
	{ruleCardWidgetCount, card0 + ".sections", widgets(slices.Repeat([]WidgetContent{Divider{}}, 101)...)},
	{ruleCardDividerStyle, card0 + ".sectionDividerStyle", cardOf(Card{Header: &CardHeader{Title: "A"}, SectionDividerStyle: "DOTTED"})},
	{ruleHeaderTitle, card0 + ".header.title", cardOf(Card{Header: &CardHeader{}})},
	{ruleHeaderImageHTTPS, card0 + ".header.imageUrl", cardOf(Card{Header: &CardHeader{Title: "A", ImageURL: "http://example.com/a.png"}})},
	{ruleHeaderImageURL, card0 + ".header.imageUrl", cardOf(Card{Header: &CardHeader{Title: "A", ImageType: ImageTypeCircle}})},
	{ruleImageTypeValue, card0 + ".header.imageType", cardOf(Card{Header: &CardHeader{Title: "A", ImageURL: "https://example.com/a.png", ImageType: "OVAL"}})},
	{ruleSectionWidgets, card0 + ".sections[0].widgets", cardOf(Card{Sections: []Section{{}}})},
	{ruleSectionCollapse, card0 + ".sections[0].collapsible", cardOf(Card{Sections: []Section{{Widgets: []Widget{{Content: Divider{}}}, UncollapsibleWidgetsCount: 1}}})},
	{ruleSectionUncollapse, card0 + ".sections[0].uncollapsibleWidgetsCount", cardOf(Card{Sections: []Section{{Widgets: []Widget{{Content: Divider{}}}, Collapsible: true, UncollapsibleWidgetsCount: 2}}})},
	{ruleWidgetContent, widget0, cardOf(Card{Sections: []Section{{Widgets: []Widget{{}}}}})},
	{ruleHorizontalAlign, widget0 + ".horizontalAlignment", cardOf(Card{Sections: []Section{{Widgets: []Widget{{Content: Divider{}, HorizontalAlignment: "LEFT"}}}}})},
	{ruleVerticalAlign, widget0 + ".decoratedText.startIconVerticalAlignment", widgets(DecoratedText{Text: "x", StartIconVerticalAlignment: "CENTER"})},

	{ruleParagraphMaxLines, widget0 + ".textParagraph.maxLines", widgets(TextParagraph{Text: "x", MaxLines: -1})},
	{ruleParagraphSyntax, widget0 + ".textParagraph.textSyntax", widgets(TextParagraph{Text: "x", TextSyntax: "PLAIN"})},
	{ruleImageURLHTTPS, widget0 + ".image.imageUrl", widgets(Image{ImageURL: "http://example.com/a.png"})},
	{ruleDecoratedText, widget0 + ".decoratedText.text", widgets(DecoratedText{})},
	{ruleDecoratedControl, widget0 + ".decoratedText", widgets(DecoratedText{Text: "x", Button: &Button{Text: "Open", OnClick: link()}, EndIcon: &Icon{KnownIcon: "STAR"}})},
	{ruleSwitchName, widget0 + ".decoratedText.switchControl.name", widgets(DecoratedText{Text: "x", SwitchControl: &SwitchControl{}})},
	{ruleSwitchType, widget0 + ".decoratedText.switchControl.controlType", widgets(DecoratedText{Text: "x", SwitchControl: &SwitchControl{Name: "on", ControlType: "RADIO"}})},

	{ruleButtonListButtons, widget0 + ".buttonList.buttons", widgets(ButtonList{})},
	{ruleButtonContent, button0, buttonWith(func(b *Button) { b.Text = "" })},
	{ruleButtonTypeValue, button0 + ".type", buttonWith(func(b *Button) { b.Type = "RAISED" })},
	{ruleButtonTypeColor, button0 + ".type", buttonWith(func(b *Button) { b.Type, b.Color = ButtonTypeOutlined, &Color{Red: 1} })},
	{ruleColorRange, button0 + ".color.red", buttonWith(func(b *Button) { b.Color = &Color{Red: 2} })},
	{ruleIconSource, button0 + ".icon", buttonWith(func(b *Button) { b.Icon = &Icon{} })},
	{ruleIconURLHTTPS, button0 + ".icon.iconUrl", buttonWith(func(b *Button) { b.Icon = &Icon{IconURL: "http://example.com/a.png"} })},
	{ruleMaterialName, button0 + ".icon.materialIcon.name", buttonWith(func(b *Button) { b.Icon = &Icon{MaterialIcon: &MaterialIcon{}} })},
	{ruleMaterialWeight, button0 + ".icon.materialIcon.weight", buttonWith(func(b *Button) { b.Icon = &Icon{MaterialIcon: &MaterialIcon{Name: "home", Weight: 150}} })},
	{ruleMaterialGrade, button0 + ".icon.materialIcon.grade", buttonWith(func(b *Button) { b.Icon = &Icon{MaterialIcon: &MaterialIcon{Name: "home", Grade: 5}} })},

	{ruleOnClickCount, button0 + ".onClick", buttonWith(func(b *Button) { b.OnClick = OnClick{} })},
	{ruleActionFunction, button0 + ".onClick.action.function", actionWith(func(a *Action) { a.Function = "" })},
	{ruleActionLoad, button0 + ".onClick.action.loadIndicator", actionWith(func(a *Action) { a.LoadIndicator = "BAR" })},
	{ruleActionInteraction, button0 + ".onClick.action.interaction", actionWith(func(a *Action) { a.Interaction = "OPEN_PAGE" })},
	{ruleActionRequired, button0 + ".onClick.action.requiredWidgets", actionWith(func(a *Action) { a.AllWidgetsAreRequired, a.RequiredWidgets = true, []string{"name"} })},
	{ruleActionRequiredName, button0 + ".onClick.action.requiredWidgets[0]", actionWith(func(a *Action) { a.RequiredWidgets = []string{""} })},
	{ruleParameterKey, button0 + ".onClick.action.parameters[0].key", actionWith(func(a *Action) { a.Parameters = []ActionParameter{{Value: "1"}} })},
	{ruleParameterUnique, button0 + ".onClick.action.parameters[1].key", actionWith(func(a *Action) { a.Parameters = []ActionParameter{{Key: "id"}, {Key: "id"}} })},
	{ruleOpenLinkURL, button0 + ".onClick.openLink.url", buttonWith(func(b *Button) { b.OnClick = OnClick{OpenLink: &OpenLink{URL: "relative"}} })},
	{ruleOverflowItems, button0 + ".onClick.overflowMenu.items", buttonWith(func(b *Button) { b.OnClick = OnClick{OverflowMenu: &OverflowMenu{}} })},
	{ruleOverflowItemText, button0 + ".onClick.overflowMenu.items[0].text", buttonWith(func(b *Button) {
		b.OnClick = OnClick{OverflowMenu: &OverflowMenu{Items: []OverflowMenuItem{{OnClick: link()}}}}
	})},
	{ruleOverflowNested, button0 + ".onClick.overflowMenu.items[0].onClick.overflowMenu", buttonWith(func(b *Button) {
		nested := OnClick{OverflowMenu: &OverflowMenu{Items: []OverflowMenuItem{{Text: "Inner", OnClick: link()}}}}
		b.OnClick = OnClick{OverflowMenu: &OverflowMenu{Items: []OverflowMenuItem{{Text: "Outer", OnClick: nested}}}}
	})},

	{ruleColumnsCount, widget0 + ".columns.columnItems", widgets(Columns{})},
	{ruleColumnWidgets, widget0 + ".columns.columnItems[0].widgets", widgets(Columns{ColumnItems: []Column{{}}})},
	{ruleColumnSizeStyle, widget0 + ".columns.columnItems[0].horizontalSizeStyle", widgets(Columns{ColumnItems: []Column{{HorizontalSizeStyle: "WIDE", Widgets: []ColumnWidget{{Content: TextParagraph{Text: "x"}}}}}})},
	{ruleColumnVertical, widget0 + ".columns.columnItems[0].verticalAlignment", widgets(Columns{ColumnItems: []Column{{VerticalAlignment: "MIDDLE", Widgets: []ColumnWidget{{Content: TextParagraph{Text: "x"}}}}}})},
	{ruleGridItems, widget0 + ".grid.items", widgets(Grid{})},
	{ruleGridColumnCount, widget0 + ".grid.columnCount", widgets(Grid{ColumnCount: -1, Items: []GridItem{{Title: "x"}}})},
	{ruleGridItemContent, widget0 + ".grid.items[0]", gridItem(GridItem{})},
	{ruleGridItemLayout, widget0 + ".grid.items[0].layout", gridItem(GridItem{Title: "x", Layout: "TEXT_LEFT"})},
	{ruleImageURIFormat, widget0 + ".grid.items[0].image.imageUri", imageItem(ImageComponent{ImageURI: "relative"})},
	{ruleCropTypeValue, widget0 + ".grid.items[0].image.cropStyle.type", imageItem(ImageComponent{ImageURI: "https://example.com/a.png", CropStyle: &ImageCropStyle{Type: "OVAL"}})},
	{ruleCropRatioValue, widget0 + ".grid.items[0].image.cropStyle.aspectRatio", imageItem(ImageComponent{ImageURI: "https://example.com/a.png", CropStyle: &ImageCropStyle{Type: ImageCropTypeRectangleCustom}})},
	{ruleCropRatioCustom, widget0 + ".grid.items[0].image.cropStyle.aspectRatio", imageItem(ImageComponent{ImageURI: "https://example.com/a.png", CropStyle: &ImageCropStyle{Type: ImageCropTypeSquare, AspectRatio: 1.5}})},
	{ruleBorderTypeValue, widget0 + ".grid.borderStyle.type", widgets(Grid{Items: []GridItem{{Title: "x"}}, BorderStyle: &BorderStyle{Type: "DASHED"}})},
	{ruleBorderStroke, widget0 + ".grid.borderStyle.strokeColor", widgets(Grid{Items: []GridItem{{Title: "x"}}, BorderStyle: &BorderStyle{Type: BorderTypeNone, StrokeColor: &Color{}}})},
	{ruleBorderRadius, widget0 + ".grid.borderStyle.cornerRadius", widgets(Grid{Items: []GridItem{{Title: "x"}}, BorderStyle: &BorderStyle{CornerRadius: -1}})},
	{ruleCarouselCards, widget0 + ".carousel.carouselCards", widgets(Carousel{})},
	{ruleCarouselWidgets, widget0 + ".carousel.carouselCards[0].widgets", widgets(Carousel{CarouselCards: []CarouselCard{{}}})},
	{ruleChipListChips, widget0 + ".chipList.chips", widgets(ChipList{})},
	{ruleChipListLayout, widget0 + ".chipList.layout", widgets(ChipList{Layout: "STACKED", Chips: []Chip{{Label: "x"}}})},
	{ruleChipContent, widget0 + ".chipList.chips[0]", widgets(ChipList{Chips: []Chip{{}}})},
}

func TestValidMessageHasNoViolations(t *testing.T) {
	action := &Action{Function: "save", Parameters: []ActionParameter{{Key: "id", Value: "1"}}, LoadIndicator: LoadIndicatorNone, RequiredWidgets: []string{"name"}}
	message := Message{CardsV2: []CardWithID{
		{CardID: "first", Card: Card{
			Header:              &CardHeader{Title: "Notice", ImageURL: "https://example.com/a.png", ImageType: ImageTypeCircle, ImageAltText: "a"},
			SectionDividerStyle: DividerStyleSolid,
			Sections: []Section{{
				Collapsible:               true,
				UncollapsibleWidgetsCount: 1,
				CollapseControl:           &CollapseControl{ExpandButton: Button{Text: "More", OnClick: link()}, CollapseButton: Button{Text: "Less", OnClick: link()}},
				Widgets: []Widget{
					{Content: TextParagraph{Text: "**x**", TextSyntax: TextSyntaxMarkdown, MaxLines: 2}, HorizontalAlignment: HorizontalAlignmentCenter},
					{Content: Image{ImageURL: "https://example.com/a.png", OnClick: &OnClick{Action: action}}},
					{Content: Divider{}},
					{Content: DecoratedText{Text: "x", StartIcon: &Icon{MaterialIcon: &MaterialIcon{Name: "home", Weight: 400, Grade: -25}}, SwitchControl: &SwitchControl{Name: "on", ControlType: SwitchControlTypeCheckBox}}},
					{Content: ButtonList{Buttons: []Button{
						{Text: "Fill", Color: &Color{Red: 1}, Type: ButtonTypeFilled, OnClick: OnClick{OverflowMenu: &OverflowMenu{Items: []OverflowMenuItem{{Text: "Item", OnClick: link()}}}}},
						{Icon: &Icon{IconURL: "https://example.com/i.png"}, OnClick: link()},
					}}},
				},
			}},
		}},
		{CardID: "second", Card: Card{Sections: []Section{{Widgets: []Widget{
			{Content: Columns{ColumnItems: []Column{{HorizontalSizeStyle: HorizontalSizeStyleFillMinimumSpace, Widgets: []ColumnWidget{{Content: ChipList{Chips: []Chip{{Label: "chip"}}}}}}}}},
			{Content: Grid{ColumnCount: 2, Items: []GridItem{{Image: &ImageComponent{ImageURI: "https://example.com/a.png", CropStyle: &ImageCropStyle{Type: ImageCropTypeRectangleCustom, AspectRatio: 1.5}, BorderStyle: &BorderStyle{Type: BorderTypeStroke, StrokeColor: &Color{Blue: 1}, CornerRadius: 4}}}}}},
			{Content: Carousel{CarouselCards: []CarouselCard{{Widgets: []NestedWidget{{Content: TextParagraph{Text: "x"}}}}}}},
		}}}}},
	}}
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
			case validation.Warning:
				require.True(t, errors.As(warningAsError, &diagnostic))
				require.Contains(t, reported, conversion.Diagnostic{Platform: "googlechat", Path: test.path, Code: test.rule.ID, Message: test.rule.Message, Severity: validation.Warning})
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

func TestNonFiniteAspectRatio(t *testing.T) {
	message := imageItem(ImageComponent{ImageURI: "https://example.com/a.png", CropStyle: &ImageCropStyle{Type: ImageCropTypeRectangleCustom, AspectRatio: math.Inf(1)}})
	require.Contains(t, violations(message), violation{widget0 + ".grid.items[0].image.cropStyle.aspectRatio", ruleCropRatioValue})
}

func TestConversionMatchesValidate(t *testing.T) {
	type click struct {
		URL string `googlechat:"openLink"`
	}
	type button struct {
		Text  string `googlechat:"part"`
		Type  string `googlechat:"part;slot=type"`
		Color struct {
			Red float64 `googlechat:"part;slot=red"`
		} `googlechat:"color"`
		Click click `googlechat:"onClick"`
	}
	source := struct {
		Title   string `googlechat:"header"`
		Body    string `googlechat:"textParagraph"`
		Image   string `googlechat:"image"`
		Buttons struct {
			Save button `googlechat:"button"`
		} `googlechat:"buttonList"`
		Decorated struct {
			Text string `googlechat:"part"`
			Icon struct {
				Material struct {
					Name   string `googlechat:"part"`
					Weight int    `googlechat:"part;slot=weight"`
				} `googlechat:"materialIcon"`
			} `googlechat:"icon;slot=startIcon"`
		} `googlechat:"decoratedText"`
	}{Title: "", Body: "x", Image: "http://example.com/a.png"}
	source.Buttons.Save.Text, source.Buttons.Save.Type = "Save", "OUTLINED"
	source.Buttons.Save.Color.Red = 1.5
	source.Buttons.Save.Click.URL = "relative"
	source.Decorated.Icon.Material.Name, source.Decorated.Icon.Material.Weight = "home", 150

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
		ruleHeaderTitle.ID, ruleImageURLHTTPS.ID, ruleButtonTypeColor.ID, ruleColorRange.ID,
		ruleOpenLinkURL.ID, ruleMaterialWeight.ID, ruleDecoratedText.ID,
	})

	_, err = ToMessage(source, conversion.WithWarningAsError())
	var diagnostic conversion.Diagnostic
	require.True(t, errors.As(err, &diagnostic))
	require.Equal(t, validation.Warning, diagnostic.Severity)
	require.Equal(t, "$.Title", diagnostic.Path)
}
