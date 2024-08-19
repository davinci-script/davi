// DaVinci Script

package interpreter

import (
	"fmt"
	"github.com/DavinciScript/Davi/interpreter/functions"
	. "github.com/DavinciScript/Davi/lexer"
	"github.com/DavinciScript/Davi/parser"
	"github.com/hokaccha/go-prettyjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type functionType interface {
	call(interp *interpreter, pos Position, args []Value) Value
	name() string
}

type userFunction struct {
	Name       string
	Parameters []string
	Ellipsis   bool
	Body       parser.Block
	Closure    map[string]Value
}

func ensureNumArgs(pos Position, name string, args []Value, required int) {
	if len(args) != required {
		plural := ""
		if required != 1 {
			plural = "s"
		}
		panic(typeError(pos, "ensure num args: %s() requires %d arg%s, got %d", name, required, plural, len(args)))
	}
}

func (f *userFunction) call(interp *interpreter, pos Position, args []Value) Value {
	if f.Ellipsis {
		ellipsisArgs := args[len(f.Parameters)-1:]
		newArgs := make([]Value, 0, len(f.Parameters)+1)
		newArgs = append(newArgs, args[:len(f.Parameters)-1]...)
		args = append(newArgs, Value(&ellipsisArgs))
	}
	ensureNumArgs(pos, f.Name, args, len(f.Parameters))
	interp.pushScope(f.Closure)
	defer interp.popScope()
	interp.pushScope(make(map[string]Value))
	defer interp.popScope()
	for i, arg := range args {
		interp.assign(f.Parameters[i], arg)
	}
	interp.stats.UserCalls++
	interp.executeBlock(f.Body)
	return Value(nil)
}

func (f *userFunction) name() string {
	if f.Name == "" {
		return "<function>"
	}
	return fmt.Sprintf("<function %s>", f.Name)
}

type builtinFunction struct {
	Function func(interp *interpreter, pos Position, args []Value) Value
	Name     string
}

func (f builtinFunction) call(interp *interpreter, pos Position, args []Value) Value {
	interp.stats.BuiltinCalls++
	return f.Function(interp, pos, args)
}

func (f builtinFunction) name() string {
	return fmt.Sprintf("<builtin %s>", f.Name)
}

var builtins = map[string]builtinFunction{
	"append":            {appendFunction, "append"},
	"args":              {argsFunction, "args"},
	"char":              {charFunction, "char"},
	"exit":              {exitFunction, "exit"},
	"find":              {findFunction, "find"},
	"int":               {intFunction, "int"},
	"join":              {joinFunction, "join"},
	"len":               {lenFunction, "len"},
	"lower":             {lowerFunction, "lower"},
	"echo":              {echoFunction, "echo"},
	"range":             {rangeFunction, "range"},
	"read":              {readFunction, "read"},
	"rune":              {runeFunction, "rune"},
	"slice":             {sliceFunction, "slice"},
	"sort":              {sortFunction, "sort"},
	"split":             {splitFunction, "split"},
	"str":               {strFunction, "str"},
	"type":              {typeFunction, "type"},
	"upper":             {upperFunction, "upper"},
	"time":              {timeFunction, "time"},
	"file_get_contents": {fileGetContentsFunction, "fileGetContents"},
	"http_register":     {httpRegisterFunction, "httpRegister"},
	"http_listen":       {httpListenFunction, "httpListen"},
}

/**
 * function: append
 * args: list, value1, value2, ...
 * return: nil
 * example: append([1, 2], 3, 4) => [1, 2, 3, 4]
 * output: [1, 2, 3, 4]
 * description: Append values to a array.
 * title: Append
 * category: Array
 */
func appendFunction(interp *interpreter, pos Position, args []Value) Value {
	if len(args) < 1 {
		panic(typeError(pos, "append() requires at least 1 arg, got %d", len(args)))
	}
	if list, ok := args[0].(*[]Value); ok {
		*list = append(*list, args[1:]...)
		return Value(nil)
	}
	panic(typeError(pos, "append() requires first argument to be list"))
}

func stringsToList(strings []string) Value {
	values := make([]Value, len(strings))
	for i, s := range strings {
		values[i] = s
	}
	return Value(&values)
}

/**
 * function: args
 * args: none
 * return: list
 * example: args() => ["arg1", "arg2"]
 * output: ["arg1", "arg2"]
 * description: Get the command-line arguments passed to the script.
 * title: Args
 * category: System
 */
