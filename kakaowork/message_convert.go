package kakaowork

import (
	"strings"

	"github.com/JSYoo5B/convertago/internal/conversion"
)

func init() {
	conversion.Register(conversion.Profile{Platform: "kakaowork", Roles: map[string]conversion.Role{
		"header": {DefaultSlot: "text", Slots: map[string]conversion.Slot{
			"text": {Repeated: true, Required: true}, "style": {},
		}},
		"text": {DefaultSlot: "text", Slots: map[string]conversion.Slot{
			"text": {Repeated: true, Required: true},
		}, Styles: []string{"bold", "italic", "strike"}, Formats: []string{"plain"}},
		"image_link": {DefaultSlot: "url", Slots: map[string]conversion.Slot{"url": {Required: true}}},
		"preview":    {DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": {Repeated: true, Required: true}}},
		"button":     {Unavailable: true}, "action": {Unavailable: true},
		"divider": {Unavailable: true}, "description": {Unavailable: true},
		"section": {Unavailable: true}, "context": {Unavailable: true},
	}})
}

// ToMessage 는 kakaowork 태그를 해석하여 Message 를 구성합니다.
// header, text, image_link, preview 를 지원하며, 블록은 필드 선언 순서로 구성됩니다.
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
		case "text":
			if err := conversion.ValidateText("kakaowork", node.Path, text, 0, 500); err != nil {
				return Message{}, err
			}
			block := TextBlock{Text: text}
			styled := false
			for _, part := range node.Parts {
				styled = styled || len(part.Style) != 0
			}
			if styled {
				for _, part := range node.Parts {
					inline := InlineStyled{Text: part.Text}
					for _, style := range part.Style {
						switch style {
						case "bold":
							inline.Bold = true
						case "italic":
							inline.Italic = true
						case "strike":
							inline.Strike = true
						}
					}
					block.Inlines = append(block.Inlines, inline)
				}
			}
			message.Blocks = append(message.Blocks, block)
		case "image_link":
			url := node.Text("url")
			if err := conversion.ValidateURL("kakaowork", node.Path, url, false); err != nil {
				return Message{}, err
			}
			message.Blocks = append(message.Blocks, ImageBlock{Url: url})
		}
	}
	return message, nil
}
