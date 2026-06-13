#!/usr/bin/env python3
"""Read `go test -bench -benchmem` output from stdin and print a markdown summary.

For each benchmark (e.g. BenchmarkSearch, BenchmarkEscape) the script averages
the runs of every sub-benchmark and produces a table ordered from fastest to
slowest, with the fastest treated as the baseline.
"""
import re
import sys
from collections import defaultdict, OrderedDict


LINE = re.compile(
    r"^(Benchmark[^/]+)/([\w_]+)-\d+\s+\d+\s+([\d.]+)\s+ns/op"
    r"(?:\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op)?"
)


def main() -> None:
    ns_runs: dict = defaultdict(list)
    bytes_runs: dict = defaultdict(list)
    allocs_runs: dict = defaultdict(list)
    top_order: list = []

    for line in sys.stdin:
        m = LINE.match(line)
        if not m:
            continue
        top, sub = m.group(1), m.group(2)
        ns = float(m.group(3))
        ns_runs[(top, sub)].append(ns)
        if m.group(4) is not None:
            bytes_runs[(top, sub)].append(int(m.group(4)))
            allocs_runs[(top, sub)].append(int(m.group(5)))
        if top not in top_order:
            top_order.append(top)

    has_mem = any(bytes_runs.values())

    for top in top_order:
        entries = []
        for (t, sub), times in ns_runs.items():
            if t != top:
                continue
            avg_ns = sum(times) / len(times)
            avg_b = (
                sum(bytes_runs[(t, sub)]) / len(bytes_runs[(t, sub)])
                if (t, sub) in bytes_runs
                else None
            )
            avg_a = (
                sum(allocs_runs[(t, sub)]) / len(allocs_runs[(t, sub)])
                if (t, sub) in allocs_runs
                else None
            )
            entries.append((sub, avg_ns, avg_b, avg_a))

        entries.sort(key=lambda e: e[1])
        if not entries:
            continue
        fastest = entries[0][1]

        print(f"### {top}\n")
        if has_mem:
            print("| Variant | ns/op | B/op | allocs/op | Relative |")
            print("|---|---:|---:|---:|---|")
        else:
            print("| Variant | ns/op | Relative |")
            print("|---|---:|---|")

        for i, (sub, ns, b, a) in enumerate(entries):
            ratio = ns / fastest
            if i == 0:
                rel = "baseline"
            else:
                rel = f"{ratio:.2f}x slower"
            if has_mem:
                b_str = f"{b:.0f}" if b is not None else ""
                a_str = f"{a:.0f}" if a is not None else ""
                print(f"| `{sub}` | {ns:.2f} | {b_str} | {a_str} | {rel} |")
            else:
                print(f"| `{sub}` | {ns:.2f} | {rel} |")
        print()


if __name__ == "__main__":
    main()
