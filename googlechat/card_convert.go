package googlechat

import (
	"github.com/JSYoo5B/convertago/internal/conversion"
	"html"
	"strings"
)

func convertHeader(node conversion.Node) (CardHeader, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	result := CardHeader{Title: r.Text("text", 1, 0), Subtitle: r.Text("subtitle", 0, 0), ImageURL: r.URL("url", true), ImageAltText: r.Text("alt", 0, 0), ImageType: ImageType(r.Enum("imageType", "SQUARE", "CIRCLE"))}
	if !node.Has("url") && (node.Has("imageType") || node.Has("alt")) {
		r.Fail("conflicting_input", "image properties require a header image URL")
	}
	return result, r.Err
}

func convertCard(node conversion.Node) (Card, error) {
	var nodes []conversion.Node
	for _, input := range node.Inputs {
		if input.Child != nil {
			nodes = append(nodes, *input.Child)
		}
	}
	card, err := convertCardNodes(nodes, node.Path)
	if err != nil {
		return Card{}, err
	}
	r := conversion.Reader{Platform: "googlechat", Node: node}
	card.SectionDividerStyle = DividerStyle(r.Enum("sectionDividerStyle", "SOLID_DIVIDER", "NO_DIVIDER"))
	return card, r.Err
}

func convertCardNodes(nodes []conversion.Node, path string) (Card, error) {
	card := Card{}
	var widgets []Widget
	flush := func() {
		if len(widgets) != 0 {
			card.Sections = append(card.Sections, Section{Widgets: widgets})
			widgets = nil
		}
	}
	for _, node := range nodes {
		switch node.Role {
		case "header":
			if card.Header != nil || len(widgets) != 0 || len(card.Sections) != 0 {
				return card, conversion.Error("googlechat", node.Path, "invalid_position", "header must precede widgets and occur only once")
			}
			header, err := convertHeader(node)
			if err != nil {
				return card, err
			}
			card.Header = &header
		case "section":
			flush()
			section, err := convertSection(node)
			if err != nil {
				return card, err
			}
			card.Sections = append(card.Sections, section)
		default:
			widget, err := convertWidget(node)
			if err != nil {
				return card, err
			}
			widgets = append(widgets, widget)
		}
	}
	flush()
	if card.Header == nil && len(card.Sections) == 0 {
		return card, conversion.Error("googlechat", path, "missing_input", "card requires a header or widgets")
	}
	total := 0
	for _, section := range card.Sections {
		for _, widget := range section.Widgets {
			total += widgetCount(widget.Content)
		}
	}
	if total > 100 {
		return card, conversion.Error("googlechat", path, "limit_exceeded", "card exceeds 100 widgets")
	}
	return card, nil
}

func widgetCount(content WidgetContent) int {
	count := 1
	switch value := content.(type) {
	case Columns:
		for _, column := range value.ColumnItems {
			for _, widget := range column.Widgets {
				count += widgetCount(widget.Content)
			}
		}
	case Carousel:
		for _, card := range value.CarouselCards {
			for _, widget := range card.Widgets {
				count += widgetCount(widget.Content)
			}
			for _, widget := range card.FooterWidgets {
				count += widgetCount(widget.Content)
			}
		}
	}
	return count
}

func convertSection(node conversion.Node) (Section, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	var header strings.Builder
	for _, part := range node.Parts {
		if part.Slot == "header" {
			if part.Format == "html" {
				header.WriteString(part.Text)
			} else {
				header.WriteString(html.EscapeString(part.Text))
			}
		}
	}
	result := Section{Header: header.String(), Collapsible: r.Bool("collapsible")}
	for _, child := range r.Count("widgets", 1, 100) {
		widget, err := convertWidget(child)
		if err != nil {
			return result, err
		}
		result.Widgets = append(result.Widgets, widget)
	}
	result.UncollapsibleWidgetsCount = r.Int("uncollapsibleWidgetsCount", 0, len(result.Widgets))
	if children := node.Children("collapseControl"); len(children) != 0 {
		child := children[0]
		c := conversion.Reader{Platform: "googlechat", Node: child}
		expand, err := convertButton(child.Children("expandButton")[0])
		if err != nil {
			return result, err
		}
		collapse, err := convertButton(child.Children("collapseButton")[0])
		if err != nil {
			return result, err
		}
		value := CollapseControl{HorizontalAlignment: HorizontalAlignment(c.Enum("horizontalAlignment", "START", "CENTER", "END")), ExpandButton: expand, CollapseButton: collapse}
		if c.Err != nil {
			return result, c.Err
		}
		result.CollapseControl = &value
	}
	if !result.Collapsible && (result.UncollapsibleWidgetsCount != 0 || result.CollapseControl != nil) {
		r.Fail("conflicting_input", "collapse properties require a collapsible section")
	}
	return result, r.Err
}
