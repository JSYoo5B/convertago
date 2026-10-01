package conversion

import (
	"net/url"
	"strings"
)

func ValidateURL(platform, path, text string, httpsOnly bool) error {
	u, err := url.ParseRequestURI(text)
	if err == nil && u.Host != "" && (strings.EqualFold(u.Scheme, "https") || (!httpsOnly && strings.EqualFold(u.Scheme, "http"))) {
		return nil
	}
	scheme := "HTTP or HTTPS"
	if httpsOnly {
		scheme = "HTTPS"
	}
	return Error(platform, path, "invalid_value", "image URL must be an absolute "+scheme+" URL")
}
