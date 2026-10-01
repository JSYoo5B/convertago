package googlechat

import (
	"encoding/json"
	"html"
	"strings"

	"github.com/JSYoo5B/convertago/internal/conversion"
)

// ToMessage builds a Message from googlechat tags in source declaration order.
// Native elements are validated against their containing widget.
// Top-level widgets form an implicit card; card and section builders select explicit layouts.
func ToMessage(input any, options ...conversion.Option) (Message, error) {
	nodes, err := conversion.Prepare(input, "googlechat", options)
	if err != nil {
		return Message{}, err
	}
	message := Message{}
	var implicit []conversion.Node
	flush := func() error {
		if len(implicit) == 0 {
			return nil
		}
		card, err := convertCardNodes(implicit, "$")
		if err != nil {
			return err
		}
		message.CardsV2 = append(message.CardsV2, CardWithID{Card: card})
		implicit = nil
		return nil
	}
	for _, node := range nodes {
		switch node.Role {
		case "text":
			message.Text += node.Text("text")
		case "fallbackText":
			message.FallbackText += node.Text("text")
		case "card", "cardWithId":
			if err := flush(); err != nil {
				return Message{}, err
			}
			wrapped := CardWithID{}
			cardNode := node
			if node.Role == "cardWithId" {
				wrapped.CardID = node.Text("cardId")
				cardNode = node.Children("card")[0]
			}
			card, err := convertCard(cardNode)
			if err != nil {
				return Message{}, err
			}
			wrapped.Card = card
			message.CardsV2 = append(message.CardsV2, wrapped)
		default:
			implicit = append(implicit, node)
		}
	}
	if err := flush(); err != nil {
		return Message{}, err
	}
	ids := map[string]bool{}
	for _, card := range message.CardsV2 {
		if len(message.CardsV2) > 1 && card.CardID == "" {
			return Message{}, conversion.Error("googlechat", "$", "missing_input", "multiple cards require distinct cardId values")
		}
		if card.CardID != "" {
			if ids[card.CardID] {
				return Message{}, conversion.Error("googlechat", "$", "duplicate_input", "duplicate cardId")
			}
			ids[card.CardID] = true
		}
	}
	if len(message.CardsV2) != 0 {
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
