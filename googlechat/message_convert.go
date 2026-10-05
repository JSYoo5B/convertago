package googlechat

import (
	"fmt"
	"html"
	"strings"

	"github.com/JSYoo5B/convertago/internal/conversion"
	"github.com/JSYoo5B/convertago/internal/validation"
)

// ToMessage builds a Message from googlechat tags in source declaration order.
// Top-level widgets form an implicit card; card and section builders select explicit layouts.
// Each native value is checked with the same rules as Validate.
func ToMessage(input any, options ...conversion.Option) (Message, error) {
	nodes, err := conversion.Prepare(input, "googlechat", options)
	if err != nil {
		return Message{}, err
	}
	c := converter{conversion.Checker{Platform: "googlechat", Options: conversion.Configure(options)}}
	message := Message{}
	var cardPaths []string
	var implicit []conversion.Node
	flush := func() error {
		if len(implicit) == 0 {
			return nil
		}
		card, err := c.cardNodes(implicit)
		if err == nil {
			card, err = checkedAt(c, implicit[0].Path, card)
		}
		if err != nil {
			return err
		}
		message.CardsV2 = append(message.CardsV2, CardWithID{Card: card})
		cardPaths = append(cardPaths, implicit[0].Path)
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
			card, err := c.card(cardNode)
			if err != nil {
				return Message{}, err
			}
			wrapped.Card = card
			message.CardsV2 = append(message.CardsV2, wrapped)
			cardPaths = append(cardPaths, node.Path)
		default:
			implicit = append(implicit, node)
		}
	}
	if err := flush(); err != nil {
		return Message{}, err
	}
	resolve := func(field string) string {
		var index int
		if _, err := fmt.Sscanf(field, "cardsV2[%d]", &index); err == nil && index < len(cardPaths) {
			return cardPaths[index]
		}
		return "$"
	}
	if err := c.CheckAt(resolve, message.check); err != nil {
		return Message{}, err
	}
	return message, nil
}

type converter struct {
	conversion.Checker
}

type checkable interface{ check(*validation.Check) }

// checked runs the rules of a value built from node, mapping native fields to slots.
func checked[T checkable](c converter, node conversion.Node, value T, err error, slots ...string) (T, error) {
	if err != nil {
		return value, err
	}
	resolve := func(field string) string {
		for i := 0; i+1 < len(slots); i += 2 {
			if field == slots[i] {
				field = slots[i+1]
			}
		}
		return node.FieldPath(field)
	}
	return value, c.CheckAt(resolve, value.check)
}

// checkedAt runs the rules of a value assembled from several nodes at one source path.
func checkedAt[T checkable](c converter, path string, value T) (T, error) {
	return value, c.CheckAt(func(string) string { return path }, value.check)
}

// paragraphText renders the text slot of node as HTML or Markdown.
func paragraphText(node conversion.Node) (string, bool, error) {
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
			return "", false, conversion.Error("googlechat", part.Path, "conflicting_format", "paragraph cannot mix Markdown with HTML or plain text")
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
	return text.String(), markdown, nil
}

func (c converter) paragraph(node conversion.Node) (TextParagraph, error) {
	text, markdown, err := paragraphText(node)
	if err != nil {
		return TextParagraph{}, err
	}
	r := conversion.Reader{Platform: "googlechat", Node: node}
	paragraph := TextParagraph{Text: text, MaxLines: r.ParseInt("maxLines")}
	if markdown {
		paragraph.TextSyntax = TextSyntaxMarkdown
	}
	return checked(c, node, paragraph, r.Err)
}

func hasStyle(styles []string, target string) bool {
	for _, style := range styles {
		if style == target {
			return true
		}
	}
	return false
}
