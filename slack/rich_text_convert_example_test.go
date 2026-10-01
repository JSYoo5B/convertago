package slack_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/slack"
)

func ExampleToMessage_richTextList() {
	type list struct {
		Style string   `slack:"part;slot=style"`
		Items []string `slack:"rich_text_section"`
	}
	source := struct {
		Body struct {
			List list `slack:"rich_text_list"`
		} `slack:"rich_text"`
	}{}
	source.Body.List = list{"ordered", []string{"First", "Second"}}
	message, err := slack.ToMessage(source)
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output:
	// {
	//   "blocks": [
	//     {
	//       "type": "rich_text",
	//       "elements": [
	//         {
	//           "type": "rich_text_list",
	//           "style": "ordered",
	//           "elements": [
	//             {
	//               "type": "rich_text_section",
	//               "elements": [
	//                 {
	//                   "type": "text",
	//                   "text": "First"
	//                 }
	//               ]
	//             },
	//             {
	//               "type": "rich_text_section",
	//               "elements": [
	//                 {
	//                   "type": "text",
	//                   "text": "Second"
	//                 }
	//               ]
	//             }
	//           ]
	//         }
	//       ]
	//     }
	//   ]
	// }
}

func ExampleToMessage_richTextMention() {
	source := struct {
		Body struct {
			Text  string `slack:"part"`
			User  string `slack:"user;style=bold"`
			Emoji string `slack:"emoji"`
		} `slack:"rich_text"`
	}{}
	source.Body.Text = "Hello "
	source.Body.User = "U123"
	source.Body.Emoji = "wave"
	message, err := slack.ToMessage(source)
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output:
	// {
	//   "blocks": [
	//     {
	//       "type": "rich_text",
	//       "elements": [
	//         {
	//           "type": "rich_text_section",
	//           "elements": [
	//             {
	//               "type": "text",
	//               "text": "Hello "
	//             },
	//             {
	//               "type": "user",
	//               "user_id": "U123",
	//               "style": {
	//                 "bold": true
	//               }
	//             },
	//             {
	//               "type": "emoji",
	//               "name": "wave"
	//             }
	//           ]
	//         }
	//       ]
	//     }
	//   ]
	// }
}
