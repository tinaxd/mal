package main

import (
	"errors"
	"fmt"
)

type Env map[string]MalValue

func InitialEnv() Env {
	env := make(Env)
	makeFunc :=
		func(f func(int64, int64) int64) MalValue {
			return MalFunc{F: func(args []MalValue) (MalValue, error) {
				if len(args) != 2 {
					return nil, errors.New("wrong number of arguments")
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
	env["+"] = makeFunc(func(a, b int64) int64 {
		return a + b
	})
	env["-"] = makeFunc(func(a, b int64) int64 {
		return a - b
	})
	env["*"] = makeFunc(func(a, b int64) int64 {
		return a * b
	})
	env["/"] = makeFunc(func(a, b int64) int64 {
		return a / b
	})
	return env
}

func EvalAst(ast MalValue, env Env) (MalValue, error) {
	switch a := ast.(type) {
	case MalSymbol:
		v, ok := env[a.Value]
		if !ok {
			return nil, errors.New("symbol not found")
		}
		return v, nil
	case MalList:
		vals := make([]MalValue, len(a.Values))
		for i, v := range a.Values {
			val, err := eval(v, env)
			if err != nil {
				return nil, err
			}
			vals[i] = val
		}
		return MalList{Values: vals}, nil
	default:
		return ast, nil
	}
}
