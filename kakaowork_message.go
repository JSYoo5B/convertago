package convertago

import (
	"errors"
	"fmt"
	"github.com/JSYoo5B/convertago/kakaowork"
	"reflect"
	"strings"
)

func MarshalKakaoworkMessage(v any) (kakaowork.Message, error) {
	Value := unwrapValue(v)
	Type := Value.Type()

	if Type.Kind() != reflect.Struct {
		return kakaowork.Message{}, errors.New("cannot marshal non-struct value")
	}

	preview, blocks := "", make([]kakaowork.BubbleBlock, 0, Type.NumField())
	for i := 0; i < Type.NumField(); i++ {
		fieldType, fieldValue := Type.Field(i), dereference(Value.Field(i))
		if !fieldType.IsExported() {
			continue
		}

		tagValues := strings.SplitN(fieldType.Tag.Get("kakaowork"), ";", 2)
		convertType := strings.ToLower(tagValues[0])
		var option string
		if len(tagValues) == 2 {
			option = tagValues[1]
		}

		switch convertType {
		case "preview":
			preview += fmt.Sprintf("%v", fieldValue)
		case "text":
			// Convert text block
		case "image", "image_link":
			// Convert image block
		case "button":
			// Convert button block
		case "divider":
			blocks = append(blocks, kakaowork.DividerBlock{})
		case "header":
			blocks = append(blocks, convertKakaoworkHeaderBlock(fieldValue, option))
		case "action":
			// Convert action block
		case "description":
			// Convert description block
		case "section":
			// Convert section block
		case "context":
			// Convert context block
		}
	}

	return kakaowork.Message{Preview: preview, Blocks: blocks}, nil
}

func convertKakaoworkHeaderBlock(v reflect.Value, option string) kakaowork.HeaderBlock {
	option = strings.ToLower(option)
	return kakaowork.HeaderBlock{
		Text:  fmt.Sprintf("%v", v),
		Style: kakaowork.HeaderStyles[kakaowork.HeaderStyle(option)],
	}
}
