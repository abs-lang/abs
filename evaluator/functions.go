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

func getFns() map[string]*object.Builtin {
	return map[string]*object.Builtin{
		// len(var:"hello")
		"len": &object.Builtin{
			Types: []string{object.STRING_OBJ, object.INTEGER_OBJ},
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return util.NewError("wrong number of arguments. got=%d, want=1", len(args))
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
				if len(args) != 1 {
					return util.NewError("wrong number of arguments. got=%d, want=1", len(args))
				}

				switch arg := args[0].(type) {
				case *object.Integer:
					r, err := rand.Int(rand.Reader, big.NewInt(arg.Value))

					if err != nil {
						return util.NewError("error occurred while calling 'rand(%d)': %s", arg.Value, err.Error())
					}

					return &object.Integer{Value: r.Int64()}
				default:
					return util.NewError("argument to `rand(...)` not supported, got %s", arg.Type())
				}
			},
		},
		// exit(code:0)
		"exit": &object.Builtin{
			Types: []string{object.INTEGER_OBJ},
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return util.NewError("wrong number of arguments. got=%d, want=1", len(args))
				}

				switch arg := args[0].(type) {
				case *object.Integer:
					os.Exit(int(arg.Value))
					return arg
				default:
					return util.NewError("argument to `exit(...)` not supported, got '%s' (%s)", arg.Inspect(), arg.Type())
				}
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
				if len(args) != 1 {
					return util.NewError("wrong number of arguments. got=%d, want=1", len(args))
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
					return util.NewError("argument to `int` not supported, got %s", args[0].Type())
				}
			},
		},
		// env(variable:"PWD")
		"env": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return util.NewError("wrong number of arguments. got=%d, want=1", len(args))
				}

				switch arg := args[0].(type) {
				case *object.String:
					return &object.String{Value: os.Getenv(arg.Value)}
				default:
					return util.NewError("argument to `env` not supported, got %s", args[0].Type())
				}
			},
		},
		// args(position:1)
		"args": &object.Builtin{
			Types: []string{object.INTEGER_OBJ},
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return util.NewError("wrong number of arguments. got=%d, want=1", len(args))
				}

				switch arg := args[0].(type) {
				case *object.Integer:
					i := arg.Value

					if int(i) > len(os.Args)-1 {
						return &object.String{Value: ""}
					}

					return &object.String{Value: os.Args[i]}
				default:
					return util.NewError("argument to `args(...)` not supported, got %s", args[0].Type())
				}
			},
		},
		// type(variable:"hello")
		"type": &object.Builtin{
			Types: []string{},
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return util.NewError("wrong number of arguments. got=%d, want=1", len(args))
				}

				return &object.String{Value: string(args[0].Type())}
			},
		},
		// split(string:"hello")
		"split": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				// TODO we're passing the wrong amount of args to these functions
				if len(args) != 2 {
					return util.NewError("wrong number of arguments. got=%d, want=1", len(args))
				}

				if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
					return util.NewError("argument to `split` must be STRING, got %s", args[0].Type())
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
		// cmd = $(ls -la)
		// cmd.ok()
		"ok": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
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
				if len(args) != 1 {
					return util.NewError("wrong number of arguments. got=%d, want=1", len(args))
				}

				if args[0].Type() != object.STRING_OBJ {
					return util.NewError("argument to `split` must be STRING, got %s", args[0].Type())
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
				if len(args) != 1 {
					return util.NewError("wrong number of arguments. got=%d, want=1", len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return util.NewError("argument to `first` must be ARRAY, got %s", args[0].Type())
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
				if len(args) != 2 {
					return util.NewError("wrong number of arguments. got=%d, want=2", len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return util.NewError("argument to `map` must be ARRAY, got %s", args[0].Type())
				}
				if args[1].Type() != object.FUNCTION_OBJ && args[1].Type() != object.BUILTIN_OBJ {
					return util.NewError("argument to `map` must be FUNCTION, got %s", args[0].Type())
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
	}

}
