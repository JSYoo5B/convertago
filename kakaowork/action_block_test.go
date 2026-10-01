package kakaowork_test

import (
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/go-playground/validator/v10"
)

func TestActionBlock_Validate(t *testing.T) {
	v := validator.New()
	kakaowork.RegisterValidation(v)
	button := kakaowork.ButtonBlock{Text: "확인", Action: kakaowork.SubmitAction{Name: "confirm"}}
	for _, count := range []int{0, 1, 2, 3, 4} {
		block := kakaowork.ActionBlock{}
		for i := 0; i < count; i++ {
			block.Elements = append(block.Elements, button)
		}
		message := kakaowork.Message{Blocks: []kakaowork.BubbleBlock{block}}
		err := v.Struct(message)
		if (err == nil) != (count == 2 || count == 3) {
			t.Fatalf("%d buttons: %v", count, err)
		}
	}
	button.Text = strings.Repeat("가", 21)
	block := kakaowork.ActionBlock{Elements: []kakaowork.ButtonBlock{button, button}}
	if err := v.Struct(block); err == nil {
		t.Fatal("nested button length was not validated")
	}
	block.Elements[0].Text = "확인"
	block.Elements[0].Action = kakaowork.OpenSystemBrowserAction{Value: "relative"}
	if err := v.Struct(block); err == nil {
		t.Fatal("nested button URL was not validated")
	}
}
