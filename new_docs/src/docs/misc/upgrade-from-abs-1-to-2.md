---
permalink: /misc/upgrade-from-abs-1-to-2
---

# Upgrading from ABS 1 to 2

It's not always possible to release backwards compatible changes,
and ABS is no exception to the rule: luckily, the major upgrade
between version 1 and 2 should be extremely painless as there have
been no syntax changes, but rather just a handful of function
changes.

## Deprecated functions

- the `slice` function you could use on arrays and strings has been removed. Use the index notation instead: `[1, 2, 3].slice(0, 1)` is equivalent to `[1, 2, 3][0:1]`
- the `contains` function you could use on arrays and strings has been removed. Use the `in` operator instead: `[1, 2, 3].contains(1)` is equivalent to `1 in [1, 2, 3]`

## Misc

The structure for decorators had to be slighly changed to allow
substantial improvements (ABS 2's decorators are 100% aligned with
Python's decorators which are extremely powerful).

Earlier, a decorator function would be declared as:

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

@log_if_slow(500)
f my_func() {
    ...
}
```

In ABS 2, a decorator must evaluate to a function that accepts
the original function and returns a new one with the "enhanced"
behaviour. It's probably easier to see it in action:

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

@log_if_slow(500)
f my_func() {
    ...
}
```

As you can see there are 2 main differences:

- the arguments to the decorator don't start with the original function anymore
- there's an additional wrapping function, accepting the original function, in the decorator

## What more?

That's really it: upgrading to ABS 2 should be an extremely painless process!

## Next

That's about it for this section!

You can now head over to read ABS' [credits](/misc/credits).
