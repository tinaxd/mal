package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrWrongFuncNArgs = errors.New("wrong number of arguments")
)

type Namespace struct {
	M map[MalSymbol]MalFunc
}

func makeSymbol(s string) MalSymbol {
	return MalSymbol{Value: s}
}

func makeFunc(f func([]MalValue) (MalValue, error)) MalFunc {
	return MalFunc{F: f}
}

func DefaultNamespace() Namespace {
	m := make(map[MalSymbol]MalFunc)
	makeF :=
		func(f func(int64, int64) int64) MalFunc {
			return MalFunc{F: func(args []MalValue) (MalValue, error) {
				if len(args) != 2 {
					return nil, ErrWrongFuncNArgs
				}
				a, ok := args[0].(MalInt)
				if !ok {
					return nil, fmt.Errorf("first argument is not an integer: %v", args[0])
				}
				b, ok := args[1].(MalInt)
				if !ok {
					return nil, fmt.Errorf("second argument is not an integer: %v", args[1])
				}
				return MalInt{Value: f(a.Value, b.Value)}, nil
			}}
		}
	m[makeSymbol("+")] = makeF(func(a, b int64) int64 {
		return a + b
	})
	m[makeSymbol("-")] = makeF(func(a, b int64) int64 {
		return a - b
	})
	m[makeSymbol("*")] = makeF(func(a, b int64) int64 {
		return a * b
	})
	m[makeSymbol("/")] = makeF(func(a, b int64) int64 {
		return a / b
	})
	m[makeSymbol("prn")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		s := PrStr(args[0], true)
		fmt.Println(s)
		return nil, nil
	})
	m[makeSymbol("list")] = makeFunc(func(args []MalValue) (MalValue, error) {
		values := make([]MalValue, len(args))
		copy(values, args)
		return MalList{Values: values}, nil
	})
	m[makeSymbol("list?")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		if args[0] == nil {
			return MalBool{Value: true}, nil
		}
		_, ok := args[0].(MalList)
		return MalBool{Value: ok}, nil
	})
	m[makeSymbol("empty?")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		if args[0] == nil {
			return MalBool{Value: true}, nil
		}
		_, ok := args[0].(MalList)
		if !ok {
			return nil, fmt.Errorf("expected MalList, got %v", args[0])
		}
		return MalBool{Value: len(args[0].(MalList).Values) == 0}, nil
	})
	m[makeSymbol("count")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		if args[0] == nil {
			return MalInt{Value: 0}, nil
		}
		l, ok := args[0].(MalList)
		if !ok {
			return nil, fmt.Errorf("expected MalList, got %v", args[0])
		}
		return MalInt{Value: int64(len(l.Values))}, nil
	})
	m[makeSymbol("=")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 2 {
			return nil, ErrWrongFuncNArgs
		}
		return MalBool{Value: malEq(args[0], args[1])}, nil
	})

	mkCmp := func(f func(int64, int64) bool) MalFunc {
		return makeFunc(func(args []MalValue) (MalValue, error) {
			if len(args) != 2 {
				return nil, ErrWrongFuncNArgs
			}
			a, ok := args[0].(MalInt)
			if !ok {
				return nil, fmt.Errorf("first argument is not an integer: %v", args[0])
			}
			b, ok := args[1].(MalInt)
			if !ok {
				return nil, fmt.Errorf("second argument is not an integer: %v", args[1])
			}
			return MalBool{Value: f(a.Value, b.Value)}, nil
		})
	}

	m[makeSymbol("<")] = mkCmp(func(a, b int64) bool { return a < b })
	m[makeSymbol("<=")] = mkCmp(func(a, b int64) bool { return a <= b })
	m[makeSymbol(">")] = mkCmp(func(a, b int64) bool { return a > b })
	m[makeSymbol(">=")] = mkCmp(func(a, b int64) bool { return a >= b })

	m[makeSymbol("read-string")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		s, ok := args[0].(MalString)
		if !ok {
			return nil, fmt.Errorf("expected MalString, got %v", args[0])
		}
		return ReadStr(s.Value)
	})

	m[makeSymbol("slurp")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		s, ok := args[0].(MalString)
		if !ok {
			return nil, fmt.Errorf("expected MalString, got %v", args[0])
		}

		// open file with filename s
		// read file contents
		f, err := os.Open(s.Value)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		return MalString{Value: string(content)}, nil
	})

	m[makeSymbol("pr-str")] = makeFunc(func(args []MalValue) (MalValue, error) {
		strs := make([]string, len(args))
		for i, a := range args {
			strs[i] = PrStr(a, true)
		}
		str := strings.Join(strs, " ")
		return MalString{Value: str}, nil
	})
	m[makeSymbol("str")] = makeFunc(func(args []MalValue) (MalValue, error) {
		strs := make([]string, len(args))
		for i, a := range args {
			strs[i] = PrStr(a, false)
		}
		str := strings.Join(strs, "")
		return MalString{Value: str}, nil
	})
	m[makeSymbol("prn")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		s := PrStr(args[0], true)
		fmt.Println(s)
		return nil, nil
	})
	m[makeSymbol("println")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		s := PrStr(args[0], false)
		fmt.Println(s)
		return nil, nil
	})

	m[makeSymbol("atom")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		return NewMalAtom(args[0]), nil
	})
	m[makeSymbol("atom?")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		_, ok := args[0].(*MalAtom)
		return MalBool{Value: ok}, nil
	})
	m[makeSymbol("deref")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		a, ok := args[0].(*MalAtom)
		if !ok {
			return nil, fmt.Errorf("expected MalAtom, got %v", args[0])
		}
		return a.Ref, nil
	})
	m[makeSymbol("reset!")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 2 {
			return nil, ErrWrongFuncNArgs
		}
		a, ok := args[0].(*MalAtom)
		if !ok {
			return nil, fmt.Errorf("expected MalAtom, got %v", args[0])
		}
		a.Ref = args[1]
		return args[1], nil
	})
	m[makeSymbol("swap!")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) < 2 {
			return nil, ErrWrongFuncNArgs
		}
		a, ok := args[0].(*MalAtom)
		if !ok {
			return nil, fmt.Errorf("expected MalAtom, got %v", args[0])
		}
		f, ok := args[1].(MalInvoke)
		if !ok {
			return nil, fmt.Errorf("expected MalFunc, got %v", args[1])
		}
		fArgs := make([]MalValue, len(args)-1)
		fArgs[0] = a.Ref
		copy(fArgs[1:], args[2:])
		newVal, err := f.Invoke(fArgs)
		if err != nil {
			return nil, err
		}
		a.Ref = newVal
		return newVal, nil
	})

	return Namespace{M: m}
}
