package object

import (
	"io"
)

// NewEnclosedEnvironment creates an environment
// with another one embedded to it, so that the
// new environment has access to identifiers stored
// in the outer one.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment(outer.Writer)
	env.outer = outer
	return env
}

// NewEnvironment creates a new environment to run
// ABS in
func NewEnvironment(w io.Writer) *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil, Writer: w}
}

// Environment represent the environment associated
// with the execution context of an ABS script: it
// holds all variables etc.
type Environment struct {
	store map[string]Object
	outer *Environment
	// Used to capture output. This is typically os.Stdout,
	// but you could capture in any io.Writer of choice
	Writer io.Writer
}

// Get returns an identifier stored within the environment
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// GetKeys returns the list of all identifiers
// stored in this environment
func (e *Environment) GetKeys() []string {
	keys := make([]string, 0, len(e.store))
	for k := range e.store {
		keys = append(keys, k)
	}

	return keys
}

// Set sets an identifier in the environment
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Delete deletes an identifier from the environment
func (e *Environment) Delete(name string) {
	delete(e.store, name)
}
