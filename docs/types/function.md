# Function

Functions are a very important element of ABS,
as they are the core of userland customizations.

A function is declared with the following syntax:

``` bash
f(x, y) {
    x + y
}
```

As you might notice, the return statement is implicity.
You can make it explicit, but we advise not to keep
your code as coincise as possible:

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

Functions can be passed as arguments to other functions:

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