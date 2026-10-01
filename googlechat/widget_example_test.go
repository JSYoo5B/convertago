package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleWidget_MarshalJSON() {
	widget := googlechat.Widget{
		Content:             googlechat.TextParagraph{Text: "Build **passed**", MaxLines: 2, TextSyntax: googlechat.TextSyntaxMarkdown},
		HorizontalAlignment: googlechat.HorizontalAlignmentCenter,
	}
	jsonBytes, err := json.MarshalIndent(widget, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "horizontalAlignment": "CENTER",
	//   "textParagraph": {
	//     "text": "Build **passed**",
	//     "maxLines": 2,
	//     "textSyntax": "MARKDOWN"
	//   }
	// }
}
