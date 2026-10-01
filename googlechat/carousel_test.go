package googlechat_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestCarousel_Validate(t *testing.T) {
	v := validator.New()
	card := googlechat.CarouselCard{Widgets: []googlechat.NestedWidget{{Content: googlechat.TextParagraph{Text: "Line"}}}}
	require.NoError(t, v.Struct(googlechat.Carousel{CarouselCards: []googlechat.CarouselCard{card}}))
	require.Error(t, v.Struct(googlechat.Carousel{}))
	require.Error(t, v.Struct(googlechat.Carousel{CarouselCards: []googlechat.CarouselCard{{}}}))
	card.FooterWidgets = []googlechat.NestedWidget{{}}
	require.Error(t, v.Struct(card))
}
