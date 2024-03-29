package main

import (
	"bufio"
	"fmt"
	"os"
)

func read(param string) (MalValue, error) {
	return ReadStr(param)
}

func eval(param MalValue) MalValue {
	return param
}

func print(param MalValue) string {
	return PrStr(param)
}

func rep(param string) (string, error) {
	step1, err := read(param)
	if err != nil {
		return "", err
	}
	step2 := eval(step1)
	step3 := print(step2)
	return step3, nil
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

		result, err := rep(userInput)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		fmt.Println(result)
	}
}
