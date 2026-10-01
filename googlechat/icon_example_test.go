package googlechat_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/googlechat"
)

func ExampleIcon_MarshalJSON() {
	icon := googlechat.Icon{
		MaterialIcon: &googlechat.MaterialIcon{Name: "check_circle", Fill: true, Weight: 500},
		AltText:      "Passed",
	}
	jsonBytes, err := json.MarshalIndent(icon, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output: {
	//   "materialIcon": {
	//     "name": "check_circle",
	//     "fill": true,
	//     "weight": 500
	//   },
	//   "altText": "Passed"
	// }
}

func ExampleMaterialIcon() {
	value := googlechat.MaterialIcon{
		Name:   "check_circle",
		Fill:   true,
		Weight: 500,
		Grade:  200,
	}
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonBytes))
	// Output:
	// {
	//   "name": "check_circle",
	//   "fill": true,
	//   "weight": 500,
	//   "grade": 200
	// }
}
