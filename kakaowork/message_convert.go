package kakaowork

import (
	"strings"

	"github.com/JSYoo5B/convertago/internal/conversion"
)

// ToMessage 는 kakaowork 태그를 해석하여 Message 를 구성합니다.
// 블록과 중첩 요소는 필드 선언 순서로 구성됩니다.
func ToMessage(input any, options ...conversion.Option) (Message, error) {
	nodes, err := conversion.Prepare(input, "kakaowork", options)
	if err != nil {
		return Message{}, err
	}
	message := Message{Blocks: make([]BubbleBlock, 0, len(nodes))}
	for _, node := range nodes {
		text := node.Text("text")
		switch node.Role {
		case "preview":
			message.Preview += text
		case "header":
			if len(message.Blocks) != 0 {
				return Message{}, conversion.Error("kakaowork", node.Path, "invalid_position", "header must be the first and only header block")
			}
			if strings.ContainsAny(text, "\r\n") {
				return Message{}, conversion.Error("kakaowork", node.Path, "invalid_value", "header does not support line breaks")
			}
			if err := conversion.ValidateText("kakaowork", node.Path, text, 0, 20); err != nil {
				return Message{}, err
			}
			style := HeaderStyleWhite
			if node.Has("style") {
				style = HeaderStyle(node.Text("style"))
				if _, exists := headerStyleConstants[style]; !exists || style == "" {
					return Message{}, conversion.Error("kakaowork", node.Path, "invalid_value", "unknown header background style")
				}
			}
			message.Blocks = append(message.Blocks, HeaderBlock{Text: text, Style: style})
		default:
			block, err := convertBlock(node)
			if err != nil {
				return Message{}, err
			}
			message.Blocks = append(message.Blocks, block)
		}
	}
	return message, nil
}
