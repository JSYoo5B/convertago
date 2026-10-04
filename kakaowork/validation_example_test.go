package kakaowork_test

import (
	"fmt"

	"github.com/JSYoo5B/convertago"
	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/go-playground/validator/v10"
)

func ExampleValidate() {
	message := kakaowork.Message{Preview: "알림", Blocks: []kakaowork.BubbleBlock{
		kakaowork.TextBlock{Text: "배포가 완료되었습니다."},
		kakaowork.HeaderBlock{Text: "알림"},
	}}
	err := kakaowork.Validate(message, convertago.WithDiagnostics(func(d convertago.Diagnostic) {
		fmt.Println(d.Severity, d.Code, d.Path)
	}))
	fmt.Println(err == nil)

	err = kakaowork.Validate(message, convertago.WithWarningAsError())
	fmt.Println(err)

	// Output:
	// warning header.position $.blocks[1]
	// true
	// convertago: kakaowork $.blocks[1]: a header block must be the first and only header [warning header.position]
}

func ExampleRegisterValidation() {
	v := validator.New()
	kakaowork.RegisterValidation(v)

	message := kakaowork.Message{Preview: "알림", Blocks: []kakaowork.BubbleBlock{
		kakaowork.HeaderBlock{Text: "알림", Style: kakaowork.HeaderStyleBlue},
	}}
	fmt.Println(v.Struct(message) == nil)

	message.Blocks[0] = kakaowork.HeaderBlock{Text: "알림", Style: "green"}
	fmt.Println(v.Struct(message) == nil)

	// Output:
	// true
	// false
}
