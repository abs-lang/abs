# Builtin function

There are many builtin functions in ABS.
Take `type`, for example:

``` bash
type(1) # NUMBER
type([]) # ARRAY
```

We'll reveal you a secret now: all string, array, number & hash functions
are actually "generic", but the syntax you see makes you think those are
specific to the string, number, etc object.

The trick is very simple; whenever the ABS' interpreter finds a method call
such as `object.func(arg)` it will actually translate it to `func(object, arg)`.

Don't believe us? Try with these examples:

``` bash
map(["1"], int) # [1]
sort([3, 2, 1]) # [1, 2, 3]
len("abc") # 3
```

At the same time, there are some builtin functions that doesn't really
make sense to call with the method notation, so we've kept them in a
"special" location in the documentation. `exit(99)`, for example, exits
the program with the status code `99`, but it would definitely look
strange to see something such as `99.exit()`.

## Generic builtin functions

### echo(var)

Prints the given variable:

``` bash
echo("hello world")
```

You can use use placeholders in your strings:

``` bash
echo("hello %s", "world")
```

### stdin()

Reads from the `stdin`:

``` bash
echo("What do you like?")
echo("Oh, you like %s!", stdin()) # This line will block until user enters some text
```

Worth to note that you can read
the `stdin` indefinitely with:

``` bash
# Will read all input to the
# stdin and output it back
for input in stdin {
    echo(input)
}

# Or from the REPL:

⧐  for input in stdin { echo((input.int() / 2).str() + "...try again:")  }
10
5...try again:
5
2.5...try again:

...
```

### exit(code [, message])

Exits the script with status `code`:

``` bash
exit(99)
```

You can specify a message that's going to be outputted right
before exiting:

``` bash
⧐  exit(99, "Got problems...")
Got problems...%
```

### rand(max)

Returns a random integer number between 0 and `max`:

``` bash
rand(10) # 7
```

### env(str)

Returns the `str` environment variable:

``` bash
env("PATH") # "/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
```

### eval(str)

Evaluates the `str` as ABS code:

``` bash
eval("1 + 1") # 2
eval('object = {"x": 10}; object.x') # 10
```

### arg(n)

Returns the `n`th argument to the current script:

``` bash
arg(0) # /usr/bin/abs
```

### type(var)

Returns the type if the given variable:

``` bash
type("") # "STRING"
type({}) # "HASH"
```

### cd() or cd(path)

Sets the current working directory to `homeDir` or the given `path`
in both Linux and Windows.

Note that the path may have a `'~/'` prefix which will be replaced
with `'homeDir/'`. Also, in Windows, any `'/'` path separator will be
replaced with `'\'` and path names are not case-sensitive.

Returns the `'/fully/expanded/path'` to the new current working directory and `path.ok`.
If `path.ok` is `false`, that means there was an error changing directory:

``` bash
path = cd()
path.ok     # true
path        # /home/user or C:\Users\user

here = pwd()
path = cd("/path/to/nowhere")
path.ok         # false
path            # 'chdir /path/to/nowhere: no such file or directory'
here == pwd()   # true

cd("~/git/abs") # /home/user/git/abs or C:\Users\user\git\abs

cd("..")        # /home/user/git or C:\Users\user\git

cd("/usr/local/bin") # /usr/local/bin

dirs = cd() && `ls`.lines()
len(dirs)   # number of directories in homeDir
```

### pwd()

Returns the path to the current working directory -- equivalent
to `env("PWD")`. 

If executed from a script this will initially be the directory
containing the script.

To change the working directory, see `cd()`.

``` bash
pwd() # /go/src/github.com/abs-lang/abs
```

### flag(str)

Returns the value of a command-line flag. Both the `--flag` and `-flag`
form are accepted, and you can specify values with `--flag=x`
as well as `--flag x`:

