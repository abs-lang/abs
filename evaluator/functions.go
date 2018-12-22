package evaluator

import (
	"abs/ast"
	"abs/lexer"
	"abs/object"
	"abs/parser"
	"abs/util"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

// Utility function that validates arguments passed to builtin
// functions.
func validateArgs(name string, args []object.Object, size int, types [][]string) object.Object {
	if len(args) != size {
		return util.NewError("wrong number of arguments to %s(...): got=%d, want=1", name, len(args))
	}

	for i, t := range types {
		if !util.Contains(t, string(args[i].Type())) {
			return util.NewError("argument %d to %s(...) is not supported (got: %s, allowed: %s)", i, name, args[i].Inspect(), strings.Join(t, ", "))
		}
	}

	return nil
}

func getFns() map[string]*object.Builtin {
	return map[string]*object.Builtin{
		// len(var:"hello")
		"len": &object.Builtin{
			Types: []string{object.STRING_OBJ, object.INTEGER_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("len", args, 1, [][]string{{object.STRING_OBJ, object.ARRAY_OBJ}})
				if err != nil {
					return err
				}

				switch arg := args[0].(type) {
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				default:
					return util.NewError("argument to `len` not supported, got %s",
						args[0].Type())
				}
			},
		},
		// rand(max:20)
		"rand": &object.Builtin{
			Types: []string{object.INTEGER_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("rand", args, 1, [][]string{{object.INTEGER_OBJ}})
				if err != nil {
					return err
				}

				arg := args[0].(*object.Integer)
				r, e := rand.Int(rand.Reader, big.NewInt(arg.Value))

				if e != nil {
					return util.NewError("error occurred while calling 'rand(%d)': %s", arg.Value, e.Error())
				}

				return &object.Integer{Value: r.Int64()}
			},
		},
		// exit(code:0)
		"exit": &object.Builtin{
			Types: []string{object.INTEGER_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("exit", args, 1, [][]string{{object.INTEGER_OBJ}})
				if err != nil {
					return err
				}

				arg := args[0].(*object.Integer)
				os.Exit(int(arg.Value))
				return arg
			},
		},
		// echo(arg:"hello")
		"echo": &object.Builtin{
			Types: []string{},
			Fn: func(args ...object.Object) object.Object {
				var arguments []interface{} = make([]interface{}, len(args)-1)
				for i, d := range args {
					if i > 0 {
						arguments[i-1] = d.Inspect()
					}
				}

				fmt.Printf(args[0].Inspect(), arguments...)
				fmt.Println("")

				return NULL
			},
		},
		// int(string:"123")
		"int": &object.Builtin{
			Types: []string{object.STRING_OBJ, object.INTEGER_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("int", args, 1, [][]string{{object.INTEGER_OBJ, object.STRING_OBJ}})
				if err != nil {
					return err
				}

				switch arg := args[0].(type) {
				case *object.Integer:
					return &object.Integer{Value: int64(arg.Value)}
				case *object.String:
					i, err := strconv.Atoi(arg.Value)

					if err != nil {
						return util.NewError("int(...) can only be called on strings which represent integers, '%s' given", arg.Value)
					}

					return &object.Integer{Value: int64(i)}
				default:
					// we will never reach here
					return util.NewError("argument to `int` not supported, got %s", args[0].Type())
				}
			},
		},
		// env(variable:"PWD")
		"env": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("env", args, 1, [][]string{{object.STRING_OBJ}})
				if err != nil {
					return err
				}

				arg := args[0].(*object.String)
				return &object.String{Value: os.Getenv(arg.Value)}
			},
		},
		// args(position:1)
		"args": &object.Builtin{
			Types: []string{object.INTEGER_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("args", args, 1, [][]string{{object.INTEGER_OBJ}})
				if err != nil {
					return err
				}

				arg := args[0].(*object.Integer)
				i := arg.Value

				if int(i) > len(os.Args)-1 {
					return &object.String{Value: ""}
				}

				return &object.String{Value: os.Args[i]}
			},
		},
		// type(variable:"hello")
		"type": &object.Builtin{
			Types: []string{},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("args", args, 1, [][]string{})
				if err != nil {
					return err
				}

				return &object.String{Value: string(args[0].Type())}
			},
		},
		// split(string:"hello")
		"split": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("args", args, 2, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}})
				if err != nil {
					return err
				}

				s := args[0].(*object.String)
				sep := args[1].(*object.String)

				parts := strings.Split(s.Value, sep.Value)
				length := len(parts)
				elements := make([]object.Object, length, length)

				for k, v := range parts {
					elements[k] = &object.String{Value: v}
				}

				return &object.Array{Elements: elements}
			},
		},
		// lines(string:"a\nb")
		"lines": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("args", args, 1, [][]string{{object.STRING_OBJ}})
				if err != nil {
					return err
				}

				s := args[0].(*object.String)
				parts := strings.Split(s.Value, "\n")
				length := len(parts)
				elements := make([]object.Object, length, length)

				for k, v := range parts {
					elements[k] = &object.String{Value: v}
				}

				return &object.Array{Elements: elements}
			},
		},
		// cmd = $(ls -la)
		// cmd.ok()
		"ok": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("args", args, 1, [][]string{{object.STRING_OBJ}})
				if err != nil {
					return err
				}

				s := args[0].(*object.String)

				return &object.Boolean{Value: s.Ok}
			},
		},
		// "{}".json()
		//
		// Converts a valid JSON document to an ABS hash.
		//
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
		//
		// This method is incomplete as it currently does not
		// support most JSON types, but rather just objects,
		// ie. "[1, 2, 3]".json() won't work.
		"json": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("args", args, 1, [][]string{{object.STRING_OBJ}})
				if err != nil {
					return err
				}

				s := args[0].(*object.String)
				env := object.NewEnvironment()
				l := lexer.New(s.Value)
				p := parser.New(l)
				hl, ok := p.ParseHashLiteral().(*ast.HashLiteral)

				if ok {
					return evalHashLiteral(hl, env)
				}

				return util.NewError("argument to `json` must be a valid JSON object, got '%s'", s.Value)
			},
		},
		// sum(array:[1, 2, 3])
		"sum": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("args", args, 1, [][]string{{object.ARRAY_OBJ}})
				if err != nil {
					return err
				}

				arr := args[0].(*object.Array)

				var sum int64 = 0

				for _, v := range arr.Elements {
					elem := v.(*object.Integer)
					sum += elem.Value
				}

				return &object.Integer{Value: int64(sum)}
			},
		},
		// map(array:[1, 2, 3], function:f(x) { x + 1 })
		"map": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("args", args, 2, [][]string{{object.ARRAY_OBJ}, {object.FUNCTION_OBJ, object.BUILTIN_OBJ}})
				if err != nil {
					return err
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				newElements := make([]object.Object, length, length)
				copy(newElements, arr.Elements)

				for k, v := range arr.Elements {
					newElements[k] = applyFunction(args[1], []object.Object{v})
				}

				return &object.Array{Elements: newElements}
			},
		},
		// contains("str", "tr")
		"contains": &object.Builtin{
			Types: []string{object.ARRAY_OBJ, object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("args", args, 2, [][]string{{object.STRING_OBJ, object.ARRAY_OBJ}, {object.STRING_OBJ}})
				if err != nil {
					return err
				}

				switch args[0].(type) {
				case *object.String:
					return &object.Boolean{Value: strings.Contains(args[0].(*object.String).Value, args[1].(*object.String).Value)}
				default:
					return &object.Boolean{Value: false}
				}
			},
		},
		// str(1)
		"str": &object.Builtin{
			Types: []string{object.INTEGER_OBJ, object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("args", args, 1, [][]string{{object.STRING_OBJ, object.INTEGER_OBJ}})
				if err != nil {
					return err
				}

				switch arg := args[0].(type) {
				case *object.String:
					return &object.String{Value: arg.Value}
				case *object.Integer:
					return &object.String{Value: strconv.Itoa(int(arg.Value))}
				default:
					return NULL
				}
			},
		},
		// Compare
		// Contains
		// ContainsAny
		// ContainsRune
		// Count
		// EqualFold
		// Fields
		// FieldsFunc
		// HasPrefix
		// HasSuffix
		// Index
		// IndexAny
		// IndexByte
		// IndexFunc
		// IndexRune
		// Join
		// LastIndex
		// LastIndexAny
		// LastIndexByte
		// LastIndexFunc
		// Map
		// Repeat
		// Replace
		// Split
		// SplitAfter
		// SplitAfterN
		// SplitN
		// Title
		// ToLower
		// ToLowerSpecial
		// ToTitle
		// ToTitleSpecial
		// ToUpper
		// ToUpperSpecial
		// Trim
		// TrimFunc
		// TrimLeft
		// TrimLeftFunc
		// TrimPrefix
		// TrimRight
		// TrimRightFunc
		// TrimSpace
		// TrimSuffix
	}

}
