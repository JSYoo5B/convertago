package kakaowork

import (
	"fmt"
	"strings"

	"github.com/JSYoo5B/convertago/internal/validation"
)

const (
	referenceBlockKit = "https://docs.kakaoi.ai/kakao_work/blockkit/"
	referenceMessages = "https://docs.kakaoi.ai/kakao_work/webapireference/messages/"
	referenceUXGuide  = referenceBlockKit + "uxguide/"
)

func rule(id string, severity validation.Severity, verified bool, reference, message string) validation.Rule {
	return validation.Rule{ID: id, Severity: severity, Verified: verified, Reference: reference, Message: message}
}

const (
	fatal    = validation.Fatal
	warning  = validation.Warning
	advisory = validation.Advisory
)

// rules lists every Kakao Work rule. The README documents each entry.
var rules = []validation.Rule{
	ruleMessageTextRequired, ruleMessageTextLength, ruleMessageBlockType,
	ruleMessageButtonCount, ruleMessageImageAfterButton, ruleMessageImageConsecutive,
	ruleMessageDividerEdge, ruleMessageDividerOnly,
	ruleHeaderPosition, ruleHeaderTextRequired, ruleHeaderTextLength, ruleHeaderTextLineBreak, ruleHeaderStyleValue,
	ruleTextTextRequired, ruleTextTextLength, ruleTextInlinesMismatch, ruleTextInlineRequired, ruleTextLinkCount,
	ruleStyledTextRequired, ruleStyledColorValue, ruleLinkTextRequired, ruleLinkURLScheme,
	ruleMentionTextRequired, ruleMentionUserID,
	ruleImageURLFormat,
	ruleButtonTextRequired, ruleButtonTextLength, ruleButtonStyleValue, ruleButtonActionRequired,
	ruleActionElementsCount,
	ruleDescriptionTermRequired, ruleDescriptionTermLength,
	ruleSystemBrowserURLFormat, ruleInAppBrowserURLFormat, ruleInAppBrowserSizeStandalone, ruleInAppBrowserSizeValue,
	ruleExternalAppValueFormat, ruleSubmitNameRequired, ruleSubmitValueRequired, ruleCallModalValueRequired,
	ruleExclusiveDefaultRequired, ruleExclusiveActionType,
}

