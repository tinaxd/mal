package main

type MalValue interface {
	MalValue()
}

type MalInvoke interface {
	Invoke([]MalValue) (MalValue, error)
	IsMacro() bool
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

type MalFunc struct {
	F     func([]MalValue) (MalValue, error)
	Macro bool
}

func (MalFunc) MalValue() {}
func (f MalFunc) Invoke(args []MalValue) (MalValue, error) {
	return f.F(args)
}
func (f MalFunc) IsMacro() bool {
	return f.Macro
}

type MalTcoFunc struct {
	Ast    MalValue
	Params []string
	Env    *Env
	Fn     MalFunc
}

func (MalTcoFunc) MalValue() {}
func (f MalTcoFunc) Invoke(args []MalValue) (MalValue, error) {
	return f.Fn.F(args)
}
func (f MalTcoFunc) IsMacro() bool {
	return f.Fn.Macro
}

type MalBool struct {
	Value bool
}

func (MalBool) MalValue() {}

type MalString struct {
	Value string
}

func (MalString) MalValue() {}

type MalAtom struct {
	Ref MalValue
}

func (*MalAtom) MalValue() {}
func NewMalAtom(v MalValue) *MalAtom {
	return &MalAtom{Ref: v}
}

func malEq(v1 MalValue, v2 MalValue) bool {
	if v1 == nil {
		return v2 == nil
	}

	switch v1 := v1.(type) {
	case MalInt:
		v2, ok := v2.(MalInt)
		if !ok {
			return false
		}
		return v1.Value == v2.Value
	case MalSymbol:
		v2, ok := v2.(MalSymbol)
		if !ok {
			return false
		}
		return v1.Value == v2.Value
	case MalFunc:
		_, ok := v2.(MalFunc)
		if !ok {
			return false
		}
		panic("unimplemented")
	case MalBool:
		v2, ok := v2.(MalBool)
		if !ok {
			return false
		}
		return v1.Value == v2.Value
	case MalString:
		v2, ok := v2.(MalString)
		if !ok {
			return false
		}
		return v1.Value == v2.Value
	case MalList:
		v2, ok := v2.(MalList)
		if !ok {
			return false
		}
		if len(v1.Values) != len(v2.Values) {
			return false
		}
		for i := range v1.Values {
			if !malEq(v1.Values[i], v2.Values[i]) {
				return false
			}
		}

		return true
	default:
		panic("unreachable")
	}
}

func isMacroCall(ast MalValue, env *Env) bool {
	switch v := ast.(type) {
	case MalList:
		if len(v.Values) > 0 {
			if sym, ok := v.Values[0].(MalSymbol); ok {
				if f, ok := env.Get(sym.Value); ok {
					if f, ok := f.(MalInvoke); ok {
						return f.IsMacro()
					}
				}
			}
		}
	}

	return false
}
