package kakaowork_test

import (
	"encoding/json"
	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSectionBlock_OptionalFields(t *testing.T) {
	t.Run("Accessory Empty", func(t *testing.T) {
		section := kakaowork.SectionBlock{
			Content:   kakaowork.TextBlock{Text: "example"},
			Accessory: nil,
			Action:    kakaowork.OpenSystemBrowserAction{Value: "http://example.com/details/999"},
		}

		var jsonBytes []byte
		var err error
		if jsonBytes, err = json.Marshal(section); err != nil {
			t.Fatal(err)
		}
		jsonMarshalString := string(jsonBytes)

		assert.NotContains(t, jsonMarshalString, "accessory")
	})

	t.Run("Action Empty", func(t *testing.T) {
		section := kakaowork.SectionBlock{
			Content:   kakaowork.TextBlock{Text: "example"},
			Accessory: &kakaowork.ImageBlock{Url: "https://something.storage.host/upload/path/filename"},
			Action:    nil,
		}

		var jsonBytes []byte
		var err error
		if jsonBytes, err = json.Marshal(section); err != nil {
			t.Fatal(err)
		}
		jsonMarshalString := string(jsonBytes)

		assert.NotContains(t, jsonMarshalString, "action")
	})
}