var (
	ruleMessageTextRequired     = rule("message.text.required", warning, false, referenceMessages, "message text is required")
	ruleMessageTextLength       = rule("message.text.length", warning, false, referenceMessages, "message text exceeds 10000 characters")
	ruleMessageBlockType        = rule("message.blocks.type", fatal, true, referenceBlockKit, "blocks accept only message bubble blocks")
	ruleMessageButtonCount      = rule("message.buttons.count", advisory, true, referenceUXGuide, "avoid four or more buttons in one message")
	ruleMessageImageAfterButton = rule("message.image_link.after_button", advisory, true, referenceUXGuide, "do not place an image link block below a button block")
	ruleMessageImageConsecutive = rule("message.image_link.consecutive", advisory, true, referenceUXGuide, "avoid consecutive image link blocks")
	ruleMessageDividerEdge      = rule("message.divider.edge", advisory, true, referenceUXGuide, "avoid a divider block at the top or bottom of a message")
	ruleMessageDividerOnly      = rule("message.divider.only", advisory, true, referenceUXGuide, "do not compose a message only of divider blocks")

	ruleHeaderPosition      = rule("header.position", warning, false, referenceBlockKit+"headerblock/", "a header block must be the first and only header")
	ruleHeaderTextRequired  = rule("header.text.required", warning, false, referenceBlockKit+"headerblock/", "header text is required")
	ruleHeaderTextLength    = rule("header.text.length", advisory, true, referenceBlockKit+"headerblock/", "header text over 20 characters is truncated")
	ruleHeaderTextLineBreak = rule("header.text.line_break", warning, false, referenceBlockKit+"headerblock/", "header text does not support line breaks")
	ruleHeaderStyleValue    = rule("header.style.value", fatal, true, referenceBlockKit+"headerblock/", "unknown header style")

	ruleTextTextRequired    = rule("text.text.required", warning, false, referenceBlockKit+"textblock/", "text is required")
	ruleTextTextLength      = rule("text.text.length", warning, false, referenceBlockKit+"textblock/", "text or its inlines exceed 500 characters")
	ruleTextInlinesMismatch = rule("text.inlines.mismatch", warning, true, referenceBlockKit+"textblock/", "text differs from its inlines, which take priority")
	ruleTextInlineRequired  = rule("text.inlines.required", fatal, true, referenceBlockKit+"textblock/", "an inline element is missing")
	ruleTextLinkCount       = rule("text.links.count", advisory, true, referenceUXGuide, "do not include two or more links in one text block")
	ruleStyledTextRequired  = rule("styled.text.required", warning, false, referenceBlockKit+"textblock/", "styled text is required")
	ruleStyledColorValue    = rule("styled.color.value", fatal, true, referenceBlockKit+"textblock/", "unknown inline color")
	ruleLinkTextRequired    = rule("link.text.required", warning, false, referenceBlockKit+"textblock/", "link text is required")
	ruleLinkURLScheme       = rule("link.url.scheme", warning, false, referenceBlockKit+"textblock/", "link URL must be an absolute HTTP, HTTPS, mailto, or tel URI")
	ruleMentionTextRequired = rule("mention.text.required", warning, false, referenceBlockKit+"textblock/", "mention text is required")
	ruleMentionUserID       = rule("mention.user_id.value", advisory, true, referenceBlockKit+"textblock/", "a mention works only for a participant, whose user ID is positive")

	ruleImageURLFormat = rule("image_link.url.format", warning, false, referenceBlockKit+"imagelinkblock/", "image URL must be an absolute HTTP or HTTPS URL")

	ruleButtonTextRequired   = rule("button.text.required", warning, false, referenceBlockKit+"buttonblock/", "button text is required")
	ruleButtonTextLength     = rule("button.text.length", advisory, true, referenceBlockKit+"buttonblock/", "button text over 20 characters is truncated")
	ruleButtonStyleValue     = rule("button.style.value", fatal, true, referenceBlockKit+"buttonblock/", "unknown button style")
	ruleButtonActionRequired = rule("button.action.required", fatal, true, referenceBlockKit+"buttonblock/", "button action is required")
	ruleActionElementsCount  = rule("action.elements.count", warning, false, referenceBlockKit+"actionblock/", "an action block holds two or three buttons")

	ruleDescriptionTermRequired = rule("description.term.required", warning, false, referenceBlockKit+"descriptionblock/", "description term is required")
	ruleDescriptionTermLength   = rule("description.term.length", warning, false, referenceBlockKit+"descriptionblock/", "description term exceeds 10 characters")

	ruleSystemBrowserURLFormat     = rule("open_system_browser.value.format", warning, false, referenceBlockKit+"buttonblock/", "browser URL must be an absolute HTTP or HTTPS URL")
	ruleInAppBrowserURLFormat      = rule("open_inapp_browser.value.format", warning, false, referenceBlockKit+"buttonblock/", "browser URL must be an absolute HTTP or HTTPS URL")
	ruleInAppBrowserSizeStandalone = rule("open_inapp_browser.size.standalone", warning, false, referenceBlockKit+"buttonblock/", "width and height apply only to a standalone window")
	ruleInAppBrowserSizeValue      = rule("open_inapp_browser.size.value", warning, false, referenceBlockKit+"buttonblock/", "width and height must be positive pixel values")
	ruleExternalAppValueFormat     = rule("open_external_app.value.format", warning, false, referenceBlockKit+"buttonblock/", "external app value must be a query string of ios and aos app URIs")
	ruleSubmitNameRequired         = rule("submit_action.name.required", warning, false, referenceBlockKit+"buttonblock/", "submit action name is required")
	ruleSubmitValueRequired        = rule("submit_action.value.required", warning, false, referenceBlockKit+"buttonblock/", "submit action value is required")
	ruleCallModalValueRequired     = rule("call_modal.value.required", warning, false, referenceBlockKit+"buttonblock/", "call modal value is required")
	ruleExclusiveDefaultRequired   = rule("exclusive.default.required", fatal, true, referenceBlockKit+"buttonblock/", "exclusive action requires a default action")
	ruleExclusiveActionType        = rule("exclusive.action.type", fatal, true, referenceBlockKit+"buttonblock/", "exclusive action slots accept only compatibility actions")
)

