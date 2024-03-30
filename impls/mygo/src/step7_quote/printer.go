package main

import (
	"strconv"
	"strings"
)

func readableString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}

func PrStr(v MalValue, readably bool) string {
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
	case MalTcoFunc:
		return "#<function>"
	case MalString:
		if readably {
			return "\"" + readableString(vv.Value) + "\""
		} else {
			return vv.Value
		}
	case MalList:
		str := "("
		for i, value := range vv.Values {
			if i != 0 {
				str += " "
			}
			str += PrStr(value, readably)
		}
		str += ")"
		return str
	case *MalAtom:
		var v MalValue
		if vv.Ref != nil {
			v = vv.Ref
		}
		return "(atom " + PrStr(v, readably) + ")"
	}

	panic("unreachable")
}
