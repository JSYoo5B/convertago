package googlechat_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleToMessage_decoratedText() {
	type icon struct {
		Name string `googlechat:"part;slot=knownIcon"`
	}
	type rowInput struct {
		Text  string `googlechat:"part;style=bold"`
		Icon  icon   `googlechat:"icon;slot=startIcon"`
		Label string `googlechat:"part;slot=topLabel"`
	}
	row := rowInput{"Ready", icon{"STAR"}, "Status"}
	message, err := googlechat.ToMessage(struct {
		Row rowInput `googlechat:"decoratedText"`
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
	//                 "decoratedText": {
	//                   "startIcon": {
	//                     "knownIcon": "STAR"
	//                   },
	//                   "topLabel": "Status",
	//                   "text": "\u003cb\u003eReady\u003c/b\u003e"
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
