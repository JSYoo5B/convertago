package googlechat_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestMessage_ValidateCardIDs(t *testing.T) {
	v := validator.New()
	googlechat.RegisterValidation(v)
	card := googlechat.Card{Header: &googlechat.CardHeader{Title: "Notice"}}
	message := googlechat.Message{CardsV2: []googlechat.CardWithID{{Card: card}}}
	require.NoError(t, v.Struct(message))
	message.CardsV2 = append(message.CardsV2, googlechat.CardWithID{CardID: "second", Card: card})
	require.Error(t, v.Struct(message))
	message.CardsV2[0].CardID = "first"
	require.NoError(t, v.Struct(message))
	message.CardsV2[1].CardID = "first"
	err := v.Struct(message)
	require.ErrorAs(t, err, new(validator.ValidationErrors))
	require.Equal(t, "Message.CardsV2[1].CardID", err.(validator.ValidationErrors)[0].Namespace())
}

func TestMessage_ValidatePayloadBytes(t *testing.T) {
	v := validator.New()
	googlechat.RegisterValidation(v)
	paragraph := &googlechat.TextParagraph{}
	message := googlechat.Message{CardsV2: []googlechat.CardWithID{{Card: googlechat.Card{
		Sections: []googlechat.Section{{Widgets: []googlechat.Widget{{Content: paragraph}}}},
	}}}}
	data, err := json.Marshal(message.CardsV2)
	require.NoError(t, err)
	available := 32*1024 - len(data)
	paragraph.Text = strings.Repeat("가", available/3) + strings.Repeat("a", available%3)
	data, err = json.Marshal(message.CardsV2)
	require.NoError(t, err)
	require.Len(t, data, 32*1024)
	require.NoError(t, v.Struct(message))

	paragraph.Text += "a"
	err = v.Struct(message)
	require.ErrorAs(t, err, new(validator.ValidationErrors))
	require.Equal(t, "max_bytes", err.(validator.ValidationErrors)[0].Tag())

	paragraph.Text = strings.Repeat("<", 6000)
	require.Error(t, v.Struct(message))
}
