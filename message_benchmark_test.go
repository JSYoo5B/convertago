package convertago_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"runtime"
	"testing"

	"github.com/JSYoo5B/convertago"
	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/JSYoo5B/convertago/internal/benchmarksource"
	"github.com/JSYoo5B/convertago/internal/conversion"
	"github.com/JSYoo5B/convertago/kakaowork"
	"github.com/JSYoo5B/convertago/slack"
	"github.com/go-playground/validator/v10"
)

// Defined types retain the generated fixtures' field layouts and tags while
// dropping their methods, so these inputs select the reflection reader.
type benchmarkReflectionFlatMessage benchmarksource.FlatMessage
type benchmarkReflectionNestedMessage benchmarksource.NestedMessage
type benchmarkReflectionDynamicMessage benchmarksource.DynamicMessage

type benchmarkMessageCase struct {
	name       string
	reflection any
	generated  any
}

func (tc benchmarkMessageCase) paths() []struct {
	name  string
	input any
} {
	return []struct {
		name  string
		input any
	}{
		{name: "Reflection", input: tc.reflection},
		{name: "Generated", input: tc.generated},
	}
}

func benchmarkMessageCases() []benchmarkMessageCase {
	flat := &benchmarksource.FlatMessage{
		Title: "Notice", Prefix: "Hello ", Name: "Jane", Count: 123,
	}
	footer := "Complete"
	nested := &benchmarksource.NestedMessage{
		Title: "Report",
		Picture: &benchmarksource.Picture{
			URL: "https://example.com/image.png", Alt: "Report image",
		},
		Items: []benchmarksource.Line{
			{Prefix: "Hello ", Name: "Jane"},
			{Prefix: "Hello ", Name: "John"},
			{Prefix: "Hello ", Name: "Alex"},
		},
		Footer: &footer,
	}
	// Both paths use the same actual interface value. The generated wrapper also
	// uses cached reflection for Content rather than statically reading Nested.
	dynamic := &benchmarksource.DynamicMessage{Content: nested}
	return []benchmarkMessageCase{
		{name: "Flat", reflection: (*benchmarkReflectionFlatMessage)(flat), generated: flat},
		{name: "Nested", reflection: (*benchmarkReflectionNestedMessage)(nested), generated: nested},
		{name: "Dynamic", reflection: (*benchmarkReflectionDynamicMessage)(dynamic), generated: dynamic},
	}
}

func BenchmarkToKakaoworkMessage(b *testing.B) {
	benchmarkToMessage(b, convertago.ToKakaoworkMessage)
}

func BenchmarkToSlackMessage(b *testing.B) {
	benchmarkToMessage(b, convertago.ToSlackMessage)
}

func BenchmarkToGoogleChatMessage(b *testing.B) {
	benchmarkToMessage(b, convertago.ToGoogleChatMessage)
}

// BenchmarkSourceFields compares the two readers through the same entry point,
// excluding node assembly and native rendering. The initial call warms any
// reflection plans, including the generated dynamic field's fallback cache.
func BenchmarkSourceFields(b *testing.B) {
	for _, platform := range []string{"kakaowork", "slack", "googlechat"} {
		b.Run(platform, func(b *testing.B) {
			benchmarkInputs(b, func(input any) ([]conversion.Field, error) {
				return conversion.Fields(input, platform)
			})
		})
	}
}

// benchmarkToMessage includes source reading, assembly, validation, and native
// rendering for both paths. Caller-side JSON serialization is excluded, while
// serialization inside a converter is included. AST analysis and generation
// happen before compilation and are excluded from these runtime benchmarks.
func benchmarkToMessage[T any](b *testing.B, convert func(any, ...convertago.Option) (T, error)) {
	benchmarkInputs(b, func(input any) (T, error) { return convert(input) })
}

func benchmarkInputs[T any](b *testing.B, read func(any) (T, error)) {
	for _, tc := range benchmarkMessageCases() {
		b.Run(tc.name, func(b *testing.B) {
			for _, path := range tc.paths() {
				b.Run(path.name, func(b *testing.B) {
					result, err := read(path.input)
					if err != nil {
						b.Fatal(err)
					}
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						result, err = read(path.input)
						if err != nil {
							b.Fatal(err)
						}
					}
					runtime.KeepAlive(result)
				})
			}
		})
	}
}

// Check both benchmark paths before relying on their performance measurements.
// JSON equality ensures the generated reader does the same useful conversion.
func TestBenchmarkMessageSources(t *testing.T) {
	v := validator.New(validator.WithRequiredStructEnabled())
	kakaowork.RegisterValidation(v)
	slack.RegisterValidation(v)
	googlechat.RegisterValidation(v)
	converters := []struct {
		platform string
		convert  func(any) (any, error)
	}{
		{"kakaowork", func(input any) (any, error) { return convertago.ToKakaoworkMessage(input) }},
		{"slack", func(input any) (any, error) { return convertago.ToSlackMessage(input) }},
		{"googlechat", func(input any) (any, error) { return convertago.ToGoogleChatMessage(input) }},
	}
	for _, tc := range benchmarkMessageCases() {
		t.Run(tc.name, func(t *testing.T) {
			if _, ok := tc.generated.(conversion.Generated); !ok {
				t.Fatal("generated fixture does not implement generated accessors")
			}
			if _, ok := tc.reflection.(conversion.Generated); ok {
				t.Fatal("reflection fixture unexpectedly implements generated accessors")
			}
			for _, converter := range converters {
				t.Run(converter.platform, func(t *testing.T) {
					fields, err := conversion.Fields(tc.reflection, converter.platform)
					if err != nil {
						t.Fatal(err)
					}
					generated, err := conversion.Fields(tc.generated, converter.platform)
					if err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(generated, fields) {
						t.Fatal("source fields differ between readers")
					}
					var results [][]byte
					for _, path := range tc.paths() {
						message, err := converter.convert(path.input)
						if err != nil {
							t.Fatal(err)
						}
						if err := v.Struct(message); err != nil {
							t.Fatalf("%s output validation: %v", path.name, err)
						}
						data, err := json.Marshal(message)
						if err != nil {
							t.Fatal(err)
						}
						results = append(results, data)
					}
					if !bytes.Equal(results[0], results[1]) {
						t.Fatalf("message JSON differs:\nreflection=%s\ngenerated=%s", results[0], results[1])
					}
				})
			}
		})
	}
}
