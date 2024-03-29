package main

import "regexp"

type Reader struct {
	Tokens   []string
	Position int
}

func NewReader(tokens []string) *Reader {
	return &Reader{Tokens: tokens, Position: 0}
}

func (r *Reader) Next() string {
	token := r.Tokens[r.Position]
	r.Position++
	return token
}

func (r *Reader) Peek() string {
	return r.Tokens[r.Position]
}

func ReadStr(input string) *Reader {
	tokens := Tokenize(input)
	return NewReader(tokens)
}

func Tokenize(input string) []string {
	re := "[\\s,]*(~@|[\\[\\]{}()'`~^@]|\"(?:\\.|[^\\\"])*\"?|;.*|[^\\s\\[\\]{}('\"`,;)]*)"
	compiled := regexp.MustCompile(re)

	rem := input
	tokens := []string{}
	for {
		token := compiled.FindStringSubmatch(rem)
		if len(token) == 0 {
			break
		}

		tokens = append(tokens, token[1])
		rem = rem[len(tokens[0]):]
	}

	return tokens
}

func (r *Reader) ReadForm() string {

}
