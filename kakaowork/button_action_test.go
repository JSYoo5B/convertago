package kakaowork_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/kakaowork"
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
