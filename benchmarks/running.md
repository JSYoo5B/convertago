# Running benchmarks

Run these commands from the repository root. See
[measured CPU and memory comparisons](README.md) for the recorded
environment, results, and individual samples.

## Fixtures

The benchmark fixtures include accessors produced by this repository's AST and
type-based generator. Regenerate them after changing the fixtures or generator:

```sh
go generate ./internal/benchmarksource
```

## Allocation benchmarks

Run the reflection cache, source reader, and messenger benchmarks with allocation reporting:

```sh
go test -run '^$' -bench 'Benchmark(SourcePlan|Fields|SourceFields|To.*Message)$' -benchmem -count=3 ./benchmarks ./internal/conversion
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

## Process CPU and memory

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
memory. Build the message benchmark binary with
`go test -c -o /tmp/convertago-message.bench ./benchmarks` to measure messenger
cases in the same way.

## Measurement runner

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
