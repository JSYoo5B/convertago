package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleTextParagraph() {
	paragraph := googlechat.TextParagraph{
		Text:       "Build **passed**",
		TextSyntax: googlechat.TextSyntaxMarkdown,
	}
	widget := googlechat.Widget{Content: paragraph}
	jsonBytes, err := json.MarshalIndent(widget, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "textParagraph": {
	//     "text": "Build **passed**",
	//     "textSyntax": "MARKDOWN"
	//   }
	// }
}
