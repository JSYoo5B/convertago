package slack

import "github.com/JSYoo5B/convertago/internal/conversion"

func convertRichText(node conversion.Node) (RichTextBlock, error) {
	r := conversion.Reader{Platform: "slack", Node: node}
	result := RichTextBlock{BlockID: r.Text("block_id", 1, 255)}
	section := RichTextSection{}
	flush := func() {
		if len(section.Elements) != 0 {
			result.Elements = append(result.Elements, section)
			section = RichTextSection{}
		}
	}
	for _, input := range node.Inputs {
		if input.Slot == "text" {
			inline, err := convertRichInline(input)
			if err != nil {
				return RichTextBlock{}, err
			}
			section.Elements = append(section.Elements, inline)
		} else if input.Slot == "elements" {
			flush()
			element, err := convertRichElement(*input.Child)
			if err != nil {
				return RichTextBlock{}, err
			}
			result.Elements = append(result.Elements, element)
		}
	}
	flush()
	if len(result.Elements) == 0 {
		r.Fail("missing_input", "rich_text requires content")
	}
	return result, r.Err
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

func convertRichInline(input conversion.Input) (RichTextInline, error) {
	if input.Part != nil {
		return TextInline{Text: input.Part.Text, Style: richStyle(input.Part.Style)}, nil
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
		result = TextInline{Text: r.Text("text", 0, 0), Style: style}
	case "link":
		result = LinkInline{URL: r.URI("url"), Text: r.Text("text", 0, 0), Unsafe: r.Bool("unsafe"), FromLLM: r.Bool("from_llm"), IsSlackURL: r.Bool("is_slack_url"), Truncated: r.Bool("truncated"), Style: style}
	case "user":
		result = UserInline{UserID: r.Text("user_id", 1, 0), FromLLM: r.Bool("from_llm"), Style: style}
	case "emoji":
		result = EmojiInline{Name: r.Text("name", 1, 0), Unicode: r.Text("unicode", 0, 0)}
	default:
		r.Fail("invalid_tag", "unsupported rich text inline")
	}
	return result, r.Err
}

func convertRichElement(node conversion.Node) (RichTextElement, error) {
	r := conversion.Reader{Platform: "slack", Node: node}
	var border *int
	if node.Has("border") {
		value := r.Int("border", 0, 1)
		border = &value
	}
	var inlines []RichTextInline
	for _, input := range node.Inputs {
		if input.Slot == "text" {
			inline, err := convertRichInline(input)
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
		value := RichTextPreformatted{Border: border, Language: r.Text("language", 1, 0)}
		for _, inline := range inlines {
			value.Elements = append(value.Elements, inline.(PreformattedInline))
		}
		result = value
	case "rich_text_list":
		value := RichTextList{Style: RichTextListStyle(r.Enum("style", "bullet", "ordered")), Indent: r.Int("indent", 0, 0), Offset: r.Int("offset", 0, 0), Border: border}
		if value.Style != "ordered" && node.Has("offset") {
			r.Fail("conflicting_input", "offset requires an ordered list")
		}
		for _, child := range r.Count("elements", 1, 0) {
			section, err := convertRichElement(child)
			if err != nil {
				return nil, err
			}
			value.Elements = append(value.Elements, section.(RichTextSection))
		}
		result = value
	default:
		r.Fail("invalid_tag", "unsupported rich text element")
	}
	return result, r.Err
}
