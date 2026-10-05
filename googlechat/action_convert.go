package googlechat

import "github.com/JSYoo5B/convertago/internal/conversion"

// click converts an explicit onClick wrapper or a direct action, openLink, or overflowMenu child.
func (c converter) click(node conversion.Node) (OnClick, error) {
	result := OnClick{}
	children := []conversion.Node{node}
	if node.Role == "onClick" {
		children = append(append(node.Children("action"), node.Children("openLink")...), node.Children("overflowMenu")...)
	}
	for _, child := range children {
		switch child.Role {
		case "openLink":
			r := conversion.Reader{Platform: "googlechat", Node: child}
			link, err := checked(c, child, OpenLink{URL: r.String("url")}, nil)
			if err != nil {
				return result, err
			}
			result.OpenLink = &link
		case "action":
			action, err := c.action(child)
			if err != nil {
				return result, err
			}
			result.Action = &action
		case "overflowMenu":
			menu, err := c.overflowMenu(child)
			if err != nil {
				return result, err
			}
			result.OverflowMenu = &menu
		default:
			return result, conversion.Error("googlechat", child.Path, "invalid_tag", "unsupported click action")
		}
	}
	return checked(c, node, result, nil)
}

func (c converter) overflowMenu(node conversion.Node) (OverflowMenu, error) {
	result := OverflowMenu{}
	for _, child := range node.Children("items") {
		r := conversion.Reader{Platform: "googlechat", Node: child}
		click, err := c.click(child.Children("onClick")[0])
		if err != nil {
			return result, err
		}
		item := OverflowMenuItem{Text: r.String("text"), OnClick: click, Disabled: r.Bool("disabled")}
		if children := child.Children("startIcon"); len(children) != 0 {
			icon, err := c.icon(children[0])
			if err != nil {
				return result, err
			}
			item.StartIcon = &icon
		}
		if item, err = checked(c, child, item, r.Err); err != nil {
			return result, err
		}
		result.Items = append(result.Items, item)
	}
	return checked(c, node, result, nil)
}

func (c converter) action(node conversion.Node) (Action, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	result := Action{Function: r.String("function"), LoadIndicator: LoadIndicator(r.String("loadIndicator")), PersistValues: r.Bool("persistValues"), Interaction: Interaction(r.String("interaction")), AllWidgetsAreRequired: r.Bool("allWidgetsAreRequired")}
	for _, part := range node.Parts {
		if part.Slot == "requiredWidgets" {
			result.RequiredWidgets = append(result.RequiredWidgets, part.Text)
		}
	}
	for _, child := range node.Children("parameters") {
		p := conversion.Reader{Platform: "googlechat", Node: child}
		result.Parameters = append(result.Parameters, ActionParameter{Key: p.String("key"), Value: p.String("value")})
	}
	return checked(c, node, result, r.Err)
}

func (c converter) color(node conversion.Node) (Color, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	return checked(c, node, Color{Red: r.ParseFloat("red"), Green: r.ParseFloat("green"), Blue: r.ParseFloat("blue")}, r.Err)
}

func (c converter) icon(node conversion.Node) (Icon, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	result := Icon{KnownIcon: r.String("knownIcon"), IconURL: r.String("iconUrl"), AltText: r.String("altText"), ImageType: ImageType(r.String("imageType"))}
	if children := node.Children("materialIcon"); len(children) != 0 {
		child := children[0]
		m := conversion.Reader{Platform: "googlechat", Node: child}
		value, err := checked(c, child, MaterialIcon{Name: m.String("name"), Fill: m.Bool("fill"), Weight: m.ParseInt("weight"), Grade: m.ParseInt("grade")}, m.Err)
		if err != nil {
			return result, err
		}
		result.MaterialIcon = &value
	}
	return checked(c, node, result, r.Err)
}

func (c converter) button(node conversion.Node) (Button, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	click, err := c.click(node.Children("onClick")[0])
	if err != nil {
		return Button{}, err
	}
	result := Button{Text: r.String("text"), OnClick: click, Disabled: r.Bool("disabled"), AltText: r.String("altText"), Type: ButtonType(r.String("type"))}
	if children := node.Children("icon"); len(children) != 0 {
		icon, err := c.icon(children[0])
		if err != nil {
			return result, err
		}
		result.Icon = &icon
	}
	if children := node.Children("color"); len(children) != 0 {
		color, err := c.color(children[0])
		if err != nil {
			return result, err
		}
		result.Color = &color
	}
	return checked(c, node, result, r.Err)
}
