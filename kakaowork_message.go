package convertago

import "github.com/JSYoo5B/convertago/kakaowork"

// ToKakaoworkMessage converts a struct using its kakaowork tags.
// It uses generated field accessors when available and cached reflection otherwise.
func ToKakaoworkMessage(input any, options ...Option) (kakaowork.Message, error) {
	return kakaowork.ToMessage(input, options...)
}
