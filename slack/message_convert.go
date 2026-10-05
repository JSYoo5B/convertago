package slack

import (
	"fmt"

	"github.com/JSYoo5B/convertago/internal/conversion"
)

// ToMessage builds a Message from slack tags in source declaration order.
// Each native value is checked with the same rules as Validate.
func ToMessage(input any, options ...conversion.Option) (Message, error) {
	nodes, err := conversion.Prepare(input, "slack", options)
	if err != nil {
		return Message{}, err
	}
	c := converter{conversion.Checker{Platform: "slack", Options: conversion.Configure(options), Slots: map[string]string{"image_url": "url", "alt_text": "alt"}}}
	message := Message{}
	textPath := "$"
	var blockPaths []string
	for _, node := range nodes {
		if node.Role == "text" {
			for _, part := range node.Parts {
				if len(part.Style) != 0 {
					return Message{}, conversion.Error("slack", part.Path, "invalid_tag", "message text cannot carry rich text styles")
				}
			}
			if message.Text == "" {
				textPath = node.Path
			}
			message.Text += node.Text("text")
			continue
		}
		block, err := c.block(node)
		if err != nil {
			return Message{}, err
		}
		message.Blocks = append(message.Blocks, block)
		blockPaths = append(blockPaths, node.Path)
	}
	resolve := func(field string) string {
		var index int
		if _, err := fmt.Sscanf(field, "blocks[%d]", &index); err == nil && index < len(blockPaths) {
			return blockPaths[index]
		}
		if field == "text" {
			return textPath
		}
		return "$"
	}
	if err := c.CheckAt(resolve, message.check); err != nil {
		return Message{}, err
	}
	return message, nil
}
