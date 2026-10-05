package googlechat

import (
	"github.com/JSYoo5B/convertago/internal/validation"

	"encoding/json"
	"fmt"
)

// WidgetContent is the content of a card widget.
// WidgetType returns the JSON property identifying that content, and String returns its text.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Widget
type WidgetContent interface {
	WidgetType() string
	String() string
	widgetContent()
	check(*validation.Check)
}

// Widget displays one content object with optional horizontal alignment.
// MarshalJSON places Content beneath the property returned by WidgetType.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Widget
type Widget struct {
	// Content contains exactly one supported widget content object.
	Content WidgetContent `json:"-"`
	// HorizontalAlignment positions the widget at the start, center, or end.
	HorizontalAlignment HorizontalAlignment `json:"-"`
}

func (w Widget) Type() string {
	if w.Content == nil {
		return ""
	}
	return w.Content.WidgetType()
}
func (w Widget) String() string {
	if w.Content == nil {
		return ""
	}
	return w.Content.String()
}
func (w Widget) MarshalJSON() ([]byte, error) {
	return marshalWidget(w.Content, w.HorizontalAlignment)
}

func marshalWidget(content WidgetContent, alignment HorizontalAlignment) ([]byte, error) {
	if content == nil {
		return nil, fmt.Errorf("googlechat: widget content is required")
	}
	payload, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	if string(payload) == "null" {
		return nil, fmt.Errorf("googlechat: widget content must not be null")
	}
	object := map[string]any{content.WidgetType(): json.RawMessage(payload)}
	if alignment != "" {
		object["horizontalAlignment"] = alignment
	}
	return json.Marshal(object)
}

// HorizontalAlignment selects the widget's horizontal position.
type HorizontalAlignment string

const (
	HorizontalAlignmentStart  HorizontalAlignment = "START"
	HorizontalAlignmentCenter HorizontalAlignment = "CENTER"
	HorizontalAlignmentEnd    HorizontalAlignment = "END"
)

// VerticalAlignment positions content vertically within its containing layout.
type VerticalAlignment string

const (
	VerticalAlignmentTop    VerticalAlignment = "TOP"
	VerticalAlignmentMiddle VerticalAlignment = "MIDDLE"
	VerticalAlignmentBottom VerticalAlignment = "BOTTOM"
)
