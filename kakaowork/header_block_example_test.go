package kakaowork_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleHeaderBlock_MarshalJSON() {
	headerBlock := kakaowork.HeaderBlock{
		Text:  "내게 요청 온 결재",
		Style: kakaowork.HeaderStyleBlue,
	}

	jsonBytes, err := json.MarshalIndent(headerBlock, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "header",
	//   "text": "내게 요청 온 결재",
	//   "style": "blue"
	// }
}
