# convertago

Convert tagged Go structs into Kakao Work Block Kit, Slack Block Kit, and Google
Chat Cards v2 messages. Each messenger uses its own tag namespace, so the same
source struct can describe messages for several platforms.

## Installation

Go 1.25 or later is required.

```sh
go get github.com/JSYoo5B/convertago
```

## Quick start

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/JSYoo5B/convertago"
)

type Notice struct {
    Title  string `kakaowork:"header" slack:"header" googlechat:"header"`
    Prefix string `kakaowork:"text;group=body" slack:"rich_text;group=body" googlechat:"textParagraph;group=body"`
    Name   string `kakaowork:"text;group=body;style=bold" slack:"rich_text;group=body;style=bold" googlechat:"textParagraph;group=body;style=bold"`
}

func main() {
    notice := Notice{Title: "Notice", Prefix: "Hello ", Name: "Jane"}
    message, err := convertago.ToSlackMessage(notice)
    if err != nil {
        panic(err)
    }
    data, err := json.Marshal(message)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(data))
}
```

`ToKakaoworkMessage`, `ToSlackMessage`, and `ToGoogleChatMessage` accept a struct
or a non-nil pointer to a struct. They return the corresponding native message
and an error. The sending application handles message delivery.

Fields contribute in declaration order. `group` joins fields into one native
node; include any spaces or newlines in the source values. An absent tag, an
empty tag, or `-` excludes a field for that platform. See the
[tag conversion guide](docs/tags.md) for nesting, slots, omission, and strict-mode
diagnostics.

## Messenger support

- [Kakao Work](kakaowork/README.md)
- [Slack](slack/README.md)
- [Google Chat](googlechat/README.md)

Each package documents its supported roles, slots, styles, and validation rules.
Its native models are also available for constructing or editing messages
directly. Conversion implementations live in the corresponding messenger package.

## Generated accessors

Known struct types can use generated field accessors. Add this directive to the
source package and run `go generate`:

```go
//go:generate go run github.com/JSYoo5B/convertago/cmd/convertago -type Notice -output zz_convertago.gen.go
```

Commit the generated file alongside the source. The conversion functions select
generated accessors automatically and otherwise use cached reflection. Dynamic
interface fields retain a reflection fallback. Both paths share tag semantics
and validation. See the [generation guide](docs/generation.md) for supported
source types, regeneration, and cross compilation.

## Further reading

| Guide | Contents |
| --- | --- |
| [Tag conversion](docs/tags.md) | Grammar, builders, ordering, groups, nesting, omissions, and diagnostics. |
| [Code generation](docs/generation.md) | Generator usage, generated-code contract, and target-specific layouts. |
| [Development](docs/development.md) | Repository layout, test placement, CI, fuzzing, and cross-build checks. |
| [Benchmark execution](benchmarks/running.md) | Allocation reporting, process CPU, and peak memory measurements. |
| [Benchmark results](benchmarks/README.md) | Recorded reflection-cache and generated-accessor comparisons. |

## License

[MIT](LICENSE).
