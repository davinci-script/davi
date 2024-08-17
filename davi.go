package main

import (
	"fmt"
	. "github.com/DavinciScript/Davi/lexer"
	"github.com/davecgh/go-spew/spew"
)

func main() {

	fmt.Println("Running davi.go")

	// Run the lexer
	lexer := NewLexer([]byte(`print(1234, "foo")`))
	for {
		pos, tok, val := lexer.Next()
		if tok == EOF {
			break
		}
		data := map[string]interface{}{
			"posLine":   pos.Line,
			"posColumn": pos.Column,
			"token":     tok,
			"value":     val,
		}
		spew.Dump(data)
		if tok == ILLEGAL {
			break
		}
	}

}
