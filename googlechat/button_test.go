package googlechat_test

import (
	"math"
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestButton_Validate(t *testing.T) {
	v := validator.New()
	googlechat.RegisterValidation(v)
	click := googlechat.OnClick{Action: &googlechat.Action{Function: "save"}}
	for _, test := range []struct {
		name   string
		button googlechat.Button
		valid  bool
	}{
		{"text", googlechat.Button{Text: "Save", OnClick: click}, true},
		{"icon", googlechat.Button{Icon: &googlechat.Icon{KnownIcon: "EMAIL"}, OnClick: click}, true},
		{"black filled", googlechat.Button{Text: "Save", Color: &googlechat.Color{}, OnClick: click}, true},
		{"missing text and icon", googlechat.Button{OnClick: click}, false},
		{"missing action", googlechat.Button{Text: "Save"}, false},
		{"color and outline", googlechat.Button{Text: "Save", Color: &googlechat.Color{}, Type: googlechat.ButtonTypeOutlined, OnClick: click}, false},
		{"unknown type", googlechat.Button{Text: "Save", Type: "UNKNOWN", OnClick: click}, false},
		{"negative color", googlechat.Button{Text: "Save", Color: &googlechat.Color{Red: -0.1}, OnClick: click}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := v.Struct(test.button)
			if test.valid {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, new(validator.ValidationErrors))
			}
		})
	}
}

func TestColor_Validate(t *testing.T) {
	v := validator.New()
	require.NoError(t, v.Struct(googlechat.Color{Red: 0, Green: 0.5, Blue: 1}))
	for _, invalid := range []float64{-0.1, 1.1, math.NaN(), math.Inf(1), math.Inf(-1)} {
		require.Error(t, v.Struct(googlechat.Color{Red: invalid}))
	}
}
