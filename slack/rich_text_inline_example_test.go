package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleRichTextStyle() {
	value := slack.RichTextStyle{Bold: true, Italic: true, Underline: true}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "bold": true,
	//   "italic": true,
	//   "underline": true
	// }
}

func ExampleTextInline_MarshalJSON() {
	value := slack.TextInline{
		Text:  "Build passed",
		Style: &slack.RichTextStyle{Bold: true, Italic: true},
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "type": "text",
	//   "text": "Build passed",
	//   "style": {
	//     "bold": true,
	//     "italic": true
	//   }
	// }
}

func ExampleLinkInline_MarshalJSON() {
	value := slack.LinkInline{
		URL:   "https://example.com/report",
		Text:  "Build report",
		Style: &slack.RichTextStyle{Underline: true},
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "type": "link",
	//   "url": "https://example.com/report",
	//   "text": "Build report",
	//   "style": {
	//     "underline": true
	//   }
	// }
}

func ExampleUserInline_MarshalJSON() {
	value := slack.UserInline{
		UserID: "U0123456",
		Style:  &slack.RichTextStyle{Bold: true},
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "type": "user",
	//   "user_id": "U0123456",
	//   "style": {
	//     "bold": true
	//   }
	// }
}

func ExampleEmojiInline_MarshalJSON() {
	value := slack.EmojiInline{Name: "wave", Unicode: "1f44b"}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "type": "emoji",
	//   "name": "wave",
	//   "unicode": "1f44b"
	// }
}
