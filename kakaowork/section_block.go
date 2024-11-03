package kakaowork

import "encoding/json"

// SectionBlock 은 TextBlock(좌측)과 ImageBlock(우측)이 조합된 레이아웃 블록으로,
// 말풍선 안에서 텍스트와 함께 추가 정보를 이미지로 표현하고자 할 때 사용됩니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/sectionblock/
type SectionBlock struct {
	// Content 에 텍스트를 표현하는 TextBlock 설정
	Content TextBlock `json:"content"`
	// Accessory 에 이미지를 표현하는 ImageBlock 설정 (선택)
	Accessory *ImageBlock `json:"accessory,omitempty"`
	// ButtonAction 에 클릭 시 수행될 ButtonAction 설정 (선택)
	Action ButtonAction `json:"action,omitempty"`
}

func (SectionBlock) Type() string     { return "section" }
func (s SectionBlock) String() string { return s.Content.String() }
func (s SectionBlock) MarshalJSON() ([]byte, error) {
	type Embed SectionBlock
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  s.Type(),
		Embed: (Embed)(s),
	})
}
