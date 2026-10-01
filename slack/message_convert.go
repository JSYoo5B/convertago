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
			block, err := convertRichText(node)
			if err != nil {
				return Message{}, err
			}
			message.Blocks = append(message.Blocks, block)
		case "text":
			for _, part := range node.Parts {
				if len(part.Style) != 0 {
					return Message{}, conversion.Error("slack", part.Path, "invalid_tag", "message text cannot carry rich text styles")
				}
			}
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
