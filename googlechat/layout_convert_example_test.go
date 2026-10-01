package googlechat_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleToMessage_columns() {
	type column struct {
		Text string `googlechat:"textParagraph"`
	}
	type rowInput struct {
		Columns []column `googlechat:"column"`
	}
	row := rowInput{[]column{{"Left"}, {"Right"}}}
	message, err := googlechat.ToMessage(struct {
		Row rowInput `googlechat:"columns"`
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
	//                 "columns": {
	//                   "columnItems": [
	//                     {
	//                       "widgets": [
	//                         {
	//                           "textParagraph": {
	//                             "text": "Left"
	//                           }
	//                         }
	//                       ]
	//                     },
	//                     {
	//                       "widgets": [
	//                         {
	//                           "textParagraph": {
	//                             "text": "Right"
	//                           }
	//                         }
	//                       ]
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

func ExampleToMessage_grid() {
	type item struct {
		Title string `googlechat:"part;slot=title"`
	}
	type rowInput struct {
		Title string `googlechat:"part;slot=title"`
		Items []item `googlechat:"gridItem"`
	}
	row := rowInput{"Items", []item{{"One"}, {"Two"}}}
	message, err := googlechat.ToMessage(struct {
		Row rowInput `googlechat:"grid"`
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
	//                 "grid": {
	//                   "title": "Items",
	//                   "items": [
	//                     {
	//                       "title": "One"
	//                     },
	//                     {
	//                       "title": "Two"
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

func ExampleToMessage_carousel() {
	type card struct {
		Text   string `googlechat:"textParagraph"`
		Footer string `googlechat:"textParagraph;slot=footerWidgets"`
	}
	type rowInput struct {
		Cards []card `googlechat:"carouselCard"`
	}
	row := rowInput{[]card{{"First", "one"}, {"Second", "two"}}}
	message, err := googlechat.ToMessage(struct {
		Row rowInput `googlechat:"carousel"`
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
	//                 "carousel": {
	//                   "carouselCards": [
	//                     {
	//                       "widgets": [
	//                         {
	//                           "textParagraph": {
	//                             "text": "First"
	//                           }
	//                         }
	//                       ],
	//                       "footerWidgets": [
	//                         {
	//                           "textParagraph": {
	//                             "text": "one"
	//                           }
	//                         }
	//                       ]
	//                     },
	//                     {
	//                       "widgets": [
	//                         {
	//                           "textParagraph": {
	//                             "text": "Second"
	//                           }
	//                         }
	//                       ],
	//                       "footerWidgets": [
	//                         {
	//                           "textParagraph": {
	//                             "text": "two"
	//                           }
	//                         }
	//                       ]
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

func ExampleToMessage_chips() {
	type rowInput struct {
		Chips []string `googlechat:"chip"`
	}
	row := rowInput{[]string{"Go", "JSON"}}
	message, err := googlechat.ToMessage(struct {
		Row rowInput `googlechat:"chipList"`
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
	//                 "chipList": {
	//                   "chips": [
	//                     {
	//                       "label": "Go"
	//                     },
	//                     {
	//                       "label": "JSON"
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
