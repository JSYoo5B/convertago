package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleToMessage() {
	source := struct {
		Text string `slack:"section"`
	}{"Slack notification"}
	message, err := slack.ToMessage(source)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output: {"blocks":[{"type":"section","text":{"type":"plain_text","text":"Slack notification"}}]}
}
