package googlechat

import "github.com/JSYoo5B/convertago/internal/conversion"

func init() {
	text := conversion.Slot{Repeated: true, Required: true}
	child := func(names []string, required, repeated bool) conversion.Slot {
		return conversion.Slot{Children: names, Required: required, Repeated: repeated}
	}
	click := []string{"onClick", "action", "openLink", "overflowMenu"}
	roles := map[string]conversion.Role{
		"header":           {DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text, "subtitle": {Repeated: true}, "url": {}, "alt": {}, "imageType": {}}, Formats: []string{"plain"}},
		"textParagraph":    {DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text, "maxLines": {}}, Styles: []string{"bold", "italic", "strike", "code", "underline"}, Formats: []string{"plain", "html", "markdown"}, FormatStyles: map[string][]string{"markdown": {"bold", "italic", "strike", "code"}}},
		"image":            {DefaultSlot: "url", Slots: map[string]conversion.Slot{"url": {Required: true}, "alt": {}, "onClick": child(click, false, false)}},
		"fallbackText":     {DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text}, Formats: []string{"plain"}},
		"text":             {DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text}, Formats: []string{"plain"}},
		"divider":          {EmptyAllowed: true},
		"buttonList":       {DefaultChildSlot: "buttons", Slots: map[string]conversion.Slot{"buttons": child([]string{"button"}, true, true)}},
		"button":           {NestedOnly: true, DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": {Repeated: true}, "icon": child([]string{"icon"}, false, false), "color": child([]string{"color"}, false, false), "onClick": child(click, true, false), "disabled": {}, "altText": {}, "type": {}}},
		"color":            {NestedOnly: true, EmptyAllowed: true, Slots: map[string]conversion.Slot{"red": {}, "green": {}, "blue": {}}},
		"onClick":          {NestedOnly: true, Slots: map[string]conversion.Slot{"action": child([]string{"action"}, false, false), "openLink": child([]string{"openLink"}, false, false), "overflowMenu": child([]string{"overflowMenu"}, false, false)}},
		"openLink":         {NestedOnly: true, DefaultSlot: "url", Slots: map[string]conversion.Slot{"url": {Required: true}}},
		"action":           {NestedOnly: true, DefaultSlot: "function", Slots: map[string]conversion.Slot{"function": {Required: true}, "parameters": child([]string{"actionParameter"}, false, true), "loadIndicator": {}, "persistValues": {}, "interaction": {}, "requiredWidgets": {Repeated: true}, "allWidgetsAreRequired": {}}},
		"actionParameter":  {NestedOnly: true, Slots: map[string]conversion.Slot{"key": {Required: true}, "value": {Required: true}}},
		"overflowMenu":     {NestedOnly: true, DefaultChildSlot: "items", Slots: map[string]conversion.Slot{"items": child([]string{"overflowMenuItem"}, true, true)}},
		"overflowMenuItem": {NestedOnly: true, DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text, "startIcon": child([]string{"icon"}, false, false), "onClick": child([]string{"onClick", "action", "openLink"}, true, false), "disabled": {}}},
		"icon":             {NestedOnly: true, Slots: map[string]conversion.Slot{"knownIcon": {}, "iconUrl": {}, "materialIcon": child([]string{"materialIcon"}, false, false), "altText": {}, "imageType": {}}},
		"materialIcon":     {NestedOnly: true, DefaultSlot: "name", Slots: map[string]conversion.Slot{"name": {Required: true}, "fill": {}, "weight": {}, "grade": {}}},
		"switchControl":    {NestedOnly: true, DefaultSlot: "name", Slots: map[string]conversion.Slot{"name": {Required: true}, "value": {}, "selected": {}, "onChangeAction": child([]string{"action"}, false, false), "controlType": {}}},
		"decoratedText":    {DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text, "startIcon": child([]string{"icon"}, false, false), "startIconVerticalAlignment": {}, "topLabel": {Repeated: true}, "topLabelText": child([]string{"textParagraph"}, false, false), "contentText": child([]string{"textParagraph"}, false, false), "wrapText": {}, "bottomLabel": {Repeated: true}, "bottomLabelText": child([]string{"textParagraph"}, false, false), "onClick": child(click, false, false), "button": child([]string{"button"}, false, false), "switchControl": child([]string{"switchControl"}, false, false), "endIcon": child([]string{"icon"}, false, false)}, Styles: []string{"bold", "italic", "strike", "code", "underline"}, Formats: []string{"plain", "html"}},
	}
	for _, name := range []string{"textParagraph", "image", "divider", "buttonList", "decoratedText"} {
		role := roles[name]
		if role.Slots == nil {
			role.Slots = map[string]conversion.Slot{}
		}
		role.Slots["horizontalAlignment"] = conversion.Slot{}
		roles[name] = role
	}
	addLayoutRoles(roles)
	conversion.Register(conversion.Profile{Platform: "googlechat", Roles: roles})
}
