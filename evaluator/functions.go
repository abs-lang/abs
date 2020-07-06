package evaluator

import (
	"bufio"
	"crypto/rand"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/abs-lang/abs/ast"
	"github.com/abs-lang/abs/lexer"
	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/parser"
	"github.com/abs-lang/abs/token"
	"github.com/abs-lang/abs/util"
	"github.com/iancoleman/strcase"
)

var scanner *bufio.Scanner
var tok token.Token
var scannerPosition int

func init() {
	scanner = bufio.NewScanner(os.Stdin)
}

/*
Here be the hairy map to all the Builtin Functions ... ARRRGH, matey
*/

func getFns() map[string]*object.Builtin {
	return map[string]*object.Builtin{
		// len(var:"hello")
		"len": &object.Builtin{
			Types: []string{object.STRING_OBJ, object.ARRAY_OBJ},
			Fn:    lenFn,
		},
		// rand(max:20)
		"rand": &object.Builtin{
			Types: []string{object.NUMBER_OBJ},
			Fn:    randFn,
		},
		// exit(code:0)
		"exit": &object.Builtin{
			Types: []string{object.NUMBER_OBJ},
			Fn:    exitFn,
		},
		// flag("my-flag")
		"flag": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    flagFn,
		},
		// pwd()
		"pwd": &object.Builtin{
			Types: []string{},
			Fn:    pwdFn,
		},
		// camel("string")
		"camel": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    camelFn,
		},
		// snake("string")
		"snake": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    snakeFn,
		},
		// kebab("string")
		"kebab": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    kebabFn,
		},
		// cd() or cd(path)
		"cd": &object.Builtin{
			Types: []string{},
			Fn:    cdFn,
		},
		// clamp(num, min, max)
		"clamp": &object.Builtin{
			Types: []string{object.NUMBER_OBJ},
			Fn:    clampFn,
		},
		// echo(arg:"hello")
		"echo": &object.Builtin{
			Types: []string{},
			Fn:    echoFn,
		},
		// int(string:"123")
		// int(number:"123")
		"int": &object.Builtin{
			Types: []string{object.STRING_OBJ, object.NUMBER_OBJ},
			Fn:    intFn,
		},
		// round(string:"123.1")
		// round(number:"123.1", 2)
		"round": &object.Builtin{
			Types: []string{object.STRING_OBJ, object.NUMBER_OBJ},
			Fn:    roundFn,
		},
		// floor(string:"123.1")
		// floor(number:123.1)
		"floor": &object.Builtin{
			Types: []string{object.STRING_OBJ, object.NUMBER_OBJ},
			Fn:    floorFn,
		},
		// ceil(string:"123.1")
		// ceil(number:123.1)
		"ceil": &object.Builtin{
			Types: []string{object.STRING_OBJ, object.NUMBER_OBJ},
			Fn:    ceilFn,
		},
		// number(string:"1.23456")
		"number": &object.Builtin{
			Types: []string{object.STRING_OBJ, object.NUMBER_OBJ},
			Fn:    numberFn,
		},
		// is_number(string:"1.23456")
		"is_number": &object.Builtin{
			Types: []string{object.STRING_OBJ, object.NUMBER_OBJ},
			Fn:    isNumberFn,
		},
		// stdin()
		"stdin": &object.Builtin{
			Next:  stdinNextFn,
			Types: []string{},
			Fn:    stdinFn,
		},
		// env(variable:"PWD") or env(string:"KEY", string:"VAL")
		"env": &object.Builtin{
			Types: []string{},
			Fn:    envFn,
		},
		// arg(position:1)
		"arg": &object.Builtin{
			Types: []string{object.NUMBER_OBJ},
			Fn:    argFn,
		},
		// args()
		"args": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    argsFn,
		},
		// type(variable:"hello")
		"type": &object.Builtin{
			Types: []string{},
			Fn:    typeFn,
		},
		// fn.call(args_array)
		"call": &object.Builtin{
			Types: []string{object.FUNCTION_OBJ, object.BUILTIN_OBJ},
			Fn:    callFn,
		},
		// chnk([...], int:2)
		"chunk": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    chunkFn,
		},
		// split(string:"hello")
		"split": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    splitFn,
		},
		// lines(string:"a\nb")
		"lines": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    linesFn,
		},
		// "{}".json()
		// Converts a valid JSON document to an ABS hash.
		"json": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    jsonFn,
		},
		// "a %s".fmt(b)
		"fmt": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    fmtFn,
		},
		// sum(array:[1, 2, 3])
		"sum": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    sumFn,
		},
		// max(array:[1, 2, 3])
		"max": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    maxFn,
		},
		// min(array:[1, 2, 3])
		"min": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    minFn,
		},
		// reduce(array:[1, 2, 3], f(){}, accumulator)
		"reduce": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    reduceFn,
		},
		// sort(array:[1, 2, 3])
		"sort": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    sortFn,
		},
		// intersect(array:[1, 2, 3], array:[1, 2, 3])
		"intersect": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    intersectFn,
		},
		// diff(array:[1, 2, 3], array:[1, 2, 3])
		"diff": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    diffFn,
		},
		// union(array:[1, 2, 3], array:[1, 2, 3])
		"union": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    unionFn,
		},
		// diff_symmetric(array:[1, 2, 3], array:[1, 2, 3])
		"diff_symmetric": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    diffSymmetricFn,
		},
		// flatten(array:[1, 2, 3])
		"flatten": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    flattenFn,
		},
		// flatten(array:[1, 2, 3])
		"flatten_deep": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    flattenDeepFn,
		},
		// partition(array:[1, 2, 3])
		"partition": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    partitionFn,
		},
		// map(array:[1, 2, 3], function:f(x) { x + 1 })
		"map": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    mapFn,
		},
		// some(array:[1, 2, 3], function:f(x) { x == 2 })
		"some": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    someFn,
		},
		// every(array:[1, 2, 3], function:f(x) { x == 2 })
		"every": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    everyFn,
		},
		// find(array:[1, 2, 3], function:f(x) { x == 2 })
		"find": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    findFn,
		},
		// filter(array:[1, 2, 3], function:f(x) { x == 2 })
		"filter": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    filterFn,
		},
		// unique(array:[1, 2, 3])
		"unique": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    uniqueFn,
		},
		// str(1)
		"str": &object.Builtin{
			Types: []string{},
			Fn:    strFn,
		},
		// any("abc", "b")
		"any": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    anyFn,
		},
		// between(number, min, max)
		"between": &object.Builtin{
			Types: []string{object.NUMBER_OBJ},
			Fn:    betweenFn,
		},
		// prefix("abc", "a")
		"prefix": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    prefixFn,
		},
		// suffix("abc", "a")
		"suffix": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    suffixFn,
		},
		// repeat("abc", 3)
		"repeat": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    repeatFn,
		},
		// replace("abc", "b", "f", -1)
		"replace": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    replaceFn,
		},
		// title("some thing")
		"title": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    titleFn,
		},
		// lower("ABC")
		"lower": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    lowerFn,
		},
		// upper("abc")
		"upper": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    upperFn,
		},
		// wait(`sleep 1 &`)
		"wait": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    waitFn,
		},
		"kill": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    killFn,
		},
		// trim("abc")
		"trim": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    trimFn,
		},
		// trim_by("abc", "c")
		"trim_by": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    trimByFn,
		},
		// index("abc", "c")
		"index": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    indexFn,
		},
		// last_index("abcc", "c")
		"last_index": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    lastIndexFn,
		},
		// shift([1,2,3])
		"shift": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    shiftFn,
		},
		// reverse([1,2,3])
		"reverse": &object.Builtin{
			Types: []string{object.ARRAY_OBJ, object.STRING_OBJ},
			Fn:    reverseFn,
		},
		// push([1,2,3], 4)
		"push": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    pushFn,
		},
		// pop([1,2,3], 4)
		"pop": &object.Builtin{
			Types: []string{object.ARRAY_OBJ, object.HASH_OBJ},
			Fn:    popFn,
		},
		// keys([1,2,3]) returns array of indices
		// keys({"a": 1, "b": 2, "c": 3}) returns array of keys
		"keys": &object.Builtin{
			Types: []string{object.ARRAY_OBJ, object.HASH_OBJ},
			Fn:    keysFn,
		},
		// values({"a": 1, "b": 2, "c": 3}) returns array of values
		"values": &object.Builtin{
			Types: []string{object.HASH_OBJ},
			Fn:    valuesFn,
		},
		// items({"a": 1, "b": 2, "c": 3}) returns array of [key, value] tuples: [[a, 1], [b, 2] [c, 3]]
		"items": &object.Builtin{
			Types: []string{object.HASH_OBJ},
			Fn:    itemsFn,
		},
		// join([1,2,3], "-")
		"join": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    joinFn,
		},
		// sleep(3000)
		"sleep": &object.Builtin{
			Types: []string{object.NUMBER_OBJ},
			Fn:    sleepFn,
		},
		// source("file.abs") -- soure a file, with access to the global environment
		"source": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    sourceFn,
		},
		// require("file.abs") -- require a file without giving it access to the global environment
		"require": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    requireFn,
		},
		// exec(command) -- execute command with interactive stdIO
		"exec": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    execFn,
		},
		// eval(code) -- evaluates code in the context of the current ABS environment
		"eval": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn:    evalFn,
		},
		// tsv([[1,2,3,4], [5,6,7,8]]) -- converts an array into a TSV string
		"tsv": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn:    tsvFn,
		},
		// unix_ms() -- returns the current unix epoch, in milliseconds
		"unix_ms": &object.Builtin{
			Types: []string{},
			Fn:    unixMsFn,
		},
	}
}

