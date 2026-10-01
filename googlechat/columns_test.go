package googlechat_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestColumns_Validate(t *testing.T) {
	v := validator.New()
	column := googlechat.Column{Widgets: []googlechat.ColumnWidget{{Content: googlechat.TextParagraph{Text: "Line"}}}}
	require.NoError(t, v.Struct(googlechat.Columns{ColumnItems: []googlechat.Column{column, column}}))
	require.Error(t, v.Struct(googlechat.Columns{ColumnItems: []googlechat.Column{column, column, column}}))
	require.Error(t, v.Struct(googlechat.Columns{}))
	column.HorizontalAlignment = "LEFT"
	require.Error(t, v.Struct(column))
	column.HorizontalAlignment = ""
	var absent *googlechat.TextParagraph
	column.Widgets[0].Content = absent
	require.Error(t, v.Struct(column))
}
