package main

import "fmt"

// やや実装のあるジェネリック関数。メソッド呼び出し＋ループ＋文字列連結を含む。
func Process[T fmt.Stringer](items []T) string {
	s := "items:\n"
	for i, item := range items {
		s += fmt.Sprintf("  [%d] %s\n", i, item.String())
	}
	return s
}

// 10 種類の異なるポインタ型（すべて Stringer を満たす）
type P1 struct{ v int }
type P2 struct{ v int }
type P3 struct{ v int }
type P4 struct{ v int }
type P5 struct{ v int }
type P6 struct{ v int }
type P7 struct{ v int }
type P8 struct{ v int }
type P9 struct{ v int }
type P10 struct{ v int }

func (p *P1) String() string  { return fmt.Sprintf("P1(%d)", p.v) }
func (p *P2) String() string  { return fmt.Sprintf("P2(%d)", p.v) }
func (p *P3) String() string  { return fmt.Sprintf("P3(%d)", p.v) }
func (p *P4) String() string  { return fmt.Sprintf("P4(%d)", p.v) }
func (p *P5) String() string  { return fmt.Sprintf("P5(%d)", p.v) }
func (p *P6) String() string  { return fmt.Sprintf("P6(%d)", p.v) }
func (p *P7) String() string  { return fmt.Sprintf("P7(%d)", p.v) }
func (p *P8) String() string  { return fmt.Sprintf("P8(%d)", p.v) }
func (p *P9) String() string  { return fmt.Sprintf("P9(%d)", p.v) }
func (p *P10) String() string { return fmt.Sprintf("P10(%d)", p.v) }

// 10 種類の異なる値型（structの underlying type が違うので別 shape）
type V1 struct{ v int }
type V2 struct{ v int }
type V3 struct{ v int }
type V4 struct{ v int }
type V5 struct{ v int }
type V6 struct{ v int }
type V7 struct{ v int }
type V8 struct{ v int }
type V9 struct{ v int }
type V10 struct{ v int }

func (v V1) String() string  { return fmt.Sprintf("V1(%d)", v.v) }
func (v V2) String() string  { return fmt.Sprintf("V2(%d)", v.v) }
func (v V3) String() string  { return fmt.Sprintf("V3(%d)", v.v) }
func (v V4) String() string  { return fmt.Sprintf("V4(%d)", v.v) }
func (v V5) String() string  { return fmt.Sprintf("V5(%d)", v.v) }
func (v V6) String() string  { return fmt.Sprintf("V6(%d)", v.v) }
func (v V7) String() string  { return fmt.Sprintf("V7(%d)", v.v) }
func (v V8) String() string  { return fmt.Sprintf("V8(%d)", v.v) }
func (v V9) String() string  { return fmt.Sprintf("V9(%d)", v.v) }
func (v V10) String() string { return fmt.Sprintf("V10(%d)", v.v) }

func main() {
	// ポインタ型 10 種類で具体化
	fmt.Print(Process([]*P1{{1}}))
	fmt.Print(Process([]*P2{{2}}))
	fmt.Print(Process([]*P3{{3}}))
	fmt.Print(Process([]*P4{{4}}))
	fmt.Print(Process([]*P5{{5}}))
	fmt.Print(Process([]*P6{{6}}))
	fmt.Print(Process([]*P7{{7}}))
	fmt.Print(Process([]*P8{{8}}))
	fmt.Print(Process([]*P9{{9}}))
	fmt.Print(Process([]*P10{{10}}))

	// 値型 10 種類で具体化
	fmt.Print(Process([]V1{{1}}))
	fmt.Print(Process([]V2{{2}}))
	fmt.Print(Process([]V3{{3}}))
	fmt.Print(Process([]V4{{4}}))
	fmt.Print(Process([]V5{{5}}))
	fmt.Print(Process([]V6{{6}}))
	fmt.Print(Process([]V7{{7}}))
	fmt.Print(Process([]V8{{8}}))
	fmt.Print(Process([]V9{{9}}))
	fmt.Print(Process([]V10{{10}}))

	diffMain()
}
