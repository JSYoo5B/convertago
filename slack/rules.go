package slack

import (
	"fmt"

	"github.com/JSYoo5B/convertago/internal/validation"
)

const (
	referenceBlockKit = "https://docs.slack.dev/reference/block-kit/"
	referenceBlocks   = referenceBlockKit + "blocks/"
	referenceElements = referenceBlockKit + "block-elements/"
	referenceObjects  = referenceBlockKit + "composition-objects/"
	referencePost     = "https://docs.slack.dev/reference/methods/chat.postMessage/"
)

const (
	fatal    = validation.Fatal
	warning  = validation.Warning
	advisory = validation.Advisory
)

func rule(id string, severity validation.Severity, verified bool, reference, message string) validation.Rule {
	return validation.Rule{ID: id, Severity: severity, Verified: verified, Reference: reference, Message: message}
}

// Block Kit constraints are verified Fatal rules: chat.postMessage rejects blocks
// that do not match Block Kit with invalid_blocks.
var (
	ruleMessageContentRequired = rule("message.content.required", fatal, true, referencePost, "a message requires text or blocks")
	ruleMessageBlocksCount     = rule("message.blocks.count", fatal, true, referenceBlocks, "a message holds at most 50 blocks")
	ruleMessageTextFallback    = rule("message.text.fallback", advisory, true, referencePost, "text is recommended as the fallback for blocks")
	ruleMessageTextLength      = rule("message.text.length", advisory, true, referencePost, "text over 40000 characters is truncated")
	ruleMarkdownTextTotal      = rule("markdown.text.total", warning, false, referenceBlocks+"markdown-block/", "markdown blocks together exceed 12000 characters")
	ruleMarkdownTextRequired   = rule("markdown.text.required", warning, false, referenceBlocks+"markdown-block/", "markdown text is required")
	ruleElementRequired        = rule("element.required", fatal, true, referenceBlocks, "a list contains a missing element")
	ruleBlockIDLength          = rule("block_id.length", fatal, true, referenceBlocks, "block_id exceeds 255 characters")

	ruleTextRequired = rule("text.text.required", fatal, true, referenceObjects+"text-object/", "text object text is required")
	ruleTextLength   = rule("text.text.length", fatal, true, referenceObjects+"text-object/", "text object text exceeds 3000 characters")

	ruleHeaderTextLength  = rule("header.text.length", fatal, true, referenceBlocks+"header-block/", "header text exceeds 150 characters")
	ruleHeaderLevelValue  = rule("header.level.value", fatal, true, referenceBlocks+"header-block/", "header level must be 1 to 4")
	ruleSectionContent    = rule("section.content.required", fatal, true, referenceBlocks+"section-block/", "a section requires text or fields")
	ruleSectionFieldCount = rule("section.fields.count", fatal, true, referenceBlocks+"section-block/", "a section holds at most 10 fields")
	ruleSectionFieldLen   = rule("section.fields.length", fatal, true, referenceBlocks+"section-block/", "a section field exceeds 2000 characters")

	ruleImageSource      = rule("image.source.required", fatal, true, referenceBlocks+"image-block/", "an image requires exactly one of image_url or slack_file")
	ruleImageURLLength   = rule("image.image_url.length", fatal, true, referenceBlocks+"image-block/", "image_url exceeds 3000 characters")
	ruleImageURLFormat   = rule("image.image_url.format", warning, false, referenceBlocks+"image-block/", "image_url must be an absolute HTTP or HTTPS URL")
	ruleImageAltRequired = rule("image.alt_text.required", warning, false, referenceBlocks+"image-block/", "alt_text is required")
	ruleImageAltLength   = rule("image.alt_text.length", fatal, true, referenceBlocks+"image-block/", "alt_text exceeds 2000 characters")
	ruleImageTitleLength = rule("image.title.length", fatal, true, referenceBlocks+"image-block/", "image title exceeds 2000 characters")
	ruleSlackFileSource  = rule("slack_file.source.required", fatal, true, referenceObjects+"slack-file-object/", "a slack_file requires exactly one of url or id")
	ruleSlackFileURL     = rule("slack_file.url.format", warning, false, referenceObjects+"slack-file-object/", "slack_file url must be an absolute HTTP or HTTPS URL")

	ruleActionsCount    = rule("actions.elements.count", fatal, true, referenceBlocks+"actions-block/", "an actions block holds at most 25 elements")
	ruleActionsRequired = rule("actions.elements.required", warning, false, referenceBlocks+"actions-block/", "an actions block requires elements")
	ruleContextCount    = rule("context.elements.count", fatal, true, referenceBlocks+"context-block/", "a context block holds at most 10 elements")
	ruleContextRequired = rule("context.elements.required", warning, false, referenceBlocks+"context-block/", "a context block requires elements")

	ruleButtonTextLength      = rule("button.text.length", fatal, true, referenceElements+"button-element/", "button text exceeds 75 characters")
	ruleButtonTextTruncation  = rule("button.text.truncation", advisory, true, referenceElements+"button-element/", "button text may truncate after about 30 characters")
	ruleButtonActionIDLength  = rule("button.action_id.length", fatal, true, referenceElements+"button-element/", "action_id exceeds 255 characters")
	ruleButtonURLLength       = rule("button.url.length", fatal, true, referenceElements+"button-element/", "button url exceeds 3000 characters")
	ruleButtonURLFormat       = rule("button.url.format", warning, false, referenceElements+"button-element/", "button url must be an absolute URI")
	ruleButtonValueLength     = rule("button.value.length", fatal, true, referenceElements+"button-element/", "button value exceeds 2000 characters")
	ruleButtonStyleValue      = rule("button.style.value", fatal, true, referenceElements+"button-element/", "button style must be primary or danger")
	ruleButtonLabelLength     = rule("button.accessibility_label.length", fatal, true, referenceElements+"button-element/", "accessibility_label exceeds 75 characters")
	ruleButtonPromptLength    = rule("button.agent_prompt.length", fatal, true, referenceElements+"button-element/", "agent_prompt exceeds 4000 characters")
	ruleButtonPromptDisplay   = rule("button.agent_prompt_display.prompt", warning, true, referenceElements+"button-element/", "agent_prompt_display takes effect only with agent_prompt")
	ruleButtonVisibleRequired = rule("button.visible_to_user_ids.required", warning, false, referenceElements+"button-element/", "visible_to_user_ids contains an empty user ID")

	ruleConfirmTitleLength   = rule("confirm.title.length", fatal, true, referenceObjects+"confirmation-dialog-object/", "confirmation title exceeds 100 characters")
	ruleConfirmTextRequired  = rule("confirm.text.required", fatal, true, referenceObjects+"confirmation-dialog-object/", "confirmation text is required")
	ruleConfirmTextLength    = rule("confirm.text.length", fatal, true, referenceObjects+"confirmation-dialog-object/", "confirmation text exceeds 300 characters")
	ruleConfirmConfirmLength = rule("confirm.confirm.length", fatal, true, referenceObjects+"confirmation-dialog-object/", "confirm label exceeds 30 characters")
	ruleConfirmDenyLength    = rule("confirm.deny.length", fatal, true, referenceObjects+"confirmation-dialog-object/", "deny label exceeds 30 characters")
	ruleConfirmStyleValue    = rule("confirm.style.value", fatal, true, referenceObjects+"confirmation-dialog-object/", "confirmation style must be primary or danger")

	ruleVideoTitleLength       = rule("video.title.length", fatal, true, referenceBlocks+"video-block/", "video title must be shorter than 200 characters")
	ruleVideoDescriptionLength = rule("video.description.length", fatal, true, referenceBlocks+"video-block/", "video description must be shorter than 200 characters")
	ruleVideoAuthorLength      = rule("video.author_name.length", fatal, true, referenceBlocks+"video-block/", "author_name must be shorter than 50 characters")
	ruleVideoURLHTTPS          = rule("video.video_url.https", fatal, true, referenceBlocks+"video-block/", "video_url must be an HTTPS URL")
	ruleVideoTitleURLHTTPS     = rule("video.title_url.https", fatal, true, referenceBlocks+"video-block/", "title_url must be an HTTPS URL")
	ruleVideoAltRequired       = rule("video.alt_text.required", warning, false, referenceBlocks+"video-block/", "video alt_text is required")
	ruleVideoThumbnailFormat   = rule("video.thumbnail_url.format", warning, false, referenceBlocks+"video-block/", "thumbnail_url must be an absolute HTTP or HTTPS URL")
	ruleVideoProviderIcon      = rule("video.provider_icon_url.format", warning, false, referenceBlocks+"video-block/", "provider_icon_url must be an absolute HTTP or HTTPS URL")

	ruleRichTextRequired    = rule("rich_text.elements.required", warning, false, referenceBlocks+"rich-text-block/", "rich text requires elements")
	ruleListStyleValue      = rule("rich_text_list.style.value", fatal, true, referenceElements+"rich-text-list-element/", "list style must be bullet or ordered")
	ruleListOffsetOrdered   = rule("rich_text_list.offset.ordered", warning, false, referenceElements+"rich-text-list-element/", "offset numbers only an ordered list")
	ruleListNumberValue     = rule("rich_text_list.number.value", warning, false, referenceElements+"rich-text-list-element/", "indent and offset must not be negative")
	ruleBorderValue         = rule("rich_text.border.value", warning, false, referenceElements+"rich-text-list-element/", "border turns the border off with 0 or on with 1")
	ruleInlineTextRequired  = rule("text_inline.text.required", warning, false, referenceElements+"text-element/", "rich text inline text is required")
	ruleLinkURLFormat       = rule("link.url.format", warning, false, referenceElements+"link-element/", "link url must be an absolute URI")
	ruleLinkStyleCode       = rule("link.style.code", fatal, true, referenceElements+"link-element/", "link style does not accept code")
	ruleUserIDRequired      = rule("user.user_id.required", warning, false, referenceElements+"user-element/", "user_id is required")
	ruleUserStyleCode       = rule("user.style.code", fatal, true, referenceElements+"user-element/", "user style does not accept code")
	ruleEmojiNameRequired   = rule("emoji.name.required", warning, false, referenceElements+"emoji-element/", "emoji name is required")
	ruleRichSectionRequired = rule("rich_text_section.elements.required", warning, false, referenceElements+"rich-text-section-element/", "a rich text section requires elements")
)

