package main

import "fmt"

// 検証1：本体が shape ごとに分かれるか
func First[T any](xs []T) T {
	return xs[0]
}

// 検証2：書き込み時、ポインタ系と非ポインタ系でライトバリアの有無が違うか
func StoreFirst[T any](xs []T, x T) {
	xs[0] = x
}

// 検証3：メソッド呼び出しが辞書経由になるか
type Stringer interface {
	String() string
}

func Greet[T Stringer](x T) string {
	return "Hi, " + x.String()
}

type UserA struct{ name string }
type UserB struct{ name string }

func (u *UserA) String() string { return u.name }
func (u *UserB) String() string { return u.name }

func main() {
	// 各 shape に該当する具体型で具体化を起こさせる
	pInts := []*int{new(int)}
	pStrs := []*string{new(string)}
	ints := []int{0}

	fmt.Println(First(pInts))
	fmt.Println(First(pStrs))
	fmt.Println(First(ints))

	StoreFirst(pInts, new(int))
	StoreFirst(pStrs, new(string))
	StoreFirst(ints, 42)

	fmt.Println(Greet(&UserA{name: "A"}))
	fmt.Println(Greet(&UserB{name: "B"}))
}
