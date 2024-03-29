package main

import (
	"bufio"
	"fmt"
	"os"
)

func read(param string) MalValue {
	return ReadStr(param)
}

func eval(param MalValue) MalValue {
	return param
}

func print(param MalValue) interface{} {
	return PrStr(param)
}

func rep(param string) interface{} {
	step1 := read(param)
	step2 := eval(step1)
	step3 := print(step2)
	return step3
}

func main() {
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

		result := rep(userInput)
		fmt.Println(result)
	}
}
