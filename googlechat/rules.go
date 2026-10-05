package googlechat

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/JSYoo5B/convertago/internal/validation"
)

const (
	referenceCards    = "https://developers.google.com/workspace/chat/api/reference/rest/v1/cards"
	referenceMessages = "https://developers.google.com/workspace/chat/api/reference/rest/v1/spaces.messages"
)

const (
	fatal   = validation.Fatal
	warning = validation.Warning
)

func rule(id string, severity validation.Severity, verified bool, anchor, message string) validation.Rule {
	reference := referenceCards + "#" + anchor
	if anchor == "Message" || anchor == "CardWithId" {
		reference = referenceMessages + "#" + anchor
	}
	return validation.Rule{ID: id, Severity: severity, Verified: verified, Reference: reference, Message: message}
}

var (
	ruleMessageCardsSize  = rule("message.cards.size", warning, false, "Message", "cards exceed 32 KB of JSON")
	ruleCardIDRequired    = rule("message.card_id.required", fatal, true, "CardWithId", "multiple cards require a cardId")
	ruleCardIDUnique      = rule("message.card_id.unique", warning, false, "CardWithId", "cardId values must be unique")
	ruleCardContent       = rule("card.content.required", fatal, true, "Card", "a card requires a header or sections")
	ruleCardWidgetCount   = rule("card.widgets.count", warning, true, "Card", "sections that push a card over 100 widgets are ignored")
	ruleCardDividerStyle  = rule("card.section_divider_style.value", fatal, true, "Card", "unknown section divider style")
	ruleHeaderTitle       = rule("header.title.required", warning, false, "CardHeader", "header title is required")
	ruleHeaderImageHTTPS  = rule("header.image_url.https", warning, false, "CardHeader", "header imageUrl must be an HTTPS URL")
	ruleHeaderImageURL    = rule("header.image.url_required", warning, false, "CardHeader", "header image properties apply only with imageUrl")
	ruleImageTypeValue    = rule("image_type.value", fatal, true, "ImageType", "unknown image type")
	ruleSectionWidgets    = rule("section.widgets.required", warning, false, "Section", "a section requires widgets")
	ruleSectionCollapse   = rule("section.collapse.collapsible", warning, true, "Section", "collapse properties apply only to a collapsible section")
	ruleSectionUncollapse = rule("section.uncollapsible.count", warning, false, "Section", "uncollapsibleWidgetsCount exceeds the section's widgets")
	ruleWidgetContent     = rule("widget.content.required", fatal, true, "Widget", "a widget requires content")
	ruleHorizontalAlign   = rule("horizontal_alignment.value", fatal, true, "HorizontalAlignment", "unknown horizontal alignment")
	ruleVerticalAlign     = rule("vertical_alignment.value", fatal, true, "VerticalAlignment", "unknown vertical alignment")

	ruleParagraphMaxLines = rule("text_paragraph.max_lines.value", warning, false, "TextParagraph", "maxLines must not be negative")
	ruleParagraphSyntax   = rule("text_paragraph.text_syntax.value", fatal, true, "TextParagraph", "unknown text syntax")
	ruleImageURLHTTPS     = rule("image.image_url.https", warning, false, "Image", "imageUrl must be an HTTPS URL")
	ruleDecoratedText     = rule("decorated_text.text.required", warning, false, "DecoratedText", "decorated text is required")
	ruleDecoratedControl  = rule("decorated_text.control.count", fatal, true, "DecoratedText", "decorated text accepts only one of button, switchControl, or endIcon")
	ruleSwitchName        = rule("switch_control.name.required", warning, false, "SwitchControl", "switch control name is required")
	ruleSwitchType        = rule("switch_control.control_type.value", fatal, true, "SwitchControl", "unknown switch control type")

	ruleButtonListButtons = rule("button_list.buttons.required", warning, false, "ButtonList", "a button list requires buttons")
	ruleButtonContent     = rule("button.content.required", warning, false, "Button", "a button requires text or an icon")
	ruleButtonTypeValue   = rule("button.type.value", fatal, true, "Button", "unknown button type")
	ruleButtonTypeColor   = rule("button.type.color", warning, true, "Button", "a color forces the FILLED type and ignores type")
	ruleColorRange        = rule("color.component.range", warning, false, "Color", "color components must be between 0 and 1")
	ruleIconSource        = rule("icon.source.count", fatal, true, "Icon", "an icon requires exactly one of knownIcon, iconUrl, or materialIcon")
	ruleIconURLHTTPS      = rule("icon.icon_url.https", warning, false, "Icon", "iconUrl must be an HTTPS URL")
	ruleMaterialName      = rule("material_icon.name.required", warning, false, "MaterialIcon", "material icon name is required")
	ruleMaterialWeight    = rule("material_icon.weight.value", warning, true, "MaterialIcon", "an unsupported weight uses the default")
	ruleMaterialGrade     = rule("material_icon.grade.value", warning, true, "MaterialIcon", "an unsupported grade uses the default")

	ruleOnClickCount       = rule("on_click.action.count", fatal, true, "OnClick", "onClick requires exactly one of action, openLink, or overflowMenu")
	ruleActionFunction     = rule("action.function.required", warning, false, "Action", "action function is required")
	ruleActionLoad         = rule("action.load_indicator.value", fatal, true, "Action", "unknown load indicator")
	ruleActionInteraction  = rule("action.interaction.value", fatal, true, "Action", "unknown interaction")
	ruleActionRequired     = rule("action.required_widgets.conflict", warning, false, "Action", "requiredWidgets and allWidgetsAreRequired cannot be combined")
	ruleActionRequiredName = rule("action.required_widgets.required", warning, false, "Action", "requiredWidgets contains an empty name")
	ruleParameterKey       = rule("action.parameters.key_required", warning, false, "ActionParameter", "action parameter key is required")
	ruleParameterUnique    = rule("action.parameters.key_unique", warning, false, "ActionParameter", "action parameter keys must be unique")
	ruleOpenLinkURL        = rule("open_link.url.format", warning, false, "OpenLink", "openLink url must be an absolute URI")
	ruleOverflowItems      = rule("overflow_menu.items.required", warning, false, "OverflowMenu", "an overflow menu requires items")
	ruleOverflowItemText   = rule("overflow_menu_item.text.required", warning, false, "OverflowMenuItem", "overflow menu item text is required")
	ruleOverflowNested     = rule("overflow_menu_item.on_click.overflow", warning, true, "OverflowMenuItem", "a nested overflow menu is dropped and disables the item")

	ruleColumnsCount    = rule("columns.column_items.count", warning, false, "Columns", "columns display one or two columns")
	ruleColumnWidgets   = rule("column.widgets.required", warning, false, "Column", "a column requires widgets")
	ruleColumnSizeStyle = rule("column.horizontal_size_style.value", fatal, true, "Column", "unknown horizontal size style")
	ruleColumnVertical  = rule("column.vertical_alignment.value", fatal, true, "Column", "unknown column vertical alignment")
	ruleGridItems       = rule("grid.items.required", warning, false, "Grid", "a grid requires items")
	ruleGridColumnCount = rule("grid.column_count.value", warning, false, "Grid", "columnCount must not be negative")
	ruleGridItemContent = rule("grid_item.content.required", warning, false, "GridItem", "a grid item requires a title, subtitle, or image")
	ruleGridItemLayout  = rule("grid_item.layout.value", fatal, true, "GridItem", "unknown grid item layout")
	ruleImageURIFormat  = rule("image_component.image_uri.format", warning, false, "ImageComponent", "imageUri must be an absolute HTTP or HTTPS URL")
	ruleCropTypeValue   = rule("image_crop.type.value", fatal, true, "ImageCropStyle", "unknown crop type")
	ruleCropRatioValue  = rule("image_crop.aspect_ratio.value", warning, false, "ImageCropStyle", "a custom rectangle requires a positive finite aspectRatio")
	ruleCropRatioCustom = rule("image_crop.aspect_ratio.custom", warning, false, "ImageCropStyle", "aspectRatio applies only to RECTANGLE_CUSTOM")
	ruleBorderTypeValue = rule("border_style.type.value", fatal, true, "BorderStyle", "unknown border type")
	ruleBorderStroke    = rule("border_style.stroke_color.stroke", warning, false, "BorderStyle", "strokeColor applies only to a STROKE border")
	ruleBorderRadius    = rule("border_style.corner_radius.value", warning, false, "BorderStyle", "cornerRadius must not be negative")
	ruleCarouselCards   = rule("carousel.cards.required", warning, false, "Carousel", "a carousel requires cards")
	ruleCarouselWidgets = rule("carousel_card.widgets.required", warning, false, "CarouselCard", "a carousel card requires widgets")
	ruleChipListChips   = rule("chip_list.chips.required", warning, false, "ChipList", "a chip list requires chips")
	ruleChipListLayout  = rule("chip_list.layout.value", fatal, true, "ChipList", "unknown chip list layout")
	ruleChipContent     = rule("chip.content.required", warning, false, "Chip", "a chip requires a label or an icon")
)

