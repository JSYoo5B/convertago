package convertago

import "github.com/JSYoo5B/convertago/googlechat"

// ToGoogleChatMessage converts a struct using its googlechat tags.
// It supplies a card and a section for the collected widgets.
func ToGoogleChatMessage(input any, options ...Option) (googlechat.Message, error) {
	return googlechat.ToMessage(input, options...)
}
