package evaluator

import (
	"testing"

	"github.com/abs-lang/abs/lexer"
	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/parser"
)

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.5", 5.5},
		{"-5.5", -5.5},
		{"5.5 + 3.7", 9.2},
		{"5.5 * 2", 11},
		{"1 / 3", 0.3333333333333333},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestEvalNumberExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.5", 5.5},
		{"1.1 + 2.1", 3.2},
		{"5.5 + 2.2", 7.7},
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"2 ** 2", 4},
		{"10 ** 0", 1},
		{"10 ** 0 - 1", 0},
		{"10 && 0", 0},
		{"10 && 1", 1},
		{"0 && 3", 0},
		{`"hello" && 10`, 10},
		{"0 <=> 1", -1},
		{"1 <=> 1", 0},
		{"2 <=> 1", 1},
		{"2 % 1", 0},
		{"3 % 2", 1},
		{"a = 5; a += 1; a", 6},
		{"a = 5; a -= 1; a", 4},
		{"a = 10; a /= 2; a", 5},
		{"a = 5; a *= 2; a", 10},
		{"a = 5; a **= 2; a", 25},
		{"a = 5; a %= 3; a", 2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"9999999999.str()", "9999999999"},
		{"12.1.str()", "12.1"},
		{"12.123456789.str()", "12.123456789"},
		{`"5"`, "5"},
		{`"5" + "5"`, "55"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"a" == "a"`, true},
		{`"a" == "b"`, false},
		{`"a" ~ "b"`, false},
		{`"a" ~ "A"`, true},
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 <= 1", true},
		{"1 > 1", false},
		{"1 >= 1", true},
		{"0 >= 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{`a = "hello"; !a.ok`, true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!!0", false},
		{`!!""`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true { 10 }", 10},
		{"if false { 10 }", nil},
		{"if 1 { 10 }", 10},
		{"if 1 < 2 { 10 }", 10},
		{"if 1 > 2 { 10 }", nil},
		{"if 1 > 2 { 10 } else { 20 }", 20},
		{"if 1 < 2 { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testNumberObject(t, evaluated, float64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestForExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"a = 0; for x = 0; x < 10; x = x + 1 { a = a + 1}; a", 10},
		{"a = 0; for x = 0; x < y; x = x + 1 { a = a + 1}; a", "identifier not found: y"},
		{"a = 0; increment = f(x) {x+1}; for x = 0; x < 10; x = increment(x) { a = a + 1}; a", 10},
		{`a = 0; for k = 0; k < 10; k = k + 1 { a = a + 1}; k`, "identifier not found: k"},
		{`k = 100; for k = 0; k < 10; k = k { k = k + 1}; k`, 100},
		{`k = 100; for k = y; k < 10; k = k { k = 9 }; k`, "identifier not found: y"},
		{`k = 100; for k = 0; k <= 10; k = k { k = y }; k`, "identifier not found: y"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testNumberObject(t, evaluated, float64(integer))
		} else {
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
				continue
			}

			if errObj.Message != tt.expected {
				t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
			}
		}
	}
}

func TestBitwiseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"1 & 1", 1},
		{"1 & 1.1", 1},
		{"1 & 0", 0},
		{`1 & ""`, "type mismatch: NUMBER & STRING"},
		{"1 ^ 1", 0},
		{"1 ^ 1.1", 0},
		{"1 ^ 0", 1},
		{`1 ^ ""`, "type mismatch: NUMBER ^ STRING"},
		{"1 | 1", 1},
		{"1 | 1.1", 1},
		{"1 | 0", 1},
		{`1 | ""`, "type mismatch: NUMBER | STRING"},
		{"1 >> 1", 0},
		{"1 >> 1.1", 0},
		{"1 >> 0", 1},
		{`1 >> ""`, "type mismatch: NUMBER >> STRING"},
		{"1 << 1", 2},
		{"1 << 1.1", 2},
		{"1 << 0", 1},
		{`1 << ""`, "type mismatch: NUMBER << STRING"},
		{"~1", -2},
		{"~1.1", -2},
		{"~0", -1},
		{`~"abc"`, "Bitwise not (~) can only be applied to numbers, got STRING (abc)"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testNumberObject(t, evaluated, float64(integer))
		} else {
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
				continue
			}

			if errObj.Message != tt.expected {
				t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
			}
		}
	}
}

func TestForInExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"a = 0; for k, x in 1 { a = a + 1}; a", "'1' is a NUMBER, not an iterable, cannot be used in for loop"},
		{"a = 0; for k, x in 1..10 { a = a + 1}; a", 10},
		{"a = 0; for x in 1 { a = a + 1}; a", "'1' is a NUMBER, not an iterable, cannot be used in for loop"},
		{"a = 0; for x in 1..10 { a = a + 1}; a", 10},
		{`a = 0; for k, v in ["x", "y", "z"] { a = a + k}; a`, 3},
		{`for k, v in ["x", "y", "z"] {}; k`, "identifier not found: k"},
		{`for k, v in ["x", "y", "z"] {}; v`, "identifier not found: v"},
		{`k = 100; for k, v in ["x", "y", "z"] {}; k`, 100},
		{`v = 100; for k, v in ["x", "y", "z"] {}; v`, 100},
		{`for k, v in ["x", "y", "z"] {k=y}; v`, "identifier not found: y"},
		{`for k, v in ["x", "y", z] {k=y}; v`, "'ERROR: identifier not found: z' is a ERROR, not an iterable, cannot be used in for loop"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testNumberObject(t, evaluated, float64(integer))
		} else {
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
				continue
			}

			if errObj.Message != tt.expected {
				t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
			}
		}
	}
}

func TestWhileExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"while true { echo x }", "identifier not found: x"},
		{"a = 0; while (a < 10) { a = a + 1 }; a", 10},
		{`a = ""; while (len(a) < 3) { a = a + "a" }; a`, "aaa"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch tt.expected.(type) {
		case int:
			testNumberObject(t, evaluated, float64(tt.expected.(int)))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				testStringObject(t, evaluated, tt.expected.(string))
				continue
			}

			if errObj.Message != tt.expected {
				t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
			}
		default:
			panic("should not reach here")
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { return 10; }", 10},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }

  return 1;
}
`,
			10,
		},
		{
			`
fn = f(x) {
  return x;
  x + 10;
};
fn(10);`,
			10,
		},
		{
			`
fn = f(x) {
   result = x + 10;
   return result;
   return 10;
};
fn(10);`,
			20,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testNumberObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: NUMBER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: NUMBER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"true + false + true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`{"name": "Abs"}[f(x) { x }];`,
			`index operator not supported: f(x) {
x
} on HASH`,
		},
		{
			`999[1]`,
			"index operator not supported: 1 on NUMBER",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestAssignStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"a = 5; a;", 5},
		{"a = 5 * 5; a;", 25},
		{"a = 5; b = a; b;", 5},
		{"a = 5; b = a; c = a + b + 5; c;", 15},
		{"[a] = [10]; a;", 10},
		{"[a] = 10; a;", "wrong assignment, expected identifier or array destructuring, got NUMBER (10)"},
		{"[a, b, c] = [1]; a;", 1},
		{"[a, b, c] = [1]; b;", nil},
		{`
tz = $(echo "10/20");
[a, b] = tz.split("/")
a.int()
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testNumberObject(t, testEval(tt.input), float64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "f(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"identity = f(x) { x; }; identity(5);", 5},
		{"identity = f(x) { return x; }; identity(5);", 5},
		{"double = f(x) { x * 2; }; double(5);", 10},
		{"add = f(x, y) { x + y; }; add(5, 5);", 10},
		{"add = f(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"f(x) { x; }(5)", 5},
		{"f(x) { x; }()", "Wrong number of arguments passed to f(x) {\nx\n}. Want [x], got []"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
				continue
			}

			if errObj.Message != tt.expected {
				t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
			}
		case float64:
			testNumberObject(t, evaluated, float64(expected))
		}
	}
}

func TestEnclosingEnvironments(t *testing.T) {
	input := `
first = 10;
second = 10;
third = 10;

ourFunction = f(first) {
  second = 20;

  first + second + third;
};

ourFunction(20) + first + second;`

	testNumberObject(t, testEval(input), float64(70))
}

func TestClosures(t *testing.T) {
	input := `
newAdder = f(x) {
  f(y) { x + y };
};

addTwo = newAdder(2);
addTwo(2);`

	testNumberObject(t, testEval(input), float64(4))
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestArrayConcatenation(t *testing.T) {
	input := `[1] + [2]`

	evaluated := testEval(input)
	arr, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	testNumberObject(t, arr.Elements[0], float64(1))
	testNumberObject(t, arr.Elements[1], float64(2))
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`[1,2,"a"].map(int)`, "int(...) can only be called on strings which represent numbers, 'a' given"},
		{`[1,2,"a"].filter(int)`, "int(...) can only be called on strings which represent numbers, 'a' given"},
		{`[1, null].sum()`, "sum(...) can only be called on an homogeneous array, got [1, null]"},
		{`[null, null].sum()`, "sum(...) can only be called on arrays of numbers, got [null, null]"},
		{`[].sum()`, 0},
		{`[1, 2].sum()`, 3},
		{`[1].map(f(x) { y = x + 1 }).str()`, "[null]"},
		{`contains("hello", "lo")`, true},
		{`contains("hello", "_")`, false},
		{`contains("hello", 1)`, false},
		{`contains(["hello"], "hello")`, true},
		{`contains([1, 2], 2)`, true},
		{`contains(["hello"], "_")`, false},
		{`contains(["hello"], 1)`, false},
		{`contains("1", 1)`, false},
		{`contains(1..100000, 1)`, true},
		{`find([1,2,3,3], f(x) {x == 3})`, 3},
		{`find([1,2], f(x) {x == "some"})`, nil},
		{`arg("o")`, "argument 0 to arg(...) is not supported (got: o, allowed: NUMBER)"},
		{`arg(3)`, ""},
		{`rand(1)`, 0},
		{`int(10)`, 10},
		{`int(10.5)`, 10},
		{`number("10")`, 10},
		{`number("10.55")`, 10.55},
		{`int("10")`, 10},
		{`int("10.5")`, 10},
		{`int("abc")`, `int(...) can only be called on strings which represent numbers, 'abc' given`},
		{`int([])`, "argument 0 to int(...) is not supported (got: [], allowed: NUMBER, STRING)"},
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument 0 to len(...) is not supported (got: 1, allowed: STRING, ARRAY)"},
		{`len("one", "two")`, "wrong number of arguments to len(...): got=2, want=1"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`echo("hello", "world!")`, nil},
		{`env("CONTEXT")`, "abs"},
		{`type("SOME")`, "STRING"},
		{`type(1)`, "NUMBER"},
		{`type({})`, "HASH"},
		{`type([])`, "ARRAY"},
		{`type("{}".json())`, "HASH"},
		{`"{\"a\": null}".json().a`, nil},
		{`type(null)`, "NULL"},
		{`"{\"k\": \"v\"}".json()["k"]`, "v"},
		{`"hello".json()`, "argument to `json` must be a valid JSON object, got 'hello'"},
		{`"\"hello".json()`, "argument to `json` must be a valid JSON object, got '\"hello'"},
		{`split("a\"b\"c", "\"")`, []string{"a", "b", "c"}},
		{`lines("a
b
c")`, []string{"a", "b", "c"}},
		{`[1, 2].sort()`, []int{1, 2}},
		{`["b", "a"].sort()`, []string{"a", "b"}},
		{`["b", 1].sort()`, "argument to `sort` must be an homogeneous array (elements of the same type), got [b, 1]"},
		{`[{}].sort()`, "cannot sort an array with given elements elements ([{}])"},
		{`[[]].sort()`, "cannot sort an array with given elements elements ([[]])"},
		{`[1, 2].some(f(x) {x == 2})`, true},
		{`[].some(f(x) {x})`, false},
		{`[1, 2].every(f(x) { return x == 2 || x == 1})`, true},
		{`[].every(f(x) {x})`, true},
		{`[1,2,3].every(f(x) {x == 1})`, false},
		{`[1,2,3].filter(f(x) {x == 1})`, []int{1}},
		{`[].shift()`, nil},
		{`[1, 2].shift()`, 1},
		{`a = [1, 2]; a.shift(); a`, []int{2}},
		{`[1, 2].reverse();`, []int{2, 1}},
		{`[1, 2].push("a");`, []interface{}{1, 2, "a"}},
		{`[1, 2].pop();`, 2},
		{`a = [1, 2]; a.pop(); a`, []int{1}},
		{`[1, 2].keys()`, []int{0, 1}},
		{`[1, 2].join("-")`, "1-2"},
		{`["a", "b"].join("-")`, "a-b"},
		{`"a".any("b")`, false},
		{`"a".any("a")`, true},
		{`"a".prefix("b")`, false},
		{`"a".prefix("a")`, true},
		{`"a".suffix("b")`, false},
		{`"a".suffix("a")`, true},
		{`"ab".index("b")`, 1},
		{`"a".index("b")`, nil},
		{`"abb".last_index("b")`, 2},
		{`"a".last_index("b")`, nil},
		{`"a".repeat(3)`, "aaa"},
		{`"a".repeat(3)`, "aaa"},
		{`"abc".slice(0, 0)`, "abc"},
		{`"abc".slice(1, 0)`, "bc"},
		{`"abc".slice(1, 2)`, "b"},
		{`"abc".slice(0, 6)`, "abc"},
		{`"abc".slice(10, 10)`, ""},
		{`"abc".slice(10, 20)`, ""},
		{`"abc".slice(-1, 0)`, "c"},
		{`"abc".slice(-20, 0)`, "abc"},
		{`"abc".slice(-20, 2)`, "ab"},
		{`"abc".slice(-1, 3)`, "c"},
		{`"abc".slice(-1, 1)`, "c"},
		{`"hello %s".fmt("world")`, "hello world"},
		{`"hello %s".fmt()`, "hello %!s(MISSING)"},
		{`"hello %s".fmt(1)`, "hello 1"},
		{`"hello %s".fmt({})`, "hello {}"},
		{`[1,2,3].slice(0, 0)`, []int{1, 2, 3}},
		{`[1,2,3].slice(1, 0)`, []int{2, 3}},
		{`[1,2,3].slice(1, 2)`, []int{2}},
		{`[1,2,3].slice(0, 6)`, []int{1, 2, 3}},
		{`[1,2,3].slice(10, 10)`, []int{}},
		{`[1,2,3].slice(10, 20)`, []int{}},
		{`[1,2,3].slice(-1, 0)`, []int{3}},
		{`[1,2,3].slice(-20, 0)`, []int{1, 2, 3}},
		{`[1,2,3].slice(-20, 2)`, []int{1, 2}},
		{`[1,2,3].slice(-1, 3)`, []int{3}},
		{`[1,2,3].slice(-1, 1)`, []int{3}},
		{`"a".replace("a", "b", -1)`, "b"},
		{`"a".str()`, "a"},
		{`1.str()`, "1"},
		{`[1].str()`, "[1]"},
		{`{"a": 10}.str()`, `{a: 10}`},
		{`"a great movie".title()`, "A Great Movie"},
		{`"A great movie".lower()`, "a great movie"},
		{`"A great movie".upper()`, "A GREAT MOVIE"},
		{`"  A great movie  ".trim()`, "A great movie"},
		{`"  A great movie  ".trim_by(" A")`, "great movie"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testNumberObject(t, evaluated, float64(expected))
		case float64:
			testNumberObject(t, evaluated, float64(expected))
		case nil:
			testNullObject(t, evaluated)
		case bool:
			testBooleanObject(t, evaluated, expected)
		case string:
			s, ok := evaluated.(*object.String)
			if ok {
				if s.Value != tt.expected {
					t.Errorf("object is not the right string. got=%s want:%s", s.Value, tt.expected)
				}

				continue
			}

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d",
					len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testNumberObject(t, array.Elements[i], float64(expectedElem))
			}
		case []string:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testStringObject(t, array.Elements[i], expectedElem)
			}
		case []interface{}:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Elements))
				continue
			}
		}
	}
}

func TestLogicalOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`1 && 0`, 0},
		{`1 && 2`, 2},
		{`"hello" && 2`, 2},
		{`"" && 2`, ""},
		{`"hello" && ""`, ""},
		{`len("hello") && 2`, 2},
		{`1 || 0`, 1},
		{`1 || 2`, 1},
		{`"hello" || 2`, "hello"},
		{`"" || 2`, 2},
		{`"hello" || ""`, "hello"},
		{`len("hello") || ""`, 5},
		{`
		(("") || ("") || (0 || 0 || 0)) || ""
`, ""},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testNumberObject(t, evaluated, float64(expected))
		case nil:
			testNullObject(t, evaluated)
		case string:
			s, ok := evaluated.(*object.String)
			if ok {
				if s.Value != tt.expected {
					t.Errorf("object is not the right string. got=%s want:%s", s.Value, tt.expected)
				}

				continue
			}

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testNumberObject(t, array.Elements[i], float64(expectedElem))
			}
		}
	}
}

func TestRangesOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`1..0`, []int{}},
		{`-1..0`, []int{-1, 0}},
		{`1..1`, []int{1}},
		{`1..2`, []int{1, 2}},
		{`len("a")..len("aa")`, []int{1, 2}},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testNumberObject(t, evaluated, float64(expected))
		case nil:
			testNullObject(t, evaluated)
		case string:
			s, ok := evaluated.(*object.String)
			if ok {
				if s.Value != tt.expected {
					t.Errorf("object is not the right string. got=%s want:%s", s.Value, tt.expected)
				}

				continue
			}

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testNumberObject(t, array.Elements[i], float64(expectedElem))
			}
		}
	}
}

