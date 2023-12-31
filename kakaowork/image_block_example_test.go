package kakaowork_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleImageBlock_MarshalJSON() {
	textBlock := kakaowork.ImageBlock{Url: "https://something.storage.host/upload/path/filename"}

	jsonBytes, err := json.MarshalIndent(textBlock, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "image_link",
	//   "url": "https://something.storage.host/upload/path/filename"
	// }
}
