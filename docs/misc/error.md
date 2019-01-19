# Errors

When using ABS, you might bump into errors within your code. When the interpreter finds an error, it will give up trying to execute the script and will exit with status code `99`.

Note that there are 2 phases of the interpreter: parser and evaluator.

When the parser phase encounters a syntax error it will continue to process the rest of the file and report all of the syntax errors it finds. 

However, the evaluator phase will quit immediately when it encounters an evaluation error. Thus, you may need to run the ABS interpreter multiple times to find all the run-time errors.

When you are running ABS interactively in the Run, Eval, Print Loop (REPL) the location of the error can only be the current line you just entered.

However, when you are running ABS over a script file, even a small one, locating errors requires more help from the interpreter. ABS now provides `[line:column]` positions as well as the error line itself following the error message. Be aware that the column position is approximate because the parser does not always accurately pin-point the location of the offending token.

For example, a file with syntax errors might look like this when the first syntax error is in line 4 somewhere around column 11.
```
$ cat examples/error-parse.abs
# there are multiple parser errors in this file

# this is a malformed identifier
m.a = 'abc'

# this is a command terminated with a semi
d/d = $(command);

# this is a command terminated with a LF
c/c = $(command)

# this is a bad infix operator
b %% c

$ abs examples/error-parse.abs
 parser errors:
	no prefix parse function for '=' found
	[4:11]	m.a = 'abc'
	no prefix parse function for '=' found
	[7:16]	d/d = $(command);
	no prefix parse function for '=' found
	[10:16]	c/c = $(command)
	no prefix parse function for '%' found
	[13:6]	b %% c

$ echo $?
99
```
Furthermore, a file with evaluation errors might look like this when the first error encountered is in line 2 somewhere around column 11:
```
$ cat examples/error-eval.abs
# there is an evaluation error on line 2
1 + "hello"
echo("should not reach here")

$ abs examples/error-eval.abs
ERROR: type mismatch: NUMBER + STRING
	[2:11]	1 + "hello"

$ echo $?
99
```

Also, you can try running the interpreter error location tests by running this bash script. 
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

## Next

That's about it for this section!

You can now head over to read [a few more technical details about ABS ](/misc/technical-details).