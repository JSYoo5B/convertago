package convertago_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago"
)

func ExampleToKakaoworkMessage() {
	source := struct {
		Preview string `kakaowork:"preview"`
		Title   string `kakaowork:"header"`
		Prefix  string `kakaowork:"text;group=body"`
		Name    string `kakaowork:"text;group=body;style=bold"`
		Image   string `kakaowork:"image_link"`
	}{"새 알림", "알림", "안녕하세요 ", "홍길동", "https://example.com/image.png"}
	message, err := convertago.ToKakaoworkMessage(source)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output: {"text":"새 알림","blocks":[{"type":"header","text":"알림","style":"white"},{"type":"text","text":"안녕하세요 홍길동","inlines":[{"type":"styled","text":"안녕하세요 "},{"type":"styled","text":"홍길동","bold":true}]},{"type":"image_link","url":"https://example.com/image.png"}]}
}
