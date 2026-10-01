package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleChip() {
	value := googlechat.Chip{
		Icon:    &googlechat.Icon{KnownIcon: "EMAIL"},
		Label:   "Send report",
		OnClick: &googlechat.OnClick{Action: &googlechat.Action{Function: "sendReport"}},
		AltText: "Email the build report",
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "icon": {
	//     "knownIcon": "EMAIL"
	//   },
	//   "label": "Send report",
	//   "onClick": {
	//     "action": {
	//       "function": "sendReport"
	//     }
	//   },
	//   "altText": "Email the build report"
	// }
}

func ExampleChipList() {
	widget := googlechat.Widget{Content: googlechat.ChipList{
		Layout: googlechat.ChipListLayoutWrapped,
		Chips: []googlechat.Chip{
			{Label: "Passed", Icon: &googlechat.Icon{KnownIcon: "STAR"}},
			{Label: "Report", OnClick: &googlechat.OnClick{OpenLink: &googlechat.OpenLink{URL: "https://example.com/report"}}},
		},
	}}
	jsonBytes, err := json.MarshalIndent(widget, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "chipList": {
	//     "layout": "WRAPPED",
	//     "chips": [
	//       {
	//         "icon": {
	//           "knownIcon": "STAR"
	//         },
	//         "label": "Passed"
	//       },
	//       {
	//         "label": "Report",
	//         "onClick": {
	//           "openLink": {
	//             "url": "https://example.com/report"
	//           }
	//         }
	//       }
	//     ]
	//   }
	// }
}