// rules lists every Google Chat rule. The README documents each entry.
var rules = []validation.Rule{
	ruleMessageCardsSize, ruleCardIDRequired, ruleCardIDUnique, ruleCardContent, ruleCardWidgetCount, ruleCardDividerStyle,
	ruleHeaderTitle, ruleHeaderImageHTTPS, ruleHeaderImageURL, ruleImageTypeValue,
	ruleSectionWidgets, ruleSectionCollapse, ruleSectionUncollapse, ruleWidgetContent, ruleHorizontalAlign, ruleVerticalAlign,
	ruleParagraphMaxLines, ruleParagraphSyntax, ruleImageURLHTTPS, ruleDecoratedText, ruleDecoratedControl, ruleSwitchName, ruleSwitchType,
	ruleButtonListButtons, ruleButtonContent, ruleButtonTypeValue, ruleButtonTypeColor, ruleColorRange,
	ruleIconSource, ruleIconURLHTTPS, ruleMaterialName, ruleMaterialWeight, ruleMaterialGrade,
	ruleOnClickCount, ruleActionFunction, ruleActionLoad, ruleActionInteraction, ruleActionRequired, ruleActionRequiredName,
	ruleParameterKey, ruleParameterUnique, ruleOpenLinkURL, ruleOverflowItems, ruleOverflowItemText, ruleOverflowNested,
	ruleColumnsCount, ruleColumnWidgets, ruleColumnSizeStyle, ruleColumnVertical,
	ruleGridItems, ruleGridColumnCount, ruleGridItemContent, ruleGridItemLayout, ruleImageURIFormat,
	ruleCropTypeValue, ruleCropRatioValue, ruleCropRatioCustom, ruleBorderTypeValue, ruleBorderStroke, ruleBorderRadius,
	ruleCarouselCards, ruleCarouselWidgets, ruleChipListChips, ruleChipListLayout, ruleChipContent,
}

