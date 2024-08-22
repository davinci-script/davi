// DaVinci Script

package parser

import (
	"fmt"
	. "github.com/DavinciScript/Davi/lexer"
	"github.com/hokaccha/go-prettyjson"
	"strconv"
)

// Error is the error type returned by ParseExpression and ParseProgram when
// they encounter a syntax error. You can use this to get the location (line
// and column) of where the error occurred, as well as the error message.
type Error struct {
	Position Position
	Message  string
}

func (e Error) Error() string {
	return fmt.Sprintf("parse error at %d:%d: %s", e.Position.Line, e.Position.Column, e.Message)
}

type parser struct {
	lexer *Lexer
	pos   Position
	tok   Token
	val   string
}

func (p *parser) next() {
	p.pos, p.tok, p.val, _ = p.lexer.Next()
	if p.tok == ILLEGAL {
		p.error("%s", p.val)
	}
}

func (p *parser) error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	panic(Error{p.pos, message})
}

func (p *parser) expect(tok Token, context ...string) {
	if p.tok != tok {
		p.error("expected %s and not %s call_from: %s", tok, p.tok, context)
	}
	p.next()
}

func (p *parser) matches(operators ...Token) bool {
	for _, operator := range operators {
		if p.tok == operator {
			return true
		}
	}
	return false
}

// program = statement*
func (p *parser) program() *Program {
	statements := p.statements(EOF)
	return &Program{statements}
}

func (p *parser) statements(end Token) Block {
	statements := Block{}
	for p.tok != end && p.tok != EOF {
		statements = append(statements, p.statement())
	}
	return statements
}

// statement = if | while | for | return | function | assign | expression
// assign    = NAME ASSIGN expression |
//
//	call subscript ASSIGN expression |
//	call dot ASSIGN expression
func (p *parser) statement() Statement {
	switch p.tok {
	case IF:
		return p.if_()
	case WHILE:
		return p.while()
	case FOR:
		return p.for_()
	case RETURN:
		return p.return_()
	case FUNCTION:
		return p.function_()
	case CLASS:
		return p.class_()
	}
	pos := p.pos
	expr := p.expression()
	if p.tok == OBJECT_OPERATOR {
		pos = p.pos
		p.expect(OBJECT_OPERATOR, "object_operator")

		methodName := p.val
		p.expect(NAME, "object_operator")

		if p.tok == LPAREN {
			pos := p.pos
			p.next()
			args := []Expression{}
			gotComma := true
			gotEllipsis := false
			for p.tok != RPAREN && p.tok != EOF && !gotEllipsis {
				if !gotComma {
					p.error("expected , between arguments")
				}
				arg := p.expression()
				args = append(args, arg)
				if p.tok == ELLIPSIS {
					gotEllipsis = true
					p.next()
				}
				if p.tok == COMMA {
					gotComma = true
					p.next()
				} else {
					gotComma = false
				}
			}
			if p.tok != RPAREN && gotEllipsis {
				p.error("can only have ... after last argument")
			}
			p.expect(RPAREN, "object_operator")

			expr = &MethodCall{pos, expr, methodName, args}
		}

		return &ExpressionStatement{pos, expr}

	}
	if p.tok == ASSIGN {
		pos = p.pos
		switch expr.(type) {
		case *Variable, *Subscript:
			p.next()
			value := p.expression()
			return &Assign{pos, expr, value}
		default:
			p.error("expected name, subscript, or dot expression on left side of =")
		}
	}
	return &ExpressionStatement{pos, expr}
}

// block = LBRACE statement* RBRACE
func (p *parser) block() Block {
	p.expect(LBRACE, "block")
	body := p.statements(RBRACE)
	p.expect(RBRACE, "block")
	return body
}

// if = IF expression block |
//
//	IF expression block ELSE block |
//	IF expression block ELSE if
func (p *parser) if_() Statement {
	pos := p.pos
	p.expect(IF, "if_")
	condition := p.expression()
	body := p.block()
	var elseBody Block
	if p.tok == ELSE {
		p.next()
		if p.tok == LBRACE {
			elseBody = p.block()
		} else if p.tok == IF {
			elseBody = Block{p.if_()}
		} else {
			p.error("expected { or if after else, not %s", p.tok)
		}
	}
	return &If{pos, condition, body, elseBody}
}

