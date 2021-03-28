package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/abs-lang/abs/ast"
	"github.com/abs-lang/abs/lexer"
	"github.com/abs-lang/abs/token"
)

const (
	_ int = iota
	LOWEST
	AND         // && or ||
	EQUALS      // == or !=
	LESSGREATER // > or <
	SUM         // + or -
	PRODUCT     // * or / or ^
	RANGE       // ..
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
	QUESTION    // some?.function() or some?.property
	DOT         // some.function() or some.property
	HIGHEST     // special preference for -x or +y
)

var precedences = map[token.TokenType]int{
	token.AND:           AND,
	token.OR:            AND,
	token.BIT_AND:       AND,
	token.BIT_XOR:       AND,
	token.BIT_RSHIFT:    AND,
	token.BIT_LSHIFT:    AND,
	token.PIPE:          AND,
	token.EQ:            EQUALS,
	token.NOT_EQ:        EQUALS,
	token.TILDE:         EQUALS,
	token.IN:            EQUALS,
	token.NOT_IN:        EQUALS,
	token.COMMA:         EQUALS,
	token.LT:            LESSGREATER,
	token.LT_EQ:         LESSGREATER,
	token.GT:            LESSGREATER,
	token.GT_EQ:         LESSGREATER,
	token.COMBINED_COMP: LESSGREATER,
	token.PLUS:          SUM,
	token.MINUS:         SUM,
	token.SLASH:         PRODUCT,
	token.ASTERISK:      PRODUCT,
	token.EXPONENT:      PRODUCT,
	token.MODULO:        PRODUCT,
	token.COMP_PLUS:     EQUALS,
	token.COMP_MINUS:    EQUALS,
	token.COMP_SLASH:    EQUALS,
	token.COMP_ASTERISK: EQUALS,
	token.COMP_EXPONENT: EQUALS,
	token.COMP_MODULO:   EQUALS,
	token.RANGE:         RANGE,
	token.LPAREN:        CALL,
	token.LBRACKET:      INDEX,
	token.QUESTION:      QUESTION,
	token.DOT:           DOT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	// support assignment to index expressions: a[0] = 1, h["a"] = 1
	prevIndexExpression *ast.IndexExpression

	// support assignment to hash property h.a = 1
	prevPropertyExpression *ast.PropertyExpression

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.NUMBER, p.ParseNumberLiteral)
	p.registerPrefix(token.STRING, p.ParseStringLiteral)
	p.registerPrefix(token.NULL, p.ParseNullLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.PLUS, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TILDE, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.ParseBoolean)
	p.registerPrefix(token.FALSE, p.ParseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.ParseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.ParseHashLiteral)
	p.registerPrefix(token.COMMAND, p.parseCommand)
	p.registerPrefix(token.BREAK, p.parseBreak)
	p.registerPrefix(token.CONTINUE, p.parseContinue)
	p.registerPrefix(token.CURRENT_ARGS, p.parseCurrentArgsLiteral)
	p.registerPrefix(token.AT, p.parseDecorator)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.QUESTION, p.parseQuestionExpression)
	p.registerInfix(token.DOT, p.parseDottedExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.EXPONENT, p.parseInfixExpression)
	p.registerInfix(token.MODULO, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.COMP_PLUS, p.parseCompoundAssignment)
	p.registerInfix(token.COMP_MINUS, p.parseCompoundAssignment)
	p.registerInfix(token.COMP_SLASH, p.parseCompoundAssignment)
	p.registerInfix(token.COMP_EXPONENT, p.parseCompoundAssignment)
	p.registerInfix(token.COMP_MODULO, p.parseCompoundAssignment)
	p.registerInfix(token.COMP_ASTERISK, p.parseCompoundAssignment)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.TILDE, p.parseInfixExpression)
	p.registerInfix(token.IN, p.parseInfixExpression)
	p.registerInfix(token.NOT_IN, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GT_EQ, p.parseInfixExpression)
	p.registerInfix(token.COMBINED_COMP, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.BIT_AND, p.parseInfixExpression)
	p.registerInfix(token.BIT_XOR, p.parseInfixExpression)
	p.registerInfix(token.PIPE, p.parseInfixExpression)
	p.registerInfix(token.BIT_RSHIFT, p.parseInfixExpression)
	p.registerInfix(token.BIT_LSHIFT, p.parseInfixExpression)
	p.registerInfix(token.RANGE, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

	if p.curTokenIs(token.ILLEGAL) {
		msg := fmt.Sprintf(`Illegal token '%s'`, p.curToken.Literal)
		p.reportError(msg, p.curToken)
	}
}

