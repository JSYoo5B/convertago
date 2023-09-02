package kakaowork

import "encoding/json"

// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/imagelinkblock/

type ImageBlock struct {
	Url string
}

func (i ImageBlock) Type() string {
	return "image_link"
}

func (i ImageBlock) String() string {
	return `{"image": "` + i.Url + `"}`
}

func (i ImageBlock) MarshalJSON() ([]byte, error) {
	type _ImageLink struct {
		Type string `json:"type"`
		Url  string `json:"url"`
	}

	return json.Marshal(_ImageLink{
		Type: i.Type(),
		Url:  i.Url,
	})
}