// while = WHILE expression block
func (p *parser) while() Statement {
	pos := p.pos
	p.expect(WHILE, "while")
	condition := p.expression()
	body := p.block()
	return &While{pos, condition, body}
}

// for = FOR NAME IN expression block
func (p *parser) for_() Statement {
	pos := p.pos
	p.expect(FOR, "for_")
	p.expect(LPAREN, "for_")
	p.expect(DOLLAR, "for_")
	name := p.val
	p.expect(NAME, "for_")
	p.expect(IN, "for_")
	iterable := p.expression()
	p.expect(RPAREN, "for_")
	body := p.block()
	return &For{pos, name, iterable, body}
}

// return = RETURN expression
func (p *parser) return_() Statement {
	pos := p.pos
	p.expect(RETURN, "return_")
	result := p.expression()
	return &Return{pos, result}
}

// class = CLASS NAME block
func (p *parser) class_() Statement {

	p.next()
	pos := p.pos
	name := p.val

	p.expect(NAME, "class_")

	// Parse the class body
	p.expect(LBRACE, "class_")
	body := []Statement{}

	for p.tok != RBRACE && p.tok != EOF {
		body = append(body, p.statement())
	}

	p.expect(RBRACE, "class_")

	return &ClassDefinition{pos, name, nil, body}
}

// function = FUNCTION NAME params block |
//
//	FUNCTION params block
func (p *parser) function_() Statement {
	pos := p.pos
	p.expect(FUNCTION, "function_")
	if p.tok == NAME {
		name := p.val
		p.next()
		params, ellipsis := p.params()
		body := p.block()
		return &FunctionDefinition{pos, name, params, ellipsis, body}
	} else {
		params, ellipsis := p.params()
		body := p.block()
		expr := &FunctionExpression{pos, params, ellipsis, body}
		return &ExpressionStatement{pos, expr}
	}
}

// params = LPAREN RPAREN |
//
//	LPAREN NAME (COMMA NAME)* ELLIPSIS? COMMA? RPAREN |
func (p *parser) params() ([]string, bool) {
	p.expect(LPAREN, "params")
	params := []string{}
	gotComma := true
	gotEllipsis := false
	for p.tok != RPAREN && p.tok != EOF && !gotEllipsis {
		if !gotComma {
			p.error("expected , between parameters")
		}
		p.expect(DOLLAR, "params")
		param := p.val
		p.expect(NAME, "params")
		params = append(params, param)
		if p.tok == ELLIPSIS {
			gotEllipsis = true
			p.next()
		}
		if p.tok == COMMA {
			gotComma = true
			p.next()
		} else {
			gotComma = false
		}
	}
	if p.tok != RPAREN && gotEllipsis {
		p.error("can only have ... after last parameter")
	}
	p.expect(RPAREN, "params")
	return params, gotEllipsis
}

func (p *parser) binary(parseFunc func() Expression, operators ...Token) Expression {
	expr := parseFunc()
	for p.matches(operators...) {
		op := p.tok
		pos := p.pos
		p.next()
		right := parseFunc()
		expr = &Binary{pos, expr, op, right}
	}
	return expr
}

// expression = and (OR and)*
func (p *parser) expression() Expression {
	return p.binary(p.and, OR)
}

// and = not (AND not)*
func (p *parser) and() Expression {
	return p.binary(p.not, AND)
}

// not = NOT not | equality
func (p *parser) not() Expression {
	if p.tok == NOT {
		pos := p.pos
		p.next()
		operand := p.not()
		return &Unary{pos, NOT, operand}
	}
	return p.equality()
}

// equality = comparison ((EQUAL | NOTEQUAL) comparison)*
func (p *parser) equality() Expression {
	return p.binary(p.comparison, EQUAL, NOTEQUAL)
}

