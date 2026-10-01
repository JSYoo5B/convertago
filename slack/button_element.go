package slack

import "encoding/json"

// ButtonElement triggers an interaction or opens a URL from a section or actions block.
//
// Reference: https://docs.slack.dev/reference/block-kit/block-elements/button-element/
type ButtonElement struct {
	// Text is a plain-text label, up to 75 characters, which may truncate near 30 characters.
	Text PlainTextObject `json:"text"`
	// ActionID identifies the interaction within the block, up to 255 characters.
	ActionID string `json:"action_id,omitempty" validate:"max=255"`
	// URL opens in the user's browser, up to 3000 characters. The interaction still needs acknowledgment.
	URL string `json:"url,omitempty" validate:"omitempty,max=3000,url"`
	// Value accompanies the interaction payload, up to 2000 characters.
	Value string `json:"value,omitempty" validate:"max=2000"`
	// Style changes the button's color. An empty value uses the default appearance.
	Style ButtonStyle `json:"style,omitempty" validate:"omitempty,oneof=primary danger"`
	// Confirm requests confirmation before executing the action.
	Confirm *ConfirmationDialogObject `json:"confirm,omitempty" validate:"omitempty"`
	// AccessibilityLabel replaces Text for screen readers, up to 75 characters.
	AccessibilityLabel string `json:"accessibility_label,omitempty" validate:"max=75"`
	// AgentPrompt opens Slackbot with this prompt when the user has Slackbot AI, up to 4000 characters.
	AgentPrompt string `json:"agent_prompt,omitempty" validate:"max=4000"`
	// AgentPromptDisplay labels the resulting Slackbot message instead of showing the full prompt.
	AgentPromptDisplay string `json:"agent_prompt_display,omitempty"`
	// VisibleToUserIDs restricts the button's audience. Omission makes it visible to everyone.
	VisibleToUserIDs []string `json:"visible_to_user_ids,omitempty" validate:"dive,required"`
}

// ButtonStyle selects an alternative button color scheme.
type ButtonStyle string

const (
	// ButtonStylePrimary emphasizes an affirmative action with green text and an outline.
	ButtonStylePrimary ButtonStyle = "primary"
	// ButtonStyleDanger marks a destructive action with red text and an outline.
	ButtonStyleDanger ButtonStyle = "danger"
)

func (e ButtonElement) Type() string   { return "button" }
func (e ButtonElement) String() string { return e.Text.String() }
func (ButtonElement) element()         {}
func (ButtonElement) actionElement()   {}
func (e ButtonElement) MarshalJSON() ([]byte, error) {
	type Embed ButtonElement
	return json.Marshal(struct {
		Type string `json:"type"`
		Embed
	}{e.Type(), Embed(e)})
}

// ConfirmationDialogObject asks the user to confirm or cancel an interaction.
//
// Reference: https://docs.slack.dev/reference/block-kit/composition-objects/confirmation-dialog-object/
type ConfirmationDialogObject struct {
	// Title is a plain-text heading, up to 100 characters.
	Title PlainTextObject `json:"title"`
	// Text explains the action, up to 300 characters.
	Text TextObject `json:"text" validate:"required"`
	// Confirm labels the confirmation button, up to 30 characters.
	Confirm PlainTextObject `json:"confirm"`
	// Deny labels the cancellation button, up to 30 characters.
	Deny PlainTextObject `json:"deny"`
	// Style colors the confirmation button. Slack uses primary when it is omitted.
	Style ButtonStyle `json:"style,omitempty" validate:"omitempty,oneof=primary danger"`
}
