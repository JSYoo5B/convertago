package googlechat_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestOnClick_Validate(t *testing.T) {
	v := validator.New()
	for _, test := range []struct {
		name  string
		click googlechat.OnClick
		valid bool
	}{
		{"action", googlechat.OnClick{Action: &googlechat.Action{Function: "save"}}, true},
		{"open link", googlechat.OnClick{OpenLink: &googlechat.OpenLink{URL: "https://example.com"}}, true},
		{"empty", googlechat.OnClick{}, false},
		{"two actions", googlechat.OnClick{Action: &googlechat.Action{Function: "save"}, OpenLink: &googlechat.OpenLink{URL: "https://example.com"}}, false},
		{"empty action", googlechat.OnClick{Action: &googlechat.Action{}}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := v.Struct(test.click)
			if test.valid {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, new(validator.ValidationErrors))
			}
		})
	}
}

func TestAction_Validate(t *testing.T) {
	v := validator.New()
	googlechat.RegisterValidation(v)
	action := googlechat.Action{Function: "save", Parameters: []googlechat.ActionParameter{
		{Key: "id", Value: "123"}, {Key: "view", Value: "edit"},
	}, AllWidgetsAreRequired: true, RequiredWidgets: []string{}}
	require.NoError(t, v.Struct(action))
	action.Parameters[1].Key = "id"
	require.Error(t, v.Struct(action))
	action.Parameters = nil
	action.RequiredWidgets = []string{"email"}
	require.Error(t, v.Struct(action))
	action.AllWidgetsAreRequired = false
	require.NoError(t, v.Struct(action))
	action.RequiredWidgets = []string{""}
	require.Error(t, v.Struct(action))
	action.RequiredWidgets = nil
	action.LoadIndicator = "BUSY"
	require.Error(t, v.Struct(action))
}

func TestOpenLink_Validate(t *testing.T) {
	v := validator.New()
	googlechat.RegisterValidation(v)
	require.NoError(t, v.Struct(googlechat.OpenLink{URL: "mailto:jane@example.com"}))
	require.Error(t, v.Struct(googlechat.OpenLink{URL: "https:#fragment"}))
	require.Error(t, v.Struct(googlechat.OpenLink{URL: "/relative"}))
}

func TestOverflowMenu_Validate(t *testing.T) {
	v := validator.New()
	googlechat.RegisterValidation(v)
	click := googlechat.OnClick{OpenLink: &googlechat.OpenLink{URL: "https://example.com"}}
	menu := googlechat.OverflowMenu{Items: []googlechat.OverflowMenuItem{{Text: "Open", OnClick: click}}}
	require.NoError(t, v.Struct(menu))
	nested := googlechat.OverflowMenu{Items: []googlechat.OverflowMenuItem{{
		Text: "More", OnClick: googlechat.OnClick{OverflowMenu: &menu},
	}}}
	require.Error(t, v.Struct(nested))
	require.Error(t, v.Struct(googlechat.OverflowMenu{}))
}
