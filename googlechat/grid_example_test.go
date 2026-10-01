package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleGridItem() {
	value := googlechat.GridItem{
		ID:       "report",
		Image:    &googlechat.ImageComponent{ImageURI: "https://example.com/report.png", AltText: "Build report"},
		Title:    "Report",
		Subtitle: "Build 42",
		Layout:   googlechat.GridItemLayoutTextBelow,
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "id": "report",
	//   "image": {
	//     "imageUri": "https://example.com/report.png",
	//     "altText": "Build report"
	//   },
	//   "title": "Report",
	//   "subtitle": "Build 42",
	//   "layout": "TEXT_BELOW"
	// }
}

func ExampleImageComponent() {
	value := googlechat.ImageComponent{
		ImageURI:    "https://example.com/avatar.png",
		AltText:     "Build bot",
		CropStyle:   &googlechat.ImageCropStyle{Type: googlechat.ImageCropTypeCircle},
		BorderStyle: &googlechat.BorderStyle{Type: googlechat.BorderTypeNone},
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "imageUri": "https://example.com/avatar.png",
	//   "altText": "Build bot",
	//   "cropStyle": {
	//     "type": "CIRCLE"
	//   },
	//   "borderStyle": {
	//     "type": "NO_BORDER"
	//   }
	// }
}

func ExampleImageCropStyle() {
	value := googlechat.ImageCropStyle{
		Type:        googlechat.ImageCropTypeRectangleCustom,
		AspectRatio: 1.5,
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "type": "RECTANGLE_CUSTOM",
	//   "aspectRatio": 1.5
	// }
}

func ExampleBorderStyle() {
	value := googlechat.BorderStyle{
		Type:         googlechat.BorderTypeStroke,
		StrokeColor:  &googlechat.Color{Red: 0.25, Green: 0.5, Blue: 1},
		CornerRadius: 4,
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "type": "STROKE",
	//   "strokeColor": {
	//     "red": 0.25,
	//     "green": 0.5,
	//     "blue": 1
	//   },
	//   "cornerRadius": 4
	// }
}

func ExampleGrid() {
	widget := googlechat.Widget{Content: googlechat.Grid{
		Title: "Build artifacts",
		Items: []googlechat.GridItem{{
			ID: "report",
			Image: &googlechat.ImageComponent{
				ImageURI:  "https://example.com/report.png",
				AltText:   "Build report",
				CropStyle: &googlechat.ImageCropStyle{Type: googlechat.ImageCropTypeRectangle4By3},
			},
			Title:  "Report",
			Layout: googlechat.GridItemLayoutTextBelow,
		}},
		BorderStyle: &googlechat.BorderStyle{Type: googlechat.BorderTypeStroke, CornerRadius: 4},
		ColumnCount: 2,
	}}
	jsonBytes, err := json.MarshalIndent(widget, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "grid": {
	//     "title": "Build artifacts",
	//     "items": [
	//       {
	//         "id": "report",
	//         "image": {
	//           "imageUri": "https://example.com/report.png",
	//           "altText": "Build report",
	//           "cropStyle": {
	//             "type": "RECTANGLE_4_3"
	//           }
	//         },
	//         "title": "Report",
	//         "layout": "TEXT_BELOW"
	//       }
	//     ],
	//     "borderStyle": {
	//       "type": "STROKE",
	//       "cornerRadius": 4
	//     },
	//     "columnCount": 2
	//   }
	// }
}