func (p *Parser) curTokenIs(typ token.TokenType) bool {
	return p.curToken.Type == typ
}

func (p *Parser) peekTokenIs(typ token.TokenType) bool {
	return p.peekToken.Type == typ
}

func (p *Parser) expectPeek(typ token.TokenType) bool {
	if p.peekTokenIs(typ) {
		p.nextToken()
		return true
	}
	p.peekError(p.curToken)
	return false
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) reportError(err string, tok token.Token) {
	// report error at token location
	lineNum, column, errorLine := p.l.ErrorLine(tok.Position)
	msg := fmt.Sprintf("%s\n\t[%d:%d]\t%s", err, lineNum, column, errorLine)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(tok token.Token) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", tok.Type, p.peekToken.Type)
	p.reportError(msg, tok)
}

func (p *Parser) noPrefixParseFnError(tok token.Token) {
	msg := fmt.Sprintf("no prefix parse function for '%s' found", tok.Literal)
	p.reportError(msg, tok)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	if p.curToken.Type == token.RETURN {
		return p.parseReturnStatement()
	}

	statement := p.parseAssignStatement()
	if statement != nil {
		return statement
	}

	return p.parseExpressionStatement()
}

// Rewinds the parser. This method
// is fairly inefficient as it starts
// from scratch.
//
// In the context of scripting, though,
// it won't cause a crazy delay as we're
// not parsing a book.
func (p *Parser) Rewind(pos int) {
	p.l.Rewind(0)

	for p.l.CurrentPosition() < pos {
		p.nextToken()
	}
}

func (p *Parser) parseDestructuringIdentifiers() []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(token.ASSIGN) {
		return list
	}

	list = append(list, p.parseIdentifier())

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.peekTokenIs(token.ASSIGN) {
		return nil
	}

	return list
}

