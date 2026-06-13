package main

import (
	"io"
	"strings"
)

// PlanetScale「Generics can make your Go code slower」(2022-03-30) の Escape ベンチ
// 元記事: https://planetscale.com/blog/generics-can-make-your-go-code-slower
// Wayback: https://web.archive.org/web/20220522015937/https://planetscale.com/blog/generics-can-make-your-go-code-slower
//
// 重要：Mono/Iface/GenericWithPointer/GenericWithExactIface/GenericWithSuperIface の差は
// 「関数シグネチャ」ではなく「呼び出し時の引数型」。同じ generic 関数を異なる interface 経由で呼ぶ。

// IBuffer：io.ByteWriter よりも広い interface（記事の定義）
type IBuffer interface {
	Write([]byte) (int, error)
	WriteByte(c byte) error
	Len() int
}

// 共通：エスケープ対象のテストデータ
var testInput = []byte(strings.Repeat("hello 'world' with `quotes` and \\backslash ", 50))

// 1. Monomorphized: *strings.Builder を直接受ける
func EscapeMonomorphized(w *strings.Builder, input []byte) {
	for _, b := range input {
		switch b {
		case '\'':
			w.WriteByte('\\')
			w.WriteByte('\'')
		case '\\':
			w.WriteByte('\\')
			w.WriteByte('\\')
		default:
			w.WriteByte(b)
		}
	}
}

// 2. Iface: io.ByteWriter で受ける
func EscapeIface(w io.ByteWriter, input []byte) {
	for _, b := range input {
		switch b {
		case '\'':
			w.WriteByte('\\')
			w.WriteByte('\'')
		case '\\':
			w.WriteByte('\\')
			w.WriteByte('\\')
		default:
			w.WriteByte(b)
		}
	}
}

// 3-5: 同じ generic 関数。呼び出し時の引数型で挙動が変わる
func EscapeGeneric[W io.ByteWriter](w W, input []byte) {
	for _, b := range input {
		switch b {
		case '\'':
			w.WriteByte('\\')
			w.WriteByte('\'')
		case '\\':
			w.WriteByte('\\')
			w.WriteByte('\\')
		default:
			w.WriteByte(b)
		}
	}
}

// SuperIface 用：制約は同じ io.ByteWriter だが、IBuffer を渡せるように
// （W の推論を IBuffer 側に倒すための関数）
func EscapeGenericSuper[W IBuffer](w W, input []byte) {
	for _, b := range input {
		switch b {
		case '\'':
			w.WriteByte('\\')
			w.WriteByte('\'')
		case '\\':
			w.WriteByte('\\')
			w.WriteByte('\\')
		default:
			w.WriteByte(b)
		}
	}
}
