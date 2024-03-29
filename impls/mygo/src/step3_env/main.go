package main

import (
	"bufio"
	"fmt"
	"os"
)

func read(param string) (MalValue, error) {
	return ReadStr(param)
}

func eval(param MalValue, env *Env) (MalValue, error) {
	switch p := param.(type) {
	case MalList:
		if len(p.Values) == 0 {
			return param, nil
		}

		rawHead := p.Values[0]
		rawArgs := p.Values[1:]
		switch h := rawHead.(type) {
		case MalSymbol:
			switch h.Value {
			case "def!":
				if len(rawArgs) != 2 {
					return nil, fmt.Errorf("wrong number of arguments")
				}
				key, ok := rawArgs[0].(MalSymbol)
				if !ok {
					return nil, fmt.Errorf("arg0 of def! must be MalSymbol, got %v", rawArgs[0])
				}
				val, err := eval(rawArgs[1], env)
				if err != nil {
					return nil, err
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
				newEnv := NewEnv(env)
				for i := 0; i < len(bindings.Values); i += 2 {
					key, ok := bindings.Values[i].(MalSymbol)
					if !ok {
						return nil, fmt.Errorf("binding key must be MalSymbol, got %v", bindings.Values[i])
					}
					val, err := eval(bindings.Values[i+1], newEnv)
					if err != nil {
						return nil, err
					}
					newEnv.Set(key.Value, val)
				}

				return eval(rawArgs[1], newEnv)
			}
		}

		evalListR, err := EvalAst(p, env)
		if err != nil {
			return nil, err
		}
		evalList := evalListR.(MalList)
		if len(evalList.Values) == 0 {
			panic("unreachable")
		}
		head := evalList.Values[0]
		args := evalList.Values[1:]

		f, ok := head.(MalFunc)
		if !ok {
			return nil, fmt.Errorf("first element of list must be MalFunc, got %v", head)
		}
		return f.F(args)
	default:
		return EvalAst(param, env)
	}
}

func print(param MalValue) string {
	return PrStr(param)
}

func rep(param string, env *Env) (string, error) {
	step1, err := read(param)
	if err != nil {
		return "", err
	}
	step2, err := eval(step1, env)
	if err != nil {
		return "", err
	}
	step3 := print(step2)
	return step3, nil
}

func main() {
	env := InitialEnv()

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
