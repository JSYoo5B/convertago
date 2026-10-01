package googlechat

import "github.com/JSYoo5B/convertago/internal/conversion"

func convertLayout(node conversion.Node) (WidgetContent, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	var result WidgetContent
	switch node.Role {
	case "columns":
		value := Columns{}
		for _, child := range r.Count("columnItems", 1, 2) {
			c := conversion.Reader{Platform: "googlechat", Node: child}
			column := Column{HorizontalSizeStyle: HorizontalSizeStyle(c.Enum("horizontalSizeStyle", "FILL_AVAILABLE_SPACE", "FILL_MINIMUM_SPACE")), HorizontalAlignment: HorizontalAlignment(c.Enum("horizontalAlignment", "START", "CENTER", "END")), VerticalAlignment: ColumnVerticalAlignment(c.Enum("verticalAlignment", "CENTER", "TOP", "BOTTOM"))}
			for _, widget := range c.Count("widgets", 1, 100) {
				content, err := nestedContent(widget)
				if err != nil {
					return nil, err
				}
				column.Widgets = append(column.Widgets, ColumnWidget{Content: content.(ColumnWidgetContent)})
			}
			if c.Err != nil {
				return nil, c.Err
			}
			value.ColumnItems = append(value.ColumnItems, column)
		}
		result = value
	case "grid":
		value := Grid{Title: r.Text("title", 0, 0), ColumnCount: r.Int("columnCount", 1, 0)}
		for _, child := range r.Count("items", 1, 0) {
			item, err := convertGridItem(child)
			if err != nil {
				return nil, err
			}
			value.Items = append(value.Items, item)
		}
		if children := node.Children("borderStyle"); len(children) != 0 {
			border, err := convertBorder(children[0])
			if err != nil {
				return nil, err
			}
			value.BorderStyle = &border
		}
		if children := node.Children("onClick"); len(children) != 0 {
			click, err := convertClick(children[0])
			if err != nil {
				return nil, err
			}
			value.OnClick = &click
		}
		result = value
	case "carousel":
		value := Carousel{}
		for _, child := range r.Count("carouselCards", 1, 0) {
			c := conversion.Reader{Platform: "googlechat", Node: child}
			card := CarouselCard{}
			for _, slot := range []struct {
				name   string
				target *[]NestedWidget
				min    int
			}{{"widgets", &card.Widgets, 1}, {"footerWidgets", &card.FooterWidgets, 0}} {
				for _, widget := range c.Count(slot.name, slot.min, 100) {
					content, err := nestedContent(widget)
					if err != nil {
						return nil, err
					}
					*slot.target = append(*slot.target, NestedWidget{Content: content.(NestedWidgetContent)})
				}
			}
			if c.Err != nil {
				return nil, c.Err
			}
			value.CarouselCards = append(value.CarouselCards, card)
		}
		result = value
	case "chipList":
		value := ChipList{Layout: ChipListLayout(r.Enum("layout", "WRAPPED", "HORIZONTAL_SCROLLABLE"))}
		for _, child := range r.Count("chips", 1, 0) {
			c := conversion.Reader{Platform: "googlechat", Node: child}
			chip := Chip{Label: c.Text("label", 0, 0), Disabled: c.Bool("disabled"), AltText: c.Text("altText", 0, 0)}
			if children := child.Children("icon"); len(children) != 0 {
				icon, err := convertIcon(children[0])
				if err != nil {
					return nil, err
				}
				chip.Icon = &icon
			}
			if children := child.Children("onClick"); len(children) != 0 {
				click, err := convertClick(children[0])
				if err != nil {
					return nil, err
				}
				chip.OnClick = &click
			}
			if chip.Label == "" && chip.Icon == nil {
				c.Fail("missing_input", "chip requires a label or icon")
			}
			if c.Err != nil {
				return nil, c.Err
			}
			value.Chips = append(value.Chips, chip)
		}
		result = value
	default:
		r.Fail("invalid_tag", "unsupported layout")
	}
	return result, r.Err
}

func nestedContent(node conversion.Node) (WidgetContent, error) {
	if node.Has("horizontalAlignment") {
		return nil, conversion.Error("googlechat", node.Path, "invalid_tag", "nested widgets do not accept horizontalAlignment")
	}
	return convertContent(node)
}

func convertGridItem(node conversion.Node) (GridItem, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	result := GridItem{ID: r.Text("id", 0, 0), Title: r.Text("title", 0, 0), Subtitle: r.Text("subtitle", 0, 0), Layout: GridItemLayout(r.Enum("layout", "TEXT_BELOW", "TEXT_ABOVE"))}
	if children := node.Children("image"); len(children) != 0 {
		child := children[0]
		i := conversion.Reader{Platform: "googlechat", Node: child}
		value := ImageComponent{ImageURI: i.URL("imageUri", false), AltText: i.Text("altText", 0, 0)}
		if crops := child.Children("cropStyle"); len(crops) != 0 {
			crop := crops[0]
			c := conversion.Reader{Platform: "googlechat", Node: crop}
			style := ImageCropStyle{Type: ImageCropType(c.Enum("type", "SQUARE", "CIRCLE", "RECTANGLE_CUSTOM", "RECTANGLE_4_3")), AspectRatio: c.Float("aspectRatio", 0, 0)}
			if style.Type == ImageCropTypeRectangleCustom {
				if !crop.Has("aspectRatio") || style.AspectRatio <= 0 {
					c.Fail("missing_input", "custom rectangle requires a positive aspectRatio")
				}
			} else if crop.Has("aspectRatio") {
				c.Fail("conflicting_input", "aspectRatio requires a custom rectangle")
			}
			if c.Err != nil {
				return result, c.Err
			}
			value.CropStyle = &style
		}
		if borders := child.Children("borderStyle"); len(borders) != 0 {
			border, err := convertBorder(borders[0])
			if err != nil {
				return result, err
			}
			value.BorderStyle = &border
		}
		if i.Err != nil {
			return result, i.Err
		}
		result.Image = &value
	}
	if result.Title == "" && result.Subtitle == "" && result.Image == nil {
		r.Fail("missing_input", "grid item requires text or an image")
	}
	return result, r.Err
}

func convertBorder(node conversion.Node) (BorderStyle, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	result := BorderStyle{Type: BorderType(r.Enum("type", "NO_BORDER", "STROKE")), CornerRadius: r.Int("cornerRadius", 0, 0)}
	if children := node.Children("strokeColor"); len(children) != 0 {
		color, err := convertColor(children[0])
		if err != nil {
			return result, err
		}
		result.StrokeColor = &color
	}
	if result.Type == BorderTypeNone && result.StrokeColor != nil {
		r.Fail("conflicting_input", "strokeColor requires a stroke border")
	}
	return result, r.Err
}
