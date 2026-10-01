package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleImageElement_MarshalJSON() {
	value := slack.ImageElement{ImageURL: "https://example.com/chart.png", AltText: "Build chart"}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "type": "image",
	//   "image_url": "https://example.com/chart.png",
	//   "alt_text": "Build chart"
	// }
}
