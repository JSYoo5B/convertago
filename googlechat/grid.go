package googlechat

import "strings"

// Grid arranges text and image items in rows and columns.
// The number of rows depends on the item count and ColumnCount.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Grid
type Grid struct {
	// Title appears in the grid header.
	Title string `json:"title,omitempty"`
	// Items contains the grid entries.
	Items []GridItem `json:"items" validate:"dive"`
	// BorderStyle applies a border to each item.
	BorderStyle *BorderStyle `json:"borderStyle,omitempty"`
	// ColumnCount sets the column count. Omission uses a default depending on the display surface.
	ColumnCount int `json:"columnCount,omitempty"`
	// OnClick is shared by the items, with each item's identifier and index added to parameters.
	OnClick *OnClick `json:"onClick,omitempty"`
}

func (Grid) WidgetType() string { return "grid" }
func (Grid) widgetContent()     {}
func (g Grid) String() string {
	var parts []string
	if g.Title != "" {
		parts = append(parts, g.Title)
	}
	for _, item := range g.Items {
		if item.Title != "" {
			parts = append(parts, item.Title)
		}
		if item.Subtitle != "" {
			parts = append(parts, item.Subtitle)
		}
	}
	return strings.Join(parts, "\n")
}

// GridItem displays text, an image, or both within a grid.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#GridItem
type GridItem struct {
	// ID is returned in the parent grid's click callback parameters.
	ID string `json:"id,omitempty"`
	// Image is the image displayed in the item.
	Image *ImageComponent `json:"image,omitempty"`
	// Title is the item's primary text.
	Title string `json:"title,omitempty"`
	// Subtitle is the item's secondary text.
	Subtitle string `json:"subtitle,omitempty"`
	// Layout places text above or below the image.
	Layout GridItemLayout `json:"layout,omitempty"`
}

// GridItemLayout selects the position of text relative to a grid image.
type GridItemLayout string

const (
	GridItemLayoutTextBelow GridItemLayout = "TEXT_BELOW"
	GridItemLayoutTextAbove GridItemLayout = "TEXT_ABOVE"
)

// ImageComponent configures a layout image with optional cropping and borders.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#ImageComponent
type ImageComponent struct {
	// ImageURI is the image URL.
	ImageURI string `json:"imageUri"`
	// AltText describes the image for accessibility.
	AltText string `json:"altText,omitempty"`
	// CropStyle sets the crop applied to the image.
	CropStyle *ImageCropStyle `json:"cropStyle,omitempty"`
	// BorderStyle sets the image border.
	BorderStyle *BorderStyle `json:"borderStyle,omitempty"`
}

// ImageCropStyle selects an image crop and an optional custom aspect ratio.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#ImageCropStyle
type ImageCropStyle struct {
	// Type selects the crop. Omission uses a square crop.
	Type ImageCropType `json:"type,omitempty"`
	// AspectRatio applies when Type is ImageCropTypeRectangleCustom.
	AspectRatio float64 `json:"aspectRatio,omitempty"`
}

// ImageCropType selects the image's crop shape.
type ImageCropType string

const (
	ImageCropTypeSquare          ImageCropType = "SQUARE"
	ImageCropTypeCircle          ImageCropType = "CIRCLE"
	ImageCropTypeRectangleCustom ImageCropType = "RECTANGLE_CUSTOM"
	ImageCropTypeRectangle4By3   ImageCropType = "RECTANGLE_4_3"
)

// BorderStyle selects the border type, color, and corner radius of an item or image.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#BorderStyle
type BorderStyle struct {
	// Type selects an outline or no border. Omission uses an outline.
	Type BorderType `json:"type,omitempty"`
	// StrokeColor specifies the outline color when Type is BorderTypeStroke.
	StrokeColor *Color `json:"strokeColor,omitempty"`
	// CornerRadius controls the rounding of the border corners.
	CornerRadius int `json:"cornerRadius,omitempty"`
}

// BorderType selects whether a border is drawn.
type BorderType string

const (
	BorderTypeNone   BorderType = "NO_BORDER"
	BorderTypeStroke BorderType = "STROKE"
)
