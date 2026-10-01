package googlechat_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestImage_Validate(t *testing.T) {
	v := validator.New()
	googlechat.RegisterValidation(v)
	require.NoError(t, v.Struct(googlechat.Image{ImageURL: "https://example.com/image.png"}))
	require.NoError(t, v.Struct(googlechat.Image{ImageURL: "HtTpS://example.com/image.png"}))
	require.Error(t, v.Struct(googlechat.Image{ImageURL: "http://example.com/image.png"}))
	require.Error(t, v.Struct(googlechat.Image{ImageURL: "/image.png"}))
	require.Error(t, v.Struct(googlechat.Image{ImageURL: "https://example.com/image.png", OnClick: &googlechat.OnClick{}}))
}
