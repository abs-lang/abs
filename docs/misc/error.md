# Errors

When using ABS, you might bump into errors within your code. When the interpreter finds an error, it will give up trying to execute the script and will exit with status code `99`.

Note that there are 2 phases of the interpreter: parser and evaluator.

When the parser phase encounters a syntax error it will continue to process the rest of the file and report all of the syntax errors it finds. 

However, the evaluator phase will quit immediately when it encounters an evaluation error. Thus, you may need to run the ABS interpreter multiple times to find all the run-time errors.

When you are running ABS interactively in the Run, Eval, Print Loop (REPL) the location of the error can only be the current line you just entered.

However, when you are running ABS over a script file (even a small one) locating errors requires more help from the interpreter. ABS now provides `[line:column]` positions as well as the error line itself following the error message.

For example, a file with syntax errors might look like this when the first syntax error is in line 4 at column 5.
```
$ cat examples/error-parse.abs
# there are multiple parser errors in this file

# this is a malformed identifier
m.a = 'abc'

# this is a command terminated with a semi
d/d = `command`;

# this is a command terminated with a LF
c/c = `command`

# this is a bad infix operator
b %% c

$ abs examples/error-parse.abs
  parser errors:
	no prefix parse function for '=' found
	[4:5]	m.a = 'abc'
	no prefix parse function for '=' found
	[7:5]	d/d = `command`;
	no prefix parse function for '=' found
	[10:5]	c/c = `command`
	no prefix parse function for '%' found
	[13:4]	b %% c

$ echo $?
99
```
Furthermore, a file with evaluation errors might look like this when the first error encountered is in line 2 at column 3:
```
$ cat examples/error-eval.abs
# there is an evaluation error on line 2
1 + "hello"
echo("should not reach here")

$ abs examples/error-eval.abs
ERROR: type mismatch: NUMBER + STRING
	[2:3]	1 + "hello"

$ echo $?
99
```

## Next

That's about it for this section!

You can now head over to read [how to configure the REPL](/misc/configuring-the-repl).