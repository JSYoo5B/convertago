# Reflection cache and generated accessor benchmarks

See the [benchmark execution guide](running.md) for commands and the CPU and
memory measurement method, and the [development guide](../docs/development.md)
for test placement and the full verification steps.

Apple M3 Pro, darwin/arm64, Go 1.27.1, GOMAXPROCS=1. Measured on 2026-10-06 at
commit `1207bd4`. Each value is the median of three runs without the race
detector.

CPU time is the sum of the user and system time returned by `os.wait4`. It
includes process startup, benchmark setup, and garbage collection, and is
measured separately from the elapsed `ns/op`. Generation and compilation finish
before measurement.

## Scope of the generated accessor comparison

The generator analyzes the AST and type information to write `ConvertagoFields`
accessors. Generated code reads known fields directly and uses parsed tag
literals with shared style arrays. Both paths go through the same entry point for
input checks and accessor selection, and share message assembly, rule checks,
and native model construction. The Generated rows measure the runtime cost of
these accessors. They exclude AST analysis, code generation, and Go compilation.

The Reflection rows use separate defined types with the same field layout and
tags, which drop the accessor methods and select reflection. Both inputs refer to
the same values. Flat contains a header, grouped text, and an integer. Nested
contains an image, a list, and pointers. Dynamic places the same Nested value in
an `any` field, which generated code also reads through cached reflection. Tests
verify each path's selection, matching source fields, and matching final JSON,
and validate the converted messages with the native rules.

Generated accessors reduced the CPU cost of reading Flat and Nested fields by
about 26 to 31 percent. Because the work after reading the source is shared, the
reduction in complete conversion was smaller: 0.9 to 8.8 percent for Flat and
Nested inputs. Generated accessors allocated 64 B and one allocation less per
Flat operation and 88 B and two allocations less per Nested operation. Dynamic
conversion differed by about 1 percent and saved 16 B and one allocation; the
reflection cost of the dynamic field remains.

## Cost of message rules

Converters now check every native value they build with the package's
[message rules](../docs/message-rules.md). Compared with the previous
measurement at commit `6ac5540`, recorded in `measurements.csv` at commit
`23c0217`, complete conversion with reflection took more process CPU and
allocated more often. Field reads, which exclude rule checks, did not change
beyond measurement variance.

| Messenger | Input | CPU s / 100k before → after | allocs/op before → after |
| --- | --- | ---: | ---: |
| Kakao Work | Flat | 0.2669 → 0.4346 | 46 → 91 |
| Kakao Work | Nested | 0.8455 → 1.1906 | 124 → 216 |
| Slack | Flat | 0.3040 → 0.5076 | 56 → 121 |
| Slack | Nested | 1.0615 → 1.5694 | 154 → 269 |
| Google Chat | Flat | 0.5960 → 0.7319 | 75 → 125 |
| Google Chat | Nested | 1.7879 → 2.0110 | 194 → 278 |

Most of the added allocations come from the check context created for each
native value and from field path strings built for nested values.

## Field reads: cached reflection and generated accessors

Each case ran 500,000 times in its own process. `BenchmarkSourceFields` compares
Reflection and Generated through the same `conversion.Fields` entry point. The
first call prepares the caches. The measurement includes reading current values
and excludes message assembly and native rendering.

| Messenger | Input | Path | ns/op | CPU s / 500k | B/op | allocs/op | Peak RSS MiB |
| --- | --- | --- | ---: | ---: | ---: | ---: | ---: |
| Kakao Work | Flat | Reflection | 397.1 | 0.2024 | 1,091 | 3 | 13.00 |
| Kakao Work | Flat | Generated | 279.3 | 0.1404 | 1,027 | 2 | 12.83 |
| Kakao Work | Nested | Reflection | 1,387.0 | 0.6847 | 3,432 | 10 | 13.06 |
| Kakao Work | Nested | Generated | 1,035.0 | 0.5089 | 3,344 | 8 | 12.89 |
| Kakao Work | Dynamic | Reflection | 1,594.0 | 0.7879 | 3,608 | 11 | 12.94 |
| Kakao Work | Dynamic | Generated | 1,555.0 | 0.7671 | 3,592 | 10 | 13.22 |
| Slack | Flat | Reflection | 389.2 | 0.1961 | 1,091 | 3 | 12.84 |
| Slack | Flat | Generated | 277.3 | 0.1414 | 1,027 | 2 | 12.66 |
| Slack | Nested | Reflection | 1,452.0 | 0.7155 | 3,656 | 10 | 12.86 |
| Slack | Nested | Generated | 1,073.0 | 0.5269 | 3,568 | 8 | 12.78 |
| Slack | Dynamic | Reflection | 1,656.0 | 0.8172 | 3,832 | 11 | 16.45 |
| Slack | Dynamic | Generated | 1,586.0 | 0.7828 | 3,816 | 10 | 13.05 |
| Google Chat | Flat | Reflection | 394.4 | 0.1998 | 1,091 | 3 | 12.69 |
| Google Chat | Flat | Generated | 280.3 | 0.1410 | 1,027 | 2 | 12.70 |
| Google Chat | Nested | Reflection | 1,501.0 | 0.7403 | 3,656 | 10 | 12.86 |
| Google Chat | Nested | Generated | 1,071.0 | 0.5332 | 3,568 | 8 | 13.06 |
| Google Chat | Dynamic | Reflection | 1,653.0 | 0.8178 | 3,832 | 11 | 13.50 |
| Google Chat | Dynamic | Generated | 1,605.0 | 0.7912 | 3,816 | 10 | 12.97 |