func (m Message) check(c *validation.Check) {
	c.When(m.Preview == "", ruleMessageTextRequired, "text")
	c.When(validation.Runes(m.Preview) > 10000, ruleMessageTextLength, "text")
	buttons, dividers := 0, 0
	for i, block := range m.Blocks {
		field := fmt.Sprintf("blocks[%d]", i)
		switch value := validation.Value(block).(type) {
		case HeaderBlock:
			c.When(i != 0, ruleHeaderPosition, field)
		case ButtonBlock:
			buttons++
		case ActionBlock:
			buttons += len(value.Elements)
		case ImageBlock:
			if i > 0 {
				switch validation.Value(m.Blocks[i-1]).(type) {
				case ButtonBlock:
					c.Fail(ruleMessageImageAfterButton, field)
				case ImageBlock:
					c.Fail(ruleMessageImageConsecutive, field)
				}
			}
		case DividerBlock:
			dividers++
			c.When(i == 0 || i == len(m.Blocks)-1, ruleMessageDividerEdge, field)
		case TextBlock, DescriptionBlock, SectionBlock, ContextBlock:
		default:
			c.Fail(ruleMessageBlockType, field)
			continue
		}
		c.Child(field, block.check)
	}
	c.When(buttons >= 4, ruleMessageButtonCount, "blocks")
	c.When(dividers != 0 && dividers == len(m.Blocks), ruleMessageDividerOnly, "blocks")
}

func (h HeaderBlock) check(c *validation.Check) {
	c.When(h.Text == "", ruleHeaderTextRequired, "text")
	c.When(validation.Runes(h.Text) > 20, ruleHeaderTextLength, "text")
	c.When(strings.ContainsAny(h.Text, "\r\n"), ruleHeaderTextLineBreak, "text")
	c.When(!headerStyles[h.Style], ruleHeaderStyleValue, "style")
}

func (t TextBlock) check(c *validation.Check) {
	c.When(t.Text == "", ruleTextTextRequired, "text")
	total, links := 0, 0
	var joined strings.Builder
	for i, inline := range t.Inlines {
		field := fmt.Sprintf("inlines[%d]", i)
		if validation.Value(inline) == nil {
			c.Fail(ruleTextInlineRequired, field)
			continue
		}
		total += validation.Runes(inline.String())
		joined.WriteString(inline.String())
		if _, ok := validation.Value(inline).(InlineLink); ok {
			links++
		}
		c.Child(field, inline.check)
	}
	c.When(validation.Runes(t.Text) > 500 || total > 500, ruleTextTextLength, "text")
	c.When(len(t.Inlines) != 0 && joined.String() != t.Text, ruleTextInlinesMismatch, "inlines")
	c.When(links >= 2, ruleTextLinkCount, "inlines")
}

func (i InlineStyled) check(c *validation.Check) {
	c.When(i.Text == "", ruleStyledTextRequired, "text")
	c.When(!inlineColors[i.Color], ruleStyledColorValue, "color")
}

func (i InlineLink) check(c *validation.Check) {
	c.When(i.Text == "", ruleLinkTextRequired, "text")
	c.When(!validation.AbsoluteURI(i.Url, "http", "https", "mailto", "tel"), ruleLinkURLScheme, "url")
}

