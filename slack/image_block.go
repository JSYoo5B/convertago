package slack

import "encoding/json"

// ImageBlock displays a PNG, JPG, JPEG, or GIF image.
//
// Reference: https://docs.slack.dev/reference/block-kit/blocks/image-block/
type ImageBlock struct {
	// ImageURL is a public image URL, up to 3000 characters. Supply this or SlackFile.
	ImageURL string `json:"image_url,omitempty"`
	// SlackFile identifies a Slack-hosted image instead of ImageURL.
	SlackFile *SlackFileObject `json:"slack_file,omitempty"`
	// AltText describes the image without markup, up to 2000 characters.
	AltText string `json:"alt_text"`
	// Title is an optional plain-text title, up to 2000 characters.
	Title *PlainTextObject `json:"title,omitempty"`
	// BlockID identifies the block, up to 255 characters. Replace it when updating a message.
	BlockID string `json:"block_id,omitempty"`
}

func (b ImageBlock) Type() string   { return "image" }
func (b ImageBlock) String() string { return b.AltText }
func (ImageBlock) block()           {}
func (b ImageBlock) MarshalJSON() ([]byte, error) {
	type Embed ImageBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{b.Type(), Embed(b)})
}

// SlackFileObject identifies an image file accessible to the posting user.
// Supply either URL or ID; Slack rejects an object containing both.
//
// Reference: https://docs.slack.dev/reference/block-kit/composition-objects/slack-file-object/
type SlackFileObject struct {
	// URL is the file's private URL or permalink.
	URL string `json:"url,omitempty"`
	// ID is the Slack file ID.
	ID string `json:"id,omitempty"`
}
