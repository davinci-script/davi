// DaVinci Script

package main

import (
	"bytes"
	"fmt"
	"github.com/DavinciScript/Davi/interpreter"
	"github.com/DavinciScript/Davi/lexer"
	"github.com/DavinciScript/Davi/parser"
	"io/ioutil"
	"os"
	"strings"
)

func main() {

	filename := os.Args[1]

	if len(os.Args) < 2 {
		fmt.Println("Usage: davi <filename>")
		os.Exit(1)
	}

	if filename == "--generate-docs" {
		interpreter.GenerateDocs()
		os.Exit(0)
	}

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

	//fmt.Print(string(input))
	//os.Exit(1)

	prog, err := parser.ParseProgram(input)
	if err != nil {
		errorMessage := fmt.Sprintf("%s", err)
		if e, ok := err.(parser.Error); ok {
			showErrorSource(input, e.Position, len(errorMessage))
		}
		fmt.Println(errorMessage)
		os.Exit(1)
	}

	_, err = interpreter.Execute(prog, &interpreter.Config{})
	if err != nil {
		errorMessage := fmt.Sprintf("%s", err)
		if e, ok := err.(interpreter.Error); ok {
			showErrorSource(input, e.Position(), len(errorMessage))
		}
		fmt.Println(errorMessage)
		os.Exit(1)
	}

}

// Show the source line and position of a parser or interpreter error
func showErrorSource(source []byte, pos lexer.Position, dividerLen int) {

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
