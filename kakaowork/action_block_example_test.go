package kakaowork_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleActionBlock_MarshalJSON() {
	block := kakaowork.ActionBlock{Elements: []kakaowork.ButtonBlock{
		{Text: "열기", Style: kakaowork.ButtonStylePrimary, Action: kakaowork.OpenSystemBrowserAction{Value: "https://example.com"}},
		{Text: "확인", Style: kakaowork.ButtonStyleDefault, Action: kakaowork.SubmitAction{Name: "confirm", Value: "yes"}},
	}}
	data, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output:
	// {
	//   "type": "action",
	//   "elements": [
	//     {
	//       "type": "button",
	//       "text": "열기",
	//       "style": "primary",
	//       "action": {
	//         "type": "open_system_browser",
	//         "value": "https://example.com"
	//       }
	//     },
	//     {
	//       "type": "button",
	//       "text": "확인",
	//       "style": "default",
	//       "action": {
	//         "type": "submit_action",
	//         "name": "confirm",
	//         "value": "yes"
	//       }
	//     }
	//   ]
	// }
}
