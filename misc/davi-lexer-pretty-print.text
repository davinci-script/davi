// DaVinci Script

package main

import (
	"bytes"
	"fmt"
	. "github.com/DavinciScript/Davi/lexer"
	"github.com/hokaccha/go-prettyjson"
	"io/ioutil"
	"os"
)

func main() {

	filename := os.Args[1]

	input, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file. Please check the file path and try again.\n")
		os.Exit(1)
	}

	// Replace <?davi with empty string
	input = bytes.Replace(input, []byte("<?davi"), []byte(""), 1)

	// Replace ?> with empty string
	input = bytes.Replace(input, []byte("?>"), []byte(""), 1)

	// Trim
	input = bytes.TrimSpace(input)

	// Run the lexer
	lexer := NewLexer(input)
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
