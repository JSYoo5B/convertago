package slack_test

import (
	"encoding/json"
	"fmt"

	"github.com/JSYoo5B/convertago/slack"
)

func ExampleToMessage_blocks() {
	source := struct {
		Fallback string `slack:"text"`
		Title    struct {
			Text  string `slack:"part"`
			Level int    `slack:"part;slot=level"`
		} `slack:"header"`
		Chart struct {
			URL   string `slack:"part;slot=url"`
			Alt   string `slack:"part;slot=alt"`
			Title string `slack:"part;slot=title"`
		} `slack:"image"`
		Divider struct{} `slack:"divider"`
		Summary string   `slack:"markdown;format=markdown"`
		Demo    struct {
			URL       string `slack:"part;slot=video_url"`
			Title     string `slack:"part;slot=title"`
			Alt       string `slack:"part;slot=alt"`
			Thumbnail string `slack:"part;slot=thumbnail_url"`
		} `slack:"video"`
	}{Fallback: "Weekly report", Summary: "**12** deploys this week"}
	source.Title.Text, source.Title.Level = "Weekly report", 2
	source.Chart.URL, source.Chart.Alt, source.Chart.Title = "https://example.com/chart.png", "Deploys per day", "Deploys"
	source.Demo.URL, source.Demo.Title = "https://example.com/embed/demo", "Release demo"
	source.Demo.Alt, source.Demo.Thumbnail = "Release demo", "https://example.com/demo.png"
	message, err := slack.ToMessage(source)
	if err != nil {
		panic(err)
	}
	data, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output:
	// {
	//   "text": "Weekly report",
	//   "blocks": [
	//     {
	//       "type": "header",
	//       "text": {
	//         "type": "plain_text",
	//         "text": "Weekly report"
	//       },
	//       "level": 2
	//     },
	//     {
	//       "type": "image",
	//       "image_url": "https://example.com/chart.png",
	//       "alt_text": "Deploys per day",
	//       "title": {
	//         "type": "plain_text",
	//         "text": "Deploys"
	//       }
	//     },
	//     {
	//       "type": "divider"
	//     },
	//     {
	//       "type": "markdown",
	//       "text": "**12** deploys this week"
	//     },
	//     {
	//       "type": "video",
	//       "title": {
	//         "type": "plain_text",
	//         "text": "Release demo"
	//       },
	//       "video_url": "https://example.com/embed/demo",
	//       "alt_text": "Release demo",
	//       "thumbnail_url": "https://example.com/demo.png"
	//     }
	//   ]
	// }
}
