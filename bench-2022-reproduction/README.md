# bench-2022-reproduction

Reproduces two well-known 2022 generics benchmarks on Go 1.26.1 and compares the numbers with the originals.

The two articles are:

- PlanetScale, Vicent Marti, [Generics can make your Go code slower](https://planetscale.com/blog/generics-can-make-your-go-code-slower) (2022-03-30)
- DoltHub, Andy Arthur, [Generics and Value Types in Golang](https://www.dolthub.com/blog/2022-04-01-fast-generics/) (2022-04-01)

## Running

```sh
go test -bench=. -benchmem -benchtime=2s -count=3
```

The repository also has a manual [Benchmark workflow](https://github.com/winebarrel/go-gcshape-verification/actions/workflows/bench.yml). Running it runs the same command on `ubuntu-latest` or `macos-latest` (selectable as an input) and prints the results in the job summary.

## Files

`dolthub.go` does binary search over `ArrayVal` (a value receiver) and `ArrayRef` (a pointer receiver), through an `Array` interface and through an `Array`-constrained generic function. The code is taken directly from the article.

`planetscale.go` defines `EscapeMonomorphized`, `EscapeIface`, `EscapeGeneric[W io.ByteWriter]`, and `EscapeGenericSuper[W IBuffer]`. The original article is a single-page application and the source could not be extracted directly, so this file is reconstructed from the prose. The structure matches what the article describes, but small differences in the function body are possible.

`bench_test.go` drives the five PlanetScale variants and the four DoltHub variants.

The five PlanetScale variants exercise the following:

| Name | What is measured |
|---|---|
| `Monomorphized` | `func Escape(*strings.Builder, []byte)`. Concrete type, no generics, no interface. |
| `Iface` | `func Escape(io.ByteWriter, []byte)`. Interface, no generics. |
| `GenericWithPointer` | `EscapeGeneric[W io.ByteWriter](w W, ...)` called with `*strings.Builder`. W is inferred as the concrete pointer type. |
| `GenericWithExactIface` | The same function, called with the value already typed as `io.ByteWriter`. W is inferred as `io.ByteWriter` itself. |
| `GenericWithSuperIface` | A second generic function constrained on the wider interface `IBuffer`, called with the value typed as `IBuffer`. This triggers the `runtime.assertI2I` path that the original article identified as the slowest case. |

## Headline result

On Go 1.26.1, darwin/arm64 (Apple M4 Pro), the cost ratios versus `Monomorphized` are quite different from the 2022 PlanetScale numbers.

| Variant | 2022 (PlanetScale) | 2026 (this repo) |
|---|---|---|
| Monomorphized | 1.00x | 1.00x |
| Iface | 1.35x | about 2.05x |
| GenericWithPointer | 1.42x | about 2.07x |
| GenericWithExactIface | 1.91x | about 1.41x |
| GenericWithSuperIface | 3.48x | about 1.41x |

The `GenericWithSuperIface` case in particular has shrunk from 3.48x to about 1.41x, which suggests the runtime's `assertI2I` path was optimised significantly between Go 1.18 and 1.26.

The DoltHub observation reproduces. Generic search over a value type runs at about 25% less time than interface-based search over the same value type. The allocation column makes the mechanism explicit: the generic value-type case is `0 B/op`, the others are `24 B/op`.

## Notes

The PlanetScale reconstruction is structural, not byte-for-byte. Treat it as the same setup, not an exact replay of the original code.

Absolute timings depend heavily on the CPU. M4 Pro is fast and modern, and the ratio against `Monomorphized` is what carries meaning across platforms.

The `strings.Builder` inside `EscapeMonomorphized` triggers internal resizing, so the allocation column here does not match the original article's allocation column. The ordering is still informative, but treat the absolute alloc numbers with caution.
