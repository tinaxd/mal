package main

import (
	"bufio"
	"fmt"
	"os"
)

func read(param interface{}) interface{} {
	return param
}

func eval(param interface{}) interface{} {
	return param
}

func print(param interface{}) interface{} {
	return param
}

func rep(param interface{}) interface{} {
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
