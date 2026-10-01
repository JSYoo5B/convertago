package slack_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/slack"
)

func ExampleToMessage_actions() {
	type button struct {
		Label string `slack:"part"`
		ID    string `slack:"part;slot=action_id"`
		URL   string `slack:"part;slot=url"`
	}
	type rowInput struct {
		Buttons []button `slack:"button"`
	}
	row := rowInput{[]button{{"Open", "open", "https://example.com"}}}
	message, err := slack.ToMessage(struct {
		Row rowInput `slack:"actions"`
	}{row})
	data, _ := json.Marshal(message)
	fmt.Println(string(data), err)
	// Output: {"blocks":[{"type":"actions","elements":[{"type":"button","text":{"type":"plain_text","text":"Open"},"action_id":"open","url":"https://example.com"}]}]} <nil>
}

func ExampleToMessage_context() {
	type image struct {
		URL string `slack:"part;slot=url"`
		Alt string `slack:"part;slot=alt"`
	}
	type rowInput struct {
		Icon image  `slack:"image"`
		Text string `slack:"plain_text"`
	}
	row := rowInput{image{"https://example.com/icon.png", "Author"}, "Posted by Jane"}
	message, err := slack.ToMessage(struct {
		Row rowInput `slack:"context"`
	}{row})
	data, _ := json.Marshal(message)
	fmt.Println(string(data), err)
	// Output: {"blocks":[{"type":"context","elements":[{"type":"image","image_url":"https://example.com/icon.png","alt_text":"Author"},{"type":"plain_text","text":"Posted by Jane"}]}]} <nil>
}
