package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleMessage() {
	message := googlechat.Message{
		Text: "Build finished",
		CardsV2: []googlechat.CardWithID{{
			CardID: "build-42",
			Card: googlechat.Card{
				Header: &googlechat.CardHeader{Title: "Build report", Subtitle: "Build 42"},
				Sections: []googlechat.Section{{
					Widgets: []googlechat.Widget{
						{Content: googlechat.TextParagraph{Text: "All checks passed."}},
						{Content: googlechat.ButtonList{Buttons: []googlechat.Button{{
							Text:    "View report",
							OnClick: googlechat.OnClick{OpenLink: &googlechat.OpenLink{URL: "https://example.com/build/42"}},
						}}}},
					},
				}},
			},
		}},
	}
	jsonBytes, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "text": "Build finished",
	//   "cardsV2": [
	//     {
	//       "cardId": "build-42",
	//       "card": {
	//         "header": {
	//           "title": "Build report",
	//           "subtitle": "Build 42"
	//         },
	//         "sections": [
	//           {
	//             "widgets": [
	//               {
	//                 "textParagraph": {
	//                   "text": "All checks passed."
	//                 }
	//               },
	//               {
	//                 "buttonList": {
	//                   "buttons": [
	//                     {
	//                       "text": "View report",
	//                       "onClick": {
	//                         "openLink": {
	//                           "url": "https://example.com/build/42"
	//                         }
	//                       }
	//                     }
	//                   ]
	//                 }
	//               }
	//             ]
	//           }
	//         ]
	//       }
	//     }
	//   ]
	// }
}

func ExampleCardWithID() {
	value := googlechat.CardWithID{
		CardID: "build-42",
		Card:   googlechat.Card{Header: &googlechat.CardHeader{Title: "Build report"}},
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "cardId": "build-42",
	//   "card": {
	//     "header": {
	//       "title": "Build report"
	//     }
	//   }
	// }
}
