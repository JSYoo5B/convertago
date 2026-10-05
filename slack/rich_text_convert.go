package slack

import "github.com/JSYoo5B/convertago/internal/conversion"

func (c converter) richText(node conversion.Node) (RichTextBlock, error) {
	r := conversion.Reader{Platform: "slack", Node: node}
	result := RichTextBlock{BlockID: r.String("block_id")}
	section := RichTextSection{}
	flush := func() error {
		if len(section.Elements) == 0 {
			return nil
		}
		checkedSection, err := checked(c, node, section, nil)
		result.Elements = append(result.Elements, checkedSection)
		section = RichTextSection{}
		return err
	}
	for _, input := range node.Inputs {
		if input.Slot == "text" {
			inline, err := c.richInline(input)
			if err != nil {
				return RichTextBlock{}, err
			}
			section.Elements = append(section.Elements, inline)
		} else if input.Slot == "elements" {
			if err := flush(); err != nil {
				return RichTextBlock{}, err
			}
			element, err := c.richElement(*input.Child)
			if err != nil {
				return RichTextBlock{}, err
			}
			result.Elements = append(result.Elements, element)
		}
	}
	if err := flush(); err != nil {
		return RichTextBlock{}, err
	}
	return checked(c, node, result, nil)
}

func richStyle(flags []string) *RichTextStyle {
	if len(flags) == 0 {
		return nil
	}
	result := &RichTextStyle{}
	for _, flag := range flags {
		switch flag {
		case "bold":
			result.Bold = true
		case "italic":
			result.Italic = true
		case "strike":
			result.Strike = true
		case "code":
			result.Code = true
		case "underline":
			result.Underline = true
		case "highlight":
			result.Highlight = true
		case "client_highlight":
			result.ClientHighlight = true
		case "unlink":
			result.Unlink = true
		}
	}
	return result
}

func (c converter) richInline(input conversion.Input) (RichTextInline, error) {
	if input.Part != nil {
		var result RichTextInline = TextInline{Text: input.Part.Text, Style: richStyle(input.Part.Style)}
		path := input.Part.Path
		return result, c.CheckAt(func(string) string { return path }, result.check)
	}
	node := *input.Child
	r := conversion.Reader{Platform: "slack", Node: node}
	var flags []string
	for _, part := range node.Parts {
		flags = append(flags, part.Style...)
	}
	style := richStyle(flags)
	var result RichTextInline
	switch node.Role {
	case "text":
		result = TextInline{Text: r.String("text"), Style: style}
	case "link":
		result = LinkInline{URL: r.String("url"), Text: r.String("text"), Unsafe: r.Bool("unsafe"), FromLLM: r.Bool("from_llm"), IsSlackURL: r.Bool("is_slack_url"), Truncated: r.Bool("truncated"), Style: style}
	case "user":
		result = UserInline{UserID: r.String("user_id"), FromLLM: r.Bool("from_llm"), Style: style}
	case "emoji":
		result = EmojiInline{Name: r.String("name"), Unicode: r.String("unicode")}
	default:
		return nil, conversion.Error("slack", node.Path, "invalid_tag", "unsupported rich text inline")
	}
	return checked(c, node, result, r.Err)
}

func (c converter) richElement(node conversion.Node) (RichTextElement, error) {
	r := conversion.Reader{Platform: "slack", Node: node}
	var border *int
	if node.Has("border") {
		value := r.ParseInt("border")
		border = &value
	}
	var inlines []RichTextInline
	for _, input := range node.Inputs {
		if input.Slot == "text" {
			inline, err := c.richInline(input)
			if err != nil {
				return nil, err
			}
			inlines = append(inlines, inline)
		}
	}
	var result RichTextElement
	switch node.Role {
	case "rich_text_section":
		result = RichTextSection{Elements: inlines}
	case "rich_text_quote":
		result = RichTextQuote{Elements: inlines, Border: border}
	case "rich_text_preformatted":
		value := RichTextPreformatted{Border: border, Language: r.String("language")}
		for _, inline := range inlines {
			value.Elements = append(value.Elements, inline.(PreformattedInline))
		}
		result = value
	case "rich_text_list":
		value := RichTextList{Style: RichTextListStyle(r.String("style")), Indent: r.ParseInt("indent"), Offset: r.ParseInt("offset"), Border: border}
		for _, child := range node.Children("elements") {
			section, err := c.richElement(child)
			if err != nil {
				return nil, err
			}
			value.Elements = append(value.Elements, section.(RichTextSection))
		}
		result = value
	default:
		return nil, conversion.Error("slack", node.Path, "invalid_tag", "unsupported rich text element")
	}
	return checked(c, node, result, r.Err)
}