/*
Here be the actual Builtin Functions
*/
// Utility function that validates arguments passed to builtin functions.
func validateArgs(tok token.Token, name string, args []object.Object, size int, types [][]string) object.Object {
	if len(args) == 0 || len(args) > size || len(args) < size {
		return newError(tok, "wrong number of arguments to %s(...): got=%d, want=%d", name, len(args), size)
	}

	for i, t := range types {
		if !util.Contains(t, string(args[i].Type())) && !util.Contains(t, object.ANY_OBJ) {
			return newError(tok, "argument %d to %s(...) is not supported (got: %s, allowed: %s)", i, name, args[i].Inspect(), strings.Join(t, ", "))
		}
	}

	return nil
}

// spec is an array of {
//   {															// signature: func(num|str, arr)
//     { NUMBER_OBJ, STRING_OBJ },	// type options for arg 0
//     { ARRAY_OBJ},								// type options for arg 1
//   },
//   {															// signature: func(num|str)
//     { NUMBER_OBJ, STRING_OBJ },	// type options for arg 0
//   },
// }
func validateVarArgs(tok token.Token, name string, args []object.Object, specs [][][]string) (object.Object, int) {
	required := -1
	max := 0

	for _, spec := range specs {
		// find the min number of arguments required
		if required == -1 || len(spec) < required {
			required = len(spec)
		}

		// find the max number of arguments supported
		if len(spec) > max {
			max = len(spec)
		}
	}

	if len(args) < required || len(args) > max {
		return newError(tok, "wrong number of arguments to %s(...): got=%d, min=%d, max=%d", name, len(args), required, max), -1
	}

	for which, spec := range specs {
		// does the number of args match this spec?
		if len(args) != len(spec) {
			continue
		}

		// do the caller's args match this spec?
		match := true
		for i, types := range spec {
			if i < len(args) && !util.Contains(types, string(args[i].Type())) {
				match = false
				break
			}
		}

		// found a match; return the index of the matched spec
		if match {
			return nil, which
		}
	}

	// no signature specs matched
	return newError(tok, usageVarArgs(name, specs)), -1
}

func usageVarArgs(name string, specs [][][]string) string {
	signatures := []string{"Wrong arguments passed to '" + name + "'. Usage:"}

	for _, spec := range specs {
		args := []string{}

		for _, types := range spec {
			args = append(args, strings.Join(types, " | "))
		}

		signatures = append(signatures, fmt.Sprintf("%s(%s)", name, strings.Join(args, ", ")))
	}

	return strings.Join(signatures, "\n")
}

// len(var:"hello")
func lenFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "len", args, 1, [][]string{{object.STRING_OBJ, object.ARRAY_OBJ}})
	if err != nil {
		return err
	}

	switch arg := args[0].(type) {
	case *object.Array:
		return &object.Number{Token: tok, Value: float64(len(arg.Elements))}
	case *object.String:
		return &object.Number{Token: tok, Value: float64(len(arg.Value))}
	default:
		return newError(tok, "argument to `len` not supported, got %s", args[0].Type())
	}
}

// rand(max:20)
func randFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "rand", args, 1, [][]string{{object.NUMBER_OBJ}})
	if err != nil {
		return err
	}

	arg := args[0].(*object.Number)
	r, e := rand.Int(rand.Reader, big.NewInt(int64(arg.Value)))

	if e != nil {
		return newError(tok, "error occurred while calling 'rand(%v)': %s", arg.Value, e.Error())
	}

	return &object.Number{Token: tok, Value: float64(r.Int64())}
}

// exit(code:0)
// exit(code:0, message:"Adios!")
func exitFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	var err object.Object
	var message string

	if len(args) == 2 {
		err = validateArgs(tok, "exit", args, 2, [][]string{{object.NUMBER_OBJ}, {object.STRING_OBJ}})
		message = args[1].(*object.String).Value
	} else {
		err = validateArgs(tok, "exit", args, 1, [][]string{{object.NUMBER_OBJ}})
	}

	if err != nil {
		return err
	}

	if message != "" {
		fmt.Fprintf(env.Writer, message)
	}

	arg := args[0].(*object.Number)
	os.Exit(int(arg.Value))
	return arg
}

// unix_ms()
func unixMsFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	return &object.Number{Value: float64(time.Now().UnixNano() / 1000000)}
}

// flag("my-flag")
func flagFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	// TODO:
	// This seems a bit more complicated than it should,
	// and I could probably use some unit testing for this.
	// In any case it's a small function so YOLO

	err := validateArgs(tok, "flag", args, 1, [][]string{{object.STRING_OBJ}})
	if err != nil {
		return err
	}

	// flag we're trying to retrieve
	name := args[0].(*object.String)
	found := false

	// Let's loop through all the arguments
	// passed to the script
	// This is O(n) but again, performance
	// is not a big deal in ABS
	for _, v := range os.Args {
		// If the flag was found in the previous
		// argument...
		if found {
			// ...and the next one is another flag
			// means we're done parsing
			// eg. --flag1 --flag2
			if strings.HasPrefix(v, "-") {
				break
			}

			// else return the next argument
			// eg --flag1 something --flag2
			return &object.String{Token: tok, Value: v}
		}

		// try to parse the flag as key=value
		parts := strings.SplitN(v, "=", 2)
		// let's just take the left-side of the flag
		left := parts[0]

		// if the left side of the current argument corresponds
		// to the flag we're looking for (both in the form of "--flag" and "-flag")...
		// ..BINGO!
		if (len(left) > 1 && left[1:] == name.Value) || (len(left) > 2 && left[2:] == name.Value) {
			if len(parts) > 1 {
				return &object.String{Token: tok, Value: parts[1]}
			} else {
				found = true
			}
		}
	}

	// If the flag was found but we got here
	// it means no value was assigned to it,
	// so let's default to true
	if found {
		return object.TRUE
	}

	// else a flag that's not found is NULL
	return NULL
}

// pwd()
func pwdFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	dir, err := os.Getwd()
	if err != nil {
		return newError(tok, err.Error())
	}
	return &object.String{Token: tok, Value: dir}
}