func oneOf[T comparable](value T, allowed ...T) bool {
	var zero T
	if value == zero {
		return true
	}
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}

func https(text string) bool { return validation.AbsoluteURI(text, "https") }

func alignment(c *validation.Check, value HorizontalAlignment, field string) {
	c.When(!oneOf(value, HorizontalAlignmentStart, HorizontalAlignmentCenter, HorizontalAlignmentEnd), ruleHorizontalAlign, field)
}

func imageType(c *validation.Check, value ImageType, field string) {
	c.When(!oneOf(value, ImageTypeSquare, ImageTypeCircle), ruleImageTypeValue, field)
}

// content checks a widget wrapper's content beneath its JSON property.
func content(c *validation.Check, field string, value WidgetContent) {
	if validation.Value(value) == nil {
		c.Fail(ruleWidgetContent, field)
		return
	}
	if field != "" {
		field += "."
	}
	c.Child(field+value.WidgetType(), value.check)
}

func (m Message) check(c *validation.Check) {
	ids := map[string]bool{}
	for i, card := range m.CardsV2 {
		field := fmt.Sprintf("cardsV2[%d]", i)
		c.When(len(m.CardsV2) > 1 && card.CardID == "", ruleCardIDRequired, field+".cardId")
		c.When(card.CardID != "" && ids[card.CardID], ruleCardIDUnique, field+".cardId")
		ids[card.CardID] = true
		c.Child(field+".card", card.Card.check)
	}
	if len(m.CardsV2) != 0 {
		if data, err := json.Marshal(m.CardsV2); err == nil && len(data) > 32*1024 {
			c.Fail(ruleMessageCardsSize, "cardsV2")
		}
	}
}

