package googlechat

import (
	"github.com/JSYoo5B/convertago/internal/conversion"
	"github.com/JSYoo5B/convertago/internal/validation"
	"github.com/go-playground/validator/v10"
)

// Validate applies the conversion rules to a manually constructed message.
// Fatal violations are returned as errors; Warning and Advisory violations go to WithDiagnostics.
// WithWarningAsError also returns Warning violations as errors. Diagnostic paths follow the sent JSON.
func Validate(message Message, options ...conversion.Option) error {
	return conversion.Validate("googlechat", message.check, options)
}

// RegisterValidation lets a go-playground validator check Google Chat models with v.Struct.
// Validator errors carry no severity, so it reports only Fatal rules; use Validate for every rule.
func RegisterValidation(v *validator.Validate) {
	validation.Register(v, Message.check)
	validation.Register(v, Card.check)
	validation.Register(v, CardHeader.check)
	validation.Register(v, Section.check)
	validation.Register(v, CollapseControl.check)
	validation.Register(v, Widget.check)
	validation.Register(v, TextParagraph.check)
	validation.Register(v, Image.check)
	validation.Register(v, DecoratedText.check)
	validation.Register(v, SwitchControl.check)
	validation.Register(v, ButtonList.check)
	validation.Register(v, Button.check)
	validation.Register(v, Color.check)
	validation.Register(v, Icon.check)
	validation.Register(v, MaterialIcon.check)
	validation.Register(v, OnClick.check)
	validation.Register(v, Action.check)
	validation.Register(v, OpenLink.check)
	validation.Register(v, OverflowMenu.check)
	validation.Register(v, OverflowMenuItem.check)
	validation.Register(v, Columns.check)
	validation.Register(v, Column.check)
	validation.Register(v, Grid.check)
	validation.Register(v, GridItem.check)
	validation.Register(v, ImageComponent.check)
	validation.Register(v, ImageCropStyle.check)
	validation.Register(v, BorderStyle.check)
	validation.Register(v, Carousel.check)
	validation.Register(v, CarouselCard.check)
	validation.Register(v, ChipList.check)
	validation.Register(v, Chip.check)
}
