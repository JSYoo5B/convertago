package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleMarkdownBlock_MarshalJSON() {
	block := slack.MarkdownBlock{Text: "## Build report\nAll checks **passed**."}
	jsonBytes, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "markdown",
	//   "text": "## Build report\nAll checks **passed**."
	// }
}
