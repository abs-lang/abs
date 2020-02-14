# Decorator

Decorators are a feature built on top of
ABS' functions -- they're not a type *per se*.

A decorator is a function that "wraps" another
function, allowing you to enhance the original
function's functionality with the decorator's
one.

An example could be a decorator that logs how
long a function takes to execute, or delays
execution altogether.

## Declaring decorators

A decorator is a plain-old function that
accepts `1 + N` arguments, where `1` is the
function being wrapped, and returns a new
function that wraps the original one:

```py
f log_if_slow(original_fn, treshold_ms) {
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
```

That's as simple as that: a named function 
that returns a new function that executes the
decorated one (`original_fn`) and returns its
result, while logging if it takes longer than
a few milliseconds.

## Using decorators

Now that we've declared our decorator, it's time
to use it, through the `@` notation:

```py
@log_if_slow(500)
f return_random_number_after_sleeping(seconds) {
    `sleep $seconds`
    return rand(1000)
}
```

and we can test our decorator has takn the stage:

```console
⧐  return_random_number_after_sleeping(0)
493
⧐  return_random_number_after_sleeping(1)
mmm, we were pretty slow...
371
```

Decorators are heavily inspired by [Python](https://www.python.org/dev/peps/pep-0318/).

## Next

That's about it for this section!

You can now head over to read a little bit about [how to install 3rd party libraries](/misc/3pl).