package slack

import (
	"encoding/json"
	"testing"

	"github.com/JSYoo5B/convertago/internal/conversion"
	"github.com/stretchr/testify/require"
)

func TestTagFixedValues(t *testing.T) {
	source := struct {
		Fallback string `slack:"text"`
		Title    string `slack:"header;level=2"`
		Summary  string `slack:"section;expand;format=mrkdwn"`
		Body     struct {
			Steps struct {
				Items []string `slack:"rich_text_section"`
			} `slack:"rich_text_list;style=ordered"`
		} `slack:"rich_text"`
	}{Fallback: "Report", Title: "Report", Summary: "*done*"}
	source.Body.Steps.Items = []string{"build", "deploy"}
	message, err := ToMessage(source, conversion.WithWarningAsError())
	require.NoError(t, err)
	data, err := json.Marshal(message.Blocks)
	require.NoError(t, err)
	require.JSONEq(t, `[
		{"type":"header","text":{"type":"plain_text","text":"Report"},"level":2},
		{"type":"section","text":{"type":"mrkdwn","text":"*done*","verbatim":true},"expand":true},
		{"type":"rich_text","elements":[{"type":"rich_text_list","style":"ordered","elements":[
			{"type":"rich_text_section","elements":[{"type":"text","text":"build"}]},
			{"type":"rich_text_section","elements":[{"type":"text","text":"deploy"}]}
		]}]}
	]`, string(data))
}
