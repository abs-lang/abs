package parser

import (
	"fmt"
	"strconv"

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
	DOT         // some.function() or some | function() or some.property
)

var precedences = map[token.TokenType]int{
	token.AND:           AND,
	token.OR:            AND,
	token.EQ:            EQUALS,
	token.NOT_EQ:        EQUALS,
	token.TILDE:         EQUALS,
	token.IN:            EQUALS,
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
	token.COMP_PLUS:     SUM,
	token.COMP_MINUS:    SUM,
	token.COMP_SLASH:    PRODUCT,
	token.COMP_ASTERISK: PRODUCT,
	token.COMP_EXPONENT: PRODUCT,
	token.COMP_MODULO:   PRODUCT,
	token.RANGE:         RANGE,
	token.LPAREN:        CALL,
	token.LBRACKET:      INDEX,
	token.DOT:           DOT,
	token.PIPE:          DOT,
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
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.ParseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.ParseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.ParseHashLiteral)
	p.registerPrefix(token.COMMAND, p.parseCommand)
	p.registerPrefix(token.COMMENT, p.parseComment)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.DOT, p.parseDottedExpression)
	p.registerInfix(token.PIPE, p.parseMethodExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
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
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GT_EQ, p.parseInfixExpression)
	p.registerInfix(token.COMBINED_COMP, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
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
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
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

// x = y
// [x] = [y]
func (p *Parser) parseAssignStatement() ast.Statement {
	stmt := &ast.AssignStatement{}

	// Is this a regular x = y assignment?
	if !p.curTokenIs(token.IDENT) {
		lexerPosition := p.l.CurrentPosition()
		// Let's figure out if we are destructuring [x] = [y]
		if !p.curTokenIs(token.LBRACKET) {
			return nil
		}

		arr := p.ParseArrayLiteral()

		if arr != nil && len(arr.(*ast.ArrayLiteral).Elements) == 0 {
			p.Rewind(lexerPosition)
			return nil
		}

		stmt.Names = arr.(*ast.ArrayLiteral).Elements

		if !p.peekTokenIs(token.ASSIGN) {
			p.Rewind(lexerPosition)
			return nil
		}
	} else {
		stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
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

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

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
		p.noPrefixParseFnError(p.curToken.Type)
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

// 1
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

// "some"
func (p *Parser) ParseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// !x
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

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
		exp := &ast.PropertyExpression{Token: t, Object: object}

		if !p.curTokenIs(token.IDENT) {
			p.errors = append(p.errors, fmt.Sprintf("property needs to be an identifier, got '%s'", p.curToken.Literal))
		}

		exp.Property = p.parseIdentifier()
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
func (p *Parser) parseBoolean() ast.Expression {
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
// } esle {
//   return y
// }
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

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

// f(x, y)
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
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

// some["thing"]
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

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
	return &ast.CommandExpression{Token: p.curToken, Value: p.curToken.Literal}
}

// We don't really have to do anything when comments
// come in, we can simply ignore them
func (p *Parser) parseComment() ast.Expression {
	return nil
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
