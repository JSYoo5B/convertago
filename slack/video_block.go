package slack

import "encoding/json"

// VideoBlock embeds a video player. The posting app needs the links.embed:write scope.
//
// Reference: https://docs.slack.dev/reference/block-kit/blocks/video-block/
type VideoBlock struct {
	// Title is a plain-text video title shorter than 200 characters.
	Title PlainTextObject `json:"title"`
	// VideoURL is an HTTPS embeddable URL matching the app's unfurl domains.
	VideoURL string `json:"video_url" validate:"required,http_url"`
	// AltText is an accessible tooltip for the video.
	AltText string `json:"alt_text" validate:"required"`
	// ThumbnailURL points to the video's thumbnail image.
	ThumbnailURL string `json:"thumbnail_url" validate:"required,http_url"`
	// TitleURL links the title to the video's non-embeddable HTTPS URL.
	TitleURL string `json:"title_url,omitempty" validate:"omitempty,http_url"`
	// Description is a plain-text description shorter than 200 characters.
	Description *PlainTextObject `json:"description,omitempty" validate:"omitempty"`
	// AuthorName identifies the author in fewer than 50 characters.
	AuthorName string `json:"author_name,omitempty" validate:"max=49"`
	// ProviderName identifies the originating app or domain.
	ProviderName string `json:"provider_name,omitempty"`
	// ProviderIconURL points to the provider's icon.
	ProviderIconURL string `json:"provider_icon_url,omitempty" validate:"omitempty,http_url"`
	// BlockID identifies the block, up to 255 characters. Replace it when updating a message.
	BlockID string `json:"block_id,omitempty" validate:"max=255"`
}

func (b VideoBlock) Type() string   { return "video" }
func (b VideoBlock) String() string { return b.Title.String() }
func (VideoBlock) block()           {}
func (b VideoBlock) MarshalJSON() ([]byte, error) {
	type Embed VideoBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{b.Type(), Embed(b)})
}