## Complete conversion per messenger

Each case ran 100,000 times in its own process. Both paths use the same input,
and the caches are prepared before conversion. The measurement includes reading
values, assembly, rule checks, and native message construction. It excludes the
caller's JSON serialization but includes serialization inside a converter.

| Messenger | Input | Path | ns/op | CPU s / 100k | B/op | allocs/op | Peak RSS MiB |
| --- | --- | --- | ---: | ---: | ---: | ---: | ---: |
| Kakao Work | Flat | Reflection | 4,350.0 | 0.4346 | 5,672 | 91 | 13.88 |
| Kakao Work | Flat | Generated | 4,163.0 | 0.4176 | 5,608 | 90 | 13.69 |
| Kakao Work | Nested | Reflection | 11,999.0 | 1.1906 | 16,121 | 216 | 14.30 |
| Kakao Work | Nested | Generated | 11,523.0 | 1.1427 | 16,033 | 214 | 14.05 |
| Kakao Work | Dynamic | Reflection | 15,202.0 | 1.5081 | 17,529 | 236 | 19.94 |
| Kakao Work | Dynamic | Generated | 15,035.0 | 1.4911 | 17,513 | 235 | 31.81 |
| Slack | Flat | Reflection | 5,085.0 | 0.5076 | 6,608 | 121 | 14.02 |
| Slack | Flat | Generated | 5,020.0 | 0.5032 | 6,544 | 120 | 14.06 |
| Slack | Nested | Reflection | 15,992.0 | 1.5694 | 17,985 | 269 | 14.36 |
| Slack | Nested | Generated | 14,442.0 | 1.4316 | 17,897 | 267 | 14.36 |
| Slack | Dynamic | Reflection | 17,888.0 | 1.7799 | 19,450 | 291 | 14.58 |
| Slack | Dynamic | Generated | 17,759.0 | 1.7620 | 19,434 | 290 | 14.48 |
| Google Chat | Flat | Reflection | 7,350.0 | 0.7319 | 8,193 | 125 | 14.25 |
| Google Chat | Flat | Generated | 7,198.0 | 0.7179 | 8,129 | 124 | 14.67 |
| Google Chat | Nested | Reflection | 20,264.0 | 2.0110 | 22,036 | 278 | 14.70 |
| Google Chat | Nested | Generated | 19,669.0 | 1.9541 | 21,948 | 276 | 14.70 |
| Google Chat | Dynamic | Reflection | 23,704.0 | 2.3556 | 23,500 | 300 | 15.05 |
| Google Chat | Dynamic | Generated | 23,924.0 | 2.3714 | 23,484 | 299 | 15.08 |

## Shared style arrays

This historical comparison measured commit `6bf6b97`, which built generated
accessor style arrays on every call, against commit `3ed514b`, which shares
them. Both binaries were built first, and each messenger, input, and count ran
three times. The second repetition ran the newer binary first to swap the run
order. Field reads ran 500,000 times and conversions 100,000 times, measured as
process CPU time. The comparison predates the removal of redundant HTTPS checks
and the message rules.

After sharing, Flat allocated 16 B and one allocation less, and Nested 48 B and
three allocations less. Dynamic reads its styled inner fields through reflection,
so its allocations did not change. The CPU values include run-to-run variance,
and differences also appeared for Dynamic, whose allocations did not change. The
clear effect of sharing styles is fewer allocations for static inputs; the CPU
difference in complete conversion was small.

