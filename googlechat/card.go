package googlechat

// Card groups a header and sections of widgets into a defined layout.
// Google Chat displays up to 100 widgets per card and ignores additional widgets.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Card
type Card struct {
	// Header appears above the card's sections.
	Header *CardHeader `json:"header,omitempty" validate:"required_without=Sections,omitempty"`
	// Sections groups widgets beneath optional section headings.
	Sections []Section `json:"sections,omitempty" validate:"dive"`
	// SectionDividerStyle controls the dividers between the header and sections.
	SectionDividerStyle DividerStyle `json:"sectionDividerStyle,omitempty" validate:"omitempty,oneof=SOLID_DIVIDER NO_DIVIDER"`
}

// CardHeader displays a title, an optional subtitle, and an optional leading image.
// With both text fields set, each occupies one line; a title alone occupies both.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#CardHeader
type CardHeader struct {
	// Title is the required heading of the card.
	Title string `json:"title" validate:"required"`
	// Subtitle appears below Title.
	Subtitle string `json:"subtitle,omitempty"`
	// ImageURL is the HTTPS URL of the header image.
	ImageURL string `json:"imageUrl,omitempty" validate:"required_with=ImageType ImageAltText,omitempty,http_url"`
	// ImageType selects a square or circular crop for the image.
	ImageType ImageType `json:"imageType,omitempty" validate:"omitempty,oneof=SQUARE CIRCLE"`
	// ImageAltText describes the image for accessibility.
	ImageAltText string `json:"imageAltText,omitempty"`
}

// Section displays its widgets vertically in the order supplied.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Section
type Section struct {
	// Header is an optional heading supporting simple HTML formatting.
	Header string `json:"header,omitempty"`
	// Widgets must contain at least one widget.
	Widgets []Widget `json:"widgets" validate:"min=1,max=100,dive"`
	// Collapsible allows users to expand and collapse some or all widgets.
	Collapsible bool `json:"collapsible,omitempty"`
	// UncollapsibleWidgetsCount keeps this many leading widgets visible when collapsed.
	UncollapsibleWidgetsCount int `json:"uncollapsibleWidgetsCount,omitempty" validate:"min=0,excluded_if=Collapsible false"`
	// CollapseControl customizes the controls when Collapsible is true.
	CollapseControl *CollapseControl `json:"collapseControl,omitempty" validate:"excluded_if=Collapsible false,omitempty"`
}

// CollapseControl customizes the expand and collapse buttons of a section.
// Both buttons must be supplied for the customization to take effect.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#CollapseControl
type CollapseControl struct {
	// HorizontalAlignment positions the controls within the section.
	HorizontalAlignment HorizontalAlignment `json:"horizontalAlignment,omitempty" validate:"omitempty,oneof=START CENTER END"`
	// ExpandButton reveals the collapsed widgets.
	ExpandButton Button `json:"expandButton"`
	// CollapseButton hides the collapsible widgets.
	CollapseButton Button `json:"collapseButton"`
}

// DividerStyle controls the dividers between card sections.
type DividerStyle string

const (
	DividerStyleSolid DividerStyle = "SOLID_DIVIDER"
	DividerStyleNone  DividerStyle = "NO_DIVIDER"
)

// ImageType selects the mask applied to an image.
type ImageType string

const (
	ImageTypeSquare ImageType = "SQUARE"
	ImageTypeCircle ImageType = "CIRCLE"
)
