package main

// DoltHub「Generics and Value Types in Golang」(2022-04-01) の Search ベンチ再現
// 元記事: https://www.dolthub.com/blog/2022-04-01-fast-generics/

type Array interface {
	Get(i int) uint64
	Len() int
}

type ArrayVal struct {
	elems []uint64
}

var _ Array = ArrayVal{}

func (v ArrayVal) Get(i int) uint64 {
	return v.elems[i]
}

func (v ArrayVal) Len() int {
	return len(v.elems)
}

type ArrayRef struct {
	elems []uint64
}

var _ Array = &ArrayRef{}

func (p *ArrayRef) Get(i int) uint64 {
	return p.elems[i]
}

func (p *ArrayRef) Len() int {
	return len(p.elems)
}

// Interface-based
func InterfaceSearch(target uint64, arr Array) int {
	i, j := 0, arr.Len()
	for i < j {
		h := int(uint(i+j) >> 1)
		if arr.Get(h) < target {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

// Generic
func GenericSearch[T Array](target uint64, arr T) int {
	i, j := 0, arr.Len()
	for i < j {
		h := int(uint(i+j) >> 1)
		if arr.Get(h) < target {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
