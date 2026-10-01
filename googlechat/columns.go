package googlechat

import "strings"

// Columns displays up to two columns, with the second wrapping on narrow screens.
// Widgets within each column are displayed vertically in the order supplied.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Columns
type Columns struct {
	// ColumnItems contains up to two columns.
	ColumnItems []Column `json:"columnItems" validate:"min=1,max=2,dive"`
}

func (Columns) WidgetType() string { return "columns" }
func (Columns) widgetContent()     {}
func (c Columns) String() string {
	var parts []string
	for _, column := range c.ColumnItems {
		for _, widget := range column.Widgets {
			if widget.Content != nil {
				parts = append(parts, widget.Content.String())
			}
		}
	}
	return strings.Join(parts, "\n")
}

// Column positions an ordered collection of widgets within a Columns layout.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Column
type Column struct {
	// HorizontalSizeStyle controls how much available width the column occupies.
	HorizontalSizeStyle HorizontalSizeStyle `json:"horizontalSizeStyle,omitempty" validate:"omitempty,oneof=FILL_AVAILABLE_SPACE FILL_MINIMUM_SPACE"`
	// HorizontalAlignment positions the column's widgets horizontally.
	HorizontalAlignment HorizontalAlignment `json:"horizontalAlignment,omitempty" validate:"omitempty,oneof=START CENTER END"`
	// VerticalAlignment positions the widgets vertically. Omission centers them.
	VerticalAlignment ColumnVerticalAlignment `json:"verticalAlignment,omitempty" validate:"omitempty,oneof=CENTER TOP BOTTOM"`
	// Widgets contains the supported column widgets in display order.
	Widgets []ColumnWidget `json:"widgets" validate:"min=1,max=100,dive"`
}

// ColumnWidgetContent is content supported within a column.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Widgets
type ColumnWidgetContent interface {
	WidgetContent
	columnWidgetContent()
}

// ColumnWidget wraps one content object accepted by a Column.
// Column widget wrappers do not carry the alignment field of a top-level Widget.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Widgets
type ColumnWidget struct {
	// Content contains a paragraph, image, decorated text, button list, or chip list.
	Content ColumnWidgetContent `json:"-" validate:"required"`
}

func (w ColumnWidget) MarshalJSON() ([]byte, error) { return marshalWidget(w.Content, "") }

// HorizontalSizeStyle controls the width allocated to a column.
type HorizontalSizeStyle string

const (
	// HorizontalSizeStyleFillAvailableSpace fills up to 70% of the card, or 50% when both columns use it.
	HorizontalSizeStyleFillAvailableSpace HorizontalSizeStyle = "FILL_AVAILABLE_SPACE"
	// HorizontalSizeStyleFillMinimumSpace uses the smallest width possible, up to 30% of the card.
	HorizontalSizeStyleFillMinimumSpace HorizontalSizeStyle = "FILL_MINIMUM_SPACE"
)

// ColumnVerticalAlignment positions widgets vertically within a column.
type ColumnVerticalAlignment string

const (
	ColumnVerticalAlignmentCenter ColumnVerticalAlignment = "CENTER"
	ColumnVerticalAlignmentTop    ColumnVerticalAlignment = "TOP"
	ColumnVerticalAlignmentBottom ColumnVerticalAlignment = "BOTTOM"
)
