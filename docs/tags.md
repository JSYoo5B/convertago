# Tag conversion

The conversion functions read each platform's tag namespace independently.
See the [quick start](../README.md#quick-start) for a basic conversion and the
[generated accessor guide](generation.md) for generating direct field readers.

## Tag grammar

Each platform uses its own tag key. The grammar is
`role;option=value;flag`, with exact, case-sensitive names and no aliases.
Absent tags, empty tags, and `-` exclude a field for that platform. Tags on
unexported fields are errors. Unknown roles, options, styles, formats, duplicate
options, duplicate styles, and malformed tags are errors even when the field is
empty or optional. Value options require a nonempty value, and flags reject one.

| Option | Meaning |
| --- | --- |
| `group=name` | Assemble contributions into one native node in this source struct. |
| `slot=name` | Supply a named input of the selected builder. |
| `style=bold,italic` | Style this field's text contributions. Available styles depend on the builder. |
| `format=plain` | Explicitly select a text syntax accepted by the builder. |
| `omitempty` | Omit an empty source value. |
| `optional` | Report a recognized unavailable feature as a Warning and skip it. |

The recognized styles are `bold`, `italic`, `strike`, `code`, `underline`,
`highlight`, `client_highlight`, and `unlink`. The recognized formats are `plain`,
`html`, `markdown`, and `mrkdwn`. Each builder accepts a subset, listed in its
package README.

Two roles are shared by every platform:

| Role | Meaning |
| --- | --- |
| `part` | Supply a scalar input to the containing builder. A top-level `part` is an error. |
| `flatten` | Emit a tagged object's children into the surrounding output. It accepts only `omitempty` and `optional`. |

## Messenger support

- [Kakao Work](../kakaowork/README.md)
- [Slack](../slack/README.md)
- [Google Chat](../googlechat/README.md)

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

Converters check the native values they build with the package's
[message rules](message-rules.md).

## Order, groups, and nesting

Fields contribute in declaration order. A group occupies the position of its
first **declared** member, even when that member is nil or omitted. Its remaining
members concatenate in declaration order. Groups cannot mix builder roles, and
inside a builder they cannot mix containing input slots.
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
a nonempty list. Styles and formats on metadata slots are unavailable features
unless that slot explicitly supports them, as described in
[diagnostics](#diagnostics-and-severities).

Strings, booleans, integers, and floats have direct text representations.
Nested fields whose type implements `encoding.TextMarshaler` or `fmt.Stringer`
use those methods, preferring `MarshalText`. The field's own type must implement
the method: a value field whose method has a pointer receiver is read as a
struct, so use a pointer field in that case. A `MarshalText` error is an
`invalid_source` error. A root struct always reads its tags. Other structs
require tagged children; arbitrary structs and maps are not implicitly formatted
as text. Maps, channels, functions, and complex numbers are unsupported source
types even with `omitempty`. Cyclic source traversals return an error.

Embedded structs do not promote their fields. Tag an embedded field with
`flatten` to read its tagged fields as if they were declared in the outer struct;
an untagged embedded field is ignored like any other untagged field.

Zero and false remain present unless `omitempty` is set. Empty scalar strings,
zero numbers, false, and zero-length lists are empty. A tagged object is empty
when all its selected fields are empty. Explicit text representations are empty
when their returned text is empty. Nil pointers, slices, and interfaces are
absent; non-nil pointers and interfaces preserve their presence even when they
contain a zero value. A present builder with no inputs is an error unless its role supports an empty
object, such as a divider. A missing required slot after omissions is an error.

## Diagnostics and severities

The currently modeled native roles are available. Recognized styles and formats
that a selected builder cannot express are unavailable. For example,
`header;style=bold;optional` can skip a header whose native object has no inline
style support. Unknown roles or invalid nesting remain errors.

Every diagnostic has a severity, as defined in the
[message rules](message-rules.md#severities). Fatal diagnostics are errors.
Warning diagnostics go to `WithDiagnostics`, and `WithWarningAsError` turns them into
errors. Advisory diagnostics always go to `WithDiagnostics`. A diagnostic carries
its platform, path, code, message, and severity.

An active unavailable feature is a Fatal error by default. With `optional` it is
a Warning, so it is skipped by default and rejected with `WithWarningAsError`.
Intentional exclusions, nil values, and `omitempty` are never diagnosed.
Typos, type mismatches, and missing required inputs are always Fatal.

Converters check the native message with the same rules as the package's native
validation. A rule violation's code is the rule ID, such as `header.text.length`,
and each package README lists its rules and their severities. Tag and source
problems are Fatal and use these codes:

| Code | Meaning |
| --- | --- |
| `invalid_tag` | Malformed tag, unknown name, tag on an unexported field, or role used in an invalid position. |
| `invalid_source` | Unsupported source type, non-struct or cyclic source, or a failing text method. |
| `invalid_value` | A source value that a slot cannot parse, such as a non-integer. |
| `group_conflict` | Group members that mix roles or containing slots. |
| `missing_input` | A required slot without inputs after omissions. |
| `duplicate_input` | Several contributions to a slot that accepts one. |
| `conflicting_input` | Inputs that cannot be combined in one slot. |
| `conflicting_format` | Text formats that cannot be mixed in one slot. |
| `invalid_position` | A builder in a position its container does not allow. |
| `unsupported_feature` | An unavailable feature; a Warning when the tag is `optional`. |

Errors are `convertago.Diagnostic` values and can be inspected with `errors.As`.
Paths start at `$`, such as `$.Items[1].Name`; static list schema errors use `[]`
when there is no specific element. `Diagnostic.Error` formats a diagnostic as
`convertago: <platform> <path>: <message> [<severity> <code>]`.

```go
message, err := convertago.ToSlackMessage(notice,
    convertago.WithWarningAsError(),
    convertago.WithDiagnostics(func(d convertago.Diagnostic) {
        log.Print(d)
    }),
)
```