func TestBuiltinProperties(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"a".ok`, false},
		{`"a".inv`, "invalid property 'inv' on type STRING"},
		{"a = $(date);\na.ok", true},
		{`{}.a`, nil},
		{`{"a": 1}.a`, 1},
		{`{1: 1}.1`, "unusable as hash key: NUMBER"},
		{`[].a`, "invalid property 'a' on type ARRAY"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case bool:
			testBooleanObject(t, evaluated, bool(expected))
		case int:
			testNumberObject(t, evaluated, float64(expected))
		case nil:
			testNullObject(t, evaluated)
		case string:
			s, ok := evaluated.(*object.String)
			if ok {
				if s.Value != tt.expected {
					t.Errorf("object is not the right string. got=%s want:%s", s.Value, tt.expected)
				}

				continue
			}

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d",
					len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testNumberObject(t, array.Elements[i], float64(expectedElem))
			}
		case []string:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testStringObject(t, array.Elements[i], expectedElem)
			}
		}
	}
}

func TestCommand(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`a = "A"; b = "B"; eee = "-e"; $(echo $eee -n $a$a$b$b$c$c)`, "AABB"},
		{`$(echo -n "123")`, "123"},
		{`$(echo -n hello world)`, "hello world"},
		{`$(echo hello world | xargs echo -n)`, "hello world"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case nil:
			testNullObject(t, evaluated)
		case string:
			stringObj, ok := evaluated.(*object.String)
			if !ok {
				t.Errorf("object is not String. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if stringObj.Value != expected {
				t.Errorf("wrong string. expected=%q, got=%q", expected, stringObj.Value)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testNumberObject(t, result.Elements[0], float64(1))
	testNumberObject(t, result.Elements[1], float64(4))
	testNumberObject(t, result.Elements[2], float64(6))
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"myArray = [1, 2, 3]; i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testNumberObject(t, evaluated, float64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]float64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testNumberObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}.foo`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{}.foo`,
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testNumberObject(t, evaluated, float64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestStringIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`"123"[10]`,
			nil,
		},
		{
			`"123"[1]`,
			"2",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		s, ok := tt.expected.(string)
		if ok {
			testStringObject(t, evaluated, s)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testNumberObject(t *testing.T, obj object.Object, expected interface{}) bool {
	result, ok := obj.(*object.Number)
	if !ok {
		t.Errorf("object is not Number. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%v, want=%v", result.Value, expected)
		return false
	}

	return true
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s",
			result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
