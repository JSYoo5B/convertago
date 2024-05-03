package kakaowork

import "encoding/json"

// DescriptionBlock 은 말풍선 안에서 짧은 단어와 그에 대한 상세 정보를 설명할 때 사용합니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/descriptionblock/
type DescriptionBlock struct {
	// Content 에 상세 내용을 TextBlock 으로 입력합니다.
	Content TextBlock `json:"content"`
	// Term 에 Content 의 내용을 요약하는 문구로 입력합니다.
	Term string `json:"term"`
	// Accent 는 Term 의 내용을 굵은 글씨체로 적용하며, 기본값은 false 입니다.
	Accent bool `json:"accent,omitempty"`
}

func (d DescriptionBlock) Type() string   { return "description" }
func (d DescriptionBlock) String() string { return d.Term + ": " + d.Content.String() }
func (d DescriptionBlock) MarshalJSON() ([]byte, error) {
	type Embed DescriptionBlock
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  d.Type(),
		Embed: (Embed)(d),
	})
}
