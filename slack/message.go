// Package slack provides Block Kit objects for composing Slack messages.
package slack

// Block is a layout block in a Slack message.
// Type returns the JSON discriminator, and String returns a textual representation.
//
// Reference: https://docs.slack.dev/reference/block-kit/blocks/
type Block interface {
	Type() string
	String() string
	MarshalJSON() ([]byte, error)
	block()
}

// Message contains text and Block Kit content for a Slack message.
//
// Reference: https://docs.slack.dev/reference/methods/chat.postMessage/
type Message struct {
	// Text is the notification and accessibility fallback when Blocks is supplied.
	Text string `json:"text,omitempty"`
	// Blocks contains up to 50 layout blocks in display order.
	Blocks []Block `json:"blocks,omitempty"`
}
