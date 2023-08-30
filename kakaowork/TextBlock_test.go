package kakaowork_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleTextBlock() {
	textBlock := kakaowork.TextBlock{
		Text: "카카오워크 텍스트 블록이 변경되었습니다. 자세한 설명은 기술 문서를 참고하시기 바랍니다. 담당자는 @ryan 입니다.",
		Inlines: []kakaowork.Inline{
			&kakaowork.InlineStyled{Text: "카카오워크 "},
			&kakaowork.InlineStyled{Text: "텍스트", Bold: true, Color: kakaowork.InlineColorRed},
			&kakaowork.InlineStyled{Text: " 블록이 변경되었습니다. 자세한 설명은 "},
			&kakaowork.InlineLink{Text: "기술 문서", Url: "https://blog.kakaowork.com/43"},
			&kakaowork.InlineStyled{Text: "를 참고하시기 바랍니다. 담당자는"},
			&kakaowork.InlineMention{Text: "@ryan", UserId: 0},
			&kakaowork.InlineStyled{Text: " 입니다."},
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
	//     "text": "카카오워크 텍스트 블록이 변경되었습니다. 자세한 설명은 기술 문서를 참고하시기 바랍니다. 담당자는 @ryan 입니다.",
	//     "inlines": [
	//         {
	//             "type": "styled",
	//             "text": "카카오워크 "
	//         },
	//         {
	//             "type": "styled",
	//             "text": "텍스트",
	//             "bold": true,
	//             "color": "red"
	//         },
	//         {
	//             "type": "styled",
	//             "text": " 블록이 변경되었습니다. 자세한 설명은 "
	//         },
	//         {
	//             "type": "link",
	//             "text": "기술 문서",
	//             "url": "https://blog.kakaowork.com/43"
	//         },
	//         {
	//             "type": "styled",
	//             "text": "를 참고하시기 바랍니다. 담당자는"
	//         },
	//         {
	//             "type": "mention",
	//             "text": "@ryan",
	//             "ref": {
	//                 "type": "kw",
	//                 "value": 0
	//             }
	//         },
	//         {
	//             "type": "styled",
	//             "text": " 입니다."
	//         }
	//     ]
	// }
}
