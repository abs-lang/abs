package evaluator

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/abs-lang/abs/lexer"
	"github.com/abs-lang/abs/object"
	"github.com/abs-lang/abs/parser"
)

func logErrorWithPosition(t *testing.T, msg string, expected interface{}) {
	errorStr := msg
	expected, _ = expected.(string)
	expectedStr := fmt.Sprintf("%s", expected)
	if strings.HasPrefix(errorStr, expectedStr) {
		// Only log when we're running the verbose tests
		if flag.Lookup("test.v").Value.String() == "true" {
			t.Log("expected error:", errorStr)
		}
	} else {
		expectedStr = fmt.Sprintf("ERROR: wrong error message. expected='%s',", expectedStr)
		t.Error(expectedStr, "\ngot=", errorStr)
	}
}
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

func TestEvalNumberAbbreviations(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5k", 5000},
		{"1m / 1M", 1},
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
		{"1e1", 10},
		{"1e-1", 0.1},
		{"1e+1", 10},
		{"1_000e1", 10000},
		{"10_000_000", 10000000},
		{"10_0.00_00", 100.0000},
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
		{`"\n"`, "\n"},
		{`"\r"`, "\r"},
		{`"\t"`, "\t"},
		{"12.123456789.str()", "12.123456789"},
		{`"nice 'escaping"`, "nice 'escaping"},
		{`'nice "escaping"`, `nice "escaping"`},
		{`'nice \'escaping`, `nice 'escaping`},
		{`"nice \"escaping"`, `nice "escaping`},
		{`"5"`, "5"},
		{`'5'`, "5"},
		{`'hello %s'.fmt("world")`, "hello world"},
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
		{`1 ~ 1`, true},
		{`1 ~ 1.5`, true},
		{`2 ~ 1.5`, false},
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

func TestLazyEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`x = false; x && x.lines`, false},
		{`y = true; y || y.lines`, true},
		{`x = true; y = false; x || (y && y.ok)`, true},
		{`true || 1..1000000000`, true},
		{`false && 1..1000000000`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func BenchmarkLazyEvaluation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testEval(`false && 1..` + strconv.Itoa(i))
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
		{"if 3 > 2 { 1 } else if 1 > 0 {2} else if 5 > 0 {3} else {4}", 1},
		{"if 1 > 2 { 1 } else if 1 > 0 {2} else if 5 > 0 {3} else {4}", 2},
		{"if 1 > 2 { 1 } else if 1 > 1 {2} else if 5 > 0 {3} else {4}", 3},
		{"if 1 > 2 { 1 } else if 1 > 1 {2} else if 5 > 10 {3} else {4}", 4},
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
		{`x = f() { for k = 0; k < 11; k = k + 1 { return k }}; x()`, 0},                           // https://github.com/abs-lang/abs/issues/303
		{`y = {}; x = f() { for k = 0; k < 11; k = k + 1 { y.x = k; return null; }}; x(); y.x`, 0}, // https://github.com/abs-lang/abs/issues/303
		{`y = {}; x = f() { for k = 0; k < 11; k = k + 1 { y.x = k; return 12; }}; x()`, 12},       // https://github.com/abs-lang/abs/issues/303
		{`y = {}; x = f() { for k = 0; k < 11; k = k + 1 { y.x = k; }}; x(); y.x`, 10},             // https://github.com/abs-lang/abs/issues/303
		{`y = []; for x in 1..10 { y.push(x) }; y.len()`, 10},                                      // https://github.com/abs-lang/abs/issues/324
		{`y = []; for x in 1..10 { if true { y.push(x) } }; y.len()`, 10},                          // https://github.com/abs-lang/abs/issues/324
		{`y = []; for x in 1..10 { if true { y.push(x); return 1 } }; y.len()`, 1},                 // https://github.com/abs-lang/abs/issues/324
		{`x = 0; for k = 0; k < 11; k = k + 1 { if k < 10 { break; }; x += k }; x`, 0},
		{`x = 0; for k = 0; k < 11; k = k + 1 { if k < 10 { continue; }; x += k }; x`, 10},
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
			logErrorWithPosition(t, errObj.Message, tt.expected)
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
			logErrorWithPosition(t, errObj.Message, tt.expected)
		}
	}
}

