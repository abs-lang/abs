<p align="center">
  <a href="https://www.abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/abs-lang/abs/blob/master/bin/abs-horizontal.png?raw=true">
  </a>
</p>

# Hash

Hashes represent a list of key-value pairs
that can conveniently be accessed with `O(1)`
cost:

``` bash
h = {"key": "val"}
h.key # "val"
h["key"] # "val"
```

Note that the `x.y` form is the preferred one, as it's more coincise
and mimics other programming languages.

Accessing an index that does not exist returns null.

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