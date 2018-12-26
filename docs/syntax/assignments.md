<p align="center">
  <a href="https://www.abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/abs-lang/abs/blob/master/bin/abs-horizontal.png?raw=true">
  </a>
</p>

# Assignments

Just like about any other language, assignments are pretty
straightforward:

``` bash
x = "hello world"
```

Array destructuring is supported, meaning you can
set multiple variables based on an array:

``` bash
[x, y, z] = ["hello world", 99, {}]
x # "hello world"
y # 99
z # {}
```

If the number of variables you're trying to set is longer
than the array, the extra variables will be set to
null:

``` bash
[x] = []
x # null
```

Note that when using destructuring, the previous line
should end with a semicolon:

``` bash
# Ok
x = "10";
[a] = [100]

# Not ok, the parser reads this `x = "10"[a]`
x = "10"
[a] = [100]
```

This behavior will be fixed soon (see [#83](https://github.com/abs-lang/abs/issues/83)).

ABS doesn't have block-specific scopes, so any new variable
declared in a block is automatically available outside as well:

``` bash
if true {
    x = "hello world"
}

echo(x) # "hello world"
```

Variables declared in native expressions, such as for loops, are the only exception to the rule,
as they get "cleared" as soon as the expression is over:

``` bash
for x in 1..10 {
    echo(x) # 1, 2, 3...
}

echo(x) # Error: x is not defined
```

Worth to note that if a variable gets re-defined within these expressions,
it will temporarely assume its new value, but will rollback to the original
one once the expression is over:

``` bash
x = "hello world"

for x in 1..10 {
    echo(x) # 1, 2, 3...
}

echo(x) # "hello world"
```

## Next

That's about it for this section!

You can now head over to read about [returning values](/syntax/return).