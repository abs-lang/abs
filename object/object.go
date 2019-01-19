package object

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/abs-lang/abs/ast"
)

type BuiltinFunction func(args ...Object) Object

type ObjectType string

const (
	NULL_OBJ  = "NULL"
	ERROR_OBJ = "ERROR"

	NUMBER_OBJ  = "NUMBER"
	BOOLEAN_OBJ = "BOOLEAN"
	STRING_OBJ  = "STRING"

	RETURN_VALUE_OBJ = "RETURN_VALUE"

	FUNCTION_OBJ = "FUNCTION"
	BUILTIN_OBJ  = "BUILTIN"

	ARRAY_OBJ = "ARRAY"
	HASH_OBJ  = "HASH"
)

type HashKey struct {
	Type  ObjectType
	Value string
}

type Hashable interface {
	HashKey() HashKey
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Iterable interface {
	Next(int) (int, Object)
	Reset()
}

type Number struct {
	Value float64
}

func (n *Number) Type() ObjectType { return NUMBER_OBJ }

// If the number we're dealing with is
// an integer, print it as such (1.0000 becomes 1).
// If it's a float, let's remove as many zeroes
// as possible (1.10000 becomes 1.1).
func (n *Number) Inspect() string {
	if n.Value == float64(int64(n.Value)) {
		return fmt.Sprintf("%d", int64(n.Value))
	}
	return strconv.FormatFloat(n.Value, 'f', -1, 64)
}
func (n *Number) ZeroValue() float64 { return float64(0) }
func (n *Number) Int() int           { return int(n.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("f")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {")
	out.WriteString(f.Body.String())
	out.WriteString("}")

	return out.String()
}

// The String is a special fella.
//
// Like ints, or bools, you might
// think it will only have a Value
// property, the string itself.
//
// TA-DA! No, we also have an Ok
// property that is used when running
// shell commands -- since the shell
// will return strings.
//
// So, look at this:
//
// cmd = $(ls -la)
// type(cmd) // STRING
// cmd.ok // TRUE
//
// cmd = $(curlzzzzz)
// type(cmd) // STRING
// cmd.ok // FALSE
type String struct {
	Value string
	Ok    *Boolean // A special property to check whether a command exited correctly
}

func (s *String) Type() ObjectType  { return STRING_OBJ }
func (s *String) Inspect() string   { return s.Value }
func (s *String) ZeroValue() string { return "" }
func (s *String) HashKey() HashKey {
	return HashKey{Type: s.Type(), Value: s.Value}
}

type Builtin struct {
	Fn       BuiltinFunction
	Next     func(int) (int, Object)
	Types    []string
	Iterable bool
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
	position int
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Next(pos int) (int, Object) {
	position := ao.position
	if len(ao.Elements) > position {
		ao.position = position + 1
		return position, ao.Elements[position]
	}

	return -1, nil
}
func (ao *Array) Reset() {
	ao.position = 0
}
func (ao *Array) Homogeneous() bool {
	if ao.Empty() {
		return true
	}

	t := ao.Elements[0].Type()
	homogeneous := true

	for _, v := range ao.Elements {
		if v.Type() != t {
			homogeneous = false
		}
	}

	return homogeneous
}
func (ao *Array) Empty() bool {
	return len(ao.Elements) == 0
}
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