func TestStringWriters(t *testing.T) {
	tests := []struct {
		input   string
		content string
		file    string
	}{
		{`"abc" > "write.txt.ignore"`, "abc", "write.txt.ignore"},
		{`"" > "append.txt.ignore"; "abc" >> "append.txt.ignore"; "abc" >> "append.txt.ignore"`, "abcabc", "append.txt.ignore"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, true)

		content, err := ioutil.ReadFile(tt.file)

		if err != nil {
			t.Errorf("unable to read file %s", tt.file)
		}

		if string(content) != tt.content {
			t.Errorf("file content is wrong: wanted %s, got %s", tt.content, string(content))
		}
	}
}

func TestStringInterpolation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`a = "123"; "abc$a"`, "abc123"},
		{`a = "123"; "abc\$a"`, "abc$a"},
		{`a = "123"; "$$a$$a$$a"`, "$123$123$123"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestForInExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`x = 0; for v in 1..10 { if v < 10 { break; }; x += v }; x`, 0},
		{`x = f() { for x in 1..10 { return x }}; x()`, 1},                                   // https://github.com/abs-lang/abs/issues/303
		{`y = {}; x = f() { for x in 1000..2000 { y.x = x; return null; }}; x(); y.x`, 1000}, // https://github.com/abs-lang/abs/issues/303
		{`y = {}; x = f() { for x in 1000..2000 { y.x = x; return 12; }}; x()`, 12},          // https://github.com/abs-lang/abs/issues/303
		{`y = {}; x = f() { for x in 1000..2000 { y.x = x; }}; x(); y.x`, 2000},              // https://github.com/abs-lang/abs/issues/303
		{`x = 0; for v in 1..10 { if v < 10 { continue; }; x += v }; x`, 10},
		{"a = 1..3; b = 0; c = 0; for x in a { b = x }; for x in a { c = x }; c", 3}, // See: https://github.com/abs-lang/abs/issues/112
		{"a = 0; for k, x in 1 { a = a + 1}; a", "'1' is a NUMBER, not an iterable, cannot be used in for loop"},
		{"a = 0; for k, x in 1..10 { a = a + 1}; a", 10},
		{"a = 0; for x in 1 { a = a + 1}; a", "'1' is a NUMBER, not an iterable, cannot be used in for loop"},
		{"a = 0; for x in 1..10 { a = a + 1}; a", 10},
		{`a = 0; for k, v in {"a": 10} { a = v}; a`, 10},
		{`a = ""; b = "abc"; for k, v in {"a": 1, "b": 2, "c": 3} { a += k}; a == b`, true},
		{`a = 0; for k, v in ["x", "y", "z"] { a = a + k}; a`, 3},
		{`for k, v in ["x", "y", "z"] {}; k`, "identifier not found: k"},
		{`for k, v in ["x", "y", "z"] {}; v`, "identifier not found: v"},
		{`k = 100; for k, v in ["x", "y", "z"] {}; k`, 100},
		{`v = 100; for k, v in ["x", "y", "z"] {}; v`, 100},
		{`for k, v in ["x", "y", "z"] {k=y}; v`, "identifier not found: y"},
		{`for k, v in ["x", "y", z] {k=y}; v`, "'ERROR: identifier not found: z"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch ev := tt.expected.(type) {
		case int:
			testNumberObject(t, evaluated, float64(ev))
		case bool:
			testBooleanObject(t, evaluated, ev)
		default:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
				continue
			}
			logErrorWithPosition(t, errObj.Message, ev)
		}
	}
}

func TestForElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"a = 0; b = 1; for v in [] { x = a } else { x = b }; x", 1},
		{"a = 100; x = 0; for i in  1..-1 { x = i } else { x = a }; x", 100},
		{"v = 100; for k, v in [] { v = 0 } else {}; v", 100},
		{"for k, v in [] {} else { x = v }; x", "identifier not found: v"},
		{"a = 0; for k, v in [] { x = a } else { x = b }; x", "identifier not found: b"},
		{"for k, v in [] { x = 0 } else { x = 100 }; z", "identifier not found: z"},
		{"for i in 1..3 { x = i } else { x = 0 }; x", 3},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch ev := tt.expected.(type) {
		case int:
			testNumberObject(t, evaluated, float64(ev))
		case bool:
			testBooleanObject(t, evaluated, ev)
		default:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
				continue
			}
			logErrorWithPosition(t, errObj.Message, ev)
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
			logErrorWithPosition(t, errObj.Message, tt.expected)
		default:
			panic("should not reach here")
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"return;", nil},
		{"return", nil},
		{"return 1", 1},
		{"fn = f() { return }; fn()", nil},
		{"fn = f() { return 1 }; fn()", 1},
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
		switch tt.expected.(type) {
		case int:
			testNumberObject(t, evaluated, float64(tt.expected.(int)))
		case nil:
			testNullObject(t, evaluated)
		default:
			panic("should not reach here")
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
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
			// `{"name": "Abs"}[f(x) {x}];`,
			`{"name": "Abs"}[f(x) {x}];`,
			`index operator not supported: f(x) {x} on HASH`,
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
		logErrorWithPosition(t, errObj.Message, tt.expected)
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
		{"a, b, c = [1]; a;", 1},
		{`a, b, c = {"a": 1}; a;`, 1},
		{"a, b, c = [1]; b;", nil},
		{`a, b, c = {"a": 1}; b;`, nil},
		{`a = 10 + 1 + 2
b, c = [1, 2]; b`, 1},
		{`a = 10 + 1 + 2
		b, c = [1, 2]; a`, 13},
		{`
		tz = "10/20"
		a, b = tz.split("/")
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
			logErrorWithPosition(t, errObj.Message, tt.expected)
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
		{"f test() { return 12 }; test()", 12},
		{"f test(x) { return x }; test(12)", 12},
		{"identity = f(x) { x; }; identity(5);", 5},
		{"identity = f(x) { return x; }; identity(5);", 5},
		{"double = f(x) { x * 2; }; double(5);", 10},
		{"add = f(x, y) { x + y; }; add(5, 5);", 10},
		{"add = f(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"f(x) { x; }(5)", 5},
		{"f(x) { x; }()", "Wrong number of arguments passed to f(x) {x}. Want [x], got []"},
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
			logErrorWithPosition(t, errObj.Message, tt.expected)
		case int:
			testNumberObject(t, evaluated, float64(expected))
		default:
			t.Fatalf("unhandled type %T", expected)
		}
	}
}

func TestDecorators(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"@decorator f hello() {}", "function 'decorator' is not defined (used as decorator)"},
		{"a = 1; @a f hello(){}", "decorator 'a' must be a function, NUMBER given"},
		{"f decorator(fn) { return f () { return fn() * 2 } }; @decorator f test() { return 1 }; test()", 2},
		{"f decorator(fn) { return f () { return fn(...) * 2 } }; @decorator f test(x) { return x }; test(1)", 2},
		{"f decorator(fn, multiplier) { return f () { return fn() * multiplier } }; @decorator(4) f test() { return 1 }; test()", 4},
		{"f decorator(fn, multiplier) { return f () { return fn(...) * multiplier } }; @decorator(1000) f test(x) { return x }; test(1)", 1000},
		{"f decorator(fn, multiplier) { return f () { return fn(...) * multiplier } }; @decorator(2) @decorator(2) f test(x) { return x }; test(1)", 4},
		{"f decorator(fn, multiplier) { return f () { return fn(...) * multiplier } }; @decorator(2) @decorator(2) @decorator(2) f test(x) { return x }; test(1)", 8},
		{"f multiply(fn, multiplier) { return f () { return fn(...) * multiplier } }; f divide(fn, div) { return f () { return fn(...) / div } }; @multiply(10) @divide(5)  f test(x) { return x }; test(1)", 2},
		{"f decorator(fn) { return f () { return fn(...) } }; @decorator() @decorator_not_existing() f test(x) { return x }; test(1)", "function 'decorator_not_existing' is not defined (used as decorator)"},
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
			logErrorWithPosition(t, errObj.Message, tt.expected)
		case int:
			testNumberObject(t, evaluated, float64(expected))
		default:
			t.Fatalf("unhandled type %T", expected)
		}
	}
}

func TestCurrentArgs(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"f wrap_env() { return env(*args) }; wrap_env('x')", ""},
		{"env('x', '2'); f wrap_env() { return env(...) }; wrap_env('x').int()", 2},
		{"env('x', '2'); f wrap_env() { return env(...) }; wrap_env('x', '5'); wrap_env('x').int()", 5},
		{"f test() { sum(... + [1]) }; test(1,1,1,1)", 5},
		{"f argsummer() { x = 0; for i in ... {x += i }; return x }; f test() { argsummer(..., 20) }; test(10, 10)", 40},
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
			logErrorWithPosition(t, errObj.Message, tt.expected)
		case int:
			testNumberObject(t, evaluated, float64(expected))
		default:
			t.Fatalf("unhandled type %T", expected)
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
				if s.Value != tt.expected.(string) {
					t.Errorf("result is not the right string for '%s'. got='%s', want='%s'", tt.input, s.Value, tt.expected)
				}

				continue
			}

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			logErrorWithPosition(t, errObj.Message, tt.expected)
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
				if s.Value != tt.expected.(string) {
					t.Errorf("result is not the right string for '%s'. got='%s', want='%s'", tt.input, s.Value, tt.expected)
				}

				continue
			}

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			logErrorWithPosition(t, errObj.Message, tt.expected)
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

func TestInExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`1 in [1]`, true},
		{`1 in []`, false},
		{`!(1 in [])`, true},
		{`1 in ["1"]`, false},
		{`"1" in [1]`, false},
		{`1 in [1, 2]`, true},
		{`"hello" in [1, 2]`, false},
		{`"str" in "string"`, true},
		{`"xyz" in "string"`, false},
		{`"abc" in ["abc", "def"]`, true},
		{`"xyz" in ["abc", "def"]`, false},
		{`"x" in {"x": 0}`, true},
		{`"y" in {"x": 0}`, false},
		{`"y" in 12`, "'in' operator not supported on NUMBER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case bool:
			testBooleanObject(t, evaluated, bool(expected))
		default:
			errObj, ok := evaluated.(*object.Error)

			if !ok {
				t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
				continue
			}
			logErrorWithPosition(t, errObj.Message, expected)
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
		{"a = $(echo hello);\na.ok", true},
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
				if s.Value != tt.expected.(string) {
					t.Errorf("result is not the right string for '%s'. got='%s', want='%s'", tt.input, s.Value, tt.expected)
				}

				continue
			}

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			logErrorWithPosition(t, errObj.Message, tt.expected)
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
	type testLine struct {
		input    string
		expected interface{}
	}
	var tests []testLine
	if runtime.GOOS == "windows" {
		// cmd.exe commands
		tests = []testLine{
			{`a = "A"; b = "B"; $(echo $a$a$b$b$c$c)`, "AABB"},
			{`$(echo 123)`, "123"},
			{`$(echo hello world)`, "hello world"},
			{"a = 'A'; b = 'B'; `echo $a$a$b$b$c$c`", "AABB"},
			{"`echo 123`", "123"},
			{"`echo hello world`", "hello world"},
		}
	} else {
		executor := "bash"
		if runtime.GOOS == "windows" {
			executor = "cmd.exe"
		}

		// bash commands
		tests = []testLine{
			{`a = "A"; b = "B"; eee = "-e"; $(echo $eee -n $a$a$b$b$c$c)`, "AABB"},
			{`$(echo -n "123")`, "123"},
			{`$(echo -n hello world)`, "hello world"},
			{`$(echo hello world | xargs echo -n)`, "hello world"},
			{`$(echo \$0)`, executor},
			{`$(echo \$CONTEXT)`, "abs"},
			{"a = 'A'; b = 'B'; eee = '-e'; `echo $eee -n $a$a$b$b$c$c`", "AABB"},
			{"`echo -n '123'`", "123"},
			{"`echo -n hello world`", "hello world"},
			{"`echo hello world | xargs echo -n`", "hello world"},
			{"`echo \\$CONTEXT`", "abs"},
			{"`sleep 0.01`", ""},
			{"`sleep 0.01`.done", true},
			{"`sleep 0.01`.ok", true},
			{"`sleep 0.01 &`", ""},
			{"`sleep 0.01 &`.done", false},
			{"`sleep 0.01 &`.ok", false},
			{"`sleep 0.01 &`.wait().ok", true},
			{"`sleep 0.01 && echo 123 &`.wait()", "123"},
			{"`sleep 0.01 && echo 123 &`.kill()", ""},
			{"`sleep 0.01 && echo 123 &`.kill().done", true},
			{"`sleep 0.01 && echo 123 &`.kill().ok", false},
			{"`echo 123; sleep 10 &`.ok", false},
			{"`echo 123; sleep 10 &`.kill().done", true},
			{"`echo 123; sleep 10 &`.kill().ok", false},
		}
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
				t.Errorf("result is not the right string for '%s'. got='%s', want='%s'", tt.input, stringObj.Value, expected)
			}
		case bool:
			booleanObj, ok := evaluated.(*object.Boolean)
			if !ok {
				t.Errorf("object is not Boolean. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if booleanObj.Value != expected {
				t.Errorf("result is not the right boolean for '%s'. got='%t', want='%t'", tt.input, booleanObj.Value, expected)
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
			"[1, 2, 3][-2]",
			2,
		},
		{
			"[1, 2, 3][-10]",
			nil,
		},
		{
			"[1, 2, 3][-3]",
			1,
		},
		{
			"[1, 2, 3][-4]",
			nil,
		},
		{
			"[1, 2, 3][-0]",
			1,
		},
		{
			"a = [1, 2, 3, 4, 5, 6, 7, 8, 9][1:-300]; a[0]",
			nil,
		},
		{
			"a = [1, 2, 3, 4, 5, 6, 7, 8, 9][1:4]; a[0] + a[1] + a[2]",
			9,
		},
		{
			"a = [1, 2, 3, 4, 5, 6, 7, 8, 9][1:4]; a.len()",
			3,
		},
		{
			"a = [1, 2, 3, 4, 5, 6, 7, 8, 9][200:3]; a[0]",
			nil,
		},
		{
			"a = [1, 2, 3, 4, 5, 6, 7, 8, 9][7:-1]; a[0]",
			8,
		},
		{
			"a = [1, 2, 3, 4, 5, 6, 7, 8, 9][100:]; a[0]",
			nil,
		},
		{
			"a = [1, 2, 3, 4, 5, 6, 7, 8, 9][0:100]; a[0]",
			1,
		},
		{
			"a = [1, 2, 3, 4, 5, 6, 7, 8, 9][-10:]; a[0]",
			1,
		},
		{
			`a = [0,1,2,3,4,5][2:5]; len(a)`,
			3,
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

func TestOptionalChaining(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`a = null; a?.b?.c`,
			nil,
		},
		{
			`a = 1; a?.b?.c`,
			nil,
		},
		{
			`a = 1; a?.b?.c()`,
			nil,
		},
		{
			`a = {"b" : {"c": 1}}; a?.b?.c`,
			1,
		},
		{
			`a = {"b": 1}; a.b`,
			1,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch evaluated.(type) {
		case *object.Number:
			testNumberObject(t, evaluated, float64(tt.expected.(int)))
		default:
			testNullObject(t, evaluated)
		}
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
		{
			`a = {"fn": null}; a.fn`,
			nil,
		},
		{
			`a = {"fn": f() { return 1 }}; a.fn()`,
			1,
		},
		{
			`a = {"fn": f(x, y) { return y * x }}; a.fn(5, 3)`,
			15,
		},
		{
			`a = {}; a.fn()`,
			`HASH does not have method 'fn()'`,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch result := evaluated.(type) {
		case *object.Number:
			testNumberObject(t, evaluated, float64(tt.expected.(int)))
		case *object.Error:
			logErrorWithPosition(t, result.Message, tt.expected)
		default:
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
			"",
		},
		{
			`"123"[1]`,
			"2",
		},
		{
			`"123"[1:]`,
			"23",
		},
		{
			`"123"[1:1]`,
			"",
		},
		{
			`"123"[:2]`,
			"12",
		},
		{
			`"123"[:-1]`,
			"12",
		},
		{
			`"123"[-2]`,
			"2",
		},
		{
			`"123"[-1]`,
			"3",
		},
		{
			`"123"[-10]`,
			"",
		},
		{
			`"123"[2:-10]`,
			"",
		},
		{
			`"123"[2:1]`,
			"",
		},
		{
			`"123"[200:]`,
			"",
		},
		{
			`"123"[0:10]`,
			"123",
		},
		{
			`"123"[-10:]`,
			"123",
		},
		{
			`"123"[-10:{}]`,
			`index ranges can only be numerical: got "{}" (type HASH)`,
		},
		{
			`"123"[3]`,
			"",
		},
		{
			`"123"[0]`,
			"1",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch result := evaluated.(type) {
		case *object.String:
			testStringObject(t, evaluated, tt.expected.(string))
		case *object.Error:
			logErrorWithPosition(t, result.Message, tt.expected)
		default:
			t.Errorf("object is not the right result. got=%s ('%+v' expected)", result.Inspect(), tt.expected)
		}
	}
}

func testEval(input string) object.Object {
	env := object.NewEnvironment(os.Stdout, "")
	lex := lexer.New(input)
	p := parser.New(lex)
	program := p.ParseProgram()

	return BeginEval(program, env, lex)
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
	_, ok := obj.(*object.Null)
	if !ok {
		t.Errorf("object is not Null. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestEvalStringSpecialChars(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		s = "a\nb\nc"
		s
		`,
			`a
b
c`,
		},
		{`
		s = "a\tb\tc"
		s
		`, `a	b	c`,
		},
		{`
		s = fmt("a\\nb\\nc\n%s\n", "x\ny\nz")
		s
		`, `a\\nb\\nc
x
y
z
`,
		},
		{`
		a = split("a\nb\nc", "\n")
		str(a)
		`, `["a", "b", "c"]`,
		},
		{`
		a = split("a\nb\nc", "\n")
		s = join(a, "\n")
		s
		`, `a
b
c`,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestEvalAssignIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		a = [1, 2, 3, 4]
		a[0] = 99
		a[1] += 10
		a += [88]
		a[2] = "string"
		a[6] = 66
		a[5] = 55
		str(a)
		`,
			`[99, 12, "string", 4, 88, 55, 66]`,
		},
		{`
		h = {"a": 1, "b": 2, "c": 3}
		h["a"] = 99
		h["a"] += 1
		h += {"c": 33, "d": 44, "e": 55}
		h["z"] = {"x": 10, "y": 20}
		h["1.23"] = "string"
		h.d = 99
		h.d += 1
		h.z.x = 66
		h.f = 1.23
		str(h)
		`, `{"1.23": "string", "a": 100, "b": 2, "c": 33, "d": 100, "e": 55, "f": 1.23, "z": {"x": 66, "y": 20}}`,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestHashFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		h = {"a": 1, "b": 2, "c": {"x": 10, "y":20}, "e": "\"test"}
		hk = h.keys()
		hk = keys(h)
		hv = h.values()
		hv = values(h)
		hi = h.items()
		hi = items(h)
		hp = h.pop("a")
		hp = pop(h, "c")
		hp = h.pop("d")
		str(h)
		`, `{"b": 2, "e": "\"test"}`,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}
