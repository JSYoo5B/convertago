package validation_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/internal/validation"
	"github.com/stretchr/testify/require"
)

func TestQueryURIs(t *testing.T) {
	for _, test := range []struct {
		text  string
		valid bool
	}{
		{"ios=kakaomap%3A%2F%2Flook&aos=kakaomap%3A%2F%2Flook", true},
		{"ios=kakaomap%3A%2F%2Flook", true},
		{"kakaomap://look", false},
		{"web=https%3A%2F%2Fexample.com", false},
		{"ios=relative", false},
		{"ios=a%3A%2F%2Fb&ios=c%3A%2F%2Fd", false},
		{"", false},
	} {
		require.Equal(t, test.valid, validation.QueryURIs(test.text, "ios", "aos"), test.text)
	}
}

func TestWalkJoinsChildPaths(t *testing.T) {
	rule := validation.Rule{ID: "leaf", Severity: validation.Warning}
	leaf := func(c *validation.Check) { c.Fail(rule, "text") }
	root := func(c *validation.Check) {
		c.Fail(rule, "")
		c.Child("blocks[0]", leaf)
		c.Child("content", leaf)
	}
	var paths []string
	validation.Walk("$", root, func(path string, v validation.Violation) {
		paths = append(paths, path)
	})
	require.Equal(t, []string{"$", "$.blocks[0].text", "$.content.text"}, paths)
}
