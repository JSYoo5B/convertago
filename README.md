# convertago

Convert tagged Go structs into Kakao Work Block Kit, Slack Block Kit, and Google
Chat Cards v2 messages. Native models remain available in `kakaowork`, `slack`,
and `googlechat`; their conversion implementations live in those packages.

```go
type Notice struct {
    Title  string `kakaowork:"header" slack:"header" googlechat:"header"`
    Prefix string `kakaowork:"text;group=body" slack:"rich_text;group=body" googlechat:"textParagraph;group=body"`
    Name   string `kakaowork:"text;group=body;style=bold" slack:"rich_text;group=body;style=bold" googlechat:"textParagraph;group=body;style=bold"`
}

notice := Notice{Title: "Notice", Prefix: "Hello ", Name: "Jane"}
message, err := convertago.ToSlackMessage(notice)
if err != nil {
    return err
}
data, err := json.Marshal(message)
```

`ToKakaoworkMessage`, `ToSlackMessage`, and `ToGoogleChatMessage` accept a struct
or a non-nil pointer to a struct. They return the corresponding native message
and an error. They do not send the message.

## Tag grammar

Each platform uses its own tag key. The grammar is
`role;option=value;flag`, with exact, case-sensitive names and no aliases.
Absent tags, empty tags, and `-` exclude a field for that platform. Tags on
unexported fields are errors. Unknown roles, options, styles, formats, duplicate
options, and malformed tags are errors even when the field is empty or optional.

| Option | Meaning |
| --- | --- |
| `group=name` | Assemble contributions into one native node in this source struct. |
| `slot=name` | Supply a named input of the selected builder. |
| `style=bold,italic` | Style this field's text contributions. Available styles depend on the builder. |
| `format=plain` | Explicitly select a text syntax accepted by the builder. |
| `omitempty` | Omit an empty source value. |
| `optional` | Skip a recognized unavailable feature with a diagnostic in normal mode. |

## Messenger support

- [Kakao Work](kakaowork/README.md)
- [Slack](slack/README.md)
- [Google Chat](googlechat/README.md)

## Builders

The tag converters cover the native models defined in this repository.
Text slots concatenate contributions without inserting a separator. Other
scalar properties normally accept one contribution. Child slots accept only
the native roles appropriate to that parent.

The default scalar slot preserves the convenient single-field form. For
containers, tag a struct whose children select native roles. A child's `slot`
belongs to its parent; a top-level field's `slot` belongs to its own builder.
A compatible default child slot is used when the profile specifies one.
Otherwise a uniquely compatible slot is inferred, and multiple matches require
an explicit `slot`.

The converters validate modeled native constraints and conflicting inputs.

## Order, groups, and nesting

Fields contribute in declaration order. A group occupies the position of its
first **declared** member, even when that member is nil or omitted. Its remaining
members concatenate in declaration order. Groups cannot mix builder roles.
Group names belong to one source struct; nested structs and repeated elements
have separate scopes. Include any spaces or newlines in the source values.

A tagged parent struct selects a builder. Its scalar children use `part` to supply that
builder's inputs, allowing domain types to differ from the native message shape.
The default slot is inherited when a child has no `slot`. Parent text style and
format apply to contributions in its default text slot, with child styles added and an explicit
child format replacing the inherited format.

```go
type Picture struct {
    URL string `kakaowork:"part;slot=url" slack:"part;slot=url" googlechat:"part;slot=url"`
    Alt string `slack:"part;slot=alt" googlechat:"part;slot=alt"`
}

type Item struct {
    Prefix string `kakaowork:"text;group=line" slack:"rich_text;group=line" googlechat:"textParagraph;group=line"`
    Name   string `kakaowork:"text;group=line;style=bold" slack:"rich_text;group=line;style=bold" googlechat:"textParagraph;group=line;style=bold"`
}

type Report struct {
    Picture Picture `kakaowork:"image_link" slack:"image" googlechat:"image"`
    Items   []Item  `kakaowork:"flatten" slack:"flatten" googlechat:"flatten"`
}
```