// camel("some string")
func camelFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	return applyStringCase("camel", strcase.ToLowerCamel, tok, env, args...)
}

// snake("some string")
func snakeFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	return applyStringCase("snake", strcase.ToSnake, tok, env, args...)
}

// kebab("some string")
func kebabFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	return applyStringCase("kebab", strcase.ToKebab, tok, env, args...)
}

func applyStringCase(fnName string, fn func(string) string, tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, fnName, args, 1, [][]string{{object.STRING_OBJ}})
	if err != nil {
		return err
	}

	return &object.String{Token: tok, Value: fn(args[0].(*object.String).Value)}
}

// cd() or cd(path) returns expanded path and path.ok
func cdFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	user, ok := user.Current()
	if ok != nil {
		return newError(tok, ok.Error())
	}
	// Default: cd to user's homeDir
	path := user.HomeDir
	if len(args) == 1 {
		// arg: rawPath
		pathStr := args[0].(*object.String)
		rawPath := pathStr.Value
		path, _ = util.ExpandPath(rawPath)
	}
	// NB. windows os.Chdir(path) will convert any '/' in path to '\', however linux will not
	error := os.Chdir(path)
	if error != nil {
		// path does not exist, return error string and !path.ok
		return &object.String{Token: tok, Value: error.Error(), Ok: &object.Boolean{Token: tok, Value: false}}
	}
	// return the full path we cd()'d into and path.ok
	// this will also test true/false for cd("path/to/somewhere") && `ls`
	dir, _ := os.Getwd()
	return &object.String{Token: tok, Value: dir, Ok: &object.Boolean{Token: tok, Value: true}}
}

// clamp(n, min, max)
func clampFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "clamp", args, 3, [][]string{{object.NUMBER_OBJ}, {object.NUMBER_OBJ}, {object.NUMBER_OBJ}})
	if err != nil {
		return err
	}

	n := args[0].(*object.Number)
	min := args[1].(*object.Number)
	max := args[2].(*object.Number)

	if min.Value >= max.Value {
		return newError(tok, "arguments to clamp(min, max) must satisfy min < max (%s < %s given)", min.Inspect(), max.Inspect())
	}

	val := n.Value

	if min.Value > n.Value {
		val = min.Value
	}

	if max.Value < n.Value {
		val = max.Value
	}

	return &object.Number{Value: val}
}

// echo(arg:"hello")
func echoFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	if len(args) == 0 {
		// allow echo() without crashing
		fmt.Fprintln(env.Writer, "")
		return NULL
	}
	var arguments []interface{} = make([]interface{}, len(args)-1)
	for i, d := range args {
		if i > 0 {
			arguments[i-1] = d.Inspect()
		}
	}

	fmt.Fprintf(env.Writer, args[0].Inspect(), arguments...)
	fmt.Fprintln(env.Writer, "")

	return NULL
}

// int(string:"123")
// int(number:123)
func intFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "int", args, 1, [][]string{{object.NUMBER_OBJ, object.STRING_OBJ}})
	if err != nil {
		return err
	}

	return applyMathFunction(tok, args[0], func(n float64) float64 {
		return float64(int64(n))
	}, "int")
}

// round(string:"123.1")
// round(number:123.1)
func roundFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	// Validate first argument
	err := validateArgs(tok, "round", args[:1], 1, [][]string{{object.NUMBER_OBJ, object.STRING_OBJ}})
	if err != nil {
		return err
	}

	decimal := float64(1)

	// If we have a second argument, let's validate it
	if len(args) > 1 {
		err := validateArgs(tok, "round", args[1:], 1, [][]string{{object.NUMBER_OBJ}})
		if err != nil {
			return err
		}

		decimal = float64(math.Pow(10, args[1].(*object.Number).Value))
	}

	return applyMathFunction(tok, args[0], func(n float64) float64 {
		return math.Round(n*decimal) / decimal
	}, "round")
}

// floor(string:"123.1")
// floor(number:123.1)
func floorFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "floor", args, 1, [][]string{{object.NUMBER_OBJ, object.STRING_OBJ}})
	if err != nil {
		return err
	}

	return applyMathFunction(tok, args[0], math.Floor, "floor")
}

// ceil(string:"123.1")
// ceil(number:123.1)
func ceilFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "ceil", args, 1, [][]string{{object.NUMBER_OBJ, object.STRING_OBJ}})
	if err != nil {
		return err
	}

	return applyMathFunction(tok, args[0], math.Ceil, "ceil")
}

// Base function to do math operations. This is here
// so that we abstract away some of the common logic
// between all math functions, for example:
// - allowing to be called on strings as well ("1.23".ceil())
// - handling errors
// NB. callers must pass the token that is used for error line reporting
func applyMathFunction(tok token.Token, arg object.Object, fn func(float64) float64, fname string) object.Object {
	switch arg := arg.(type) {
	case *object.Number:
		return &object.Number{Token: tok, Value: float64(fn(arg.Value))}
	case *object.String:
		i, err := strconv.ParseFloat(arg.Value, 64)

		if err != nil {
			return newError(tok, "%s(...) can only be called on strings which represent numbers, '%s' given", fname, arg.Value)
		}

		return &object.Number{Token: tok, Value: float64(fn(i))}
	default:
		// we should never reach here since our callers should validate
		// the type of the arguments
		return newError(tok, "argument to `%s` not supported, got %s", fname, arg.Type())
	}
}

// number(string:"1.23456")
func numberFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "number", args, 1, [][]string{{object.NUMBER_OBJ, object.STRING_OBJ}})
	if err != nil {
		return err
	}

	switch arg := args[0].(type) {
	case *object.Number:
		return arg
	case *object.String:
		i, err := strconv.ParseFloat(arg.Value, 64)

		if err != nil {
			return newError(tok, "number(...) can only be called on strings which represent numbers, '%s' given", arg.Value)
		}

		return &object.Number{Token: tok, Value: i}
	default:
		// we will never reach here
		return newError(tok, "argument to `number` not supported, got %s", args[0].Type())
	}
}

// is_number(string:"1.23456")
func isNumberFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "number", args, 1, [][]string{{object.NUMBER_OBJ, object.STRING_OBJ}})
	if err != nil {
		return err
	}

	switch arg := args[0].(type) {
	case *object.Number:
		return &object.Boolean{Token: tok, Value: true}
	case *object.String:
		return &object.Boolean{Token: tok, Value: util.IsNumber(arg.Value)}
	default:
		// we will never reach here
		return newError(tok, "argument to `is_number` not supported, got %s", args[0].Type())
	}
}

// stdin() -- implemented with 2 functions
func stdinFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	v := scanner.Scan()

	if !v {
		return EOF
	}

	return &object.String{Token: tok, Value: scanner.Text()}
}
func stdinNextFn() (object.Object, object.Object) {
	v := scanner.Scan()

	if !v {
		return nil, EOF
	}

	defer func() {
		scannerPosition += 1
	}()
	return &object.Number{Value: float64(scannerPosition)}, &object.String{Token: tok, Value: scanner.Text()}
}

// env(variable:"PWD") or env(string:"KEY", string:"VAL")
func envFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err, spec := validateVarArgs(tok, "env", args, [][][]string{
		{{object.STRING_OBJ}, {object.STRING_OBJ}},
		{{object.STRING_OBJ}},
	})

	if err != nil {
		return err
	}

	key := args[0].(*object.String)

	if spec == 0 {
		val := args[1].(*object.String)
		os.Setenv(key.Value, val.Value)
	}

	return &object.String{Token: tok, Value: os.Getenv(key.Value)}
}

// arg(position:1)
func argFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "arg", args, 1, [][]string{{object.NUMBER_OBJ}})
	if err != nil {
		return err
	}

	arg := args[0].(*object.Number)
	i := arg.Int()

	if i > len(os.Args)-1 || i < 0 {
		return &object.String{Token: tok, Value: ""}
	}

	return &object.String{Token: tok, Value: os.Args[i]}
}

