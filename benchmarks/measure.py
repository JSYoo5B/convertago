"""Measure Go benchmark allocations and process CPU usage on macOS or Linux."""

import argparse
import csv
import os
from pathlib import Path
import re
import subprocess
import sys
import tempfile


RESULT = re.compile(
    r"^(Benchmark\S+)\s+(\d+)\s+([\d.]+) ns/op\s+"
    r"(\d+) B/op\s+(\d+) allocs/op$",
    re.MULTILINE,
)
COLUMNS = [
    "benchmark", "repetition", "iterations", "ns_per_op", "bytes_per_op",
    "allocs_per_op", "user_cpu_seconds", "system_cpu_seconds",
    "total_cpu_seconds", "peak_rss_bytes",
]


def parse_results(output):
    return [dict(
        benchmark=name, iterations=int(iterations), ns_per_op=float(elapsed),
        bytes_per_op=int(allocated), allocs_per_op=int(allocations),
    ) for name, iterations, elapsed, allocated, allocations in RESULT.findall(output)]


def measure(binary, name, iterations, environment, directory):
    # Files avoid a full stdout/stderr pipe blocking the child before wait4.
    with tempfile.TemporaryFile(dir=directory) as output, tempfile.TemporaryFile(dir=directory) as errors:
        process = subprocess.Popen([
            str(binary), "-test.run=^$", "-test.bench=^" + name + "$",
            f"-test.benchtime={iterations}x", "-test.benchmem", "-test.cpu=1",
        ], stdout=output, stderr=errors, env=environment)
        _, status, usage = os.wait4(process.pid, 0)
        process.returncode = os.waitstatus_to_exitcode(status)
        output.seek(0)
        text = output.read().decode()
        errors.seek(0)
        diagnostics = errors.read().decode()
    if process.returncode != 0:
        raise RuntimeError(f"{name} failed:\n{text}\n{diagnostics}")
    results = parse_results(text)
    if len(results) != 1 or results[0]["benchmark"] != name or results[0]["iterations"] != iterations:
        raise RuntimeError(f"unexpected benchmark selection:\n{text}")
    result = results[0]
    result.update(
        user_cpu_seconds=usage.ru_utime,
        system_cpu_seconds=usage.ru_stime,
        total_cpu_seconds=usage.ru_utime + usage.ru_stime,
        peak_rss_bytes=usage.ru_maxrss if sys.platform == "darwin" else usage.ru_maxrss * 1024,
    )
    return result


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output", required=True, type=Path, help="destination CSV")
    parser.add_argument("--count", type=int, default=3)
    args = parser.parse_args()
    if args.count < 1:
        parser.error("--count must be positive")
    if sys.platform != "darwin" and not sys.platform.startswith("linux"):
        parser.error("process CPU and RSS measurement requires macOS or Linux")

    root = Path(__file__).resolve().parent.parent
    environment = dict(os.environ, GOMAXPROCS="1")
    records = []
    with tempfile.TemporaryDirectory(prefix="convertago-bench-") as temporary:
        directory = Path(temporary)
        for binary, package in (("conversion", "./internal/conversion"), ("message", ".")):
            subprocess.run([
                "go", "test", "-mod=readonly", "-c", "-o", str(directory / binary), package,
            ], cwd=root, env=environment, check=True)

        # Cold Fields pauses its timer for cache eviction. Keep process CPU and
        # RSS blank because those totals would also include the timer harness.
        output = subprocess.check_output([
            str(directory / "conversion"), "-test.run=^$", "-test.bench=^BenchmarkFields$",
            "-test.benchtime=100ms", f"-test.count={args.count}", "-test.benchmem", "-test.cpu=1",
        ], env=environment, text=True)
        repetitions = {}
        for record in parse_results(output):
            name = record["benchmark"]
            repetitions[name] = repetitions.get(name, 0) + 1
            record["repetition"] = repetitions[name]
            records.append(record)
        if len(repetitions) != 6 or any(count != args.count for count in repetitions.values()):
            raise RuntimeError(f"unexpected Fields results:\n{output}")

        cases = []
        for shape in ("Flat", "Nested"):
            for mode in ("Cached", "Uncached"):
                cases.append(("conversion", f"BenchmarkSourcePlan/{shape}/{mode}", 2000000))
        for platform in ("kakaowork", "slack", "googlechat"):
            for shape in ("Flat", "Nested", "Dynamic"):
                for mode in ("Reflection", "Generated"):
                    cases.append(("message", f"BenchmarkSourceFields/{platform}/{shape}/{mode}", 500000))
        for messenger in ("Kakaowork", "Slack", "GoogleChat"):
            for shape in ("Flat", "Nested", "Dynamic"):
                for mode in ("Reflection", "Generated"):
                    cases.append(("message", f"BenchmarkTo{messenger}Message/{shape}/{mode}", 100000))
        for repetition in range(1, args.count + 1):
            for binary, name, iterations in cases:
                result = measure(directory / binary, name, iterations, environment, directory)
                result["repetition"] = repetition
                records.append(result)
                print(f"{repetition}/{args.count} {name}: CPU {result['total_cpu_seconds']:.4f}s", flush=True)

    with args.output.open("w", newline="") as output:
        writer = csv.DictWriter(output, fieldnames=COLUMNS, lineterminator="\n")
        writer.writeheader()
        writer.writerows(records)
    print(f"Wrote {len(records)} samples to {args.output}")


if __name__ == "__main__":
    main()
