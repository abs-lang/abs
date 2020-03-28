package evaluator

import (
	"testing"

	"github.com/abs-lang/abs/object"
)

type Tests struct {
	input    string
	expected interface{}
}

func TestUnique(t *testing.T) {
	tests := []Tests{
		{`[1,2,3,3,2,1].unique()`, []int{1, 2, 3}},
	}

	testBuiltinFunction(tests, t)
}

func TestMap(t *testing.T) {
	tests := []Tests{
		{`[1,2,"a"].map(int)`, "int(...) can only be called on strings which represent numbers, 'a' given"},
		{`[1].map(f(x) { y = x + 1 }).str()`, "[null]"},
		{`(0..99).map( f(i) { arg(i) } ).filter( f(i) { i != "" } ).len() == args().len()`, true},
	}

	testBuiltinFunction(tests, t)
}

func TestSum(t *testing.T) {
	tests := []Tests{
		{`[1, null].sum()`, "sum(...) can only be called on an homogeneous array, got [1, null]"},
		{`[null, null].sum()`, "sum(...) can only be called on arrays of numbers, got [null, null]"},
		{`[].sum()`, 0},
		{`[1, 2].sum()`, 3},
	}

	testBuiltinFunction(tests, t)
}

func TestArgs(t *testing.T) {
	tests := []Tests{
		{`arg("o")`, "argument 0 to arg(...) is not supported (got: o, allowed: NUMBER)"},
		{`arg(99)`, ""},
		{`arg(-1)`, ""},
		{`arg(0) == args()[0]`, true},
		{`arg(1) == args()[1]`, true},
		{`arg(2) == args()[2]`, true},
	}

	testBuiltinFunction(tests, t)
}

func TestIsNumber(t *testing.T) {
	tests := []Tests{
		{`is_number("aaa")`, false},
		{`is_number("123")`, true},
		{`is_number("123.33")`, true},
		{`is_number(123)`, true},
		{`is_number(123.33)`, true},
	}

	testBuiltinFunction(tests, t)
}

func TestSlice(t *testing.T) {
	tests := []Tests{
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
	}

	testBuiltinFunction(tests, t)
}

func TestType(t *testing.T) {
	tests := []Tests{
		{`type("SOME")`, "STRING"},
		{`type(1)`, "NUMBER"},
		{`type({})`, "HASH"},
		{`type([])`, "ARRAY"},
		{`type("{}".json())`, "HASH"},
		{`type(null)`, "NULL"},
	}

	testBuiltinFunction(tests, t)
}

