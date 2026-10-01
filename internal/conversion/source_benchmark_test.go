package conversion

import (
	"reflect"
	"runtime"
	"testing"
)

type benchmarkFlatSource struct {
	Title   string `test:"header"`
	Prefix  string `test:"text;group=body"`
	Name    string `test:"text;group=body;style=bold"`
	Count   int    `test:"text"`
	Enabled bool   `test:"text;omitempty"`
}

type benchmarkTextSource struct {
	Prefix string `test:"part;group=line"`
	Name   string `test:"part;group=line;style=bold"`
}

type benchmarkItemSource struct {
	Text  benchmarkTextSource `test:"text"`
	Count int                 `test:"text"`
}

type benchmarkNestedSource struct {
	Title  string                `test:"header"`
	Items  []benchmarkItemSource `test:"flatten"`
	Footer *string               `test:"text;omitempty"`
}

type benchmarkDynamicSource struct {
	Content any `test:"flatten"`
}

// BenchmarkSourcePlan isolates the cached type metadata from source values.
// Cached measures loadPlan hits. Uncached compiles and validates the same type
// on every iteration, excluding cache insertion and eviction. Both branches can
// run a fixed iteration count for process CPU comparisons without timer pauses.
func BenchmarkSourcePlan(b *testing.B) {
	for _, tc := range []struct {
		name   string
		typeOf reflect.Type
	}{
		{name: "Flat", typeOf: reflect.TypeOf(benchmarkFlatSource{})},
		{name: "Nested", typeOf: reflect.TypeOf(benchmarkNestedSource{})},
	} {
		b.Run(tc.name, func(b *testing.B) {
			key := planKey{tc.typeOf, testProfile.Platform}
			b.Cleanup(func() { sourcePlans.Delete(key) })
			b.Run("Cached", func(b *testing.B) {
				entry := loadPlan(tc.typeOf, testProfile)
				if entry.err != nil {
					b.Fatal(entry.err)
				}
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					entry = loadPlan(tc.typeOf, testProfile)
					if entry.err != nil {
						b.Fatal(entry.err)
					}
				}
				runtime.KeepAlive(entry)
			})
			b.Run("Uncached", func(b *testing.B) {
				var entry cachedPlan
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					plan := compileObjectPlan(tc.typeOf, testProfile)
					entry = cachedPlan{plan: plan, err: ValidateShape(testProfile, plan.shape)}
					if entry.err != nil {
						b.Fatal(entry.err)
					}
				}
				runtime.KeepAlive(entry)
			})
		})
	}
}

// BenchmarkFields measures reflection source reading, excluding node assembly
// and native message rendering. These types have no generated accessors.
// Cold evicts only the fixture's root and dynamic type plans before each call.
// Cache eviction is outside the timer; compilation, static validation, cache
// insertion, and reading current values are included. Warm reuses those plans.
func BenchmarkFields(b *testing.B) {
	footer := "Complete"
	nested := &benchmarkNestedSource{
		Title: "Report",
		Items: []benchmarkItemSource{
			{Text: benchmarkTextSource{Prefix: "Hello ", Name: "Jane"}, Count: 123},
			{Text: benchmarkTextSource{Prefix: "Hello ", Name: "John"}, Count: 456},
			{Text: benchmarkTextSource{Prefix: "Hello ", Name: "Alex"}, Count: 789},
		},
		Footer: &footer,
	}
	cases := []struct {
		name         string
		input        any
		fieldCount   int
		dynamicTypes []reflect.Type
	}{
		{
			name: "Flat",
			input: &benchmarkFlatSource{
				Title: "Notice", Prefix: "Hello ", Name: "Jane", Count: 123, Enabled: true,
			},
			fieldCount: 5,
		},
		{name: "Nested", input: nested, fieldCount: 3},
		{
			name: "Dynamic", input: &benchmarkDynamicSource{Content: nested}, fieldCount: 1,
			dynamicTypes: []reflect.Type{reflect.TypeOf(nested)},
		},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			rootType := reflect.TypeOf(tc.input).Elem()
			evict := func() {
				sourcePlans.Delete(planKey{rootType, testProfile.Platform})
				for _, typeOf := range tc.dynamicTypes {
					rawPlans.Delete(planKey{typeOf, testProfile.Platform})
				}
			}
			for _, mode := range []string{"Cold", "Warm"} {
				b.Run(mode, func(b *testing.B) {
					evict()
					b.Cleanup(evict)
					fields, err := Fields(tc.input, testProfile.Platform)
					if err != nil {
						b.Fatal(err)
					}
					if len(fields) != tc.fieldCount {
						b.Fatalf("Fields returned %d fields, want %d", len(fields), tc.fieldCount)
					}
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						if mode == "Cold" {
							b.StopTimer()
							evict()
							b.StartTimer()
						}
						fields, err = Fields(tc.input, testProfile.Platform)
						if err != nil {
							b.Fatal(err)
						}
					}
					runtime.KeepAlive(fields)
				})
			}
		})
	}
}
