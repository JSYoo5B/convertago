package googlechat

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

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
	v.RegisterStructValidation(validateHeader, CardHeader{})
	v.RegisterStructValidation(validateSection, Section{})
	v.RegisterStructValidation(validateCard, Card{})
	v.RegisterStructValidation(validateMessage, Message{})
	v.RegisterStructValidation(validateCrop, ImageCropStyle{})
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

func validateHeader(sl validator.StructLevel) {
	header := sl.Current().Interface().(CardHeader)
	validateHTTPS(sl, "ImageURL", header.ImageURL)
}

func validateSection(sl validator.StructLevel) {
	section := sl.Current().Interface().(Section)
	if section.UncollapsibleWidgetsCount > len(section.Widgets) {
		sl.ReportError(section.UncollapsibleWidgetsCount, "UncollapsibleWidgetsCount", "UncollapsibleWidgetsCount", "max", strconv.Itoa(len(section.Widgets)))
	}
}

func validateCard(sl validator.StructLevel) {
	card := sl.Current().Interface().(Card)
	if card.Header == nil && len(card.Sections) == 0 {
		sl.ReportError(card.Header, "Header", "Header", "required_without", "Sections")
	}
	total := 0
	for _, section := range card.Sections {
		for _, widget := range section.Widgets {
			total += countWidgets(widget.Content)
		}
	}
	if total > 100 {
		sl.ReportError(card.Sections, "Sections", "Sections", "max_widgets", "100")
	}
}

func countWidgets(content any) int {
	value := validation.Value(content)
	if value == nil {
		return 0
	}
	count := 1
	switch value := value.(type) {
	case Columns:
		for _, column := range value.ColumnItems {
			for _, widget := range column.Widgets {
				count += countWidgets(widget.Content)
			}
		}
	case Carousel:
		for _, card := range value.CarouselCards {
			for _, widget := range card.Widgets {
				count += countWidgets(widget.Content)
			}
			for _, widget := range card.FooterWidgets {
				count += countWidgets(widget.Content)
			}
		}
	}
	return count
}

func validateMessage(sl validator.StructLevel) {
	message := sl.Current().Interface().(Message)
	ids := make(map[string]bool, len(message.CardsV2))
	for i, card := range message.CardsV2 {
		field := fmt.Sprintf("CardsV2[%d].CardID", i)
		if len(message.CardsV2) > 1 && card.CardID == "" {
			sl.ReportError(card.CardID, field, field, "required", "")
		}
		if card.CardID != "" {
			if ids[card.CardID] {
				sl.ReportError(card.CardID, field, field, "unique", "")
			}
			ids[card.CardID] = true
		}
	}
	if len(message.CardsV2) == 0 {
		return
	}
	data, err := json.Marshal(message.CardsV2)
	if err != nil {
		sl.ReportError(message.CardsV2, "CardsV2", "CardsV2", "json", "")
	} else if len(data) > 32*1024 {
		sl.ReportError(message.CardsV2, "CardsV2", "CardsV2", "max_bytes", "32768")
	}
}

func validateCrop(sl validator.StructLevel) {
	crop := sl.Current().Interface().(ImageCropStyle)
	if math.IsNaN(crop.AspectRatio) || math.IsInf(crop.AspectRatio, 0) {
		sl.ReportError(crop.AspectRatio, "AspectRatio", "AspectRatio", "finite", "")
	}
}
