package kakaowork

type MessageBubbleBlock interface {
	Type() string
	String() string
	MarshalJSON() ([]byte, error)
}
