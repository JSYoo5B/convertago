package sample

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/JSYoo5B/convertago"
	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
)

type reflectionNotice Notice
type reflectionVerboseNotice VerboseNotice

func TestGeneratedAndReflectionAgree(t *testing.T) {
	zero := 0
	base := Notice{
		Title: "Notice", Between: "between", Name: "Jane", Image: &Picture{"https://example.com/a.png", "photo"},
		Items: []Item{{"first ", "A"}, {"second ", "B"}}, Numbers: [2]int{0, 42}, Flags: []bool{false, true},
		Fraction: 0.1, Pointer: &zero, Label: "label", Caption: Caption("caption"), Time: time.Unix(0, 0).UTC(),
		Optional: "skip this unsupported header style",
		Alert:    Alert{"urgent"}, Wide: "wide",
		Ignored: map[string]string{"no tag": "ignore"},
	}
	for _, test := range []struct {
		name   string
		change func(*Notice)
	}{
		{"normal", func(*Notice) {}},
		{"nil pointers", func(n *Notice) { n.Image, n.Pointer = nil, nil }},
		{"without fixed fields", func(n *Notice) { n.Wide = "" }},
		{"omitted optional feature", func(n *Notice) { n.Optional = "" }},
		{"invalid URL", func(n *Notice) { n.Image = &Picture{"relative", "photo"} }},
		{"missing alt", func(n *Notice) { n.Image = &Picture{"https://example.com/a.png", ""} }},
		{"marshal error", func(n *Notice) { n.Label = "bad" }},
		{"empty marshal text", func(n *Notice) { n.Label = "" }},
		{"typed nil interface", func(n *Notice) { var caption *Caption; n.Caption = caption }},
		{"dynamic struct", func(n *Notice) {
			n.Dynamic = struct {
				Text string `kakaowork:"text" slack:"rich_text" googlechat:"textParagraph"`
			}{"dynamic"}
		}},
		{"dynamic typo", func(n *Notice) {
			n.Dynamic = struct {
				Text string `kakaowork:"typo" slack:"typo" googlechat:"typo"`
			}{}
		}},
		{"dynamic list typo", func(n *Notice) {
			n.Dynamic = []any{struct {
				Text string `kakaowork:"typo" slack:"typo" googlechat:"typo"`
			}{}}
		}},
		{"dynamic map", func(n *Notice) { n.Dynamic = map[string]string{"bad": "input"} }},
		{"dynamic slice cycle", func(n *Notice) {
			items := []any{nil}
			items[0] = items
			n.Dynamic = items
		}},
		{"pointer cycle", func(n *Notice) { n.Next = n }},
	} {
		t.Run(test.name, func(t *testing.T) {
			notice := base
			test.change(&notice)
			compare(t, notice, reflectionNotice(notice))
		})
	}
	compare(t, VerboseNotice{Text: "root tags"}, reflectionVerboseNotice{Text: "root tags"})
}