func TestLen(t *testing.T) {
	tests := []Tests{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument 0 to len(...) is not supported (got: 1, allowed: STRING, ARRAY)"},
		{`len("one", "two")`, "wrong number of arguments to len(...): got=2, want=1"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
	}

	testBuiltinFunction(tests, t)
}

func TestInt(t *testing.T) {
	tests := []Tests{
		{`int("10")`, 10},
		{`int("10.5")`, 10},
		{`int("abc")`, `int(...) can only be called on strings which represent numbers, 'abc' given`},
		{`int([])`, "argument 0 to int(...) is not supported (got: [], allowed: NUMBER, STRING)"},
	}

	testBuiltinFunction(tests, t)
}

func TestFind(t *testing.T) {
	tests := []Tests{
		{`find([1,2,3,3], f(x) {x == 3})`, 3},
		{`find([1,2], f(x) {x == "some"})`, nil},
		{`find([{}, {}], f(x) {x.y == 1})`, nil},
		{`x = find([{}, {"y": 1, "z": 10}, {}], f(x) {x.y == 1}); x.z`, 10},
		{`x = find([{}, {"y": 1, "z": 10}, {}], {"y": 1}); x.z`, 10},
		{`x = find([{}, {"y": {}, "z": 10}, {}], {"y": {}}); x.z`, 10},
		{`find([{}, {"y": "1", "z": 10}, {}], {"y": 1})`, nil},
	}

	testBuiltinFunction(tests, t)
}

func TestPrefix(t *testing.T) {
	tests := []Tests{
		{`"a".prefix("b")`, false},
		{`"a".prefix("a")`, true},
	}

	testBuiltinFunction(tests, t)
}

func TestJson(t *testing.T) {
	tests := []Tests{
		{`"{\"a\": null}".json().a`, nil},
		{`"{\"k\": \"v\"}".json()["k"]`, "v"},
		{`''.json()`, ""},
		{`'         '.json()`, ""},
		{`"2".json()`, 2},
		{`'"2"'.json()`, "2"},
		{`'true'.json()`, true},
		{`'null'.json()`, nil},
		{`'"hello"'.json()`, "hello"},
		{`'[1, 2, 3]'.json()`, []int{1, 2, 3}},
		{`'"hello'.json()`, "argument to `json` must be a valid JSON object, got '\"hello'"},
	}

	testBuiltinFunction(tests, t)
}

func TestRand(t *testing.T) {
	tests := []Tests{
		{`rand(1)`, 0},
	}

	testBuiltinFunction(tests, t)
}

func TestSplit(t *testing.T) {
	tests := []Tests{
		{`split("a\"b\"c", "\"")`, []string{"a", "b", "c"}},
		{`split("a b c", " ")`, []string{"a", "b", "c"}},
		{`split("a b c")`, []string{"a", "b", "c"}},
	}

	testBuiltinFunction(tests, t)
}

func TestFmt(t *testing.T) {
	tests := []Tests{
		{`"hello %s".fmt("world")`, "hello world"},
		{`"hello %s".fmt()`, "hello %!s(MISSING)"},
		{`"hello %s".fmt(1)`, "hello 1"},
		{`"hello %s".fmt({})`, "hello {}"},
	}

	testBuiltinFunction(tests, t)
}

func TestReplace(t *testing.T) {
	tests := []Tests{
		{`"a".replace("a", "b", -1)`, "b"},
		{`"a".replace("a", "b")`, "b"},
		{`"ac".replace(["a", "c"], "b", -1)`, "bb"},
		{`"ac".replace(["a", "c"], "b")`, "bb"},
	}

	testBuiltinFunction(tests, t)
}

func TestCeil(t *testing.T) {
	tests := []Tests{
		{`1.ceil()`, 1},
		{`1.ceil()`, 1},
		{`1.23.ceil()`, 2},
		{`1.66.ceil()`, 2},
		{`"1.23".ceil()`, 2},
		{`"1.66".ceil()`, 2},
	}

	testBuiltinFunction(tests, t)
}

func TestFloor(t *testing.T) {
	tests := []Tests{
		{`1.floor()`, 1},
		{`1.floor()`, 1},
		{`1.23.floor()`, 1},
		{`1.66.floor()`, 1},
		{`"1.23".floor()`, 1},
		{`"1.66".floor()`, 1},
	}

	testBuiltinFunction(tests, t)
}

func TestRound(t *testing.T) {
	tests := []Tests{
		{`1.round()`, 1},
		{`1.round(2)`, 1.00},
		{`1.23.round(1)`, 1.2},
		{`1.66.round(1)`, 1.7},
		{`"1.23".round(1)`, 1.2},
		{`"1.66".round(1)`, 1.7},
	}

	testBuiltinFunction(tests, t)
}

func TestStr(t *testing.T) {
	tests := []Tests{
		{`"a".str()`, "a"},
		{`1.str()`, "1"},
		{`[1].str()`, "[1]"},
		{`{"a": 10}.str()`, `{"a": 10}`},
	}

	testBuiltinFunction(tests, t)
}

func TestTsv(t *testing.T) {
	tests := []Tests{
		{`[[1,2,3], [2,3,4]].tsv()`, "1\t2\t3\n2\t3\t4"},
		{`[1].tsv()`, "tsv() must be called on an array of arrays or objects, such as [[1, 2, 3], [4, 5, 6]], '[1]' given"},
		{`[{"c": 3, "b": "hello"}, {"b": 20, "c": 0}].tsv()`, "b\tc\nhello\t3\n20\t0"},
		{`[[1,2,3], [2,3,4]].tsv(",")`, "1,2,3\n2,3,4"},
		{`[[1,2,3], [2]].tsv(",")`, "1,2,3\n2"},
		{`[[1,2,3], [2,3,4]].tsv("abc")`, "1a2a3\n2a3a4"},
		{`[[1,2,3], [2,3,4]].tsv("")`, "the separator argument to the tsv() function needs to be a valid character, '' given"},
		{`[{"c": 3, "b": "hello"}, {"b": 20, "c": 0}].tsv("\t", ["c", "b", "a"])`, "c\tb\ta\n3\thello\tnull\n0\t20\tnull"},
	}

	testBuiltinFunction(tests, t)
}

func TestCall(t *testing.T) {
	tests := []Tests{
		{`adder = f (a, b) { return a + b }; adder.call([5, 5])`, 10},
		{`int.call(["12"])`, 12},
	}

	testBuiltinFunction(tests, t)
}

func TestNumber(t *testing.T) {
	tests := []Tests{
		{`number("aaa")`, "number(...) can only be called on strings which represent numbers, 'aaa' given"},
		{`number("10")`, 10},
		{`number("10.55")`, 10.55},
	}

	testBuiltinFunction(tests, t)
}

func TestEnv(t *testing.T) {
	tests := []Tests{
		{`env("CONTEXT")`, "abs"},
		{`env("FOO")`, ""},
		{`env("FOO", "bar")`, "bar"},
	}

	testBuiltinFunction(tests, t)
}

func TestFilter(t *testing.T) {
	tests := []Tests{
		{`[1,2,"a"].filter(int)`, "int(...) can only be called on strings which represent numbers, 'a' given"},
		{`[1,2,3].filter(f(x) {x == 1})`, []int{1}},
	}

	testBuiltinFunction(tests, t)
}

func TestEcho(t *testing.T) {
	tests := []Tests{
		{`echo("hello", "world!")`, nil},
	}

	testBuiltinFunction(tests, t)
}

func TestSort(t *testing.T) {
	tests := []Tests{
		{`[1, 2].sort()`, []int{1, 2}},
		{`["b", "a"].sort()`, []string{"a", "b"}},
		{`["b", 1].sort()`, `argument to 'sort' must be an homogeneous array (elements of the same type), got ["b", 1]`},
		{`[{}].sort()`, "cannot sort an array with given elements elements ([{}])"},
		{`[[]].sort()`, "cannot sort an array with given elements elements ([[]])"},
	}

	testBuiltinFunction(tests, t)
}

func TestSource(t *testing.T) {
	tests := []Tests{
		{`"a = 2; return 10" >> "test-ignore-source-vs-require.abs"; a = 1; x = source("test-ignore-source-vs-require.abs"); a`, 2},
		{`"a = 2; return 10" >> "test-ignore-source-vs-require.abs"; a = 1; x = source("test-ignore-source-vs-require.abs"); x`, 10},
	}

	testBuiltinFunction(tests, t)
}

func TestRequire(t *testing.T) {
	tests := []Tests{
		{`"a = 2; return 10" >> "test-ignore-source-vs-require.abs"; a = 1; x = require("test-ignore-source-vs-require.abs"); a`, 1},
		{`"a = 2; return 10" >> "test-ignore-source-vs-require.abs"; a = 1; x = require("test-ignore-source-vs-require.abs"); x`, 10},
	}

	testBuiltinFunction(tests, t)
}

func TestSleep(t *testing.T) {
	tests := []Tests{
		{`sleep(1000)`, nil},
		{`sleep(0.01)`, nil},
	}

	testBuiltinFunction(tests, t)
}

func TestSome(t *testing.T) {
	tests := []Tests{
		{`[1, 2].some(f(x) {x == 2})`, true},
		{`[].some(f(x) {x})`, false},
	}

	testBuiltinFunction(tests, t)
}

func TestEvery(t *testing.T) {
	tests := []Tests{
		{`[1, 2].every(f(x) { return x == 2 || x == 1})`, true},
		{`[].every(f(x) {x})`, true},
		{`[1,2,3].every(f(x) {x == 1})`, false},
	}

	testBuiltinFunction(tests, t)
}

func TestShift(t *testing.T) {
	tests := []Tests{
		{`[].shift()`, nil},
		{`[1, 2].shift()`, 1},
		{`a = [1, 2]; a.shift(); a`, []int{2}},
	}

	testBuiltinFunction(tests, t)
}

func TestReverse(t *testing.T) {
	tests := []Tests{
		{`[1, 2].reverse();`, []int{2, 1}},
		{`"abc".reverse();`, "cba"},
	}

	testBuiltinFunction(tests, t)
}

func TestPush(t *testing.T) {
	tests := []Tests{
		{`[1, 2].push("a");`, []interface{}{1, 2, "a"}},
	}

	testBuiltinFunction(tests, t)
}

func TestPop(t *testing.T) {
	tests := []Tests{
		{`[1, 2].pop();`, 2},
		{`a = [1, 2]; a.pop(); a`, []int{1}},
	}

	testBuiltinFunction(tests, t)
}

func TestKeys(t *testing.T) {
	tests := []Tests{
		{`[1, 2].keys()`, []int{0, 1}},
		{`{'a': 1}.keys()`, []string{"a"}},
	}

	testBuiltinFunction(tests, t)
}

func TestJoin(t *testing.T) {
	tests := []Tests{
		{`[1, 2].join("-")`, "1-2"},
		{`["a", "b"].join("-")`, "a-b"},
		{`["a", "b"].join()`, "ab"},
	}

	testBuiltinFunction(tests, t)
}

func TestAny(t *testing.T) {
	tests := []Tests{
		{`"a".any("b")`, false},
		{`"a".any("a")`, true},
	}

	testBuiltinFunction(tests, t)
}

func TestSuffix(t *testing.T) {
	tests := []Tests{
		{`"a".suffix("b")`, false},
		{`"a".suffix("a")`, true},
	}

	testBuiltinFunction(tests, t)
}

func TestIndex(t *testing.T) {
	tests := []Tests{
		{`"ab".index("b")`, 1},
		{`"a".index("b")`, nil},
	}

	testBuiltinFunction(tests, t)
}

func TestLastIndex(t *testing.T) {
	tests := []Tests{
		{`"abb".last_index("b")`, 2},
		{`"a".last_index("b")`, nil},
	}

	testBuiltinFunction(tests, t)
}

func TestRepeat(t *testing.T) {
	tests := []Tests{
		{`"a".repeat(3)`, "aaa"},
		{`"a".repeat(3)`, "aaa"},
	}

	testBuiltinFunction(tests, t)
}

func TestTitle(t *testing.T) {
	tests := []Tests{
		{`"a great movie".title()`, "A Great Movie"},
	}

	testBuiltinFunction(tests, t)
}

func TestLower(t *testing.T) {
	tests := []Tests{
		{`"A great movie".lower()`, "a great movie"},
	}

	testBuiltinFunction(tests, t)
}

func TestUpper(t *testing.T) {
	tests := []Tests{
		{`"A great movie".upper()`, "A GREAT MOVIE"},
	}

	testBuiltinFunction(tests, t)
}

func TestTrim(t *testing.T) {
	tests := []Tests{
		{`"  A great movie  ".trim()`, "A great movie"},
	}

	testBuiltinFunction(tests, t)
}

func TestTrimBy(t *testing.T) {
	tests := []Tests{
		{`"  A great movie  ".trim_by(" A")`, "great movie"},
	}

	testBuiltinFunction(tests, t)
}

func TestEval(t *testing.T) {
	tests := []Tests{
		{`a = 1; eval("a")`, 1},
	}

	testBuiltinFunction(tests, t)
}

func TestMisc(t *testing.T) {
	tests := []Tests{
		{`pwd().split("").reverse().slice(0, 33).reverse().join("").replace("\\", "/", -1).suffix("/evaluator")`, true}, // Little trick to get travis to run this test, as the base path is not /go/src/
		{`cwd = cd(); cwd == pwd()`, true},
		{`cwd = cd("path/to/nowhere"); cwd == pwd()`, false},
		{`lines("a
b
c")`, []string{"a", "b", "c"}},
		{`$()`, ""},
	}

	testBuiltinFunction(tests, t)
}

func TestChunk(t *testing.T) {
	tests := []Tests{
		{`chunk([1,2,3,4,5])`, "wrong number of arguments to chunk(...): got=1, want=2"},
		{`x = chunk([1,2,3,4,5], 2); len(x)`, 3},
		{`x = chunk([1,2,3,4,5], 2); x[0]`, []int{1, 2}},
		{`x = chunk([1,2,3,4,5], 2); x[1]`, []int{3, 4}},
		{`x = chunk([1,2,3,4,5], 2); x[2]`, []int{5}},
		{`x = chunk([1,2,3,4,5], 0);`, "argument to chunk must be a positive integer, got '0'"},
		{`x = chunk([1,2,3,4,5], -1);`, "argument to chunk must be a positive integer, got '-1'"},
		{`x = chunk([1,2,3,4,5], -1.5);`, "argument to chunk must be a positive integer, got '-1.5'"},
		{`x = chunk([1,2,3,4,5], 1.5);`, "argument to chunk must be a positive integer, got '1.5'"},
		{`x = chunk([], 10); len(x)`, 0},
		{`x = chunk([], 10); x`, []int{}},
	}

	testBuiltinFunction(tests, t)
}

func TestBetween(t *testing.T) {
	tests := []Tests{
		{`1.between(0, 2)`, true},
		{`1.between(0, 1.1)`, true},
		{`1.between(0, 0.9)`, false},
		{`1.between(1, 0)`, "arguments to between(min, max) must satisfy min < max (1 < 0 given)"},
		{`1.between(1, 2)`, true},
	}

	testBuiltinFunction(tests, t)
}

func TestClamp(t *testing.T) {
	tests := []Tests{
		{`2.clamp(0, 10)`, 2},
		{`2.clamp(2, 10)`, 2},
		{`2.clamp(3, 10)`, 3},
		{`2.clamp(0, 3)`, 2},
		{`2.clamp(2, 3)`, 2},
		{`2.clamp(3, 3)`, "arguments to clamp(min, max) must satisfy min < max"},
		{`2.clamp(3, 10)`, 3},
		{`2.clamp(0, 1)`, 1},
		{`2.clamp(0, 2)`, 2},
		{`2.clamp(1.5, 2.5)`, 2},
		{`2.clamp(2.1, 2.5)`, 2.1},
		{`2.5.clamp(2.1, 2.3)`, 2.3},
	}

	testBuiltinFunction(tests, t)
}

func TestCamel(t *testing.T) {
	tests := []Tests{
		{`"long cool woman in a black dress".camel()`, "longCoolWomanInABlackDress"},
		{`"  long cool woman in a black dress   ".camel()`, "longCoolWomanInABlackDress"},
		{`"  long cool woman in a_black dress   ".camel()`, "longCoolWomanInABlackDress"},
	}

	testBuiltinFunction(tests, t)
}

func TestSnake(t *testing.T) {
	tests := []Tests{
		{`"long cool woman in a black dress".snake()`, "long_cool_woman_in_a_black_dress"},
		{`"  long cool woman in a black dress   ".snake()`, "long_cool_woman_in_a_black_dress"},
		{`"  long cool woman in a_black dress   ".snake()`, "long_cool_woman_in_a_black_dress"},
	}

	testBuiltinFunction(tests, t)
}

func TestKebab(t *testing.T) {
	tests := []Tests{
		{`"long cool woman in a black dress".kebab()`, "long-cool-woman-in-a-black-dress"},
		{`"  long cool woman in a black dress   ".kebab()`, "long-cool-woman-in-a-black-dress"},
		{`"  long cool woman in a_black dress   ".kebab()`, "long-cool-woman-in-a-black-dress"},
	}

	testBuiltinFunction(tests, t)
}

func TestIntersect(t *testing.T) {
	tests := []Tests{
		{`[1,2,3].intersect([])`, []int{}},
		{`[1,2,3].intersect([3])`, []int{3}},
		{`[1,2,3].intersect([3, 1])`, []int{1, 3}},
		{`[1,2,3].intersect([1,2,3,4])`, []int{1, 2, 3}},
	}

	testBuiltinFunction(tests, t)
}

func TestDiff(t *testing.T) {
	tests := []Tests{
		{`[1,2,3].diff([])`, []int{1, 2, 3}},
		{`[1,2,3].diff([3])`, []int{1, 2}},
		{`[1,2,3].diff([3, 1])`, []int{2}},
		{`[1,2,3].diff([1,2,3,4])`, []int{}},
	}

	testBuiltinFunction(tests, t)
}

func TestDiffSymmetric(t *testing.T) {
	tests := []Tests{
		{`[1,2,3].diff_symmetric([])`, []int{1, 2, 3}},
		{`[1,2,3].diff_symmetric([3])`, []int{1, 2}},
		{`[1,2,3].diff_symmetric([3, 1])`, []int{2}},
		{`[1,2,3].diff_symmetric([1,2,3,4])`, []int{4}},
	}

	testBuiltinFunction(tests, t)
}

func TestUnion(t *testing.T) {
	tests := []Tests{
		{`[1, 2, 3].union([1, 2, 3, 4])`, []int{1, 2, 3, 4}},
		{`[1, 2, 3].union([3])`, []int{1, 2, 3}},
		{`[].union([3, 1])`, []int{3, 1}},
		{`[1, 2].union([3, 4])`, []int{1, 2, 3, 4}},
	}

	testBuiltinFunction(tests, t)
}

func TestFlatten(t *testing.T) {
	tests := []Tests{
		{`[1, 2, 3].flatten()`, []int{1, 2, 3}},
		{`[1, 2, [3]].flatten()`, []int{1, 2, 3}},
		{`[1, 2, [3, 4]].flatten()`, []int{1, 2, 3, 4}},
		{`[[1, 2], [3, 4]].flatten()`, []int{1, 2, 3, 4}},
	}

	testBuiltinFunction(tests, t)
}

func TestFlattenDeep(t *testing.T) {
	tests := []Tests{
		{`[1, 2, 3].flatten_deep()`, []int{1, 2, 3}},
		{`[1, 2, [3]].flatten_deep()`, []int{1, 2, 3}},
		{`[1, 2, [3, 4]].flatten_deep()`, []int{1, 2, 3, 4}},
		{`[[1, 2], [3, 4]].flatten_deep()`, []int{1, 2, 3, 4}},
		{`[[1, [2]], [3, 4]].flatten_deep()`, []int{1, 2, 3, 4}},
		{`[[[1, [2]], [3, 4]]].flatten_deep()`, []int{1, 2, 3, 4}},
	}

	testBuiltinFunction(tests, t)
}

func TestMax(t *testing.T) {
	tests := []Tests{
		{`[].max()`, nil},
		{`[-10].max()`, -10},
		{`[-10, 0, 100, 9].max()`, 100},
		{`[-10, 0, 100, 9, 100.1].max()`, 100.1},
		{`[-10, {}, 100, 9].max()`, "max(...) can only be called on an homogeneous array, got [-10, {}, 100, 9]"},
	}

	testBuiltinFunction(tests, t)
}

func TestMin(t *testing.T) {
	tests := []Tests{
		{`[].min()`, nil},
		{`[-10].min()`, -10},
		{`[-10, 0, 100, 9].min()`, -10},
		{`[-10, 0, 100, 9, -10.5].min()`, -10.5},
		{`[-10, {}, 100, 9].min()`, "min(...) can only be called on an homogeneous array, got [-10, {}, 100, 9]"},
	}

	testBuiltinFunction(tests, t)
}

func testBuiltinFunction(tests []Tests, t *testing.T) {
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
