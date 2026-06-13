package main

import "fmt"

// 検証4：複数メソッドの境界 interface — dict に何個入るか
type RW interface {
	Read() string
	Write(string)
}

func DoRW[T RW](x T) {
	s := x.Read()
	x.Write(s)
}

// 検証5：本体が make[T] するケース — 型情報が要る
func MakeAndFirst[T any](n int) T {
	xs := make([]T, n)
	var zero T
	if n > 0 {
		xs[0] = zero
	}
	return xs[0]
}

// 検証6：ネストしたジェネリック呼び出し — サブ辞書が要る
func Outer[T fmt.Stringer](x T) string {
	return Inner(x)
}

func Inner[T fmt.Stringer](x T) string {
	return x.String()
}

type Box struct{ v int }

func (b *Box) Read() string   { return fmt.Sprint(b.v) }
func (b *Box) Write(s string) { fmt.Println("write:", s) }

func init() {
	// 各具体化を強制
	b := &Box{v: 1}
	DoRW(b)

	xs := MakeAndFirst[int](3)
	fmt.Println(xs)

	u := &UserA{name: "X"}
	fmt.Println(Outer(u))
}
