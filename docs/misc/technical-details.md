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
You can run the interpreter error location tests by invoking this bash script: `tests/test-abs.sh`.
```
$ tests/test-abs.sh

=======================================
Test Parser
tests/test-parser.abs
 parser errors:
	no prefix parse function for '=' found
	[4:11]	m.a = 'abc'
	no prefix parse function for '=' found
	[7:16]	d/d = $(command);
	no prefix parse function for '=' found
	[10:16]	c/c = $(command)
	no prefix parse function for '%' found
	[13:6]	b %% c
	no prefix parse function for '&&' found
	[22:4]	&&||!-/*5;
	no prefix parse function for 'OR' found
	[22:5]	&&||!-/*5;
	no prefix parse function for '/' found
	[22:8]	&&||!-/*5;
	no prefix parse function for '<=>' found
	[25:3]	<=>
	expected next token to be ], got , instead
	[44:3]	[1, 2];
	no prefix parse function for ',' found
	[44:5]	[1, 2];
	no prefix parse function for ']' found
	[44:7]	[1, 2];
	no prefix parse function for '%' found
	[68:2]	~%
	no prefix parse function for '-=' found
	[70:2]	-=
	no prefix parse function for '/=' found
	[72:2]	/=
	no prefix parse function for '%=' found
	[74:2]	%=
	no prefix parse function for '^' found
	[79:4]	&^>><<
	no prefix parse function for '<<' found
	[79:6]	&^>><<
	Illegal token '$111'
	[80:4]	$111
	no prefix parse function for 'ILLEGAL' found
	[76:7]	1.str()
Exit code: 99

=======================================
Test Eval()
tests/test-eval.abs
ERROR: type mismatch: STRING + NUMBER
	[8:35]	    s = s + 1   # this is a comment
Exit code: 99

ERROR: invalid property 'junk' on type ARRAY
	[14:6]	    a.junk
Exit code: 99

ERROR: index operator not supported: f(x) {x} on HASH
	[19:29]	    {"name": "Abs"}[f(x) {x}];  
Exit code: 99

```
## Roadmap

We're currently working on the [preview-4](https://github.com/abs-lang/abs/milestone/7).

## Next

That's about it for this section!

You can now head over to read ABS' [credits](/misc/credits).