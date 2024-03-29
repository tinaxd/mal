package main

import "strconv"

func PrStr(v MalValue) string {
	if v == nil {
		return "nil"
	}

	switch vv := v.(type) {
	case MalSymbol:
		return vv.Value
	case MalInt:
		return strconv.FormatInt(vv.Value, 10)
	case MalBool:
		if vv.Value {
			return "true"
		} else {
			return "false"
		}
	case MalFunc:
		return "#<function>"
	case MalList:
		str := "("
		for i, value := range vv.Values {
			if i != 0 {
				str += " "
			}
			str += PrStr(value)
		}
		str += ")"
		return str
	}

	panic("unreachable")
}
