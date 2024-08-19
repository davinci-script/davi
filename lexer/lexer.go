// DaVinci Script

package lexer

import (
	"fmt"
	"unicode/utf8"
)

// Token is the type of a single token
type Token int

const (
	// Stop tokens
	ILLEGAL Token = iota
	EOF

	// Single-character tokens
	ASSIGN
	COLON
	SEMI
	COMMA
	DIVIDE
	DOT
	GT
	LBRACE
	LBRACKET
	LPAREN
	LT
	MINUS
	MODULO
	PLUS
	RBRACE
	RBRACKET
	RPAREN
	TIMES
	QUESTION
	DOLLAR

	// Two-character tokens
	EQUAL
	GTE
	LTE
	NOTEQUAL

	// Three-character tokens
	ELLIPSIS

	// Keywords
	AND
	ELSE
	FALSE
	FOR
	FUNCTION
	IF
	IN
	NIL
	NOT
	OR
	RETURN
	TRUE
	WHILE
	CLASS
	EXTENDS
	PUBLIC
	PRIVATE
	PROTECTED
	STATIC
	ABSTRACT
	FINAL
	CONST
	NEW

	// Literals and identifiers
	INT
	NAME
	STR
)

var keywordTokens = map[string]Token{
	"and":       AND,
	"else":      ELSE,
	"false":     FALSE,
	"for":       FOR,
	"function":  FUNCTION,
	"if":        IF,
	"in":        IN,
	"nil":       NIL,
	"not":       NOT,
	"or":        OR,
	"return":    RETURN,
	"true":      TRUE,
	"while":     WHILE,
	"class":     CLASS,
	"extends":   EXTENDS,
	"public":    PUBLIC,
	"private":   PRIVATE,
	"protected": PROTECTED,
	"static":    STATIC,
	"abstract":  ABSTRACT,
	"final":     FINAL,
	"const":     CONST,
	"new":       NEW,
}

var tokenNames = map[Token]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	ASSIGN:   "=",
	COLON:    ":",
	SEMI:     ";",
	COMMA:    ",",
	DIVIDE:   "/",
	DOT:      ".",
	GT:       ">",
	LBRACE:   "{",
	LBRACKET: "[",
	LPAREN:   "(",
	LT:       "<",
	MINUS:    "-",
	MODULO:   "%",
	PLUS:     "+",
	RBRACE:   "}",
	RBRACKET: "]",
	RPAREN:   ")",
	TIMES:    "*",
	QUESTION: "?",
	DOLLAR:   "$",

	EQUAL:    "==",
	GTE:      ">=",
	LTE:      "<=",
	NOTEQUAL: "!=",

	ELLIPSIS: "...",

	AND:      "and",
	ELSE:     "else",
	FALSE:    "false",
	FOR:      "for",
	FUNCTION: "function",
	IF:       "if",
	IN:       "in",
	NIL:      "nil",
	NOT:      "not",
	OR:       "or",
	RETURN:   "return",
	TRUE:     "true",
	WHILE:    "while",

	INT:  "int",
	NAME: "name",
	STR:  "str",
}

func (t Token) String() string {
	return tokenNames[t]
}

// Position stores the line and column a token starts at
type Position struct {
	Line   int
	Column int
}

// Lexer parses input source code to a stream of tokens. Use
// NewLexer() to actually create a tokenizer, and Next() to get the next
// token in the input.
type Lexer struct {
	input    []byte
	offset   int
	ch       rune
	errorMsg string
	pos      Position
	nextPos  Position
}

// NewLexer returns a new tokenizer that works off the given input.
func NewLexer(input []byte) *Lexer {
	l := new(Lexer)
	l.input = input
	l.nextPos.Line = 1
	l.nextPos.Column = 1
	l.next()
	return l
}

