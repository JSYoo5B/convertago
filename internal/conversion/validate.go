package conversion

import "github.com/JSYoo5B/convertago/internal/validation"

func ValidateURL(platform, path, text string, httpsOnly bool) error {
	var valid bool
	if httpsOnly {
		valid = validation.AbsoluteURI(text, "https")
	} else {
		valid = validation.AbsoluteURI(text, "http", "https")
	}
	if valid {
		return nil
	}
	scheme := "HTTP or HTTPS"
	if httpsOnly {
		scheme = "HTTPS"
	}
	return Error(platform, path, "invalid_value", "URL must be an absolute "+scheme+" URL")
}