func argsFunction(interp *interpreter, pos Position, args []Value) Value {
	ensureNumArgs(pos, "args", args, 0)
	return stringsToList(interp.args)
}

/**
 * function: char
 * args: string
 * return: str
 * example: char(65) => "A"
 * output: "A"
 * description: Convert an ASCII code to a character.
 * title: Char
 * category: String
 */
func charFunction(interp *interpreter, pos Position, args []Value) Value {
	ensureNumArgs(pos, "char", args, 1)
	if code, ok := args[0].(int); ok {
		return string(code)
	}
	panic(typeError(pos, "char() requires an int, not %s", typeName(args[0])))
}

/**
 * function: exit
 * args: [code]
 * return: nil
 * example: exit(1)
 * output: exit status 1
 * description: Exit the script with an optional exit code.
 * title: Exit
 * category: System
 */
func exitFunction(interp *interpreter, pos Position, args []Value) Value {
	if len(args) > 1 {
		panic(typeError(pos, "exit() requires 0 or 1 args, got %d", len(args)))
	}
	code := 0
	if len(args) > 0 {
		arg, ok := args[0].(int)
		if !ok {
			panic(typeError(pos, "exit() requires an int, not %s", typeName(args[0])))
		}
		code = arg
	}
	interp.exit(code)
	return Value(nil)
}

/**
 * function: find
 * args: haystack, needle
 * return: int
 * example: find("hello", "e") => 1
 * output: 1
 * description: Find the first occurrence of a substring in a string or a value in a list.
 * title: Find
 * category: String
 */
func findFunction(interp *interpreter, pos Position, args []Value) Value {
	ensureNumArgs(pos, "find", args, 2)
	switch haystack := args[0].(type) {
	case string:
		if needle, ok := args[1].(string); ok {
			return Value(strings.Index(haystack, needle))
		}
		panic(typeError(pos, "find() on str requires second argument to be a str"))
	case *[]Value:
		needle := args[1]
		for i, v := range *haystack {
			if evalEqual(pos, needle, v).(bool) {
				return Value(i)
			}
		}
		return Value(-1)
	default:
		panic(typeError(pos, "find() requires first argument to be a str or list"))
	}
}

/**
 * function: int
 * args: value
 * return: int
 * example: int("42") => 42
 * output: 42
 * description: Convert a value to an integer.
 * title: Int
 * category: Conversion
 */
func intFunction(interp *interpreter, pos Position, args []Value) Value {
	ensureNumArgs(pos, "int", args, 1)
	switch arg := args[0].(type) {
	case int:
		return args[0]
	case string:
		i, err := strconv.Atoi(arg)
		if err != nil {
			return Value(nil)
		}
		return Value(i)
	default:
		panic(typeError(pos, "int() requires an int or a str"))
	}
}

/**
 * function: join
 * args: list, separator
 * return: str
 * example: join(["a", "b", "c"], ", ") => "a, b, c"
 * output: "a, b, c"
 * description: Join a list of strings into a single string with a separator.
 * title: Join
 * category: String
 */
func joinFunction(interp *interpreter, pos Position, args []Value) Value {
	ensureNumArgs(pos, "join", args, 2)
	sep, ok := args[1].(string)
	if !ok {
		panic(typeError(pos, "join() requires separator to be a str"))
	}
	if list, ok := args[0].(*[]Value); ok {
		strs := make([]string, len(*list))
		for i, v := range *list {
			s, ok := v.(string)
			if !ok {
				panic(typeError(pos, "join() requires all list elements to be strs"))
			}
			strs[i] = s
		}
		joined := strings.Join(strs, sep)
		return Value(joined)
	}
	panic(typeError(pos, "join() requires first argument to be a list"))
}

/**
 * function: len
 * args: value
 * return: int
 * example: len("hello") => 5
 * output: 5
 * description: Get the length of a string, list, or map.
 * title: Length
 * category: String
 */
func lenFunction(interp *interpreter, pos Position, args []Value) Value {
	ensureNumArgs(pos, "len", args, 1)
	var length int
	switch arg := args[0].(type) {
	case string:
		length = len(arg)
	case *[]Value:
		length = len(*arg)
	case map[string]Value:
		length = len(arg)
	default:
		panic(typeError(pos, "len() requires a str, list, or map"))
	}
	return Value(length)
}

/**
 * function: lower
 * args: string
 * return: str
 * example: lower("HELLO") => "hello"
 * output: "hello"
 * description: Convert a string to lowercase.
 * title: Lower
 * category: String
 */
