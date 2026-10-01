package kakaowork

import (
	"fmt"
	"net/url"
	"unicode/utf8"

	"github.com/JSYoo5B/convertago/internal/validation"
	"github.com/go-playground/validator/v10"
)

// RegisterValidation 은 필드 태그만으로 표현할 수 없는 카카오워크 구조체의 검증을 v 에 등록합니다.
// 헤더 위치, 인라인 텍스트의 전체 길이, 링크 스킴과 버튼 액션의 조합을 검사합니다.
// v 를 사용하기 전에 등록하고, v.Struct 로 메시지나 개별 구조체를 검사합니다.
func RegisterValidation(v *validator.Validate) {
	v.RegisterStructValidation(validateMessage, Message{})
	v.RegisterStructValidation(validateTextBlock, TextBlock{})
	v.RegisterStructValidation(validateInlineLink, InlineLink{})
	v.RegisterStructValidation(validateExternalAction, OpenExternalAppAction{})
	v.RegisterStructValidation(validateExclusiveAction, ExclusiveAction{})
}

func validateMessage(sl validator.StructLevel) {
	message := sl.Current().Interface().(Message)
	for i, block := range message.Blocks {
		field := fmt.Sprintf("Blocks[%d]", i)
		switch validation.Value(block).(type) {
		case HeaderBlock:
			if i != 0 {
				sl.ReportError(block, field, field, "header_position", "")
			}
		case TextBlock, ImageBlock, ButtonBlock, DividerBlock, DescriptionBlock, SectionBlock, ContextBlock:
		case nil:
			// The slice's required tag reports nil blocks.
		default:
			sl.ReportError(block, field, field, "block_type", "")
		}
	}
}

func validateTextBlock(sl validator.StructLevel) {
	block := sl.Current().Interface().(TextBlock)
	total := 0
	for _, inline := range block.Inlines {
		switch value := validation.Value(inline).(type) {
		case InlineStyled:
			total += utf8.RuneCountInString(value.Text)
		case InlineLink:
			total += utf8.RuneCountInString(value.Text)
		case InlineMention:
			total += utf8.RuneCountInString(value.Text)
		}
	}
	if total > 500 {
		sl.ReportError(block.Inlines, "Inlines", "Inlines", "max", "500")
	}
}

func validateInlineLink(sl validator.StructLevel) {
	link := sl.Current().Interface().(InlineLink)
	if !validation.AbsoluteURI(link.Url, "http", "https", "mailto", "tel") {
		sl.ReportError(link.Url, "Url", "Url", "link_uri", "")
	}
}

func validateExternalAction(sl validator.StructLevel) {
	action := sl.Current().Interface().(OpenExternalAppAction)
	if validation.AbsoluteURI(action.Value) {
		return
	}
	values, err := url.ParseQuery(action.Value)
	valid := err == nil && len(values) != 0
	for key, destinations := range values {
		valid = valid && (key == "ios" || key == "aos") && len(destinations) == 1 && validation.AbsoluteURI(destinations[0])
	}
	if !valid {
		sl.ReportError(action.Value, "Value", "Value", "app_uri", "")
	}
}

func validateExclusiveAction(sl validator.StructLevel) {
	action := sl.Current().Interface().(ExclusiveAction)
	for _, field := range []struct {
		name   string
		action ButtonAction
	}{
		{"Default", action.Default}, {"Pc", action.Pc}, {"Mobile", action.Mobile},
		{"Windows", action.Windows}, {"MacOs", action.MacOs}, {"Android", action.Android}, {"Ios", action.Ios},
	} {
		if _, nested := validation.Value(field.action).(ExclusiveAction); nested {
			sl.ReportError(field.action, field.name, field.name, "action_type", "")
		}
	}
}
