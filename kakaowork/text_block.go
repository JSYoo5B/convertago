package kakaowork

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// TextBlock 은 말풍선에서 가장 기본인 텍스트를 표현하는 블록입니다.
// 단독으로도 존재할 수 있지만, 레이아웃 블록을 구성하는 엘리먼트의 속성으로 사용되기도 합니다.
//
// Reference: https://docs.kakaoi.ai/kakao_work/blockkit/textblock/
type TextBlock struct {
	// Text 에 표현할 텍스트를 입력.
	// 전체 텍스트를 \n을 포함하여 기술,
	// Inlines 서식을 적용할 텍스트까지 포함하여 작성 필요,
	// Inlines 와 정합성이 맞지 않을 경우, Inlines 의 String()를 우선으로 적용함
	Text string `json:"text"`
	// Inlines 에 다양한 추가 서식을 적용하여 텍스트 및 스타일 설정 가능.
	// 적용 가능한 서식: InlineStyled, InlineLink, InlineMention 참고
	Inlines []Inline `json:"inlines,omitempty"`
}

func (TextBlock) Type() string { return "text" }
func (t TextBlock) String() string {
	if t.Inlines == nil {
		return t.Text
	} else {
		var inlineTexts []string
		for _, inline := range t.Inlines {
			inlineTexts = append(inlineTexts, inline.String())
		}
		return strings.Join(inlineTexts, "")
	}
}
func (t TextBlock) MarshalJSON() ([]byte, error) {
	type Embed TextBlock
	return json.Marshal(&struct {
		Type string `json:"type"`
		Embed
	}{
		Type:  t.Type(),
		Embed: (Embed)(t),
	})
}

func ConvertTextBlock(target any) (result TextBlock, err error) {
	Type, Value := reflect.TypeOf(target), reflect.ValueOf(target)
	for i := 0; i < Type.NumField(); i++ {
		fieldType, fieldValue := Type.Field(i), Value.Field(i)
		tag := fieldType.Tag.Get("kakaowork")
		if !fieldType.IsExported() {
			continue
		} else if tag == "" {
			continue
		}

		tags := strings.Split(tag, ";")
		if len(tags) == 0 ||
			(tags[0] != "TextBlock" && tags[0] != "text") {
			continue
		}

		result, err = convertTextField(result, tags[1:], fieldValue)
		if err != nil {
			return TextBlock{}, err
		}
	}

	return result, nil
}

func convertTextField(previous TextBlock, _ []string, value reflect.Value) (result TextBlock, err error) {
	result = previous
	result.Text += fmt.Sprintf("%v", value)

	return result, nil
}
