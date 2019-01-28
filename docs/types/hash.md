# Hash

Hashes represent a list of key-value pairs that can conveniently be accessed with `O(1)` cost:

``` bash
h = {"key": "val"}
h.key # "val"
h["key"] # "val"
```

Note that the `hash.key` hash property form is the preferred one, as it's more coincise and mimics other programming languages.

Accessing a key that does not exist returns null.

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

It is also possible to extend a hash using the `+=` operator with another hash. Note that existing keys in the left side will be replaced with the same key on the right side.
```bash
h = {"a": 1, "b": 2, "c": 3}
h # {a: 1, b: 2, c: 3}

# extending a hash by += compound operator
h += {"c": 33, "d": 4, "e": 5}
h # {a: 1, b: 2, c: 33, d: 4, e: 5}
```

If the left side is a `hash["key"]` or `hash.key` and the right side is a hash, then the resulting hash will have a new nested hash at the `hash.new`. This includes `hash["newkey"]` or `hash.newkey` as well.
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
{"k": "v"}.str() # "{k: v}"
```

Note that hashes are set to receive a substantial
boost in the future, as most hash-related
functions didn't make the cut to ABS' first public
release (see [#36](https://github.com/abs-lang/abs/issues/36)).

## Next

That's about it for this section!

You can now head over to read about [functions](/types/function).