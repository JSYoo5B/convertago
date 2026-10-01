package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleVideoBlock_MarshalJSON() {
	block := slack.VideoBlock{
		Title:        slack.PlainTextObject{Text: "Build walkthrough"},
		VideoURL:     "https://example.com/embed/walkthrough",
		AltText:      "Build walkthrough video",
		ThumbnailURL: "https://example.com/walkthrough.png",
		TitleURL:     "https://example.com/watch/walkthrough",
		Description:  &slack.PlainTextObject{Text: "Review the build results."},
	}
	jsonBytes, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "video",
	//   "title": {
	//     "type": "plain_text",
	//     "text": "Build walkthrough"
	//   },
	//   "video_url": "https://example.com/embed/walkthrough",
	//   "alt_text": "Build walkthrough video",
	//   "thumbnail_url": "https://example.com/walkthrough.png",
	//   "title_url": "https://example.com/watch/walkthrough",
	//   "description": {
	//     "type": "plain_text",
	//     "text": "Review the build results."
	//   }
	// }
}
