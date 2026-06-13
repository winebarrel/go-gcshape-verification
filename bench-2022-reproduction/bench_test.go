package main

import (
	"io"
	"strings"
	"testing"
)

// === DoltHub Search ベンチ ===

var elems []uint64

func init() {
	elems = make([]uint64, 500_000)
	for i := range elems {
		elems[i] = uint64(i)
	}
}

func BenchmarkSearch(b *testing.B) {
	b.Run("interface_search_value_type", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			arr := ArrayVal{elems: elems}
			target := uint64(i % arr.Len())
			InterfaceSearch(target, arr)
		}
	})

	b.Run("generic_search_value_type", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			arr := ArrayVal{elems: elems}
			target := uint64(i % arr.Len())
			GenericSearch(target, arr)
		}
	})

	b.Run("interface_search_reference_type", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			arr := &ArrayRef{elems: elems}
			target := uint64(i % arr.Len())
			InterfaceSearch(target, arr)
		}
	})

	b.Run("generic_search_reference_type", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			arr := &ArrayRef{elems: elems}
			target := uint64(i % arr.Len())
			GenericSearch(target, arr)
		}
	})
}

// === PlanetScale Escape ベンチ ===

func BenchmarkEscape(b *testing.B) {
	var sb strings.Builder
	sb.Grow(len(testInput) * 2)

	b.Run("Monomorphized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sb.Reset()
			EscapeMonomorphized(&sb, testInput)
		}
	})

	b.Run("Iface", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sb.Reset()
			EscapeIface(&sb, testInput)
		}
	})

	b.Run("GenericWithPointer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sb.Reset()
			EscapeGeneric(&sb, testInput) // W は *strings.Builder と推論される
		}
	})

	b.Run("GenericWithExactIface", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sb.Reset()
			var w io.ByteWriter = &sb
			EscapeGeneric(w, testInput) // W は io.ByteWriter と推論される
		}
	})

	b.Run("GenericWithSuperIface", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sb.Reset()
			var w IBuffer = &sb
			EscapeGenericSuper(w, testInput) // W は IBuffer と推論される（assertI2I 経路）
		}
	})
}
