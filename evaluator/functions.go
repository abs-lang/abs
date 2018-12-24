package evaluator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/abs-lang/abs/ast"
	"github.com/abs-lang/abs/lexer"
	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/parser"
	"github.com/abs-lang/abs/util"
)

// Utility function that validates arguments passed to builtin
// functions.
func validateArgs(name string, args []object.Object, size int, types [][]string) object.Object {
	if len(args) != size {
		return newError("wrong number of arguments to %s(...): got=%d, want=%d", name, len(args), size)
	}

	for i, t := range types {
		if !util.Contains(t, string(args[i].Type())) {
			return newError("argument %d to %s(...) is not supported (got: %s, allowed: %s)", i, name, args[i].Inspect(), strings.Join(t, ", "))
		}
	}

	return nil
}

func getFns() map[string]*object.Builtin {
	return map[string]*object.Builtin{
		// len(var:"hello")
		"len": &object.Builtin{
			Types: []string{object.STRING_OBJ, object.ARRAY_OBJ},
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
					return newError("argument to `len` not supported, got %s",
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
					return newError("error occurred while calling 'rand(%d)': %s", arg.Value, e.Error())
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
						return newError("int(...) can only be called on strings which represent integers, '%s' given", arg.Value)
					}

					return &object.Integer{Value: int64(i)}
				default:
					// we will never reach here
					return newError("argument to `int` not supported, got %s", args[0].Type())
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
		// arg(position:1)
		"arg": &object.Builtin{
			Types: []string{object.INTEGER_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("arg", args, 1, [][]string{{object.INTEGER_OBJ}})
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
				err := validateArgs("type", args, 1, [][]string{})
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
				err := validateArgs("split", args, 2, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}})
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
				err := validateArgs("lines", args, 1, [][]string{{object.STRING_OBJ}})
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
				err := validateArgs("json", args, 1, [][]string{{object.STRING_OBJ}})
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

				return newError("argument to `json` must be a valid JSON object, got '%s'", s.Value)
			},
		},
		// "a %s".fmt(b)
		"fmt": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				// err := validateArgs("fmt", args, 1, [][]string{{object.STRING_OBJ}})
				// if err != nil {
				// return err
				// }

				list := []interface{}{}

				for _, s := range args[1:] {
					list = append(list, s.Inspect())
				}

				return &object.String{Value: fmt.Sprintf(args[0].(*object.String).Value, list...)}
			},
		},
		// sum(array:[1, 2, 3])
		"sum": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("sum", args, 1, [][]string{{object.ARRAY_OBJ}})
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
		// sort(array:[1, 2, 3])
		"sort": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("sort", args, 1, [][]string{{object.ARRAY_OBJ}})
				if err != nil {
					return err
				}

				arr := args[0].(*object.Array)
				elements := arr.Elements

				if len(elements) == 0 {
					return arr
				}

				if !arr.Homogeneous() {
					return newError("argument to `sort` must be an homogeneous array (elements of the same type), got %s", arr.Inspect())
				}

				switch elements[0].(type) {
				case *object.Integer:
					a := []int{}
					for _, v := range elements {
						a = append(a, int(v.(*object.Integer).Value))
					}
					sort.Ints(a)

					o := []object.Object{}

					for _, v := range a {
						o = append(o, &object.Integer{Value: int64(v)})
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
						o = append(o, &object.String{Value: v})
					}
					return &object.Array{Elements: o}
				default:
					return newError("cannot sort an array with given elements elements (%s)", arr.Inspect())
				}
			},
		},
		// map(array:[1, 2, 3], function:f(x) { x + 1 })
		"map": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("map", args, 2, [][]string{{object.ARRAY_OBJ}, {object.FUNCTION_OBJ, object.BUILTIN_OBJ}})
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
		// some(array:[1, 2, 3], function:f(x) { x == 2 })
		"some": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("some", args, 2, [][]string{{object.ARRAY_OBJ}, {object.FUNCTION_OBJ, object.BUILTIN_OBJ}})
				if err != nil {
					return err
				}

				var result bool

				arr := args[0].(*object.Array)

				for _, v := range arr.Elements {
					r := applyFunction(args[1], []object.Object{v})

					if isTruthy(r) {
						result = true
						break
					}
				}

				return &object.Boolean{Value: result}
			},
		},
		// every(array:[1, 2, 3], function:f(x) { x == 2 })
		"every": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("every", args, 2, [][]string{{object.ARRAY_OBJ}, {object.FUNCTION_OBJ, object.BUILTIN_OBJ}})
				if err != nil {
					return err
				}

				result := true

				arr := args[0].(*object.Array)

				for _, v := range arr.Elements {
					r := applyFunction(args[1], []object.Object{v})

					if !isTruthy(r) {
						result = false
					}
				}

				return &object.Boolean{Value: result}
			},
		},
		// find(array:[1, 2, 3], function:f(x) { x == 2 })
		"find": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("find", args, 2, [][]string{{object.ARRAY_OBJ}, {object.FUNCTION_OBJ, object.BUILTIN_OBJ}})
				if err != nil {
					return err
				}

				arr := args[0].(*object.Array)

				for _, v := range arr.Elements {
					r := applyFunction(args[1], []object.Object{v})

					if isTruthy(r) {
						return v
					}
				}

				return NULL
			},
		},
		// filter(array:[1, 2, 3], function:f(x) { x == 2 })
		"filter": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("filter", args, 2, [][]string{{object.ARRAY_OBJ}, {object.FUNCTION_OBJ, object.BUILTIN_OBJ}})
				if err != nil {
					return err
				}

				result := []object.Object{}
				arr := args[0].(*object.Array)

				for _, v := range arr.Elements {
					r := applyFunction(args[1], []object.Object{v})

					if isTruthy(r) {
						result = append(result, v)
					}
				}

				return &object.Array{Elements: result}
			},
		},
		// contains("str", "tr")
		"contains": &object.Builtin{
			Types: []string{object.ARRAY_OBJ, object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("contains", args, 2, [][]string{{object.STRING_OBJ, object.ARRAY_OBJ}, {object.STRING_OBJ, object.INTEGER_OBJ}})
				if err != nil {
					return err
				}

				switch arg := args[0].(type) {
				case *object.String:
					needle, ok := args[1].(*object.String)

					if ok {
						return &object.Boolean{Value: strings.Contains(arg.Value, needle.Value)}
					}
				case *object.Array:
					var found bool

					switch needle := args[1].(type) {
					case *object.String:
						for _, v := range arg.Elements {
							if v.Inspect() == needle.Value {
								found = true
								break // Let's get outta here!
							}
						}

						return &object.Boolean{Value: found}
					case *object.Integer:
						for _, v := range arg.Elements {
							// Quite ghetto but also the easiest way out
							// Instead of doing type checking on the argument,
							// we received back its string representation.
							// If they match, we then check that its type was
							// integer.
							if v.Inspect() == strconv.Itoa(int(needle.Value)) && v.Type() == object.INTEGER_OBJ {
								found = true
								break // Let's get outta here!
							}
						}

						return &object.Boolean{Value: found}
					}
				}

				return &object.Boolean{Value: false}
			},
		},
		// str(1)
		"str": &object.Builtin{
			Types: []string{},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("str", args, 1, [][]string{})
				if err != nil {
					return err
				}

				return &object.String{Value: args[0].Inspect()}
			},
		},
		// any("abc", "b")
		"any": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("any", args, 2, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}})
				if err != nil {
					return err
				}

				return &object.Boolean{Value: strings.ContainsAny(args[0].(*object.String).Value, args[1].(*object.String).Value)}
			},
		},
		// prefix("abc", "a")
		"prefix": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("prefix", args, 2, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}})
				if err != nil {
					return err
				}

				return &object.Boolean{Value: strings.HasPrefix(args[0].(*object.String).Value, args[1].(*object.String).Value)}
			},
		},
		// suffix("abc", "a")
		"suffix": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("suffix", args, 2, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}})
				if err != nil {
					return err
				}

				return &object.Boolean{Value: strings.HasSuffix(args[0].(*object.String).Value, args[1].(*object.String).Value)}
			},
		},
		// repeat("abc", 3)
		"repeat": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("repeat", args, 2, [][]string{{object.STRING_OBJ}, {object.INTEGER_OBJ}})
				if err != nil {
					return err
				}

				return &object.String{Value: strings.Repeat(args[0].(*object.String).Value, int(args[1].(*object.Integer).Value))}
			},
		},
		// replace("abc", "b", "f", -1)
		"replace": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("replace", args, 4, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}, {object.STRING_OBJ}, {object.INTEGER_OBJ}})
				if err != nil {
					return err
				}

				return &object.String{Value: strings.Replace(args[0].(*object.String).Value, args[1].(*object.String).Value, args[2].(*object.String).Value, int(args[3].(*object.Integer).Value))}
			},
		},
		// title("some thing")
		"title": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("title", args, 1, [][]string{{object.STRING_OBJ}})
				if err != nil {
					return err
				}

				return &object.String{Value: strings.Title(args[0].(*object.String).Value)}
			},
		},
		// lower("ABC")
		"lower": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("lower", args, 1, [][]string{{object.STRING_OBJ}})
				if err != nil {
					return err
				}

				return &object.String{Value: strings.ToLower(args[0].(*object.String).Value)}
			},
		},
		// upper("abc")
		"upper": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("upper", args, 1, [][]string{{object.STRING_OBJ}})
				if err != nil {
					return err
				}

				return &object.String{Value: strings.ToUpper(args[0].(*object.String).Value)}
			},
		},
		// trim("abc")
		"trim": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("trim", args, 1, [][]string{{object.STRING_OBJ}})
				if err != nil {
					return err
				}

				return &object.String{Value: strings.Trim(args[0].(*object.String).Value, " ")}
			},
		},
		// trim_by("abc", "c")
		"trim_by": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("trim_by", args, 2, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}})
				if err != nil {
					return err
				}

				return &object.String{Value: strings.Trim(args[0].(*object.String).Value, args[1].(*object.String).Value)}
			},
		},
		// index("abc", "c")
		"index": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("index", args, 2, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}})
				if err != nil {
					return err
				}

				i := strings.Index(args[0].(*object.String).Value, args[1].(*object.String).Value)

				if i == -1 {
					return NULL
				}

				return &object.Integer{Value: int64(i)}
			},
		},
		// last_index("abcc", "c")
		"last_index": &object.Builtin{
			Types: []string{object.STRING_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("last_index", args, 2, [][]string{{object.STRING_OBJ}, {object.STRING_OBJ}})
				if err != nil {
					return err
				}

				i := strings.LastIndex(args[0].(*object.String).Value, args[1].(*object.String).Value)

				if i == -1 {
					return NULL
				}

				return &object.Integer{Value: int64(i)}
			},
		},
		// slice("abcc", 0, -1)
		"slice": &object.Builtin{
			Types: []string{object.STRING_OBJ, object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("slice", args, 3, [][]string{{object.STRING_OBJ, object.ARRAY_OBJ}, {object.INTEGER_OBJ}, {object.INTEGER_OBJ}})
				if err != nil {
					return err
				}

				start := int(args[1].(*object.Integer).Value)
				end := int(args[2].(*object.Integer).Value)

				switch arg := args[0].(type) {
				case *object.String:
					s := arg.Value
					start, end := sliceStartAndEnd(len(s), start, end)

					return &object.String{Value: s[start:end]}
				case *object.Array:
					start, end := sliceStartAndEnd(len(arg.Elements), start, end)

					return &object.Array{Elements: arg.Elements[start:end]}
				}

				return NULL
			},
		},
		// shift([1,2,3])
		"shift": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("shift", args, 1, [][]string{{object.ARRAY_OBJ}})
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
			},
		},
		// reverse([1,2,3])
		"reverse": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("reverse", args, 1, [][]string{{object.ARRAY_OBJ}})
				if err != nil {
					return err
				}

				array := args[0].(*object.Array)

				for i, j := 0, len(array.Elements)-1; i < j; i, j = i+1, j-1 {
					array.Elements[i], array.Elements[j] = array.Elements[j], array.Elements[i]
				}

				return array
			},
		},
		// push([1,2,3], 4)
		"push": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("push", args, 2, [][]string{{object.ARRAY_OBJ}, {object.NULL_OBJ, object.ARRAY_OBJ, object.INTEGER_OBJ, object.STRING_OBJ, object.HASH_OBJ}})
				if err != nil {
					return err
				}

				array := args[0].(*object.Array)
				array.Elements = append(array.Elements, args[1])

				return array
			},
		},
		// pop([1,2,3], 4)
		"pop": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("pop", args, 1, [][]string{{object.ARRAY_OBJ}})
				if err != nil {
					return err
				}

				array := args[0].(*object.Array)
				elem := array.Elements[len(array.Elements)-1]
				array.Elements = array.Elements[0 : len(array.Elements)-1]

				return elem
			},
		},
		// keys([1,2,3])
		"keys": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("keys", args, 1, [][]string{{object.ARRAY_OBJ}})
				if err != nil {
					return err
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				newElements := make([]object.Object, length, length)

				for k, _ := range arr.Elements {
					newElements[k] = &object.Integer{Value: int64(k)}
				}

				return &object.Array{Elements: newElements}
			},
		},
		// join([1,2,3], "-")
		"join": &object.Builtin{
			Types: []string{object.ARRAY_OBJ},
			Fn: func(args ...object.Object) object.Object {
				err := validateArgs("join", args, 2, [][]string{{object.ARRAY_OBJ}, {object.STRING_OBJ}})
				if err != nil {
					return err
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				newElements := make([]string, length, length)

				for k, v := range arr.Elements {
					newElements[k] = v.Inspect()
				}

				return &object.String{Value: strings.Join(newElements, args[1].(*object.String).Value)}
			},
		},
	}
}

// Clamps starts and end arguments to the slice
// function. When you slice "abc" you can have
// start=10 and end -20...
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
