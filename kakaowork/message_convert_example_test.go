package kakaowork_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleToMessage() {
	source := struct {
		Text string `kakaowork:"text"`
	}{"카카오워크 알림"}
	message, err := kakaowork.ToMessage(source)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output: {"text":"","blocks":[{"type":"text","text":"카카오워크 알림"}]}
}
