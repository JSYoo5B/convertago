package slack

import "github.com/JSYoo5B/convertago/internal/conversion"

func init() {
	textSlot := map[string]conversion.Slot{"text": {Repeated: true, Required: true}}
	conversion.Register(conversion.Profile{Platform: "slack", Roles: map[string]conversion.Role{
		"header":    {DefaultSlot: "text", Slots: textSlot, Formats: []string{"plain"}},
		"section":   {DefaultSlot: "text", Slots: textSlot, Formats: []string{"plain", "mrkdwn"}},
		"rich_text": {DefaultSlot: "text", Slots: textSlot, Styles: []string{"bold", "italic", "strike", "code", "underline"}, Formats: []string{"plain"}},
		"image": {DefaultSlot: "url", Slots: map[string]conversion.Slot{
			"url": {Required: true}, "alt": {Required: true}, "title": {},
		}},
		"actions": {Unavailable: true}, "context": {Unavailable: true},
		"divider": {Unavailable: true}, "markdown": {Unavailable: true}, "video": {Unavailable: true},
	}})
}

// ToMessage builds a Message from slack tags in source declaration order.
// It supports header, section, rich_text, and image builders.
func ToMessage(input any, options ...conversion.Option) (Message, error) {
	nodes, err := conversion.Prepare(input, "slack", options)
	if err != nil {
		return Message{}, err
	}
	message := Message{}
	for _, node := range nodes {
		text := node.Text("text")
		switch node.Role {
		case "header":
			if err := conversion.ValidateText("slack", node.Path, text, 1, 150); err != nil {
				return Message{}, err
			}
			message.Blocks = append(message.Blocks, HeaderBlock{Text: PlainTextObject{Text: text}})
		case "section":
			if err := conversion.ValidateText("slack", node.Path, text, 1, 3000); err != nil {
				return Message{}, err
			}
			format := "plain"
			for i, part := range node.Parts {
				partFormat := part.Format
				if partFormat == "" {
					partFormat = "plain"
				}
				if i == 0 {
					format = partFormat
				} else if partFormat != format {
					return Message{}, conversion.Error("slack", part.Path, "conflicting_format", "section cannot mix plain text and mrkdwn")
				}
			}
			var object TextObject = PlainTextObject{Text: text}
			if format == "mrkdwn" {
				object = MrkdwnTextObject{Text: text, Verbatim: true}
			}
			message.Blocks = append(message.Blocks, SectionBlock{Text: object})
		case "rich_text":
			section := RichTextSection{}
			for _, part := range node.Parts {
				if err := conversion.ValidateText("slack", part.Path, part.Text, 0, 0); err != nil {
					return Message{}, err
				}
				inline := TextInline{Text: part.Text}
				if len(part.Style) != 0 {
					style := &RichTextStyle{}
					for _, flag := range part.Style {
						switch flag {
						case "bold":
							style.Bold = true
						case "italic":
							style.Italic = true
						case "strike":
							style.Strike = true
						case "code":
							style.Code = true
						case "underline":
							style.Underline = true
						}
					}
					inline.Style = style
				}
				section.Elements = append(section.Elements, inline)
			}
			message.Blocks = append(message.Blocks, RichTextBlock{Elements: []RichTextElement{section}})
		case "image":
			block := ImageBlock{ImageURL: node.Text("url"), AltText: node.Text("alt")}
			if err := conversion.ValidateURL("slack", node.Path, block.ImageURL, false); err != nil {
				return Message{}, err
			}
			if err := conversion.ValidateText("slack", node.Path, block.ImageURL, 1, 3000); err != nil {
				return Message{}, err
			}
			if err := conversion.ValidateText("slack", node.Path, block.AltText, 1, 2000); err != nil {
				return Message{}, err
			}
			if node.Has("title") {
				title := node.Text("title")
				if err := conversion.ValidateText("slack", node.Path, title, 1, 2000); err != nil {
					return Message{}, err
				}
				block.Title = &PlainTextObject{Text: title}
			}
			message.Blocks = append(message.Blocks, block)
		}
		if len(message.Blocks) > 50 {
			return Message{}, conversion.Error("slack", node.Path, "limit_exceeded", "message exceeds 50 blocks")
		}
	}
	return message, nil
}
