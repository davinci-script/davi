// DaVinci Script

package main

import (
	"bytes"
	"fmt"
	"github.com/DavinciScript/Davi/lexer"
	"github.com/DavinciScript/Davi/parser"
	"github.com/hokaccha/go-prettyjson"
	"io/ioutil"
	"os"
	"strings"
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

	prog, err := parser.ParseProgram(input)
	if err != nil {
		errorMessage := fmt.Sprintf("%s", err)
		if e, ok := err.(parser.Error); ok {
			showErrorSourcePrettyPrint(input, e.Position, len(errorMessage))
		}
		fmt.Println(errorMessage)
		os.Exit(1)
	}

	formatter := prettyjson.NewFormatter()
	output, _ := formatter.Marshal(prog)
	fmt.Println(string(output))
}

// Show the source line and position of a parser or interpreter error
func showErrorSourcePrettyPrint(source []byte, pos lexer.Position, dividerLen int) {
	divider := strings.Repeat("-", dividerLen)
	if divider != "" {
		fmt.Println(divider)
	}
	lines := bytes.Split(source, []byte{'\n'})
	errorLine := string(lines[pos.Line-1])
	numTabs := strings.Count(errorLine[:pos.Column-1], "\t")
	fmt.Println(strings.Replace(errorLine, "\t", "    ", -1))
	fmt.Println(strings.Repeat(" ", pos.Column-1) + strings.Repeat("   ", numTabs) + "^")
	if divider != "" {
		fmt.Println(divider)
	}
}
