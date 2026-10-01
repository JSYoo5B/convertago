package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleButtonList() {
	widget := googlechat.Widget{Content: googlechat.ButtonList{Buttons: []googlechat.Button{{
		Text: "Approve",
		OnClick: googlechat.OnClick{Action: &googlechat.Action{
			Function:      "approveBuild",
			Parameters:    []googlechat.ActionParameter{{Key: "build", Value: "42"}},
			LoadIndicator: googlechat.LoadIndicatorNone,
		}},
		Type: googlechat.ButtonTypeFilled,
	}}}}
	jsonBytes, err := json.MarshalIndent(widget, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "buttonList": {
	//     "buttons": [
	//       {
	//         "text": "Approve",
	//         "onClick": {
	//           "action": {
	//             "function": "approveBuild",
	//             "parameters": [
	//               {
	//                 "key": "build",
	//                 "value": "42"
	//               }
	//             ],
	//             "loadIndicator": "NONE"
	//           }
	//         },
	//         "type": "FILLED"
	//       }
	//     ]
	//   }
	// }
}

func ExampleButton() {
	value := googlechat.Button{
		Text:    "View report",
		Color:   &googlechat.Color{Blue: 1},
		OnClick: googlechat.OnClick{OpenLink: &googlechat.OpenLink{URL: "https://example.com/report"}},
		AltText: "Open the build report",
		Type:    googlechat.ButtonTypeFilled,
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "text": "View report",
	//   "color": {
	//     "red": 0,
	//     "green": 0,
	//     "blue": 1
	//   },
	//   "onClick": {
	//     "openLink": {
	//       "url": "https://example.com/report"
	//     }
	//   },
	//   "altText": "Open the build report",
	//   "type": "FILLED"
	// }
}

func ExampleColor() {
	value := googlechat.Color{Red: 0.25, Green: 0.5, Blue: 1}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "red": 0.25,
	//   "green": 0.5,
	//   "blue": 1
	// }
}
