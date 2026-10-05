package slack

import (
	"github.com/JSYoo5B/convertago/internal/conversion"
	"github.com/JSYoo5B/convertago/internal/validation"
	"github.com/go-playground/validator/v10"
)

// Validate applies the conversion rules to a manually constructed message.
// Fatal violations are returned as errors; Warning and Advisory violations go to WithDiagnostics.
// WithWarningAsError also returns Warning violations as errors. Diagnostic paths follow the sent JSON.
func Validate(message Message, options ...conversion.Option) error {
	return conversion.Validate("slack", message.check, options)
}

// RegisterValidation lets a go-playground validator check Slack models with v.Struct.
// Validator errors carry no severity, so it reports only Fatal rules; use Validate for every rule.
func RegisterValidation(v *validator.Validate) {
	validation.Register(v, Message.check)
	validation.Register(v, PlainTextObject.check)
	validation.Register(v, MrkdwnTextObject.check)
	validation.Register(v, HeaderBlock.check)
	validation.Register(v, SectionBlock.check)
	validation.Register(v, ImageBlock.check)
	validation.Register(v, ImageElement.check)
	validation.Register(v, SlackFileObject.check)
	validation.Register(v, ActionsBlock.check)
	validation.Register(v, ContextBlock.check)
	validation.Register(v, DividerBlock.check)
	validation.Register(v, MarkdownBlock.check)
	validation.Register(v, ButtonElement.check)
	validation.Register(v, ConfirmationDialogObject.check)
	validation.Register(v, VideoBlock.check)
	validation.Register(v, RichTextBlock.check)
	validation.Register(v, RichTextSection.check)
	validation.Register(v, RichTextList.check)
	validation.Register(v, RichTextPreformatted.check)
	validation.Register(v, RichTextQuote.check)
	validation.Register(v, TextInline.check)
	validation.Register(v, LinkInline.check)
	validation.Register(v, UserInline.check)
	validation.Register(v, EmojiInline.check)
}
