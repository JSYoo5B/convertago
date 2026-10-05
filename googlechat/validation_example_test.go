package googlechat_test

import (
	"fmt"

	"github.com/JSYoo5B/convertago"
	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/go-playground/validator/v10"
)

func ExampleValidate() {
	message := googlechat.Message{CardsV2: []googlechat.CardWithID{{Card: googlechat.Card{
		Sections: []googlechat.Section{{Widgets: []googlechat.Widget{
			{Content: googlechat.Image{ImageURL: "http://example.com/image.png"}},
		}}},
	}}}}
	err := googlechat.Validate(message, convertago.WithDiagnostics(func(d convertago.Diagnostic) {
		fmt.Println(d.Severity, d.Code, d.Path)
	}))
	fmt.Println(err == nil)

	err = googlechat.Validate(message, convertago.WithWarningAsError())
	fmt.Println(err)

	// Output:
	// warning image.image_url.https $.cardsV2[0].card.sections[0].widgets[0].image.imageUrl
	// true
	// convertago: googlechat $.cardsV2[0].card.sections[0].widgets[0].image.imageUrl: imageUrl must be an HTTPS URL [warning image.image_url.https]
}

func ExampleRegisterValidation() {
	v := validator.New()
	googlechat.RegisterValidation(v)

	icon := googlechat.Icon{KnownIcon: "STAR"}
	fmt.Println(v.Struct(icon) == nil)

	icon.IconURL = "https://example.com/icon.png"
	fmt.Println(v.Struct(icon) == nil)

	// Output:
	// true
	// false
}
