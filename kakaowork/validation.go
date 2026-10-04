package kakaowork

import (
	"github.com/JSYoo5B/convertago/internal/conversion"
	"github.com/JSYoo5B/convertago/internal/validation"
	"github.com/go-playground/validator/v10"
)

// Validate 는 직접 구성한 메시지에 변환과 같은 규칙을 적용합니다.
// Fatal 위반은 오류로 반환하고, Warning 과 Advisory 위반은 WithDiagnostics 로 전달합니다.
// WithWarningAsError 를 지정하면 Warning 위반도 오류로 반환합니다. 진단 경로는 전송되는 JSON 의 경로입니다.
func Validate(message Message, options ...conversion.Option) error {
	return conversion.Validate("kakaowork", message.check, options)
}

// RegisterValidation 은 go-playground validator 에서 카카오워크 구조체를 검사하도록 v 에 등록합니다.
// validator 오류에는 심각도가 없으므로 Fatal 규칙만 보고합니다. 전체 규칙은 Validate 로 검사합니다.
func RegisterValidation(v *validator.Validate) {
	validation.Register(v, Message.check)
	validation.Register(v, HeaderBlock.check)
	validation.Register(v, TextBlock.check)
	validation.Register(v, InlineStyled.check)
	validation.Register(v, InlineLink.check)
	validation.Register(v, InlineMention.check)
	validation.Register(v, ImageBlock.check)
	validation.Register(v, ButtonBlock.check)
	validation.Register(v, ActionBlock.check)
	validation.Register(v, DescriptionBlock.check)
	validation.Register(v, SectionBlock.check)
	validation.Register(v, ContextBlock.check)
	validation.Register(v, OpenSystemBrowserAction.check)
	validation.Register(v, OpenInAppBrowserAction.check)
	validation.Register(v, OpenExternalAppAction.check)
	validation.Register(v, SubmitAction.check)
	validation.Register(v, CallModalAction.check)
	validation.Register(v, ExclusiveAction.check)
}
