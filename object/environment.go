package object

import (
	"io"
	"sort"
)

// NewEnclosedEnvironment creates an environment
// with another one embedded to it, so that the
// new environment has access to identifiers stored
// in the outer one.
func NewEnclosedEnvironment(outer *Environment, args []Object) *Environment {
	env := NewEnvironment(outer.Writer, outer.Dir, outer.Version)
	env.outer = outer
	env.CurrentArgs = args
	return env
}

// NewEnvironment creates a new environment to run
// ABS in, specifying a writer for the output of the
// program and the base dir (which is used to require
// other scripts)
func NewEnvironment(w io.Writer, dir string, version string) *Environment {
	s := make(map[string]Object)
	// e.Version and ABS_VERSION are duplicate, we should
	// find a better way to manage it
	e := &Environment{store: s, outer: nil, Writer: w, Dir: dir, Version: version}
	e.Set("ABS_VERSION", &String{Value: e.Version})

	return e
}

// Environment represent the environment associated
// with the execution context of an ABS script: it
// holds all variables etc.
type Environment struct {
	store map[string]Object
	// Arguments this environment was created in.
	// When we call function(1, 2, 3), a new environment
	// for the function to execute is created, and 1/2/3
	// are recorded as arguments for this environment.
	//
	// Later, if we need to access the arguments passed
	// to the function, we can refer back to them
	// through env.CurrentArgs. This is how ... is
	// implemented.
	CurrentArgs []Object
	outer       *Environment
	// Used to capture output. This is typically os.Stdout,
	// but you could capture in any io.Writer of choice
	Writer io.Writer
	// Dir represents the directory from which we're executing code.
	// It starts as the directory from which we invoke the ABS
	// executable, but changes when we call require("...") as each
	// require call resets the dir to its own directory, so that
	// relative imports work.
	//
	// If we have script A and B in /tmp, A can require("B")
	// wihout having to specify its full absolute path
	// eg. require("/tmp/B")
	Dir string
	// Version of the ABS runtime
	Version string
}

// Get returns an identifier stored within the environment
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) GetHere(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// GetKeys returns the list of all identifiers
// stored in this environment
func (e *Environment) GetKeys() []string {
	keys := make([]string, 0, len(e.store))
	for k := range e.store {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

// Set sets an identifier in the environment
func (e *Environment) SetOuter(name string, val Object) Object {
	_, found := e.GetHere(name)

	if found {
		e.store[name] = val
	} else {
		e.outer.store[name] = val
	}

	return val
}
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val

	return val
}

// Delete deletes an identifier from the environment
func (e *Environment) Delete(name string) {
	delete(e.store, name)
}
