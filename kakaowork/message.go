package kakaowork

// BubbleBlock 은 Message.Blocks 에 넣을 타입 캐스팅 인터페이스입니다.
// 해당 인터페이스를 구현하는 구조체는 설정 가능한 속성만 노출해야합니다.
// 고정된 속성들은 MarshalJSON 에서 추가 처리되어야 합니다.
type BubbleBlock interface {
	// Type 은 메시지 블록 타입을 반환해야합니다. JSON 의 "type" 속성을 나타냅니다.
	Type() string
	// String 은 구조체 변환 실패 혹은 값 구별을 위해 구현합니다.
	String() string
	// MarshalJSON 은 고정 속성값들을 숨기면서 원래 사양에 맞게 JSON 변환을 제공할 수 있어야 합니다.
	// (대표적으로 "type" 속성)
	MarshalJSON() ([]byte, error)
}

type Message struct {
	// Preview 는 알림과 채팅 미리보기에서 사용할 간단한 텍스트를 입력
	Preview string `json:"text"`
	// Blocks 은 실제 메시지 내용을 기술
	Blocks []BubbleBlock `json:"blocks"`
}
