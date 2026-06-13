package main

import "fmt"

// 各 D 型は構造が違う（フィールド名や型や数が違う）→ underlying type が違う
type D1 struct{ a int }
type D2 struct{ b int }
type D3 struct {
	x, y int
}
type D4 struct{ s string }
type D5 struct {
	x int
	y string
}

func (d D1) String() string { return fmt.Sprintf("D1(%d)", d.a) }
func (d D2) String() string { return fmt.Sprintf("D2(%d)", d.b) }
func (d D3) String() string { return fmt.Sprintf("D3(%d,%d)", d.x, d.y) }
func (d D4) String() string { return fmt.Sprintf("D4(%s)", d.s) }
func (d D5) String() string { return fmt.Sprintf("D5(%d,%s)", d.x, d.y) }

func diffMain() {
	fmt.Print(Process([]D1{{1}}))
	fmt.Print(Process([]D2{{2}}))
	fmt.Print(Process([]D3{{1, 2}}))
	fmt.Print(Process([]D4{{"a"}}))
	fmt.Print(Process([]D5{{1, "b"}}))
}
