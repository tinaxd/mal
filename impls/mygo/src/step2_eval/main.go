package main

import (
	"bufio"
	"fmt"
	"os"
)

func read(param string) (MalValue, error) {
	return ReadStr(param)
}

func eval(param MalValue, env Env) (MalValue, error) {
	switch p := param.(type) {
	case MalList:
		if len(p.Values) == 0 {
			return param, nil
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
			return nil, fmt.Errorf("first element is not a function: %v", head)
		}
		return f.F(args)
	default:
		return EvalAst(param, env)
	}
}

func print(param MalValue) string {
	return PrStr(param)
}

func rep(param string, env Env) (string, error) {
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
