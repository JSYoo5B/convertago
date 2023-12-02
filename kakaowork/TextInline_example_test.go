package kakaowork_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleInlineStyled() {
	textBlock := kakaowork.TextBlock{
		Text: "카카오워크 텍스트 블록 변경",
		Inlines: []kakaowork.Inline{
			&kakaowork.InlineStyled{Text: "카카오워크 "},
			&kakaowork.InlineStyled{Text: "텍스트", Bold: true},
			&kakaowork.InlineStyled{Text: " 블록 변경", Color: kakaowork.InlineColorRed},
		},
	}

	jsonBytes, err := json.MarshalIndent(textBlock, "", "    ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//     "type": "text",
	//     "text": "카카오워크 텍스트 블록 변경",
	//     "inlines": [
	//         {
	//             "type": "styled",
	//             "text": "카카오워크 "
	//         },
	//         {
	//             "type": "styled",
	//             "text": "텍스트",
	//             "bold": true
	//         },
	//         {
	//             "type": "styled",
	//             "text": " 블록 변경",
	//             "color": "red"
	//         }
	//     ]
	// }
}

func ExampleInlineLink() {
	textBlock := kakaowork.TextBlock{
		Text: "카카오워크 기술문서",
		Inlines: []kakaowork.Inline{
			&kakaowork.InlineLink{Text: "카카오워크 기술 문서", Url: "https://blog.kakaowork.com/43"},
		},
	}

	jsonBytes, err := json.MarshalIndent(textBlock, "", "    ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//     "type": "text",
	//     "text": "카카오워크 기술문서",
	//     "inlines": [
	//         {
	//             "type": "link",
	//             "text": "카카오워크 기술 문서",
	//             "url": "https://blog.kakaowork.com/43"
	//         }
	//     ]
	// }
}

func ExampleInlineMention() {
	textBlock := kakaowork.TextBlock{
		Text: "@ryan",
		Inlines: []kakaowork.Inline{
			&kakaowork.InlineMention{Text: "@ryan", UserId: 800},
		},
	}

	jsonBytes, err := json.MarshalIndent(textBlock, "", "    ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//     "type": "text",
	//     "text": "@ryan",
	//     "inlines": [
	//         {
	//             "type": "mention",
	//             "text": "@ryan",
	//             "ref": {
	//                 "type": "kw",
	//                 "value": 800
	//             }
	//         }
	//     ]
	// }
}
