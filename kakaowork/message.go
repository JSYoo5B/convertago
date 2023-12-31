package kakaowork

// BubbleBlock is type casting interface for Message.Blocks.
// Implementing structs should expose only customizable properties.
// Other fixed properties should be handled during marshaling to JSON.
type BubbleBlock interface {
	// Type returns the message block type, equivalent to the JSON "type" property's value.
	Type() string
	// String is required for struct conversion fallback and value identification.
	String() string
	// MarshalJSON is overridden to hide fixed properties, like the "type" property.
	MarshalJSON() ([]byte, error)
}

type Message struct {
	// Preview describes simple text for notification and chat preview.
	Preview string `json:"text"`
	// Blocks describes the actual message contents.
	Blocks []BubbleBlock `json:"blocks"`
}
