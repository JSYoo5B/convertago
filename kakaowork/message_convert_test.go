package kakaowork

import (
	"strings"
	"testing"
)

func TestToMessageValidatesKakaoworkConstraints(t *testing.T) {
	for _, source := range []any{
		struct {
			Text string `kakaowork:"text"`
		}{strings.Repeat("가", 501)},
		struct {
			Text string `kakaowork:"header"`
		}{strings.Repeat("가", 21)},
		struct {
			Text string `kakaowork:"header"`
		}{"two\nlines"},
		struct {
			Text  string `kakaowork:"text"`
			Title string `kakaowork:"header"`
		}{"body", "late"},
		struct {
			URL string `kakaowork:"image_link"`
		}{"/relative.png"},
	} {
		if _, err := ToMessage(source); err == nil {
			t.Fatalf("expected an error for %#v", source)
		}
	}
	source := struct {
		Text string `kakaowork:"text"`
	}{strings.Repeat("가", 500)}
	if _, err := ToMessage(source); err != nil {
		t.Fatalf("500 Unicode characters must be accepted: %v", err)
	}
}