func (l *Lexer) next() {
	l.pos = l.nextPos
	ch, size := utf8.DecodeRune(l.input[l.offset:])
	if size == 0 {
		l.ch = -1
		return
	}
	if ch == utf8.RuneError {
		l.ch = -1
		l.errorMsg = fmt.Sprintf("invalid UTF-8 byte 0x%02x", l.input[l.offset])
		return
	}
	if ch == '\n' {
		l.nextPos.Line++
		l.nextPos.Column = 1
	} else {
		l.nextPos.Column++
	}
	l.ch = ch
	l.offset += size
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
			l.next()
		}
		if !(l.ch == '/' && l.offset < len(l.input) && l.input[l.offset] == '/') {
			break
		}
		// Skip //-prefixed comment (to end of line or end of input)
		l.next()
		l.next()
		for l.ch != '\n' && l.ch >= 0 {
			l.next()
		}
		l.next()
	}
}

func isNameStart(ch rune) bool {
	return ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// Next() returns the position, token type, and token value of the next token
// in the source. For ordinary tokens, the token value is empty. For INT,
// NAME, and STR tokens, it's the number or string value. For an ILLEGAL
// token, it's the error message.
func (l *Lexer) Next() (Position, Token, string, string) {
	l.skipWhitespaceAndComments()
	if l.ch < 0 {
		if l.errorMsg != "" {
			return l.pos, ILLEGAL, l.errorMsg, ""
		}
		return l.pos, EOF, "", ""
	}

	pos := l.pos
	token := ILLEGAL
	value := ""

	ch := l.ch
	l.next()

	// Names (identifiers) and keywords
	if isNameStart(ch) {
		runes := []rune{ch}
		for isNameStart(l.ch) || (l.ch >= '0' && l.ch <= '9') {
			runes = append(runes, l.ch)
			l.next()
		}
		name := string(runes)
		token, isKeyword := keywordTokens[name]
		if !isKeyword {
			token = NAME
			value = name
		}
		return pos, token, value, string(ch)
	}

	switch ch {
	case ':':
		token = COLON
	case ';':
		token = SEMI
	case ',':
		token = COMMA
	case '/':
		token = DIVIDE
	case '{':
		token = LBRACE
	case '[':
		token = LBRACKET
	case '(':
		token = LPAREN
	case '-':
		token = MINUS
	case '%':
		token = MODULO
	case '+':
		token = PLUS
	case '}':
		token = RBRACE
	case ']':
		token = RBRACKET
	case ')':
		token = RPAREN
	case '*':
		token = TIMES
	case '?':
		token = QUESTION
	case '$':
		token = DOLLAR

	case '=':
		if l.ch == '=' {
			l.next()
			token = EQUAL
		} else {
			token = ASSIGN
		}
	case '!':
		if l.ch == '=' {
			l.next()
			token = NOTEQUAL
		} else {
			token = ILLEGAL
			value = fmt.Sprintf("expected != instead of !%c", l.ch)
		}
	case '<':
		if l.ch == '=' {
			l.next()
			token = LTE
		} else {
			token = LT
		}
	case '>':
		if l.ch == '=' {
			l.next()
			token = GTE
		} else {
			token = GT
		}

	case '.':
		if l.ch == '.' {
			l.next()
			if l.ch != '.' {
				return pos, ILLEGAL, "unexpected ..", string(ch)
			}
			l.next()
			token = ELLIPSIS
		} else {
			token = DOT
		}

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		runes := []rune{ch}
		for l.ch >= '0' && l.ch <= '9' {
			runes = append(runes, l.ch)
			l.next()
		}
		token = INT
		value = string(runes)

	case '"':
		runes := []rune{}
		for l.ch != '"' {
			c := l.ch
			if c < 0 {
				return pos, ILLEGAL, "didn't find end quote in string", string(ch)
			}
			if c == '\r' || c == '\n' {
				return pos, ILLEGAL, "can't have newline in string", string(ch)
			}
			if c == '\\' {
				l.next()
				switch l.ch {
				case '"', '\\':
					c = l.ch
				case 't':
					c = '\t'
				case 'r':
					c = '\r'
				case 'n':
					c = '\n'
				default:
					return pos, ILLEGAL, fmt.Sprintf("invalid string escape \\%c", l.ch), string(ch)
				}
			}
			runes = append(runes, c)
			l.next()
		}
		l.next()
		token = STR
		value = string(runes)

	default:
		token = ILLEGAL
		value = fmt.Sprintf("unexpected %c", ch)
	}
	return pos, token, value, string(ch)
}
