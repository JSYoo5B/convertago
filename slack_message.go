package convertago

import "github.com/JSYoo5B/convertago/slack"

// ToSlackMessage converts a struct using its slack tags.
// It uses generated field accessors when available and cached reflection otherwise.
func ToSlackMessage(input any, options ...Option) (slack.Message, error) {
	return slack.ToMessage(input, options...)
}