`flatten` emits a tagged object's children into the surrounding output. Inside a
builder it supplies child parts to that builder. Flattening a scalar is an error.
A slice or array repeats in element order. An ungrouped builder creates a node
per element; a grouped field contributes its elements to its group's one node.
Inside a builder, scalar lists concatenate parts. Native child builders use a
parent slot that accepts their role; an explicit `slot` selects among several
accepted positions. Child lists repeat native elements in declaration order.
Empty container structs can represent dividers. Nil elements are absent.
`omitempty` on a list omits an empty list and retains zero-valued elements in
a nonempty list. Styles and formats on metadata slots are rejected unless that
slot explicitly supports them.

Strings, booleans, integers, and floats have direct text representations.
Nested fields that explicitly implement `encoding.TextMarshaler` or `fmt.Stringer`
use those methods, preferring `MarshalText`. A root struct always reads its tags.
Other structs require tagged children; arbitrary structs and maps are not
implicitly formatted as text. Cyclic source traversals return an error.

Zero and false remain present unless `omitempty` is set. Empty scalar strings,
zero numbers, false, and zero-length lists are empty. A tagged object is empty
when all its selected fields are empty. Explicit text representations are empty
when their returned text is empty. Nil pointers, slices, and interfaces are
absent; non-nil pointers and interfaces preserve their presence even when they
contain a zero value. A present builder with no inputs is an error unless its role supports an empty
object, such as a divider. A missing required slot after omissions is an error.

## Diagnostics and strict mode

The currently modeled native roles are available. Recognized styles and formats
that a selected builder cannot express are unavailable. For example,
`header;style=bold;optional` can skip a header whose native object has no inline
style support. Unknown roles or invalid nesting remain errors.

An active unavailable feature errors by default. With `optional` it is skipped
in normal mode, and `WithDiagnostics` receives its platform, source path, code,
and message. `WithStrict` rejects the active feature even when optional.
Intentional exclusions, nil values, and `omitempty` still work in strict mode.
Typos, type mismatches, missing required inputs, and native validation failures
always error. Errors are `convertago.Diagnostic` values and can be inspected with
`errors.As`. Paths start at `$`, such as `$.Items[1].Name`; static list schema
errors use `[]` when there is no specific element.

```go
message, err := convertago.ToSlackMessage(notice,
    convertago.WithStrict(),
    convertago.WithDiagnostics(func(d convertago.Diagnostic) {
        log.Print(d)
    }),
)
```

## Generated accessors

For known struct types, add a directive to their source package:

```go
//go:generate go run github.com/JSYoo5B/convertago/cmd/convertago -type Notice -output zz_convertago.gen.go
```

Run `go generate` and commit the generated file alongside the source. Multiple
root types can be supplied as a comma-separated `-type` list. The generator uses
the active Go package files, AST, compiled dependency export data, and `go/types`
to resolve imports, aliases, field types, and text methods. It validates tags with
the same platform profiles as reflection before writing the output. Regeneration
excludes the previous output from type checking and refuses to replace handwritten
files. Generation currently targets non-generic defined struct types in packages
without cgo.

Generated `ConvertagoFields` methods access known fields directly, check nil
pointers, iterate lists, and emit parsed tag metadata. The ordinary `To…Message`
functions select these methods automatically. Interface-valued fields use a
cached reflection fallback for their actual runtime types. Without generated
methods, the entire source uses cached reflection plans by type and platform.
Both paths share ordering, grouping, assembly, validation, and diagnostics.

`SourceField`, `SourceValue`, `SourceTag`, `SourceState`, `SourceObject`, and
`SourceMarshaled` form the generated-code contract. Applications normally use
tags, the conversion functions, and native message types instead of constructing
these source representations themselves. Rerun generation after changing tags
or field types.

JSON examples are attached to each conversion function in its corresponding
`*_example_test.go`. Generator integration tests compile a separate consumer
package and compare JSON, error diagnostics, and skip diagnostics for all three
platforms against reflection.

## Verification

CI runs tests, the race detector, and `go vet` on Go 1.25 and the current stable
release. It regenerates the checked-in benchmark accessors and rejects a diff.
It also runs bounded fuzz checks and exercises each benchmark without timing
thresholds.