// comparison = addition ((LT | LTE | GT | GTE | IN) addition)*
func (p *parser) comparison() Expression {
	return p.binary(p.addition, LT, LTE, GT, GTE, IN)
}

// addition = multiply ((PLUS | MINUS) multiply)*
func (p *parser) addition() Expression {
	return p.binary(p.multiply, PLUS, MINUS)
}

// multiply = negative ((TIMES | DIVIDE | MODULO) negative)*
func (p *parser) multiply() Expression {
	return p.binary(p.negative, TIMES, DIVIDE, MODULO)
}

// negative = MINUS negative | call
func (p *parser) negative() Expression {
	if p.tok == MINUS {
		pos := p.pos
		p.next()
		operand := p.negative()
		return &Unary{pos, MINUS, operand}
	}
	return p.call()
}

// call      = primary (args | subscript | dot)*
// args      = LPAREN RPAREN |
//
//	LPAREN expression (COMMA expression)* ELLIPSIS? COMMA? RPAREN)
//
// subscript = LBRACKET expression RBRACKET
// dot       = DOT NAME
func (p *parser) call() Expression {
	expr := p.primary()
	for p.matches(LPAREN, LBRACKET, DOT) {
		if p.tok == LPAREN {
			pos := p.pos
			p.next()
			args := []Expression{}
			gotComma := true
			gotEllipsis := false
			for p.tok != RPAREN && p.tok != EOF && !gotEllipsis {
				if !gotComma {
					p.error("expected , between arguments")
				}
				arg := p.expression()
				args = append(args, arg)
				if p.tok == ELLIPSIS {
					gotEllipsis = true
					p.next()
				}
				if p.tok == COMMA {
					gotComma = true
					p.next()
				} else {
					gotComma = false
				}
			}
			if p.tok != RPAREN && gotEllipsis {
				p.error("can only have ... after last argument")
			}
			p.expect(RPAREN, "call")

			expr = &Call{pos, expr, args, gotEllipsis}
		} else if p.tok == LBRACKET {
			pos := p.pos
			p.next()
			subscript := p.expression()
			p.expect(RBRACKET, "call")
			//p.expect(SEMI)
			expr = &Subscript{pos, expr, subscript}
		} else {
			pos := p.pos
			p.next()
			subscript := &Literal{p.pos, p.val}
			if p.tok == STR {
				p.next()
				if p.tok == DOT {
					p.next()
				}
				p.primary()
			} else {
				p.expect(NAME, "call-name")
				expr = &Subscript{pos, expr, subscript}
			}
		}
	}
	return expr
}