func lowerFunction(interp *interpreter, pos Position, args []Value) Value {
	ensureNumArgs(pos, "lower", args, 1)
	if s, ok := args[0].(string); ok {
		return Value(strings.ToLower(s))
	}
	panic(typeError(pos, "lower() requires a str"))
}

/**
 * function: echo
 * args: value1, value2, ...
 * return: nil
 * example: echo("hello", 42) => hello 42
 * output: hello 42
 * description: Print values to the standard output.
 * title: Echo
 * category: System
 */
func echoFunction(interp *interpreter, pos Position, args []Value) Value {
	strs := make([]interface{}, len(args))
	for i, a := range args {
		strs[i] = toString(a, false)
	}
	fmt.Fprintln(interp.stdout, strs...)
	return Value(nil)
}

/**
 * function: range
 * args: n
 * return: list
 * example: range(3) => [0, 1, 2]
 * output: [0, 1, 2]
 * description: Generate a list of integers from 0 to n-1.
 * title: Range
 * category: Array
 */
func rangeFunction(interp *interpreter, pos Position, args []Value) Value {
	ensureNumArgs(pos, "range", args, 1)
	if n, ok := args[0].(int); ok {
		if n < 0 {
			panic(valueError(pos, "range() argument must not be negative"))
		}
		nums := make([]Value, n)
		for i := 0; i < n; i++ {
			nums[i] = i
		}
		return Value(&nums)
	}
	panic(typeError(pos, "range() requires an int"))
}

/**
 * function: read
 * args: [filename]
 * return: str
 * example: read("file.txt") => "contents of file.txt"
 * output: "contents of file.txt"
 * description: Read the contents of a file or standard input.
 * title: Read
 * category: System
 */
func readFunction(interp *interpreter, pos Position, args []Value) Value {
	if len(args) > 1 {
		panic(typeError(pos, "read() requires 0 or 1 args, got %d", len(args)))
	}
	var b []byte
	var err error
	if len(args) == 0 {
		b, err = ioutil.ReadAll(interp.stdin)
	} else {
		filename, ok := args[0].(string)
		if !ok {
			panic(typeError(pos, "read() argument must be a str"))
		}
		b, err = ioutil.ReadFile(filename)
	}
	if err != nil {
		panic(runtimeError(pos, "read() error: %v", err))
	}
	return Value(string(b))
}

/**
 * function: rune
 * args: str
 * return: int
 * example: rune("A") => 65
 * output: 65
 * description: Convert a 1-character string to an ASCII code.
 * title: Rune
 * category: String
 */
func runeFunction(interp *interpreter, pos Position, args []Value) Value {
	ensureNumArgs(pos, "rune", args, 1)
	if s, ok := args[0].(string); ok {
		runes := []rune(s)
		if len(runes) != 1 {
			panic(valueError(pos, "rune() requires a 1-character str"))
		}
		return Value(int(runes[0]))
	}
	panic(typeError(pos, "rune() requires a str"))
}

/**
 * function: slice
 * args: str or list, start, end
 * return: str or list
 * example: slice("hello", 1, 3) => "el"
 * output: "el"
 * description: Get a substring or sublist from a string or list.
 * title: Slice
 * category: Array
 */
func sliceFunction(interp *interpreter, pos Position, args []Value) Value {
	ensureNumArgs(pos, "slice", args, 3)
	start, sok := args[1].(int)
	end, eok := args[2].(int)
	if !sok || !eok {
		panic(typeError(pos, "slice() requires start and end to be ints"))
	}
	switch s := args[0].(type) {
	case string:
		if start < 0 || end > len(s) || start > end {
			panic(valueError(pos, "slice() start or end out of bounds"))
		}
		return Value(s[start:end])
	case *[]Value:
		if start < 0 || end > len(*s) || start > end {
			panic(valueError(pos, "slice() start or end out of bounds"))
		}
		result := make([]Value, end-start)
		copy(result, (*s)[start:end])
		return Value(&result)
	default:
		panic(typeError(pos, "slice() requires first argument to be a str or list"))
	}
}

/**
 * function: sort
 * args: list, [key]
 * return: nil
 * example: sort([3, 1, 2]) => [1, 2, 3]
 * output: [1, 2, 3]
 * description: Sort a list of values.
 * title: Sort
 * category: Array
 */
