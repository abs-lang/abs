# Hash

Hashes represent a list of key-value pairs that can conveniently be accessed with `O(1)` cost:

``` bash
h = {"key": "val"}
h.key # "val"
h["key"] # "val"
```

Note that the `hash.key` hash property form is the preferred one, as it's more coincise and mimics other programming languages.

Accessing a key that does not exist returns null.

An individual hash element may be assigned to via its `hash["key"]`
index or its property `hash.key`. This includes compound operators 
such as `+=`. Note that a new key may be created as well using `hash["newkey"]` or `hash.newkey`:

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

It is also possible to extend a hash using the `+=` operator
with another hash. Note that any existing keys on the left side 
will be replaced with the same key from the right side:

```bash
h = {"a": 1, "b": 2, "c": 3}
h   # {a: 1, b: 2, c: 3}

# extending a hash by += compound operator
h += {"c": 33, "d": 4, "e": 5}
h   # {a: 1, b: 2, c: 33, d: 4, e: 5}
```

In a similar way, we can make a **shallow** copy of a hash using
the `+` operator with an empty hash. Be careful, the empty hash 
must be on the left side of the `+` operator:

```bash
a = {"a": 1, "b": 2, "c": 3}
a   # {a: 1, b: 2, c: 3}

# shallow copy a hash using the + operator with an empty hash
# note well that the empty hash must be on the left side of the +
b = {} + a
b   # {a: 1, b: 2, c: 3}

# modify the shallow copy without changing the original
b.a = 99
b   # {a: 99, b: 2, c: 3}
a   # {a: 1, b: 2, c: 3}
```

If the left side is a `hash["key"]` or `hash.key` and the
right side is a hash, then the resulting hash will have a
new nested hash at `hash.newkey`. This includes `hash["newkey"]`
or `hash.newKey` as well:

```bash
h = {"a": 1, "b": 2, "c": 3}
h # {a: 1, b: 2, c: 3}

# nested hash assigned to hash.key
h.c = {"x": 10, "y": 20}
h # {a: 1, b: 2, c: {x: 10, y: 20}}

# nested hash assigned to hash.newkey
h.z = {"xx": 11, "yy": 21}
h # {a: 1, b: 2, c: {x: 10, y: 20}, z: {xx: 11, yy: 21}}
```

## Supported functions

### str()

Returns the string representation of the hash:

``` bash
h = {"k": "v"}
h.str() # "{k: v}"
str(h)  # "{k: v}"
```

### keys()

Returns an array of keys to the hash. 

Note well that only the first level keys are returned.

``` bash
h = {"a": 1, "b": 2, "c": 3}
h.keys() # [a, b, c]
keys(h) # [a, b, c]
```

### values()

Returns an array of values in the hash. 

Note well that only the first level values are returned.

``` bash
h = {"a": 1, "b": 2, "c": 3}
h.values()  # [1, 2, 3]
values(h)   # [1, 2, 3]
```

### items()

Returns an array of [key, value] tuples for each item in the hash.

Note well that only the first level items are returned.

``` bash
h = {"a": 1, "b": 2, "c": 3}
h.items()   # [[a, 1], [b, 2], [c, 3]]
items(h)    # [[a, 1], [b, 2], [c, 3]]
```

### pop(k)

Removes and returns the matching `{"key": value}` item from the hash. If the key does not exist `hash.pop("key")` returns `null`.

Note well that only the first level items can be popped.

``` bash
h = {"a": 1, "b": 2, "c": {"x": 10, "y":20}}

h.pop("a")  # {a: 1}
h   # {b: 2, c: {x: 10, y: 20}}

pop(h, "c")  # {c: {x: 10, y: 20}}
h   # {b: 2}

pop(h, "d") # null
h   # {b: 2}

```

## User-defined functions

A useful property of being able to assign keys of any type to an hash
results in the ability to define objects with custom functions, such as:

``` bash
hash = {"greeter": f(name) { return "Hello $name!" }}
hash.greeter("Sally") # "Hello Sally!"
```

## Next

That's about it for this section!

You can now head over to read about [functions](/types/function).