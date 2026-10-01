package googlechat_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestTextParagraph_Validate(t *testing.T) {
	v := validator.New()
	require.NoError(t, v.Struct(googlechat.TextParagraph{Text: "All lines"}))
	require.NoError(t, v.Struct(googlechat.TextParagraph{Text: "Markdown", TextSyntax: googlechat.TextSyntaxMarkdown, MaxLines: 2}))
	require.Error(t, v.Struct(googlechat.TextParagraph{MaxLines: -1}))
	require.Error(t, v.Struct(googlechat.TextParagraph{TextSyntax: "PLAIN"}))
}
