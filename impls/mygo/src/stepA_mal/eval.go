package main

import (
	"fmt"
)

func InitialEnv() *Env {
	env, err := NewEnv(nil, nil, nil)
	if err != nil {
		panic("unreachable")
	}

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
			return nil, fmt.Errorf("'%s' not found", a.Value)
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
		return MalList{Values: vals, Vector: a.Vector}, nil
	case *MalMap:
		kvs := make([]MalValue, 0)
		for _, kv := range a.Iter() {
			kvs = append(kvs, kv.Key)

			v, err := eval(kv.Value, replEnv, env)
			if err != nil {
				return nil, err
			}
			kvs = append(kvs, v)
		}
		return NewMapFromList(kvs)
	default:
		return ast, nil
	}
}
