package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleMessage() {
	message := slack.Message{
		Text: "Build finished. View the report.",
		Blocks: []slack.Block{
			slack.HeaderBlock{Text: slack.PlainTextObject{Text: "Build finished"}},
			slack.SectionBlock{
				Text: slack.MrkdwnTextObject{Text: "The build *passed*."},
				Accessory: slack.ButtonElement{
					Text: slack.PlainTextObject{Text: "View report"},
					URL:  "https://example.com/build/42",
				},
			},
		},
	}

	jsonBytes, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "text": "Build finished. View the report.",
	//   "blocks": [
	//     {
	//       "type": "header",
	//       "text": {
	//         "type": "plain_text",
	//         "text": "Build finished"
	//       }
	//     },
	//     {
	//       "type": "section",
	//       "text": {
	//         "type": "mrkdwn",
	//         "text": "The build *passed*."
	//       },
	//       "accessory": {
	//         "type": "button",
	//         "text": {
	//           "type": "plain_text",
	//           "text": "View report"
	//         },
	//         "url": "https://example.com/build/42"
	//       }
	//     }
	//   ]
	// }
}