// rules lists every Slack rule. The README documents each entry.
var rules = []validation.Rule{
	ruleMessageContentRequired, ruleMessageBlocksCount, ruleMessageTextFallback, ruleMessageTextLength,
	ruleMarkdownTextTotal, ruleMarkdownTextRequired, ruleElementRequired, ruleBlockIDLength,
	ruleTextRequired, ruleTextLength,
	ruleHeaderTextLength, ruleHeaderLevelValue, ruleSectionContent, ruleSectionFieldCount, ruleSectionFieldLen,
	ruleImageSource, ruleImageURLLength, ruleImageURLFormat, ruleImageAltRequired, ruleImageAltLength, ruleImageTitleLength,
	ruleSlackFileSource, ruleSlackFileURL,
	ruleActionsCount, ruleActionsRequired, ruleContextCount, ruleContextRequired,
	ruleButtonTextLength, ruleButtonTextTruncation, ruleButtonActionIDLength, ruleButtonURLLength, ruleButtonURLFormat,
	ruleButtonValueLength, ruleButtonStyleValue, ruleButtonLabelLength, ruleButtonPromptLength, ruleButtonPromptDisplay, ruleButtonVisibleRequired,
	ruleConfirmTitleLength, ruleConfirmTextRequired, ruleConfirmTextLength, ruleConfirmConfirmLength, ruleConfirmDenyLength, ruleConfirmStyleValue,
	ruleVideoTitleLength, ruleVideoDescriptionLength, ruleVideoAuthorLength, ruleVideoURLHTTPS, ruleVideoTitleURLHTTPS,
	ruleVideoAltRequired, ruleVideoThumbnailFormat, ruleVideoProviderIcon,
	ruleRichTextRequired, ruleListStyleValue, ruleListOffsetOrdered, ruleListNumberValue, ruleBorderValue,
	ruleInlineTextRequired, ruleLinkURLFormat, ruleLinkStyleCode, ruleUserIDRequired, ruleUserStyleCode, ruleEmojiNameRequired, ruleRichSectionRequired,
}

