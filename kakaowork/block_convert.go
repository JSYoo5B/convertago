package kakaowork

import "github.com/JSYoo5B/convertago/internal/conversion"

func convertBlock(node conversion.Node) (BubbleBlock, error) {
	r := conversion.Reader{Platform: "kakaowork", Node: node}
	var block BubbleBlock
	switch node.Role {
	case "text":
		return convertText(node)
	case "image_link":
		block = ImageBlock{Url: r.URL("url", false)}
	case "divider":
		block = DividerBlock{}
	case "button":
		action, err := convertAction(node.Children("action")[0])
		if err != nil {
			return nil, err
		}
		style := ButtonStyle(r.Enum("style", "default", "primary", "danger"))
		if style == "" {
			style = ButtonStyleDefault
		}
		block = ButtonBlock{Text: r.Text("text", 1, 20), Style: style, Action: action}
	case "action":
		result := ActionBlock{}
		for _, child := range r.Count("elements", 2, 3) {
			button, err := convertBlock(child)
			if err != nil {
				return nil, err
			}
			result.Elements = append(result.Elements, button.(ButtonBlock))
		}
		block = result
	case "description", "section", "context":
		content, err := convertText(node.Children("content")[0])
		if err != nil {
			return nil, err
		}
		switch node.Role {
		case "description":
			block = DescriptionBlock{Content: content, Term: r.Text("term", 1, 10), Accent: r.Bool("accent")}
		case "section":
			result := SectionBlock{Content: content}
			if children := node.Children("accessory"); len(children) != 0 {
				image, err := convertBlock(children[0])
				if err != nil {
					return nil, err
				}
				value := image.(ImageBlock)
				result.Accessory = &value
			}
			if children := node.Children("action"); len(children) != 0 {
				action, err := convertAction(children[0])
				if err != nil {
					return nil, err
				}
				result.Action = action
			}
			block = result
		case "context":
			image, err := convertBlock(node.Children("image")[0])
			if err != nil {
				return nil, err
			}
			block = ContextBlock{Content: content, Image: image.(ImageBlock)}
		}
	default:
		r.Fail("invalid_tag", "unsupported top-level block")
	}
	return block, r.Err
}

func convertText(node conversion.Node) (TextBlock, error) {
	result := TextBlock{}
	styled := false
	for _, input := range node.Inputs {
		var inline Inline
		if input.Part != nil {
			part := input.Part
			item := InlineStyled{Text: part.Text}
			for _, flag := range part.Style {
				switch flag {
				case "bold":
					item.Bold = true
				case "italic":
					item.Italic = true
				case "strike":
					item.Strike = true
				}
			}
			styled = styled || len(part.Style) != 0
			inline = item
		} else {
			child := *input.Child
			r := conversion.Reader{Platform: "kakaowork", Node: child}
			text := r.Text("text", 0, 500)
			switch child.Role {
			case "styled":
				item := InlineStyled{Text: text, Color: InlineColor(r.Enum("color", "default", "red", "blue", "grey")), Bold: r.Bool("bold"), Italic: r.Bool("italic"), Strike: r.Bool("strike")}
				for _, part := range child.Parts {
					if part.Slot == "text" {
						for _, flag := range part.Style {
							switch flag {
							case "bold":
								item.Bold = true
							case "italic":
								item.Italic = true
							case "strike":
								item.Strike = true
							}
						}
					}
				}
				inline = item
			case "link":
				inline = InlineLink{Text: text, Url: r.URI("url", "http", "https", "mailto", "tel")}
			case "mention":
				inline = InlineMention{Text: text, UserId: r.Int("user_id", 1, 0)}
			}
			if r.Err != nil {
				return TextBlock{}, r.Err
			}
			styled = true
		}
		result.Text += inline.String()
		result.Inlines = append(result.Inlines, inline)
	}
	if err := conversion.ValidateText("kakaowork", node.Path, result.Text, 0, 500); err != nil {
		return TextBlock{}, err
	}
	if !styled {
		result.Inlines = nil
	}
	return result, nil
}

func convertAction(node conversion.Node) (ButtonAction, error) {
	r := conversion.Reader{Platform: "kakaowork", Node: node}
	name := r.Text("name", 0, 0)
	var result ButtonAction
	switch node.Role {
	case "open_system_browser":
		result = OpenSystemBrowserAction{Name: name, Value: r.URL("value", false)}
	case "open_inapp_browser":
		action := OpenInAppBrowserAction{Name: name, Value: r.URL("value", false), Standalone: r.Bool("standalone"), Width: r.Int("width", 1, 0), Height: r.Int("height", 1, 0)}
		if !action.Standalone && (node.Has("width") || node.Has("height")) {
			r.Fail("conflicting_input", "width and height require standalone")
		}
		result = action
	case "open_external_app":
		result = OpenExternalAppAction{Name: name, Value: r.URI("value")}
	case "submit_action":
		result = SubmitAction{Name: r.Text("name", 1, 0), Value: r.Text("value", 0, 0)}
	case "call_modal":
		result = CallModalAction{Name: name, Value: r.Text("value", 0, 0)}
	case "exclusive":
		action := ExclusiveAction{}
		for _, slot := range []struct {
			name   string
			target *ButtonAction
		}{{"default", &action.Default}, {"pc", &action.Pc}, {"mobile", &action.Mobile}, {"windows", &action.Windows}, {"macos", &action.MacOs}, {"android", &action.Android}, {"ios", &action.Ios}} {
			if children := node.Children(slot.name); len(children) != 0 {
				value, err := convertAction(children[0])
				if err != nil {
					return nil, err
				}
				*slot.target = value
			}
		}
		result = action
	default:
		r.Fail("invalid_tag", "unsupported button action")
	}
	return result, r.Err
}
