package main

import (
	"fmt"
)

func InitialEnv() *Env {
	env := NewEnv(nil, nil, nil)

	ns := DefaultNamespace()
	for k, v := range ns.M {
		env.Set(k.Value, v)
	}

	return env
}

func EvalAst(ast MalValue, replEnv *Env, env *Env) (MalValue, error) {
	switch a := ast.(type) {
	case MalSymbol:
		v, ok := env.Get(a.Value)
		if !ok {
			return nil, fmt.Errorf("symbol '%s' not found", a.Value)
		}
		return v, nil
	case MalList:
		vals := make([]MalValue, len(a.Values))
		for i, v := range a.Values {
			val, err := eval(v, replEnv, env)
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
