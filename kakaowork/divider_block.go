package kakaowork

import "encoding/json"

// DividerBlock 은 말풍선 안에서 구분 선을 표현하는 블록입니다.
// BubbleBlock 들 사이에 자유롭게 위치할 수 있으며, Message.Blocks 최상단과 최하단에 배치하는 것은 지양합니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/dividerblock/
type DividerBlock struct{}

func (d DividerBlock) Type() string   { return "divider" }
func (d DividerBlock) String() string { return "---" }
func (d DividerBlock) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type string `json:"type"`
	}{
		Type: d.Type(),
	})
}
