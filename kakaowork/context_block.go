package kakaowork

import "encoding/json"

// ContextBlock 은 ImageBlock(좌측)과 TextBlock(우측)이 조합된 블록으로
// 추가적인 정보를 덧붙여 표현하기 위해 사용됩니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/contextblock/
type ContextBlock struct {
	// Content 는 필수 요소이며, 흐린 색의 폰트로 표시,
	// TextBlock 의 InlineLink 를 활용하여 링크 적용 가능
	Content TextBlock `json:"content"`
	// Image 는 필수 요소이며, 24px * 24px 사이즈로 고정
	Image ImageBlock `json:"image"`
}

func (ContextBlock) Type() string     { return "context" }
func (c ContextBlock) String() string { return c.Content.String() }
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