func runes(text string) int { return validation.Runes(text) }

func httpURL(text string) bool { return validation.AbsoluteURI(text, "http", "https") }

// textLength reports rule when a text object's text exceeds maximum characters.
func textLength(c *validation.Check, text TextObject, maximum int, rule validation.Rule, field string) {
	if value := validation.Value(text); value != nil && runes(text.String()) > maximum {
		c.Fail(rule, field)
	}
}

func blockID(c *validation.Check, id string) {
	c.When(runes(id) > 255, ruleBlockIDLength, "block_id")
}

func (m Message) check(c *validation.Check) {
	c.When(m.Text == "" && len(m.Blocks) == 0, ruleMessageContentRequired, "")
	c.When(m.Text == "" && len(m.Blocks) != 0, ruleMessageTextFallback, "text")
	c.When(runes(m.Text) > 40000, ruleMessageTextLength, "text")
	c.When(len(m.Blocks) > 50, ruleMessageBlocksCount, "blocks")
	markdown := 0
	for i, block := range m.Blocks {
		field := fmt.Sprintf("blocks[%d]", i)
		if validation.Value(block) == nil {
			c.Fail(ruleElementRequired, field)
			continue
		}
		if value, ok := validation.Value(block).(MarkdownBlock); ok {
			markdown += runes(value.Text)
		}
		c.Child(field, block.check)
	}
	c.When(markdown > 12000, ruleMarkdownTextTotal, "blocks")
}

