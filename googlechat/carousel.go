package googlechat

import "strings"

// Carousel displays a collection of cards with previous and next navigation.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Carousel
type Carousel struct {
	// CarouselCards contains the cards displayed by the carousel.
	CarouselCards []CarouselCard `json:"carouselCards" validate:"dive"`
}

func (Carousel) WidgetType() string { return "carousel" }
func (Carousel) widgetContent()     {}
func (c Carousel) String() string {
	var parts []string
	for _, card := range c.CarouselCards {
		for _, widget := range card.Widgets {
			if widget.Content != nil {
				parts = append(parts, widget.Content.String())
			}
		}
		for _, widget := range card.FooterWidgets {
			if widget.Content != nil {
				parts = append(parts, widget.Content.String())
			}
		}
	}
	return strings.Join(parts, "\n")
}

// CarouselCard groups widgets and optional footer widgets in a carousel item.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#CarouselCard
type CarouselCard struct {
	// Widgets contains the card's content in display order.
	Widgets []NestedWidget `json:"widgets" validate:"dive"`
	// FooterWidgets contains content shown at the bottom of the card.
	FooterWidgets []NestedWidget `json:"footerWidgets,omitempty" validate:"dive"`
}

// NestedWidgetContent is a paragraph, image, or button list accepted by a CarouselCard.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#NestedWidget
type NestedWidgetContent interface {
	WidgetContent
	nestedWidgetContent()
}

// NestedWidget wraps one content object in a containing layout such as a carousel card.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#NestedWidget
type NestedWidget struct {
	// Content contains a paragraph, image, or button list.
	Content NestedWidgetContent `json:"-"`
}

func (w NestedWidget) MarshalJSON() ([]byte, error) { return marshalWidget(w.Content, "") }