func (card Card) check(c *validation.Check) {
	c.When(card.Header == nil && len(card.Sections) == 0, ruleCardContent, "")
	c.When(!oneOf(card.SectionDividerStyle, DividerStyleSolid, DividerStyleNone), ruleCardDividerStyle, "sectionDividerStyle")
	total := 0
	for _, section := range card.Sections {
		for _, widget := range section.Widgets {
			total += countWidgets(widget.Content)
		}
	}
	c.When(total > 100, ruleCardWidgetCount, "sections")
	if card.Header != nil {
		c.Child("header", card.Header.check)
	}
	for i, section := range card.Sections {
		c.Child(fmt.Sprintf("sections[%d]", i), section.check)
	}
}

func countWidgets(value WidgetContent) int {
	switch value := validation.Value(value).(type) {
	case nil:
		return 0
	case Columns:
		count := 1
		for _, column := range value.ColumnItems {
			for _, widget := range column.Widgets {
				count += countWidgets(widget.Content)
			}
		}
		return count
	case Carousel:
		count := 1
		for _, card := range value.CarouselCards {
			for _, widget := range append(append([]NestedWidget{}, card.Widgets...), card.FooterWidgets...) {
				count += countWidgets(widget.Content)
			}
		}
		return count
	}
	return 1
}

func (h CardHeader) check(c *validation.Check) {
	c.When(h.Title == "", ruleHeaderTitle, "title")
	c.When(h.ImageURL != "" && !https(h.ImageURL), ruleHeaderImageHTTPS, "imageUrl")
	c.When(h.ImageURL == "" && (h.ImageType != "" || h.ImageAltText != ""), ruleHeaderImageURL, "imageUrl")
	imageType(c, h.ImageType, "imageType")
}

func (s Section) check(c *validation.Check) {
	c.When(len(s.Widgets) == 0, ruleSectionWidgets, "widgets")
	c.When(!s.Collapsible && (s.UncollapsibleWidgetsCount != 0 || s.CollapseControl != nil), ruleSectionCollapse, "collapsible")
	c.When(s.UncollapsibleWidgetsCount < 0 || s.UncollapsibleWidgetsCount > len(s.Widgets), ruleSectionUncollapse, "uncollapsibleWidgetsCount")
	for i, widget := range s.Widgets {
		c.Child(fmt.Sprintf("widgets[%d]", i), widget.check)
	}
	if s.CollapseControl != nil {
		c.Child("collapseControl", s.CollapseControl.check)
	}
}

func (cc CollapseControl) check(c *validation.Check) {
	alignment(c, cc.HorizontalAlignment, "horizontalAlignment")
	c.Child("expandButton", cc.ExpandButton.check)
	c.Child("collapseButton", cc.CollapseButton.check)
}

func (w Widget) check(c *validation.Check) {
	alignment(c, w.HorizontalAlignment, "horizontalAlignment")
	content(c, "", w.Content)
}

func (t TextParagraph) check(c *validation.Check) {
	c.When(t.MaxLines < 0, ruleParagraphMaxLines, "maxLines")
	c.When(!oneOf(t.TextSyntax, TextSyntaxHTML, TextSyntaxMarkdown), ruleParagraphSyntax, "textSyntax")
}

