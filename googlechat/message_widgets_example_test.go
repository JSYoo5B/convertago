package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleToMessage_widgets() {
	source := struct {
		Text     string `googlechat:"text"`
		Fallback string `googlechat:"fallbackText"`
		Banner   struct {
			URL string `googlechat:"part"`
			Alt string `googlechat:"part;slot=alt"`
		} `googlechat:"image"`
		Divider struct{} `googlechat:"divider"`
		Body    string   `googlechat:"textParagraph"`
	}{Text: "Weekly report", Fallback: "Weekly report card", Body: "12 deploys this week"}
	source.Banner.URL, source.Banner.Alt = "https://example.com/banner.png", "Deploy chart"
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
	//   "text": "Weekly report",
	//   "cardsV2": [
	//     {
	//       "card": {
	//         "sections": [
	//           {
	//             "widgets": [
	//               {
	//                 "image": {
	//                   "imageUrl": "https://example.com/banner.png",
	//                   "altText": "Deploy chart"
	//                 }
	//               },
	//               {
	//                 "divider": {}
	//               },
	//               {
	//                 "textParagraph": {
	//                   "text": "12 deploys this week"
	//                 }
	//               }
	//             ]
	//           }
	//         ]
	//       }
	//     }
	//   ],
	//   "fallbackText": "Weekly report card"
	// }
}
