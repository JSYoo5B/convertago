package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleDecoratedText() {
	widget := googlechat.Widget{Content: googlechat.DecoratedText{
		StartIcon: &googlechat.Icon{KnownIcon: "EMAIL"},
		TopLabel:  "Report recipient",
		Text:      "reviewer@example.com",
		WrapText:  true,
		SwitchControl: &googlechat.SwitchControl{
			Name:        "send_email",
			Value:       "yes",
			Selected:    true,
			ControlType: googlechat.SwitchControlTypeCheckBox,
		},
	}}
	jsonBytes, err := json.MarshalIndent(widget, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "decoratedText": {
	//     "startIcon": {
	//       "knownIcon": "EMAIL"
	//     },
	//     "topLabel": "Report recipient",
	//     "text": "reviewer@example.com",
	//     "wrapText": true,
	//     "switchControl": {
	//       "name": "send_email",
	//       "value": "yes",
	//       "selected": true,
	//       "controlType": "CHECK_BOX"
	//     }
	//   }
	// }
}

func ExampleSwitchControl() {
	value := googlechat.SwitchControl{
		Name:           "send_email",
		Value:          "yes",
		Selected:       true,
		OnChangeAction: &googlechat.Action{Function: "updateEmailPreference"},
		ControlType:    googlechat.SwitchControlTypeCheckBox,
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "name": "send_email",
	//   "value": "yes",
	//   "selected": true,
	//   "onChangeAction": {
	//     "function": "updateEmailPreference"
	//   },
	//   "controlType": "CHECK_BOX"
	// }
}
