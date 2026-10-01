package googlechat_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleToMessage_buttons() {
	type button struct {
		Text string `googlechat:"part"`
		Link string `googlechat:"openLink"`
	}
	type rowInput struct {
		Buttons []button `googlechat:"button"`
	}
	row := rowInput{[]button{{"Open", "https://example.com"}}}
	message, err := googlechat.ToMessage(struct {
		Row rowInput `googlechat:"buttonList"`
	}{row})
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
	//             "widgets": [
	//               {
	//                 "buttonList": {
	//                   "buttons": [
	//                     {
	//                       "text": "Open",
	//                       "onClick": {
	//                         "openLink": {
	//                           "url": "https://example.com"
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
