package conversion

import (
	"fmt"
	"unicode/utf8"
)

func ValidateText(platform, path, text string, minimum, maximum int) error {
	if !utf8.ValidString(text) {
		return Error(platform, path, "invalid_value", "text must be valid UTF-8")
	}
	length := utf8.RuneCountInString(text)
	if length < minimum {
		return Error(platform, path, "invalid_value", fmt.Sprintf("text must contain at least %d characters", minimum))
	}
	if maximum > 0 && length > maximum {
		return Error(platform, path, "invalid_value", fmt.Sprintf("text exceeds %d characters", maximum))
	}
	return nil
}
