package kakaowork

import (
	"strings"

	"github.com/JSYoo5B/convertago/internal/conversion"
)

func (c converter) block(node conversion.Node) (BubbleBlock, error) {
	r := conversion.Reader{Platform: "kakaowork", Node: node}
	switch node.Role {
	case "header":
		return checked[BubbleBlock](c, node, HeaderBlock{Text: r.String("text"), Style: HeaderStyle(r.String("style"))}, nil)
	case "text":
		block, err := c.text(node)
		return block, err
	case "image_link":
		return checked[BubbleBlock](c, node, ImageBlock{Url: r.String("url")}, nil)
	case "divider":
		return DividerBlock{}, nil
	case "button":
		block, err := c.button(node)
		return block, err
	case "action":
		result := ActionBlock{}
		for _, child := range node.Children("elements") {
			button, err := c.button(child)
			if err != nil {
				return nil, err
			}
			result.Elements = append(result.Elements, button)
		}
		return checked[BubbleBlock](c, node, result, nil)
	case "description", "section", "context":
		content, err := c.text(node.Children("content")[0])
		if err != nil {
			return nil, err
		}
		switch node.Role {
		case "description":
			block := DescriptionBlock{Content: content, Term: r.String("term"), Accent: r.Bool("accent")}
			return checked[BubbleBlock](c, node, block, r.Err)
		case "section":
			result := SectionBlock{Content: content}
			if children := node.Children("accessory"); len(children) != 0 {
				image, err := checked(c, children[0], ImageBlock{Url: children[0].Text("url")}, nil)
				if err != nil {
					return nil, err
				}
				result.Accessory = &image
			}
			if children := node.Children("action"); len(children) != 0 {
				action, err := c.action(children[0])
				if err != nil {
					return nil, err
				}
				result.Action = action
			}
			return checked[BubbleBlock](c, node, result, nil)
		default:
			image := node.Children("image")[0]
			block := ContextBlock{Content: content, Image: ImageBlock{Url: image.Text("url")}}
			if _, err := checked(c, image, block.Image, nil); err != nil {
				return nil, err
			}
			return checked[BubbleBlock](c, node, block, nil)
		}
	}
	return nil, conversion.Error("kakaowork", node.Path, "invalid_tag", "unsupported top-level block")
}

func (c converter) button(node conversion.Node) (ButtonBlock, error) {
	r := conversion.Reader{Platform: "kakaowork", Node: node}
	action, err := c.action(node.Children("action")[0])
	if err != nil {
		return ButtonBlock{}, err
	}
	return checked(c, node, ButtonBlock{Text: r.String("text"), Style: ButtonStyle(r.String("style")), Action: action}, nil)
}

func (c converter) text(node conversion.Node) (TextBlock, error) {
	result := TextBlock{}
	styled := false
	var sources []conversion.Node
	for _, input := range node.Inputs {
		var inline Inline
		source := conversion.Node{Path: node.Path}
		if input.Part != nil {
			item := InlineStyled{Text: input.Part.Text}
			applyStyles(&item, input.Part.Style)
			styled = styled || len(input.Part.Style) != 0
			inline = item
			source.Path = input.Part.Path
		} else {
			child := *input.Child
			r := conversion.Reader{Platform: "kakaowork", Node: child}
			text := r.String("text")
			switch child.Role {
			case "styled":
				item := InlineStyled{Text: text, Color: InlineColor(r.String("color")), Bold: r.Bool("bold"), Italic: r.Bool("italic"), Strike: r.Bool("strike")}
				for _, part := range child.Parts {
					if part.Slot == "text" {
						applyStyles(&item, part.Style)
					}
				}
				inline = item
			case "link":
				inline = InlineLink{Text: text, Url: r.String("url")}
			case "mention":
				inline = InlineMention{Text: text, UserId: r.ParseInt("user_id")}
			}
			if r.Err != nil {
				return TextBlock{}, r.Err
			}
			styled = true
			source = child
		}
		result.Text += inline.String()
		result.Inlines = append(result.Inlines, inline)
		sources = append(sources, source)
	}
	if !styled {
		result.Inlines = nil
	}
	for i, inline := range result.Inlines {
		if err := c.Check(sources[i], inline.check); err != nil {
			return TextBlock{}, err
		}
	}
	resolve := func(field string) string {
		if strings.HasPrefix(field, "inlines[") {
			field = "text" + strings.TrimPrefix(field, "inlines")
		}
		return node.FieldPath(field)
	}
	return result, c.CheckAt(resolve, result.check)
}

func applyStyles(item *InlineStyled, styles []string) {
	for _, flag := range styles {
		switch flag {
		case "bold":
			item.Bold = true
		case "italic":
			item.Italic = true
		case "strike":
			item.Strike = true
		}
	}
}

func (c converter) action(node conversion.Node) (ButtonAction, error) {
	r := conversion.Reader{Platform: "kakaowork", Node: node}
	name, value := r.String("name"), r.String("value")
	var result ButtonAction
	switch node.Role {
	case "open_system_browser":
		result = OpenSystemBrowserAction{Name: name, Value: value}
	case "open_inapp_browser":
		result = OpenInAppBrowserAction{Name: name, Value: value, Standalone: r.Bool("standalone"), Width: r.ParseInt("width"), Height: r.ParseInt("height")}
	case "open_external_app":
		result = OpenExternalAppAction{Name: name, Value: value}
	case "submit_action":
		result = SubmitAction{Name: name, Value: value}
	case "call_modal":
		result = CallModalAction{Name: name, Value: value}
	case "exclusive":
		action := ExclusiveAction{}
		for _, slot := range []struct {
			name   string
			target *ButtonAction
		}{{"default", &action.Default}, {"pc", &action.Pc}, {"mobile", &action.Mobile}, {"windows", &action.Windows}, {"macos", &action.MacOs}, {"android", &action.Android}, {"ios", &action.Ios}} {
			if children := node.Children(slot.name); len(children) != 0 {
				value, err := c.action(children[0])
				if err != nil {
					return nil, err
				}
				*slot.target = value
			}
		}
		result = action
	default:
		return nil, conversion.Error("kakaowork", node.Path, "invalid_tag", "unsupported button action")
	}
	if r.Err != nil {
		return nil, r.Err
	}
	return result, c.Check(node, result.check)
}
