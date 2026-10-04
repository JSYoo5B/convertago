package kakaowork

import "github.com/JSYoo5B/convertago/internal/validation"

// BubbleBlock 은 카카오워크 구조체가 공통으로 구현하는 인터페이스입니다.
// kakaowork 패키지에서 제공하는 구현체만 직접 구현할 수 있습니다.
// 버튼 액션과 인라인도 이 인터페이스를 구현하지만, Message.Blocks 에는 말풍선 블록만 넣을 수 있습니다.
// 해당 인터페이스를 구현하는 구조체는 설정 가능한 속성만 노출해야합니다.
// 고정된 속성들은 MarshalJSON 에서 추가 처리되어야 합니다.
type BubbleBlock interface {
	// Type 은 블록의 JSON "type" 속성을 반환합니다.
	// 버튼 액션은 "action", 인라인은 "inline" 을 반환하며, JSON 타입은 ActionType 과 InlineType 이 반환합니다.
	Type() string
	// String 은 구조체 변환 실패 혹은 값 구별을 위해 구현합니다.
	String() string
	// MarshalJSON 은 고정 속성값들을 숨기면서 원래 사양에 맞게 JSON 변환을 제공할 수 있어야 합니다.
	// (대표적으로 "type" 속성)
	MarshalJSON() ([]byte, error)
	bubbleBlock()
	check(*validation.Check)
}

// Message 는 카카오워크로 전송할 Block Kit 메시지입니다.
// 직접 구성한 메시지는 Validate 로 검사할 수 있습니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/webapireference/messages/
type Message struct {
	// Preview 는 알림과 채팅 미리보기에서 사용할 간단한 텍스트를 입력
	Preview string `json:"text"`
	// Blocks 은 실제 메시지 내용을 기술
	Blocks []BubbleBlock `json:"blocks" validate:"dive"`
}
