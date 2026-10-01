package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExamplePlainTextObject_MarshalJSON() {
	emoji := false
	text := slack.PlainTextObject{Text: "Build finished", Emoji: &emoji}
	jsonBytes, err := json.MarshalIndent(text, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "plain_text",
	//   "text": "Build finished",
	//   "emoji": false
	// }
}

func ExampleMrkdwnTextObject_MarshalJSON() {
	text := slack.MrkdwnTextObject{Text: "Build *finished*", Verbatim: true}
	jsonBytes, err := json.MarshalIndent(text, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "mrkdwn",
	//   "text": "Build *finished*",
	//   "verbatim": true
	// }
}
