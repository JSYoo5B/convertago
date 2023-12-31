package kakaowork_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleContextBlock_MarshalJSON() {
	contextBlock := kakaowork.ContextBlock{
		Content: kakaowork.TextBlock{
			Text: "카카오 판교 오피스",
			Inlines: []kakaowork.Inline{
				&kakaowork.InlineLink{Text: "카카오 판교 오피스", Url: "http://kko.to/RRWQwZQj0"},
			},
		},
		Image: kakaowork.ImageBlock{Url: "https://something.storage.host/upload/path/filename"},
	}

	jsonBytes, err := json.MarshalIndent(contextBlock, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "context",
	//   "content": {
	//     "type": "text",
	//     "text": "카카오 판교 오피스",
	//     "inlines": [
	//       {
	//         "type": "link",
	//         "text": "카카오 판교 오피스",
	//         "url": "http://kko.to/RRWQwZQj0"
	//       }
	//     ]
	//   },
	//   "image": {
	//     "type": "image_link",
	//     "url": "https://something.storage.host/upload/path/filename"
	//   }
	// }
}
