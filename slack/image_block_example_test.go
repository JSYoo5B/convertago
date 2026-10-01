package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleImageBlock_MarshalJSON() {
	block := slack.ImageBlock{
		SlackFile: &slack.SlackFileObject{ID: "F0123456"},
		AltText:   "Build duration chart",
		Title:     &slack.PlainTextObject{Text: "Build duration"},
	}
	jsonBytes, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "image",
	//   "slack_file": {
	//     "id": "F0123456"
	//   },
	//   "alt_text": "Build duration chart",
	//   "title": {
	//     "type": "plain_text",
	//     "text": "Build duration"
	//   }
	// }
}

func ExampleSlackFileObject() {
	value := slack.SlackFileObject{ID: "F0123456"}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "id": "F0123456"
	// }
}
