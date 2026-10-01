package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleSectionBlock_MarshalJSON() {
	block := slack.SectionBlock{
		Fields: []slack.TextObject{
			slack.MrkdwnTextObject{Text: "*Status*\nPassed"},
			slack.PlainTextObject{Text: "Duration: 12 seconds"},
		},
	}
	jsonBytes, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "section",
	//   "fields": [
	//     {
	//       "type": "mrkdwn",
	//       "text": "*Status*\nPassed"
	//     },
	//     {
	//       "type": "plain_text",
	//       "text": "Duration: 12 seconds"
	//     }
	//   ]
	// }
}
