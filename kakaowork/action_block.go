package kakaowork

import "encoding/json"

// ActionBlock 은 말풍선 안에서 여러 개의 버튼을 하나의 행에 표현할 수 있는 레이아웃 블록입니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/actionblock/
type ActionBlock struct {
	// Buttons 에 ButtonBlock 목록을 입력
	Buttons []ButtonBlock `json:"elements"`
}

func (ActionBlock) BubbleType() string { return "action" }
func (a ActionBlock) MarshalJSON() ([]byte, error) {
	type Embed ActionBlock
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  a.BubbleType(),
		Embed: (Embed)(a),
	})
}
