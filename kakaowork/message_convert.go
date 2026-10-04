package kakaowork

import (
	"fmt"

	"github.com/JSYoo5B/convertago/internal/conversion"
	"github.com/JSYoo5B/convertago/internal/validation"
)

// ToMessage 는 kakaowork 태그를 해석하여 Message 를 구성합니다.
// 블록과 중첩 요소는 필드 선언 순서로 구성되며, 각 요소는 Validate 와 같은 규칙으로 검사합니다.
func ToMessage(input any, options ...conversion.Option) (Message, error) {
	nodes, err := conversion.Prepare(input, "kakaowork", options)
	if err != nil {
		return Message{}, err
	}
	c := converter{conversion.Checker{Platform: "kakaowork", Options: conversion.Configure(options)}}
	message := Message{Blocks: make([]BubbleBlock, 0, len(nodes))}
	previewPath := "$"
	var blockPaths []string
	for _, node := range nodes {
		if node.Role == "preview" {
			if message.Preview == "" {
				previewPath = node.Path
			}
			message.Preview += node.Text("text")
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
			return previewPath
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

// checked runs the rules of a value built from node.
func checked[T interface{ check(*validation.Check) }](c converter, node conversion.Node, value T, err error) (T, error) {
	if err != nil {
		return value, err
	}
	return value, c.Check(node, value.check)
}
