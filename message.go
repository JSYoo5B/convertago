package convertago

import (
	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/JSYoo5B/convertago/slack"
)

// ToKakaoworkMessage converts a struct using its kakaowork tags.
// It uses generated field accessors when available and cached reflection otherwise.
func ToKakaoworkMessage(input any, options ...Option) (kakaowork.Message, error) {
	return kakaowork.ToMessage(input, options...)
}

// ToSlackMessage converts a struct using its slack tags.
// It uses generated field accessors when available and cached reflection otherwise.
func ToSlackMessage(input any, options ...Option) (slack.Message, error) {
	return slack.ToMessage(input, options...)
}

// ToGoogleChatMessage converts a struct using its googlechat tags.
// It supplies a card and a section for the collected widgets.
func ToGoogleChatMessage(input any, options ...Option) (googlechat.Message, error) {
	return googlechat.ToMessage(input, options...)
}
