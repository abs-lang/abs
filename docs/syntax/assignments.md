# Assignments

Just like about any other language, assignments are pretty straightforward:

``` bash
x = "hello world"
```

Array destructuring is supported, meaning you can set multiple variables based on an array:

``` bash
x, y, z = ["hello world", 99, {}]
x # "hello world"
y # 99
z # {}
```

If the number of variables you're trying to set is longer than the array, the extra variables will be set to null:

``` bash
x, y = [1]
y # null
```

An individual array element may be assigned to via its `array[index]`. This includes compound operators such as `+=`. Also an array can be extended by assigning to an index beyond its current length.
```bash
a = [1, 2, 3, 4]
a # [1, 2, 3, 4]

# index assignment
a[0] = 99
a # [99, 2, 3, 4]

# compound assignment
a[0] += 1
a # [100, 2, 3, 4]

# extending an array; note intervening nulls are created if needed
a[5] = 55
a # [100, 2, 3, 4, null, 55]
a[4] = 44
a # [100, 2, 3, 4, 44, 55]
```

An individual hash element may be assigned to via its `hash["key"]` index or its property `hash.key`. This includes compound operators such as `+=`. Note that a new key may be created as well using `hash["newkey"]` or `hash.newkey`.
```bash
h = {"a": 1, "b": 2, "c": 3}
h # {a: 1, b: 2, c: 3}

# index assignment
h["a"] = 99
h # {a: 99, b: 2, c: 3}

# property assignment
h.a # 99
h.a = 88
h # {a: 88, b: 2, c: 3}

# compound operator assignment to property
h.a += 1
h.a # 89
h # {a: 88, b: 2, c: 3}

# create new keys via index or property
h["x"] = 10
h.y = 20
h # {a: 88, b: 2, c: 3, x: 10, y: 20}
```

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
it will temporaraly assume its new value, but will rollback to the original
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