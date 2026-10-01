package slack

import "github.com/JSYoo5B/convertago/internal/conversion"

// ToMessage builds a Message from slack tags in source declaration order.
// Nested elements are validated against their containing block.
func ToMessage(input any, options ...conversion.Option) (Message, error) {
	nodes, err := conversion.Prepare(input, "slack", options)
	if err != nil {
		return Message{}, err
	}
	message := Message{}
	for _, node := range nodes {
		text := node.Text("text")
		switch node.Role {
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
		case "text":
			message.Text += text
		default:
			block, err := convertBlock(node)
			if err != nil {
				return Message{}, err
			}
			message.Blocks = append(message.Blocks, block)
		}
		if len(message.Blocks) > 50 {
			return Message{}, conversion.Error("slack", node.Path, "limit_exceeded", "message exceeds 50 blocks")
		}
	}
	var markdown string
	for _, block := range message.Blocks {
		if value, ok := block.(MarkdownBlock); ok {
			markdown += value.Text
		}
	}
	if err := conversion.ValidateText("slack", "$", markdown, 0, 12000); err != nil {
		return Message{}, err
	}
	return message, nil
}
