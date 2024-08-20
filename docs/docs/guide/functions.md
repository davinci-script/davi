# Functions 

## Array

### Sort

```php
sort(list, [key])
```

Sort a list of values.

#### Example

```php
sort([3, 1, 2])

// output: [1, 2, 3]
```

### Append

```php
append(list, value1, value2, ...)
```

Append values to a array.

#### Example

```php
append([1, 2], 3, 4)

// output: [1, 2, 3, 4]
```

### Range

```php
range(n)
```

Generate a list of integers from 0 to n-1.

#### Example

```php
range(3)

// output: [0, 1, 2]
```

### Slice

```php
slice(str or list, start, end)
```

Get a substring or sublist from a string or list.

#### Example

```php
slice("hello", 1, 3)

// output: "el"
```

## String

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

### Join

```php
join(list, separator)
```

Join a list of strings into a single string with a separator.

#### Example

```php
join(["a", "b", "c"], ", ")

// output: "a, b, c"
```

### Length

```php
len(value)
```

Get the length of a string, list, or map.

#### Example

```php
len("hello")

// output: 5
```

### Type

```php
type(value)
```

Get the type of a value as a string.

#### Example

```php
type(42)

// output: "int"
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

### Explode

```php
explode([separator], string)
```

Explode a string into a list of substrings. It's the same as split() with the arguments reversed.

#### Example

```php
explode(", ", "a, b, c")

// output: ["a", "b", "c"]
```

### Up First

```php
upFirst(string)
```

Convert the first character of a string to uppercase.

#### Example

```php
upFirst("hello")

// output: "Hello"
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

### Rune

```php
rune(str)
```

Convert a 1-character string to an ASCII code.

#### Example

```php
rune("A")

// output: 65
```

### Split

```php
split(string, [separator])
```

Split a string into a list of substrings.

#### Example

```php
split("a, b, c", ", ")

// output: ["a", "b", "c"]
```

### Find

```php
find(haystack, needle)
```

Find the first occurrence of a substring in a string or a value in a list.

#### Example

```php
find("hello", "e")

// output: 1
```

### Char

```php
char(string)
```

Convert an ASCII code to a character.

#### Example

```php
char(65)

// output: "A"
```

### Lower

```php
lower(string)
```

Convert a string to lowercase.

#### Example

```php
lower("HELLO")

// output: "hello"
```

### Upper

```php
upper(string)
```

Convert a string to uppercase.

#### Example

```php
upper("hello")

// output: "HELLO"
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

## System

### Echo

```php
echo(value1, value2, ...)
```

Print values to the standard output.

#### Example

```php
echo("hello", 42)

// output: hello 42
```

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

### Time

```php
time(none)
```

Get the current date and time as a string.

#### Example

```php
time()

// output: "2018-01-01 12
```

### Read

```php
read([filename])
```

Read the contents of a file or standard input.

#### Example

```php
read("file.txt")

// output: "contents of file.txt"
```

### Args

```php
args(none)
```

Get the command-line arguments passed to the script.

#### Example

```php
args()

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

## Conversion

### Str

```php
str(value)
```

Convert a value to a string.

#### Example

```php
str([1, 2, 3])

// output: "[1, 2, 3]"
```

### Int

```php
int(value)
```

Convert a value to an integer.

#### Example

```php
int("42")

// output: 42
```

