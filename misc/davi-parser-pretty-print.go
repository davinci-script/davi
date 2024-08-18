package main

import (
	"bytes"
	"fmt"
	"github.com/DavinciScript/Davi/lexer"
	"github.com/DavinciScript/Davi/parser"
	"github.com/hokaccha/go-prettyjson"
	"os"
	"strings"
)

func main() {

	fmt.Println("Running davi.go")

	input := []byte(`

	// This is a comment
	//$firstMessage = "Hello World";
	//$secondMessage = "Hello DavinciScript";
	//
	//function person($name, $age) {
	//
	//}
	//
	//echo "$firstMessage";
	//echo ($secondMessage, $firstMessage);
    
    echo("qko");
	
	`)

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