func (t PlainTextObject) check(c *validation.Check) {
	c.When(t.Text == "", ruleTextRequired, "text")
	c.When(runes(t.Text) > 3000, ruleTextLength, "text")
}

func (t MrkdwnTextObject) check(c *validation.Check) {
	c.When(t.Text == "", ruleTextRequired, "text")
	c.When(runes(t.Text) > 3000, ruleTextLength, "text")
}

func (b HeaderBlock) check(c *validation.Check) {
	textLength(c, b.Text, 150, ruleHeaderTextLength, "text")
	c.When(b.Level != 0 && (b.Level < 1 || b.Level > 4), ruleHeaderLevelValue, "level")
	blockID(c, b.BlockID)
	c.Child("text", b.Text.check)
}

func (b SectionBlock) check(c *validation.Check) {
	c.When(validation.Value(b.Text) == nil && len(b.Fields) == 0, ruleSectionContent, "")
	c.When(len(b.Fields) > 10, ruleSectionFieldCount, "fields")
	blockID(c, b.BlockID)
	if validation.Value(b.Text) != nil {
		c.Child("text", b.Text.check)
	}
	for i, field := range b.Fields {
		name := fmt.Sprintf("fields[%d]", i)
		if validation.Value(field) == nil {
			c.Fail(ruleElementRequired, name)
			continue
		}
		textLength(c, field, 2000, ruleSectionFieldLen, name)
		c.Child(name, field.check)
	}
	if validation.Value(b.Accessory) != nil {
		c.Child("accessory", b.Accessory.check)
	}
}

func checkImage(c *validation.Check, url string, file *SlackFileObject, alt string) {
	c.When((url == "") == (file == nil), ruleImageSource, "")
	if url != "" {
		c.When(runes(url) > 3000, ruleImageURLLength, "image_url")
		c.When(!httpURL(url), ruleImageURLFormat, "image_url")
	}
	c.When(alt == "", ruleImageAltRequired, "alt_text")
	c.When(runes(alt) > 2000, ruleImageAltLength, "alt_text")
	if file != nil {
		c.Child("slack_file", file.check)
	}
}

func (b ImageBlock) check(c *validation.Check) {
	checkImage(c, b.ImageURL, b.SlackFile, b.AltText)
	blockID(c, b.BlockID)
	if b.Title != nil {
		c.When(runes(b.Title.Text) > 2000, ruleImageTitleLength, "title")
		c.Child("title", b.Title.check)
	}
}

