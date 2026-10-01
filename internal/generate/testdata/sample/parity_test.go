package sample

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/JSYoo5B/convertago"
)

type reflectionNotice Notice
type reflectionVerboseNotice VerboseNotice

func TestGeneratedAndReflectionAgree(t *testing.T) {
	zero := 0
	base := Notice{
		Title: "Notice", Between: "between", Name: "Jane", Image: &Picture{"https://example.com/a.png", "photo"},
		Items: []Item{{"first ", "A"}, {"second ", "B"}}, Numbers: [2]int{0, 42}, Flags: []bool{false, true},
		Fraction: 0.1, Pointer: &zero, Label: "label", Caption: Caption("caption"), Time: time.Unix(0, 0).UTC(),
		Optional: "skip this unimplemented converter role",
		Ignored:  map[string]string{"no tag": "ignore"},
	}
	for _, test := range []struct {
		name   string
		change func(*Notice)
	}{
		{"normal", func(*Notice) {}},
		{"nil pointers", func(n *Notice) { n.Image, n.Pointer = nil, nil }},
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
	converters := []struct {
		name string
		fn   func(any, ...convertago.Option) (any, error)
	}{
		{"kakaowork", func(v any, o ...convertago.Option) (any, error) { return convertago.ToKakaoworkMessage(v, o...) }},
		{"slack", func(v any, o ...convertago.Option) (any, error) { return convertago.ToSlackMessage(v, o...) }},
		{"googlechat", func(v any, o ...convertago.Option) (any, error) { return convertago.ToGoogleChatMessage(v, o...) }},
	}
	for _, converter := range converters {
		for _, strict := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/strict=%t", converter.name, strict), func(t *testing.T) {
				var gotDiagnostics, wantDiagnostics []convertago.Diagnostic
				gotOptions := []convertago.Option{convertago.WithDiagnostics(func(d convertago.Diagnostic) { gotDiagnostics = append(gotDiagnostics, d) })}
				wantOptions := []convertago.Option{convertago.WithDiagnostics(func(d convertago.Diagnostic) { wantDiagnostics = append(wantDiagnostics, d) })}
				if strict {
					gotOptions = append(gotOptions, convertago.WithStrict())
					wantOptions = append(wantOptions, convertago.WithStrict())
				}
				got, gotErr := converter.fn(generated, gotOptions...)
				want, wantErr := converter.fn(reflected, wantOptions...)
				if !reflect.DeepEqual(gotErr, wantErr) || !reflect.DeepEqual(gotDiagnostics, wantDiagnostics) {
					t.Fatalf("errors: generated=%#v reflected=%#v\ndiagnostics: generated=%#v reflected=%#v", gotErr, wantErr, gotDiagnostics, wantDiagnostics)
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
