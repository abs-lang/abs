package ast

import (
	"bytes"
	"strings"

	"github.com/abs-lang/abs/token"
)

// The base Node interface
type Node interface {
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

// Represents the whole program
// as a bunch of statements
type Program struct {
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
	Token    token.Token // the token.ASSIGN token
	Name     *Identifier
	Names    []Expression
	Index    *IndexExpression    // support assignment to indexed expressions: a[0] = 1, h["a"] = 1
	Property *PropertyExpression // support assignment to hash properties: h.a = 1
	Value    Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) String() string {
	var out bytes.Buffer

	if as.Name != nil {
		out.WriteString(as.Name.String())
	} else if len(as.Names) > 0 {
		out.WriteString(as.Names[0].String())
		for i := 1; i < len(as.Names); i++ {
			out.WriteString(", ")
			out.WriteString(as.Names[i].String())
		}
	} else if as.Index != nil {
		out.WriteString(as.Index.String())
	} else if as.Property != nil {
		out.WriteString(as.Property.String())
	}

	out.WriteString(" = ")

	if as.Value != nil {
		out.WriteString(as.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type BreakStatement struct {
	Token token.Token // the 'break' token
}

func (bs *BreakStatement) expressionNode()      {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string {
	return "break;"
}

type ContinueStatement struct {
	Token token.Token // the 'continue' token
}

func (cs *ContinueStatement) expressionNode()      {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string {
	return "continue;"
}

type ReturnStatement struct {
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
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (nl *NumberLiteral) expressionNode()      {}
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NumberLiteral) String() string       { return nl.Token.Literal }

type PrefixExpression struct {
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
	Token     token.Token // The operator token, e.g. .
	Object    Expression
	Method    Expression
	Arguments []Expression
	Optional  bool
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
	if me.Optional {
		out.WriteString("?")
	}
	out.WriteString(".")
	out.WriteString(me.Method.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// A scenario is used to
// represent a path within an
// IF block (if x = 2 { return x}).
// It has a condition (x = 2)
// and a consenquence (return x).
type Scenario struct {
	Condition   Expression
	Consequence *BlockStatement
}

type IfExpression struct {
	Token     token.Token // The 'if' token
	Scenarios []*Scenario
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	for i, s := range ie.Scenarios {
		if i != 0 {
			out.WriteString("else")
			out.WriteString(" ")
		}

		out.WriteString("if")
		out.WriteString(s.Condition.String())
		out.WriteString(" ")
		out.WriteString(s.Consequence.String())
	}

	return out.String()
}

type WhileExpression struct {
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
	Token       token.Token     // The 'for' token
	Block       *BlockStatement // The block executed inside the for loop
	Iterable    Expression      // An expression that should return an iterable ([1, 2, 3] or x in 1..10)
	Key         string
	Value       string
	Alternative *BlockStatement
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

	if fie.Alternative != nil {
		out.WriteString("else")
		out.WriteString(fie.Alternative.String())
	}

	return out.String()
}

type ForExpression struct {
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
	Token token.Token // The command itself
	Value string
}

func (ce *CommandExpression) expressionNode()      {}
func (ce *CommandExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CommandExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Token.Literal)

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Name       string      // identifier for this function
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

type Decorator struct {
	Token     token.Token // @
	Name      string
	Arguments []Expression
	Decorated *FunctionLiteral
}

func (dc *Decorator) expressionNode()      {}
func (dc *Decorator) TokenLiteral() string { return dc.Token.Literal }
func (dc *Decorator) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range dc.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(dc.TokenLiteral())
	out.WriteString(dc.Name)
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(") ")
	out.WriteString(dc.Decorated.String())

	return out.String()
}

type CurrentArgsLiteral struct {
	Token token.Token // ...
}

func (cal *CurrentArgsLiteral) expressionNode()      {}
func (cal *CurrentArgsLiteral) TokenLiteral() string { return cal.Token.Literal }
func (cal *CurrentArgsLiteral) String() string {
	return "..."
}

type CallExpression struct {
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
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type NullLiteral struct {
	Token token.Token
}

func (nl *NullLiteral) expressionNode()      {}
func (nl *NullLiteral) TokenLiteral() string { return "null" }
func (nl *NullLiteral) String() string       { return "null" }

type ArrayLiteral struct {
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

// IndexExpression allows accessing a single index, or a range,
// over a string or an array.
//
// array[1:10]	-> left[index:end]
// array[1] 	-> left[index]
// string[1] 	-> left[index]
type IndexExpression struct {
	Token   token.Token // The [ token
	Left    Expression  // the argument on which the index is access eg array of array[1]
	Index   Expression  // the left-most index eg. 1 in array[1] or array[1:10]
	IsRange bool        // whether the expression is a range (1:10)
	End     Expression  // the end of the range, if the expression is a range
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")

	if ie.IsRange {
		out.WriteString(ie.Index.String() + ":" + ie.End.String())
	} else {
		out.WriteString(ie.Index.String())
	}

	out.WriteString("])")

	return out.String()
}

type PropertyExpression struct {
	Token    token.Token // The . token
	Object   Expression
	Property Expression
	Optional bool
}

func (pe *PropertyExpression) expressionNode()      {}
func (pe *PropertyExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PropertyExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Object.String())

	if pe.Optional {
		out.WriteString("?")
	}

	out.WriteString(".")
	out.WriteString(pe.Property.String())
	out.WriteString(")")

	return out.String()
}

type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
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
