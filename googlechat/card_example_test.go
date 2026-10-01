package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleCard() {
	card := googlechat.Card{
		Header: &googlechat.CardHeader{
			Title:        "Build report",
			ImageURL:     "https://example.com/avatar.png",
			ImageType:    googlechat.ImageTypeCircle,
			ImageAltText: "Build bot",
		},
		Sections: []googlechat.Section{{
			Header: "Details",
			Widgets: []googlechat.Widget{
				{Content: googlechat.TextParagraph{Text: "Build passed"}},
				{Content: googlechat.TextParagraph{Text: "Duration: 12 seconds"}},
			},
			Collapsible:               true,
			UncollapsibleWidgetsCount: 1,
		}},
		SectionDividerStyle: googlechat.DividerStyleNone,
	}
	jsonBytes, err := json.MarshalIndent(card, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "header": {
	//     "title": "Build report",
	//     "imageUrl": "https://example.com/avatar.png",
	//     "imageType": "CIRCLE",
	//     "imageAltText": "Build bot"
	//   },
	//   "sections": [
	//     {
	//       "header": "Details",
	//       "widgets": [
	//         {
	//           "textParagraph": {
	//             "text": "Build passed"
	//           }
	//         },
	//         {
	//           "textParagraph": {
	//             "text": "Duration: 12 seconds"
	//           }
	//         }
	//       ],
	//       "collapsible": true,
	//       "uncollapsibleWidgetsCount": 1
	//     }
	//   ],
	//   "sectionDividerStyle": "NO_DIVIDER"
	// }
}

func ExampleCardHeader() {
	value := googlechat.CardHeader{
		Title:        "Build report",
		Subtitle:     "Build 42",
		ImageURL:     "https://example.com/avatar.png",
		ImageType:    googlechat.ImageTypeCircle,
		ImageAltText: "Build bot",
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "title": "Build report",
	//   "subtitle": "Build 42",
	//   "imageUrl": "https://example.com/avatar.png",
	//   "imageType": "CIRCLE",
	//   "imageAltText": "Build bot"
	// }
}

func ExampleSection() {
	value := googlechat.Section{
		Header: "Details",
		Widgets: []googlechat.Widget{
			{Content: googlechat.TextParagraph{Text: "Build passed"}},
			{Content: googlechat.TextParagraph{Text: "Duration: 12 seconds"}},
		},
		Collapsible:               true,
		UncollapsibleWidgetsCount: 1,
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "header": "Details",
	//   "widgets": [
	//     {
	//       "textParagraph": {
	//         "text": "Build passed"
	//       }
	//     },
	//     {
	//       "textParagraph": {
	//         "text": "Duration: 12 seconds"
	//       }
	//     }
	//   ],
	//   "collapsible": true,
	//   "uncollapsibleWidgetsCount": 1
	// }
}

func ExampleCollapseControl() {
	value := googlechat.CollapseControl{
		HorizontalAlignment: googlechat.HorizontalAlignmentEnd,
		ExpandButton: googlechat.Button{
			Text:    "Show details",
			OnClick: googlechat.OnClick{Action: &googlechat.Action{Function: "expandSection"}},
			Type:    googlechat.ButtonTypeBorderless,
		},
		CollapseButton: googlechat.Button{
			Text:    "Hide details",
			OnClick: googlechat.OnClick{Action: &googlechat.Action{Function: "collapseSection"}},
			Type:    googlechat.ButtonTypeBorderless,
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
	//   "horizontalAlignment": "END",
	//   "expandButton": {
	//     "text": "Show details",
	//     "onClick": {
	//       "action": {
	//         "function": "expandSection"
	//       }
	//     },
	//     "type": "BORDERLESS"
	//   },
	//   "collapseButton": {
	//     "text": "Hide details",
	//     "onClick": {
	//       "action": {
	//         "function": "collapseSection"
	//       }
	//     },
	//     "type": "BORDERLESS"
	//   }
	// }
}
