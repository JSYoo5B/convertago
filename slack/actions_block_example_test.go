package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleActionsBlock_MarshalJSON() {
	block := slack.ActionsBlock{Elements: []slack.ActionElement{
		slack.ButtonElement{Text: slack.PlainTextObject{Text: "Approve"}, ActionID: "approve", Style: slack.ButtonStylePrimary},
	}}
	jsonBytes, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "actions",
	//   "elements": [
	//     {
	//       "type": "button",
	//       "text": {
	//         "type": "plain_text",
	//         "text": "Approve"
	//       },
	//       "action_id": "approve",
	//       "style": "primary"
	//     }
	//   ]
	// }
}
