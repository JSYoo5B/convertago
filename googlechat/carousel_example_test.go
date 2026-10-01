package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleCarousel() {
	widget := googlechat.Widget{Content: googlechat.Carousel{CarouselCards: []googlechat.CarouselCard{
		{
			Widgets: []googlechat.NestedWidget{{Content: googlechat.TextParagraph{Text: "Build passed"}}},
			FooterWidgets: []googlechat.NestedWidget{{Content: googlechat.ButtonList{Buttons: []googlechat.Button{{
				Text:    "Report",
				OnClick: googlechat.OnClick{OpenLink: &googlechat.OpenLink{URL: "https://example.com/report"}},
			}}}}},
		},
		{
			Widgets: []googlechat.NestedWidget{{Content: googlechat.Image{ImageURL: "https://example.com/chart.png", AltText: "Build chart"}}},
		},
	}}}
	jsonBytes, err := json.MarshalIndent(widget, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "carousel": {
	//     "carouselCards": [
	//       {
	//         "widgets": [
	//           {
	//             "textParagraph": {
	//               "text": "Build passed"
	//             }
	//           }
	//         ],
	//         "footerWidgets": [
	//           {
	//             "buttonList": {
	//               "buttons": [
	//                 {
	//                   "text": "Report",
	//                   "onClick": {
	//                     "openLink": {
	//                       "url": "https://example.com/report"
	//                     }
	//                   }
	//                 }
	//               ]
	//             }
	//           }
	//         ]
	//       },
	//       {
	//         "widgets": [
	//           {
	//             "image": {
	//               "imageUrl": "https://example.com/chart.png",
	//               "altText": "Build chart"
	//             }
	//           }
	//         ]
	//       }
	//     ]
	//   }
	// }
}

func ExampleCarouselCard() {
	value := googlechat.CarouselCard{
		Widgets: []googlechat.NestedWidget{
			{Content: googlechat.TextParagraph{Text: "Build passed"}},
		},
		FooterWidgets: []googlechat.NestedWidget{
			{Content: googlechat.ButtonList{Buttons: []googlechat.Button{{
				Text:    "View report",
				OnClick: googlechat.OnClick{OpenLink: &googlechat.OpenLink{URL: "https://example.com/report"}},
			}}}},
		},
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "widgets": [
	//     {
	//       "textParagraph": {
	//         "text": "Build passed"
	//       }
	//     }
	//   ],
	//   "footerWidgets": [
	//     {
	//       "buttonList": {
	//         "buttons": [
	//           {
	//             "text": "View report",
	//             "onClick": {
	//               "openLink": {
	//                 "url": "https://example.com/report"
	//               }
	//             }
	//           }
	//         ]
	//       }
	//     }
	//   ]
	// }
}

func ExampleNestedWidget_MarshalJSON() {
	value := googlechat.NestedWidget{
		Content: googlechat.Image{ImageURL: "https://example.com/chart.png", AltText: "Build chart"},
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "image": {
	//     "imageUrl": "https://example.com/chart.png",
	//     "altText": "Build chart"
	//   }
	// }
}