// args()
func argsFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	length := len(os.Args)
	result := make([]object.Object, length, length)

	for i, v := range os.Args {
		result[i] = &object.String{Token: tok, Value: v}
	}

	return &object.Array{Elements: result}
}

// type(variable:"hello")
func typeFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "type", args, 1, [][]string{})
	if err != nil {
		return err
	}

	return &object.String{Token: tok, Value: string(args[0].Type())}
}

// fn.call(args_array)
func callFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "call", args, 2, [][]string{{object.FUNCTION_OBJ, object.BUILTIN_OBJ}, {object.ARRAY_OBJ}})
	if err != nil {
		return err
	}

	return applyFunction(tok, args[0], env, args[1].(*object.Array).Elements)
}

// chunk([...], integer:2)
func chunkFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "chunk", args, 2, [][]string{{object.ARRAY_OBJ}, {object.NUMBER_OBJ}})
	if err != nil {
		return err
	}

	number := args[1].(*object.Number)
	size := int(number.Value)

	if size < 1 || !number.IsInt() {
		return newError(tok, "argument to chunk must be a positive integer, got '%s'", number.Inspect())
	}

	var chunks []object.Object
	elements := args[0].(*object.Array).Elements

	for i := 0; i < len(elements); i += size {
		end := i + size

		if end > len(elements) {
			end = len(elements)
		}

		chunks = append(chunks, &object.Array{Elements: elements[i:end]})
	}

	return &object.Array{Elements: chunks}
}

// split(string:"hello world!", sep:" ")
func splitFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err, spec := validateVarArgs(tok, "split", args, [][][]string{
		{{object.STRING_OBJ}, {object.STRING_OBJ}},
		{{object.STRING_OBJ}},
	})

	if err != nil {
		return err
	}

	s := args[0].(*object.String)

	sep := " "
	if spec == 0 {
		sep = args[1].(*object.String).Value
	}

	parts := strings.Split(s.Value, sep)
	length := len(parts)
	elements := make([]object.Object, length, length)

	for k, v := range parts {
		elements[k] = &object.String{Token: tok, Value: v}
	}

	return &object.Array{Elements: elements}
}

// lines(string:"a\nb")
func linesFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "lines", args, 1, [][]string{{object.STRING_OBJ}})
	if err != nil {
		return err
	}

	s := args[0].(*object.String)
	parts := strings.FieldsFunc(s.Value, func(r rune) bool {
		return r == '\n' || r == '\r' || r == '\f'
	})
	length := len(parts)
	elements := make([]object.Object, length, length)

	for k, v := range parts {
		elements[k] = &object.String{Token: tok, Value: v}
	}

	return &object.Array{Elements: elements}
}

// "{}".json()
// Converts a valid JSON document to an ABS hash.
func jsonFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	// One interesting thing here is that we're creating
	// a new environment from scratch, whereas it might
	// be interesting to use the existing one. That would
	// allow to do things like:
	//
	// x = 10
	// '{"key": x}'.json()["key"] // 10
	//
	// Also, we're instantiating a new lexer & parser from
	// scratch, so this is a tad slow.

	err := validateArgs(tok, "json", args, 1, [][]string{{object.STRING_OBJ}})
	if err != nil {
		return err
	}

	s := args[0].(*object.String)
	str := strings.TrimSpace(s.Value)
	env = object.NewEnvironment(env.Writer, env.Dir, env.Version)
	l := lexer.New(str)
	p := parser.New(l)
	var node ast.Node
	ok := false

	// JSON types:
	// - objects
	// - arrays
	// - number
	// - string
	// - null
	// - bool
	if len(str) != 0 {
		switch str[0] {
		case '{':
			node, ok = p.ParseHashLiteral().(*ast.HashLiteral)
		case '[':
			node, ok = p.ParseArrayLiteral().(*ast.ArrayLiteral)
		}
	}

	// if str is empty, the length will be 0
	// we can parse it the same way as string literal
	if len(str) == 0 || (str[0] == '"' && str[len(str)-1] == '"') {
		node, ok = p.ParseStringLiteral().(*ast.StringLiteral)
	}

	if util.IsNumber(str) {
		node, ok = p.ParseNumberLiteral().(*ast.NumberLiteral)
	}

	if str == "false" || str == "true" {
		node, ok = p.ParseBoolean().(*ast.Boolean)
	}

	if str == "null" {
		return NULL
	}

	if ok {
		return Eval(node, env)
	}

	return newError(tok, "argument to `json` must be a valid JSON object, got '%s'", s.Value)

}

// "a %s".fmt(b)
func fmtFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	list := []interface{}{}

	for _, s := range args[1:] {
		list = append(list, s.Inspect())
	}

	return &object.String{Token: tok, Value: fmt.Sprintf(args[0].(*object.String).Value, list...)}
}

// sum(array:[1, 2, 3])
func sumFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "sum", args, 1, [][]string{{object.ARRAY_OBJ}})
	if err != nil {
		return err
	}

	arr := args[0].(*object.Array)
	if arr.Empty() {
		return &object.Number{Token: tok, Value: float64(0)}
	}

	if !arr.Homogeneous() {
		return newError(tok, "sum(...) can only be called on an homogeneous array, got %s", arr.Inspect())
	}

	if arr.Elements[0].Type() != object.NUMBER_OBJ {
		return newError(tok, "sum(...) can only be called on arrays of numbers, got %s", arr.Inspect())
	}

	var sum float64 = 0

	for _, v := range arr.Elements {
		elem := v.(*object.Number)
		sum += elem.Value
	}

	return &object.Number{Token: tok, Value: sum}
}

// max(array:[1, 2, 3])
func maxFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "max", args, 1, [][]string{{object.ARRAY_OBJ}})
	if err != nil {
		return err
	}

	arr := args[0].(*object.Array)
	if arr.Empty() {
		return object.NULL
	}

	if !arr.Homogeneous() {
		return newError(tok, "max(...) can only be called on an homogeneous array, got %s", arr.Inspect())
	}

	if arr.Elements[0].Type() != object.NUMBER_OBJ {
		return newError(tok, "max(...) can only be called on arrays of numbers, got %s", arr.Inspect())
	}

	max := arr.Elements[0].(*object.Number).Value

	for _, v := range arr.Elements[1:] {
		elem := v.(*object.Number)

		if elem.Value > max {
			max = elem.Value
		}
	}

	return &object.Number{Token: tok, Value: max}
}

// min(array:[1, 2, 3])
func minFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "min", args, 1, [][]string{{object.ARRAY_OBJ}})
	if err != nil {
		return err
	}

	arr := args[0].(*object.Array)
	if arr.Empty() {
		return object.NULL
	}

	if !arr.Homogeneous() {
		return newError(tok, "min(...) can only be called on an homogeneous array, got %s", arr.Inspect())
	}

	if arr.Elements[0].Type() != object.NUMBER_OBJ {
		return newError(tok, "min(...) can only be called on arrays of numbers, got %s", arr.Inspect())
	}

	min := arr.Elements[0].(*object.Number).Value

	for _, v := range arr.Elements[1:] {
		elem := v.(*object.Number)

		if elem.Value < min {
			min = elem.Value
		}
	}

	return &object.Number{Token: tok, Value: min}
}

// reduce(array:[1, 2, 3], f(){}, accumulator)
func reduceFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "reduce", args, 3, [][]string{{object.ARRAY_OBJ}, {object.FUNCTION_OBJ}, {object.ANY_OBJ}})
	if err != nil {
		return err
	}

	accumulator := args[2]

	for _, v := range args[0].(*object.Array).Elements {
		accumulator = applyFunction(tok, args[1].(*object.Function), env, []object.Object{accumulator, v})
	}

	return accumulator
}

