package googlechat_test

import (
	"math"
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestGrid_Validate(t *testing.T) {
	v := validator.New()
	grid := googlechat.Grid{Items: []googlechat.GridItem{{Subtitle: "Item"}}}
	require.NoError(t, v.Struct(grid))
	grid.ColumnCount = -1
	require.Error(t, v.Struct(grid))
	grid.ColumnCount = 0
	grid.Items = []googlechat.GridItem{{}}
	require.Error(t, v.Struct(grid))
	grid.Items = []googlechat.GridItem{{Image: &googlechat.ImageComponent{ImageURI: "http://example.com/image"}}}
	require.NoError(t, v.Struct(grid))
	grid.Items[0].Image.ImageURI = "/image"
	require.Error(t, v.Struct(grid))
}

func TestImageCropStyle_Validate(t *testing.T) {
	v := validator.New()
	googlechat.RegisterValidation(v)
	require.NoError(t, v.Struct(googlechat.ImageCropStyle{}))
	require.NoError(t, v.Struct(googlechat.ImageCropStyle{Type: googlechat.ImageCropTypeRectangleCustom, AspectRatio: 16.0 / 9}))
	require.Error(t, v.Struct(googlechat.ImageCropStyle{Type: googlechat.ImageCropTypeRectangleCustom}))
	require.Error(t, v.Struct(googlechat.ImageCropStyle{Type: googlechat.ImageCropTypeSquare, AspectRatio: 2}))
	require.Error(t, v.Struct(googlechat.ImageCropStyle{Type: googlechat.ImageCropTypeRectangleCustom, AspectRatio: math.Inf(1)}))
}

func TestBorderStyle_Validate(t *testing.T) {
	v := validator.New()
	require.NoError(t, v.Struct(googlechat.BorderStyle{Type: googlechat.BorderTypeStroke, StrokeColor: &googlechat.Color{}, CornerRadius: 4}))
	require.Error(t, v.Struct(googlechat.BorderStyle{Type: googlechat.BorderTypeNone, StrokeColor: &googlechat.Color{}}))
	require.Error(t, v.Struct(googlechat.BorderStyle{CornerRadius: -1}))
}
