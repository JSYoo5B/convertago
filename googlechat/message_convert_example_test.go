package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleToMessage() {
	source := struct {
		Text string `googlechat:"textParagraph"`
	}{"Google Chat notification"}
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
	//             "widgets": [
	//               {
	//                 "textParagraph": {
	//                   "text": "Google Chat notification"
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