func TestConcurrentGeneratedConversions(t *testing.T) {
	notice := Notice{Title: "Notice", Between: "body", Name: "Jane", Items: []Item{{"hello ", "A"}}}
	var work sync.WaitGroup
	for i := 0; i < 16; i++ {
		work.Add(1)
		go func() {
			defer work.Done()
			for j := 0; j < 10; j++ {
				if _, err := convertago.ToSlackMessage(&notice); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	work.Wait()
}

func compare(t *testing.T, generated, reflected any) {
	t.Helper()
	v := validator.New(validator.WithRequiredStructEnabled())
	kakaowork.RegisterValidation(v)
	slack.RegisterValidation(v)
	googlechat.RegisterValidation(v)
	converters := []struct {
		name string
		fn   func(any, ...convertago.Option) (any, error)
	}{
		{"kakaowork", func(v any, o ...convertago.Option) (any, error) { return convertago.ToKakaoworkMessage(v, o...) }},
		{"slack", func(v any, o ...convertago.Option) (any, error) { return convertago.ToSlackMessage(v, o...) }},
		{"googlechat", func(v any, o ...convertago.Option) (any, error) { return convertago.ToGoogleChatMessage(v, o...) }},
	}
	for _, converter := range converters {
		for _, warningAsError := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/warningAsError=%t", converter.name, warningAsError), func(t *testing.T) {
				var gotDiagnostics, wantDiagnostics []convertago.Diagnostic
				gotOptions := []convertago.Option{convertago.WithDiagnostics(func(d convertago.Diagnostic) { gotDiagnostics = append(gotDiagnostics, d) })}
				wantOptions := []convertago.Option{convertago.WithDiagnostics(func(d convertago.Diagnostic) { wantDiagnostics = append(wantDiagnostics, d) })}
				if warningAsError {
					gotOptions = append(gotOptions, convertago.WithWarningAsError())
					wantOptions = append(wantOptions, convertago.WithWarningAsError())
				}
				got, gotErr := converter.fn(generated, gotOptions...)
				want, wantErr := converter.fn(reflected, wantOptions...)
				if !reflect.DeepEqual(gotErr, wantErr) || !reflect.DeepEqual(gotDiagnostics, wantDiagnostics) {
					t.Fatalf("errors: generated=%#v reflected=%#v\ndiagnostics: generated=%#v reflected=%#v", gotErr, wantErr, gotDiagnostics, wantDiagnostics)
				}
				if gotErr == nil {
					for _, message := range []any{got, want} {
						if err := v.Struct(message); err != nil {
							t.Fatalf("converted output validation: %v", err)
						}
					}
				}
				gotJSON, err := json.Marshal(got)
				if err != nil {
					t.Fatal(err)
				}
				wantJSON, err := json.Marshal(want)
				if err != nil {
					t.Fatal(err)
				}
				if string(gotJSON) != string(wantJSON) {
					t.Fatalf("JSON differs:\ngenerated=%s\nreflected=%s", gotJSON, wantJSON)
				}
			})
		}
	}
}

type reflectionLayoutNotice LayoutNotice

func TestGeneratedNestedLayoutsAgree(t *testing.T) {
	zero := 0
	base := LayoutNotice{
		Header:  "Notice",
		Body:    RichInput{"See ", LinkInput{"docs", "https://example.com"}, " now"},
		Buttons: ButtonRow{[]ButtonInput{{"Confirm", ActionInput{"confirm", "yes"}}, {"Cancel", ActionInput{"cancel", "no"}}}},
		Rich:    SlackRich{Before: "Before", List: SlackList{Style: "ordered", Border: &zero, Items: []string{"one", "two"}}, Quote: "Quote", After: "After"},
		Cards: []GoogleWrappedCard{
			{"first", GoogleCard{"First", []GoogleSection{{"Summary", GoogleColumns{[]GoogleColumn{{"left"}, {"right"}}}}}}},
			{"second", GoogleCard{"Second", []GoogleSection{{"Details", GoogleColumns{[]GoogleColumn{{"detail"}}}}}}},
		},
	}
	// Confirm the normal case succeeds, rather than comparing two failing paths.
	if _, err := convertago.ToKakaoworkMessage(base); err != nil {
		t.Fatal(err)
	}
	if _, err := convertago.ToSlackMessage(base); err != nil {
		t.Fatal(err)
	}
	if _, err := convertago.ToGoogleChatMessage(base); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name   string
		change func(*LayoutNotice)
	}{
		{"normal", func(*LayoutNotice) {}},
		{"invalid link", func(n *LayoutNotice) { n.Body.Link.URL = "relative" }},
		{"empty action row", func(n *LayoutNotice) { n.Buttons.Buttons = nil }},
		{"too many Kakao buttons", func(n *LayoutNotice) { n.Buttons.Buttons = append(n.Buttons.Buttons, n.Buttons.Buttons...) }},
		{"invalid list style", func(n *LayoutNotice) { n.Rich.List.Style = "unknown" }},
		{"missing list border", func(n *LayoutNotice) { n.Rich.List.Border = nil }},
		{"invalid column count", func(n *LayoutNotice) {
			n.Cards[0].Card.Sections[0].Columns.Columns = append(n.Cards[0].Card.Sections[0].Columns.Columns, GoogleColumn{"third"})
		}},
		{"duplicate card ID", func(n *LayoutNotice) { n.Cards[1].ID = "first" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			notice := base
			notice.Cards = append([]GoogleWrappedCard(nil), base.Cards...)
			notice.Cards[0].Card.Sections = append([]GoogleSection(nil), base.Cards[0].Card.Sections...)
			test.change(&notice)
			compare(t, notice, reflectionLayoutNotice(notice))
		})
	}
}
