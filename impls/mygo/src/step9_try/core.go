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
		lst, ok := args[0].(MalList)
		return MalBool{Value: ok && !lst.IsVector()}, nil
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
		for i, arg := range args {
			s := PrStr(arg, true)

			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(s)
		}
		fmt.Println()
		return nil, nil
	})
	m[makeSymbol("println")] = makeFunc(func(args []MalValue) (MalValue, error) {
		for i, arg := range args {
			s := PrStr(arg, false)

			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(s)
		}
		fmt.Println()
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

	m[makeSymbol("cons")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 2 {
			return nil, ErrWrongFuncNArgs
		}
		l, ok := args[1].(MalList)
		if !ok {
			return nil, fmt.Errorf("expected MalList, got %v", args[1])
		}
		values := make([]MalValue, len(l.Values)+1)
		values[0] = args[0]
		copy(values[1:], l.Values)
		return MalList{Values: values}, nil
	})
	m[makeSymbol("concat")] = makeFunc(func(args []MalValue) (MalValue, error) {
		values := make([]MalValue, 0)
		for _, a := range args {
			l, ok := a.(MalList)
			if !ok {
				return nil, fmt.Errorf("expected MalList, got %v", a)
			}
			values = append(values, l.Values...)
		}
		return MalList{Values: values}, nil
	})
	m[makeSymbol("nth")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 2 {
			return nil, ErrWrongFuncNArgs
		}
		l, ok := args[0].(MalList)
		if !ok {
			return nil, fmt.Errorf("expected MalList, got %v", args[0])
		}
		i, ok := args[1].(MalInt)
		if !ok {
			return nil, fmt.Errorf("expected MalInt, got %v", args[1])
		}
		if i.Value < 0 || i.Value >= int64(len(l.Values)) {
			return nil, fmt.Errorf("index out of range: %d", i.Value)
		}
		return l.Values[i.Value], nil
	})
	m[makeSymbol("first")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		if args[0] == nil {
			return nil, nil
		}
		l, ok := args[0].(MalList)
		if !ok {
			return nil, fmt.Errorf("expected MalList, got %v", args[0])
		}
		if len(l.Values) == 0 {
			return nil, nil
		}
		return l.Values[0], nil
	})
	m[makeSymbol("rest")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		if args[0] == nil {
			return NewList([]MalValue{}), nil
		}
		l, ok := args[0].(MalList)
		if !ok {
			return nil, fmt.Errorf("expected MalList, got %v", args[0])
		}
		if len(l.Values) == 0 {
			return MalList{Values: []MalValue{}}, nil
		}
		return MalList{Values: l.Values[1:]}, nil
	})

	m[makeSymbol("throw")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		return nil, NewErrorFromValue(args[0])
	})

	m[makeSymbol("apply")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) < 2 {
			return nil, ErrWrongFuncNArgs
		}
		f, ok := args[0].(MalInvoke)
		if !ok {
			return nil, fmt.Errorf("expected MalFunc, got %v", args[0])
		}

		lastList, ok := args[len(args)-1].(MalList)
		if !ok {
			return nil, fmt.Errorf("expected MalList or MalVector, got %v", args[len(args)-1])
		}

		fArgs := make([]MalValue, 0)
		for i := 1; i < len(args)-1; i++ {
			fArgs = append(fArgs, args[i])
		}
		fArgs = append(fArgs, lastList.Values...)
		return f.Invoke(fArgs)
	})

	m[makeSymbol("map")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 2 {
			return nil, ErrWrongFuncNArgs
		}
		f, ok := args[0].(MalInvoke)
		if !ok {
			return nil, fmt.Errorf("expected MalFunc, got %v", args[0])
		}

		l, ok := args[1].(MalList)
		if !ok {
			return nil, fmt.Errorf("expected MalList or MalVector, got %v", args[1])
		}

		values := make([]MalValue, len(l.Values))
		for i, v := range l.Values {
			result, err := f.Invoke([]MalValue{v})
			if err != nil {
				return nil, err
			}
			values[i] = result
		}
		return MalList{Values: values}, nil
	})

	m[makeSymbol("symbol")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		s, ok := args[0].(MalString)
		if !ok {
			return nil, fmt.Errorf("expected MalString, got %v", args[0])
		}
		return makeSymbol(s.Value), nil
	})

	m[makeSymbol("keyword")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		s, ok := args[0].(MalString)
		if !ok {
			return nil, fmt.Errorf("expected MalString, got %v", args[0])
		}
		if s.IsKeyword() {
			return s, nil
		}
		return NewKeyword(s.Value), nil
	})

	m[makeSymbol("vector")] = makeFunc(func(args []MalValue) (MalValue, error) {
		values := make([]MalValue, len(args))
		copy(values, args)
		return NewVector(values), nil
	})

	m[makeSymbol("hash-map")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args)%2 != 0 {
			return nil, fmt.Errorf("expected even number of arguments, got %d", len(args))
		}
		return NewMapFromList(args)
	})

	m[makeSymbol("assoc")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("expected at least 3 arguments, got %d", len(args))
		}
		m, ok := args[0].(*MalMap)
		if !ok {
			return nil, fmt.Errorf("expected MalMap, got %v", args[0])
		}
		if len(args)%2 != 1 {
			return nil, fmt.Errorf("expected even number of arguments, got %d", len(args))
		}
		newMap := CloneMap(m)
		for i := 1; i < len(args); i += 2 {
			newMap.Set(args[i], args[i+1])
		}
		return newMap, nil
	})

	m[makeSymbol("dissoc")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) < 2 {
			return nil, ErrWrongFuncNArgs
		}
		m, ok := args[0].(*MalMap)
		if !ok {
			return nil, fmt.Errorf("expected MalMap, got %v", args[0])
		}

		newMap := CloneMap(m)
		for _, k := range args[1:] {
			newMap.Del(k)
		}
		return newMap, nil
	})

	m[makeSymbol("get")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 2 {
			return nil, ErrWrongFuncNArgs
		}
		if args[0] == nil {
			return nil, nil
		}
		m, ok := args[0].(*MalMap)
		if !ok {
			return nil, fmt.Errorf("expected MalMap, got %v", args[0])
		}
		v, ok := m.Get(args[1])
		if !ok {
			return nil, nil
		}
		return v, nil
	})

	m[makeSymbol("contains?")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 2 {
			return nil, ErrWrongFuncNArgs
		}
		m, ok := args[0].(*MalMap)
		if !ok {
			return nil, fmt.Errorf("expected MalMap, got %v", args[0])
		}
		_, ok = m.Get(args[1])
		return NewBool(ok), nil
	})

	m[makeSymbol("keys")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		m, ok := args[0].(*MalMap)
		if !ok {
			return nil, fmt.Errorf("expected MalMap, got %v", args[0])
		}

		keys := make([]MalValue, 0)
		for _, kv := range m.Iter() {
			keys = append(keys, kv.Key)
		}
		return NewList(keys), nil
	})

	m[makeSymbol("vals")] = makeFunc(func(args []MalValue) (MalValue, error) {
		if len(args) != 1 {
			return nil, ErrWrongFuncNArgs
		}
		m, ok := args[0].(*MalMap)
		if !ok {
			return nil, fmt.Errorf("expected MalMap, got %v", args[0])
		}

		vals := make([]MalValue, 0)
		for _, kv := range m.Iter() {
			vals = append(vals, kv.Value)
		}
		return NewList(vals), nil
	})

	onePred := func(f func(MalValue) bool) MalFunc {
		return makeFunc(func(args []MalValue) (MalValue, error) {
			if len(args) != 1 {
				return nil, ErrWrongFuncNArgs
			}
			return MalBool{Value: f(args[0])}, nil
		})
	}
	m[makeSymbol("nil?")] = onePred(func(v MalValue) bool {
		return v == nil
	})
	m[makeSymbol("true?")] = onePred(func(v MalValue) bool {
		b, ok := v.(MalBool)
		return ok && b.Value
	})
	m[makeSymbol("false?")] = onePred(func(v MalValue) bool {
		b, ok := v.(MalBool)
		return ok && !b.Value
	})
	m[makeSymbol("symbol?")] = onePred(func(v MalValue) bool {
		_, ok := v.(MalSymbol)
		return ok
	})
	m[makeSymbol("keyword?")] = onePred(func(v MalValue) bool {
		kw, ok := v.(MalString)
		return ok && kw.IsKeyword()
	})
	m[makeSymbol("vector?")] = onePred(func(v MalValue) bool {
		l, ok := v.(MalList)
		return ok && l.IsVector()
	})
	m[makeSymbol("sequential?")] = onePred(func(v MalValue) bool {
		_, ok := v.(MalList)
		return ok
	})
	m[makeSymbol("map?")] = onePred(func(v MalValue) bool {
		_, ok := v.(*MalMap)
		return ok
	})

	return Namespace{M: m}
}

