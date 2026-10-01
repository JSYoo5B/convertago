package kakaowork_test

import (
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestDescriptionBlock_Validate(t *testing.T) {
	v := validator.New()
	block := kakaowork.DescriptionBlock{Term: strings.Repeat("가", 10), Content: kakaowork.TextBlock{Text: "설명"}}
	require.NoError(t, v.Struct(block))
	block.Term += "나"
	require.Error(t, v.Struct(block))
	block.Term = ""
	require.Error(t, v.Struct(block))
}
