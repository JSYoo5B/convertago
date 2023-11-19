package kakaowork_test

import (
	"encoding/json"
	"fmt"
	"github.com/JSYoo5B/convertago/kakaowork"
)

func ExampleDividerBlock() {
	dividerBlock := kakaowork.DividerBlock{}

	jsonBytes, err := json.MarshalIndent(dividerBlock, "", "    ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonBytes))
	// Output: {
	//     "type": "divider"
	// }
}
