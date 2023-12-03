package kakaowork_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleSectionBlock() {
	sectionBlock := kakaowork.SectionBlock{
		Content: kakaowork.TextBlock{
			Text: "블록킷\n구성\n정책",
			Inlines: []kakaowork.Inline{
				kakaowork.InlineStyled{Text: "블록킷\n", Bold: true},
				kakaowork.InlineStyled{Text: "구성\n"},
				kakaowork.InlineStyled{Text: "정책", Italic: true},
			},
		},
		Accessory: &kakaowork.ImageBlock{Url: "https://something.storage.host/upload/path/filename"},
		Action:    kakaowork.OpenSystemBrowserAction{Value: "http://example.com/details/999"},
	}

	jsonBytes, err := json.MarshalIndent(sectionBlock, "", "    ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//     "type": "section",
	//     "content": {
	//         "type": "text",
	//         "text": "블록킷\n구성\n정책",
	//         "inlines": [
	//             {
	//                 "type": "styled",
	//                 "text": "블록킷\n",
	//                 "bold": true
	//             },
	//             {
	//                 "type": "styled",
	//                 "text": "구성\n"
	//             },
	//             {
	//                 "type": "styled",
	//                 "text": "정책",
	//                 "italic": true
	//             }
	//         ]
	//     },
	//     "accessory": {
	//         "type": "image_link",
	//         "url": "https://something.storage.host/upload/path/filename"
	//     },
	//     "action": {
	//         "type": "open_system_browser",
	//         "value": "http://example.com/details/999"
	//     }
	// }
}
