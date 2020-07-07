package evaluator

import (
	"testing"

	"github.com/abs-lang/abs/object"
)

type tests struct {
	input    string
	expected interface{}
}

func TestRuntime(t *testing.T) {
	tests := []tests{
		{`require('@runtime').version`, "test_version"},
		{`require('@runtime').name`, "abs"},
	}

	testStdLib(tests, t)
}

func TestUtil(t *testing.T) {
	tests := []tests{
		{`memo = require('@util').memoize; @memo(1) f test() { return 1 }; test()`, 1},
		{`memo = require('@util').memoize; @memo(1) f test(n) { return 1 + n }; test(10)`, 11},
		{`memo = require('@util').memoize; @memo(1) f test(n, m) { return n + m }; test(10, 5)`, 15},
		{`memo = require('@util').memoize; x = {"y": 0}; @memo(0) f test() { x.y += 1 }; test(); x.y`, 1},
		{`memo = require('@util').memoize; x = {"y": 0}; @memo(0) f test() { x.y += 1 }; test(); test(); test(); x.y`, 3},
		{`memo = require('@util').memoize; x = {"y": 0}; @memo(10) f test() { x.y += 1 }; test(); test(); test(); x.y`, 1},
		{`memo = require('@util').memoize; x = {"y": 0}; @memo(10) f test() { x.y += 1 }; @memo(10) f test2() { x.y += 1 }; test(); test(); test2(); x.y`, 2},
		{`memo = require('@util').memoize; x = {"y": 0}; @memo(0.250) f test() { x.y += 1 }; test(); test(); test(); x.y`, 1},
		{`memo = require('@util').memoize; x = {"y": 0}; @memo(0.250) f test() { x.y += 1 }; test(); test(); sleep(251); test(); x.y`, 2},
	}

	testStdLib(tests, t)
}

func testStdLib(tests []tests, t *testing.T) {
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
