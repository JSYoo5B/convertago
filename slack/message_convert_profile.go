package slack

import "github.com/JSYoo5B/convertago/internal/conversion"

func init() {
	text := conversion.Slot{Repeated: true, Required: true}
	objects := []string{"plain_text", "mrkdwn"}
	child := func(names []string, required, repeated bool) conversion.Slot {
		return conversion.Slot{Children: names, Required: required, Repeated: repeated}
	}
	roles := map[string]conversion.Role{
		"text":       {DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text}, Formats: []string{"plain"}},
		"header":     {DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": {Required: true, Repeated: true, Scalar: true, Children: []string{"plain_text"}}, "block_id": {}, "level": {}}, Formats: []string{"plain"}},
		"section":    {DefaultSlot: "text", DefaultChildSlot: "text", Slots: map[string]conversion.Slot{"text": {Repeated: true, Scalar: true, Children: objects}, "fields": child(objects, false, true), "accessory": child([]string{"button", "image"}, false, false), "block_id": {}, "expand": {}}, Formats: []string{"plain", "mrkdwn"}},
		"rich_text":  {DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text}, Styles: []string{"bold", "italic", "strike", "code", "underline"}, Formats: []string{"plain"}},
		"image":      {DefaultSlot: "url", Slots: map[string]conversion.Slot{"url": {}, "slack_file": child([]string{"slack_file"}, false, false), "alt": {Required: true}, "title": {Scalar: true, Children: []string{"plain_text"}}, "block_id": {}}},
		"slack_file": {NestedOnly: true, Slots: map[string]conversion.Slot{"url": {}, "id": {}}},
		"plain_text": {NestedOnly: true, DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text, "emoji": {}}, Formats: []string{"plain"}},
		"mrkdwn":     {NestedOnly: true, DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text, "verbatim": {}}, Formats: []string{"mrkdwn"}},
		"actions":    {DefaultChildSlot: "elements", Slots: map[string]conversion.Slot{"elements": child([]string{"button"}, true, true), "block_id": {}}},
		"context":    {DefaultChildSlot: "elements", Slots: map[string]conversion.Slot{"elements": child([]string{"plain_text", "mrkdwn", "image"}, true, true), "block_id": {}}},
		"divider":    {EmptyAllowed: true, Slots: map[string]conversion.Slot{"block_id": {}}},
		"markdown":   {DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text}, Formats: []string{"markdown"}},
		"video":      {DefaultSlot: "video_url", Slots: map[string]conversion.Slot{"title": {Required: true, Scalar: true, Children: []string{"plain_text"}}, "video_url": {Required: true}, "alt": {Required: true}, "thumbnail_url": {Required: true}, "title_url": {}, "description": {Scalar: true, Children: []string{"plain_text"}}, "author_name": {}, "provider_name": {}, "provider_icon_url": {}, "block_id": {}}},
		"button":     {NestedOnly: true, DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": {Required: true, Scalar: true, Children: []string{"plain_text"}}, "action_id": {}, "url": {}, "value": {}, "style": {}, "confirm": child([]string{"confirm"}, false, false), "accessibility_label": {}, "agent_prompt": {}, "agent_prompt_display": {}, "visible_to_user_ids": {Repeated: true}}},
		"confirm":    {NestedOnly: true, DefaultSlot: "text", Slots: map[string]conversion.Slot{"title": {Required: true, Scalar: true, Children: []string{"plain_text"}}, "text": {Required: true, Repeated: true, Scalar: true, Children: objects}, "confirm": {Required: true, Scalar: true, Children: []string{"plain_text"}}, "deny": {Required: true, Scalar: true, Children: []string{"plain_text"}}, "style": {}}, Formats: []string{"plain", "mrkdwn"}},
	}
	styles := []string{"bold", "italic", "strike", "code", "underline", "highlight", "client_highlight", "unlink"}
	inline := []string{"text", "link", "user", "emoji"}
	inlineSlot := conversion.Slot{Repeated: true, Required: true, Scalar: true, Children: inline}
	roles["text"] = conversion.Role{DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": text}, Styles: styles, Formats: []string{"plain"}}
	roles["rich_text"] = conversion.Role{DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": {Repeated: true, Scalar: true, Children: inline}, "elements": child([]string{"rich_text_section", "rich_text_list", "rich_text_preformatted", "rich_text_quote"}, false, true), "block_id": {}}, Styles: styles, Formats: []string{"plain"}}
	roles["rich_text_section"] = conversion.Role{NestedOnly: true, DefaultSlot: "text", DefaultChildSlot: "text", Slots: map[string]conversion.Slot{"text": inlineSlot}, Styles: styles, Formats: []string{"plain"}}
	roles["rich_text_quote"] = conversion.Role{NestedOnly: true, DefaultSlot: "text", DefaultChildSlot: "text", Slots: map[string]conversion.Slot{"text": inlineSlot, "border": {}}, Styles: styles, Formats: []string{"plain"}}
	roles["rich_text_preformatted"] = conversion.Role{NestedOnly: true, DefaultSlot: "text", DefaultChildSlot: "text", Slots: map[string]conversion.Slot{"text": {Repeated: true, Required: true, Scalar: true, Children: []string{"text", "link"}}, "border": {}, "language": {}}, Styles: styles, Formats: []string{"plain"}}
	roles["rich_text_list"] = conversion.Role{NestedOnly: true, DefaultChildSlot: "elements", Slots: map[string]conversion.Slot{"elements": child([]string{"rich_text_section"}, true, true), "style": {Required: true}, "indent": {}, "offset": {}, "border": {}}}
	roles["link"] = conversion.Role{NestedOnly: true, DefaultSlot: "url", Slots: map[string]conversion.Slot{"url": {Required: true}, "text": {Repeated: true, Styles: []string{"bold", "italic", "strike", "underline", "highlight", "client_highlight", "unlink"}}, "unsafe": {}, "from_llm": {}, "is_slack_url": {}, "truncated": {}}, Styles: []string{"bold", "italic", "strike", "underline", "highlight", "client_highlight", "unlink"}}
	roles["user"] = conversion.Role{NestedOnly: true, DefaultSlot: "user_id", Slots: map[string]conversion.Slot{"user_id": {Required: true}, "from_llm": {}}, Styles: roles["link"].Styles}
	roles["emoji"] = conversion.Role{NestedOnly: true, DefaultSlot: "name", Slots: map[string]conversion.Slot{"name": {Required: true}, "unicode": {}}}
	conversion.Register(conversion.Profile{Platform: "slack", Roles: roles})
}