// sort(array:[1, 2, 3])
func sortFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "sort", args, 1, [][]string{{object.ARRAY_OBJ}})
	if err != nil {
		return err
	}

	arr := args[0].(*object.Array)
	elements := arr.Elements

	if len(elements) == 0 {
		return arr
	}

	if !arr.Homogeneous() {
		return newError(tok, "argument to 'sort' must be an homogeneous array (elements of the same type), got %s", arr.Inspect())
	}

	switch elements[0].(type) {
	case *object.Number:
		a := []float64{}
		for _, v := range elements {
			a = append(a, v.(*object.Number).Value)
		}
		sort.Float64s(a)

		o := []object.Object{}

		for _, v := range a {
			o = append(o, &object.Number{Token: tok, Value: v})
		}
		return &object.Array{Elements: o}
	case *object.String:
		a := []string{}
		for _, v := range elements {
			a = append(a, v.(*object.String).Value)
		}
		sort.Strings(a)

		o := []object.Object{}

		for _, v := range a {
			o = append(o, &object.String{Token: tok, Value: v})
		}
		return &object.Array{Elements: o}
	default:
		return newError(tok, "cannot sort an array with given elements elements (%s)", arr.Inspect())
	}
}

// intersect(array:[1, 2, 3], array:[1, 2, 3])
func intersectFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "intersect", args, 2, [][]string{{object.ARRAY_OBJ}, {object.ARRAY_OBJ}})
	if err != nil {
		return err
	}

	left := args[0].(*object.Array).Elements
	right := args[1].(*object.Array).Elements
	found := map[string]object.Object{}
	intersection := []object.Object{}

	for _, o := range right {
		found[object.GenerateEqualityString(o)] = o
	}

	for _, o := range left {
		element, ok := found[object.GenerateEqualityString(o)]

		if ok {
			intersection = append(intersection, element)
		}
	}

	return &object.Array{Elements: intersection}
}

// diff(array:[1, 2, 3], array:[1, 2, 3])
func diff(symmetric bool, fnName string, tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, fnName, args, 2, [][]string{{object.ARRAY_OBJ}, {object.ARRAY_OBJ}})
	if err != nil {
		return err
	}

	left := args[0].(*object.Array).Elements
	right := args[1].(*object.Array).Elements
	foundRight := map[string]object.Object{}
	difference := []object.Object{}

	for _, o := range right {
		foundRight[object.GenerateEqualityString(o)] = o
	}

	for _, o := range left {
		_, ok := foundRight[object.GenerateEqualityString(o)]

		if !ok {
			difference = append(difference, o)
		}
	}

	if symmetric {
		// If the did is symmetric, we simply re-run this function with the arrays swapped
		// so diff_sym(a, b) = diff(a, b) + diff(b, a)
		difference = append(difference, diff(false, fnName, tok, env, args[1], args[0]).(*object.Array).Elements...)
	}

	return &object.Array{Elements: difference}
}

func diffFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	return diff(false, "diff", tok, env, args...)
}

// diff_symmetric(array:[1, 2, 3], array:[1, 2, 3])
func diffSymmetricFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	return diff(true, "diff_symmetric", tok, env, args...)
}

// union(array:[1, 2, 3], array:[1, 2, 3])
func unionFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "union", args, 2, [][]string{{object.ARRAY_OBJ}, {object.ARRAY_OBJ}})
	if err != nil {
		return err
	}

	left := args[0].(*object.Array).Elements
	right := args[1].(*object.Array).Elements

	union := []object.Object{}

	for _, v := range left {
		union = append(union, v)
	}

	m := util.Mapify(left)

	for _, v := range right {
		_, found := m[object.GenerateEqualityString(v)]

		if !found {
			union = append(union, v)
		}
	}

	return &object.Array{Elements: union}
}

// flatten(array:[1, 2, 3])
func flattenFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	return flatten("flatten", false, tok, env, args...)
}

// flatten_deep(array:[1, 2, 3])
func flattenDeepFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	return flatten("flatten_deep", true, tok, env, args...)
}

func flatten(fnName string, deep bool, tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, fnName, args, 1, [][]string{{object.ARRAY_OBJ}})
	if err != nil {
		return err
	}

	originalElements := args[0].(*object.Array).Elements
	elements := []object.Object{}

	for _, v := range originalElements {
		switch e := v.(type) {
		case *object.Array:
			if deep {
				elements = append(elements, flattenDeepFn(tok, env, e).(*object.Array).Elements...)
			} else {
				for _, x := range e.Elements {
					elements = append(elements, x)
				}
			}
		default:
			elements = append(elements, e)
		}
	}

	return &object.Array{Elements: elements}
}

func partitionFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "partition", args, 2, [][]string{{object.ARRAY_OBJ}, {object.FUNCTION_OBJ, object.BUILTIN_OBJ}})
	if err != nil {
		return err
	}

	partitions := map[string][]object.Object{}
	elements := args[0].(*object.Array).Elements
	// This will allows us to preserve the order
	// of partitions based on the order of elements.
	//
	// When we run the partitioning function, we store
	// it's results in a map of result{list_of_values...}.
	// When we loop over that map, Go doesn't guarantee
	// order of results (https://nathanleclaire.com/blog/2014/04/27/a-surprising-feature-of-golang-that-colored-me-impressed/)
	// but we want to, so
	// we use the partitionOrder list to extract values
	// from the map based on the order they were
	// inserted in.
	partitionOrder := []string{}
	scanned := map[string]bool{}

	for _, v := range elements {
		res := applyFunction(tok, args[1], env, []object.Object{v})
		eqs := object.GenerateEqualityString(res)

		partitions[eqs] = append(partitions[eqs], v)

		if _, ok := scanned[eqs]; !ok {
			partitionOrder = append(partitionOrder, eqs)
			scanned[eqs] = true
		}
	}

	result := &object.Array{Elements: []object.Object{}}
	for _, eqs := range partitionOrder {
		partition := partitions[eqs]
		result.Elements = append(result.Elements, &object.Array{Elements: partition})
	}

	return result
}

// map(array:[1, 2, 3], function:f(x) { x + 1 })
func mapFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "map", args, 2, [][]string{{object.ARRAY_OBJ}, {object.FUNCTION_OBJ, object.BUILTIN_OBJ}})
	if err != nil {
		return err
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	newElements := make([]object.Object, length, length)
	copy(newElements, arr.Elements)

	for k, v := range arr.Elements {
		evaluated := applyFunction(tok, args[1], env, []object.Object{v})

		if isError(evaluated) {
			return evaluated
		}
		newElements[k] = evaluated
	}

	return &object.Array{Elements: newElements}
}

// some(array:[1, 2, 3], function:f(x) { x == 2 })
func someFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "some", args, 2, [][]string{{object.ARRAY_OBJ}, {object.FUNCTION_OBJ, object.BUILTIN_OBJ}})
	if err != nil {
		return err
	}

	var result bool

	arr := args[0].(*object.Array)

	for _, v := range arr.Elements {
		r := applyFunction(tok, args[1], env, []object.Object{v})

		if isTruthy(r) {
			result = true
			break
		}
	}

	return &object.Boolean{Token: tok, Value: result}
}

// every(array:[1, 2, 3], function:f(x) { x == 2 })
func everyFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "every", args, 2, [][]string{{object.ARRAY_OBJ}, {object.FUNCTION_OBJ, object.BUILTIN_OBJ}})
	if err != nil {
		return err
	}

	result := true

	arr := args[0].(*object.Array)

	for _, v := range arr.Elements {
		r := applyFunction(tok, args[1], env, []object.Object{v})

		if !isTruthy(r) {
			result = false
		}
	}

	return &object.Boolean{Token: tok, Value: result}
}

