package main

type MalValue struct {
	Type  MalValueType
	Value interface{}
}

type MalValueType int

const (
	TMalList   MalValueType = iota // []MalValue
	TMalInt                        // int
	TMalSymbol                     // string
)
