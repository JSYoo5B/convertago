package googlechat

import "github.com/JSYoo5B/convertago/internal/conversion"

func convertWidget(node conversion.Node) (Widget, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	widget := Widget{HorizontalAlignment: HorizontalAlignment(r.Enum("horizontalAlignment", "START", "CENTER", "END"))}
	content, err := convertContent(node)
	if err != nil {
		return widget, err
	}
	widget.Content = content
	return widget, r.Err
}

func convertContent(node conversion.Node) (WidgetContent, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	var result WidgetContent
	switch node.Role {
	case "textParagraph":
		return convertParagraph(node)
	case "divider":
		result = Divider{}
	case "image":
		value := Image{ImageURL: r.URL("url", true), AltText: r.Text("alt", 0, 0)}
		if children := node.Children("onClick"); len(children) != 0 {
			click, err := convertClick(children[0])
			if err != nil {
				return nil, err
			}
			value.OnClick = &click
		}
		result = value
	case "buttonList":
		value := ButtonList{}
		for _, child := range r.Count("buttons", 1, 0) {
			button, err := convertButton(child)
			if err != nil {
				return nil, err
			}
			value.Buttons = append(value.Buttons, button)
		}
		result = value
	case "decoratedText":
		return convertDecorated(node)
	case "columns", "grid", "carousel", "chipList":
		return convertLayout(node)
	default:
		r.Fail("invalid_tag", "unsupported widget")
	}
	return result, r.Err
}

func convertDecorated(node conversion.Node) (DecoratedText, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	paragraph, err := convertParagraph(node)
	if err != nil {
		return DecoratedText{}, err
	}
	result := DecoratedText{Text: paragraph.Text, TopLabel: r.Text("topLabel", 0, 0), BottomLabel: r.Text("bottomLabel", 0, 0), WrapText: r.Bool("wrapText"), StartIconVerticalAlignment: VerticalAlignment(r.Enum("startIconVerticalAlignment", "TOP", "MIDDLE", "BOTTOM"))}
	for _, slot := range []struct {
		name   string
		target **Icon
	}{{"startIcon", &result.StartIcon}, {"endIcon", &result.EndIcon}} {
		if children := node.Children(slot.name); len(children) != 0 {
			icon, err := convertIcon(children[0])
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
			text, err := convertParagraph(children[0])
			if err != nil {
				return result, err
			}
			*slot.target = &text
		}
	}
	if children := node.Children("onClick"); len(children) != 0 {
		click, err := convertClick(children[0])
		if err != nil {
			return result, err
		}
		result.OnClick = &click
	}
	count := 0
	if result.EndIcon != nil {
		count++
	}
	if children := node.Children("button"); len(children) != 0 {
		count++
		button, err := convertButton(children[0])
		if err != nil {
			return result, err
		}
		result.Button = &button
	}
	if children := node.Children("switchControl"); len(children) != 0 {
		count++
		child := children[0]
		s := conversion.Reader{Platform: "googlechat", Node: child}
		value := SwitchControl{Name: s.Text("name", 1, 0), Value: s.Text("value", 0, 0), Selected: s.Bool("selected"), ControlType: SwitchControlType(s.Enum("controlType", "SWITCH", "CHECK_BOX"))}
		if actions := child.Children("onChangeAction"); len(actions) != 0 {
			action, err := convertAction(actions[0])
			if err != nil {
				return result, err
			}
			value.OnChangeAction = &action
		}
		if s.Err != nil {
			return result, s.Err
		}
		result.SwitchControl = &value
	}
	if count > 1 {
		r.Fail("conflicting_input", "decoratedText accepts only one trailing control")
	}
	return result, r.Err
}
