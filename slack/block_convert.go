package slack

import (
	"github.com/JSYoo5B/convertago/internal/conversion"
	"github.com/JSYoo5B/convertago/internal/validation"
)

type converter struct {
	conversion.Checker
}

// checked runs the rules of a value built from node.
func checked[T interface{ check(*validation.Check) }](c converter, node conversion.Node, value T, err error) (T, error) {
	if err != nil {
		return value, err
	}
	return value, c.Check(node, value.check)
}

func (c converter) block(node conversion.Node) (Block, error) {
	r := conversion.Reader{Platform: "slack", Node: node}
	id := r.String("block_id")
	switch node.Role {
	case "header":
		text, err := c.text(node, "text", true)
		if err != nil {
			return nil, err
		}
		return checked[Block](c, node, HeaderBlock{Text: text.(PlainTextObject), BlockID: id, Level: r.ParseInt("level")}, r.Err)
	case "section":
		result := SectionBlock{BlockID: id, Expand: r.Bool("expand")}
		if node.Has("text") {
			text, err := c.text(node, "text", false)
			if err != nil {
				return nil, err
			}
			result.Text = text
		}
		for _, child := range node.Children("fields") {
			text, err := c.textObject(child)
			if err != nil {
				return nil, err
			}
			result.Fields = append(result.Fields, text)
		}
		for _, child := range node.Children("accessory") {
			var err error
			if child.Role == "button" {
				result.Accessory, err = c.button(child)
			} else {
				result.Accessory, err = c.imageElement(child)
			}
			if err != nil {
				return nil, err
			}
		}
		return checked[Block](c, node, result, r.Err)
	case "image":
		file, err := c.slackFile(node)
		if err != nil {
			return nil, err
		}
		result := ImageBlock{ImageURL: r.String("url"), SlackFile: file, AltText: r.String("alt"), BlockID: id}
		if node.Has("title") {
			text, err := c.text(node, "title", true)
			if err != nil {
				return nil, err
			}
			title := text.(PlainTextObject)
			result.Title = &title
		}
		return checked[Block](c, node, result, nil)
	case "divider":
		return checked[Block](c, node, DividerBlock{BlockID: id}, nil)
	case "markdown":
		return checked[Block](c, node, MarkdownBlock{Text: r.String("text")}, nil)
	case "actions":
		result := ActionsBlock{BlockID: id}
		for _, child := range node.Children("elements") {
			button, err := c.button(child)
			if err != nil {
				return nil, err
			}
			result.Elements = append(result.Elements, button)
		}
		return checked[Block](c, node, result, nil)
	case "context":
		result := ContextBlock{BlockID: id}
		for _, child := range node.Children("elements") {
			var element ContextElement
			var err error
			if child.Role == "image" {
				element, err = c.imageElement(child)
			} else {
				var text TextObject
				text, err = c.textObject(child)
				element, _ = text.(ContextElement)
			}
			if err != nil {
				return nil, err
			}
			result.Elements = append(result.Elements, element)
		}
		return checked[Block](c, node, result, nil)
	case "video":
		title, err := c.text(node, "title", true)
		if err != nil {
			return nil, err
		}
		result := VideoBlock{Title: title.(PlainTextObject), BlockID: id, VideoURL: r.String("video_url"), AltText: r.String("alt"), ThumbnailURL: r.String("thumbnail_url"), TitleURL: r.String("title_url"), AuthorName: r.String("author_name"), ProviderName: r.String("provider_name"), ProviderIconURL: r.String("provider_icon_url")}
		if node.Has("description") {
			text, err := c.text(node, "description", true)
			if err != nil {
				return nil, err
			}
			value := text.(PlainTextObject)
			result.Description = &value
		}
		return checked[Block](c, node, result, nil)
	case "rich_text":
		block, err := c.richText(node)
		return block, err
	}
	return nil, conversion.Error("slack", node.Path, "invalid_tag", "unsupported top-level block")
}

