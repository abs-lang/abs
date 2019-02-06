package object

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/abs-lang/abs/ast"
	"github.com/abs-lang/abs/token"
)

type BuiltinFunction func(tok token.Token, args ...Object) Object

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
	Token token.Token
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
	Next() (Object, Object)
	Reset()
}

type Number struct {
	Token token.Token
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
	Token token.Token
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct {
	Token token.Token
}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Token token.Token
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
	Token      token.Token
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
	Token token.Token
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
	Token    token.Token
	Fn       BuiltinFunction
	Next     func() (Object, Object)
	Types    []string
	Iterable bool
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Token    token.Token
	Elements []Object
	position int
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Next() (Object, Object) {
	position := ao.position
	if len(ao.Elements) > position {
		ao.position = position + 1
		return &Number{Value: float64(position)}, ao.Elements[position]
	}

	return nil, nil
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
	Token    token.Token
	Pairs    map[HashKey]HashPair
	Position int
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) GetPair(key string) (HashPair, bool) {
	record, ok := h.Pairs[HashKey{Type: "STRING", Value: key}]

	return record, ok
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	// create stable key ordered output
	sort.Strings(pairs)

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// Pretty convoluted logic here we could
// refactor.
// First we sort the hash keys alphabetically
// and then we loop through them until
// we reach the required position within
// the loop.
func (h *Hash) Next() (Object, Object) {
	curPosition := 0
	pairs := make(map[string]HashPair)
	var keys []string
	for _, v := range h.Pairs {
		pairs[v.Key.Inspect()] = v
		keys = append(keys, v.Key.Inspect())
	}

	sort.Strings(keys)

	for _, k := range keys {
		if h.Position == curPosition {
			h.Position += 1
			return pairs[k].Key, pairs[k].Value
		}

		curPosition += 1
	}

	return nil, nil
}
func (h *Hash) Reset() {
	h.Position = 0
}
