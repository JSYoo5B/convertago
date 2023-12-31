package kakaowork_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleOpenSystemBrowserAction_MarshalJSON() {
	buttonBlock := kakaowork.ButtonBlock{
		Text:  "자세히보기",
		Style: kakaowork.ButtonStylePrimary,
		Action: kakaowork.OpenSystemBrowserAction{
			Name:  "button1",
			Value: "http://example.com/details/999",
		},
	}

	jsonBytes, err := json.MarshalIndent(buttonBlock, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "button",
	//   "text": "자세히보기",
	//   "style": "primary",
	//   "action": {
	//     "type": "open_system_browser",
	//     "name": "button1",
	//     "value": "http://example.com/details/999"
	//   }
	// }
}

func ExampleOpenInAppBrowserAction_MarshalJSON() {
	buttonBlock := kakaowork.ButtonBlock{
		Text:  "새 창에서 보기",
		Style: kakaowork.ButtonStylePrimary,
		Action: kakaowork.OpenInAppBrowserAction{
			Name:       "button1",
			Value:      "http://example.com/details/999",
			Standalone: true,
			Width:      980,
			Height:     720,
		},
	}

	jsonBytes, err := json.MarshalIndent(buttonBlock, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "button",
	//   "text": "새 창에서 보기",
	//   "style": "primary",
	//   "action": {
	//     "type": "open_inapp_browser",
	//     "name": "button1",
	//     "value": "http://example.com/details/999",
	//     "standalone": true,
	//     "width": 980,
	//     "height": 720
	//   }
	// }
}

func ExampleOpenExternalAppAction_MarshalJSON() {
	buttonBlock := kakaowork.ButtonBlock{
		Text: "카카오 지도에서 보기",
		Action: kakaowork.OpenExternalAppAction{
			Value: "ios=kakaomap%3A%2F%2Flook%3Fp%3D37.537229%2C127.005515",
		},
	}

	jsonBytes, err := json.MarshalIndent(buttonBlock, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "button",
	//   "text": "카카오 지도에서 보기",
	//   "action": {
	//     "type": "open_external_app",
	//     "value": "ios=kakaomap%3A%2F%2Flook%3Fp%3D37.537229%2C127.005515"
	//   }
	// }
}

func ExampleSubmitAction_MarshalJSON() {
	buttonBlock := kakaowork.ButtonBlock{
		Text: "일정 수락",
		Action: kakaowork.SubmitAction{
			Name:  "accept",
			Value: "event_id=20210301-0024",
		},
	}

	jsonBytes, err := json.MarshalIndent(buttonBlock, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "button",
	//   "text": "일정 수락",
	//   "action": {
	//     "type": "submit_action",
	//     "name": "accept",
	//     "value": "event_id=20210301-0024"
	//   }
	// }
}

func ExampleCallModalAction_MarshalJSON() {
	buttonBlock := kakaowork.ButtonBlock{
		Text: "결재창 띄우기",
		Action: kakaowork.CallModalAction{
			Value: "number=20200401-PR-0024",
		},
	}

	jsonBytes, err := json.MarshalIndent(buttonBlock, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "button",
	//   "text": "결재창 띄우기",
	//   "action": {
	//     "type": "call_modal",
	//     "value": "number=20200401-PR-0024"
	//   }
	// }
}

func ExampleExclusiveAction_MarshalJSON() {
	buttonBlock := kakaowork.ButtonBlock{
		Text:  "자세히보기",
		Style: kakaowork.ButtonStylePrimary,
		Action: kakaowork.ExclusiveAction{
			Default: kakaowork.OpenSystemBrowserAction{
				Value: "https://example.com",
			},
			Ios: kakaowork.OpenInAppBrowserAction{
				Value: "http://m.example.com",
			},
			Android: kakaowork.OpenExternalAppAction{
				Value: "geo:37.537229,127.005515",
			},
		},
	}

	jsonBytes, err := json.MarshalIndent(buttonBlock, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//   "type": "button",
	//   "text": "자세히보기",
	//   "style": "primary",
	//   "action": {
	//     "type": "exclusive",
	//     "default": {
	//       "type": "open_system_browser",
	//       "value": "https://example.com"
	//     },
	//     "android": {
	//       "type": "open_external_app",
	//       "value": "geo:37.537229,127.005515"
	//     },
	//     "ios": {
	//       "type": "open_inapp_browser",
	//       "value": "http://m.example.com"
	//     }
	//   }
	// }
}
