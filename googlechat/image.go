package googlechat

// Image displays a URL-hosted image with an optional click action.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Image
type Image struct {
	// ImageURL is the HTTPS URL hosting the image.
	ImageURL string `json:"imageUrl"`
	// OnClick runs when the image is clicked.
	OnClick *OnClick `json:"onClick,omitempty"`
	// AltText describes the image for accessibility.
	AltText string `json:"altText,omitempty"`
}

func (Image) WidgetType() string   { return "image" }
func (i Image) String() string     { return i.AltText }
func (Image) widgetContent()       {}
func (Image) columnWidgetContent() {}
func (Image) nestedWidgetContent() {}

// Divider draws a horizontal line between widgets.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Divider
type Divider struct{}

func (Divider) WidgetType() string { return "divider" }
func (Divider) String() string     { return "" }
func (Divider) widgetContent()     {}
