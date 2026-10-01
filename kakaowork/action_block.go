package kakaowork

import (
	"encoding/json"
	"strings"
)

// ActionBlock 은 말풍선 안에서 두 개 또는 세 개의 버튼을 한 행에 나란히 배치합니다.
// 각 버튼의 레이블은 최대 20자이며, 기본 속성은 ButtonBlock 과 동일합니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/actionblock/
type ActionBlock struct {
	// Elements 는 한 행에 배치할 ButtonBlock 의 목록입니다.
	Elements []ButtonBlock `json:"elements"`
}

func (b ActionBlock) Type() string { return "action" }
func (b ActionBlock) String() string {
	var text []string
	for _, button := range b.Elements {
		text = append(text, button.String())
	}
	return strings.Join(text, " ")
}
func (ActionBlock) bubbleBlock() {}
func (b ActionBlock) MarshalJSON() ([]byte, error) {
	type Embed ActionBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{b.Type(), Embed(b)})
}