| Messenger | Input | Operation | CPU s before → after | B/op before → after | allocs/op before → after |
| --- | --- | --- | ---: | ---: | ---: |
| Kakao Work | Flat | Field reads | 0.1539 → 0.1421 | 1,043 → 1,027 | 3 → 2 |
| Kakao Work | Flat | Conversion | 0.2572 → 0.2542 | 3,440 → 3,424 | 46 → 45 |
| Kakao Work | Nested | Field reads | 0.5256 → 0.5079 | 3,392 → 3,344 | 11 → 8 |
| Kakao Work | Nested | Conversion | 0.8191 → 0.7958 | 11,160 → 11,112 | 125 → 122 |
| Kakao Work | Dynamic | Field reads | 0.7591 → 0.7403 | 3,592 → 3,592 | 10 → 10 |
| Kakao Work | Dynamic | Conversion | 1.1316 → 1.1349 | 12,584 → 12,584 | 143 → 143 |
| Slack | Flat | Field reads | 0.1452 → 0.1402 | 1,043 → 1,027 | 3 → 2 |
| Slack | Flat | Conversion | 0.2938 → 0.2906 | 3,736 → 3,720 | 56 → 55 |
| Slack | Nested | Field reads | 0.5417 → 0.5212 | 3,616 → 3,568 | 11 → 8 |
| Slack | Nested | Conversion | 1.0168 → 1.0153 | 12,912 → 12,864 | 155 → 152 |
| Slack | Dynamic | Field reads | 0.7874 → 0.7734 | 3,816 → 3,816 | 10 → 10 |
| Slack | Dynamic | Conversion | 1.3869 → 1.3766 | 14,400 → 14,400 | 175 → 175 |
| Google Chat | Flat | Field reads | 0.1493 → 0.1407 | 1,043 → 1,027 | 3 → 2 |
| Google Chat | Flat | Conversion | 0.5784 → 0.5814 | 5,912 → 5,896 | 75 → 74 |
| Google Chat | Nested | Field reads | 0.5379 → 0.5249 | 3,616 → 3,568 | 11 → 8 |
| Google Chat | Nested | Conversion | 1.7560 → 1.7422 | 18,283 → 18,235 | 197 → 194 |
| Google Chat | Dynamic | Field reads | 0.8100 → 0.7686 | 3,816 → 3,816 | 10 → 10 |
| Google Chat | Dynamic | Conversion | 2.1468 → 2.1116 | 19,771 → 19,771 | 217 → 217 |

## Type metadata cache

Each case ran 2,000,000 times in its own process with the internal shared test
profile. Cached looks up an existing type plan, and Uncached parses tags and
validates the structure on every call. Uncached excludes cache insertion and
removal. Reading values and building messages are excluded.

| Input | Cache | ns/op | CPU s / 2M | B/op | allocs/op | Peak RSS MiB |
| --- | --- | ---: | ---: | ---: | ---: | ---: |
| Flat | Cached | 25.6 | 0.0564 | 0 | 0 | 7.75 |
| Flat | Uncached | 2,592.0 | 5.0767 | 5,224 | 27 | 13.34 |
| Nested | Cached | 25.2 | 0.0544 | 0 | 0 | 7.66 |
| Nested | Uncached | 3,307.0 | 6.5315 | 5,008 | 44 | 13.17 |

## Reading current values: cold and warm

The internal shared test profile inputs ran for 100ms, three times each. Cold
removes the root and dynamic type plans of the input before each `Fields` call;
the removal is excluded from timing and allocation measurement. Pausing and
resuming the timer distorts process CPU time, so this table compares the elapsed
time and allocations reported by Go benchmarks. The inputs differ from the
per-messenger SourceFields inputs above.

| Input | Cache | ns/op | B/op | allocs/op |
| --- | --- | ---: | ---: | ---: |
| Flat | Cold | 4,026.0 | 6,552 | 33 |
| Flat | Warm | 400.1 | 1,219 | 3 |
| Nested | Cold | 6,129.0 | 9,360 | 62 |
| Nested | Warm | 1,431.0 | 4,225 | 15 |
| Dynamic | Cold | 5,895.0 | 10,208 | 67 |
| Dynamic | Warm | 1,693.0 | 4,417 | 16 |

## Memory and reproduction

`B/op` and `allocs/op` are allocations per call. Peak RSS is the highest value
for the whole process, including the Go runtime and test initialization, and
varies between runs. These values are separate from the memory the caches
retain. Cache hits and generated accessors still allocate to read current values
and build messages. Shared style arrays are static package metadata; copy a
returned tag's styles before editing them.

All 138 samples are recorded in [measurements.csv](measurements.csv). Fields
rows leave the CPU and RSS columns empty because of their timer pauses. The shared
style comparison recorded 54 samples each in
[static-tags-before.csv](static-tags-before.csv) and
[static-tags-after.csv](static-tags-after.csv). Each table shows the median of
its CSV.

The [measurement runner](measure.py) builds the binaries first and runs both
paths with the same counts on macOS or Linux. It updates the CSV only after every
measurement succeeds.

```sh
go generate ./internal/benchmarksource
python3 benchmarks/measure.py --output benchmarks/measurements.csv
```

To reproduce a historical measurement, check out its commit and run the
`measure.py` from that commit, because the benchmark packages have moved since.
See the [benchmark execution guide](running.md) for ordinary Go benchmark
commands. CI exercises every benchmark once without timing thresholds;
performance values are observations from this environment and are not pass
conditions.
