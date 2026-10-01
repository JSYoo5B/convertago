package slack_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/slack"
)

func ExampleToMessage_richTextList() {
	type list struct {
		Style string   `slack:"part;slot=style"`
		Items []string `slack:"rich_text_section"`
	}
	source := struct {
		Body struct {
			List list `slack:"rich_text_list"`
		} `slack:"rich_text"`
	}{}
	source.Body.List = list{"ordered", []string{"First", "Second"}}
	message, err := slack.ToMessage(source)
	data, _ := json.Marshal(message)
	fmt.Println(string(data), err)
	// Output: {"blocks":[{"type":"rich_text","elements":[{"type":"rich_text_list","style":"ordered","elements":[{"type":"rich_text_section","elements":[{"type":"text","text":"First"}]},{"type":"rich_text_section","elements":[{"type":"text","text":"Second"}]}]}]}]} <nil>
}

func ExampleToMessage_richTextMention() {
	source := struct {
		Body struct {
			Text  string `slack:"part"`
			User  string `slack:"user;style=bold"`
			Emoji string `slack:"emoji"`
		} `slack:"rich_text"`
	}{}
	source.Body.Text = "Hello "
	source.Body.User = "U123"
	source.Body.Emoji = "wave"
	message, err := slack.ToMessage(source)
	data, _ := json.Marshal(message)
	fmt.Println(string(data), err)
	// Output: {"blocks":[{"type":"rich_text","elements":[{"type":"rich_text_section","elements":[{"type":"text","text":"Hello "},{"type":"user","user_id":"U123","style":{"bold":true}},{"type":"emoji","name":"wave"}]}]}]} <nil>
}
