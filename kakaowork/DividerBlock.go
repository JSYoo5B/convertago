package kakaowork

import "encoding/json"

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/dividerblock/

type DividerBlock struct{}

func (d DividerBlock) Type() string {
	return "divider"
}

func (d DividerBlock) String() string {
	return "---"
}

func (d DividerBlock) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type string `json:"type"`
	}{
		Type: d.Type(),
	})
}
