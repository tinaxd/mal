package main

import (
	"errors"
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

func (r *Reader) Next() (string, error) {
	if r.Position >= len(r.Tokens) {
		return "", errors.New("unexpected EOF")
	}
	token := r.Tokens[r.Position]
	r.Position++
	return token, nil
}

func (r *Reader) Peek() (string, error) {
	if r.Position >= len(r.Tokens) {
		return "", errors.New("unexpected EOF")
	}
	return r.Tokens[r.Position], nil
}

func ReadStr(input string) (MalValue, error) {
	tokens := Tokenize(input)
	r := NewReader(tokens)
	return r.ReadForm()
}

func Tokenize(input string) []string {
	re := `[\s,]*(~@|[\[\]{}()'\x60~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"\x60,;)]*)`
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

func (r *Reader) ReadForm() (MalValue, error) {
	peek, err := r.Peek()
	if err != nil {
		return nil, err
	}
	if peek == "(" {
		return r.ReadList(false)
	} else if peek == "[" {
		return r.ReadList(true)
	} else if peek == "'" {
		r.Next() // consume "'"
		form, err := r.ReadForm()
		if err != nil {
			return nil, err
		}
		return MalList{Values: []MalValue{makeSymbol("quote"), form}}, nil
	} else if peek == "`" {
		r.Next() // consume "`"
		form, err := r.ReadForm()
		if err != nil {
			return nil, err
		}
		return MalList{Values: []MalValue{makeSymbol("quasiquote"), form}}, nil
	} else if peek == "~" {
		r.Next() // consume "~"
		form, err := r.ReadForm()
		if err != nil {
			return nil, err
		}
		return MalList{Values: []MalValue{makeSymbol("unquote"), form}}, nil
	} else if peek == "~@" {
		r.Next() // consume "~@"
		form, err := r.ReadForm()
		if err != nil {
			return nil, err
		}
		return MalList{Values: []MalValue{makeSymbol("splice-unquote"), form}}, nil
	} else {
		return r.ReadAtom()
	}
}

func (r *Reader) ReadList(isVector bool) (MalValue, error) {
	r.Next() // consume "("
	values := []MalValue{}
	for {
		peek, err := r.Peek()
		if err != nil {
			return nil, err
		}
		if peek == ")" {
			if isVector {
				return nil, errors.New("unexpected `)`")
			}
			r.Next() // consume ")"
			break
		}
		if peek == "]" {
			if !isVector {
				return nil, errors.New("unexpected `]`")
			}
			r.Next() // consume "]"
			break
		}
		form, err := r.ReadForm()
		if err != nil {
			return nil, err
		}
		values = append(values, form)
	}

	if isVector {
		return NewVector(values), nil
	} else {
		return NewList(values), nil
	}
}

func (r *Reader) ReadAtom() (MalValue, error) {
	token, err := r.Next()
	if err != nil {
		return nil, err
	}
	// if token[0] is a digit
	if (token[0] >= '0' && token[0] <= '9') || (len(token) > 1 && token[0] == '-' && (token[1] >= '0' && token[1] <= '9')) {
		integer, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			panic(err)
		}
		return MalInt{Value: integer}, nil
	} else if token == "true" {
		return MalBool{Value: true}, nil
	} else if token == "false" {
		return MalBool{Value: false}, nil
	} else if token == "nil" {
		return nil, nil
	} else if token[0] == ':' {
		substr := token[1:]
		return NewKeyword(substr), nil
	} else if token[0] == '"' {
		substr := token[1 : len(token)-1]
		substr = readString(substr)
		return MalString{Value: substr}, nil
	} else {
		return MalSymbol{Value: token}, nil
	}
}

func readString(s string) string {
	result := ""
	backslash := false
	for _, ch := range s {
		if !backslash {
			if ch == '\\' {
				backslash = true
			} else {
				result += string(ch)
			}
			continue
		}

		switch ch {
		case '\\':
			result += "\\"
		case 'n':
			result += "\n"
		case '"':
			result += "\""
		default:
			panic("unsupported escape character")
		}
		backslash = false
	}

	return result
}
