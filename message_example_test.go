package convertago_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago"
)

func ExampleToKakaoworkMessage() {
	source := struct {
		Preview string `kakaowork:"preview"`
		Title   string `kakaowork:"header"`
		Prefix  string `kakaowork:"text;group=body"`
		Name    string `kakaowork:"text;group=body;style=bold"`
		Image   string `kakaowork:"image_link"`
	}{"새 알림", "알림", "안녕하세요 ", "홍길동", "https://example.com/image.png"}
	message, err := convertago.ToKakaoworkMessage(source)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output: {"text":"새 알림","blocks":[{"type":"header","text":"알림","style":"white"},{"type":"text","text":"안녕하세요 홍길동","inlines":[{"type":"styled","text":"안녕하세요 "},{"type":"styled","text":"홍길동","bold":true}]},{"type":"image_link","url":"https://example.com/image.png"}]}
}

func ExampleToSlackMessage() {
	source := struct {
		Title  string `slack:"header"`
		Prefix string `slack:"rich_text;group=body"`
		Name   string `slack:"rich_text;group=body;style=bold"`
		Image  struct {
			URL string `slack:"part;slot=url"`
			Alt string `slack:"part;slot=alt"`
		} `slack:"image"`
	}{Title: "Notice", Prefix: "Hello ", Name: "Jane"}
	source.Image.URL = "https://example.com/image.png"
	source.Image.Alt = "Team photo"
	message, err := convertago.ToSlackMessage(source)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output: {"blocks":[{"type":"header","text":{"type":"plain_text","text":"Notice"}},{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[{"type":"text","text":"Hello "},{"type":"text","text":"Jane","style":{"bold":true}}]}]},{"type":"image","image_url":"https://example.com/image.png","alt_text":"Team photo"}]}
}

func ExampleToGoogleChatMessage() {
	source := struct {
		Title  string `googlechat:"header"`
		Prefix string `googlechat:"textParagraph;group=body"`
		Name   string `googlechat:"textParagraph;group=body;style=bold"`
		Image  string `googlechat:"image"`
	}{"Notice", "Hello ", "Jane", "https://example.com/image.png"}
	message, err := convertago.ToGoogleChatMessage(source)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output: {"cardsV2":[{"card":{"header":{"title":"Notice"},"sections":[{"widgets":[{"textParagraph":{"text":"Hello \u003cb\u003eJane\u003c/b\u003e"}},{"image":{"imageUrl":"https://example.com/image.png"}}]}]}}]}
}
