package main

import (
	"fmt"
	. "github.com/DavinciScript/Davi/lexer"
	"github.com/hokaccha/go-prettyjson"
)

func daviLexerPrettyPrint() {

	fmt.Println("Running davi.go")

	// Run the lexer
	lexer := NewLexer([]byte(`

$age = 20;

`))
	for {
		pos, tok, val, ch := lexer.Next()
		if tok == EOF {
			break
		}
		data := map[string]interface{}{
			"posLine":   pos.Line,
			"posColumn": pos.Column,
			"token":     tok,
			"value":     val,
			"char":      ch,
		}

		formatter := prettyjson.NewFormatter()
		output, _ := formatter.Marshal(data)
		fmt.Println(string(output))

		if tok == ILLEGAL {
			break
		}
	}

}
