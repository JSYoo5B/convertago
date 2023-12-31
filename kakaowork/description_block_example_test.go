package kakaowork_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleDescriptionBlock_MarshalJSON() {
	descriptionBlock := kakaowork.DescriptionBlock{
		Term: "일시",
		Content: kakaowork.TextBlock{
			Text: "2020년 2월 22일 오후 2시",
			Inlines: []kakaowork.Inline{
				&kakaowork.InlineStyled{Text: "2020년 2월 22일 오후 2시", Bold: true},
			},
		},
		Accent: true,
	}

	jsonBytes, err := json.MarshalIndent(descriptionBlock, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "description",
	//   "content": {
	//     "type": "text",
	//     "text": "2020년 2월 22일 오후 2시",
	//     "inlines": [
	//       {
	//         "type": "styled",
	//         "text": "2020년 2월 22일 오후 2시",
	//         "bold": true
	//       }
	//     ]
	//   },
	//   "term": "일시",
	//   "accent": true
	// }
}
