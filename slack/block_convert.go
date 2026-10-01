package slack

import "github.com/JSYoo5B/convertago/internal/conversion"

func convertBlock(node conversion.Node) (Block, error) {
	r := conversion.Reader{Platform: "slack", Node: node}
	id := r.Text("block_id", 1, 255)
	var block Block
	switch node.Role {
	case "header":
		text, err := slotText(node, "text", 150, true)
		if err != nil {
			return nil, err
		}
		block = HeaderBlock{Text: text.(PlainTextObject), BlockID: id, Level: r.Int("level", 1, 4)}
	case "section":
		result := SectionBlock{BlockID: id, Expand: r.Bool("expand")}
		if node.Has("text") {
			text, err := slotText(node, "text", 3000, false)
			if err != nil {
				return nil, err
			}
			result.Text = text
		}
		for _, child := range r.Count("fields", 0, 10) {
			text, err := convertTextObject(child, 2000)
			if err != nil {
				return nil, err
			}
			result.Fields = append(result.Fields, text)
		}
		if result.Text == nil && len(result.Fields) == 0 {
			r.Fail("missing_input", "section requires text or fields")
		}
		for _, child := range node.Children("accessory") {
			if child.Role == "button" {
				button, err := convertButton(child)
				if err != nil {
					return nil, err
				}
				result.Accessory = button
			} else {
				image, err := convertImage(child, true)
				if err != nil {
					return nil, err
				}
				result.Accessory = ImageElement{ImageURL: image.ImageURL, SlackFile: image.SlackFile, AltText: image.AltText}
			}
		}
		block = result
	case "image":
		return convertImage(node, false)
	case "divider":
		block = DividerBlock{BlockID: id}
	case "markdown":
		block = MarkdownBlock{Text: r.Text("text", 1, 12000)}
	case "actions":
		result := ActionsBlock{BlockID: id}
		for _, child := range r.Count("elements", 1, 25) {
			button, err := convertButton(child)
			if err != nil {
				return nil, err
			}
			result.Elements = append(result.Elements, button)
		}
		block = result
	case "context":
		result := ContextBlock{BlockID: id}
		for _, child := range r.Count("elements", 1, 10) {
			if child.Role == "image" {
				image, err := convertImage(child, true)
				if err != nil {
					return nil, err
				}
				result.Elements = append(result.Elements, ImageElement{ImageURL: image.ImageURL, SlackFile: image.SlackFile, AltText: image.AltText})
			} else {
				text, err := convertTextObject(child, 3000)
				if err != nil {
					return nil, err
				}
				result.Elements = append(result.Elements, text.(ContextElement))
			}
		}
		block = result
	case "video":
		title, err := slotText(node, "title", 199, true)
		if err != nil {
			return nil, err
		}
		result := VideoBlock{Title: title.(PlainTextObject), BlockID: id, VideoURL: r.URL("video_url", true), AltText: r.Text("alt", 1, 0), ThumbnailURL: r.URL("thumbnail_url", false), TitleURL: r.URL("title_url", true), AuthorName: r.Text("author_name", 1, 49), ProviderName: r.Text("provider_name", 1, 0), ProviderIconURL: r.URL("provider_icon_url", false)}
		if node.Has("description") {
			text, err := slotText(node, "description", 199, true)
			if err != nil {
				return nil, err
			}
			value := text.(PlainTextObject)
			result.Description = &value
		}
		block = result
	default:
		r.Fail("invalid_tag", "unsupported top-level block")
	}
	return block, r.Err
}

func convertTextObject(node conversion.Node, max int) (TextObject, error) {
	r := conversion.Reader{Platform: "slack", Node: node}
	text := r.Text("text", 1, max)
	if node.Role == "plain_text" {
		result := PlainTextObject{Text: text}
		if node.Has("emoji") {
			value := r.Bool("emoji")
			result.Emoji = &value
		}
		return result, r.Err
	}
	return MrkdwnTextObject{Text: text, Verbatim: r.Bool("verbatim")}, r.Err
}

