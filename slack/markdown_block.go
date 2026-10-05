package slack

import "encoding/json"

// MarkdownBlock lets Slack render standard Markdown from apps using platform AI features.
// Slack may translate one MarkdownBlock into multiple layout blocks.
//
// Reference: https://docs.slack.dev/reference/block-kit/blocks/markdown-block/
type MarkdownBlock struct {
	// Text contains standard Markdown. All markdown blocks together have a 12000-character limit.
	Text string `json:"text"`
}

func (b MarkdownBlock) Type() string   { return "markdown" }
func (b MarkdownBlock) String() string { return b.Text }
func (MarkdownBlock) block()           {}
func (b MarkdownBlock) MarshalJSON() ([]byte, error) {
	type Embed MarkdownBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{b.Type(), Embed(b)})
}
