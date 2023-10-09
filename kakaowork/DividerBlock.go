package kakaowork

import "encoding/json"

type DividerBlock struct{}

func (d DividerBlock) Type() string {
	return "divider"
}

func (d DividerBlock) String() string {
	return "---"
}

func (d DividerBlock) MarshalJSON() ([]byte, error) {
	type _DividerBlock struct {
		Type string `json:"type"`
	}

	return json.Marshal(_DividerBlock{
		Type: d.Type(),
	})
}
