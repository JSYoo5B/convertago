package googlechat

import (
	"encoding/json"
	"html"
	"strings"

	"github.com/JSYoo5B/convertago/internal/conversion"
)

// ToMessage builds a Message from googlechat tags in source declaration order.
// Native elements are validated against their containing widget.
// Widgets are placed in one section of one card, with a header when supplied.
func ToMessage(input any, options ...conversion.Option) (Message, error) {
	nodes, err := conversion.Prepare(input, "googlechat", options)
	if err != nil {
		return Message{}, err
	}
	message := Message{}
	card := Card{}
	var widgets []Widget
	for _, node := range nodes {
		switch node.Role {
		case "text":
			message.Text += node.Text("text")
		case "fallbackText":
			message.FallbackText += node.Text("text")
		case "header":
			if card.Header != nil || len(widgets) != 0 {
				return Message{}, conversion.Error("googlechat", node.Path, "invalid_position", "header must precede widgets and occur only once")
			}
			header := &CardHeader{Title: node.Text("text"), Subtitle: node.Text("subtitle"), ImageURL: node.Text("url"), ImageAltText: node.Text("alt"), ImageType: ImageType(node.Text("imageType"))}
			if err := conversion.ValidateText("googlechat", node.Path, header.Title, 1, 0); err != nil {
				return Message{}, err
			}
			if node.Has("url") {
				if err := conversion.ValidateURL("googlechat", node.Path, header.ImageURL, true); err != nil {
					return Message{}, err
				}
			}
			card.Header = header
		default:
			widget, err := convertWidget(node)
			if err != nil {
				return Message{}, err
			}
			widgets = append(widgets, widget)
		}
		if len(widgets) > 100 {
			return Message{}, conversion.Error("googlechat", node.Path, "limit_exceeded", "card exceeds 100 widgets")
		}
	}
	if len(widgets) != 0 {
		card.Sections = []Section{{Widgets: widgets}}
	}
	if card.Header != nil || len(widgets) != 0 {
		message.CardsV2 = []CardWithID{{Card: card}}
		data, err := json.Marshal(message.CardsV2)
		if err != nil {
			return Message{}, err
		}
		if len(data) > 32*1024 {
			return Message{}, conversion.Error("googlechat", "$", "limit_exceeded", "cards exceed 32 KB of JSON")
		}
	}
	return message, nil
}

func convertParagraph(node conversion.Node) (TextParagraph, error) {
	markdown := false
	for _, part := range node.Parts {
		if part.Slot == "text" {
			markdown = part.Format == "markdown"
			break
		}
	}
	var text strings.Builder
	for _, part := range node.Parts {
		if part.Slot != "text" {
			continue
		}
		if (part.Format == "markdown") != markdown {
			return TextParagraph{}, conversion.Error("googlechat", part.Path, "conflicting_format", "paragraph cannot mix Markdown with HTML or plain text")
		}
		if err := conversion.ValidateText("googlechat", part.Path, part.Text, 0, 0); err != nil {
			return TextParagraph{}, err
		}
		content := part.Text
		if !markdown && part.Format != "html" {
			content = strings.ReplaceAll(html.EscapeString(content), "\n", "<br>")
		}
		// Use a fixed nesting order so style flag order does not change the JSON.
		for _, style := range []string{"underline", "code", "strike", "italic", "bold"} {
			if !hasStyle(part.Style, style) {
				continue
			}
			if markdown {
				marker := map[string]string{"bold": "**", "italic": "*", "strike": "~", "code": "`"}[style]
				content = marker + content + marker
			} else {
				tag := map[string]string{"bold": "b", "italic": "i", "strike": "s", "code": "code", "underline": "u"}[style]
				content = "<" + tag + ">" + content + "</" + tag + ">"
			}
		}
		text.WriteString(content)
	}
	reader := conversion.Reader{Platform: "googlechat", Node: node}
	paragraph := TextParagraph{Text: text.String(), MaxLines: reader.Int("maxLines", 0, 0)}
	if markdown {
		paragraph.TextSyntax = TextSyntaxMarkdown
	}
	return paragraph, reader.Err
}

func hasStyle(styles []string, target string) bool {
	for _, style := range styles {
		if style == target {
			return true
		}
	}
	return false
}
