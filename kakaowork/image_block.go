package kakaowork

import "encoding/json"

// ImageBlock 은 말풍선 안에서 이미지를 표현하는 블록으로, 레이아웃 블록을 구성하는 엘리먼트의 속성으로 사용되기도 합니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/imagelinkblock/
type ImageBlock struct {
	// Url 에 출력할 이미지의 URL 을 입력합니다.
	Url string `json:"url" validate:"required,url"`
}

func (i ImageBlock) Type() string   { return "image_link" }
func (i ImageBlock) String() string { return `{"image": "` + i.Url + `"}` }
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
