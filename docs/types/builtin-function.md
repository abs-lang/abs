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

### exit(code)

Exists the script with status `code`:

``` bash
exit(99)
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

Returns the `'/fully/expanded/path'` to the new current working directory and `path.ok`. If `path.ok` is `false`, then the `path`
is a `null string` and the current working directory is unchanged.
This supports compound tests such as ``cd() && `ls` ``.

``` bash
path = cd()
path.ok     # true
path        # /home/user or C:\Users\user

here = pwd()
path = cd("/path/to/nowhere")
path.ok         # false
path            # null string
here == pwd()   # true

cd("~/git/abs") # /home/user/git/abs or C:\Users\user\git\abs

cd("..")        # /home/user/git or C:\Users\user\git

cd("/usr/local/bin") # /usr/local/bin

dirs = cd() && `ls`.lines()
len(dirs)   # number of directories in homeDir

dirs = cd("/path/to/nowhere") && `ls`.lines()
len(dirs)   # 0
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

## Next

That's about it for this section!

You can now head over to read a little bit about [errors](/misc/error).