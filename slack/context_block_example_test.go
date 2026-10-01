package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleContextBlock_MarshalJSON() {
	block := slack.ContextBlock{Elements: []slack.ContextElement{
		slack.ImageElement{ImageURL: "https://example.com/avatar.png", AltText: "Build bot"},
		slack.PlainTextObject{Text: "Updated just now"},
	}}
	jsonBytes, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "context",
	//   "elements": [
	//     {
	//       "type": "image",
	//       "image_url": "https://example.com/avatar.png",
	//       "alt_text": "Build bot"
	//     },
	//     {
	//       "type": "plain_text",
	//       "text": "Updated just now"
	//     }
	//   ]
	// }
}
