---
permalink: /types/decorator
---

# Decorator

Decorators are a feature built on top of
ABS' functions -- they're not a type _per se_
but they do have their own _syntactic sugar_.

A decorator is a function that "wraps" another
function, allowing you to enhance the original
function's functionality with the decorator's
one.

An example could be a decorator that logs how
long a function takes to execute, or delays
execution altogether.

## Simple decorators

A decorator is a plain-old function that
accepts the original function and returns a new
function that wraps the original one with its
own behaviour. After defining it, you can
"decorate" other functions through the convenient
`@` syntax:

```py
f uppercase(fn) {
    return f() {
        return fn(...).upper()
    }
}

@uppercase
f stringer(x) {
    return x.str()
}

stringer({}) # "{}"
stringer(12) # "12"
stringer("hello") # "HELLO"
```

As you see, `stringer`'s behaviour has been altered:
it will now output uppercase strings.

## Decorators with arguments

As we've just seen, a decorator simply needs to
be a function that accepts the original
function and returns a new one, "enhancing"
the original behavior. If you wish to
configure decorators with arguments, it
is as simple as adding another level
of "wrapping":

```py
f log_if_slow(treshold_ms) {
    return f(original_fn) {
        return f() {
            start = `date +%s%3N`.int()
            res = original_fn(...)
            end = `date +%s%3N`.int()

            if end - start > treshold_ms {
                echo("mmm, we were pretty slow...")
            }

            return res
        }
    }
}
```

That's as simple as that: a named function
that returns a new function that executes the
decorated one (`original_fn`) and returns its
result, while logging if it takes longer than
a few milliseconds.

Now that we've declared our decorator, it's time
to use it, through the `@` notation:

```py
@log_if_slow(500)
f return_random_number_after_sleeping(seconds) {
    `sleep $seconds`
    return rand(1000)
}
```

and we can test our decorator has taken the stage:

```console
⧐  return_random_number_after_sleeping(0)
493
⧐  return_random_number_after_sleeping(1)
mmm, we were pretty slow...
371
```

Decorators are heavily inspired by [Python](https://www.python.org/dev/peps/pep-0318/) -- if you wish to understand
how they work more in depth we'd recommend reading this [primer on Python decorators](https://realpython.com/primer-on-python-decorators).