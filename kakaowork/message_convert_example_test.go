package kakaowork_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleToMessage() {
	source := struct {
		Preview string `kakaowork:"preview"`
		Text    string `kakaowork:"text"`
	}{"새 알림", "카카오워크 알림"}
	message, err := kakaowork.ToMessage(source)
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output:
	// {
	//   "text": "새 알림",
	//   "blocks": [
	//     {
	//       "type": "text",
	//       "text": "카카오워크 알림"
	//     }
	//   ]
	// }
}
