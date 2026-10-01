package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleImage() {
	widget := googlechat.Widget{Content: googlechat.Image{
		ImageURL: "https://example.com/chart.png",
		OnClick:  &googlechat.OnClick{OpenLink: &googlechat.OpenLink{URL: "https://example.com/report"}},
		AltText:  "Build duration chart",
	}}
	jsonBytes, err := json.MarshalIndent(widget, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "image": {
	//     "imageUrl": "https://example.com/chart.png",
	//     "onClick": {
	//       "openLink": {
	//         "url": "https://example.com/report"
	//       }
	//     },
	//     "altText": "Build duration chart"
	//   }
	// }
}

func ExampleDivider() {
	jsonBytes, err := json.MarshalIndent(googlechat.Widget{Content: googlechat.Divider{}}, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "divider": {}
	// }
}
