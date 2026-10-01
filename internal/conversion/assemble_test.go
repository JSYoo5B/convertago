package conversion

import (
	"errors"
	"reflect"
	"sync"
	"testing"
)

func TestGroupAnchorSurvivesOmission(t *testing.T) {
	source := struct {
		First  string `test:"text;group=body;omitempty"`
		Middle string `test:"header"`
		Last   string `test:"text;group=body;style=bold"`
	}{Middle: "heading", Last: "name"}
	nodes, err := Prepare(source, "test", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 2 || nodes[0].Path != "$.First" || nodes[0].Text("text") != "name" || nodes[1].Role != "header" || !reflect.DeepEqual(nodes[0].Parts[0].Style, []string{"bold"}) {
		t.Fatalf("unexpected nodes: %#v", nodes)
	}
}

func TestNestedListsAndIndependentScopes(t *testing.T) {
	type line struct {
		Start string `test:"part;group=line"`
		End   string `test:"part;group=line;style=bold"`
	}
	type item struct {
		Title string `test:"header"`
		Body  line   `test:"text"`
	}
	source := struct {
		Items []item `test:"flatten"`
	}{[]item{{"one", line{"hello ", "A"}}, {"two", line{"hello ", "B"}}}}
	nodes, err := Prepare(source, "test", nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"one", "hello A", "two", "hello B"}
	var got []string
	for _, node := range nodes {
		got = append(got, node.Text("text"))
	}
	if !reflect.DeepEqual(got, want) || nodes[3].Parts[1].Path != "$.Items[1].Body.End" {
		t.Fatalf("nodes = %#v", nodes)
	}
}

func TestScalarOmissionAndPresence(t *testing.T) {
	zero := 0
	source := struct {
		Zero int    `test:"text"`
		No   bool   `test:"text"`
		Skip int    `test:"text;omitempty"`
		Ptr  *int   `test:"text;omitempty"`
		Nil  *int   `test:"text"`
		List []bool `test:"text"`
	}{Ptr: &zero, List: []bool{false, true}}
	nodes, err := Prepare(source, "test", nil)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, node := range nodes {
		got = append(got, node.Text("text"))
	}
	if !reflect.DeepEqual(got, []string{"0", "false", "0", "false", "true"}) {
		t.Fatalf("texts = %#v", got)
	}
}

func TestInputErrorsHaveSourcePaths(t *testing.T) {
	tests := []struct {
		name   string
		source any
		path   string
		code   string
	}{
		{"missing slot", struct {
			URL string `test:"image"`
		}{"https://example.com/a.png"}, "$.URL", "missing_input"},
		{"duplicate slot", struct {
			A string `test:"image;group=photo"`
			B string `test:"image;group=photo"`
		}{"a", "b"}, "$.B", "duplicate_input"},
		{"group conflict while empty", struct {
			A string `test:"text;group=x;omitempty"`
			B string `test:"header;group=x;omitempty"`
		}{}, "$.B", "group_conflict"},
		{"bad tag beneath nil", struct {
			Body *struct {
				Text string `test:"prat"`
			} `test:"text"`
		}{}, "$.Body.Text", "invalid_tag"},
		{"bad slot beneath nil", struct {
			Body *struct {
				Text string `test:"part;slot=typo"`
			} `test:"text"`
		}{}, "$.Body.Text", "invalid_tag"},
		{"map even when empty", struct {
			Data map[string]string `test:"text;omitempty"`
		}{}, "$.Data", "invalid_source"},
		{"scalar flatten", struct {
			Text string `test:"flatten;omitempty"`
		}{}, "$.Text", "invalid_source"},
		{"part at root", struct {
			Text string `test:"part"`
		}{}, "$.Text", "invalid_tag"},
		{"dynamic bad tag", struct {
			Data any `test:"flatten"`
		}{struct {
			Text string `test:"typo"`
		}{}}, "$.Data.Text", "invalid_tag"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Prepare(test.source, "test", nil)
			var diagnostic Diagnostic
			if !errors.As(err, &diagnostic) || diagnostic.Platform != "test" || diagnostic.Path != test.path || diagnostic.Code != test.code {
				t.Fatalf("error = %#v, want %s at %s", err, test.code, test.path)
			}
		})
	}
}

func TestOptionalAndStrict(t *testing.T) {
	source := struct {
		Skip string `test:"button;optional"`
		Keep string `test:"text"`
	}{"click", "body"}
	var diagnostics []Diagnostic
	nodes, err := Prepare(source, "test", []Option{WithDiagnostics(func(d Diagnostic) { diagnostics = append(diagnostics, d) })})
	if err != nil || len(nodes) != 1 || len(diagnostics) != 1 || diagnostics[0].Path != "$.Skip" {
		t.Fatalf("nodes = %#v, diagnostics = %#v, error = %v", nodes, diagnostics, err)
	}
	if _, err := Prepare(source, "test", []Option{WithStrict()}); err == nil {
		t.Fatal("strict must reject optional unsupported features")
	}
	source2 := struct {
		Skip string `test:"text;style=underline;optional"`
	}{"body"}
	if nodes, err := Prepare(source2, "test", nil); err != nil || len(nodes) != 0 {
		t.Fatalf("nodes = %#v, error = %v", nodes, err)
	}
}

func TestStrictPreservesIntentionalOmissions(t *testing.T) {
	source := struct {
		Missing *string `test:"button;optional"`
		Empty   string  `test:"button;optional;omitempty"`
		Ignored string  `test:"-"`
		NoTag   string
	}{}
	var diagnostics []Diagnostic
	nodes, err := Prepare(source, "test", []Option{WithStrict(), WithDiagnostics(func(d Diagnostic) { diagnostics = append(diagnostics, d) })})
	if err != nil || len(nodes) != 0 || len(diagnostics) != 0 {
		t.Fatalf("nodes = %#v, diagnostics = %#v, error = %v", nodes, diagnostics, err)
	}
}

func TestDynamicListDoesNotHideMalformedTags(t *testing.T) {
	source := struct {
		Items any `test:"flatten"`
	}{[]any{struct {
		Text string `test:"typo"`
	}{}}}
	_, err := Prepare(source, "test", nil)
	var diagnostic Diagnostic
	if !errors.As(err, &diagnostic) || diagnostic.Path != "$.Items[0].Text" || diagnostic.Code != "invalid_tag" {
		t.Fatalf("error = %#v", err)
	}
}

type recursiveSource struct {
	Text string           `test:"text"`
	Next *recursiveSource `test:"flatten"`
}

func TestCyclesAndSharedPointers(t *testing.T) {
	source := &recursiveSource{Text: "body"}
	source.Next = source
	if _, err := Prepare(source, "test", nil); err == nil {
		t.Fatal("expected cycle error")
	}
	source.Next = nil
	root := struct {
		Items []*recursiveSource `test:"flatten"`
	}{[]*recursiveSource{source, source}}
	if nodes, err := Prepare(root, "test", nil); err != nil || len(nodes) != 2 {
		t.Fatalf("nodes = %#v, error = %v", nodes, err)
	}
}

func TestCachedPlansAreSafeForConcurrentConversion(t *testing.T) {
	source := struct {
		Text string `test:"text"`
	}{"hello"}
	var tasks sync.WaitGroup
	for i := 0; i < 16; i++ {
		tasks.Add(1)
		go func() {
			defer tasks.Done()
			for j := 0; j < 10; j++ {
				if nodes, err := Prepare(source, "test", nil); err != nil || nodes[0].Text("text") != "hello" {
					t.Errorf("nodes = %#v, error = %v", nodes, err)
				}
			}
		}()
	}
	tasks.Wait()
}
