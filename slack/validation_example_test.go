package slack_test

import (
	"fmt"
	"strings"

	"github.com/JSYoo5B/convertago"
	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
)

func ExampleValidate() {
	message := slack.Message{Blocks: []slack.Block{
		slack.ImageBlock{ImageURL: "/images/deploy.png", AltText: "Deployment"},
	}}
	err := slack.Validate(message, convertago.WithDiagnostics(func(d convertago.Diagnostic) {
		fmt.Println(d.Severity, d.Code, d.Path)
	}))
	fmt.Println(err == nil)

	err = slack.Validate(message, convertago.WithWarningAsError())
	fmt.Println(err)

	// Output:
	// advisory message.text.fallback $.text
	// warning image.image_url.format $.blocks[0].image_url
	// true
	// convertago: slack $.blocks[0].image_url: image_url must be an absolute HTTP or HTTPS URL [warning image.image_url.format]
}

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
