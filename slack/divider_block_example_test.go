package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleDividerBlock_MarshalJSON() {
	jsonBytes, err := json.MarshalIndent(slack.DividerBlock{}, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "divider"
	// }
}
