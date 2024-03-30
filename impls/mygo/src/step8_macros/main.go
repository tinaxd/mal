package main

import (
	"bufio"
	"fmt"
	"os"
)

func macroexpand(ast MalValue, env *Env) (MalValue, error) {
	for isMacroCall(ast, env) {
		lst := ast.(MalList)
		sym := lst.Values[0].(MalSymbol)
		macroV, ok := env.Get(sym.Value)
		if !ok {
			panic("unreachable")
		}
		macro := macroV.(MalInvoke)

		args := lst.Values[1:]
		expanded, err := macro.Invoke(args)
		if err != nil {
			return nil, fmt.Errorf("error while expanding macro: %w", err)
		}

		ast = expanded
	}

	return ast, nil
}

func read(param string) (MalValue, error) {
	return ReadStr(param)
}

func eval(param MalValue, replEnv *Env, env *Env) (MalValue, error) {
	for {
		switch p := param.(type) {
		case MalList:
			if len(p.Values) == 0 {
				return param, nil
			}

			expanded, err := macroexpand(p, env)
			if err != nil {
				return nil, err
			}
			switch expanded.(type) {
			case MalList:
				param = expanded
				p = param.(MalList)
			default:
				evaled, err := EvalAst(expanded, replEnv, env)
				if err != nil {
					return nil, err
				}
				return evaled, nil
			}

			rawHead := p.Values[0]
			rawArgs := p.Values[1:]
			switch h := rawHead.(type) {
			case MalSymbol:
				// special forms
				switch h.Value {
				case "macroexpand":
					if len(rawArgs) != 1 {
						return nil, fmt.Errorf("wrong number of arguments for macroexpand")
					}
					expanded, err := macroexpand(rawArgs[0], env)
					if err != nil {
						return nil, err
					}
					return expanded, nil
				case "eval":
					if len(rawArgs) != 1 {
						return nil, fmt.Errorf("wrong number of arguments eval")
					}
					arg0, err := eval(rawArgs[0], replEnv, env)
					if err != nil {
						return nil, err
					}
					return eval(arg0, replEnv, replEnv) // evaluate in replEnv
				case "defmacro!":
					fallthrough
				case "def!":
					if len(rawArgs) != 2 {
						return nil, fmt.Errorf("wrong number of arguments")
					}
					key, ok := rawArgs[0].(MalSymbol)
					if !ok {
						return nil, fmt.Errorf("arg0 of def! must be MalSymbol, got %v", rawArgs[0])
					}
					val, err := eval(rawArgs[1], replEnv, env)
					if err != nil {
						return nil, err
					}

					if h.Value == "defmacro!" {
						switch f := val.(type) {
						case MalFunc:
							f.Macro = true
							env.Set(key.Value, f)
							return f, nil
						case MalTcoFunc:
							f.Fn.Macro = true
							env.Set(key.Value, f)
							return f, nil
						default:
							return nil, fmt.Errorf("defmacro! must be a function")
						}
					}

					env.Set(key.Value, val)
					return val, nil
				case "let*":
					if len(rawArgs) != 2 {
						return nil, fmt.Errorf("wrong number of arguments")
					}

					bindings, ok := rawArgs[0].(MalList)
					if !ok {
						return nil, fmt.Errorf("arg0 of let* must be MalList, got %v", rawArgs[0])
					}
					if len(bindings.Values)%2 != 0 {
						return nil, fmt.Errorf("bindings must be even, got %v", bindings)
					}
					env = NewEnv(env, nil, nil)
					for i := 0; i < len(bindings.Values); i += 2 {
						key, ok := bindings.Values[i].(MalSymbol)
						if !ok {
							return nil, fmt.Errorf("binding key must be MalSymbol, got %v", bindings.Values[i])
						}
						val, err := eval(bindings.Values[i+1], replEnv, env)
						if err != nil {
							return nil, err
						}
						env.Set(key.Value, val)
					}

					param = rawArgs[1]
					continue
				case "do":
					if len(rawArgs) == 0 {
						return nil, fmt.Errorf("wrong number of arguments for do")
					}
					for _, arg := range rawArgs[:len(rawArgs)-1] {
						_, err := eval(arg, replEnv, env)
						if err != nil {
							return nil, err
						}
					}
					param = rawArgs[len(rawArgs)-1]
					continue
				case "if":
					if len(rawArgs) != 2 && len(rawArgs) != 3 {
						return nil, fmt.Errorf("wrong number of arguments for if")
					}
					cond, err := eval(rawArgs[0], replEnv, env)
					if err != nil {
						return nil, err
					}

					truthy := true
					if cond == nil {
						truthy = false
					} else if b, ok := cond.(MalBool); ok {
						truthy = b.Value
					}

					if truthy {
						param = rawArgs[1]
						continue
					} else {
						if len(rawArgs) != 3 {
							return nil, nil
						} else {
							param = rawArgs[2]
							continue
						}
					}
				case "fn*":
					if len(rawArgs) != 2 {
						return nil, fmt.Errorf("wrong number of arguments for fn*")
					}

					params, ok := rawArgs[0].(MalList)
					if !ok {
						return nil, fmt.Errorf("first argument of fn* must be MalList, got %v", rawArgs[0])
					}
					paramStrs := make([]string, len(params.Values))
					for i, p := range params.Values {
						sym, ok := p.(MalSymbol)
						if !ok {
							return nil, fmt.Errorf("parameter must be MalSymbol, got %v", p)
						}
						paramStrs[i] = sym.Value
					}

					fn := func(args []MalValue) (MalValue, error) {
						if len(args) != len(paramStrs) {
							return nil, fmt.Errorf("wrong number of arguments")
						}
						newEnv := NewEnv(env, paramStrs, args)
						return eval(rawArgs[1], replEnv, newEnv)
					}
					return MalTcoFunc{Ast: rawArgs[1], Params: paramStrs, Env: env, Fn: MalFunc{F: fn}}, nil
				case "quote":
					if len(rawArgs) != 1 {
						return nil, fmt.Errorf("wrong number of arguments for quote")
					}
					return rawArgs[0], nil
				case "quasiquoteexpand":
					if len(rawArgs) != 1 {
						return nil, fmt.Errorf("wrong number of arguments for quasiquoteexpand")
					}
					return quasiquote(rawArgs[0])
				case "quasiquote":
					if len(rawArgs) != 1 {
						return nil, fmt.Errorf("wrong number of arguments for quasiquote")
					}
					q, err := quasiquote(rawArgs[0])
					if err != nil {
						return nil, err
					}
					param = q
					continue
				}
			}

			evalListR, err := EvalAst(p, replEnv, env)
			if err != nil {
				return nil, err
			}
			evalList := evalListR.(MalList)
			if len(evalList.Values) == 0 {
				panic("unreachable")
			}
			head := evalList.Values[0]
			args := evalList.Values[1:]

			switch f := head.(type) {
			case MalFunc:
				return f.F(args)
			case MalTcoFunc:
				param = f.Ast
				env = NewEnv(f.Env, f.Params, args)
				continue
			}
		default:
			return EvalAst(param, replEnv, env)
		}
	}
}

func print(param MalValue) string {
	return PrStr(param, true)
}

func rep(param string, env *Env) (string, error) {
	step1, err := read(param)
	if err != nil {
		return "", err
	}
	step2, err := eval(step1, env, env)
	if err != nil {
		return "", err
	}
	step3 := print(step2)
	return step3, nil
}

func main() {
	env := InitialEnv()

	rep("(def! load-file (fn* (f) (eval (read-string (str \"(do \" (slurp f) \"\nnil)\")))))", env)
	rep("(def! not (fn* (a) (if a false true)))", env)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		prompt := "user> "
		fmt.Print(prompt)

		scanner.Scan()
		userInput := scanner.Text()

		if userInput == "" {
			fmt.Println("")
			break
		}

		result, err := rep(userInput, env)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		fmt.Println(result)
	}
}
