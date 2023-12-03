package kakaowork_test

import (
	"bytes"
	"encoding/json"
	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html"
	"testing"
)

func TestAmpersandEscapeIssue(t *testing.T) {
	openExtAction := kakaowork.OpenExternalAppAction{
		Value: "ios=kakaomap%3A%2F%2Flook%3Fp%3D37.537229%2C127.005515&aos=kakaomap%3A%2F%2Flook%3Fp%3D37.537229%2C127.005515",
	}
	require.Contains(t, openExtAction.Value, "&")
	require.NotContains(t, openExtAction.Value, `\u0026`)

	marshaledJson := `
	{
		"type": "open_external_app",
		"value": "ios=kakaomap%3A%2F%2Flook%3Fp%3D37.537229%2C127.005515\u0026aos=kakaomap%3A%2F%2Flook%3Fp%3D37.537229%2C127.005515"
	}`
	require.Contains(t, marshaledJson, `\u0026`)
	require.NotContains(t, marshaledJson, "&")

	t.Run("Default marshal with ampersand", func(t *testing.T) {
		var jsonBytes []byte
		var err error
		if jsonBytes, err = json.Marshal(openExtAction); err != nil {
			t.Fatal(err)
		}
		jsonMarshalString := string(jsonBytes)

		assert.Contains(t, jsonMarshalString, `\u0026`)
		assert.NotContains(t, jsonMarshalString, "&")
	})

	t.Run("New Encoder marshaling with ampersand", func(t *testing.T) {
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(&openExtAction); err != nil {
			t.Fatal(err)
		}
		newEncoderString := html.UnescapeString(buf.String())

		assert.Contains(t, newEncoderString, `\u0026`)
		assert.NotContains(t, newEncoderString, "&")
	})

	t.Run("Unmarshal", func(t *testing.T) {
		var unmarshalledJson map[string]any
		var exists bool
		var value string
		if err := json.Unmarshal([]byte(marshaledJson), &unmarshalledJson); err != nil {
			t.Fatal(err)
		}

		_, exists = unmarshalledJson["value"]
		require.True(t, exists)
		value, exists = unmarshalledJson["value"].(string)
		require.True(t, exists)

		assert.Contains(t, value, `&`)
		assert.NotContains(t, value, `\u0026`)
	})
}
