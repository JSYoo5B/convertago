package kakaowork

import (
	"slices"
	"testing"

	"github.com/JSYoo5B/convertago/internal/conversion"
)

func TestConvertLayouts(t *testing.T) {
	type content struct {
		Text string `kakaowork:"part"`
	}
	type image struct {
		URL string `kakaowork:"part;slot=url"`
	}
	type details struct {
		Body   content `kakaowork:"text"`
		Term   string  `kakaowork:"part;slot=term"`
		Accent bool    `kakaowork:"part;slot=accent"`
	}
	type context struct {
		Body  content `kakaowork:"text"`
		Image image   `kakaowork:"image_link;slot=image"`
	}
	source := struct {
		Details details  `kakaowork:"description"`
		Divider struct{} `kakaowork:"divider"`
		Context context  `kakaowork:"context"`
	}{Details: details{content{"done"}, "state", true}, Context: context{content{"author"}, image{"https://example.com/image.png"}}}
	message, err := ToMessage(source)
	if err != nil || len(message.Blocks) != 3 {
		t.Fatalf("message=%#v err=%v", message, err)
	}
	if d := message.Blocks[0].(DescriptionBlock); !d.Accent || d.Content.Text != "done" {
		t.Fatalf("description=%#v", d)
	}
}

func TestTextInlineOrderAndLimit(t *testing.T) {
	type mention struct {
		Text string `kakaowork:"part"`
		ID   int    `kakaowork:"part;slot=user_id"`
	}
	source := struct {
		Body struct {
			Before string  `kakaowork:"part"`
			User   mention `kakaowork:"mention"`
			After  string  `kakaowork:"part"`
		} `kakaowork:"text"`
	}{}
	source.Body.Before = "hello "
	source.Body.User = mention{"user", 7}
	source.Body.After = "!"
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	block := message.Blocks[0].(TextBlock)
	if block.Text != "hello user!" || len(block.Inlines) != 3 {
		t.Fatalf("block=%#v", block)
	}
	source.Body.User.ID = 0
	var codes []string
	if _, err := ToMessage(source, conversion.WithDiagnostics(func(d conversion.Diagnostic) { codes = append(codes, d.Code) })); err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(codes, ruleMentionUserID.ID) {
		t.Fatalf("diagnostics=%v", codes)
	}
}

func TestEveryKakaoworkRoleAvailable(t *testing.T) {
	profile, _ := conversion.Lookup("kakaowork")
	for name, role := range profile.Roles {
		if role.Unavailable {
			t.Errorf("unimplemented role %s", name)
		}
	}
}

func TestKakaoworkActionVariants(t *testing.T) {
	type value struct {
		Name  string `kakaowork:"part;slot=name"`
		Value string `kakaowork:"part;slot=value"`
	}
	type exclusive struct {
		Default value `kakaowork:"open_system_browser"`
		Mobile  value `kakaowork:"open_external_app;slot=mobile"`
	}
	type actions struct {
		System    *value     `kakaowork:"open_system_browser"`
		App       *value     `kakaowork:"open_external_app"`
		Submit    *value     `kakaowork:"submit_action"`
		Modal     *value     `kakaowork:"call_modal"`
		Exclusive *exclusive `kakaowork:"exclusive"`
	}
	type button struct {
		Text    string  `kakaowork:"part"`
		Actions actions `kakaowork:"flatten"`
	}
	source := struct {
		Buttons []button `kakaowork:"button"`
	}{[]button{
		{"System", actions{System: &value{"open", "https://example.com"}}},
		{"App", actions{App: &value{"app", "exampleapp://open"}}},
		{"Submit", actions{Submit: &value{"confirm", "yes"}}},
		{"Modal", actions{Modal: &value{"modal", "form"}}},
		{"Exclusive", actions{Exclusive: &exclusive{value{"web", "https://example.com"}, value{"app", "exampleapp://open"}}}},
	}}
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"open_system_browser", "open_external_app", "submit_action", "call_modal", "exclusive"}
	for i, block := range message.Blocks {
		if block.(ButtonBlock).Action.ActionType() != want[i] {
			t.Fatalf("action %d=%#v", i, block)
		}
	}
	action := message.Blocks[4].(ButtonBlock).Action.(ExclusiveAction)
	if action.Mobile.ActionType() != "open_external_app" {
		t.Fatal("exclusive mobile override lost")
	}
}

func TestKakaoworkStyledColorAndSectionAction(t *testing.T) {
	type styled struct {
		Text  string `kakaowork:"part"`
		Color string `kakaowork:"part;slot=color"`
		Bold  bool   `kakaowork:"part;slot=bold"`
	}
	type text struct {
		Styled styled `kakaowork:"styled"`
	}
	type action struct {
		Value string `kakaowork:"part"`
	}
	source := struct {
		Section struct {
			Text   text   `kakaowork:"text"`
			Action action `kakaowork:"open_system_browser;slot=action"`
		} `kakaowork:"section"`
	}{}
	source.Section.Text.Styled = styled{"color", "red", true}
	source.Section.Action.Value = "https://example.com"
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	section := message.Blocks[0].(SectionBlock)
	inline := section.Content.Inlines[0].(InlineStyled)
	if !inline.Bold || inline.Color != InlineColorRed || section.Action == nil {
		t.Fatalf("section=%#v inline=%#v", section, inline)
	}
	source.Section.Text.Styled.Color = "unknown"
	if _, err := ToMessage(source); err == nil {
		t.Fatal("unknown inline color accepted")
	}
}
