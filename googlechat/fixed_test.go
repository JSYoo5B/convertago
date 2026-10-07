package googlechat

import (
	"encoding/json"
	"testing"

	"github.com/JSYoo5B/convertago/internal/conversion"
	"github.com/stretchr/testify/require"
)

func TestTagFixedValues(t *testing.T) {
	type link struct {
		URL string `googlechat:"openLink"`
	}
	source := struct {
		Body    string `googlechat:"textParagraph;horizontalAlignment=CENTER"`
		Details string `googlechat:"decoratedText;wrapText"`
		Buttons struct {
			Save struct {
				Text  string `googlechat:"part"`
				Click link   `googlechat:"onClick"`
			} `googlechat:"button;type=FILLED_TONAL;disabled"`
		} `googlechat:"buttonList"`
	}{Body: "Centered", Details: "Wrapped"}
	source.Buttons.Save.Text = "Save"
	source.Buttons.Save.Click.URL = "https://example.com"
	message, err := ToMessage(source, conversion.WithWarningAsError())
	require.NoError(t, err)
	data, err := json.Marshal(message.CardsV2[0].Card.Sections[0].Widgets)
	require.NoError(t, err)
	require.JSONEq(t, `[
		{"textParagraph":{"text":"Centered"},"horizontalAlignment":"CENTER"},
		{"decoratedText":{"text":"Wrapped","wrapText":true}},
		{"buttonList":{"buttons":[{"text":"Save","onClick":{"openLink":{"url":"https://example.com"}},"disabled":true,"type":"FILLED_TONAL"}]}}
	]`, string(data))
}
