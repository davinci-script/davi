package main

import (
	"fmt"
)

func main() {

	fmt.Println("Running davi.go")

	// Run the lexer
	lexer := NewLexer([]byte(`print(1234, "foo") @`))
	for {
		pos, tok, val := lexer.Next()
		if tok == EOF {
			break
		}
		fmt.Printf("%d:%d %s %q\n", pos.Line, pos.Column, tok, val)
		if tok == ILLEGAL {
			break
		}
	}

}