// find(array:[1, 2, 3], function:f(x) { x == 2 })
func findFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "find", args, 2, [][]string{{object.ARRAY_OBJ}, {object.FUNCTION_OBJ, object.BUILTIN_OBJ, object.HASH_OBJ}})
	if err != nil {
		return err
	}

	arr := args[0].(*object.Array)

	switch predicate := args[1].(type) {
	case *object.Hash:
		for _, v := range arr.Elements {
			v, ok := v.(*object.Hash)

			if !ok {
				continue
			}

			match := true
			for k, pair := range predicate.Pairs {
				toCompare, ok := v.GetPair(k.Value)
				if !ok {
					match = false
					continue
				}

				if !object.Equal(pair.Value, toCompare.Value) {
					match = false
				}
			}

			if match {
				return v
			}
		}
	default:
		for _, v := range arr.Elements {
			r := applyFunction(tok, predicate, env, []object.Object{v})

			if isTruthy(r) {
				return v
			}
		}
	}

	return NULL
}

// filter(array:[1, 2, 3], function:f(x) { x == 2 })
func filterFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "filter", args, 2, [][]string{{object.ARRAY_OBJ}, {object.FUNCTION_OBJ, object.BUILTIN_OBJ}})
	if err != nil {
		return err
	}

	result := []object.Object{}
	arr := args[0].(*object.Array)

	for _, v := range arr.Elements {
		evaluated := applyFunction(tok, args[1], env, []object.Object{v})

		if isError(evaluated) {
			return evaluated
		}

		if isTruthy(evaluated) {
			result = append(result, v)
		}
	}

	return &object.Array{Elements: result}
}

// unique(array:[1, 2, 3])
func uniqueFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "unique", args, 1, [][]string{{object.ARRAY_OBJ}})
	if err != nil {
		return err
	}

	result := []object.Object{}
	arr := args[0].(*object.Array)
	existingElements := map[string]bool{}

	for _, v := range arr.Elements {
		key := object.GenerateEqualityString(v)

		if _, ok := existingElements[key]; !ok {
			existingElements[key] = true
			result = append(result, v)
		}
	}

	return &object.Array{Elements: result}
}

// str(1)
func strFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "str", args, 1, [][]string{})
	if err != nil {
		return err
	}

	return &object.String{Token: tok, Value: args[0].Inspect()}
}

// any("abc", "b")
func anyFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "any", args, 2, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}})
	if err != nil {
		return err
	}

	return &object.Boolean{Token: tok, Value: strings.ContainsAny(args[0].(*object.String).Value, args[1].(*object.String).Value)}
}

// between(10, 0, 100)
func betweenFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "between", args, 3, [][]string{{object.NUMBER_OBJ}, {object.NUMBER_OBJ}, {object.NUMBER_OBJ}})
	if err != nil {
		return err
	}

	n := args[0].(*object.Number)
	min := args[1].(*object.Number)
	max := args[2].(*object.Number)

	if min.Value >= max.Value {
		return newError(tok, "arguments to between(min, max) must satisfy min < max (%s < %s given)", min.Inspect(), max.Inspect())
	}

	return &object.Boolean{Token: tok, Value: ((min.Value <= n.Value) && (n.Value <= max.Value))}
}

// prefix("abc", "a")
func prefixFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "prefix", args, 2, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}})
	if err != nil {
		return err
	}

	return &object.Boolean{Token: tok, Value: strings.HasPrefix(args[0].(*object.String).Value, args[1].(*object.String).Value)}
}

// suffix("abc", "a")
func suffixFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "suffix", args, 2, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}})
	if err != nil {
		return err
	}

	return &object.Boolean{Token: tok, Value: strings.HasSuffix(args[0].(*object.String).Value, args[1].(*object.String).Value)}
}

// repeat("abc", 3)
func repeatFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "repeat", args, 2, [][]string{{object.STRING_OBJ}, {object.NUMBER_OBJ}})
	if err != nil {
		return err
	}

	return &object.String{Token: tok, Value: strings.Repeat(args[0].(*object.String).Value, int(args[1].(*object.Number).Value))}
}

// replace("abd", "d", "c") --> short form
// replace("abd", "d", "c", -1)
// replace("abc", ["a", "b"], "c", -1)
func replaceFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	var err object.Object

	// Support short form
	if len(args) == 3 {
		err = validateArgs(tok, "replace", args, 3, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ, object.ARRAY_OBJ}, {object.STRING_OBJ}})
	} else {
		err = validateArgs(tok, "replace", args, 4, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ, object.ARRAY_OBJ}, {object.STRING_OBJ}, {object.NUMBER_OBJ}})
	}

	if err != nil {
		return err
	}

	original := args[0].(*object.String).Value
	replacement := args[2].(*object.String).Value

	n := -1

	if len(args) == 4 {
		n = int(args[3].(*object.Number).Value)
	}

	if characters, ok := args[1].(*object.Array); ok {
		for _, c := range characters.Elements {
			original = strings.Replace(original, c.Inspect(), replacement, n)
		}

		return &object.String{Token: tok, Value: original}
	}

	return &object.String{Token: tok, Value: strings.Replace(original, args[1].(*object.String).Value, replacement, n)}
}

// title("some thing")
func titleFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "title", args, 1, [][]string{{object.STRING_OBJ}})
	if err != nil {
		return err
	}

	return &object.String{Token: tok, Value: strings.Title(args[0].(*object.String).Value)}
}

// lower("ABC")
func lowerFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "lower", args, 1, [][]string{{object.STRING_OBJ}})
	if err != nil {
		return err
	}

	return &object.String{Token: tok, Value: strings.ToLower(args[0].(*object.String).Value)}
}

// upper("abc")
func upperFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "upper", args, 1, [][]string{{object.STRING_OBJ}})
	if err != nil {
		return err
	}

	return &object.String{Token: tok, Value: strings.ToUpper(args[0].(*object.String).Value)}
}

// wait(`sleep 10 &`)
func waitFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "wait", args, 1, [][]string{{object.STRING_OBJ}})
	if err != nil {
		return err
	}

	cmd := args[0].(*object.String)

	if cmd.Cmd == nil {
		return cmd
	}

	cmd.Wait()
	return cmd
}

// kill(`sleep 10 &`)
func killFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "kill", args, 1, [][]string{{object.STRING_OBJ}})
	if err != nil {
		return err
	}

	cmd := args[0].(*object.String)

	if cmd.Cmd == nil {
		return cmd
	}

	errCmdKill := cmd.Kill()

	if errCmdKill != nil {
		return newError(tok, "Error killing command %s with error %s", cmd.Inspect(), errCmdKill.Error())
	}
	return cmd
}

// trim("abc")
func trimFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "trim", args, 1, [][]string{{object.STRING_OBJ}})
	if err != nil {
		return err
	}

	return &object.String{Token: tok, Value: strings.Trim(args[0].(*object.String).Value, " ")}
}

// trim_by("abc", "c")
func trimByFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "trim_by", args, 2, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}})
	if err != nil {
		return err
	}

	return &object.String{Token: tok, Value: strings.Trim(args[0].(*object.String).Value, args[1].(*object.String).Value)}
}

// index("abc", "c")
func indexFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "index", args, 2, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}})
	if err != nil {
		return err
	}

	i := strings.Index(args[0].(*object.String).Value, args[1].(*object.String).Value)

	if i == -1 {
		return NULL
	}

	return &object.Number{Token: tok, Value: float64(i)}
}

// last_index("abcc", "c")
func lastIndexFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "last_index", args, 2, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}})
	if err != nil {
		return err
	}

	i := strings.LastIndex(args[0].(*object.String).Value, args[1].(*object.String).Value)

	if i == -1 {
		return NULL
	}

	return &object.Number{Token: tok, Value: float64(i)}
}

// Clamps start and end arguments to the slice
// function. When you slice "abc" you can have
// start 10 and end -20...
func sliceStartAndEnd(l int, start int, end int) (int, int) {
	if end == 0 {
		end = l
	}

	if start > l {
		start = l
	}

	if start < 0 {
		newStart := l + start
		if newStart < 0 {
			start = 0
		} else {
			start = newStart
		}
	}

	if end > l || start > end {
		end = l
	}

	return start, end
}

