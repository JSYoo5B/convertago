package kakaowork_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleToMessage_actions() {
	type action struct {
		Name  string `kakaowork:"part;slot=name"`
		Value string `kakaowork:"part;slot=value"`
	}
	type button struct {
		Label  string `kakaowork:"part"`
		Action action `kakaowork:"submit_action"`
	}
	type rowInput struct {
		Buttons []button `kakaowork:"button"`
	}
	row := rowInput{[]button{{"확인", action{"confirm", "yes"}}, {"취소", action{"cancel", "no"}}}}
	message, err := kakaowork.ToMessage(struct {
		Row rowInput `kakaowork:"action"`
	}{row})
	data, _ := json.Marshal(message)
	fmt.Println(string(data), err)
	// Output: {"text":"","blocks":[{"type":"action","elements":[{"type":"button","text":"확인","style":"default","action":{"type":"submit_action","name":"confirm","value":"yes"}},{"type":"button","text":"취소","style":"default","action":{"type":"submit_action","name":"cancel","value":"no"}}]}]} <nil>
}

func ExampleToMessage_inlines() {
	type link struct {
		Label string `kakaowork:"part"`
		URL   string `kakaowork:"part;slot=url"`
	}
	type bodyInput struct {
		Before string `kakaowork:"part"`
		Link   link   `kakaowork:"link"`
		After  string `kakaowork:"part;style=bold"`
	}
	body := bodyInput{"문서: ", link{"열기", "https://example.com"}, " 완료"}
	message, err := kakaowork.ToMessage(struct {
		Body bodyInput `kakaowork:"text"`
	}{body})
	data, _ := json.Marshal(message)
	fmt.Println(string(data), err)
	// Output: {"text":"","blocks":[{"type":"text","text":"문서: 열기 완료","inlines":[{"type":"styled","text":"문서: "},{"type":"link","text":"열기","url":"https://example.com"},{"type":"styled","text":" 완료","bold":true}]}]} <nil>
}
