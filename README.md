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

## Initial builders

The native packages contain more models than the initial tag converters support.
The converters currently support the following builders and input slots.
Text slots concatenate contributions without inserting a separator.

| Platform | Role | Default slot | Other slots |
| --- | --- | --- | --- |
| Kakao Work | `header` | `text` | `style`, a single background color: white, blue, red, or yellow |
| Kakao Work | `text` | `text` | — |
| Kakao Work | `image_link` | `url` | — |
| Kakao Work | `preview` | `text` | — |
| Slack | `header` | `text` | — |
| Slack | `section` | `text` | — |
| Slack | `rich_text` | `text` | — |
| Slack | `image` | `url` | `alt` required; `title` optional |
| Google Chat | `header` | `text` | `subtitle`, `url`, `alt` optional |
| Google Chat | `textParagraph` | `text` | — |
| Google Chat | `image` | `url` | `alt` optional |
| Google Chat | `fallbackText` | `text` | — |

Every builder requires its default slot. Only text contributions and Google Chat
header subtitles can repeat. Duplicate URL, alt, title, or background inputs are
errors. Image URLs must be absolute HTTP or HTTPS URLs; Google Chat requires
HTTPS. Slack alt text must be non-empty and is never fabricated from the URL.

Kakao Work supports `bold`, `italic`, and `strike` on `text`. Slack supports
`bold`, `italic`, `strike`, `code`, and `underline` on `rich_text`. Google Chat
supports those styles on `textParagraph`; Markdown does not support underline.
Style flags describe inline text, while Kakao Work's header `slot=style` supplies
its background color.

Strings are literal by default. Kakao Work `text` and Slack `rich_text` also
accept `format=plain`. Slack `section` accepts `plain` or explicit `mrkdwn`;
a grouped section cannot mix the two. Google Chat `textParagraph` escapes plain
text into HTML and converts newlines into `<br>`. Explicit `html` and `markdown`
retain the supplied markup. HTML can include escaped plain contributions;
Markdown cannot share a paragraph with plain or HTML contributions.

Converters validate supported native constraints, including Kakao Work's header
position and text lengths, Slack's text lengths and 50-block message limit, and
Google Chat's 100-widget and 32 KB card limits. Google Chat places widgets in one
section of one card. Its single header must precede widgets. Values that violate
these rules return an error; content is not truncated or silently reordered.

## Order, groups, and nesting

Fields contribute in declaration order. A group occupies the position of its
first **declared** member, even when that member is nil or omitted. Its remaining
members concatenate in declaration order. Groups cannot mix builder roles.
Group names belong to one source struct; nested structs and repeated elements
have separate scopes. Include any spaces or newlines in the source values.

A tagged parent struct selects a builder. Its scalar children use `part` to supply that
builder's inputs, allowing domain types to differ from the native message shape.
The default slot is inherited when a child has no `slot`. Parent text style and
format apply to its child contributions, with child styles added and an explicit
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
contain a zero value. A present builder with no inputs, or a missing required
slot after omissions, is an error.

## Diagnostics and strict mode

Recognized unavailable converter roles are Kakao Work `button`, `action`,
`divider`, `description`, `section`, and `context`; Slack `actions`, `context`,
`divider`, `markdown`, and `video`; Google Chat `decoratedText`, `buttonList`,
`divider`, `columns`, `grid`, `carousel`, and `chipList`. Recognized style or
format names that a selected builder cannot express are also unavailable.

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
