# Functions 

## String

### Split

```php
split(string, [separator])
```

Split a string into a list of substrings.

#### Example

```php
split("a, b, c", ", ") => ["a", "b", "c"]

// output: ["a", "b", "c"]
```

### Type

```php
type(value)
```

Get the type of a value as a string.

#### Example

```php
type(42) => "int"

// output: "int"
```

### Lower First

```php
lowerFirst(string)
```

Convert the first character of a string to lowercase.

#### Example

```php
lowerFirst("Hello")

// output: "hello"
```

### Char

```php
char(string)
```

Convert an ASCII code to a character.

#### Example

```php
char(65) => "A"

// output: "A"
```

### Rune

```php
rune(str)
```

Convert a 1-character string to an ASCII code.

#### Example

```php
rune("A") => 65

// output: 65
```

### Explode

```php
explode([separator], string)
```

Explode a string into a list of substrings. It's the same as split() with the arguments reversed.

#### Example

```php
explode(", ", "a, b, c") => ["a", "b", "c"]

// output: ["a", "b", "c"]
```

### Lower Words

```php
lowerWords(string)
```

Convert all words in a string to lowercase.

#### Example

```php
lowerWords("Hello, World!")

// output: "hello, world!"
```

### Dot Case

```php
dotCase(string)
```

Convert a string to dot.case.

#### Example

```php
dotCase("Hello, World!")

// output: "hello.world"
```

### Lower

```php
lower(string)
```

Convert a string to lowercase.

#### Example

```php
lower("HELLO") => "hello"

// output: "hello"
```

### Up Words

```php
upWords(string)
```

Convert all words in a string to uppercase.

#### Example

```php
upWords("hello, world!")

// output: "Hello, World!"
```

### Camel Case

```php
camelCase(string)
```

Convert a string to camelCase.

#### Example

```php
camelCase("Hello, World!")

// output: "helloWorld"
```

### Find

```php
find(haystack, needle)
```

Find the first occurrence of a substring in a string or a value in a list.

#### Example

```php
find("hello", "e") => 1

// output: 1
```

### Kebab Case

```php
kebabCase(string)
```

Convert a string to kebab-case.

#### Example

```php
kebabCase("Hello, World!")

// output: "hello-world"
```

### Join

```php
join(list, separator)
```

Join a list of strings into a single string with a separator.

#### Example

```php
join(["a", "b", "c"], ", ") => "a, b, c"

// output: "a, b, c"
```

### Snake Case

```php
snakeCase(string)
```

Convert a string to snake_case.

#### Example

```php
snakeCase("Hello, World!")

// output: "hello_world"
```

### Pascal Case

```php
pascalCase(string)
```

Convert a string to PascalCase.

#### Example

```php
pascalCase("Hello, World!")

// output: "HelloWorld"
```

### Length

```php
len(value)
```

Get the length of a string, list, or map.

#### Example

```php
len("hello") => 5

// output: 5
```

### Upper

```php
upper(string)
```

Convert a string to uppercase.

#### Example

```php
upper("hello") => "HELLO"

// output: "HELLO"
```

### Up First

```php
upFirst(string)
```

Convert the first character of a string to uppercase.

#### Example

```php
upFirst("hello") => "Hello"

// output: "Hello"
```

## Conversion

### Int

```php
int(value)
```

Convert a value to an integer.

#### Example

```php
int("42") => 42

// output: 42
```

### Str

```php
str(value)
```

Convert a value to a string.

#### Example

```php
str([1, 2, 3]) => "[1, 2, 3]"

// output: "[1, 2, 3]"
```

## Array

### Append

```php
append(list, value1, value2, ...)
```

Append values to a array.

#### Example

```php
append([1, 2], 3, 4) => [1, 2, 3, 4]

// output: [1, 2, 3, 4]
```

### Sort

```php
sort(list, [key])
```

Sort a list of values.

#### Example

```php
sort([3, 1, 2]) => [1, 2, 3]

// output: [1, 2, 3]
```

### Range

```php
range(n)
```

Generate a list of integers from 0 to n-1.

#### Example

```php
range(3) => [0, 1, 2]

// output: [0, 1, 2]
```

### Slice

```php
slice(str or list, start, end)
```

Get a substring or sublist from a string or list.

#### Example

```php
slice("hello", 1, 3) => "el"

// output: "el"
```

## System

### Exit

```php
exit([code])
```

Exit the script with an optional exit code.

#### Example

```php
exit(1)

// output: exit status 1
```

### Echo

```php
echo(value1, value2, ...)
```

Print values to the standard output.

#### Example

```php
echo("hello", 42) => hello 42

// output: hello 42
```

### Read

```php
read([filename])
```

Read the contents of a file or standard input.

#### Example

```php
read("file.txt") => "contents of file.txt"

// output: "contents of file.txt"
```

### Time

```php
time(none)
```

Get the current date and time as a string.

#### Example

```php
time() => "2018-01-01 12

// output: "2018-01-01 12
```

### Args

```php
args(none)
```

Get the command-line arguments passed to the script.

#### Example

```php
args() => ["arg1", "arg2"]

// output: ["arg1", "arg2"]
```

## File System

### File Get Contents

```php
fileGetContents(url)
```

Get the contents of a file or URL.

#### Example

```php
fileGetContents("http

// output: "..."
```

## HTTP

### HTTP Register

```php
httpRegister(pattern, handler)
```

Register a handler function for a URL pattern.

#### Example

```php
httpRegister("/", func() { return "Hello, World!" })

// output: "Hello, World!"
```

### HTTP Listen

```php
httpListen(portOrAddress)
```

Start the HTTP server.

#### Example

```php
httpListen("

// output: Server is starting on http
```