// shift([1,2,3]) removes and returns first value or null if array is empty
func shiftFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "shift", args, 1, [][]string{{object.ARRAY_OBJ}})
	if err != nil {
		return err
	}

	array := args[0].(*object.Array)
	if len(array.Elements) == 0 {
		return NULL
	}
	e := array.Elements[0]
	array.Elements = append(array.Elements[:0], array.Elements[1:]...)

	return e
}

// reverse([1,2,3])
func reverseFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err, spec := validateVarArgs(tok, "reverse", args, [][][]string{
		{{object.ARRAY_OBJ}},
		{{object.STRING_OBJ}},
	})

	if err != nil {
		return err
	}

	if spec == 0 {
		// array
		array := args[0].(*object.Array)

		for i, j := 0, len(array.Elements)-1; i < j; i, j = i+1, j-1 {
			array.Elements[i], array.Elements[j] = array.Elements[j], array.Elements[i]
		}

		return array
	} else {
		// string
		str := []rune(args[0].(*object.String).Value)

		for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
			str[i], str[j] = str[j], str[i]
		}

		return &object.String{Token: tok, Value: string(str)}
	}
}

// push([1,2,3], 4)
func pushFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "push", args, 2, [][]string{{object.ARRAY_OBJ}, {object.NULL_OBJ,
		object.ARRAY_OBJ, object.NUMBER_OBJ, object.STRING_OBJ, object.HASH_OBJ}})
	if err != nil {
		return err
	}

	array := args[0].(*object.Array)
	array.Elements = append(array.Elements, args[1])

	return array
}

// pop([1,2,3]) removes and returns last value or null if array is empty
// pop({"a":1, "b":2, "c":3}, "a") removes and returns {"key": value} or null if key not found
func popFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	// pop has 2 signatures: pop(array), and pop(hash, key)
	var err object.Object
	if len(args) > 0 {
		if args[0].Type() == object.ARRAY_OBJ {
			err = validateArgs(tok, "pop", args, 1, [][]string{{object.ARRAY_OBJ}})
		} else if args[0].Type() == object.HASH_OBJ {
			err = validateArgs(tok, "pop", args, 2, [][]string{{object.HASH_OBJ}})
		}
	}
	if err != nil {
		return err
	}
	if len(args) < 1 {
		return NULL
	}
	switch arg := args[0].(type) {
	case *object.Array:
		if len(arg.Elements) > 0 {
			elem := arg.Elements[len(arg.Elements)-1]
			arg.Elements = arg.Elements[0 : len(arg.Elements)-1]
			return elem
		}
	case *object.Hash:
		if len(args) == 2 {
			key := args[1].(object.Hashable)
			hashKey := key.HashKey()
			item, ok := arg.Pairs[hashKey]
			if ok {
				pairs := make(map[object.HashKey]object.HashPair)
				pairs[hashKey] = item
				delete(arg.Pairs, hashKey)
				return &object.Hash{Pairs: pairs}
			}
		}
	}
	return NULL
}

// keys([1,2,3]) returns array of indices
// keys({"a": 1, "b": 2, "c": 3}) returns array of keys
func keysFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "keys", args, 1, [][]string{{object.ARRAY_OBJ, object.HASH_OBJ}})
	if err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Array:
		length := len(arg.Elements)
		newElements := make([]object.Object, length, length)
		for k := range arg.Elements {
			newElements[k] = &object.Number{Token: tok, Value: float64(k)}
		}
		return &object.Array{Elements: newElements}
	case *object.Hash:
		pairs := arg.Pairs
		keys := []object.Object{}
		for _, pair := range pairs {
			key := pair.Key
			keys = append(keys, key)
		}
		return &object.Array{Elements: keys}
	}
	return NULL
}

// values({"a": 1, "b": 2, "c": 3}) returns array of values
func valuesFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "values", args, 1, [][]string{{object.HASH_OBJ}})
	if err != nil {
		return err
	}
	hash := args[0].(*object.Hash)
	pairs := hash.Pairs
	values := []object.Object{}
	for _, pair := range pairs {
		value := pair.Value
		values = append(values, value)
	}
	return &object.Array{Elements: values}
}

// items({"a": 1, "b": 2, "c": 3}) returns array of [key, value] tuples: [[a, 1], [b, 2] [c, 3]]
func itemsFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "items", args, 1, [][]string{{object.HASH_OBJ}})
	if err != nil {
		return err
	}
	hash := args[0].(*object.Hash)
	pairs := hash.Pairs
	items := []object.Object{}
	for _, pair := range pairs {
		key := pair.Key
		value := pair.Value
		item := &object.Array{Elements: []object.Object{key, value}}
		items = append(items, item)
	}
	return &object.Array{Elements: items}
}

func joinFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err, spec := validateVarArgs(tok, "join", args, [][][]string{
		{{object.ARRAY_OBJ}, {object.STRING_OBJ}},
		{{object.ARRAY_OBJ}},
	})

	if err != nil {
		return err
	}

	glue := ""
	if spec == 0 {
		glue = args[1].(*object.String).Value
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)
	newElements := make([]string, length, length)

	for k, v := range arr.Elements {
		newElements[k] = v.Inspect()
	}

	return &object.String{Token: tok, Value: strings.Join(newElements, glue)}
}

func sleepFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "sleep", args, 1, [][]string{{object.NUMBER_OBJ}})
	if err != nil {
		return err
	}

	ms := args[0].(*object.Number)
	time.Sleep(time.Duration(ms.Value) * time.Millisecond)

	return NULL
}

// source("file.abs")
const ABS_SOURCE_DEPTH = "10"

var sourceDepth, _ = strconv.Atoi(ABS_SOURCE_DEPTH)
var sourceLevel = 0

func sourceFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	file, _ := util.ExpandPath(args[0].Inspect())
	return doSource(tok, env, file, args...)
}

// require("file.abs")
var history = make(map[string]string)

var packageAliases map[string]string
var packageAliasesLoaded bool

func requireFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	if !packageAliasesLoaded {
		a, err := ioutil.ReadFile("./packages.abs.json")

		// We couldn't open the packages, file, possibly doesn't exists
		// and the code shouldn't fail
		if err == nil {
			// Try to decode the packages file:
			// if an error occurs we will simply
			// ignore it
			json.Unmarshal(a, &packageAliases)
		}

		packageAliasesLoaded = true
	}

	file := util.UnaliasPath(args[0].Inspect(), packageAliases)

	if !strings.HasPrefix(file, "@") {
		file = filepath.Join(env.Dir, file)
	}

	e := object.NewEnvironment(env.Writer, filepath.Dir(file), env.Version)
	return doSource(tok, e, file, args...)
}

