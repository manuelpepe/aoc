package main

import (
	"3/lexer"
	"3/token"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var debug = false

func log(msg string, args ...any) {
	if debug {
		fmt.Printf(msg, args...)
	}
}

var MUL = []byte{'m', 'u', 'l'}

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	lexer := lexer.NewLexer(string(data))

	log("parsing tokens...\n")
	tokens := make([]token.Token, 0)
	tok := lexer.NextToken()
	for {
		if tok.Type == token.EOF {
			break
		}
		tokens = append(tokens, tok)
		tok = lexer.NextToken()

	}
	log("parsed %d tokens\n", len(tokens))

	sol1(tokens)
	sol2(tokens)

}

func sol1(tokens []token.Token) {
	var total int64

	for ix, tok := range tokens {
		if len(tok.Literal) < 3 {
			continue
		}

		if tok.Literal[len(tok.Literal)-3:] != "mul" {
			continue
		}

		if tokens[ix+1].Literal != "(" {
			continue
		}

		if tokens[ix+2].Type != token.INT {
			continue
		}

		if tokens[ix+3].Literal != "," {
			continue
		}

		if tokens[ix+4].Type != token.INT {
			continue
		}

		if tokens[ix+5].Literal != ")" {
			continue
		}

		a, err := strconv.ParseInt(tokens[ix+2].Literal, 10, 64)
		if err != nil {
			panic(err)
		}

		b, err := strconv.ParseInt(tokens[ix+4].Literal, 10, 64)
		if err != nil {
			panic(err)
		}

		log("found mul(%d,%d)\n", a, b)

		total += a * b
	}

	fmt.Printf("Total is: %d\n", total)

}

func sol2(tokens []token.Token) {

	var total int64
	doing := true

	for ix, tok := range tokens {
		if strings.HasSuffix(tok.Literal, "do") {
			doing = true
		}

		if strings.HasSuffix(tok.Literal, "don't") {
			doing = false
		}

		if doing && strings.HasSuffix(tok.Literal, "mul") {
			res, err := parseMul(ix, tokens)
			if err != nil {
				continue
			}
			total += res
		}

	}

	fmt.Printf("Total with conditionals is: %d\n", total)
}

func parseMul(ix int, tokens []token.Token) (int64, error) {
	if tokens[ix+1].Literal != "(" {
		return 0, fmt.Errorf("invalid mul")
	}

	if tokens[ix+2].Type != token.INT {
		return 0, fmt.Errorf("invalid mul")
	}

	if tokens[ix+3].Literal != "," {
		return 0, fmt.Errorf("invalid mul")
	}

	if tokens[ix+4].Type != token.INT {
		return 0, fmt.Errorf("invalid mul")
	}

	if tokens[ix+5].Literal != ")" {
		return 0, fmt.Errorf("invalid mul")
	}

	a, err := strconv.ParseInt(tokens[ix+2].Literal, 10, 64)
	if err != nil {
		panic(err)
	}

	b, err := strconv.ParseInt(tokens[ix+4].Literal, 10, 64)
	if err != nil {
		panic(err)
	}

	log("found mul(%d,%d)\n", a, b)
	return a * b, nil
}
