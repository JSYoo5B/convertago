package kakaowork

import "encoding/json"

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/imagelinkblock/

type ImageBlock struct {
	Url string `json:"url" validate:"required,url"`
}

func (i ImageBlock) Type() string {
	return "image_link"
}

func (i ImageBlock) String() string {
	return `{"image": "` + i.Url + `"}`
}

func (i ImageBlock) MarshalJSON() ([]byte, error) {
	type Embed ImageBlock
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  i.Type(),
		Embed: (Embed)(i),
	})
}
