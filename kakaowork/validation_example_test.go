package kakaowork_test

import (
	"fmt"

	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/go-playground/validator/v10"
)

func ExampleRegisterValidation() {
	v := validator.New()
	kakaowork.RegisterValidation(v)

	message := kakaowork.Message{Blocks: []kakaowork.BubbleBlock{
		kakaowork.HeaderBlock{Text: "알림"},
		kakaowork.TextBlock{Text: "배포가 완료되었습니다."},
	}}
	fmt.Println(v.Struct(message) == nil)

	message.Blocks[0], message.Blocks[1] = message.Blocks[1], message.Blocks[0]
	fmt.Println(v.Struct(message) == nil)

	// Output:
	// true
	// false
}
