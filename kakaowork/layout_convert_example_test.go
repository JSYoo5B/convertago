package kakaowork_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleToMessage_header() {
	source := struct {
		Preview string `kakaowork:"preview"`
		Title   struct {
			Text  string `kakaowork:"part"`
			Style string `kakaowork:"part;slot=style"`
		} `kakaowork:"header"`
		Body string `kakaowork:"text"`
	}{Preview: "결재 요청", Body: "휴가 신청서가 도착했습니다."}
	source.Title.Text, source.Title.Style = "결재 요청", "blue"
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
	//   "text": "결재 요청",
	//   "blocks": [
	//     {
	//       "type": "header",
	//       "text": "결재 요청",
	//       "style": "blue"
	//     },
	//     {
	//       "type": "text",
	//       "text": "휴가 신청서가 도착했습니다."
	//     }
	//   ]
	// }
}

func ExampleToMessage_layouts() {
	type content struct {
		Text string `kakaowork:"part"`
	}
	type image struct {
		URL string `kakaowork:"part;slot=url"`
	}
	type browser struct {
		URL string `kakaowork:"part"`
	}
	source := struct {
		Preview string `kakaowork:"preview"`
		Banner  string `kakaowork:"image_link"`
		Details struct {
			Body   content `kakaowork:"text"`
			Term   string  `kakaowork:"part;slot=term"`
			Accent bool    `kakaowork:"part;slot=accent"`
		} `kakaowork:"description"`
		Divider struct{} `kakaowork:"divider"`
		Summary struct {
			Body      content `kakaowork:"text"`
			Thumbnail image   `kakaowork:"image_link;slot=accessory"`
			Open      browser `kakaowork:"open_system_browser;slot=action"`
		} `kakaowork:"section"`
		Author struct {
			Body   content `kakaowork:"text"`
			Avatar image   `kakaowork:"image_link;slot=image"`
		} `kakaowork:"context"`
	}{Preview: "배포 결과", Banner: "https://example.com/banner.png"}
	source.Details.Body.Text, source.Details.Term, source.Details.Accent = "성공", "상태", true
	source.Summary.Body.Text = "변경 사항 12건"
	source.Summary.Thumbnail.URL = "https://example.com/thumb.png"
	source.Summary.Open.URL = "https://example.com/deploys/1"
	source.Author.Body.Text = "배포 봇"
	source.Author.Avatar.URL = "https://example.com/bot.png"
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
	//   "text": "배포 결과",
	//   "blocks": [
	//     {
	//       "type": "image_link",
	//       "url": "https://example.com/banner.png"
	//     },
	//     {
	//       "type": "description",
	//       "content": {
	//         "type": "text",
	//         "text": "성공"
	//       },
	//       "term": "상태",
	//       "accent": true
	//     },
	//     {
	//       "type": "divider"
	//     },
	//     {
	//       "type": "section",
	//       "content": {
	//         "type": "text",
	//         "text": "변경 사항 12건"
	//       },
	//       "accessory": {
	//         "type": "image_link",
	//         "url": "https://example.com/thumb.png"
	//       },
	//       "action": {
	//         "type": "open_system_browser",
	//         "value": "https://example.com/deploys/1"
	//       }
	//     },
	//     {
	//       "type": "context",
	//       "content": {
	//         "type": "text",
	//         "text": "배포 봇"
	//       },
	//       "image": {
	//         "type": "image_link",
	//         "url": "https://example.com/bot.png"
	//       }
	//     }
	//   ]
	// }
}

func ExampleToMessage_buttonActions() {
	type browser struct {
		URL        string `kakaowork:"part"`
		Standalone bool   `kakaowork:"part;slot=standalone"`
		Width      int    `kakaowork:"part;slot=width"`
		Height     int    `kakaowork:"part;slot=height"`
	}
	type exclusive struct {
		Default struct {
			URL string `kakaowork:"part"`
		} `kakaowork:"open_system_browser"`
		Mobile struct {
			Value string `kakaowork:"part"`
		} `kakaowork:"open_external_app;slot=mobile"`
	}
	source := struct {
		Preview string `kakaowork:"preview"`
		Popup   struct {
			Label  string  `kakaowork:"part"`
			Action browser `kakaowork:"open_inapp_browser"`
		} `kakaowork:"button"`
		Map struct {
			Label  string    `kakaowork:"part"`
			Style  string    `kakaowork:"part;slot=style"`
			Action exclusive `kakaowork:"exclusive"`
		} `kakaowork:"button"`
		Approve struct {
			Label  string `kakaowork:"part"`
			Action struct {
				Value string `kakaowork:"part"`
			} `kakaowork:"call_modal"`
		} `kakaowork:"button"`
	}{Preview: "회의 안내"}
	source.Popup.Label = "회의록"
	source.Popup.Action = browser{"https://example.com/notes", true, 800, 600}
	source.Map.Label, source.Map.Style = "위치 보기", "primary"
	source.Map.Action.Default.URL = "https://map.kakao.com"
	source.Map.Action.Mobile.Value = "ios=kakaomap%3A%2F%2Flook&aos=kakaomap%3A%2F%2Flook"
	source.Approve.Label = "참석 응답"
	source.Approve.Action.Value = "meeting=42"
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
	//   "text": "회의 안내",
	//   "blocks": [
	//     {
	//       "type": "button",
	//       "text": "회의록",
	//       "style": "default",
	//       "action": {
	//         "type": "open_inapp_browser",
	//         "value": "https://example.com/notes",
	//         "standalone": true,
	//         "width": 800,
	//         "height": 600
	//       }
	//     },
	//     {
	//       "type": "button",
	//       "text": "위치 보기",
	//       "style": "primary",
	//       "action": {
	//         "type": "exclusive",
	//         "default": {
	//           "type": "open_system_browser",
	//           "value": "https://map.kakao.com"
	//         },
	//         "mobile": {
	//           "type": "open_external_app",
	//           "value": "ios=kakaomap%3A%2F%2Flook\u0026aos=kakaomap%3A%2F%2Flook"
	//         }
	//       }
	//     },
	//     {
	//       "type": "button",
	//       "text": "참석 응답",
	//       "style": "default",
	//       "action": {
	//         "type": "call_modal",
	//         "value": "meeting=42"
	//       }
	//     }
	//   ]
	// }
}
