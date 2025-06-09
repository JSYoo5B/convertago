package kakaowork

import "encoding/json"

// DividerBlock 은 말풍선 안에서 구분 선을 표현하는 블록입니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/dividerblock/
type DividerBlock struct{}

func (DividerBlock) BubbleType() string { return "divider" }
func (d DividerBlock) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type string `json:"type"`
	}{
		Type: d.BubbleType(),
	})
}
