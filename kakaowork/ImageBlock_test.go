package kakaowork_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleImageBlock() {
	textBlock := kakaowork.ImageBlock{Url: "https://something.storage.host/upload/path/filename"}

	jsonBytes, err := json.MarshalIndent(textBlock, "", "    ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//     "type": "image_link",
	//     "url": "https://something.storage.host/upload/path/filename"
	// }
}

func TestImageBlock_Validate(t *testing.T) {
	v := validator.New()

	t.Run("Valid url", func(t *testing.T) {
		imageBlock := kakaowork.ImageBlock{
			Url: "https://picsum.photos/200/300",
		}

		err := v.Struct(imageBlock)
		assert.Nil(t, err)
	})

	t.Run("Invalid url", func(t *testing.T) {
		imageBlock := kakaowork.ImageBlock{
			Url: "picsum.photos/200/300",
		}

		err := v.Struct(imageBlock)
		assert.Error(t, err)
	})
}
