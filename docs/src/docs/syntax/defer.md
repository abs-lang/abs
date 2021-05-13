---
permalink: /syntax/defer
---

# Defer <Badge text="experimental" type="warning"/>

Sometimes it is very helpful to guarantee a certain function is executed
regardless of what code path we take: you can use the `defer` keyword for
this.

```py
echo(1)
defer echo(3)
echo(2)
# 1
# 2
# 3
```

When you schedule a function to be deferred, it will executed right at
the end of the current scope. A `defer` inside a function will then
execute at the end of that function itself:

```py
echo(1)
f fn() {
    defer echo(3)
    echo(2)
}
fn()
echo(4)
# 1
# 2
# 3
# 4
```

You can `defer` any callable: a function call, a method or even a system
command. This can be very helpful if you need to run a cleanup function
right before wrapping up with your code:

```sh
defer `rm my-file.txt`
"some text" > "my-file.txt"

...
...
"some other text" >> "my-file.txt"
```

In this case, you will be guaranteed to execute the command that removes
`my-file.txt` before the program closes.

Be aware that code that is deferred does not have access to the return value
of its scope, and will supress errors -- if a `defer` block messes up you're
not going to see any error. This behavior is experimental, but we would most
likely like to give this kind of control through [try...catch...finally](https://github.com/abs-lang/abs/issues/118).