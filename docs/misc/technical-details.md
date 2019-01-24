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

The best way to start developing *for* ABS is to clone the repository
and run a `make build`: this will build a docker container with all
the necessary dependencies for developing locally (for experienced
Gophers: you might want to skip this altogether as your environment
will probably work perfectly).

With `make run` you can get inside a container built for ABS'
development, and `make test` will run all tests.

## Testing
### Interpreter Error Location
You can run the interpreter error location tests by invoking this bash script: `tests/test-abs.sh`. This script iterates over the `test-parser.abs` and `test-eval.abs` test scripts.
```
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
String handling tests can be run from `abs tests/test-strings.abs`
```
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

## Roadmap

We're currently working on [1.0](https://github.com/abs-lang/abs/milestone/5).

## Next

That's about it for this section!

You can now head over to read ABS' [credits](/misc/credits).