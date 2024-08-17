package main

import (
	"bytes"
	"fmt"
	"github.com/DavinciScript/Davi/interpreter"
	"github.com/DavinciScript/Davi/lexer"
	"github.com/DavinciScript/Davi/parser"
	"os"
	"strings"
	"time"
)

func main() {

	input := []byte(`

		//$url = "https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&current=temperature_2m,wind_speed_10m&hourly=temperature_2m,relative_humidity_2m,wind_speed_10m";
		//fileGetContents($url);

		httpRegister("GET", "/hello");
		httpListen(":3030");
		
	`)

	prog, err := parser.ParseProgram(input)
	if err != nil {
		errorMessage := fmt.Sprintf("%s", err)
		if e, ok := err.(parser.Error); ok {
			showErrorSource(input, e.Position, len(errorMessage))
		}
		fmt.Println(errorMessage)
		os.Exit(1)
	}

	startTime := time.Now()
	stats, err := interpreter.Execute(prog, &interpreter.Config{})
	if err != nil {
		errorMessage := fmt.Sprintf("%s", err)
		if e, ok := err.(interpreter.Error); ok {
			showErrorSource(input, e.Position(), len(errorMessage))
		}
		fmt.Println(errorMessage)
		os.Exit(1)
	}
	showStats := false
	if showStats {
		elapsed := time.Since(startTime)
		fmt.Printf("%s elapsed: %d ops (%.0f/s), %d builtin calls (%.0f/s), %d user calls (%.0f/s)\n",
			elapsed,
			stats.Ops, float64(stats.Ops)/elapsed.Seconds(),
			stats.BuiltinCalls, float64(stats.BuiltinCalls)/elapsed.Seconds(),
			stats.UserCalls, float64(stats.UserCalls)/elapsed.Seconds(),
		)
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
