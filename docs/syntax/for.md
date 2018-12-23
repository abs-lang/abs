<p align="center">
  <a href="https://abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/odino/abs/blob/master/bin/abs-horizontal.png?raw=true">
  </a>
</p>

# For loops

ABS supports 2 types of `for` loops: the "standard" form and
the "in" one.

## Standard form

A standard `for` loop takes the canonical form:

``` bash
for x = 0; x < 10; x = x + 1 {
    echo("Looping...")
}
```

It's important to understand scoping in this form:
if the loop declares an identifier that was already
declared earlier, it will be "temporarely" overwritten
in the loop, but will assume its original value
afterwards.

Code is better than 1000 words:

``` bash
x = "hello world"

for x = 0; x < 10; x = x + 1 {
    # x is 0, 1, 2...
}

echo(x) # "hello world"
```

Similarly, a variable declared on the `loop` (not inside)
will cease to exist after the loop is done:

``` bash
for x = 0; x < 10; x = x + 1 {
    # x is 0, 1, 2...
}

echo(x) # x is not defined here
```

Finally, variables declared inside the loop will instead
keep living afterwards:

``` bash
for x = 0; x < 10; x = x + 1 {
    y = x
}

echo(y) # 9
```

## In form

The "in" form of `for` loops allows you to iterate over
an array:

``` bash
for x in [1, 2, 3] {
    # x is 1, 2, 3
}
```

Both key and values are available in the loop:

``` bash
for k, v in [1, 2, 3] {
    # k is 0, 1, 2
    # v is 1, 2, 3
}
```

In terms of scoping, the "in" form follows the same rules
as the standard one, meaning that:

``` bash
k = "hello world"

for k, v in [1, 2, 3] {
    # k is 0, 1, 2
    # v is 1, 2, 3
}

echo(k) # "hello world"
echo(v) # v is not defined
```

## Next

That's about it for this section!

You can now head over to read about [for loops](/syntax/while).