func (e ImageElement) check(c *validation.Check) {
	checkImage(c, e.ImageURL, e.SlackFile, e.AltText)
}

func (f SlackFileObject) check(c *validation.Check) {
	c.When((f.URL == "") == (f.ID == ""), ruleSlackFileSource, "")
	c.When(f.URL != "" && !httpURL(f.URL), ruleSlackFileURL, "url")
}

func (b ActionsBlock) check(c *validation.Check) {
	c.When(len(b.Elements) == 0, ruleActionsRequired, "elements")
	c.When(len(b.Elements) > 25, ruleActionsCount, "elements")
	blockID(c, b.BlockID)
	for i, element := range b.Elements {
		name := fmt.Sprintf("elements[%d]", i)
		if validation.Value(element) == nil {
			c.Fail(ruleElementRequired, name)
			continue
		}
		c.Child(name, element.check)
	}
}

func (b ContextBlock) check(c *validation.Check) {
	c.When(len(b.Elements) == 0, ruleContextRequired, "elements")
	c.When(len(b.Elements) > 10, ruleContextCount, "elements")
	blockID(c, b.BlockID)
	for i, element := range b.Elements {
		name := fmt.Sprintf("elements[%d]", i)
		if validation.Value(element) == nil {
			c.Fail(ruleElementRequired, name)
			continue
		}
		c.Child(name, element.check)
	}
}

func (b DividerBlock) check(c *validation.Check) { blockID(c, b.BlockID) }

func (b MarkdownBlock) check(c *validation.Check) {
	c.When(b.Text == "", ruleMarkdownTextRequired, "text")
}

func (e ButtonElement) check(c *validation.Check) {
	c.When(runes(e.Text.Text) > 75, ruleButtonTextLength, "text")
	c.When(runes(e.Text.Text) > 30, ruleButtonTextTruncation, "text")
	c.When(runes(e.ActionID) > 255, ruleButtonActionIDLength, "action_id")
	if e.URL != "" {
		c.When(runes(e.URL) > 3000, ruleButtonURLLength, "url")
		c.When(!validation.AbsoluteURI(e.URL), ruleButtonURLFormat, "url")
	}
	c.When(runes(e.Value) > 2000, ruleButtonValueLength, "value")
	c.When(!buttonStyles[e.Style], ruleButtonStyleValue, "style")
	c.When(runes(e.AccessibilityLabel) > 75, ruleButtonLabelLength, "accessibility_label")
	c.When(runes(e.AgentPrompt) > 4000, ruleButtonPromptLength, "agent_prompt")
	c.When(e.AgentPromptDisplay != "" && e.AgentPrompt == "", ruleButtonPromptDisplay, "agent_prompt_display")
	for i, id := range e.VisibleToUserIDs {
		c.When(id == "", ruleButtonVisibleRequired, fmt.Sprintf("visible_to_user_ids[%d]", i))
	}
	c.Child("text", e.Text.check)
	if e.Confirm != nil {
		c.Child("confirm", e.Confirm.check)
	}
}

func (d ConfirmationDialogObject) check(c *validation.Check) {
	textLength(c, d.Title, 100, ruleConfirmTitleLength, "title")
	textLength(c, d.Confirm, 30, ruleConfirmConfirmLength, "confirm")
	textLength(c, d.Deny, 30, ruleConfirmDenyLength, "deny")
	c.When(!buttonStyles[d.Style], ruleConfirmStyleValue, "style")
	c.Child("title", d.Title.check)
	if validation.Value(d.Text) == nil {
		c.Fail(ruleConfirmTextRequired, "text")
	} else {
		textLength(c, d.Text, 300, ruleConfirmTextLength, "text")
		c.Child("text", d.Text.check)
	}
	c.Child("confirm", d.Confirm.check)
	c.Child("deny", d.Deny.check)
}

