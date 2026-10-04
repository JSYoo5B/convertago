# Development

## Repository layout

| Location | Responsibility |
| --- | --- |
| `message.go` | Public message conversion functions for all three platforms. |
| `message_example_test.go` | Executable godoc examples for those conversion functions. |
| `conversion.go` | Public conversion options and diagnostics. |
| `generated_source.go` | Public contract used by generated accessors. |
| `kakaowork/`, `slack/`, `googlechat/` | Native models, platform conversion and validation, package tests, and function-specific examples. |
| `internal/conversion/` | Shared tag reading and assembly, with unit tests and reflection-cache benchmarks. |
| `internal/generate/` | Code generation, consumer compilation tests, and cross-compilation tests. |
| `internal/validation/` | Rule severities, rule checks and traversal, the go-playground validator adapter, and shared URI checks. |
| `internal/integration/` | Tests spanning messenger packages: shared validator configuration and generated/reflection fuzz parity. |
| `internal/benchmarksource/` | Shared source fixtures and their checked-in generated accessors. |
| `benchmarks/` | Message and source-reader benchmarks, the measurement runner, and recorded results. |
| `cmd/convertago/` | Generator command-line entry point. |
| `docs/` | Message rules, documentation policy, platform guide, tag rules, code generation, and development guidance. |

Keep public declarations in their existing packages to preserve application
imports and generated-code compatibility. Write examples and other documentation
as described in the [documentation policy](documentation-policy.md).

Unit tests that exercise a package's implementation stay in that package.
Checks that combine messenger packages belong in `internal/integration`.
Performance comparisons that span platforms belong in `benchmarks`; benchmarks
requiring private conversion internals stay in `internal/conversion`.
Generator consumer tests stay in `internal/generate` with their source fixtures
in `testdata`.

## Local checks

From the repository root, run:

```sh
go test -mod=readonly -count=1 ./...
go test -mod=readonly -race -count=1 ./...
go vet ./...
go generate ./internal/benchmarksource
git diff --exit-code -- internal/benchmarksource/zz_convertago.gen.go
```

`go test ./...` includes the public API examples, platform unit tests, integration
fuzz seeds, benchmark parity checks, and generator consumer tests. Bounded fuzz
runs and benchmark execution are separate checks.

Generator tests compile a separate consumer package and compare JSON, error
diagnostics, and skip diagnostics for all three platforms against reflection.

## CI

CI runs tests, the race detector, and `go vet` on Go 1.25 and the current stable
release. It regenerates the checked-in benchmark accessors and rejects a diff.
It also runs bounded fuzz checks and exercises each benchmark without timing
thresholds.

Generator integration tests compile consumer packages for Linux 386, Linux
arm64, and Windows amd64 without executing the target binaries. They cover both
portable host-generated accessors and target-specific layouts. A separate CI
matrix builds all packages for those targets and Darwin arm64 with cgo disabled.

## Fuzzing

`FuzzParse` checks canonical tag round trips. `FuzzGeneratedReflectionConversion`
compares generated and reflected diagnostics and JSON for flat, nested, and
dynamic inputs across all three platforms, then validates successful native
messages. Run longer local fuzz sessions with:

```sh
go test ./internal/conversion -run '^$' -fuzz '^FuzzParse$' -fuzztime=1m
go test ./internal/integration -run '^$' -fuzz '^FuzzGeneratedReflectionConversion$' -fuzztime=1m
```

## Cross-build checks

CI compiles all packages for Linux 386, Linux arm64, Windows amd64, and Darwin
arm64. To compile the same targets locally without running foreign binaries:

```sh
GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -mod=readonly ./...
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -mod=readonly ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -mod=readonly ./...
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -mod=readonly ./...
```

See [generation for cross-compilation targets](generation.md#cross-compilation)
when the source layout itself depends on the target.

## Benchmarking

See the [benchmark execution guide](../benchmarks/running.md) for allocation,
CPU, and memory measurements, and the [recorded results](../benchmarks/README.md)
for previous measurements. Timings are informational, with no pass/fail thresholds.