// primary = NAME | INT | STR | TRUE | FALSE | NIL | list | map |
//
//	FUNCTION params block |
//	LPAREN expression RPAREN
func (p *parser) primary() Expression {
	switch p.tok {
	case NAME:
		name := p.val
		pos := p.pos
		//function := &Variable{pos, name}
		p.next()

		//if name == "echo" {
		//	value := p.expression()
		//	args := []Expression{value}
		//	return &Call{pos, function, args, false}
		//}

		return &Variable{pos, name}
	case DOLLAR:
		p.expect(DOLLAR, "primary")
		name := p.val
		pos := p.pos
		p.next()
		return &Variable{pos, name}
	case INT:
		val := p.val
		pos := p.pos
		p.next()
		n, err := strconv.Atoi(val)
		if err != nil {
			// Tokenizer should never give us this
			panic(fmt.Sprintf("tokenizer gave INT token that isn't an int: %s", val))
		}
		return &Literal{pos, n}
	case STR:
		val := p.val
		pos := p.pos
		p.next()
		return &Literal{pos, val}
	case TRUE:
		pos := p.pos
		p.next()
		return &Literal{pos, true}
	case FALSE:
		pos := p.pos
		p.next()
		return &Literal{pos, false}
	case NIL:
		pos := p.pos
		p.next()
		return &Literal{pos, nil}
	case LBRACKET:
		return p.list()
	case LBRACE:
		return p.map_()
	case FUNCTION:
		pos := p.pos
		p.next()
		args, ellipsis := p.params()
		body := p.block()
		return &FunctionExpression{pos, args, ellipsis, body}
	case LPAREN:
		p.next()
		expr := p.expression()
		p.expect(RPAREN, "primary")
		//p.expect(SEMI)
		return expr
	case SEMI:
		p.next()
		pos := p.pos
		return &SemiTag{pos}
	case NEW:
		pos := p.pos
		p.expect(NEW, "new") // Move past the 'NEW' token

		// Expect the name of the class or type to be instantiated
		if p.tok != NAME {
			p.error("expected class name after 'new'")
			return nil
		}
		className := p.val
		p.expect(NAME, "new")
		args, _ := p.params()

		return &NewExpression{pos, className, args}
	//case OBJECT_OPERATOR:
	//	pos := p.pos
	//	p.next() // Move past the OBJECT_OPERATOR token
	//
	//	if p.tok != NAME {
	//		p.error("expected a method or property name after '->'")
	//		return nil
	//	}
	//
	//	methodName := p.val
	//	p.next() // Move past the method name
	//
	//	//// Assuming you want to handle method calls like $test->call(args...)
	//	if p.tok == LPAREN {
	//		p.next() // Skip '('
	//		args := []Expression{}
	//		p.expect(RPAREN) // Ensure method call is closed with ')'
	//		return &MethodCall{pos, nil, methodName, args}
	//	}
	//
	//	return &PropertyAccess{pos, nil, methodName}

	default:
		formatter := prettyjson.NewFormatter()
		output, _ := formatter.Marshal(p.val)
		fmt.Println(string(output))
		p.error("expected expression, ___%s___ ", p.tok)
		return nil
	}
}

// list = LBRACKET RBRACKET |
//
//	LBRACKET expression (COMMA expression)* COMMA? RBRACKET
func (p *parser) list() Expression {
	pos := p.pos
	p.expect(LBRACKET, "list")
	values := []Expression{}
	gotComma := true
	for p.tok != RBRACKET && p.tok != EOF {
		if !gotComma {
			p.error("expected , between list elements")
		}
		value := p.expression()
		values = append(values, value)
		if p.tok == COMMA {
			gotComma = true
			p.next()
		} else {
			gotComma = false
		}
	}
	p.expect(RBRACKET, "list")
	//p.expect(SEMI)
	return &List{pos, values}
}

// map = LBRACE RBRACE |
//
//	LBRACE expression COLON expression
//	       (COMMA expression COLON expression)* COMMA? RBRACE
func (p *parser) map_() Expression {
	pos := p.pos
	p.expect(LBRACE, "map_")
	items := []MapItem{}
	gotComma := true
	for p.tok != RBRACE && p.tok != EOF {
		if !gotComma {
			p.error("expected , between map items")
		}
		key := p.expression()
		p.expect(COLON, "map_")
		value := p.expression()
		items = append(items, MapItem{key, value})
		if p.tok == COMMA {
			gotComma = true
			p.next()
		} else {
			gotComma = false
		}
	}
	p.expect(RBRACE, "map_")
	return &Map{pos, items}
}

// ParseExpression parses a single expression into an Expression interface
// (can be one of many expression types). If the expression parses correctly,
// return an Expression and nil. If there's a syntax error, return nil and
// a parser.Error value.
func ParseExpression(input []byte) (e Expression, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Convert to parser.Error or re-panic
			err = r.(Error)
		}
	}()
	l := NewLexer(input)
	p := parser{lexer: l}
	p.next()
	return p.expression(), nil
}

// ParseProgram parses an entire program and returns a *Program (which is
// basically a list of statements). If the program parses correctly, return
// a *Program and nil. If there's a syntax error, return nil and a
// parser.Error value.
func ParseProgram(input []byte) (prog *Program, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Convert to parser.Error or re-panic
			err = r.(Error)
		}
	}()
	l := NewLexer(input)
	p := parser{lexer: l}
	p.next()
	return p.program(), nil
}
