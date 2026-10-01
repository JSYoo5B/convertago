package convertago_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago"
)

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
