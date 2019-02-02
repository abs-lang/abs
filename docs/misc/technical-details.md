# A few technical details...

The ABS interpreter is built with Golang version `1.11`, and is mostly based on [the interpreter book](https://interpreterbook.com/) written by [Thorsten Ball](https://twitter.com/thorstenball).

ABS is extremely different from Monkey, the "fictional" language the reader builds throughout the book, but the base structure (lexer, parser, evaluator) are still very much based on Thorsten's work.

## Why Go?

There are multiple reasons Go is the ideal choice for ABS, in no
particular order:

* portability, as our goal is to be able to deliver ABS to
multiple platforms without any hassle
* performance, as ABS' itself is not big on performance: having the
interpreter based on a fast platform allows us to recover
something back
* "strict" language, suitable for the purpose of making sure
syntax / parser errors are easily caught
* rich standard library, which allows to ship most of the ABS'
interpreter without relying on many external dependencies

## Development & contributing

Please see [github.com/abs-lang/abs/blob/master/CONTRIBUTING.md](https://github.com/abs-lang/abs/blob/master/CONTRIBUTING.md)

## Testing

### Interpreter Error Location

You can run the interpreter error location tests by invoking this bash script: `tests/test-abs.sh`. This script iterates over the `test-parser.abs` and `test-eval.abs` test scripts.
```bash
$ tests/test-abs.sh
=======================================
Test Parser
tests/test-parser.abs
 parser errors:
	no prefix parse function for '=' found
	[4:5]	m.a = 'abc'
	no prefix parse function for '=' found
	[7:5]	d/d = $(command);
	no prefix parse function for '=' found
	[10:5]	c/c = $(command)
	no prefix parse function for '%' found
	[13:4]	b %% c
	no prefix parse function for '&&' found
	[22:1]	&&||!-/*5;
	no prefix parse function for '||' found
	[22:3]	&&||!-/*5;
	no prefix parse function for '/' found
	[22:7]	&&||!-/*5;
	no prefix parse function for '<=>' found
	[25:2]	<=>
	expected next token to be NUMBER, got , instead
	[44:2]	[1, 2];
	no prefix parse function for ',' found
	[44:3]	[1, 2];
	no prefix parse function for ']' found
	[44:6]	[1, 2];
	no prefix parse function for '%' found
	[68:2]	~%
	no prefix parse function for '-=' found
	[70:1]	-=
	no prefix parse function for '/=' found
	[72:1]	/=
	no prefix parse function for '%=' found
	[74:1]	%=
	no prefix parse function for '^' found
	[79:2]	&^>><<
	no prefix parse function for '<<' found
	[79:5]	&^>><<
	Illegal token '$111'
	[80:1]	$111
	no prefix parse function for '$111' found
	[80:1]	$111
Exit code: 99

=======================================
Test Eval()
tests/test-eval.abs
ERROR: type mismatch: STRING + NUMBER
	[8:11]	    s = s + 1   # this is a comment
Exit code: 99

ERROR: invalid property 'junk' on type ARRAY
	[14:6]	    a.junk
Exit code: 99

ERROR: index operator not supported: f(x) {x} on HASH
	[19:20]	    {"name": "Abs"}[f(x) {x}];  
Exit code: 99
```

### String tests

String handling tests can be run from `abs tests/test-strings.abs`:

```bash
$ abs tests/test-strings.abs
=====================
>>> Testing string with expanded LFs:
echo("a\nb\nc")
a
b
c
=====================
>>> Testing string with expanded TABs:
echo("a\tb\tc")
a	b	c
=====================
>>> Testing string with expanded CRs:
echo("a\rb\rc")
c
=====================
>>> Testing string with mixed expanded LFs and escaped LFs:
echo("a\\nb\\nc\n%s\n", "x\ny\nz")
a\\nb\\nc
x
y
z

=====================
>>> Testing string with multiple escapes:
echo("hel\\\\lo")
hel\\\\lo
=====================
>>> Testing split and join strings with expanded LFs:
s = split("a\nb\nc", "\n")
echo(s)
[a, b, c]
ss = join(s, "\n")
echo(ss)
a
b
c
=====================
>>> Testing split and join strings with literal LFs:
s = split('a\nb\nc', '\n')
echo(s)
[a, b, c]
ss = join(s, '\n')
echo(ss)
a\nb\nc
```

### Array and hash assignment tests 

Array and hash assignment tests can be run from `abs tests/test-assign-index.abs `:

```bash
=====================
Test assignment to array indexed objects
>>> a = [1, 2, 3, 4]
[1, 2, 3, 4]
>>> a[0] = 99
[99, 2, 3, 4]
>>> a[1] += 10
[99, 12, 3, 4]
>>> a += [88]
[99, 12, 3, 4, 88]
>>> a[2] = "string"
[99, 12, string, 4, 88]
>>> a[6] = 66
[99, 12, string, 4, 88, null, 66]
>>> a[5] = 55
[99, 12, string, 4, 88, 55, 66]
=====================
Test assignment to hash indexed objects
>>> h = {"a": 1, "b": 2, "c": 3}
{a: 1, b: 2, c: 3}
>>> h["a"] = 99
{a: 99, b: 2, c: 3}
>>> h["a"] += 1
{a: 100, b: 2, c: 3}
>>> h += {"c": 33, "d": 44, "e": 55}
{a: 100, b: 2, c: 33, d: 44, e: 55}
h["z"] = {"x": 10, "y": 20}
{a: 100, b: 2, c: 33, d: 44, e: 55, z: {x: 10, y: 20}}
h["1.23"] = "string"
{1.23: string, a: 100, b: 2, c: 33, d: 44, e: 55, z: {x: 10, y: 20}}
h.d = 99
{1.23: string, a: 100, b: 2, c: 33, d: 99, e: 55, z: {x: 10, y: 20}}
h.d += 1
{1.23: string, a: 100, b: 2, c: 33, d: 100, e: 55, z: {x: 10, y: 20}}
h.z.x = 66
{1.23: string, a: 100, b: 2, c: 33, d: 100, e: 55, z: {x: 66, y: 20}}
h.f = 88
{1.23: string, a: 100, b: 2, c: 33, d: 100, e: 55, f: 88, z: {x: 66, y: 20}}
=====================
Error: assign to non-hash property
s = "string"
s.ok = true
ERROR: can only assign to hash property, got STRING
	[66:2]	s.ok = true
=====================
Error: add number to null hash property
h.g += 1
ERROR: type mismatch: NULL + NUMBER
	[72:5]	h.g += 1
=====================
Error: add number to null hash element
>>> h["g"] += 1
ERROR: type mismatch: NULL + NUMBER
	[78:8]	h["g"] += 1
```

## Roadmap

We're currently working on [1.2](https://github.com/abs-lang/abs/milestone/9).

## Next

That's about it for this section!

You can now head over to read ABS' [credits](/misc/credits).