func doSource(tok token.Token, env *object.Environment, fileName string, args ...object.Object) object.Object {
	err := validateArgs(tok, "source", args, 1, [][]string{{object.STRING_OBJ}})
	if err != nil {
		// reset the source level
		sourceLevel = 0
		return err
	}

	// get configured source depth if any
	sourceDepthStr := util.GetEnvVar(env, "ABS_SOURCE_DEPTH", ABS_SOURCE_DEPTH)
	sourceDepth, _ = strconv.Atoi(sourceDepthStr)

	// limit source file inclusion depth
	if sourceLevel >= sourceDepth {
		// reset the source level
		sourceLevel = 0
		// use errObj.Message instead of errObj.Inspect() to avoid nested "ERROR: " prefixes
		errObj := newError(tok, "maximum source file inclusion depth exceeded at %d levels", sourceDepth)
		errObj = &object.Error{Message: errObj.Message}
		return errObj
	}
	// mark this source level
	sourceLevel++

	var code []byte
	var error error

	// Manage std library requires starting with
	// a '@' eg. require('@runtime')
	if strings.HasPrefix(fileName, "@") {
		code, error = Asset("stdlib/" + fileName[1:])
	} else {
		// load the source file
		code, error = ioutil.ReadFile(fileName)
	}

	if error != nil {
		// reset the source level
		sourceLevel = 0
		// cannot read source file
		return newError(tok, "cannot read source file: %s:\n%s", fileName, error.Error())
	}
	// parse it
	l := lexer.New(string(code))
	p := parser.New(l)
	program := p.ParseProgram()
	errors := p.Errors()
	if len(errors) != 0 {
		// reset the source level
		sourceLevel = 0
		errMsg := fmt.Sprintf("%s", " parser errors:\n")
		for _, msg := range errors {
			errMsg += fmt.Sprintf("%s", "\t"+msg+"\n")
		}
		return newError(tok, "error found in source file: %s\n%s", fileName, errMsg)
	}
	// invoke BeginEval() passing in the sourced program, env, and our lexer
	// we save the current global lexer and restore it after we return from BeginEval()
	// NB. saving the lexer allows error line numbers to be relative to any nested source files
	savedLexer := lex
	evaluated := BeginEval(program, env, l)
	lex = savedLexer
	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		// use errObj.Message instead of errObj.Inspect() to avoid nested "ERROR: " prefixes
		evalErrMsg := evaluated.(*object.Error).Message
		sourceErrMsg := newError(tok, "error found in eval block: %s", fileName).Message
		errObj := &object.Error{Message: fmt.Sprintf("%s\n\t%s", sourceErrMsg, evalErrMsg)}
		return errObj
	}
	// restore this source level
	sourceLevel--

	return evaluated
}

func evalFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "eval", args, 1, [][]string{{object.STRING_OBJ}})
	if err != nil {
		return err
	}

	// parse it
	l := lexer.New(string(args[0].Inspect()))
	p := parser.New(l)
	program := p.ParseProgram()
	errors := p.Errors()
	if len(errors) != 0 {
		errMsg := fmt.Sprintf("%s", " parser errors:\n")
		for _, msg := range errors {
			errMsg += fmt.Sprintf("%s", "\t"+msg+"\n")
		}
		return newError(tok, "error found in eval block: %s\n%s", args[0].Inspect(), errMsg)
	}
	// invoke BeginEval() passing in the sourced program, env, and our lexer
	// we save the current global lexer and restore it after we return from BeginEval()
	// NB. saving the lexer allows error line numbers to be relative to any nested source files
	savedLexer := lex
	evaluated := BeginEval(program, env, l)
	lex = savedLexer

	if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
		// use errObj.Message instead of errObj.Inspect() to avoid nested "ERROR: " prefixes
		evalErrMsg := evaluated.(*object.Error).Message
		sourceErrMsg := newError(tok, "error found in eval block: %s", args[0].Inspect()).Message
		errObj := &object.Error{Message: fmt.Sprintf("%s\n\t%s", sourceErrMsg, evalErrMsg)}
		return errObj
	}

	return evaluated
}

// [[1,2], [3,4]].tsv()
// [{"a": 1, "b": 2}, {"b": 3, "c": 4}].tsv()
func tsvFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	// all arguments were passed
	if len(args) == 3 {
		err := validateArgs(tok, "tsv", args, 3, [][]string{{object.ARRAY_OBJ}, {object.STRING_OBJ}, {object.ARRAY_OBJ}})
		if err != nil {
			return err
		}
	}

	// If no header was passed, let's set it to empty list by default
	if len(args) == 2 {
		err := validateArgs(tok, "tsv", args, 2, [][]string{{object.ARRAY_OBJ}, {object.STRING_OBJ}})
		if err != nil {
			return err
		}
		args = append(args, &object.Array{Elements: []object.Object{}})
	}

	// If no separator and header was passed, let's set them to tab and empty list by default
	if len(args) == 1 {
		err := validateArgs(tok, "tsv", args, 1, [][]string{{object.ARRAY_OBJ}})
		if err != nil {
			return err
		}
		args = append(args, &object.String{Value: "\t"})
		args = append(args, &object.Array{Elements: []object.Object{}})
	}

	array := args[0].(*object.Array)
	separator := args[1].(*object.String).Value

	if len(separator) < 1 {
		return newError(tok, "the separator argument to the tsv() function needs to be a valid character, '%s' given", separator)
	}
	// the final outut
	out := &strings.Builder{}
	tsv := csv.NewWriter(out)
	tsv.Comma = rune(separator[0])

	// whether our array is made of ALL arrays or ALL hashes
	var isArray bool
	var isHash bool
	homogeneous := array.Homogeneous()

	if len(array.Elements) > 0 {
		_, isArray = array.Elements[0].(*object.Array)
		_, isHash = array.Elements[0].(*object.Hash)
	}

	// if the array is not homogeneous, we cannot process it
	if !homogeneous || (!isArray && !isHash) {
		return newError(tok, "tsv() must be called on an array of arrays or objects, such as [[1, 2, 3], [4, 5, 6]], '%s' given as argument", array.Inspect())
	}

	headerObj := args[2].(*object.Array)
	header := []string{}

	if len(headerObj.Elements) > 0 {
		for _, v := range headerObj.Elements {
			header = append(header, v.Inspect())
		}
	} else if isHash {
		// if our array is made of hashes, we will include a header in
		// our TSV output, made of all possible keys found in every object
		for _, rows := range array.Elements {
			for _, pair := range rows.(*object.Hash).Pairs {
				header = append(header, pair.Key.Inspect())
			}
		}

		// When no header is provided, we will simply
		// use the list of keys from all object, alphabetically
		// sorted
		header = util.UniqueStrings(header)
		sort.Strings(header)
	}

	if len(header) > 0 {
		err := tsv.Write(header)

		if err != nil {
			return newError(tok, err.Error())
		}
	}

	for _, row := range array.Elements {
		// Row values
		values := []string{}

		// In the case of an array, creating the row is fairly
		// straightforward: we loop through the elements and extract
		// their value
		if isArray {
			for _, element := range row.(*object.Array).Elements {
				values = append(values, element.Inspect())
			}

		}

		// In case of an hash, we want to extract values based on
		// the header. If a key is not present in an hash, we will
		// simply set it to null
		if isHash {
			for _, key := range header {
				pair, ok := row.(*object.Hash).GetPair(key)
				var value object.Object

				if ok {
					value = pair.Value
				} else {
					value = NULL
				}

				values = append(values, value.Inspect())
			}
		}

		// Add the row to the final output, by concatenating
		// it with the given separator
		err := tsv.Write(values)

		if err != nil {
			return newError(tok, err.Error())
		}
	}

	tsv.Flush()
	return &object.String{Value: strings.TrimSpace(out.String())}
}

func execFn(tok token.Token, env *object.Environment, args ...object.Object) object.Object {
	err := validateArgs(tok, "exec", args, 1, [][]string{{object.STRING_OBJ}})
	if err != nil {
		return err
	}
	cmd := args[0].Inspect()
	cmd = strings.Trim(cmd, " ")

	// interpolate any $vars in the cmd string
	cmd = util.InterpolateStringVars(cmd, env)

	var commands []string
	var executor string
	if runtime.GOOS == "windows" {
		commands = []string{"/C", cmd}
		executor = "cmd.exe"
	} else { //assume it's linux, darwin, freebsd, openbsd, solaris, etc
		// invoke bash commands with login option (-l) --
		// this allows the use of commands in $PATH
		commands = []string{"-lc", cmd}
		executor = "bash"
	}
	// set up command to execute using our stdIO
	c := exec.Command(executor, commands...)
	c.Env = os.Environ()
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	// N.B. that a bash command may end with '&' --
	// in this case bash will launch it as a daemon process and then exit c.Run() immediately
	// this may require pkill to terminate the daemon process using the pid
	runErr := c.Run()

	if runErr != nil {
		return &object.String{Value: runErr.Error()}
	}
	return NULL
}
