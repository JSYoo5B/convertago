package googlechat

import (
	"html"
	"strings"

	"github.com/JSYoo5B/convertago/internal/conversion"
)

func (c converter) header(node conversion.Node) (CardHeader, error) {
	r := conversion.Reader{Platform: "googlechat", Node: node}
	result := CardHeader{Title: r.String("text"), Subtitle: r.String("subtitle"), ImageURL: r.String("url"), ImageAltText: r.String("alt"), ImageType: ImageType(r.String("imageType"))}
	return checked(c, node, result, nil, "title", "text", "imageUrl", "url", "imageAltText", "alt")
}

func (c converter) card(node conversion.Node) (Card, error) {
	var nodes []conversion.Node
	for _, input := range node.Inputs {
		if input.Child != nil {
			nodes = append(nodes, *input.Child)
		}
	}
	card, err := c.cardNodes(nodes)
	if err != nil {
		return Card{}, err
	}
	r := conversion.Reader{Platform: "googlechat", Node: node}
	card.SectionDividerStyle = DividerStyle(r.String("sectionDividerStyle"))
	return checked(c, node, card, nil)
}

// cardNodes assembles a header, explicit sections, and implicit sections of widgets.
func (c converter) cardNodes(nodes []conversion.Node) (Card, error) {
	card := Card{}
	var widgets []Widget
	var widgetPath string
	flush := func() error {
		if len(widgets) == 0 {
			return nil
		}
		section, err := checkedAt(c, widgetPath, Section{Widgets: widgets})
		card.Sections = append(card.Sections, section)
		widgets = nil
		return err
	}
	for _, node := range nodes {
		switch node.Role {
		case "header":
			if card.Header != nil || len(widgets) != 0 || len(card.Sections) != 0 {
				return card, conversion.Error("googlechat", node.Path, "invalid_position", "header must precede widgets and occur only once")
			}
			header, err := c.header(node)
			if err != nil {
				return card, err
			}
			card.Header = &header
		case "section":
			if err := flush(); err != nil {
				return card, err
			}
			section, err := c.section(node)
			if err != nil {
				return card, err
			}
			card.Sections = append(card.Sections, section)
		default:
			widget, err := c.widget(node)
			if err != nil {
				return card, err
			}
			if len(widgets) == 0 {
				widgetPath = node.Path
			}
			widgets = append(widgets, widget)
		}
	}
	return card, flush()
}

func (c converter) section(node conversion.Node) (Section, error) {
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
	result := Section{Header: header.String(), Collapsible: r.Bool("collapsible"), UncollapsibleWidgetsCount: r.ParseInt("uncollapsibleWidgetsCount")}
	for _, child := range node.Children("widgets") {
		widget, err := c.widget(child)
		if err != nil {
			return result, err
		}
		result.Widgets = append(result.Widgets, widget)
	}
	if children := node.Children("collapseControl"); len(children) != 0 {
		child := children[0]
		expand, err := c.button(child.Children("expandButton")[0])
		if err != nil {
			return result, err
		}
		collapse, err := c.button(child.Children("collapseButton")[0])
		if err != nil {
			return result, err
		}
		cr := conversion.Reader{Platform: "googlechat", Node: child}
		value, err := checked(c, child, CollapseControl{HorizontalAlignment: HorizontalAlignment(cr.String("horizontalAlignment")), ExpandButton: expand, CollapseButton: collapse}, nil)
		if err != nil {
			return result, err
		}
		result.CollapseControl = &value
	}
	return checked(c, node, result, r.Err)
}
