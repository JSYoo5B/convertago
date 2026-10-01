package convertago_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago"
)

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
