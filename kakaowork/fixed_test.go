package kakaowork

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/JSYoo5B/convertago/internal/conversion"
	"github.com/stretchr/testify/require"
)

func TestTagFixedValues(t *testing.T) {
	type browser struct {
		URL string `kakaowork:"part"`
	}
	type button struct {
		Label  string  `kakaowork:"part"`
		Style  string  `kakaowork:"part;slot=style;omitempty"`
		Action browser `kakaowork:"open_inapp_browser;standalone"`
	}
	source := struct {
		Preview string `kakaowork:"preview"`
		Title   string `kakaowork:"header;style=blue"`
		Body    struct {
			Word struct {
				Text string `kakaowork:"part"`
			} `kakaowork:"styled;color=red;bold;italic"`
		} `kakaowork:"text"`
		Default  button `kakaowork:"button;style=primary"`
		Override button `kakaowork:"button;style=primary"`
	}{Preview: "알림", Title: "알림"}
	source.Body.Word.Text = "긴급"
	source.Default = button{Label: "열기", Action: browser{"https://example.com"}}
	source.Override = button{Label: "삭제", Style: "danger", Action: browser{"https://example.com"}}

	message, err := ToMessage(source, conversion.WithWarningAsError())
	require.NoError(t, err)
	data, err := json.Marshal(message.Blocks)
	require.NoError(t, err)
	require.JSONEq(t, `[
		{"type":"header","text":"알림","style":"blue"},
		{"type":"text","text":"긴급","inlines":[{"type":"styled","text":"긴급","bold":true,"italic":true,"color":"red"}]},
		{"type":"button","text":"열기","style":"primary","action":{"type":"open_inapp_browser","value":"https://example.com","standalone":true}},
		{"type":"button","text":"삭제","style":"danger","action":{"type":"open_inapp_browser","value":"https://example.com","standalone":true}}
	]`, string(data))
}

func TestTagFixedValueUsesRuleCode(t *testing.T) {
	source := struct {
		Title string `kakaowork:"header;style=green"`
	}{"알림"}
	_, err := ToMessage(source)
	var diagnostic conversion.Diagnostic
	require.True(t, errors.As(err, &diagnostic))
	require.Equal(t, ruleHeaderStyleValue.ID, diagnostic.Code)
	require.Equal(t, "$.Title", diagnostic.Path)
}