func (i Image) check(c *validation.Check) {
	c.When(!https(i.ImageURL), ruleImageURLHTTPS, "imageUrl")
	if i.OnClick != nil {
		c.Child("onClick", i.OnClick.check)
	}
}

func (Divider) check(*validation.Check) {}

func (t DecoratedText) check(c *validation.Check) {
	c.When(t.Text == "", ruleDecoratedText, "text")
	controls := 0
	for _, present := range []bool{t.Button != nil, t.SwitchControl != nil, t.EndIcon != nil} {
		if present {
			controls++
		}
	}
	c.When(controls > 1, ruleDecoratedControl, "")
	c.When(!oneOf(t.StartIconVerticalAlignment, VerticalAlignmentTop, VerticalAlignmentMiddle, VerticalAlignmentBottom), ruleVerticalAlign, "startIconVerticalAlignment")
	if t.StartIcon != nil {
		c.Child("startIcon", t.StartIcon.check)
	}
	for _, label := range []struct {
		field string
		text  *TextParagraph
	}{{"topLabelText", t.TopLabelText}, {"contentText", t.ContentText}, {"bottomLabelText", t.BottomLabelText}} {
		if label.text != nil {
			c.Child(label.field, label.text.check)
		}
	}
	if t.OnClick != nil {
		c.Child("onClick", t.OnClick.check)
	}
	if t.Button != nil {
		c.Child("button", t.Button.check)
	}
	if t.SwitchControl != nil {
		c.Child("switchControl", t.SwitchControl.check)
	}
	if t.EndIcon != nil {
		c.Child("endIcon", t.EndIcon.check)
	}
}

func (s SwitchControl) check(c *validation.Check) {
	c.When(s.Name == "", ruleSwitchName, "name")
	c.When(!oneOf(s.ControlType, SwitchControlTypeSwitch, SwitchControlTypeCheckBox), ruleSwitchType, "controlType")
	if s.OnChangeAction != nil {
		c.Child("onChangeAction", s.OnChangeAction.check)
	}
}

func (b ButtonList) check(c *validation.Check) {
	c.When(len(b.Buttons) == 0, ruleButtonListButtons, "buttons")
	for i, button := range b.Buttons {
		c.Child(fmt.Sprintf("buttons[%d]", i), button.check)
	}
}

func (b Button) check(c *validation.Check) {
	c.When(b.Text == "" && b.Icon == nil, ruleButtonContent, "")
	c.When(!oneOf(b.Type, ButtonTypeOutlined, ButtonTypeFilled, ButtonTypeFilledTonal, ButtonTypeBorderless), ruleButtonTypeValue, "type")
	c.When(b.Color != nil && b.Type != "" && b.Type != ButtonTypeFilled, ruleButtonTypeColor, "type")
	if b.Icon != nil {
		c.Child("icon", b.Icon.check)
	}
	if b.Color != nil {
		c.Child("color", b.Color.check)
	}
	c.Child("onClick", b.OnClick.check)
}

func (color Color) check(c *validation.Check) {
	for _, component := range []struct {
		field string
		value float64
	}{{"red", color.Red}, {"green", color.Green}, {"blue", color.Blue}} {
		c.When(!(component.value >= 0 && component.value <= 1), ruleColorRange, component.field)
	}
}

func (i Icon) check(c *validation.Check) {
	sources := 0
	for _, present := range []bool{i.KnownIcon != "", i.IconURL != "", i.MaterialIcon != nil} {
		if present {
			sources++
		}
	}
	c.When(sources != 1, ruleIconSource, "")
	c.When(i.IconURL != "" && !https(i.IconURL), ruleIconURLHTTPS, "iconUrl")
	imageType(c, i.ImageType, "imageType")
	if i.MaterialIcon != nil {
		c.Child("materialIcon", i.MaterialIcon.check)
	}
}

