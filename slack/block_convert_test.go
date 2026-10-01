package slack

import (
	"strings"
	"testing"
)

func TestSectionFieldsAndAccessory(t *testing.T) {
	type button struct {
		Text string `slack:"part"`
		ID   string `slack:"part;slot=action_id"`
	}
	source := struct {
		Section struct {
			Fields []string `slack:"plain_text;slot=fields"`
			Button button   `slack:"button;slot=accessory"`
			Expand bool     `slack:"part;slot=expand"`
		} `slack:"section"`
	}{}
	source.Section.Fields = []string{"one", "two"}
	source.Section.Button = button{"open", "action"}
	source.Section.Expand = true
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	section := message.Blocks[0].(SectionBlock)
	if len(section.Fields) != 2 || section.Text != nil || !section.Expand || section.Accessory.(ButtonElement).ActionID != "action" {
		t.Fatalf("section=%#v", section)
	}
	source.Section.Fields = make([]string, 11)
	for i := range source.Section.Fields {
		source.Section.Fields[i] = "field"
	}
	if _, err := ToMessage(source); err == nil {
		t.Fatal("accepted too many section fields")
	}
}

func TestSlackFileSourceValidation(t *testing.T) {
	type file struct {
		ID  string `slack:"part;slot=id;omitempty"`
		URL string `slack:"part;slot=url;omitempty"`
	}
	source := struct {
		Image struct {
			File file   `slack:"slack_file"`
			Alt  string `slack:"part;slot=alt"`
		} `slack:"image"`
	}{}
	source.Image.File.ID = "F123"
	source.Image.Alt = "photo"
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	if message.Blocks[0].(ImageBlock).SlackFile.ID != "F123" {
		t.Fatal("file ID lost")
	}
	source.Image.File.URL = "https://example.com/image.png"
	if _, err := ToMessage(source); err == nil {
		t.Fatal("ambiguous file source accepted")
	}
}

func TestMarkdownCumulativeLimit(t *testing.T) {
	source := struct {
		Text []string `slack:"markdown"`
	}{[]string{strings.Repeat("x", 6000), strings.Repeat("x", 6000)}}
	if _, err := ToMessage(source); err != nil {
		t.Fatal(err)
	}
	source.Text[1] += "x"
	if _, err := ToMessage(source); err == nil {
		t.Fatal("accepted more than 12000 total characters")
	}
}

func TestImageElementRejectsBlockProperties(t *testing.T) {
	type image struct {
		URL   string `slack:"part;slot=url"`
		Alt   string `slack:"part;slot=alt"`
		Title string `slack:"part;slot=title"`
	}
	source := struct {
		Context struct {
			Image image `slack:"image"`
		} `slack:"context"`
	}{}
	source.Context.Image = image{"https://example.com/image.png", "image", "title"}
	if _, err := ToMessage(source); err == nil {
		t.Fatal("silently dropped image element title")
	}
}

func TestConfirmationTextAndExplicitFalse(t *testing.T) {
	type plain struct {
		Text  string `slack:"part"`
		Emoji bool   `slack:"part;slot=emoji"`
	}
	type confirmation struct {
		Title   string `slack:"part;slot=title"`
		Text    string `slack:"part;format=mrkdwn"`
		Confirm string `slack:"part;slot=confirm"`
		Deny    string `slack:"part;slot=deny"`
	}
	type button struct {
		Text    plain        `slack:"plain_text;slot=text"`
		Confirm confirmation `slack:"confirm"`
		Users   []string     `slack:"part;slot=visible_to_user_ids"`
	}
	source := struct {
		Actions struct {
			Button button `slack:"button"`
		} `slack:"actions"`
	}{}
	source.Actions.Button = button{plain{"Send", false}, confirmation{"Confirm", "*Send?*", "Yes", "No"}, []string{"U1", "U2"}}
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	native := message.Blocks[0].(ActionsBlock).Elements[0].(ButtonElement)
	if native.Text.Emoji == nil || *native.Text.Emoji || native.Confirm.Text.(MrkdwnTextObject).Text != "*Send?*" || len(native.VisibleToUserIDs) != 2 {
		t.Fatalf("button=%#v", native)
	}
	source.Actions.Button.Confirm.Confirm = strings.Repeat("x", 31)
	if _, err := ToMessage(source); err == nil {
		t.Fatal("confirmation label too long")
	}
}

func TestVideoNativeProperties(t *testing.T) {
	source := struct {
		Video struct {
			Title       string `slack:"part;slot=title"`
			URL         string `slack:"part;slot=video_url"`
			Thumbnail   string `slack:"part;slot=thumbnail_url"`
			Alt         string `slack:"part;slot=alt"`
			Description string `slack:"part;slot=description"`
			Author      string `slack:"part;slot=author_name"`
		} `slack:"video"`
	}{}
	source.Video.Title = "Video"
	source.Video.URL = "https://example.com/embed"
	source.Video.Thumbnail = "https://example.com/image.png"
	source.Video.Alt = "Video preview"
	source.Video.Description = "Description"
	source.Video.Author = "Author"
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	native := message.Blocks[0].(VideoBlock)
	if native.Description.Text != "Description" || native.AuthorName != "Author" {
		t.Fatalf("video=%#v", native)
	}
	source.Video.URL = "http://example.com/embed"
	if _, err := ToMessage(source); err == nil {
		t.Fatal("non-HTTPS embed accepted")
	}
}