func sortFunction(interp *interpreter, pos Position, args []Value) Value {
	if len(args) != 1 && len(args) != 2 {
		panic(typeError(pos, "sort() requires 1 or 2 args, got %d", len(args)))
	}
	list, ok := args[0].(*[]Value)
	if !ok {
		panic(typeError(pos, "sort() requires first argument to be a list"))
	}
	if len(*list) <= 1 {
		return Value(nil)
	}
	if len(args) == 1 {
		sort.SliceStable(*list, func(i, j int) bool {
			return evalLess(pos, (*list)[i], (*list)[j]).(bool)
		})
	} else {
		keyFunc, ok := args[1].(functionType)
		if !ok {
			panic(typeError(pos, "sort() requires second argument to be a function"))
		}
		// Decorate, sort, undecorate (so we only call key function
		// once per element)
		type pair struct {
			value Value
			key   Value
		}
		pairs := make([]pair, len(*list))
		for i, v := range *list {
			key := interp.callFunction(pos, keyFunc, []Value{v})
			pairs[i] = pair{v, key}
		}
		sort.SliceStable(pairs, func(i, j int) bool {
			return evalLess(pos, pairs[i].key, pairs[j].key).(bool)
		})
		values := make([]Value, len(pairs))
		for i, p := range pairs {
			values[i] = p.value
		}
		*list = values
	}
	return Value(nil)
}

/**
 * function: split
 * args: string, [separator]
 * return: list
 * example: split("a, b, c", ", ") => ["a", "b", "c"]
 * output: ["a", "b", "c"]
 * description: Split a string into a list of substrings.
 * title: Split
 * category: String
 */
func splitFunction(interp *interpreter, pos Position, args []Value) Value {
	if len(args) != 1 && len(args) != 2 {
		panic(typeError(pos, "split() requires 1 or 2 args, got %d", len(args)))
	}
	str, ok := args[0].(string)
	if !ok {
		panic(typeError(pos, "split() requires first argument to be a str"))
	}
	var parts []string
	if len(args) == 1 || args[1] == nil {
		parts = strings.Fields(str)
	} else if sep, ok := args[1].(string); ok {
		parts = strings.Split(str, sep)
	} else {
		panic(typeError(pos, "split() requires separator to be a str or nil"))
	}
	return stringsToList(parts)
}

/**
 * function: str
 * args: value
 * return: str
 * example: str([1, 2, 3]) => "[1, 2, 3]"
 * output: "[1, 2, 3]"
 * description: Convert a value to a string.
 * title: Str
 * category: Conversion
 */
func toString(value Value, quoteStr bool) string {
	var s string
	switch v := value.(type) {
	case nil:
		s = "nil"
	case bool:
		if v {
			s = "true"
		} else {
			s = "false"
		}
	case int:
		s = fmt.Sprintf("%d", v)
	case string:
		if quoteStr {
			s = fmt.Sprintf("%q", v)
		} else {
			s = v
		}
	case *[]Value:
		strs := make([]string, len(*v))
		for i, v := range *v {
			strs[i] = toString(v, true)
		}
		s = fmt.Sprintf("[%s]", strings.Join(strs, ", "))
	case map[string]Value:
		strs := make([]string, 0, len(v))
		for k, v := range v {
			item := fmt.Sprintf("%q: %s", k, toString(v, true))
			strs = append(strs, item)
		}
		sort.Strings(strs) // Ensure str(output) is consistent
		s = fmt.Sprintf("{%s}", strings.Join(strs, ", "))
	case functionType:
		s = v.name()
	default:
		// Interpreter should never give us this
		panic(fmt.Sprintf("str() got unexpected type %T", v))
	}
	return s
}

/**
 * function: str
 * args: value
 * return: str
 * example: str([1, 2, 3]) => "[1, 2, 3]"
 * output: "[1, 2, 3]"
 * description: Convert a value to a string.
 * title: Str
 * category: Conversion
 */
func strFunction(interp *interpreter, pos Position, args []Value) Value {
	ensureNumArgs(pos, "str", args, 1)
	return Value(toString(args[0], false))
}

func typeName(v Value) string {
	var t string

	switch v.(type) {
	case nil:
		t = "nil"
	case bool:
		t = "bool"
	case int:
		t = "int"
	case string:
		t = "str"
	case *[]Value:
		t = "list"
	case map[string]Value:
		t = "map"
	case functionType:
		t = "function"
	case *userObject:
		t = "object"

	default:
		// Interpreter should never give us this
		panic(fmt.Sprintf("type() got unexpected type %T", v))
	}
	return t
}

