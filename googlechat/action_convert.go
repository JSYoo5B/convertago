package googlechat

import "github.com/JSYoo5B/convertago/internal/conversion"

func convertClick(node conversion.Node) (OnClick, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	result := OnClick{}
	switch node.Role {
	case "onClick":
		children := append(node.Children("action"), node.Children("openLink")...)
		children = append(children, node.Children("overflowMenu")...)
		if len(children) != 1 {
			r.Fail("conflicting_input", "onClick requires exactly one action")
			return result, r.Err
		}
		return convertClick(children[0])
	case "openLink":
		result.OpenLink = &OpenLink{URL: r.URI("url")}
	case "action":
		value, err := convertAction(node)
		if err != nil {
			return result, err
		}
		result.Action = &value
	case "overflowMenu":
		value := OverflowMenu{}
		for _, child := range r.Count("items", 1, 0) {
			itemReader := conversion.Reader{Platform: "googlechat", Node: child}
			click, err := convertClick(child.Children("onClick")[0])
			if err != nil {
				return result, err
			}
			if click.OverflowMenu != nil {
				return result, conversion.Error("googlechat", child.Path, "conflicting_input", "overflow menu items cannot open overflow menus")
			}
			item := OverflowMenuItem{Text: itemReader.Text("text", 1, 0), OnClick: click, Disabled: itemReader.Bool("disabled")}
			if children := child.Children("startIcon"); len(children) != 0 {
				icon, err := convertIcon(children[0])
				if err != nil {
					return result, err
				}
				item.StartIcon = &icon
			}
			if itemReader.Err != nil {
				return result, itemReader.Err
			}
			value.Items = append(value.Items, item)
		}
		result.OverflowMenu = &value
	default:
		r.Fail("invalid_tag", "unsupported click action")
	}
	return result, r.Err
}

func convertAction(node conversion.Node) (Action, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	result := Action{Function: r.Text("function", 1, 0), LoadIndicator: LoadIndicator(r.Enum("loadIndicator", "SPINNER", "NONE")), PersistValues: r.Bool("persistValues"), Interaction: Interaction(r.Enum("interaction", "OPEN_DIALOG")), AllWidgetsAreRequired: r.Bool("allWidgetsAreRequired")}
	for _, part := range node.Parts {
		if part.Slot == "requiredWidgets" {
			if err := conversion.ValidateText("googlechat", part.Path, part.Text, 1, 0); err != nil {
				return result, err
			}
			result.RequiredWidgets = append(result.RequiredWidgets, part.Text)
		}
	}
	if result.AllWidgetsAreRequired && len(result.RequiredWidgets) != 0 {
		r.Fail("conflicting_input", "requiredWidgets and allWidgetsAreRequired cannot both be supplied")
	}
	keys := map[string]bool{}
	for _, child := range node.Children("parameters") {
		p := conversion.Reader{Platform: "googlechat", Node: child}
		item := ActionParameter{Key: p.Text("key", 1, 0), Value: p.Text("value", 0, 0)}
		if p.Err != nil {
			return result, p.Err
		}
		if keys[item.Key] {
			return result, conversion.Error("googlechat", child.Path, "duplicate_input", "duplicate action parameter key")
		}
		keys[item.Key] = true
		result.Parameters = append(result.Parameters, item)
	}
	return result, r.Err
}

func convertColor(node conversion.Node) (Color, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	result := Color{Red: r.Float("red", 0, 1), Green: r.Float("green", 0, 1), Blue: r.Float("blue", 0, 1)}
	return result, r.Err
}

func convertIcon(node conversion.Node) (Icon, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	result := Icon{KnownIcon: r.Text("knownIcon", 1, 0), IconURL: r.URL("iconUrl", true), AltText: r.Text("altText", 0, 0), ImageType: ImageType(r.Enum("imageType", "SQUARE", "CIRCLE"))}
	count := 0
	if node.Has("knownIcon") {
		count++
	}
	if node.Has("iconUrl") {
		count++
	}
	if children := node.Children("materialIcon"); len(children) != 0 {
		count++
		child := children[0]
		m := conversion.Reader{Platform: "googlechat", Node: child}
		value := MaterialIcon{Name: m.Text("name", 1, 0), Fill: m.Bool("fill"), Weight: m.Int("weight", 0, 700), Grade: m.Int("grade", -25, 200)}
		if value.Grade != -25 && value.Grade != 0 && value.Grade != 200 {
			m.Fail("invalid_value", "material icon grade must be -25, 0, or 200")
		}
		if child.Has("weight") && value.Weight%100 != 0 {
			m.Fail("invalid_value", "material icon weight must be a multiple of 100")
		}
		if m.Err != nil {
			return result, m.Err
		}
		result.MaterialIcon = &value
	}
	if count != 1 {
		r.Fail("conflicting_input", "icon requires exactly one source")
	}
	return result, r.Err
}

func convertButton(node conversion.Node) (Button, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	click, err := convertClick(node.Children("onClick")[0])
	if err != nil {
		return Button{}, err
	}
	result := Button{Text: r.Text("text", 0, 0), OnClick: click, Disabled: r.Bool("disabled"), AltText: r.Text("altText", 0, 0), Type: ButtonType(r.Enum("type", "OUTLINED", "FILLED", "FILLED_TONAL", "BORDERLESS"))}
	if children := node.Children("icon"); len(children) != 0 {
		icon, err := convertIcon(children[0])
		if err != nil {
			return result, err
		}
		result.Icon = &icon
	}
	if children := node.Children("color"); len(children) != 0 {
		color, err := convertColor(children[0])
		if err != nil {
			return result, err
		}
		result.Color = &color
		if result.Type != "" && result.Type != ButtonTypeFilled {
			r.Fail("conflicting_input", "a color requires the filled button type")
		}
	}
	if result.Text == "" && result.Icon == nil {
		r.Fail("missing_input", "button requires text or an icon")
	}
	return result, r.Err
}
