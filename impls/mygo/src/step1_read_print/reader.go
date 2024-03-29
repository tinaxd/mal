package main

import (
	"regexp"
	"strconv"
)

type Reader struct {
	Tokens   []string
	Position int
}

func NewReader(tokens []string) *Reader {
	return &Reader{Tokens: tokens, Position: 0}
}

func (r *Reader) Next() string {
	if r.Position >= len(r.Tokens) {
		panic("unexpected EOF")
	}
	token := r.Tokens[r.Position]
	r.Position++
	return token
}

func (r *Reader) Peek() string {
	if r.Position >= len(r.Tokens) {
		panic("unexpected EOF")
	}
	return r.Tokens[r.Position]
}

func ReadStr(input string) MalValue {
	tokens := Tokenize(input)
	r := NewReader(tokens)
	return r.ReadForm()
}

func Tokenize(input string) []string {
	re := "[\\s,]*(~@|[\\[\\]{}()'`~^@]|\"(?:\\.|[^\\\"])*\"?|;.*|[^\\s\\[\\]{}('\"`,;)]*)"
	compiled := regexp.MustCompile(re)

	rem := input
	tokens := []string{}
	for {
		// log.Printf("current rem: `%v`", rem)
		token := compiled.FindStringSubmatch(rem)
		if len(token) == 0 || token[1] == "" {
			break
		}
		// log.Printf("found token: %s", token[1])
		// log.Printf("token: %v", token)
		// log.Printf("len(token): %d", len(token))
		// log.Printf("len(token[0]): %d", len(token[0]))

		tokens = append(tokens, token[1])
		rem = rem[len(token[0]):]
	}

	return tokens
}

func (r *Reader) ReadForm() MalValue {
	if r.Peek() == "(" {
		return r.ReadList()
	} else {
		return r.ReadAtom()
	}
}

func (r *Reader) ReadList() MalValue {
	r.Next() // consume "("
	values := []MalValue{}
	for {
		if r.Peek() == ")" {
			r.Next() // consume ")"
			break
		}
		values = append(values, r.ReadForm())
	}

	return MalList{Values: values}
}

func (r *Reader) ReadAtom() MalValue {
	token := r.Next()
	// if token[0] is a digit
	if token[0] >= '0' && token[0] <= '9' {
		integer, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			panic(err)
		}
		return MalInt{Value: integer}
	} else {
		return MalSymbol{Value: token}
	}
}
