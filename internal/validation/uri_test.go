package validation_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/internal/validation"
)

func TestAbsoluteURI(t *testing.T) {
	for _, test := range []struct {
		text    string
		schemes []string
		valid   bool
	}{
		{"https://example.com/path?q=1#section", []string{"https"}, true},
		{"HTTPS://example.com", []string{"https"}, true},
		{"http://example.com", []string{"https"}, false},
		{"mailto:hello@example.com", []string{"mailto", "tel"}, true},
		{"tel:+821012345678", []string{"mailto", "tel"}, true},
		{"kakaomap://look?p=1,2", nil, true},
		{"https:example.com", nil, false},
		{"https:///path", nil, false},
		{"https://example.com/%ZZ", nil, false},
		{"https://example.com/a\nb", nil, false},
		{"//example.com", nil, false},
		{"relative", nil, false},
		{"mailto:", nil, false},
		{"", nil, false},
	} {
		t.Run(test.text, func(t *testing.T) {
			if got := validation.AbsoluteURI(test.text, test.schemes...); got != test.valid {
				t.Fatalf("AbsoluteURI = %t, want %t", got, test.valid)
			}
		})
	}
}
