package kakaowork

import "github.com/JSYoo5B/convertago/internal/conversion"

func init() {
	text := conversion.Slot{Repeated: true, Required: true}
	actions := []string{"open_system_browser", "open_inapp_browser", "open_external_app", "submit_action", "call_modal", "exclusive"}
	child := func(roles []string, required bool) conversion.Slot {
		return conversion.Slot{Children: roles, Required: required}
	}
	roles := map[string]conversion.Role{
		"header":      {DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text, "style": {}}},
		"text":        {DefaultSlot: "text", DefaultChildSlot: "text", Slots: map[string]conversion.Slot{"text": {Repeated: true, Required: true, Scalar: true, Children: []string{"styled", "link", "mention"}}}, Styles: []string{"bold", "italic", "strike"}, Formats: []string{"plain"}},
		"image_link":  {DefaultSlot: "url", Slots: map[string]conversion.Slot{"url": {Required: true}}},
		"preview":     {DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text}},
		"divider":     {EmptyAllowed: true},
		"styled":      {NestedOnly: true, DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text, "color": {}, "bold": {}, "italic": {}, "strike": {}}, Styles: []string{"bold", "italic", "strike"}, Formats: []string{"plain"}},
		"link":        {NestedOnly: true, DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text, "url": {Required: true}}},
		"mention":     {NestedOnly: true, DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text, "user_id": {Required: true}}},
		"button":      {DefaultSlot: "text", DefaultChildSlot: "action", Slots: map[string]conversion.Slot{"text": text, "style": {}, "action": child(actions, true)}},
		"action":      {DefaultChildSlot: "elements", Slots: map[string]conversion.Slot{"elements": {Children: []string{"button"}, Repeated: true, Required: true}}},
		"description": {DefaultChildSlot: "content", Slots: map[string]conversion.Slot{"content": child([]string{"text"}, true), "term": {Required: true}, "accent": {}}},
		"section":     {DefaultChildSlot: "content", Slots: map[string]conversion.Slot{"content": child([]string{"text"}, true), "accessory": child([]string{"image_link"}, false), "action": child(actions, false)}},
		"context":     {DefaultChildSlot: "content", Slots: map[string]conversion.Slot{"content": child([]string{"text"}, true), "image": child([]string{"image_link"}, true)}},
	}
	for _, name := range actions[:5] {
		slots := map[string]conversion.Slot{"name": {}, "value": {Required: true}}
		if name == "submit_action" {
			slots["name"] = conversion.Slot{Required: true}
		}
		if name == "open_inapp_browser" {
			slots["standalone"] = conversion.Slot{}
			slots["width"] = conversion.Slot{}
			slots["height"] = conversion.Slot{}
		}
		roles[name] = conversion.Role{NestedOnly: true, DefaultSlot: "value", Slots: slots}
	}
	slots := map[string]conversion.Slot{}
	for _, name := range []string{"default", "pc", "mobile", "windows", "macos", "android", "ios"} {
		slots[name] = child(actions[:5], name == "default")
	}
	roles["exclusive"] = conversion.Role{NestedOnly: true, DefaultChildSlot: "default", Slots: slots}
	conversion.Register(conversion.Profile{Platform: "kakaowork", Roles: roles})
}
