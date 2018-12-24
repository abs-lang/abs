<p align="center">
  <a href="https://www.abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/abs-lang/abs/blob/master/bin/abs-horizontal.png?raw=true">
  </a>
</p>

# Builtin function

There are many builtin functions in ABS.
Take `type`, for example:

``` bash
type(1) # INTEGER
type([]) # ARRAY
```

We'll reveal you a secret now: all string, array, integer, hash functions
are actually "generic", but the syntax you see makes you think those are
specific to the string, integer, etc object.

The trick is very simple; whenever the ABS' interpreter seem a method call
such as `object.func(arg)` it will actually translate it to `func(object, arg)`.

Don't believe us? Try with these examples:

``` bash
map(["1"], int) # [1]
cmd = $(date)
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

### exit(code)

Exists the script with status `code`:

``` bash
exit(99)
```

### rand(max)

Returns a random number between 0 and `max`:

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

## Next

That's about it for this section!

You can now head over to read a little bit about [errors](/misc/error).