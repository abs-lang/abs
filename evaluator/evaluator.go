package evaluator

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/abs-lang/abs/ast"
	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/util"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	Fns   map[string]*object.Builtin
)

func init() {
	Fns = getFns()
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.ForInExpression:
		return evalForInExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.MethodExpression:
		o := Eval(node.Object, env)
		if isError(o) {
			return o
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyMethod(o, node.Method.String(), args)

	case *ast.PropertyExpression:
		return evalPropertyExpression(node, env)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.CommandExpression:
		return evalCommandExpression(node.Value, env)

	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(
	block *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	// 1 && 2
	// We will first verify left is truthy and,
	// if so, proceed to check whether right is
	// also truthy.
	// At the end of the process we will return
	// right, without any implicit bool conversion.
	case operator == "&&":
		if !isTruthy(left) {
			return left
		}

		return right
	// 1 || 2
	// We will first verify left is truthy, and
	// return it if so. If not, we will return
	// right, without any implicit bool conversion
	// (which allows short-circuiting).
	case operator == "||":
		if isTruthy(left) {
			return left
		}

		return right
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.ARRAY_OBJ && right.Type() == object.ARRAY_OBJ:
		return evalArrayInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		switch o := right.(type) {
		case *object.String:
			if o.Value == o.ZeroValue() {
				return TRUE
			}

			return FALSE
		case *object.Integer:
			if o.Value == o.ZeroValue() {
				return TRUE
			}

			return FALSE
		default:
			return FALSE
		}
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "**":
		// TODO this does not support floats
		return &object.Integer{Value: int64(math.Pow(float64(leftVal), float64(rightVal)))}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	// A range results in an array of integers from left to right
	case "..":
		a := make([]object.Object, 0)

		for i := leftVal; i <= rightVal; i++ {
			a = append(a, &object.Integer{Value: int64(i)})
		}
		return &object.Array{Elements: a}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	if operator == "+" {
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return &object.String{Value: leftVal + rightVal}
	}

	if operator == "==" {
		return &object.Boolean{Value: left.(*object.String).Value == right.(*object.String).Value}
	}

	if operator == "~" {
		return &object.Boolean{Value: strings.ToLower(left.(*object.String).Value) == strings.ToLower(right.(*object.String).Value)}
	}

	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalArrayInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	if operator == "+" {
		leftVal := left.(*object.Array).Elements
		rightVal := right.(*object.Array).Elements
		return &object.Array{Elements: append(leftVal, rightVal...)}
	}

	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalIfExpression(
	ie *ast.IfExpression,
	env *object.Environment,
) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalWhileExpression(
	we *ast.WhileExpression,
	env *object.Environment,
) object.Object {
	condition := Eval(we.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		Eval(we.Consequence, env)
		evalWhileExpression(we, env)
	}
	return NULL
}

// for x = 0; x < 10; x++ {x}
func evalForExpression(
	fe *ast.ForExpression,
	env *object.Environment,
) object.Object {
	// Let's figure out if the foor loop is using a variable that's
	// already been declared. If so, let's keep it aside for now.
	existingIdentifier, identifierExisted := env.Get(fe.Identifier)

	// Eval the starter (x = 0)
	err := Eval(fe.Starter, env)
	if isError(err) {
		return err
	}

	// This represents whether the for condition holds true
	holds := true

	// Final cleanup: we remove the x from the environment. If
	// it was already declared before the foor loop, we restore
	// it to its original value
	defer func() {
		if identifierExisted {
			env.Set(fe.Identifier, existingIdentifier)
		} else {
			env.Delete(fe.Identifier)
		}
	}()

	// When for is while...
	for holds {
		// Evaluate the for condition
		evaluated := Eval(fe.Condition, env)
		if isError(evaluated) {
			return evaluated
		}

		// If truthy, execute the block and the closer
		if isTruthy(evaluated) {
			err = Eval(fe.Block, env)
			if isError(err) {
				return err
			}

			err = Eval(fe.Closer, env)
			if isError(err) {
				return err
			}

			continue
		}

		// If not, let's break out of the loop
		holds = false
	}

	return NULL
}

// for k,v in 1..10 {v}
func evalForInExpression(
	fie *ast.ForInExpression,
	env *object.Environment,
) object.Object {
	iterable := Eval(fie.Iterable, env)
	switch i := iterable.(type) {
	case *object.Array:
		// If "k" and "v" were already declared, let's keep
		// them aside...
		existingKeyIdentifier, okk := env.Get(fie.Key)
		existingValueIdentifier, okv := env.Get(fie.Value)

		// ...so that we can restore them after the for
		// loop is over
		defer func() {
			if okk {
				env.Set(fie.Key, existingKeyIdentifier)
			} else {
				env.Delete(fie.Key)
			}

			if okv {
				env.Set(fie.Value, existingValueIdentifier)
			} else {
				env.Delete(fie.Value)
			}
		}()

		// Iterate through the array
		for k, v := range i.Elements {
			// set the special k v variables in the
			// environment
			env.Set(fie.Key, &object.Integer{Value: int64(k)})
			env.Set(fie.Value, v)
			err := Eval(fie.Block, env)

			if isError(err) {
				return err
			}
		}

		return NULL
	default:
		return newError("'%s' is not an array, cannot be used in for loop", i.Inspect())
	}
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := Fns[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

// This is the core of ABS's logical
// evaluation, and epic quirks we'll
// remember for years are to be found
// here.
func isTruthy(obj object.Object) bool {
	switch v := obj.(type) {
	// A null is always false
	case *object.Null:
		return false
	case *object.Boolean:
		return v.Value
	// An integer is truthy
	// unless it's 0
	case *object.Integer:
		return v.Value != v.ZeroValue()
	// A string is truthy
	// unless is empty
	case *object.String:
		return v.Value != v.ZeroValue()
	// Everything else is truthy
	//
	// NOTE: we might regret this
	// in the future
	//
	// NOTE 2: yolo!
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// Property expression (x.y) evaluator.
//
// Here we have a special case, as strings
// have an .ok property when they're the result
// of a command.
//
// Else we will try to parse the property
// as an index of an hash.
//
// If that doesn't work, we'll spectacularly
// give up.
func evalPropertyExpression(pe *ast.PropertyExpression, env *object.Environment) object.Object {
	o := Eval(pe.Object, env)
	if isError(o) {
		return o
	}

	switch obj := o.(type) {
	case *object.String:
		// Special .ok property of commands
		if pe.Property.String() == "ok" {
			if obj.Ok != nil {
				return obj.Ok
			}

			return FALSE
		}
	case *object.Hash:
		return evalHashIndexExpression(obj, &object.String{Value: pe.Property.String()})
	}

	return newError("invalid property '%s' on type %s", pe.Property.String(), o.Type())
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func applyMethod(o object.Object, method string, args []object.Object) object.Object {
	f, ok := Fns[method]

	if !ok {
		return newError("%s does not have method '%s()'", o.Type(), method)
	}

	if !util.Contains(f.Types, string(o.Type())) && len(f.Types) != 0 {
		return newError("cannot call method '%s()' on '%s'", method, o.Type())
	}

	args = append([]object.Object{o}, args...)
	return f.Fn(args...)
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ && index.Type() == object.STRING_OBJ:
		return evalHashIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s on %s", index.Inspect(), left.Type())
	}
}

func evalStringIndexExpression(array, index object.Object) object.Object {
	stringObject := array.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return &object.String{Value: string(stringObject.Value[idx])}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalCommandExpression(cmd string, env *object.Environment) object.Object {
	re := regexp.MustCompile("\\$([a-zA-Z_]{1,})")
	cmd = re.ReplaceAllStringFunc(cmd, func(m string) string {
		v, ok := env.Get(m[1:])

		if !ok {
			return ""
		}

		return v.Inspect()
	})

	commands := []string{"-c", cmd}
	c := exec.Command("bash", commands...)
	c.Env = os.Environ()
	var out bytes.Buffer
	var stderr bytes.Buffer
	c.Stdout = &out
	c.Stderr = &stderr
	err := c.Run()

	if err != nil {
		return &object.String{Value: stderr.String(), Ok: FALSE}
	}

	return &object.String{Value: out.String(), Ok: TRUE}
}
