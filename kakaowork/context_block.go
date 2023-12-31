package kakaowork

import "encoding/json"

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/contextblock/

type ContextBlock struct {
	Content TextBlock  `json:"content"`
	Image   ImageBlock `json:"image"`
}

func (c ContextBlock) Type() string {
	return "context"
}

func (c ContextBlock) String() string {
	return c.Content.String()
}

func (c ContextBlock) MarshalJSON() ([]byte, error) {
	type Embed ContextBlock
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  c.Type(),
		Embed: (Embed)(c),
	})
}
