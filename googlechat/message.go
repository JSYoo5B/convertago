// Package googlechat provides Cards v2 objects for composing Google Chat messages.
package googlechat

// Message contains text and card interfaces sent by a Google Chat app.
// The message's cards have a combined size limit of 32 KB.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/spaces.messages
type Message struct {
	// Text is the message's plain-text body, which can use Chat's text formatting.
	Text string `json:"text,omitempty"`
	// CardsV2 contains the cards displayed in the message.
	CardsV2 []CardWithID `json:"cardsV2,omitempty" validate:"dive"`
	// FallbackText describes the cards when a client cannot display them.
	FallbackText string `json:"fallbackText,omitempty"`
}

// CardWithID associates a card with an identifier used when updating a message.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/spaces.messages#CardWithId
type CardWithID struct {
	// CardID uniquely identifies a card within the message and is required when there are multiple cards.
	CardID string `json:"cardId,omitempty"`
	// Card is the card interface to display.
	Card Card `json:"card"`
}
