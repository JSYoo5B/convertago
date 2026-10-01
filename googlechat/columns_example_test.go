package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleColumns() {
	widget := googlechat.Widget{Content: googlechat.Columns{ColumnItems: []googlechat.Column{
		{
			HorizontalSizeStyle: googlechat.HorizontalSizeStyleFillAvailableSpace,
			VerticalAlignment:   googlechat.ColumnVerticalAlignmentTop,
			Widgets:             []googlechat.ColumnWidget{{Content: googlechat.TextParagraph{Text: "Status: passed"}}},
		},
		{
			HorizontalAlignment: googlechat.HorizontalAlignmentEnd,
			Widgets:             []googlechat.ColumnWidget{{Content: googlechat.TextParagraph{Text: "Duration: 12 seconds"}}},
		},
	}}}
	jsonBytes, err := json.MarshalIndent(widget, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "columns": {
	//     "columnItems": [
	//       {
	//         "horizontalSizeStyle": "FILL_AVAILABLE_SPACE",
	//         "verticalAlignment": "TOP",
	//         "widgets": [
	//           {
	//             "textParagraph": {
	//               "text": "Status: passed"
	//             }
	//           }
	//         ]
	//       },
	//       {
	//         "horizontalAlignment": "END",
	//         "widgets": [
	//           {
	//             "textParagraph": {
	//               "text": "Duration: 12 seconds"
	//             }
	//           }
	//         ]
	//       }
	//     ]
	//   }
	// }
}

func ExampleColumn() {
	value := googlechat.Column{
		HorizontalSizeStyle: googlechat.HorizontalSizeStyleFillAvailableSpace,
		HorizontalAlignment: googlechat.HorizontalAlignmentStart,
		VerticalAlignment:   googlechat.ColumnVerticalAlignmentTop,
		Widgets: []googlechat.ColumnWidget{
			{Content: googlechat.TextParagraph{Text: "Build passed"}},
			{Content: googlechat.Image{ImageURL: "https://example.com/chart.png", AltText: "Build chart"}},
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
	//   "horizontalSizeStyle": "FILL_AVAILABLE_SPACE",
	//   "horizontalAlignment": "START",
	//   "verticalAlignment": "TOP",
	//   "widgets": [
	//     {
	//       "textParagraph": {
	//         "text": "Build passed"
	//       }
	//     },
	//     {
	//       "image": {
	//         "imageUrl": "https://example.com/chart.png",
	//         "altText": "Build chart"
	//       }
	//     }
	//   ]
	// }
}

func ExampleColumnWidget_MarshalJSON() {
	value := googlechat.ColumnWidget{
		Content: googlechat.TextParagraph{Text: "Build **passed**", TextSyntax: googlechat.TextSyntaxMarkdown},
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "textParagraph": {
	//     "text": "Build **passed**",
	//     "textSyntax": "MARKDOWN"
	//   }
	// }
}