`FuzzParse` checks canonical tag round trips. `FuzzGeneratedReflectionConversion`
compares generated and reflected diagnostics and JSON for flat, nested, and
dynamic inputs across all three platforms, then validates successful native
messages. Run longer local fuzz sessions with:

```sh
go test ./internal/conversion -run '^$' -fuzz '^FuzzParse$' -fuzztime=1m
go test . -run '^$' -fuzz '^FuzzGeneratedReflectionConversion$' -fuzztime=1m
```

## Benchmarks

See [measured CPU and memory comparisons](benchmarks/README.md) for the recorded
environment, results, and individual samples.

The benchmark fixtures include accessors produced by this repository's AST and
type-based generator. Regenerate them after changing the fixtures or generator:

```sh
go generate ./internal/benchmarksource
```

Run the reflection cache, source reader, and messenger benchmarks with allocation reporting:

```sh
go test -run '^$' -bench 'Benchmark(SourcePlan|Fields|SourceFields|To.*Message)$' -benchmem -count=3 . ./internal/conversion
```

`BenchmarkSourcePlan` compares a cached type-plan lookup with compiling and
statically validating the same type on every call. Its uncached case excludes
cache insertion. `BenchmarkFields` measures the source reader with cold and warm
plans, including a dynamic interface value's separate type cache. The cold case
evicts only the fixture's plans outside the timer, then measures compilation,
validation, insertion, and reading current values. Its timer pauses add overhead
to the benchmark process, so use `BenchmarkSourcePlan` for process CPU comparisons.

`BenchmarkSourceFields` compares cached reflection with generated accessors for
each platform through the same source-reading entry point, excluding assembly
and native rendering. The generated and reflection inputs share field layouts,
tags, and values. Defined copies of the generated types drop their methods to
select reflection. A test verifies both reader dispatch and matching source
fields and native JSON before relying on the measurements.

`BenchmarkToKakaoworkMessage`, `BenchmarkToSlackMessage`, and
`BenchmarkToGoogleChatMessage` measure complete conversion after warming the
caches. Each flat, nested, and dynamic case has `Reflection` and `Generated`
sub-benchmarks using the same source values. AST analysis and code generation
happen before compilation and are excluded from runtime measurements. Generated
accessors read known fields directly and retain cached reflection for `any`
fields. Assembly and native validation remain shared. Caller-side JSON
marshaling is excluded; serialization inside a converter is included.

`ns/op` reports elapsed time per operation. `B/op` and `allocs/op` report allocated
bytes and allocation counts per operation, rather than retained cache memory or
peak resident memory. To compare CPU time, compile the benchmark binary once and
run each case with the same fixed iteration count:

```sh
go test -c -o /tmp/convertago-conversion.bench ./internal/conversion
/usr/bin/time -p /tmp/convertago-conversion.bench -test.run '^$' \
    -test.bench '^BenchmarkSourcePlan/Nested/Cached$' \
    -test.benchtime=2000000x -test.benchmem
/usr/bin/time -p /tmp/convertago-conversion.bench -test.run '^$' \
    -test.bench '^BenchmarkSourcePlan/Nested/Uncached$' \
    -test.benchtime=2000000x -test.benchmem
```

Compare the sum of `user` and `sys` CPU seconds. These process totals include
startup, benchmark setup, and garbage collection. Repeat measurements under the
same Go version, machine, and `GOMAXPROCS`; timings are informational and are not
test pass/fail thresholds. On macOS, `/usr/bin/time -l` also reports peak resident
memory. Build the root package's test binary to measure the messenger cases in
the same way.

The optional Python 3.9+ runner builds the test binaries once, measures process
CPU and peak RSS on macOS or Linux, and records three runs of every case:

```sh
python3 benchmarks/measure.py --output benchmarks/measurements.csv
```

It uses `GOMAXPROCS=1` with fixed counts of 2,000,000 type-plan operations,
500,000 source reads, and 100,000 message conversions. Cold/warm `Fields` cases
use 100ms and leave process CPU and RSS blank because their timer pauses add
benchmark harness overhead. The runner excludes generation and compilation
from process resource measurements and overwrites the selected CSV only after
all runs succeed.
