package googlechat_test

import (
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/go-playground/validator/v10"
)

func ExampleRegisterValidation() {
	v := validator.New()
	googlechat.RegisterValidation(v)

	image := googlechat.Image{ImageURL: "https://example.com/image.png"}
	fmt.Println(v.Struct(image) == nil)

	image.ImageURL = "http://example.com/image.png"
	fmt.Println(v.Struct(image) == nil)

	// Output:
	// true
	// false
}