func (m MaterialIcon) check(c *validation.Check) {
	c.When(m.Name == "", ruleMaterialName, "name")
	c.When(!oneOf(m.Weight, 100, 200, 300, 400, 500, 600, 700), ruleMaterialWeight, "weight")
	c.When(m.Grade != -25 && m.Grade != 0 && m.Grade != 200, ruleMaterialGrade, "grade")
}

func (o OnClick) check(c *validation.Check) {
	actions := 0
	for _, present := range []bool{o.Action != nil, o.OpenLink != nil, o.OverflowMenu != nil} {
		if present {
			actions++
		}
	}
	c.When(actions != 1, ruleOnClickCount, "")
	if o.Action != nil {
		c.Child("action", o.Action.check)
	}
	if o.OpenLink != nil {
		c.Child("openLink", o.OpenLink.check)
	}
	if o.OverflowMenu != nil {
		c.Child("overflowMenu", o.OverflowMenu.check)
	}
}

func (a Action) check(c *validation.Check) {
	c.When(a.Function == "", ruleActionFunction, "function")
	c.When(!oneOf(a.LoadIndicator, LoadIndicatorSpinner, LoadIndicatorNone), ruleActionLoad, "loadIndicator")
	c.When(!oneOf(a.Interaction, InteractionOpenDialog), ruleActionInteraction, "interaction")
	c.When(a.AllWidgetsAreRequired && len(a.RequiredWidgets) != 0, ruleActionRequired, "requiredWidgets")
	for i, name := range a.RequiredWidgets {
		c.When(name == "", ruleActionRequiredName, fmt.Sprintf("requiredWidgets[%d]", i))
	}
	keys := map[string]bool{}
	for i, parameter := range a.Parameters {
		field := fmt.Sprintf("parameters[%d].key", i)
		c.When(parameter.Key == "", ruleParameterKey, field)
		c.When(parameter.Key != "" && keys[parameter.Key], ruleParameterUnique, field)
		keys[parameter.Key] = true
	}
}

func (o OpenLink) check(c *validation.Check) {
	c.When(!validation.AbsoluteURI(o.URL), ruleOpenLinkURL, "url")
}

func (o OverflowMenu) check(c *validation.Check) {
	c.When(len(o.Items) == 0, ruleOverflowItems, "items")
	for i, item := range o.Items {
		c.Child(fmt.Sprintf("items[%d]", i), item.check)
	}
}

func (i OverflowMenuItem) check(c *validation.Check) {
	c.When(i.Text == "", ruleOverflowItemText, "text")
	c.When(i.OnClick.OverflowMenu != nil, ruleOverflowNested, "onClick.overflowMenu")
	if i.StartIcon != nil {
		c.Child("startIcon", i.StartIcon.check)
	}
	c.Child("onClick", i.OnClick.check)
}

func (columns Columns) check(c *validation.Check) {
	c.When(len(columns.ColumnItems) == 0 || len(columns.ColumnItems) > 2, ruleColumnsCount, "columnItems")
	for i, column := range columns.ColumnItems {
		c.Child(fmt.Sprintf("columnItems[%d]", i), column.check)
	}
}

func (column Column) check(c *validation.Check) {
	c.When(len(column.Widgets) == 0, ruleColumnWidgets, "widgets")
	c.When(!oneOf(column.HorizontalSizeStyle, HorizontalSizeStyleFillAvailableSpace, HorizontalSizeStyleFillMinimumSpace), ruleColumnSizeStyle, "horizontalSizeStyle")
	alignment(c, column.HorizontalAlignment, "horizontalAlignment")
	c.When(!oneOf(column.VerticalAlignment, ColumnVerticalAlignmentCenter, ColumnVerticalAlignmentTop, ColumnVerticalAlignmentBottom), ruleColumnVertical, "verticalAlignment")
	for i, widget := range column.Widgets {
		content(c, fmt.Sprintf("widgets[%d]", i), widget.Content)
	}
}

