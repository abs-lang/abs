package ast

import (
	"bytes"
	"strings"

	"github.com/abs-lang/abs/lexer"
	"github.com/abs-lang/abs/token"
)

// The base Node interface
type Node interface {
	PositionInterface
	TokenLiteral() string
	String() string
}

// All statement nodes implement this
type Statement interface {
	Node
	statementNode()
}

// All expression nodes implement this
type Expression interface {
	Node
	expressionNode()
}

// PositionInterface all nodes have this interface
type PositionInterface interface {
	SetPosition(int, *lexer.Lexer)
	GetLine() (int, int, string)
}

// Position defines the offset into the file where this node was seen
type Position struct {
	Position int          // records our Position in the file
	Lexer    *lexer.Lexer // pointer to our Lexer instance for error location in Eval()
}

// SetPosition sets the Position and Lexer for this node
func (p *Position) SetPosition(pos int, lex *lexer.Lexer) {
	p.Position = pos
	p.Lexer = lex
}

// GetLine returns (lineNum, column, thisLine) from this node's Position
func (p *Position) GetLine() (int, int, string) {
	pos := p.Position
	lineNum, column, thisLine := p.Lexer.GetLinePos(pos)
	return lineNum, column, thisLine
}

// Represents the whole program
// as a bunch of statements
type Program struct {
	Position   // our position in the file
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Statements
type AssignStatement struct {
	Position             // our position in the file
	Token    token.Token // the token.ASSIGN token
	Name     *Identifier
	Names    []Expression
	Value    Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) String() string {
	var out bytes.Buffer

	if as.Name != nil {
		out.WriteString(as.Name.String())
	}

	out.WriteString(" = ")

	if as.Value != nil {
		out.WriteString(as.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Position                // our position in the file
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Position               // our position in the file
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type BlockStatement struct {
	Position               // our position in the file
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Expressions
type Identifier struct {
	Position             // our position in the file
	Token    token.Token // the token.IDENT token
	Value    string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type Boolean struct {
	Position // our position in the file
	Token    token.Token
	Value    bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type NumberLiteral struct {
	Position // our position in the file
	Token    token.Token
	Value    float64
}

func (nl *NumberLiteral) expressionNode()      {}
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NumberLiteral) String() string       { return nl.Token.Literal }

type PrefixExpression struct {
	Position             // our position in the file
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Position             // our position in the file
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type CompoundAssignment struct {
	Position             // our position in the file
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ca *CompoundAssignment) expressionNode()      {}
func (ca *CompoundAssignment) TokenLiteral() string { return ca.Token.Literal }
func (ca *CompoundAssignment) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ca.Left.String())
	out.WriteString(" " + ca.Operator + " ")
	out.WriteString(ca.Right.String())
	out.WriteString(")")

	return out.String()
}

type MethodExpression struct {
	Position              // our position in the file
	Token     token.Token // The operator token, e.g. .
	Object    Expression
	Method    Expression
	Arguments []Expression
}

func (me *MethodExpression) expressionNode()      {}
func (me *MethodExpression) TokenLiteral() string { return me.Token.Literal }
func (me *MethodExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range me.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(me.Object.String())
	out.WriteString(".")
	out.WriteString(me.Method.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Position                // our position in the file
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type WhileExpression struct {
	Position                // our position in the file
	Token       token.Token // The 'while' token
	Condition   Expression
	Consequence *BlockStatement
}

func (ie *WhileExpression) expressionNode()      {}
func (ie *WhileExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString("while")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	return out.String()
}

type ForInExpression struct {
	Position                 // our position in the file
	Token    token.Token     // The 'for' token
	Block    *BlockStatement // The block executed inside the for loop
	Iterable Expression      // An expression that should return an iterable ([1, 2, 3] or x in 1..10)
	Key      string
	Value    string
}

func (fie *ForInExpression) expressionNode()      {}
func (fie *ForInExpression) TokenLiteral() string { return fie.Token.Literal }
func (fie *ForInExpression) String() string {
	var out bytes.Buffer

	out.WriteString("for ")

	if fie.Key != "" {
		out.WriteString(fie.Key + ", ")
	}
	out.WriteString(fie.Value)
	out.WriteString(" in ")
	out.WriteString(fie.Iterable.String())
	out.WriteString(fie.Block.String())

	return out.String()
}

type ForExpression struct {
	Position                   // our position in the file
	Token      token.Token     // The 'for' token
	Identifier string          // "x"
	Starter    Statement       // x = 0
	Closer     Statement       // x++
	Condition  Expression      // x < 1
	Block      *BlockStatement // The block executed inside the for loop
}

func (fe *ForExpression) expressionNode()      {}
func (fe *ForExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *ForExpression) String() string {
	var out bytes.Buffer

	out.WriteString("for ")

	out.WriteString(fe.Starter.String())
	out.WriteString(";")
	out.WriteString(fe.Condition.String())
	out.WriteString(";")
	out.WriteString(fe.Closer.String())
	out.WriteString(";")
	out.WriteString(fe.Block.String())

	return out.String()
}

type CommandExpression struct {
	Position             // our position in the file
	Token    token.Token // The command itself
	Value    string
}

func (ce *CommandExpression) expressionNode()      {}
func (ce *CommandExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CommandExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Token.Literal)

	return out.String()
}

type FunctionLiteral struct {
	Position               // our position in the file
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Position              // our position in the file
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Position // our position in the file
	Token    token.Token
	Value    string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type NullLiteral struct {
	Position // our position in the file
	Token    token.Token
}

func (nl *NullLiteral) expressionNode()      {}
func (nl *NullLiteral) TokenLiteral() string { return "null" }
func (nl *NullLiteral) String() string       { return "null" }

type ArrayLiteral struct {
	Position             // our position in the file
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Position             // our position in the file
	Token    token.Token // The [ token
	Left     Expression
	Index    Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type PropertyExpression struct {
	Position             // our position in the file
	Token    token.Token // The . token
	Object   Expression
	Property Expression
}

func (pe *PropertyExpression) expressionNode()      {}
func (pe *PropertyExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PropertyExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Object.String())
	out.WriteString(".")
	out.WriteString(pe.Property.String())
	out.WriteString(")")

	return out.String()
}

type HashLiteral struct {
	Position             // our position in the file
	Token    token.Token // the '{' token
	Pairs    map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