func (i InlineMention) check(c *validation.Check) {
	c.When(i.Text == "", ruleMentionTextRequired, "text")
	c.When(i.UserId <= 0, ruleMentionUserID, "ref.value")
}

func (i ImageBlock) check(c *validation.Check) {
	c.When(!validation.AbsoluteURI(i.Url, "http", "https"), ruleImageURLFormat, "url")
}

func (b ButtonBlock) check(c *validation.Check) {
	c.When(b.Text == "", ruleButtonTextRequired, "text")
	c.When(validation.Runes(b.Text) > 20, ruleButtonTextLength, "text")
	c.When(!buttonStyles[b.Style], ruleButtonStyleValue, "style")
	if validation.Value(b.Action) == nil {
		c.Fail(ruleButtonActionRequired, "action")
		return
	}
	c.Child("action", b.Action.check)
}

func (a ActionBlock) check(c *validation.Check) {
	c.When(len(a.Elements) < 2 || len(a.Elements) > 3, ruleActionElementsCount, "elements")
	for i, button := range a.Elements {
		c.Child(fmt.Sprintf("elements[%d]", i), button.check)
	}
}

func (DividerBlock) check(*validation.Check) {}

func (d DescriptionBlock) check(c *validation.Check) {
	c.When(d.Term == "", ruleDescriptionTermRequired, "term")
	c.When(validation.Runes(d.Term) > 10, ruleDescriptionTermLength, "term")
	c.Child("content", d.Content.check)
}

func (s SectionBlock) check(c *validation.Check) {
	c.Child("content", s.Content.check)
	if s.Accessory != nil {
		c.Child("accessory", s.Accessory.check)
	}
	if validation.Value(s.Action) != nil {
		c.Child("action", s.Action.check)
	}
}

func (b ContextBlock) check(c *validation.Check) {
	c.Child("content", b.Content.check)
	c.Child("image", b.Image.check)
}

func (o OpenSystemBrowserAction) check(c *validation.Check) {
	c.When(!validation.AbsoluteURI(o.Value, "http", "https"), ruleSystemBrowserURLFormat, "value")
}

func (o OpenInAppBrowserAction) check(c *validation.Check) {
	c.When(!validation.AbsoluteURI(o.Value, "http", "https"), ruleInAppBrowserURLFormat, "value")
	c.When(!o.Standalone && (o.Width != 0 || o.Height != 0), ruleInAppBrowserSizeStandalone, "standalone")
	c.When(o.Width < 0, ruleInAppBrowserSizeValue, "width")
	c.When(o.Height < 0, ruleInAppBrowserSizeValue, "height")
}

func (o OpenExternalAppAction) check(c *validation.Check) {
	c.When(!validation.QueryURIs(o.Value, "ios", "aos"), ruleExternalAppValueFormat, "value")
}

func (s SubmitAction) check(c *validation.Check) {
	c.When(s.Name == "", ruleSubmitNameRequired, "name")
	c.When(s.Value == "", ruleSubmitValueRequired, "value")
}

func (m CallModalAction) check(c *validation.Check) {
	c.When(m.Value == "", ruleCallModalValueRequired, "value")
}

func (e ExclusiveAction) check(c *validation.Check) {
	for _, slot := range []struct {
		field  string
		action ButtonAction
	}{{"default", e.Default}, {"pc", e.Pc}, {"mobile", e.Mobile}, {"windows", e.Windows}, {"macos", e.MacOs}, {"android", e.Android}, {"ios", e.Ios}} {
		switch validation.Value(slot.action).(type) {
		case nil:
			c.When(slot.field == "default", ruleExclusiveDefaultRequired, slot.field)
		case ExclusiveAction:
			c.Fail(ruleExclusiveActionType, slot.field)
		default:
			c.Child(slot.field, slot.action.check)
		}
	}
}