func (b VideoBlock) check(c *validation.Check) {
	c.When(runes(b.Title.Text) >= 200, ruleVideoTitleLength, "title")
	c.When(b.Description != nil && runes(b.Description.Text) >= 200, ruleVideoDescriptionLength, "description")
	c.When(runes(b.AuthorName) >= 50, ruleVideoAuthorLength, "author_name")
	c.When(!validation.AbsoluteURI(b.VideoURL, "https"), ruleVideoURLHTTPS, "video_url")
	c.When(b.TitleURL != "" && !validation.AbsoluteURI(b.TitleURL, "https"), ruleVideoTitleURLHTTPS, "title_url")
	c.When(b.AltText == "", ruleVideoAltRequired, "alt_text")
	c.When(!httpURL(b.ThumbnailURL), ruleVideoThumbnailFormat, "thumbnail_url")
	c.When(b.ProviderIconURL != "" && !httpURL(b.ProviderIconURL), ruleVideoProviderIcon, "provider_icon_url")
	blockID(c, b.BlockID)
	c.Child("title", b.Title.check)
	if b.Description != nil {
		c.Child("description", b.Description.check)
	}
}

func (b RichTextBlock) check(c *validation.Check) {
	c.When(len(b.Elements) == 0, ruleRichTextRequired, "elements")
	blockID(c, b.BlockID)
	for i, element := range b.Elements {
		name := fmt.Sprintf("elements[%d]", i)
		if validation.Value(element) == nil {
			c.Fail(ruleElementRequired, name)
			continue
		}
		c.Child(name, element.check)
	}
}

func checkInlines[T interface{ check(*validation.Check) }](c *validation.Check, elements []T) {
	for i, inline := range elements {
		name := fmt.Sprintf("elements[%d]", i)
		if validation.Value(inline) == nil {
			c.Fail(ruleElementRequired, name)
			continue
		}
		c.Child(name, inline.check)
	}
}

func checkBorder(c *validation.Check, border *int) {
	c.When(border != nil && *border != 0 && *border != 1, ruleBorderValue, "border")
}

func (e RichTextSection) check(c *validation.Check) {
	c.When(len(e.Elements) == 0, ruleRichSectionRequired, "elements")
	checkInlines(c, e.Elements)
}

func (e RichTextList) check(c *validation.Check) {
	c.When(e.Style != RichTextListStyleBullet && e.Style != RichTextListStyleOrdered, ruleListStyleValue, "style")
	c.When(len(e.Elements) == 0, ruleRichSectionRequired, "elements")
	c.When(e.Offset != 0 && e.Style != RichTextListStyleOrdered, ruleListOffsetOrdered, "offset")
	c.When(e.Indent < 0, ruleListNumberValue, "indent")
	c.When(e.Offset < 0, ruleListNumberValue, "offset")
	checkBorder(c, e.Border)
	for i, section := range e.Elements {
		c.Child(fmt.Sprintf("elements[%d]", i), section.check)
	}
}

func (e RichTextPreformatted) check(c *validation.Check) {
	c.When(len(e.Elements) == 0, ruleRichSectionRequired, "elements")
	checkBorder(c, e.Border)
	checkInlines(c, e.Elements)
}

func (e RichTextQuote) check(c *validation.Check) {
	c.When(len(e.Elements) == 0, ruleRichSectionRequired, "elements")
	checkBorder(c, e.Border)
	checkInlines(c, e.Elements)
}

func (e TextInline) check(c *validation.Check) {
	c.When(e.Text == "", ruleInlineTextRequired, "text")
}

func (e LinkInline) check(c *validation.Check) {
	c.When(!validation.AbsoluteURI(e.URL), ruleLinkURLFormat, "url")
	c.When(e.Style != nil && e.Style.Code, ruleLinkStyleCode, "style.code")
}

func (e UserInline) check(c *validation.Check) {
	c.When(e.UserID == "", ruleUserIDRequired, "user_id")
	c.When(e.Style != nil && e.Style.Code, ruleUserStyleCode, "style.code")
}

func (e EmojiInline) check(c *validation.Check) {
	c.When(e.Name == "", ruleEmojiNameRequired, "name")
}

var buttonStyles = map[ButtonStyle]bool{"": true, ButtonStylePrimary: true, ButtonStyleDanger: true}