func (g Grid) check(c *validation.Check) {
	c.When(len(g.Items) == 0, ruleGridItems, "items")
	c.When(g.ColumnCount < 0, ruleGridColumnCount, "columnCount")
	for i, item := range g.Items {
		c.Child(fmt.Sprintf("items[%d]", i), item.check)
	}
	if g.BorderStyle != nil {
		c.Child("borderStyle", g.BorderStyle.check)
	}
	if g.OnClick != nil {
		c.Child("onClick", g.OnClick.check)
	}
}

func (i GridItem) check(c *validation.Check) {
	c.When(i.Title == "" && i.Subtitle == "" && i.Image == nil, ruleGridItemContent, "")
	c.When(!oneOf(i.Layout, GridItemLayoutTextBelow, GridItemLayoutTextAbove), ruleGridItemLayout, "layout")
	if i.Image != nil {
		c.Child("image", i.Image.check)
	}
}

func (i ImageComponent) check(c *validation.Check) {
	c.When(!validation.AbsoluteURI(i.ImageURI, "http", "https"), ruleImageURIFormat, "imageUri")
	if i.CropStyle != nil {
		c.Child("cropStyle", i.CropStyle.check)
	}
	if i.BorderStyle != nil {
		c.Child("borderStyle", i.BorderStyle.check)
	}
}

func (s ImageCropStyle) check(c *validation.Check) {
	c.When(!oneOf(s.Type, ImageCropTypeSquare, ImageCropTypeCircle, ImageCropTypeRectangleCustom, ImageCropTypeRectangle4By3), ruleCropTypeValue, "type")
	finite := !math.IsNaN(s.AspectRatio) && !math.IsInf(s.AspectRatio, 0)
	if s.Type == ImageCropTypeRectangleCustom {
		c.When(!finite || s.AspectRatio <= 0, ruleCropRatioValue, "aspectRatio")
	} else {
		c.When(s.AspectRatio != 0, ruleCropRatioCustom, "aspectRatio")
	}
}

func (b BorderStyle) check(c *validation.Check) {
	c.When(!oneOf(b.Type, BorderTypeNone, BorderTypeStroke), ruleBorderTypeValue, "type")
	c.When(b.Type == BorderTypeNone && b.StrokeColor != nil, ruleBorderStroke, "strokeColor")
	c.When(b.CornerRadius < 0, ruleBorderRadius, "cornerRadius")
	if b.StrokeColor != nil {
		c.Child("strokeColor", b.StrokeColor.check)
	}
}

func (carousel Carousel) check(c *validation.Check) {
	c.When(len(carousel.CarouselCards) == 0, ruleCarouselCards, "carouselCards")
	for i, card := range carousel.CarouselCards {
		c.Child(fmt.Sprintf("carouselCards[%d]", i), card.check)
	}
}

func (card CarouselCard) check(c *validation.Check) {
	c.When(len(card.Widgets) == 0, ruleCarouselWidgets, "widgets")
	for i, widget := range card.Widgets {
		content(c, fmt.Sprintf("widgets[%d]", i), widget.Content)
	}
	for i, widget := range card.FooterWidgets {
		content(c, fmt.Sprintf("footerWidgets[%d]", i), widget.Content)
	}
}

func (l ChipList) check(c *validation.Check) {
	c.When(len(l.Chips) == 0, ruleChipListChips, "chips")
	c.When(!oneOf(l.Layout, ChipListLayoutWrapped, ChipListLayoutHorizontalScrollable), ruleChipListLayout, "layout")
	for i, chip := range l.Chips {
		c.Child(fmt.Sprintf("chips[%d]", i), chip.check)
	}
}

func (chip Chip) check(c *validation.Check) {
	c.When(chip.Label == "" && chip.Icon == nil, ruleChipContent, "")
	if chip.Icon != nil {
		c.Child("icon", chip.Icon.check)
	}
	if chip.OnClick != nil {
		c.Child("onClick", chip.OnClick.check)
	}
}
