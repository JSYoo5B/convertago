package googlechat_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestCardHeader_Validate(t *testing.T) {
	v := validator.New()
	googlechat.RegisterValidation(v)
	require.NoError(t, v.Struct(googlechat.CardHeader{Title: "Notice"}))
	require.NoError(t, v.Struct(googlechat.CardHeader{Title: "Notice", ImageURL: "https://example.com/avatar.png", ImageType: googlechat.ImageTypeCircle}))
	require.Error(t, v.Struct(googlechat.CardHeader{}))
	require.Error(t, v.Struct(googlechat.CardHeader{Title: "Notice", ImageType: googlechat.ImageTypeCircle}))
	require.Error(t, v.Struct(googlechat.CardHeader{Title: "Notice", ImageURL: "http://example.com/avatar.png"}))
}

func TestSection_Validate(t *testing.T) {
	v := validator.New()
	googlechat.RegisterValidation(v)
	section := googlechat.Section{Widgets: []googlechat.Widget{{Content: googlechat.TextParagraph{Text: "Notice"}}}}
	require.NoError(t, v.Struct(section))
	section.UncollapsibleWidgetsCount = 1
	require.Error(t, v.Struct(section))
	section.Collapsible = true
	require.NoError(t, v.Struct(section))
	section.UncollapsibleWidgetsCount = 2
	require.Error(t, v.Struct(section))
	section.UncollapsibleWidgetsCount = 0
	section.CollapseControl = &googlechat.CollapseControl{}
	require.Error(t, v.Struct(section))
	require.Error(t, v.Struct(googlechat.Section{}))
}

func TestCard_ValidateWidgetCount(t *testing.T) {
	v := validator.New()
	googlechat.RegisterValidation(v)
	widgets := make([]googlechat.Widget, 51)
	for i := range widgets {
		widgets[i] = googlechat.Widget{Content: &googlechat.TextParagraph{Text: "Line"}}
	}
	card := googlechat.Card{Sections: []googlechat.Section{
		{Widgets: widgets[:50]}, {Widgets: widgets[:50]},
	}}
	require.NoError(t, v.Struct(card))
	card.Sections[1].Widgets = widgets
	err := v.Struct(card)
	require.ErrorAs(t, err, new(validator.ValidationErrors))
	require.Equal(t, "max_widgets", err.(validator.ValidationErrors)[0].Tag())
	require.Error(t, v.Struct(googlechat.Card{Sections: []googlechat.Section{}}))
}

func TestCard_ValidateNestedWidgetCount(t *testing.T) {
	v := validator.New()
	googlechat.RegisterValidation(v)
	columnWidgets := make([]googlechat.ColumnWidget, 100)
	for i := range columnWidgets {
		columnWidgets[i] = googlechat.ColumnWidget{Content: &googlechat.TextParagraph{Text: "Line"}}
	}
	columns := googlechat.Columns{ColumnItems: []googlechat.Column{{Widgets: columnWidgets[:99]}}}
	card := googlechat.Card{Sections: []googlechat.Section{{Widgets: []googlechat.Widget{{Content: &columns}}}}}
	require.NoError(t, v.Struct(card))
	columns.ColumnItems[0].Widgets = columnWidgets
	require.Error(t, v.Struct(card))

	nested := make([]googlechat.NestedWidget, 99)
	for i := range nested {
		nested[i] = googlechat.NestedWidget{Content: googlechat.TextParagraph{Text: "Line"}}
	}
	carousel := googlechat.Carousel{CarouselCards: []googlechat.CarouselCard{{
		Widgets: nested[:98], FooterWidgets: nested[:1],
	}}}
	card.Sections[0].Widgets[0].Content = &carousel
	require.NoError(t, v.Struct(card))
	carousel.CarouselCards[0].FooterWidgets = nested[:2]
	require.Error(t, v.Struct(card))
}
