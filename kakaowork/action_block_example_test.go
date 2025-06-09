package kakaowork_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleActionBlock_MarshalJSON() {
	actionBlock := kakaowork.ActionBlock{
		Buttons: []kakaowork.ButtonBlock{
			{
				Text:  "승인",
				Style: kakaowork.ButtonStylePrimary,
				Action: kakaowork.SubmitAction{
					Name:  "approve",
					Value: "true",
				},
			},
			{
				Text:  "반려",
				Style: kakaowork.ButtonStyleDefault,
				Action: kakaowork.SubmitAction{
					Name:  "approve",
					Value: "false",
				},
			},
		},
	}

	jsonBytes, err := json.MarshalIndent(actionBlock, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "action",
	//   "elements": [
	//     {
	//       "type": "button",
	//       "text": "승인",
	//       "style": "primary",
	//       "action": {
	//         "type": "submit_action",
	//         "name": "approve",
	//         "value": "true"
	//       }
	//     },
	//     {
	//       "type": "button",
	//       "text": "반려",
	//       "style": "default",
	//       "action": {
	//         "type": "submit_action",
	//         "name": "approve",
	//         "value": "false"
	//       }
	//     }
	//   ]
	// }
}
