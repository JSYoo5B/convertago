package googlechat_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestDecoratedText_Validate(t *testing.T) {
	v := validator.New()
	control := googlechat.SwitchControl{Name: "enabled", ControlType: googlechat.SwitchControlTypeCheckBox}
	text := googlechat.DecoratedText{Text: "Enabled", SwitchControl: &control}
	require.NoError(t, v.Struct(text))
	text.EndIcon = &googlechat.Icon{KnownIcon: "EMAIL"}
	require.Error(t, v.Struct(text))
	text.EndIcon = nil
	control.ControlType = "RADIO"
	require.Error(t, v.Struct(text))
	control.ControlType = ""
	control.Name = ""
	require.Error(t, v.Struct(text))
}