``` bash
$ abs --test --test2 2 --test3=3 --test4 -test5
Hello user, welcome to the ABS programming language!
Type 'quit' when you're done, 'help' if you get lost!
⧐  flag("test")
true
⧐  flag("test2")
2
⧐  flag("test3")
3
⧐  flag("test4")
true
⧐  flag("test5")
true
⧐  flag("test6")
⧐  
```

If a flag value is not set, it will default to `true`.
The value of a flag that does not exist is `NULL`.

In all other cases `flag(...)` returns the literal string
value of the flag:

``` bash
$ abs --number 10
Hello user, welcome to the ABS programming language!
Type 'quit' when you're done, 'help' if you get lost!
⧐  n = flag("number")
⧐  n
10
⧐  type(n)
STRING
```

### sleep(ms)

Halts the process for as many `ms` you specified:

``` bash
sleep(1000) # sleeps for 1 second
```

### require(path_to_file.abs)

Evaluates the script at `path_to_file.abs`, and makes
its return value available to the caller.

For example, suppose we have a `module.abs` file:

``` bash
adder = f(a, b) { a + b }
multiplier = f(a, b) { a * b }

return {"adder": adder, "multiplier": multiplier}
```

and a `main.abs` such as:

``` bash
mod = require("module.abs")

echo(mod.adder(1, 2)) # 3
```

This is mostly useful to create external library
functions, like NPM modules or PIP packages, that
do not have access to the global environment. Any
variable set outside of the module will not be
available inside it, and vice-versa. The only
variable available to the caller (the script requiring
the module) is the module's return value.

Note that `require` uses paths that are relative to
the current script. Say that you have 2 files (`a.abs` and `b.abs`)
in the `/tmp` folder, `a.abs` can `require("./b.abs")`
without having to specify the full path (eg. `require("/tmp/b.abs")`).

### source(path_to_file.abs)

Evaluates the script at `path_to_file.abs` in the context of the 
ABS global environment. The results of any expressions in the file
become available to other commands in the REPL command line or to other
scripts in the current script execution chain. 

This is very similar to `require`, but allows the module to access
and edit the global environment. Any variable set inside the module
will also be available outside of it.

This is most useful for creating library functions in a startup script,
or variables that can be used by many other scripts. Often these library functions
are loaded via the ABS Init File `~/.absrc` (see [ABS Init File](/introduction/how-to-run-abs-code)).

For example:
```bash
$ cat ~/abs/lib/library.abs
# Useful function library ~/abs/lib/library.abs
adder = f(n, i) { n + i }

$ cat ~/.absrc
# ABS init file ~/.absrc
source("~/abs/lib/library.abs")

$ abs
Hello user, welcome to the ABS (1.9.0) programming language!
Type 'quit' when you are done, 'help' if you get lost!
⧐ adder(1, 2)
3
⧐
```

In addition to source file inclusion in scripts, you can also use
`source()` in the interactive REPL to load a script being
debugged. When the loaded script completes, the REPL command line
will have access to all variables and functions evaluated in the
script.

For example:
```bash
⧐  source("~/git/abs/tests/test-strings.abs")
...
=====================
>>> Testing split and join strings with expanded LFs:
s = split("a\nb\nc", "\n")
echo(s)
[a, b, c]
...
⧐  s
[a, b, c]
⧐ 
```

Note well that nested source files must not create a circular
inclusion condition. You can configure the intended source file
inclusion depth using the `ABS_SOURCE_DEPTH` OS or ABS environment
variables. The default is `ABS_SOURCE_DEPTH=10`. This will prevent
a panic in the ABS interpreter if there is an unintended circular
source inclusion.

For example an ABS Init File may contain:
```bash
ABS_SOURCE_DEPTH = 15
source("~/path/to/abs/lib")
```

This will limit the source inclusion depth to 15 levels for this
`source()` statement and will also apply to future `source()`
statements until changed.

## Next

That's about it for this section!

You can now head over to read a little bit about [how to install 3rd party libraries](/misc/3pl).