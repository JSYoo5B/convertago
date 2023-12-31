package kakaowork

import "encoding/json"

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/sectionblock/

type SectionBlock struct {
	Content   TextBlock    `json:"content"`
	Accessory *ImageBlock  `json:"accessory,omitempty"`
	Action    ButtonAction `json:"action,omitempty"`
}

func (s SectionBlock) Type() string {
	return "section"
}

func (s SectionBlock) String() string {
	return s.Content.String()
}

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
