package googlechat

import "github.com/JSYoo5B/convertago/internal/conversion"

func (c converter) widget(node conversion.Node) (Widget, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	content, err := c.content(node)
	if err != nil {
		return Widget{}, err
	}
	return checked(c, node, Widget{Content: content, HorizontalAlignment: HorizontalAlignment(r.String("horizontalAlignment"))}, nil)
}

func (c converter) content(node conversion.Node) (WidgetContent, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	switch node.Role {
	case "textParagraph":
		return c.paragraph(node)
	case "divider":
		return Divider{}, nil
	case "image":
		value := Image{ImageURL: r.String("url"), AltText: r.String("alt")}
		if children := node.Children("onClick"); len(children) != 0 {
			click, err := c.click(children[0])
			if err != nil {
				return nil, err
			}
			value.OnClick = &click
		}
		return checked(c, node, value, nil, "imageUrl", "url", "altText", "alt")
	case "buttonList":
		value := ButtonList{}
		for _, child := range node.Children("buttons") {
			button, err := c.button(child)
			if err != nil {
				return nil, err
			}
			value.Buttons = append(value.Buttons, button)
		}
		return checked(c, node, value, nil)
	case "decoratedText":
		return c.decorated(node)
	case "columns", "grid", "carousel", "chipList":
		return c.layout(node)
	}
	return nil, conversion.Error("googlechat", node.Path, "invalid_tag", "unsupported widget")
}

func (c converter) decorated(node conversion.Node) (DecoratedText, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	text, _, err := paragraphText(node)
	if err != nil {
		return DecoratedText{}, err
	}
	result := DecoratedText{Text: text, TopLabel: r.String("topLabel"), BottomLabel: r.String("bottomLabel"), WrapText: r.Bool("wrapText"), StartIconVerticalAlignment: VerticalAlignment(r.String("startIconVerticalAlignment"))}
	for _, slot := range []struct {
		name   string
		target **Icon
	}{{"startIcon", &result.StartIcon}, {"endIcon", &result.EndIcon}} {
		if children := node.Children(slot.name); len(children) != 0 {
			icon, err := c.icon(children[0])
			if err != nil {
				return result, err
			}
			*slot.target = &icon
		}
	}
	for _, slot := range []struct {
		name   string
		target **TextParagraph
	}{{"topLabelText", &result.TopLabelText}, {"contentText", &result.ContentText}, {"bottomLabelText", &result.BottomLabelText}} {
		if children := node.Children(slot.name); len(children) != 0 {
			text, err := c.paragraph(children[0])
			if err != nil {
				return result, err
			}
			*slot.target = &text
		}
	}
	if children := node.Children("onClick"); len(children) != 0 {
		click, err := c.click(children[0])
		if err != nil {
			return result, err
		}
		result.OnClick = &click
	}
	if children := node.Children("button"); len(children) != 0 {
		button, err := c.button(children[0])
		if err != nil {
			return result, err
		}
		result.Button = &button
	}
	if children := node.Children("switchControl"); len(children) != 0 {
		child := children[0]
		s := conversion.Reader{Platform: "googlechat", Node: child}
		value := SwitchControl{Name: s.String("name"), Value: s.String("value"), Selected: s.Bool("selected"), ControlType: SwitchControlType(s.String("controlType"))}
		if actions := child.Children("onChangeAction"); len(actions) != 0 {
			action, err := c.action(actions[0])
			if err != nil {
				return result, err
			}
			value.OnChangeAction = &action
		}
		value, err := checked(c, child, value, s.Err)
		if err != nil {
			return result, err
		}
		result.SwitchControl = &value
	}
	return checked(c, node, result, r.Err)
}
