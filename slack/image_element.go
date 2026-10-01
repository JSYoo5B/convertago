package slack

import "encoding/json"

// ImageElement displays an image inside a section or context block.
//
// Reference: https://docs.slack.dev/reference/block-kit/block-elements/image-element/
type ImageElement struct {
	// ImageURL is a public image URL, up to 3000 characters. Supply this or SlackFile.
	ImageURL string `json:"image_url,omitempty" validate:"required_without=SlackFile,excluded_with=SlackFile,omitempty,http_url,max=3000"`
	// SlackFile identifies a Slack-hosted image instead of ImageURL.
	SlackFile *SlackFileObject `json:"slack_file,omitempty" validate:"required_without=ImageURL,excluded_with=ImageURL,omitempty"`
	// AltText describes the image without markup for accessibility.
	AltText string `json:"alt_text" validate:"required,max=2000"`
}

func (e ImageElement) Type() string   { return "image" }
func (e ImageElement) String() string { return e.AltText }
func (ImageElement) element()         {}
func (ImageElement) contextElement()  {}
func (e ImageElement) MarshalJSON() ([]byte, error) {
	type Embed ImageElement
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{e.Type(), Embed(e)})
}
