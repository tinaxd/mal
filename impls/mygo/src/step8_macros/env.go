package main

type Env struct {
	M     map[string]MalValue
	Outer *Env
}

func NewEnv(outer *Env, binds []string, exprs []MalValue) *Env {
	env := &Env{M: make(map[string]MalValue), Outer: outer}

	for i, bind := range binds {
		if bind == "&" {
			// set the rest
			env.Set(binds[i+1], MalList{Values: exprs[i:]})
			break
		}
		env.Set(bind, exprs[i])
	}

	return env
}

func (e *Env) Set(key string, val MalValue) {
	e.M[key] = val
}

func (e *Env) Find(key string) (*Env, bool) {
	if _, ok := e.M[key]; ok {
		return e, true
	}
	if e.Outer == nil {
		return nil, false
	}
	return e.Outer.Find(key)
}

func (e *Env) Get(key string) (MalValue, bool) {
	env, ok := e.Find(key)
	if !ok {
		return nil, false
	}
	return env.M[key], true
}
