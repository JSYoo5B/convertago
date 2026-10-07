package googlechat

import "github.com/JSYoo5B/convertago/internal/conversion"

func addLayoutRoles(roles map[string]conversion.Role) {
	child := func(names []string, required, repeated bool) conversion.Slot {
		return conversion.Slot{Children: names, Required: required, Repeated: repeated}
	}
	widgets := []string{"textParagraph", "image", "decoratedText", "buttonList", "divider", "columns", "grid", "carousel", "chipList"}
	clicks := []string{"onClick", "action", "openLink", "overflowMenu"}
	roles["cardWithId"] = conversion.Role{Slots: map[string]conversion.Slot{"cardId": {}, "card": child([]string{"card"}, true, false)}}
	roles["card"] = conversion.Role{Slots: map[string]conversion.Slot{"header": child([]string{"header"}, false, false), "sections": child([]string{"section"}, false, true), "widgets": child(widgets, false, true), "sectionDividerStyle": enum(dividerStyles, ruleCardDividerStyle)}}
	roles["section"] = conversion.Role{DefaultSlot: "header", DefaultChildSlot: "widgets", Slots: map[string]conversion.Slot{"header": {Repeated: true}, "widgets": child(widgets, true, true), "collapsible": flag, "uncollapsibleWidgetsCount": {}, "collapseControl": child([]string{"collapseControl"}, false, false)}, Formats: []string{"plain", "html"}}
	roles["collapseControl"] = conversion.Role{NestedOnly: true, Slots: map[string]conversion.Slot{"horizontalAlignment": enum(horizontalAlignments, ruleHorizontalAlign), "expandButton": child([]string{"button"}, true, false), "collapseButton": child([]string{"button"}, true, false)}}
	roles["columns"] = conversion.Role{DefaultChildSlot: "columnItems", Slots: map[string]conversion.Slot{"columnItems": child([]string{"column"}, true, true)}}
	roles["column"] = conversion.Role{NestedOnly: true, DefaultChildSlot: "widgets", Slots: map[string]conversion.Slot{"widgets": child([]string{"textParagraph", "image", "decoratedText", "buttonList", "chipList"}, true, true), "horizontalSizeStyle": enum(horizontalSizeStyles, ruleColumnSizeStyle), "horizontalAlignment": enum(horizontalAlignments, ruleHorizontalAlign), "verticalAlignment": enum(columnVerticalAlignments, ruleColumnVertical)}}
	roles["grid"] = conversion.Role{DefaultChildSlot: "items", Slots: map[string]conversion.Slot{"title": {}, "items": child([]string{"gridItem"}, true, true), "borderStyle": child([]string{"borderStyle"}, false, false), "columnCount": {}, "onClick": child(clicks, false, false)}}
	roles["gridItem"] = conversion.Role{NestedOnly: true, DefaultSlot: "title", Slots: map[string]conversion.Slot{"id": {}, "image": child([]string{"imageComponent"}, false, false), "title": {}, "subtitle": {}, "layout": enum(gridItemLayouts, ruleGridItemLayout)}}
	roles["imageComponent"] = conversion.Role{NestedOnly: true, DefaultSlot: "imageUri", Slots: map[string]conversion.Slot{"imageUri": {Required: true}, "altText": {}, "cropStyle": child([]string{"imageCropStyle"}, false, false), "borderStyle": child([]string{"borderStyle"}, false, false)}}
	roles["imageCropStyle"] = conversion.Role{NestedOnly: true, EmptyAllowed: true, Slots: map[string]conversion.Slot{"type": enum(imageCropTypes, ruleCropTypeValue), "aspectRatio": {}}}
	roles["borderStyle"] = conversion.Role{NestedOnly: true, EmptyAllowed: true, Slots: map[string]conversion.Slot{"type": enum(borderTypes, ruleBorderTypeValue), "strokeColor": child([]string{"color"}, false, false), "cornerRadius": {}}}
	roles["carousel"] = conversion.Role{DefaultChildSlot: "carouselCards", Slots: map[string]conversion.Slot{"carouselCards": child([]string{"carouselCard"}, true, true)}}
	roles["carouselCard"] = conversion.Role{NestedOnly: true, DefaultChildSlot: "widgets", Slots: map[string]conversion.Slot{"widgets": child([]string{"textParagraph", "image", "buttonList"}, true, true), "footerWidgets": child([]string{"textParagraph", "image", "buttonList"}, false, true)}}
	roles["chipList"] = conversion.Role{DefaultChildSlot: "chips", Slots: map[string]conversion.Slot{"layout": enum(chipListLayouts, ruleChipListLayout), "chips": child([]string{"chip"}, true, true)}}
	roles["chip"] = conversion.Role{NestedOnly: true, DefaultSlot: "label", Slots: map[string]conversion.Slot{"icon": child([]string{"icon"}, false, false), "label": {}, "onClick": child(clicks, false, false), "disabled": flag, "altText": {}}}
	for _, name := range []string{"columns", "grid", "carousel", "chipList"} {
		role := roles[name]
		role.Slots["horizontalAlignment"] = enum(horizontalAlignments, ruleHorizontalAlign)
		roles[name] = role
	}
}
