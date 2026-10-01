package slack

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/JSYoo5B/convertago/internal/validation"
	"github.com/go-playground/validator/v10"
)

// RegisterValidation adds container-specific text limits and message constraints to v.
// Call it before using v, then validate messages or individual models with v.Struct.
// Field tags apply the limits that do not depend on a containing model.
func RegisterValidation(v *validator.Validate) {
	v.RegisterStructValidation(validateHeader, HeaderBlock{})
	v.RegisterStructValidation(validateSection, SectionBlock{})
	v.RegisterStructValidation(validateImage, ImageBlock{})
	v.RegisterStructValidation(validateButton, ButtonElement{})
	v.RegisterStructValidation(validateConfirmation, ConfirmationDialogObject{})
	v.RegisterStructValidation(validateVideo, VideoBlock{})
	v.RegisterStructValidation(validateMessage, Message{})
	v.RegisterStructValidation(validateLink, LinkInline{})
	v.RegisterStructValidation(validateUser, UserInline{})
}

func validateText(sl validator.StructLevel, field string, text string, maximum int) {
	if utf8.RuneCountInString(text) > maximum {
		sl.ReportError(text, field, field, "max", strconv.Itoa(maximum))
	}
}

func validateHeader(sl validator.StructLevel) {
	header := sl.Current().Interface().(HeaderBlock)
	validateText(sl, "Text.Text", header.Text.Text, 150)
}

func validateSection(sl validator.StructLevel) {
	section := sl.Current().Interface().(SectionBlock)
	if validation.Value(section.Text) == nil && len(section.Fields) == 0 {
		sl.ReportError(section.Text, "Text", "Text", "required_without", "Fields")
	}
	for i, field := range section.Fields {
		switch text := validation.Value(field).(type) {
		case PlainTextObject:
			validateText(sl, fmt.Sprintf("Fields[%d].Text", i), text.Text, 2000)
		case MrkdwnTextObject:
			validateText(sl, fmt.Sprintf("Fields[%d].Text", i), text.Text, 2000)
		}
	}
}

func validateImage(sl validator.StructLevel) {
	image := sl.Current().Interface().(ImageBlock)
	if image.Title != nil {
		validateText(sl, "Title.Text", image.Title.Text, 2000)
	}
}

func validateButton(sl validator.StructLevel) {
	button := sl.Current().Interface().(ButtonElement)
	validateText(sl, "Text.Text", button.Text.Text, 75)
	if button.URL != "" && !validation.AbsoluteURI(button.URL) {
		sl.ReportError(button.URL, "URL", "URL", "absolute_uri", "")
	}
}

func validateConfirmation(sl validator.StructLevel) {
	dialog := sl.Current().Interface().(ConfirmationDialogObject)
	validateText(sl, "Title.Text", dialog.Title.Text, 100)
	validateText(sl, "Confirm.Text", dialog.Confirm.Text, 30)
	validateText(sl, "Deny.Text", dialog.Deny.Text, 30)
	switch text := validation.Value(dialog.Text).(type) {
	case PlainTextObject:
		validateText(sl, "Text.Text", text.Text, 300)
	case MrkdwnTextObject:
		validateText(sl, "Text.Text", text.Text, 300)
	}
}

func validateVideo(sl validator.StructLevel) {
	video := sl.Current().Interface().(VideoBlock)
	validateText(sl, "Title.Text", video.Title.Text, 199)
	if video.Description != nil {
		validateText(sl, "Description.Text", video.Description.Text, 199)
	}
	for _, field := range []struct {
		name string
		url  string
	}{{"VideoURL", video.VideoURL}, {"TitleURL", video.TitleURL}} {
		if field.url != "" && !validation.AbsoluteURI(field.url, "https") {
			sl.ReportError(field.url, field.name, field.name, "https_url", "")
		}
	}
}

func validateMessage(sl validator.StructLevel) {
	message := sl.Current().Interface().(Message)
	total := 0
	for _, block := range message.Blocks {
		if markdown, ok := validation.Value(block).(MarkdownBlock); ok {
			total += utf8.RuneCountInString(markdown.Text)
		}
	}
	if total > 12000 {
		sl.ReportError(message.Blocks, "Blocks", "Blocks", "markdown_max", "12000")
	}
}

func validateLink(sl validator.StructLevel) {
	link := sl.Current().Interface().(LinkInline)
	if !validation.AbsoluteURI(link.URL) {
		sl.ReportError(link.URL, "URL", "URL", "absolute_uri", "")
	}
	if link.Style != nil && link.Style.Code {
		sl.ReportError(link.Style.Code, "Style.Code", "Style.Code", "excluded", "")
	}
}

func validateUser(sl validator.StructLevel) {
	user := sl.Current().Interface().(UserInline)
	if user.Style != nil && user.Style.Code {
		sl.ReportError(user.Style.Code, "Style.Code", "Style.Code", "excluded", "")
	}
}
