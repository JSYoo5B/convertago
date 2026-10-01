package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleButtonElement_MarshalJSON() {
	button := slack.ButtonElement{
		Text:     slack.PlainTextObject{Text: "Delete"},
		ActionID: "delete",
		Style:    slack.ButtonStyleDanger,
		Confirm: &slack.ConfirmationDialogObject{
			Title:   slack.PlainTextObject{Text: "Delete report?"},
			Text:    slack.PlainTextObject{Text: "This action removes the report."},
			Confirm: slack.PlainTextObject{Text: "Delete"},
			Deny:    slack.PlainTextObject{Text: "Cancel"},
		},
	}
	jsonBytes, err := json.MarshalIndent(button, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "button",
	//   "text": {
	//     "type": "plain_text",
	//     "text": "Delete"
	//   },
	//   "action_id": "delete",
	//   "style": "danger",
	//   "confirm": {
	//     "title": {
	//       "type": "plain_text",
	//       "text": "Delete report?"
	//     },
	//     "text": {
	//       "type": "plain_text",
	//       "text": "This action removes the report."
	//     },
	//     "confirm": {
	//       "type": "plain_text",
	//       "text": "Delete"
	//     },
	//     "deny": {
	//       "type": "plain_text",
	//       "text": "Cancel"
	//     }
	//   }
	// }
}

func ExampleConfirmationDialogObject() {
	value := slack.ConfirmationDialogObject{
		Title:   slack.PlainTextObject{Text: "Delete report?"},
		Text:    slack.PlainTextObject{Text: "The report will be removed."},
		Confirm: slack.PlainTextObject{Text: "Delete"},
		Deny:    slack.PlainTextObject{Text: "Cancel"},
		Style:   slack.ButtonStyleDanger,
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "title": {
	//     "type": "plain_text",
	//     "text": "Delete report?"
	//   },
	//   "text": {
	//     "type": "plain_text",
	//     "text": "The report will be removed."
	//   },
	//   "confirm": {
	//     "type": "plain_text",
	//     "text": "Delete"
	//   },
	//   "deny": {
	//     "type": "plain_text",
	//     "text": "Cancel"
	//   },
	//   "style": "danger"
	// }
}
