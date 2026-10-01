package googlechat

import (
	"encoding/json"
	"fmt"
)

// OnClick performs one action when a user clicks an interactive card component.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#OnClick
type OnClick struct {
	// Action invokes an app function. Supply exactly one of Action, OpenLink, or OverflowMenu.
	Action *Action `json:"action,omitempty"`
	// OpenLink opens a hyperlink.
	OpenLink *OpenLink `json:"openLink,omitempty"`
	// OverflowMenu opens a menu of secondary actions.
	OverflowMenu *OverflowMenu `json:"overflowMenu,omitempty"`
}

func (a OnClick) MarshalJSON() ([]byte, error) {
	count := 0
	if a.Action != nil {
		count++
	}
	if a.OpenLink != nil {
		count++
	}
	if a.OverflowMenu != nil {
		count++
	}
	if count != 1 {
		return nil, fmt.Errorf("googlechat: OnClick requires exactly one action")
	}
	type Embed OnClick
	return json.Marshal(Embed(a))
}

// OpenLink opens a hyperlink when the containing component is activated.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#OpenLink
type OpenLink struct {
	// URL is the link to open.
	URL string `json:"url"`
}

// Action invokes a custom function with optional parameters and form requirements.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Action
type Action struct {
	// Function identifies the app function to invoke.
	Function string `json:"function"`
	// Parameters supplies string values to the function.
	Parameters []ActionParameter `json:"parameters,omitempty"`
	// LoadIndicator selects the loading feedback. Omission uses a spinner.
	LoadIndicator LoadIndicator `json:"loadIndicator,omitempty"`
	// PersistValues retains form values after the action. Card messages also need an UPDATE_MESSAGE
	// response using the original CardID. LoadIndicatorNone allows edits while the action runs.
	PersistValues bool `json:"persistValues,omitempty"`
	// Interaction allows a card-message button to open a dialog.
	Interaction Interaction `json:"interaction,omitempty"`
	// RequiredWidgets names inputs that must have values before submission.
	RequiredWidgets []string `json:"requiredWidgets,omitempty"`
	// AllWidgetsAreRequired requires values for all inputs before submission.
	AllWidgetsAreRequired bool `json:"allWidgetsAreRequired,omitempty"`
}

// ActionParameter supplies a named string value to an Action.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#ActionParameter
type ActionParameter struct {
	// Key names the parameter passed to the function.
	Key string `json:"key"`
	// Value is the parameter's string value.
	Value string `json:"value"`
}

// LoadIndicator selects feedback while an action runs.
type LoadIndicator string

const (
	LoadIndicatorSpinner LoadIndicator = "SPINNER"
	LoadIndicatorNone    LoadIndicator = "NONE"
)

// Interaction selects a special response to a card-message interaction.
type Interaction string

const InteractionOpenDialog Interaction = "OPEN_DIALOG"

// OverflowMenu presents secondary actions in a pop-up menu.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#OverflowMenu
type OverflowMenu struct {
	// Items contains the actions offered by the menu.
	Items []OverflowMenuItem `json:"items"`
}

// OverflowMenuItem displays text and an optional icon for a menu action.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#OverflowMenuItem
type OverflowMenuItem struct {
	// StartIcon appears before the text.
	StartIcon *Icon `json:"startIcon,omitempty"`
	// Text is the required label describing the menu item.
	Text string `json:"text"`
	// OnClick is the required action. An OverflowMenu action is dropped and disables the item.
	OnClick OnClick `json:"onClick"`
	// Disabled prevents the menu option from responding to user actions.
	Disabled bool `json:"disabled,omitempty"`
}
