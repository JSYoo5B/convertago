package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleHeaderBlock_MarshalJSON() {
	block := slack.HeaderBlock{Text: slack.PlainTextObject{Text: "Build report"}, BlockID: "header-42", Level: 2}
	jsonBytes, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "header",
	//   "text": {
	//     "type": "plain_text",
	//     "text": "Build report"
	//   },
	//   "block_id": "header-42",
	//   "level": 2
	// }
}