func slotText(node conversion.Node, slot string, max int, plain bool) (TextObject, error) {
	children := node.Children(slot)
	var text string
	format := ""
	count := 0
	for _, part := range node.Parts {
		if part.Slot == slot {
			count++
			current := part.Format
			if current == "" {
				current = "plain"
			}
			if format != "" && current != format {
				return nil, conversion.Error("slack", part.Path, "conflicting_format", "text slot cannot mix plain and mrkdwn")
			}
			format = current
			text += part.Text
		}
	}
	if len(children) != 0 {
		if count != 0 || len(children) != 1 {
			return nil, conversion.Error("slack", node.Path, "conflicting_input", "text slot requires one object or scalar parts")
		}
		if plain && children[0].Role != "plain_text" {
			return nil, conversion.Error("slack", children[0].Path, "invalid_tag", "plain text is required")
		}
		return convertTextObject(children[0], max)
	}
	if err := conversion.ValidateText("slack", node.Path, text, 1, max); err != nil {
		return nil, err
	}
	if format == "mrkdwn" {
		if plain {
			return nil, conversion.Error("slack", node.Path, "invalid_tag", "plain text is required")
		}
		return MrkdwnTextObject{Text: text, Verbatim: true}, nil
	}
	return PlainTextObject{Text: text}, nil
}

func convertImage(node conversion.Node, element bool) (ImageBlock, error) {
	r := conversion.Reader{Platform: "slack", Node: node}
	imageURL := r.URL("url", false)
	r.Text("url", 1, 3000)
	alt := r.Text("alt", 1, 2000)
	var file *SlackFileObject
	if children := node.Children("slack_file"); len(children) != 0 {
		child := children[0]
		f := conversion.Reader{Platform: "slack", Node: child}
		value := SlackFileObject{URL: f.URL("url", false), ID: f.Text("id", 1, 0)}
		if child.Has("url") == child.Has("id") {
			f.Fail("conflicting_input", "slack_file requires exactly one of url or id")
		}
		if f.Err != nil {
			return ImageBlock{}, f.Err
		}
		file = &value
	}
	if node.Has("url") == (file != nil) {
		r.Fail("conflicting_input", "image requires exactly one of url or slack_file")
	}
	if element && (node.Has("title") || node.Has("block_id")) {
		r.Fail("invalid_tag", "image elements do not accept title or block_id")
	}
	result := ImageBlock{ImageURL: imageURL, SlackFile: file, AltText: alt, BlockID: r.Text("block_id", 1, 255)}
	if node.Has("title") {
		text, err := slotText(node, "title", 2000, true)
		if err != nil {
			return ImageBlock{}, err
		}
		title := text.(PlainTextObject)
		result.Title = &title
	}
	return result, r.Err
}

func convertButton(node conversion.Node) (ButtonElement, error) {
	r := conversion.Reader{Platform: "slack", Node: node}
	text, err := slotText(node, "text", 75, true)
	if err != nil {
		return ButtonElement{}, err
	}
	result := ButtonElement{Text: text.(PlainTextObject), ActionID: r.Text("action_id", 1, 255), URL: r.URI("url"), Value: r.Text("value", 0, 2000), Style: ButtonStyle(r.Enum("style", "primary", "danger")), AccessibilityLabel: r.Text("accessibility_label", 1, 75), AgentPrompt: r.Text("agent_prompt", 1, 4000), AgentPromptDisplay: r.Text("agent_prompt_display", 1, 0)}
	r.Text("url", 1, 3000)
	for _, part := range node.Parts {
		if part.Slot == "visible_to_user_ids" {
			if err := conversion.ValidateText("slack", part.Path, part.Text, 1, 0); err != nil {
				return ButtonElement{}, err
			}
			result.VisibleToUserIDs = append(result.VisibleToUserIDs, part.Text)
		}
	}
	if children := node.Children("confirm"); len(children) != 0 {
		child := children[0]
		c := conversion.Reader{Platform: "slack", Node: child}
		value := ConfirmationDialogObject{Style: ButtonStyle(c.Enum("style", "primary", "danger"))}
		for _, slot := range []struct {
			name  string
			max   int
			plain bool
		}{{"title", 100, true}, {"text", 300, false}, {"confirm", 30, true}, {"deny", 30, true}} {
			text, err := slotText(child, slot.name, slot.max, slot.plain)
			if err != nil {
				return ButtonElement{}, err
			}
			switch slot.name {
			case "title":
				value.Title = text.(PlainTextObject)
			case "text":
				value.Text = text
			case "confirm":
				value.Confirm = text.(PlainTextObject)
			case "deny":
				value.Deny = text.(PlainTextObject)
			}
		}
		if c.Err != nil {
			return ButtonElement{}, c.Err
		}
		result.Confirm = &value
	}
	return result, r.Err
}
