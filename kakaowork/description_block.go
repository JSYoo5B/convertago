package kakaowork

import "encoding/json"

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/descriptionblock/

type DescriptionBlock struct {
	Content TextBlock `json:"content"`
	Term    string    `json:"term"`
	Accent  bool      `json:"accent,omitempty"`
}

func (d DescriptionBlock) Type() string {
	return "description"
}

func (d DescriptionBlock) String() string {
	return d.Term + ": " + d.Content.String()
}

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