func quasiquote(ast MalValue) (MalValue, error) {
	switch q := ast.(type) {
	case MalList:
		if len(q.Values) > 0 {
			sym, ok := q.Values[0].(MalSymbol)
			if ok && sym.Value == "unquote" {
				if len(q.Values) != 2 {
					return nil, fmt.Errorf("wrong number of arguments for unquote")
				}
				return q.Values[1], nil
			}

			result := MalList{Values: make([]MalValue, 0)}
			for i := len(q.Values) - 1; i >= 0; i-- {
				elt := q.Values[i]
				switch e := elt.(type) {
				case MalList:
					if len(e.Values) > 0 {
						sym, ok := e.Values[0].(MalSymbol)

						if ok && sym.Value == "splice-unquote" {
							if len(e.Values) != 2 {
								return nil, fmt.Errorf("wrong number of arguments for splice-unquote")
							}
							result = MalList{
								Values: []MalValue{
									makeSymbol("concat"),
									e.Values[1],
									result,
								},
							}
							continue
						}
					}
				}
				eltQuasi, err := quasiquote(elt)
				if err != nil {
					return nil, err
				}
				result = MalList{
					Values: []MalValue{
						makeSymbol("cons"),
						eltQuasi,
						result,
					},
				}
				continue
			}

			return result, nil
		}
	case MalSymbol:
		return MalList{
			Values: []MalValue{
				makeSymbol("quote"),
				ast,
			},
		}, nil
	}

	return ast, nil
}
