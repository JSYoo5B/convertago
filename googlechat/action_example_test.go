package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleOnClick_MarshalJSON() {
	value := googlechat.OnClick{
		OpenLink: &googlechat.OpenLink{URL: "https://example.com/report"},
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "openLink": {
	//     "url": "https://example.com/report"
	//   }
	// }
}

func ExampleOpenLink() {
	value := googlechat.OpenLink{URL: "https://example.com/report"}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "url": "https://example.com/report"
	// }
}

func ExampleAction() {
	value := googlechat.Action{
		Function:        "approveBuild",
		Parameters:      []googlechat.ActionParameter{{Key: "build", Value: "42"}},
		LoadIndicator:   googlechat.LoadIndicatorNone,
		RequiredWidgets: []string{"approval"},
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "function": "approveBuild",
	//   "parameters": [
	//     {
	//       "key": "build",
	//       "value": "42"
	//     }
	//   ],
	//   "loadIndicator": "NONE",
	//   "requiredWidgets": [
	//     "approval"
	//   ]
	// }
}

func ExampleActionParameter() {
	value := googlechat.ActionParameter{Key: "build", Value: "42"}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "key": "build",
	//   "value": "42"
	// }
}

func ExampleOverflowMenu() {
	value := googlechat.OverflowMenu{Items: []googlechat.OverflowMenuItem{
		{
			Text:    "View report",
			OnClick: googlechat.OnClick{OpenLink: &googlechat.OpenLink{URL: "https://example.com/report"}},
		},
		{
			Text:    "Send report",
			OnClick: googlechat.OnClick{Action: &googlechat.Action{Function: "sendReport"}},
		},
	}}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "items": [
	//     {
	//       "text": "View report",
	//       "onClick": {
	//         "openLink": {
	//           "url": "https://example.com/report"
	//         }
	//       }
	//     },
	//     {
	//       "text": "Send report",
	//       "onClick": {
	//         "action": {
	//           "function": "sendReport"
	//         }
	//       }
	//     }
	//   ]
	// }
}

func ExampleOverflowMenuItem() {
	value := googlechat.OverflowMenuItem{
		StartIcon: &googlechat.Icon{KnownIcon: "EMAIL"},
		Text:      "Send report",
		OnClick:   googlechat.OnClick{Action: &googlechat.Action{Function: "sendReport"}},
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "startIcon": {
	//     "knownIcon": "EMAIL"
	//   },
	//   "text": "Send report",
	//   "onClick": {
	//     "action": {
	//       "function": "sendReport"
	//     }
	//   }
	// }
}
