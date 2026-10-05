package googlechat

import "github.com/JSYoo5B/convertago/internal/conversion"

func (c converter) layout(node conversion.Node) (WidgetContent, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	switch node.Role {
	case "columns":
		value := Columns{}
		for _, child := range node.Children("columnItems") {
			cr := conversion.Reader{Platform: "googlechat", Node: child}
			column := Column{HorizontalSizeStyle: HorizontalSizeStyle(cr.String("horizontalSizeStyle")), HorizontalAlignment: HorizontalAlignment(cr.String("horizontalAlignment")), VerticalAlignment: ColumnVerticalAlignment(cr.String("verticalAlignment"))}
			for _, widget := range child.Children("widgets") {
				content, err := c.nestedContent(widget)
				if err != nil {
					return nil, err
				}
				column.Widgets = append(column.Widgets, ColumnWidget{Content: content.(ColumnWidgetContent)})
			}
			column, err := checked(c, child, column, nil)
			if err != nil {
				return nil, err
			}
			value.ColumnItems = append(value.ColumnItems, column)
		}
		return checked(c, node, value, nil)
	case "grid":
		value := Grid{Title: r.String("title"), ColumnCount: r.ParseInt("columnCount")}
		for _, child := range node.Children("items") {
			item, err := c.gridItem(child)
			if err != nil {
				return nil, err
			}
			value.Items = append(value.Items, item)
		}
		if children := node.Children("borderStyle"); len(children) != 0 {
			border, err := c.border(children[0])
			if err != nil {
				return nil, err
			}
			value.BorderStyle = &border
		}
		if children := node.Children("onClick"); len(children) != 0 {
			click, err := c.click(children[0])
			if err != nil {
				return nil, err
			}
			value.OnClick = &click
		}
		return checked(c, node, value, r.Err)
	case "carousel":
		value := Carousel{}
		for _, child := range node.Children("carouselCards") {
			card := CarouselCard{}
			for _, slot := range []struct {
				name   string
				target *[]NestedWidget
			}{{"widgets", &card.Widgets}, {"footerWidgets", &card.FooterWidgets}} {
				for _, widget := range child.Children(slot.name) {
					content, err := c.nestedContent(widget)
					if err != nil {
						return nil, err
					}
					*slot.target = append(*slot.target, NestedWidget{Content: content.(NestedWidgetContent)})
				}
			}
			card, err := checked(c, child, card, nil)
			if err != nil {
				return nil, err
			}
			value.CarouselCards = append(value.CarouselCards, card)
		}
		return checked(c, node, value, nil)
	case "chipList":
		value := ChipList{Layout: ChipListLayout(r.String("layout"))}
		for _, child := range node.Children("chips") {
			cr := conversion.Reader{Platform: "googlechat", Node: child}
			chip := Chip{Label: cr.String("label"), Disabled: cr.Bool("disabled"), AltText: cr.String("altText")}
			if children := child.Children("icon"); len(children) != 0 {
				icon, err := c.icon(children[0])
				if err != nil {
					return nil, err
				}
				chip.Icon = &icon
			}
			if children := child.Children("onClick"); len(children) != 0 {
				click, err := c.click(children[0])
				if err != nil {
					return nil, err
				}
				chip.OnClick = &click
			}
			chip, err := checked(c, child, chip, cr.Err)
			if err != nil {
				return nil, err
			}
			value.Chips = append(value.Chips, chip)
		}
		return checked(c, node, value, nil)
	}
	return nil, conversion.Error("googlechat", node.Path, "invalid_tag", "unsupported layout")
}

func (c converter) nestedContent(node conversion.Node) (WidgetContent, error) {
	if node.Has("horizontalAlignment") {
		return nil, conversion.Error("googlechat", node.Path, "invalid_tag", "nested widgets do not accept horizontalAlignment")
	}
	return c.content(node)
}

func (c converter) gridItem(node conversion.Node) (GridItem, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	result := GridItem{ID: r.String("id"), Title: r.String("title"), Subtitle: r.String("subtitle"), Layout: GridItemLayout(r.String("layout"))}
	if children := node.Children("image"); len(children) != 0 {
		child := children[0]
		i := conversion.Reader{Platform: "googlechat", Node: child}
		value := ImageComponent{ImageURI: i.String("imageUri"), AltText: i.String("altText")}
		if crops := child.Children("cropStyle"); len(crops) != 0 {
			cr := conversion.Reader{Platform: "googlechat", Node: crops[0]}
			style, err := checked(c, crops[0], ImageCropStyle{Type: ImageCropType(cr.String("type")), AspectRatio: cr.ParseFloat("aspectRatio")}, cr.Err)
			if err != nil {
				return result, err
			}
			value.CropStyle = &style
		}
		if borders := child.Children("borderStyle"); len(borders) != 0 {
			border, err := c.border(borders[0])
			if err != nil {
				return result, err
			}
			value.BorderStyle = &border
		}
		value, err := checked(c, child, value, i.Err)
		if err != nil {
			return result, err
		}
		result.Image = &value
	}
	return checked(c, node, result, r.Err)
}

func (c converter) border(node conversion.Node) (BorderStyle, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	result := BorderStyle{Type: BorderType(r.String("type")), CornerRadius: r.ParseInt("cornerRadius")}
	if children := node.Children("strokeColor"); len(children) != 0 {
		color, err := c.color(children[0])
		if err != nil {
			return result, err
		}
		result.StrokeColor = &color
	}
	return checked(c, node, result, r.Err)
}
