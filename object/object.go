package object

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"

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

var (
	NULL  = &Null{}
	EOF   = &Error{Message: "EOF"}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
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
	Json() string
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
func (n *Number) Json() string       { return n.Inspect() }
func (n *Number) ZeroValue() float64 { return float64(0) }
func (n *Number) Int() int           { return int(n.Value) }

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Json() string     { return b.Inspect() }

type Null struct {
	Token token.Token
}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }
func (n *Null) Json() string     { return n.Inspect() }

type ReturnValue struct {
	Token token.Token
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Json() string     { return rv.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Json() string     { return e.Inspect() }

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

func (f *Function) Json() string { return f.Inspect() }

// The String is a special fella.
//
// Like ints, or bools, you might
// think it will only have a Value
// property, the string itself.
//
// TA-DA! No, we also have an Ok and Done
// properties that are used when running
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
//
// cmd = $(sleep 10 &)
// type(cmd) // STRING
// cmd.done // FALSE
// cmd.wait() // ...
// cmd.done // TRUE
type String struct {
	Token  token.Token
	Value  string
	Ok     *Boolean  // A special property to check whether a command exited correctly
	Cmd    *exec.Cmd // A special property to access the underlying command
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
	Done   *Boolean
	mux    *sync.Mutex
}

func (s *String) Type() ObjectType  { return STRING_OBJ }
func (s *String) Inspect() string   { return s.Value }
func (s *String) Json() string      { return `"` + s.Inspect() + `"` }
func (s *String) ZeroValue() string { return "" }
func (s *String) HashKey() HashKey {
	return HashKey{Type: s.Type(), Value: s.Value}
}

// Function that ensure a mutex
// instance is created on the
// string
func (s *String) mustHaveMutex() {
	if s.mux == nil {
		s.mux = &sync.Mutex{}
	}
}

// To be called when the command
// is done. Releases the internal
// mutex.
func (s *String) SetDone() {
	s.mustHaveMutex()
	s.mux.Unlock()
}

// To be called when the command
// is starting in background, so
// that anyone accessing it will
// be blocked.
func (s *String) SetRunning() {
	s.mustHaveMutex()
	s.mux.Lock()
}

// To be called when we want to
// wait on the background command
// to be done.
func (s *String) Wait() {
	s.mustHaveMutex()
	s.mux.Lock()
	s.mux.Unlock()
}

// To be called when we want to
// kill the background command
func (s *String) Kill() error {
	err := s.Cmd.Process.Kill()

	// The command value includes output and possible error
	// We might want to change this
	output := s.Stdout.String()
	outputErr := s.Stderr.String()
	s.Value = strings.TrimSpace(output) + strings.TrimSpace(outputErr)

	if err != nil {
		return err
	}

	s.Done = TRUE
	return nil
}

// Sets the result of the underlying command
// on the string.
// 3 things are set:
// - the string itself (output of the command)
// - str.ok
// - str.done
func (s *String) SetCmdResult(Ok *Boolean) {
	s.Ok = Ok
	var output string

	if Ok.Value {
		output = s.Stdout.String()
	} else {
		output = s.Stderr.String()
	}

	// trim space at both ends of out.String(); works in both linux and windows
	s.Value = strings.TrimSpace(output)
	s.Done = TRUE
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
func (b *Builtin) Json() string     { return b.Inspect() }

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
		elements = append(elements, e.Json())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (ao *Array) Json() string { return ao.Inspect() }

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
		pairs = append(pairs, fmt.Sprintf(`%s: %s`, pair.Key.Json(), pair.Value.Json()))
	}
	// create stable key ordered output
	sort.Strings(pairs)

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
func (h *Hash) Json() string { return h.Inspect() }

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
