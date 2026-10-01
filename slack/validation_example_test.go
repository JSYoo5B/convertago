package slack_test

import (
	"fmt"
	"strings"

	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
)

func ExampleRegisterValidation() {
	v := validator.New()
	slack.RegisterValidation(v)

	header := slack.HeaderBlock{Text: slack.PlainTextObject{Text: "Deployment complete"}}
	fmt.Println(v.Struct(header) == nil)

	header.Text.Text = strings.Repeat("a", 151)
	fmt.Println(v.Struct(header) == nil)

	// Output:
	// true
	// false
}