/**
 * function: type
 * args: value
 * return: str
 * example: type(42) => "int"
 * output: "int"
 * description: Get the type of a value as a string.
 * title: Type
 * category: String
 */
func typeFunction(interp *interpreter, pos Position, args []Value) Value {
	ensureNumArgs(pos, "type", args, 1)
	return Value(typeName(args[0]))
}

/**
 * function: upper
 * args: string
 * return: str
 * example: upper("hello") => "HELLO"
 * output: "HELLO"
 * description: Convert a string to uppercase.
 * title: Upper
 * category: String
 */
func upperFunction(interp *interpreter, pos Position, args []Value) Value {
	ensureNumArgs(pos, "upper", args, 1)
	if s, ok := args[0].(string); ok {
		return Value(strings.ToUpper(s))
	}
	panic(typeError(pos, "upper() requires a str"))
}

/**
 * function: time
 * args: none
 * return: str
 * example: time() => "2018-01-01 12:00:00"
 * output: "2018-01-01 12:00:00"
 * description: Get the current date and time as a string.
 * title: Time
 * category: System
 */
func timeFunction(interp *interpreter, pos Position, args []Value) Value {
	ensureNumArgs(pos, "time", args, 0)

	dt := time.Now()

	return Value(dt.String())
}

/**
 * function: fileGetContents
 * args: url
 * return: str
 * example: fileGetContents("http://example.com") => "..."
 * output: "..."
 * description: Get the contents of a file or URL.
 * title: File Get Contents
 * category: File System
 */
func fileGetContentsFunction(interp *interpreter, pos Position, args []Value) Value {

	ensureNumArgs(pos, "fileGetContents", args, 1)
	if s, ok := args[0].(string); ok {

		_, err := url.ParseRequestURI(s)
		if err == nil {
			data, err := functions.GetContentFromUrl(s)
			if err != nil {
				panic(runtimeError(pos, "fileGetContents() error: %v", err))
			} else {
				return Value(string(data))
			}
		} else {
			fmt.Fprintln(interp.stdout, s)
		}

		return Value(nil)
	}

	panic(typeError(pos, "fileGetContents() requires a str"))
}

/**
 * function: httpRegister
 * args: pattern, handler
 * return: nil
 * example: httpRegister("/", func() { return "Hello, World!" })
 * output: "Hello, World!"
 * description: Register a handler function for a URL pattern.
 * title: HTTP Register
 * category: HTTP
 */
func httpRegisterFunction(interp *interpreter, pos Position, args []Value) Value {

	ensureNumArgs(pos, "httpRegister", args, 2)
	//
	if len(args) != 1 && len(args) != 2 {
		panic(typeError(pos, "httpRegisterFunction() requires 2 args, got %d", len(args)))
	}

	handlerFunction, ok := args[1].(functionType)
	if !ok {
		panic(typeError(pos, "httpRegisterFunction() requires second argument to be a function"))
	}

	formatter := prettyjson.NewFormatter()
	output, _ := formatter.Marshal(handlerFunction)
	fmt.Println(string(output))

	pattern := args[0].(string)

	getRoot := func(w http.ResponseWriter, r *http.Request) {
		outputFunction := interp.callFunction(pos, handlerFunction, []Value{})
		fmt.Fprintln(w, outputFunction)
	}
	http.HandleFunc(pattern, getRoot)

	return Value(nil)

}

/**
 * function: httpListen
 * args: portOrAddress
 * return: nil
 * example: httpListen(":8080")
 * output: Server is starting on http://localhost:8080...
 * description: Start the HTTP server.
 * title: HTTP Listen
 * category: HTTP
 */
func httpListenFunction(interp *interpreter, pos Position, args []Value) Value {

	ensureNumArgs(pos, "httpListen", args, 1)

	if len(args) != 1 {
		panic(typeError(pos, "httpRegisterFunction() requires 1 arg, got %d", len(args)))
	}

	portOrAddress, ok := args[0].(string)
	if !ok {
		panic(typeError(pos, "httpListenFunction() requires first argument to be a string"))
	}

	if !strings.Contains(portOrAddress, ":") {
		fmt.Printf("Server is starting on %s...\n", portOrAddress)
	} else {
		fmt.Printf("Server is starting on http://localhost%s...\n", portOrAddress)
	}

	err := http.ListenAndServe(portOrAddress, nil)
	if err != nil {
		panic(runtimeError(pos, "httpListen() error: %v", err))
	}

	return Value(nil)
}