// assign to variable: x = y
// destructuring assignment: x, y = [z, zz]
// assign to index expressions: a[0] = 1, h["a"] = 1
// assign to hash property expressions: h.a = 1
func (p *Parser) parseAssignStatement() ast.Statement {
	stmt := &ast.AssignStatement{}

	// Is this a regular x = y assignment?
	if p.peekTokenIs(token.COMMA) {
		lexerPosition := p.l.CurrentPosition()
		// Let's figure out if we are destructuring x, y = [z, zz]
		if !p.curTokenIs(token.IDENT) {
			return nil
		}

		stmt.Names = p.parseDestructuringIdentifiers()

		if !p.peekTokenIs(token.ASSIGN) {
			p.Rewind(lexerPosition)
			return nil
		}
	} else if p.curTokenIs(token.IDENT) {
		stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	} else if p.curTokenIs(token.ASSIGN) {
		stmt.Token = p.curToken
		if p.prevIndexExpression != nil {
			// support assignment to indexed expressions: a[0] = 1, h["a"] = 1
			stmt.Index = p.prevIndexExpression
			p.nextToken()
			stmt.Value = p.parseExpression(LOWEST)
			// consume the IndexExpression
			p.prevIndexExpression = nil

			if p.peekTokenIs(token.SEMICOLON) {
				p.nextToken()
			}

			return stmt
		}
		if p.prevPropertyExpression != nil {
			// support assignment to hash properties: h.a = 1
			stmt.Property = p.prevPropertyExpression
			p.nextToken()
			stmt.Value = p.parseExpression(LOWEST)
			// consume the PropertyExpression
			p.prevPropertyExpression = nil

			if p.peekTokenIs(token.SEMICOLON) {
				p.nextToken()
			}

			return stmt
		}
	}

	if !p.peekTokenIs(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Token = p.curToken
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// return x
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	returnToken := p.curToken

	// return;
	if p.peekTokenIs(token.SEMICOLON) {
		stmt.ReturnValue = &ast.NullLiteral{Token: p.curToken}
	} else if p.peekTokenIs(token.RBRACE) || p.peekTokenIs(token.EOF) {
		// return
		stmt.ReturnValue = &ast.NullLiteral{Token: returnToken}
	} else {
		// return xyz
		p.nextToken()
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// (x * y) + z
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

// var
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// 1 or 1.1 or 1k
func (p *Parser) ParseNumberLiteral() ast.Expression {
	lit := &ast.NumberLiteral{Token: p.curToken}
	var abbr float64
	var ok bool
	number := p.curToken.Literal

	// Check if the last character of this number is an abbreviation
	if abbr, ok = token.NumberAbbreviations[strings.ToLower(string(number[len(number)-1]))]; ok {
		number = p.curToken.Literal[:len(p.curToken.Literal)-1]
	}

	value, err := strconv.ParseFloat(number, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as number", number)
		p.reportError(msg, p.curToken)
		return nil
	}

	if abbr != 0 {
		value *= abbr
	}

	lit.Value = value

	return lit
}

// "some"
func (p *Parser) ParseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// null
func (p *Parser) ParseNullLiteral() ast.Expression {
	return &ast.NullLiteral{Token: p.curToken}
}

// !x
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	precedence := PREFIX

	// When +- are used as prefixes, we want them to have
	// the highest priority, so that -5.clamp(4, 5) is read
	// as (-5).clamp(4, 5) = 4 instead of
	// -(5.clamp(4,5)) = -5
	if p.curTokenIs(token.PLUS) || p.curTokenIs(token.MINUS) {
		precedence = HIGHEST
	}
	p.nextToken()

	expression.Right = p.parseExpression(precedence)

	return expression
}

// x * x
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// x += x
func (p *Parser) parseCompoundAssignment(left ast.Expression) ast.Expression {
	expression := &ast.CompoundAssignment{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// some.function() or some.property
func (p *Parser) parseDottedExpression(object ast.Expression) ast.Expression {
	t := p.curToken
	precedence := p.curPrecedence()
	p.nextToken()

	// Here we try to figure out if
	// we're in front of a method or
	// a property accessor.
	//
	// If the token after the identifier
	// is a (, then we're expecting this
	// to me a method (x.f()), else we'll
	// assume it's a property (x.p).
	if p.peekTokenIs(token.LPAREN) {
		exp := &ast.MethodExpression{Token: t, Object: object}
		exp.Method = p.parseExpression(precedence)
		p.nextToken()
		exp.Arguments = p.parseExpressionList(token.RPAREN)
		return exp
	} else {
		// support assignment to hash property h.a = 1
		exp := &ast.PropertyExpression{Token: t, Object: object}
		exp.Property = p.parseIdentifier()
		p.prevPropertyExpression = exp
		p.prevIndexExpression = nil
		return exp
	}
}

// some?.function() or some?.property
// Here we skip the "?" and parse the expression as a regular dotted one.
// When we're back, we mark the expression as optional.
func (p *Parser) parseQuestionExpression(object ast.Expression) ast.Expression {
	p.nextToken()
	exp := p.parseDottedExpression(object)

	switch res := exp.(type) {
	case *ast.PropertyExpression:
		res.Optional = true
		return res
	case *ast.MethodExpression:
		res.Optional = true
		return res
	default:
		return exp
	}
}

// some.function()
func (p *Parser) parseMethodExpression(object ast.Expression) ast.Expression {
	exp := &ast.MethodExpression{Token: p.curToken, Object: object}
	precedence := p.curPrecedence()
	p.nextToken()
	exp.Method = p.parseExpression(precedence)
	p.nextToken()
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

// true
func (p *Parser) ParseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// if x {
//   return x
// } else if y {
//   return y
// } else {
//   return z
// }
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	scenarios := []*ast.Scenario{}

	p.nextToken()
	scenario := &ast.Scenario{}
	scenario.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	scenario.Consequence = p.parseBlockStatement()
	scenarios = append(scenarios, scenario)
	// If we encounter ELSEs then let's add more
	// scenarios to our expression.
	for p.peekTokenIs(token.ELSE) {
		p.nextToken()
		p.nextToken()
		scenario := &ast.Scenario{}

		// ELSE IF
		if p.curTokenIs(token.IF) {
			p.nextToken()
			scenario.Condition = p.parseExpression(LOWEST)

			if !p.expectPeek(token.LBRACE) {
				return nil
			}
		} else {
			// This is a regular ELSE block.
			//
			// In order not to have a weird data structure
			// representing an IF expression, we simply define
			// it as a list of scenarios.
			// In case a simple ELSE if encountered, we set the
			// condition of this scenario to true, so that it always
			// evaluates to true.
			tok := &token.Token{Position: -99, Literal: "true", Type: token.LookupIdent(token.TRUE)}
			scenario.Condition = &ast.Boolean{Token: *tok, Value: true}
		}

		scenario.Consequence = p.parseBlockStatement()
		scenarios = append(scenarios, scenario)
	}

	expression.Scenarios = scenarios
	return expression
}

// while true {
// 	echo("true")
// }
func (p *Parser) parseWhileExpression() ast.Expression {
	expression := &ast.WhileExpression{Token: p.curToken}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()
	return expression
}

// We first try parsing the code as a regular for loop.
// If we realize this is a for .. in we will then switch
// around.
func (p *Parser) parseForExpression() ast.Expression {
	expression := &ast.ForExpression{Token: p.curToken}
	p.nextToken()

	if !p.curTokenIs(token.IDENT) {
		return nil
	}

	if !p.peekTokenIs(token.ASSIGN) {
		return p.parseForInExpression(expression)
	}

	expression.Identifier = p.curToken.Literal
	expression.Starter = p.parseAssignStatement()

	if expression.Starter == nil {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if expression.Condition == nil {
		return nil
	}
	p.nextToken()
	p.nextToken()
	expression.Closer = p.parseAssignStatement()
	if expression.Closer == nil {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Block = p.parseBlockStatement()

	return expression
}

// for x in [1,2,3] {
// 	echo("true")
// }
func (p *Parser) parseForInExpression(initialExpression *ast.ForExpression) ast.Expression {
	expression := &ast.ForInExpression{Token: initialExpression.Token}

	if !p.curTokenIs(token.IDENT) {
		return nil
	}

	val := p.curToken.Literal
	var key string
	p.nextToken()

	if p.curTokenIs(token.COMMA) {
		p.nextToken()

		if !p.curTokenIs(token.IDENT) {
			return nil
		}

		key = val
		val = p.curToken.Literal
		p.nextToken()
	}

	expression.Key = key
	expression.Value = val

	if !p.curTokenIs(token.IN) {
		return nil
	}
	p.nextToken()

	expression.Iterable = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Block = p.parseBlockStatement()

	// for x in [] {
	// 	echo("shouldn't be here")
	// } else {
	//  echo("ok")
	// }
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

// { x + 1 }
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// f() {
//   return 1
// }
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = p.curToken.Literal
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

// @decorator
// @decorator(1, 2)
func (p *Parser) parseDecorator() ast.Expression {
	dc := &ast.Decorator{Token: p.curToken}
	// A decorator always preceeds a
	// named function or another decorator,
	// so once we're done
	// with our "own" parsing, we can defer
	// to parsing the next statement, which
	// should either return another decorator
	// or a function.
	defer (func() {
		p.nextToken()
		exp := p.parseExpressionStatement()

		switch fn := exp.Expression.(type) {
		case *ast.FunctionLiteral:
			if fn.Name == "" {
				p.reportError("a decorator should decorate a named function", dc.Token)
			}

			dc.Decorated = fn
		case *ast.Decorator:
			dc.Decorated = fn
		default:
			p.reportError("a decorator should decorate a named function", dc.Token)
		}
	})()

	p.nextToken()
	exp := p.parseExpressionStatement()
	dc.Expression = exp.Expression

	return dc
}

// ...
func (p *Parser) parseCurrentArgsLiteral() ast.Expression {
	return &ast.CurrentArgsLiteral{Token: p.curToken}
}

// f(x, y = 2)
func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	parameters := []*ast.Parameter{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return parameters
	}

	p.nextToken()

	param, foundOptionalParameter := p.parseFunctionParameter()
	parameters = append(parameters, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		param, optional := p.parseFunctionParameter()

		if foundOptionalParameter && !optional {
			p.reportError("found mandatory parameter after optional one", p.curToken)
		}

		if optional {
			foundOptionalParameter = true
		}

		parameters = append(parameters, param)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return parameters
}

// parse a single function parameter
// x
// x = 2
func (p *Parser) parseFunctionParameter() (param *ast.Parameter, optional bool) {
	// first, parse the identifier
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// if we find a comma or the closing parenthesis, then this
	// parameter is done eg. fn(x, y, z)
	if p.peekTokenIs(token.COMMA) || p.peekTokenIs(token.RPAREN) {
		return &ast.Parameter{Identifier: ident, Default: nil}, false
	}

	// else, we are in front of an optional parameter
	// fn(x = 2)
	// if the next token is not an assignment, though, there's
	// a major problem
	if !p.peekTokenIs(token.ASSIGN) {
		p.reportError("invalid parameter format", p.curToken)
		return &ast.Parameter{Identifier: ident, Default: nil}, false
	}

	// skip to the =
	p.nextToken()
	// skip to the default value of the parameter
	p.nextToken()
	// parse this default value as an expression
	// this allows for funny stuff like:
	// fn(x = 1)
	// fn(x = "")
	// fn(x = null)
	// fn(x = {})
	// fn(x = [1, 2, 3, 4])
	// fn(x = [1, 2, 3, 4].filter(f(x) {x > 2}) <--- very funny but ¯\_(ツ)_/¯
	exp := p.parseExpression(LOWEST)

	return &ast.Parameter{Identifier: ident, Default: exp}, true
}

// function()
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// [1, 2, 3]
func (p *Parser) ParseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

// some["thing"] or some[1:10]
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	if p.peekTokenIs(token.COLON) {
		exp.Index = &ast.NumberLiteral{Value: 0, Token: token.Token{Type: token.NUMBER, Position: 0, Literal: "0"}}
		exp.IsRange = true
	} else {
		p.nextToken()
		exp.Index = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.COLON) {
		exp.IsRange = true
		p.nextToken()

		if p.peekTokenIs(token.RBRACKET) {
			exp.End = nil
		} else {
			p.nextToken()
			exp.End = p.parseExpression(LOWEST)
		}
	}

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	// support assignment to index expression: a[0] = 1
	p.prevIndexExpression = exp
	p.prevPropertyExpression = nil

	return exp
}

// {"a": "b"}
func (p *Parser) ParseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseCommand() ast.Expression {
	cmd := &ast.CommandExpression{Token: p.curToken, Value: p.curToken.Literal}
	return cmd
}

// We don't really have to do anything when comments
// come in, we can simply ignore them
func (p *Parser) parseComment() ast.Expression {
	return nil
}

func (p *Parser) parseBreak() ast.Expression {
	return &ast.BreakStatement{Token: p.curToken}
}

func (p *Parser) parseContinue() ast.Expression {
	return &ast.ContinueStatement{Token: p.curToken}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
