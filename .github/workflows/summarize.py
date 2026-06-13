#!/usr/bin/env python3
"""Read one or more `go test -bench -benchmem` outputs and print a markdown summary.

Usage:
    summarize.py                          # single input from stdin
    summarize.py LABEL=FILE [LABEL=FILE]  # multiple inputs, combined table per benchmark

For each top-level benchmark (e.g. BenchmarkSearch) the runs of each sub-
benchmark are averaged. Within each input the fastest sub-benchmark becomes
the baseline. With multiple inputs the table places one ns/op + Relative
column pair per input.
"""
import re
import sys
from collections import defaultdict


LINE = re.compile(
    r"^(Benchmark[^/]+)/([\w_]+)-\d+\s+\d+\s+([\d.]+)\s+ns/op"
    r"(?:\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op)?"
)


def parse(stream):
    ns_runs = defaultdict(list)
    bytes_runs = defaultdict(list)
    allocs_runs = defaultdict(list)
    top_order = []

    for line in stream:
        m = LINE.match(line)
        if not m:
            continue
        top, sub = m.group(1), m.group(2)
        ns_runs[(top, sub)].append(float(m.group(3)))
        if m.group(4) is not None:
            bytes_runs[(top, sub)].append(int(m.group(4)))
            allocs_runs[(top, sub)].append(int(m.group(5)))
        if top not in top_order:
            top_order.append(top)

    avgs = {}
    for key, times in ns_runs.items():
        avgs[key] = {
            "ns": sum(times) / len(times),
            "bytes": (sum(bytes_runs[key]) / len(bytes_runs[key])
                      if key in bytes_runs else None),
            "allocs": (sum(allocs_runs[key]) / len(allocs_runs[key])
                       if key in allocs_runs else None),
        }
    return top_order, avgs


def fmt_rel(ratio):
    if abs(ratio - 1.0) < 0.005:
        return "baseline"
    return f"{ratio:.2f}x slower"


def print_single(top_order, avgs):
    has_mem = any(v["bytes"] is not None for v in avgs.values())
    for top in top_order:
        print(f"### {top}\n")
        if has_mem:
            print("| Variant | ns/op | B/op | allocs/op | Relative |")
            print("|---|---:|---:|---:|---|")
        else:
            print("| Variant | ns/op | Relative |")
            print("|---|---:|---|")

        entries = [(s, avgs[(top, s)]) for (t, s) in avgs if t == top]
        entries.sort(key=lambda e: e[1]["ns"])
        if not entries:
            continue
        baseline = entries[0][1]["ns"]
        for sub, v in entries:
            rel = fmt_rel(v["ns"] / baseline)
            if has_mem:
                b = f"{v['bytes']:.0f}" if v["bytes"] is not None else ""
                a = f"{v['allocs']:.0f}" if v["allocs"] is not None else ""
                print(f"| `{sub}` | {v['ns']:.2f} | {b} | {a} | {rel} |")
            else:
                print(f"| `{sub}` | {v['ns']:.2f} | {rel} |")
        print()


def print_combined(sources):
    all_tops = []
    for _, top_order, _ in sources:
        for t in top_order:
            if t not in all_tops:
                all_tops.append(t)

    for top in all_tops:
        print(f"### {top}\n")

        header = ["Variant"]
        sep = ["---"]
        for label, _, _ in sources:
            header += [f"{label} ns/op", f"{label} Relative"]
            sep += ["---:", "---"]
        print("| " + " | ".join(header) + " |")
        print("|" + "|".join(sep) + "|")

        baselines = []
        for _, _, avgs in sources:
            ns_values = [v["ns"] for (t, _), v in avgs.items() if t == top]
            baselines.append(min(ns_values) if ns_values else None)

        subs_seen = []
        for _, _, avgs in sources:
            for (t, s) in avgs:
                if t == top and s not in subs_seen:
                    subs_seen.append(s)

        first_avgs = sources[0][2]
        subs_seen.sort(
            key=lambda s: first_avgs.get((top, s), {}).get("ns", float("inf"))
        )

        for sub in subs_seen:
            cells = [f"`{sub}`"]
            for (label, _, avgs), baseline in zip(sources, baselines):
                v = avgs.get((top, sub))
                if v is None or baseline is None:
                    cells += ["-", "-"]
                else:
                    cells += [f"{v['ns']:.2f}", fmt_rel(v["ns"] / baseline)]
            print("| " + " | ".join(cells) + " |")
        print()


def main():
    args = sys.argv[1:]
    if not args:
        top_order, avgs = parse(sys.stdin)
        print_single(top_order, avgs)
        return

    sources = []
    for arg in args:
        label, _, path = arg.partition("=")
        if not path:
            print(
                f"Usage: {sys.argv[0]} LABEL=FILE [LABEL=FILE ...]",
                file=sys.stderr,
            )
            sys.exit(2)
        with open(path) as f:
            top_order, avgs = parse(f)
        sources.append((label, top_order, avgs))

    print_combined(sources)


if __name__ == "__main__":
    main()
