package googlechat

// TextParagraph displays formatted text with optional line truncation.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#TextParagraph
type TextParagraph struct {
	// Text is the paragraph's content.
	Text string `json:"text"`
	// MaxLines hides excess lines behind an expansion control. Zero displays all lines.
	MaxLines int `json:"maxLines,omitempty"`
	// TextSyntax selects HTML or Markdown. Omission uses HTML.
	TextSyntax TextSyntax `json:"textSyntax,omitempty"`
}

func (TextParagraph) WidgetType() string   { return "textParagraph" }
func (t TextParagraph) String() string     { return t.Text }
func (TextParagraph) widgetContent()       {}
func (TextParagraph) columnWidgetContent() {}
func (TextParagraph) nestedWidgetContent() {}

// TextSyntax selects the formatting syntax of a TextParagraph.
type TextSyntax string

const (
	TextSyntaxHTML     TextSyntax = "HTML"
	TextSyntaxMarkdown TextSyntax = "MARKDOWN"
)
