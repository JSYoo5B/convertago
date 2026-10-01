package googlechat

import (
	"github.com/JSYoo5B/convertago/internal/validation"
	"github.com/go-playground/validator/v10"
)

// RegisterValidation adds native model constraints that cannot be expressed by field tags to v.
// Call it before using v, then validate messages or individual models with v.Struct.
// Field tags validate scalar values, required inputs, and mutually exclusive fields.
func RegisterValidation(v *validator.Validate) {
	v.RegisterStructValidation(validateOpenLink, OpenLink{})
	v.RegisterStructValidation(validateAction, Action{})
	v.RegisterStructValidation(validateOverflowItem, OverflowMenuItem{})
	v.RegisterStructValidation(validateButton, Button{})
	v.RegisterStructValidation(validateIcon, Icon{})
	v.RegisterStructValidation(validateImage, Image{})
}

func validateHTTPS(sl validator.StructLevel, field, text string) {
	if text != "" && !validation.AbsoluteURI(text, "https") {
		sl.ReportError(text, field, field, "https_url", "")
	}
}

func validateOpenLink(sl validator.StructLevel) {
	link := sl.Current().Interface().(OpenLink)
	if !validation.AbsoluteURI(link.URL) {
		sl.ReportError(link.URL, "URL", "URL", "absolute_uri", "")
	}
}

func validateAction(sl validator.StructLevel) {
	action := sl.Current().Interface().(Action)
	if action.AllWidgetsAreRequired && len(action.RequiredWidgets) != 0 {
		sl.ReportError(action.RequiredWidgets, "RequiredWidgets", "RequiredWidgets", "excluded_if", "AllWidgetsAreRequired true")
	}
}

func validateOverflowItem(sl validator.StructLevel) {
	item := sl.Current().Interface().(OverflowMenuItem)
	if item.OnClick.OverflowMenu != nil {
		sl.ReportError(item.OnClick.OverflowMenu, "OnClick.OverflowMenu", "OnClick.OverflowMenu", "excluded", "")
	}
}

func validateButton(sl validator.StructLevel) {
	button := sl.Current().Interface().(Button)
	if button.Color != nil && button.Type != "" && button.Type != ButtonTypeFilled {
		sl.ReportError(button.Type, "Type", "Type", "color_type", "FILLED")
	}
}

func validateIcon(sl validator.StructLevel) {
	icon := sl.Current().Interface().(Icon)
	validateHTTPS(sl, "IconURL", icon.IconURL)
}

func validateImage(sl validator.StructLevel) {
	image := sl.Current().Interface().(Image)
	validateHTTPS(sl, "ImageURL", image.ImageURL)
}
