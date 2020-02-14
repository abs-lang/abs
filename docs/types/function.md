# Function

Functions are a very important element of ABS,
as they are the core of userland customizations.

A function is declared with the following syntax:

``` bash
f(x, y) {
    x + y
}
```

As you might notice, the return statement is implicit.
You can make it explicit, but we advise not to, in order
to keep your code as concise as possible:

``` bash
f(x, y) {
    return x + y
}
```

Most languages use a more "explicit" identifier for
functions (such as `function` or `func`), but ABS
favors `f` for 2 main reasons:

* brevity
* resembles the standard mathematical notation everyone is used to (*x ↦ f(x)*)

Functions must be called with the right number of arguments:

``` bash
fn = f(x) { x }
fn()
# ERROR: Wrong number of arguments passed to f(x) {
# x
# }. Want [x], got []
```

They can be passed as arguments to other functions:

``` bash
[1, 2, 3].map(f(x){ x + 1}) # [2, 3, 4]
```

and they can be assigned to variables as well:

``` bash
func = f(x){ x + 1}
[1, 2, 3].map(func) # [2, 3, 4]
```

Scoping is an important topic to cover when dealing with
functions:

``` bash
a = 10
func = f(x){ x + a }

f(1) # 11
a = 20
f(1) # 21
```

ABS supports closures just like mainstream languages:

``` bash
func = f(x) {
    f(y) {
        x + 1
    }
}

# can also be expressed as

func = f(x) {
    return f(y) {
        return x + 1
    }
}
```

## Named functions

You can create named functions by specifying an identifier
after the `f` keyword:

``` bash
f greeter(name) {
    echo("Hello $name!")
}

greeter(`whoami`) # "Hello root!"
```

As an alternative, you can manually assign
a function declaration to a variable, though
this is not the recommended approach:

``` bash
greeter = f (name) {
    echo("Hello $name!")
}

greeter(`whoami`) # "Hello root!"
```

Named functions are the basis of [decorators](/types/decorators).

## Accessing function arguments

Functions can receive a dynamic number of arguments,
and arguments can be "packed" through the special
`...` variable:

```py
f sum_numbers() {
    s = 0
    for x in ... {
        s += x
    }

    return s
}

sum_numbers(1) # 1
sum_numbers(1, 2, 3) # 6
```

`...` is a special variable that acts
like an array, so you can loop and slice
it however you want:

```py
f first_arg() {
    if ....len() > 0 {
        return ...[0]
    }

    return "No first arg"
}

first_arg() # "No first arg"
first_arg(1) # 1
```

When you pass `...` directly to a function,
it will be unpacked:

```py
f echo_wrapper() {
    echo(...)
}

echo_wrapper("hello %s", "root") # "hello root"
```

and you can add additional arguments as well:

```py
f echo_wrapper() {
    echo(..., "root")
}

echo_wrapper("hello %s %s", "sir") # "hello sir root"
```

## Supported functions

### str()

Returns the string representation of the function:

``` bash
f(x){}.str()
# f(x) {
#
# }
```

## Next

That's about it for this section!

You can now head over to read about [builtin functions](/types/builtin-function).