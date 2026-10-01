package googlechat_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleToMessage_cards() {
	type card struct {
		Title string `googlechat:"header"`
		Body  string `googlechat:"textParagraph"`
	}
	type wrapped struct {
		ID   string `googlechat:"part;slot=cardId"`
		Card card   `googlechat:"card"`
	}
	source := struct {
		Cards []wrapped `googlechat:"cardWithId"`
	}{[]wrapped{{"first", card{"First", "one"}}, {"second", card{"Second", "two"}}}}
	message, err := googlechat.ToMessage(source)
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
	//   "cardsV2": [
	//     {
	//       "cardId": "first",
	//       "card": {
	//         "header": {
	//           "title": "First"
	//         },
	//         "sections": [
	//           {
	//             "widgets": [
	//               {
	//                 "textParagraph": {
	//                   "text": "one"
	//                 }
	//               }
	//             ]
	//           }
	//         ]
	//       }
	//     },
	//     {
	//       "cardId": "second",
	//       "card": {
	//         "header": {
	//           "title": "Second"
	//         },
	//         "sections": [
	//           {
	//             "widgets": [
	//               {
	//                 "textParagraph": {
	//                   "text": "two"
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

func ExampleToMessage_sections() {
	type section struct {
		Header      string `googlechat:"part;slot=header"`
		Text        string `googlechat:"textParagraph"`
		Collapsible bool   `googlechat:"part;slot=collapsible"`
	}
	source := struct {
		Sections []section `googlechat:"section"`
	}{[]section{{"Summary", "Ready", false}, {"Details", "Completed", true}}}
	message, err := googlechat.ToMessage(source)
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
	//   "cardsV2": [
	//     {
	//       "card": {
	//         "sections": [
	//           {
	//             "header": "Summary",
	//             "widgets": [
	//               {
	//                 "textParagraph": {
	//                   "text": "Ready"
	//                 }
	//               }
	//             ]
	//           },
	//           {
	//             "header": "Details",
	//             "widgets": [
	//               {
	//                 "textParagraph": {
	//                   "text": "Completed"
	//                 }
	//               }
	//             ],
	//             "collapsible": true
	//           }
	//         ]
	//       }
	//     }
	//   ]
	// }
}