// textObject converts a nested plain_text or mrkdwn builder.
func (c converter) textObject(node conversion.Node) (TextObject, error) {
	r := conversion.Reader{Platform: "slack", Node: node}
	var result TextObject
	if node.Role == "plain_text" {
		value := PlainTextObject{Text: r.String("text")}
		if node.Has("emoji") {
			emoji := r.Bool("emoji")
			value.Emoji = &emoji
		}
		result = value
	} else {
		result = MrkdwnTextObject{Text: r.String("text"), Verbatim: r.Bool("verbatim")}
	}
	return checked(c, node, result, r.Err)
}

// text converts a slot holding either one text object child or scalar parts.
func (c converter) text(node conversion.Node, slot string, plain bool) (TextObject, error) {
	children := node.Children(slot)
	var text, format string
	count := 0
	for _, part := range node.Parts {
		if part.Slot != slot {
			continue
		}
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
	if len(children) != 0 {
		if count != 0 || len(children) != 1 {
			return nil, conversion.Error("slack", node.Path, "conflicting_input", "text slot requires one object or scalar parts")
		}
		if plain && children[0].Role != "plain_text" {
			return nil, conversion.Error("slack", children[0].Path, "invalid_tag", "plain text is required")
		}
		return c.textObject(children[0])
	}
	var result TextObject = PlainTextObject{Text: text}
	if format == "mrkdwn" {
		if plain {
			return nil, conversion.Error("slack", node.Path, "invalid_tag", "plain text is required")
		}
		result = MrkdwnTextObject{Text: text, Verbatim: true}
	}
	path := node.FieldPath(slot)
	return result, c.CheckAt(func(string) string { return path }, result.check)
}

func (c converter) slackFile(node conversion.Node) (*SlackFileObject, error) {
	children := node.Children("slack_file")
	if len(children) == 0 {
		return nil, nil
	}
	f := conversion.Reader{Platform: "slack", Node: children[0]}
	value, err := checked(c, children[0], SlackFileObject{URL: f.String("url"), ID: f.String("id")}, nil)
	return &value, err
}

func (c converter) imageElement(node conversion.Node) (ImageElement, error) {
	if node.Has("title") || node.Has("block_id") {
		return ImageElement{}, conversion.Error("slack", node.Path, "invalid_tag", "image elements do not accept title or block_id")
	}
	file, err := c.slackFile(node)
	if err != nil {
		return ImageElement{}, err
	}
	r := conversion.Reader{Platform: "slack", Node: node}
	return checked(c, node, ImageElement{ImageURL: r.String("url"), SlackFile: file, AltText: r.String("alt")}, nil)
}

func (c converter) button(node conversion.Node) (ButtonElement, error) {
	r := conversion.Reader{Platform: "slack", Node: node}
	text, err := c.text(node, "text", true)
	if err != nil {
		return ButtonElement{}, err
	}
	result := ButtonElement{Text: text.(PlainTextObject), ActionID: r.String("action_id"), URL: r.String("url"), Value: r.String("value"), Style: ButtonStyle(r.String("style")), AccessibilityLabel: r.String("accessibility_label"), AgentPrompt: r.String("agent_prompt"), AgentPromptDisplay: r.String("agent_prompt_display")}
	for _, part := range node.Parts {
		if part.Slot == "visible_to_user_ids" {
			result.VisibleToUserIDs = append(result.VisibleToUserIDs, part.Text)
		}
	}
	if children := node.Children("confirm"); len(children) != 0 {
		child := children[0]
		d := conversion.Reader{Platform: "slack", Node: child}
		dialog := ConfirmationDialogObject{Style: ButtonStyle(d.String("style"))}
		for _, slot := range []struct {
			name  string
			plain bool
		}{{"title", true}, {"text", false}, {"confirm", true}, {"deny", true}} {
			text, err := c.text(child, slot.name, slot.plain)
			if err != nil {
				return ButtonElement{}, err
			}
			switch slot.name {
			case "title":
				dialog.Title = text.(PlainTextObject)
			case "text":
				dialog.Text = text
			case "confirm":
				dialog.Confirm = text.(PlainTextObject)
			case "deny":
				dialog.Deny = text.(PlainTextObject)
			}
		}
		if dialog, err = checked(c, child, dialog, nil); err != nil {
			return ButtonElement{}, err
		}
		result.Confirm = &dialog
	}
	return checked(c, node, result, r.Err)
}
