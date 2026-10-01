package googlechat

import (
	"encoding/json"
	"html"
	"strings"

	"github.com/JSYoo5B/convertago/internal/conversion"
)

func init() {
	conversion.Register(conversion.Profile{Platform: "googlechat", Roles: map[string]conversion.Role{
		"header": {DefaultSlot: "text", Slots: map[string]conversion.Slot{
			"text": {Repeated: true, Required: true}, "subtitle": {Repeated: true}, "url": {}, "alt": {},
		}, Formats: []string{"plain"}},
		"textParagraph": {DefaultSlot: "text", Slots: map[string]conversion.Slot{
			"text": {Repeated: true, Required: true},
		}, Styles: []string{"bold", "italic", "strike", "code", "underline"},
			Formats:      []string{"plain", "html", "markdown"},
			FormatStyles: map[string][]string{"markdown": {"bold", "italic", "strike", "code"}},
		},
		"image":         {DefaultSlot: "url", Slots: map[string]conversion.Slot{"url": {Required: true}, "alt": {}}},
		"fallbackText":  {DefaultSlot: "text", Slots: map[string]conversion.Slot{"text": {Repeated: true, Required: true}}, Formats: []string{"plain"}},
		"decoratedText": {Unavailable: true}, "buttonList": {Unavailable: true},
		"divider": {Unavailable: true}, "columns": {Unavailable: true}, "grid": {Unavailable: true},
		"carousel": {Unavailable: true}, "chipList": {Unavailable: true},
	}})
}

// ToMessage builds a Message from googlechat tags in source declaration order.
// It supports header, textParagraph, image, and fallbackText builders.
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
		case "fallbackText":
			message.FallbackText += node.Text("text")
		case "header":
			if card.Header != nil || len(widgets) != 0 {
				return Message{}, conversion.Error("googlechat", node.Path, "invalid_position", "header must precede widgets and occur only once")
			}
			header := &CardHeader{Title: node.Text("text"), Subtitle: node.Text("subtitle"), ImageURL: node.Text("url"), ImageAltText: node.Text("alt")}
			if err := conversion.ValidateText("googlechat", node.Path, header.Title, 1, 0); err != nil {
				return Message{}, err
			}
			if node.Has("url") {
				if err := conversion.ValidateURL("googlechat", node.Path, header.ImageURL, true); err != nil {
					return Message{}, err
				}
			}
			card.Header = header
		case "textParagraph":
			paragraph, err := convertParagraph(node)
			if err != nil {
				return Message{}, err
			}
			widgets = append(widgets, Widget{Content: paragraph})
		case "image":
			image := Image{ImageURL: node.Text("url"), AltText: node.Text("alt")}
			if err := conversion.ValidateURL("googlechat", node.Path, image.ImageURL, true); err != nil {
				return Message{}, err
			}
			widgets = append(widgets, Widget{Content: image})
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
	markdown := node.Parts[0].Format == "markdown"
	var text strings.Builder
	for _, part := range node.Parts {
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
	paragraph := TextParagraph{Text: text.String()}
	if markdown {
		paragraph.TextSyntax = TextSyntaxMarkdown
	}
	return paragraph, nil
}

func hasStyle(styles []string, target string) bool {
	for _, style := range styles {
		if style == target {
			return true
		}
	}
	return false
}
