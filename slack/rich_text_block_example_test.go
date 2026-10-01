package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleRichTextBlock_MarshalJSON() {
	block := slack.RichTextBlock{Elements: []slack.RichTextElement{
		slack.RichTextSection{Elements: []slack.RichTextInline{
			slack.TextInline{Text: "Build ", Style: &slack.RichTextStyle{Bold: true}},
			slack.LinkInline{URL: "https://example.com/build/42", Text: "report"},
			slack.TextInline{Text: " for "},
			slack.UserInline{UserID: "U0123456"},
			slack.EmojiInline{Name: "white_check_mark"},
		}},
		slack.RichTextList{
			Style: slack.RichTextListStyleBullet,
			Elements: []slack.RichTextSection{
				{Elements: []slack.RichTextInline{slack.TextInline{Text: "Tests passed"}}},
				{Elements: []slack.RichTextInline{slack.TextInline{Text: "Artifacts published"}}},
			},
		},
	}}
	jsonBytes, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "rich_text",
	//   "elements": [
	//     {
	//       "type": "rich_text_section",
	//       "elements": [
	//         {
	//           "type": "text",
	//           "text": "Build ",
	//           "style": {
	//             "bold": true
	//           }
	//         },
	//         {
	//           "type": "link",
	//           "url": "https://example.com/build/42",
	//           "text": "report"
	//         },
	//         {
	//           "type": "text",
	//           "text": " for "
	//         },
	//         {
	//           "type": "user",
	//           "user_id": "U0123456"
	//         },
	//         {
	//           "type": "emoji",
	//           "name": "white_check_mark"
	//         }
	//       ]
	//     },
	//     {
	//       "type": "rich_text_list",
	//       "style": "bullet",
	//       "elements": [
	//         {
	//           "type": "rich_text_section",
	//           "elements": [
	//             {
	//               "type": "text",
	//               "text": "Tests passed"
	//             }
	//           ]
	//         },
	//         {
	//           "type": "rich_text_section",
	//           "elements": [
	//             {
	//               "type": "text",
	//               "text": "Artifacts published"
	//             }
	//           ]
	//         }
	//       ]
	//     }
	//   ]
	// }
}

func ExampleRichTextPreformatted_MarshalJSON() {
	border := 0
	region := slack.RichTextPreformatted{
		Elements: []slack.PreformattedInline{slack.TextInline{Text: "go test ./..."}},
		Border:   &border,
		Language: "shell",
	}
	jsonBytes, err := json.MarshalIndent(region, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "rich_text_preformatted",
	//   "elements": [
	//     {
	//       "type": "text",
	//       "text": "go test ./..."
	//     }
	//   ],
	//   "border": 0,
	//   "language": "shell"
	// }
}

func ExampleRichTextQuote_MarshalJSON() {
	quote := slack.RichTextQuote{Elements: []slack.RichTextInline{slack.TextInline{Text: "All checks passed."}}}
	jsonBytes, err := json.MarshalIndent(quote, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "rich_text_quote",
	//   "elements": [
	//     {
	//       "type": "text",
	//       "text": "All checks passed."
	//     }
	//   ]
	// }
}

func ExampleRichTextSection_MarshalJSON() {
	value := slack.RichTextSection{Elements: []slack.RichTextInline{
		slack.TextInline{Text: "Hello "},
		slack.UserInline{UserID: "U0123456"},
	}}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "type": "rich_text_section",
	//   "elements": [
	//     {
	//       "type": "text",
	//       "text": "Hello "
	//     },
	//     {
	//       "type": "user",
	//       "user_id": "U0123456"
	//     }
	//   ]
	// }
}

func ExampleRichTextList_MarshalJSON() {
	border := 0
	value := slack.RichTextList{
		Style: slack.RichTextListStyleOrdered,
		Elements: []slack.RichTextSection{
			{Elements: []slack.RichTextInline{slack.TextInline{Text: "Run tests"}}},
			{Elements: []slack.RichTextInline{slack.TextInline{Text: "Review report"}}},
		},
		Indent: 1,
		Offset: 4,
		Border: &border,
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "type": "rich_text_list",
	//   "style": "ordered",
	//   "elements": [
	//     {
	//       "type": "rich_text_section",
	//       "elements": [
	//         {
	//           "type": "text",
	//           "text": "Run tests"
	//         }
	//       ]
	//     },
	//     {
	//       "type": "rich_text_section",
	//       "elements": [
	//         {
	//           "type": "text",
	//           "text": "Review report"
	//         }
	//       ]
	//     }
	//   ],
	//   "indent": 1,
	//   "offset": 4,
	//   "border": 0
	// }
}
