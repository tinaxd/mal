package main

type MalValue interface {
	MalValue()
}

type MalValueType int

const (
	TMalList   MalValueType = iota // []MalValue
	TMalInt                        // int
	TMalSymbol                     // string
)

type MalList struct {
	Values []MalValue
}

func (MalList) MalValue() {}

type MalInt struct {
	Value int64
}

func (MalInt) MalValue() {}

type MalSymbol struct {
	Value string
}

func (MalSymbol) MalValue() {}
