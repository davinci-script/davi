package main

import (
	"encoding/json"
	"fmt"
	. "github.com/DavinciScript/Davi/lexer"
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
		out, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(out))
		if tok == ILLEGAL {
			break
		}
	}

}
