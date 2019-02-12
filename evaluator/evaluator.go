package evaluator

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/abs-lang/abs/ast"
	"github.com/abs-lang/abs/lexer"
	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/token"
	"github.com/abs-lang/abs/util"
)

var (
	NULL  = &object.Null{}
	EOF   = &object.Error{Message: "EOF"}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	Fns   map[string]*object.Builtin
)

// This program's lexer used for error location in Eval(program)
var lex *lexer.Lexer

func init() {
	Fns = getFns()
}

func newError(tok token.Token, format string, a ...interface{}) *object.Error {
	// get the token position from the error node and append the offending line to the error message
	lineNum, column, errorLine := lex.ErrorLine(tok.Position)
	errorPosition := fmt.Sprintf("\n\t[%d:%d]\t%s", lineNum, column, errorLine)
	return &object.Error{Message: fmt.Sprintf(format, a...) + errorPosition}
}

// BeginEval (program, env, lexer) object.Object
// REPL and testing modules call this function to init the global lexer pointer for error location
// NB. Eval(node, env) is recursive
func BeginEval(program ast.Node, env *object.Environment, lexer *lexer.Lexer) object.Object {
	// global lexer
	lex = lexer
	// run the evaluator
	return Eval(program, env)
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
		err := evalAssignment(node, env)

		if isError(err) {
			return err
		}

		return NULL
	// Expressions
	case *ast.NumberLiteral:
		return &object.Number{Token: node.Token, Value: node.Value}

	case *ast.NullLiteral:
		return NULL

	case *ast.StringLiteral:
		return &object.String{Token: node.Token, Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Token, node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Token, node.Operator, left, right)

	case *ast.CompoundAssignment:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		// multi-character operators like "+=" and "**=" are reduced to "+" or "**" for evalInfixExpression()
		op := node.Operator
		if len(op) >= 2 {
			op = op[:len(op)-1]
		}
		// get the result of the infix operation
		expr := evalInfixExpression(node.Token, op, left, right)
		if isError(expr) {
			return expr
		}
		switch nodeLeft := node.Left.(type) {
		case *ast.Identifier:
			env.Set(nodeLeft.String(), expr)
			return NULL
		case *ast.IndexExpression:
			// support index assignment expressions: a[0] += 1, h["a"] += 1
			return evalIndexAssignment(nodeLeft, expr, env)
		case *ast.PropertyExpression:
			// support assignment to hash property: h.a += 1
			return evalPropertyAssignment(nodeLeft, expr, env)
		}
		// otherwise
		env.Set(node.Left.String(), expr)
		return NULL

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
		return &object.Function{Token: node.Token, Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(node.Token, function, args)

	case *ast.MethodExpression:
		o := Eval(node.Object, env)
		if isError(o) {
			return o
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyMethod(node.Token, o, node.Method.String(), args)

	case *ast.PropertyExpression:
		return evalPropertyExpression(node, env)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Token: node.Token, Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(node.Token, left, index)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.CommandExpression:
		return evalCommandExpression(node.Token, node.Value, env)

	}

	return NULL
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

// support index assignment expressions: a[0] = 1, h["a"] = 1
func evalIndexAssignment(iex *ast.IndexExpression, expr object.Object, env *object.Environment) object.Object {
	leftObj := Eval(iex.Left, env)
	index := Eval(iex.Index, env)
	if leftObj.Type() == object.ARRAY_OBJ {
		arrayObject := leftObj.(*object.Array)
		idx := index.(*object.Number).Int()
		elems := arrayObject.Elements
		if idx < 0 {
			return newError(iex.Token, "index out of range: %d", idx)
		}
		if idx >= len(elems) {
			// expand the array by appending Null objects
			for i := len(elems); i <= idx; i++ {
				elems = append(elems, NULL)
			}
			arrayObject.Elements = elems
		}
		elems[idx] = expr
		return NULL
	}
	if leftObj.Type() == object.HASH_OBJ {
		hashObject := leftObj.(*object.Hash)
		key, ok := index.(object.Hashable)
		if !ok {
			return newError(iex.Token, "unusable as hash key: %s", index.Type())
		}
		hashed := key.HashKey()
		pair := object.HashPair{Key: index, Value: expr}
		hashObject.Pairs[hashed] = pair
		return NULL
	}
	return NULL
}

// support assignment to hash property: h.a = 1
func evalPropertyAssignment(pex *ast.PropertyExpression, expr object.Object, env *object.Environment) object.Object {
	leftObj := Eval(pex.Object, env)
	if leftObj.Type() == object.HASH_OBJ {
		hashObject := leftObj.(*object.Hash)
		prop := &object.String{Token: pex.Token, Value: pex.Property.String()}
		hashed := prop.HashKey()
		pair := object.HashPair{Key: prop, Value: expr}
		hashObject.Pairs[hashed] = pair
		return NULL
	}
	return newError(pex.Token, "can only assign to hash property, got %s", leftObj.Type())
}

func evalAssignment(as *ast.AssignStatement, env *object.Environment) object.Object {
	val := Eval(as.Value, env)
	if isError(val) {
		return val
	}

	// regular assignment x = 0
	if as.Name != nil {
		env.Set(as.Name.Value, val)
		return nil
	}

	// destructuring x, y = [1, 2]
	if len(as.Names) > 0 {
		switch v := val.(type) {
		case *object.Array:
			elements := v.Elements
			for i, name := range as.Names {
				if i < len(elements) {
					env.Set(name.String(), elements[i])
					continue
				}

				env.Set(name.String(), NULL)
			}
		case *object.Hash:
			for _, name := range as.Names {
				x, ok := v.GetPair(name.String())

				if ok {
					env.Set(name.String(), x.Value)
				} else {
					env.Set(name.String(), NULL)
				}
			}
		default:
			return newError(as.Token, "wrong assignment, expected identifier or array destructuring, got %s (%s)", val.Type(), val.Inspect())
		}

		return nil
	}
	// support assignment to indexed expressions: a[0] = 1, h["a"] = 1
	if as.Index != nil {
		return evalIndexAssignment(as.Index, val, env)
	}
	// support assignment to hash property h.a = 1
	if as.Property != nil {
		return evalPropertyAssignment(as.Property, val, env)
	}

	return nil
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(tok token.Token, operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(tok, right)
	case "~":
		return evalTildePrefixOperatorExpression(tok, right)
	default:
		return newError(tok, "unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(
	tok token.Token, operator string,
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
	case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
		return evalNumberInfixExpression(tok, operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(tok, operator, left, right)
	case left.Type() == object.ARRAY_OBJ && right.Type() == object.ARRAY_OBJ:
		return evalArrayInfixExpression(tok, operator, left, right)
	case left.Type() == object.HASH_OBJ && right.Type() == object.HASH_OBJ:
		return evalHashInfixExpression(tok, operator, left, right)
	case operator == "in":
		return evalInExpression(left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError(tok, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		case *object.Number:
			if o.Value == o.ZeroValue() {
				return TRUE
			}

			return FALSE
		default:
			return FALSE
		}
	}
}

func evalTildePrefixOperatorExpression(tok token.Token, right object.Object) object.Object {
	switch o := right.(type) {
	case *object.Number:
		return &object.Number{Value: float64(^int64(o.Value))}
	default:
		return newError(tok, "Bitwise not (~) can only be applied to numbers, got %s (%s)", o.Type(), o.Inspect())
	}
}

func evalMinusPrefixOperatorExpression(tok token.Token, right object.Object) object.Object {
	if right.Type() != object.NUMBER_OBJ {
		return newError(tok, "unknown operator: -%s", right.Type())
	}

	value := right.(*object.Number).Value
	return &object.Number{Value: -value}
}

func evalNumberInfixExpression(
	tok token.Token, operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Number).Value
	rightVal := right.(*object.Number).Value
	switch operator {
	case "+":
		return &object.Number{Token: tok, Value: leftVal + rightVal}
	case "-":
		return &object.Number{Token: tok, Value: leftVal - rightVal}
	case "*":
		return &object.Number{Token: tok, Value: leftVal * rightVal}
	case "/":
		return &object.Number{Token: tok, Value: leftVal / rightVal}
	case "**":
		// TODO this does not support floats
		return &object.Number{Token: tok, Value: math.Pow(leftVal, rightVal)}
	case "%":
		return &object.Number{Token: tok, Value: math.Mod(leftVal, rightVal)}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<=>":
		i := &object.Number{Token: tok}

		if leftVal == rightVal {
			i.Value = 0
		} else if leftVal > rightVal {
			i.Value = 1
		} else {
			i.Value = -1
		}

		return i
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "&":
		return &object.Number{Token: tok, Value: float64(int64(leftVal) & int64(rightVal))}
	case "|":
		return &object.Number{Token: tok, Value: float64(int64(leftVal) | int64(rightVal))}
	case ">>":
		return &object.Number{Token: tok, Value: float64(uint64(leftVal) >> uint64(rightVal))}
	case "<<":
		return &object.Number{Token: tok, Value: float64(uint64(leftVal) << uint64(rightVal))}
	case "^":
		return &object.Number{Token: tok, Value: float64(int64(leftVal) ^ int64(rightVal))}
	case "~":
		return &object.Boolean{Token: tok, Value: int64(leftVal) == int64(rightVal)}
	// A range results in an array of integers from left to right
	case "..":
		a := make([]object.Object, 0)

		for i := leftVal; i <= rightVal; i++ {
			a = append(a, &object.Number{Token: tok, Value: float64(i)})
		}
		return &object.Array{Token: tok, Elements: a}
	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(
	tok token.Token,
	operator string,
	left, right object.Object,
) object.Object {

	if operator == "+" {
		leftVal := left.(*object.String).Value
		rightVal := right.(*object.String).Value
		return &object.String{Token: tok, Value: leftVal + rightVal}
	}

	if operator == "==" {
		return &object.Boolean{Token: tok, Value: left.(*object.String).Value == right.(*object.String).Value}
	}

	if operator == "!=" {
		return &object.Boolean{Token: tok, Value: left.(*object.String).Value != right.(*object.String).Value}
	}

	if operator == "~" {
		return &object.Boolean{Token: tok, Value: strings.ToLower(left.(*object.String).Value) == strings.ToLower(right.(*object.String).Value)}
	}

	if operator == "in" {
		return evalInExpression(left, right)
	}

	return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalArrayInfixExpression(
	tok token.Token,
	operator string,
	left, right object.Object,
) object.Object {
	if operator == "+" {
		leftVal := left.(*object.Array).Elements
		rightVal := right.(*object.Array).Elements
		return &object.Array{Token: tok, Elements: append(leftVal, rightVal...)}
	}

	return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalHashInfixExpression(
	tok token.Token,
	operator string,
	left, right object.Object,
) object.Object {
	leftHashObject := left.(*object.Hash)
	rightHashObject := right.(*object.Hash)
	if operator == "+" {
		leftVal := leftHashObject.Pairs
		rightVal := rightHashObject.Pairs
		for _, rightPair := range rightVal {
			key := rightPair.Key
			hashed := key.(object.Hashable).HashKey()
			leftVal[hashed] = object.HashPair{Key: key, Value: rightPair.Value}
		}
		return &object.Hash{Token: tok, Pairs: leftVal}
	}

	return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalInExpression(
	left, right object.Object,
) object.Object {
	var found bool

	switch rightObj := right.(type) {
	case *object.Array:
		switch needle := left.(type) {
		case *object.String:
			for _, v := range rightObj.Elements {
				if v.Inspect() == needle.Value && v.Type() == object.STRING_OBJ {
					found = true
					break // Let's get outta here!
				}
			}
		case *object.Number:
			for _, v := range rightObj.Elements {
				// Quite ghetto but also the easiest way out
				// Instead of doing type checking on the argument,
				// we received back its string representation.
				// If they match, we then check that its type was
				// integer.
				if v.Inspect() == strconv.Itoa(int(needle.Value)) && v.Type() == object.NUMBER_OBJ {
					found = true
					break // Let's get outta here!
				}
			}
		}
	case *object.String:
		if left.Type() == object.STRING_OBJ {
			found = strings.Contains(right.Inspect(), left.Inspect())
		}
	case *object.Hash:
		if left.Type() == object.STRING_OBJ {
			_, ok := rightObj.GetPair(left.(*object.String).Value)
			found = ok
		}
	default:
		return newError(tok, "'in' operator not supported on %s", right.Type())
	}

	return &object.Boolean{Token: tok, Value: found}
}

func evalIfExpression(
	ie *ast.IfExpression,
	env *object.Environment,
) object.Object {
	for _, scenario := range ie.Scenarios {
		condition := Eval(scenario.Condition, env)

		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			return Eval(scenario.Consequence, env)
		}
	}

	return NULL
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
		evaluated := Eval(we.Consequence, env)

		if isError(evaluated) {
			return evaluated
		}

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

	switch i := iterable.(type) {
	case object.Iterable:
		defer func() {
			i.Reset()
		}()

		return loopIterable(i.Next, env, fie)
	case *object.Builtin:
		if i.Next == nil {
			return newError(fie.Token, "builtin function cannot be used in loop")
		}

		return loopIterable(i.Next, env, fie)
	default:
		return newError(fie.Token, "'%s' is a %s, not an iterable, cannot be used in for loop", i.Inspect(), i.Type())
	}
}

func loopIterable(next func() (object.Object, object.Object), env *object.Environment, fie *ast.ForInExpression) object.Object {
	k, v := next()

	if k == nil || v == EOF {
		return NULL
	}

	// set the special k v variables in the
	// environment
	env.Set(fie.Key, k)
	env.Set(fie.Value, v)
	err := Eval(fie.Block, env)

	if isError(err) {
		return err
	}

	if k != nil {
		return loopIterable(next, env, fie)
	}

	return NULL
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

	return newError(node.Token, "identifier not found: "+node.Value)
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
	case *object.Number:
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
		return evalHashIndexExpression(obj.Token, obj, &object.String{Token: pe.Token, Value: pe.Property.String()})
	}

	return newError(pe.Token, "invalid property '%s' on type %s", pe.Property.String(), o.Type())
}

func applyFunction(tok token.Token, fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		extendedEnv, err := extendFunctionEnv(fn, args)

		if err != nil {
			return err
		}
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(tok, args...)

	default:
		return newError(tok, "not a function: %s", fn.Type())
	}
}

func applyMethod(tok token.Token, o object.Object, method string, args []object.Object) object.Object {
	f, ok := Fns[method]

	if !ok {
		return newError(tok, "%s does not have method '%s()'", o.Type(), method)
	}

	if !util.Contains(f.Types, string(o.Type())) && len(f.Types) != 0 {
		return newError(tok, "cannot call method '%s()' on '%s'", method, o.Type())
	}

	args = append([]object.Object{o}, args...)
	return f.Fn(tok, args...)
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) (*object.Environment, *object.Error) {
	env := object.NewEnclosedEnvironment(fn.Env)

	if len(args) != len(fn.Parameters) {
		return nil, newError(fn.Token, "Wrong number of arguments passed to %s. Want %s, got %s", fn.Inspect(), fn.Parameters, args)
	}

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env, nil
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalIndexExpression(tok token.Token, left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.NUMBER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ && index.Type() == object.STRING_OBJ:
		return evalHashIndexExpression(tok, left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.NUMBER_OBJ:
		return evalStringIndexExpression(tok, left, index)
	default:
		return newError(tok, "index operator not supported: %s on %s", index.Inspect(), left.Type())
	}
}

func evalStringIndexExpression(tok token.Token, array, index object.Object) object.Object {
	stringObject := array.(*object.String)
	idx := index.(*object.Number).Int()
	max := len(stringObject.Value) - 1

	if idx < 0 || idx > max {
		return NULL
	}

	return &object.String{Token: tok, Value: string(stringObject.Value[idx])}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Number).Int()
	max := len(arrayObject.Elements) - 1

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
			return newError(node.Token, "unusable as hash key: %s", key.Type())
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

func evalHashIndexExpression(tok token.Token, hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError(tok, "unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalCommandExpression(tok token.Token, cmd string, env *object.Environment) object.Object {
	// Match all strings preceded by
	// a $ or a \$
	re := regexp.MustCompile("(\\\\)?\\$([a-zA-Z_]{1,})")
	cmd = re.ReplaceAllStringFunc(cmd, func(m string) string {
		// If the string starts with a backslash,
		// that's an escape, so we should replace
		// it with the remaining portion of the match.
		// \$VAR becomes $VAR
		if string(m[0]) == "\\" {
			return m[1:]
		}

		// If the string starts with $, then
		// it's an interpolation. Let's
		// replace $VAR with the variable
		// named VAR in the ABS' environment.
		// If the variable is not found, we
		// just dump an empty string
		v, ok := env.Get(m[1:])

		if !ok {
			return ""
		}

		return v.Inspect()
	})

	// thanks to @haifenghuang
	var commands []string
	var executor string
	if runtime.GOOS == "windows" {
		commands = []string{"/C", cmd}
		executor = "cmd.exe"
	} else { //assume it's linux, darwin, freebsd, openbsd, solaris, etc
		commands = []string{"-c", cmd}
		executor = "bash"
	}
	c := exec.Command(executor, commands...)
	c.Env = os.Environ()
	var out bytes.Buffer
	var stderr bytes.Buffer
	c.Stdin = os.Stdin
	c.Stdout = &out
	c.Stderr = &stderr
	err := c.Run()

	if err != nil {
		return &object.String{Token: tok, Value: stderr.String(), Ok: FALSE}
	}

	return &object.String{Token: tok, Value: strings.Trim(out.String(), "\n"), Ok: TRUE}
}
