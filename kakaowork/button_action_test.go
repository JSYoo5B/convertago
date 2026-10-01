package kakaowork_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

var (
	_ kakaowork.ButtonAction = kakaowork.OpenSystemBrowserAction{}
	_ kakaowork.ButtonAction = (*kakaowork.OpenSystemBrowserAction)(nil)
	_ kakaowork.ButtonAction = kakaowork.OpenInAppBrowserAction{}
	_ kakaowork.ButtonAction = (*kakaowork.OpenInAppBrowserAction)(nil)
	_ kakaowork.ButtonAction = kakaowork.OpenExternalAppAction{}
	_ kakaowork.ButtonAction = (*kakaowork.OpenExternalAppAction)(nil)
	_ kakaowork.ButtonAction = kakaowork.SubmitAction{}
	_ kakaowork.ButtonAction = (*kakaowork.SubmitAction)(nil)
	_ kakaowork.ButtonAction = kakaowork.CallModalAction{}
	_ kakaowork.ButtonAction = (*kakaowork.CallModalAction)(nil)
	_ kakaowork.ButtonAction = kakaowork.ExclusiveAction{}
	_ kakaowork.ButtonAction = (*kakaowork.ExclusiveAction)(nil)
)

type externalButtonAction struct {
	kakaowork.BubbleBlock
}

func (externalButtonAction) ActionType() string { return "submit_action" }
func (externalButtonAction) buttonAction()      {}

func TestButtonActionRejectsExternalImplementation(t *testing.T) {
	value := externalButtonAction{BubbleBlock: kakaowork.TextBlock{Text: "custom"}}
	for _, candidate := range []any{value, &value} {
		if _, ok := candidate.(kakaowork.ButtonAction); ok {
			t.Errorf("external type %T implements ButtonAction", candidate)
		}
	}
}

func TestButtonAction_Validate(t *testing.T) {
	v := validator.New()
	kakaowork.RegisterValidation(v)
	var absent *kakaowork.SubmitAction
	for _, test := range []struct {
		name   string
		action any
		valid  bool
	}{
		{"browser", kakaowork.OpenSystemBrowserAction{Value: "https://example.com"}, true},
		{"browser scheme", kakaowork.OpenSystemBrowserAction{Value: "ftp://example.com"}, false},
		{"browser relative", kakaowork.OpenSystemBrowserAction{Value: "/page"}, false},
		{"in-app default dimensions", kakaowork.OpenInAppBrowserAction{Value: "https://example.com"}, true},
		{"standalone dimensions", kakaowork.OpenInAppBrowserAction{Value: "https://example.com", Standalone: true, Width: 980, Height: 720}, true},
		{"dimensions without standalone", kakaowork.OpenInAppBrowserAction{Value: "https://example.com", Width: 980}, false},
		{"negative dimension", kakaowork.OpenInAppBrowserAction{Value: "https://example.com", Standalone: true, Height: -1}, false},
		{"app scheme", kakaowork.OpenExternalAppAction{Value: "geo:37.537229,127.005515"}, true},
		{"platform app schemes", kakaowork.OpenExternalAppAction{Value: "ios=kakaomap%3A%2F%2Flook&aos=geo%3A37.5%2C127"}, true},
		{"relative app destination", kakaowork.OpenExternalAppAction{Value: "ios=relative"}, false},
		{"duplicate platform", kakaowork.OpenExternalAppAction{Value: "ios=geo%3A1&ios=geo%3A2"}, false},
		{"submit name", kakaowork.SubmitAction{Name: "accept"}, true},
		{"missing submit name", kakaowork.SubmitAction{Value: "event=1"}, false},
		{"exclusive default", kakaowork.ExclusiveAction{Default: &kakaowork.SubmitAction{Name: "accept"}}, true},
		{"missing exclusive default", kakaowork.ExclusiveAction{}, false},
		{"typed nil default", kakaowork.ExclusiveAction{Default: absent}, false},
		{"nested exclusive", kakaowork.ExclusiveAction{Default: &kakaowork.ExclusiveAction{Default: kakaowork.SubmitAction{Name: "accept"}}}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := v.Struct(test.action)
			if test.valid {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, new(validator.ValidationErrors))
			}
		})
	}
}
