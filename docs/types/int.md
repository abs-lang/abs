<p align="center">
  <a href="https://abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/abs-lang/abs/blob/master/bin/abs-horizontal.png?raw=true">
  </a>
</p>

# Integer

Integers are very straightforward data structures:

``` bash
123456789
```

Most mathematical operators are designed to work
with them

``` bash
2 ** 5 / 1 / 3 + 19
```

Note that integers have what we call a "zero value":
a value that evaluates to `false` when casted to boolean:

``` bash
!!0 # false
```

## Supported functions

### int()

Identity:

``` bash
99.int() # 99
```

### str()

Returns a string containing the integer:

``` bash
99.str() # "99"
```

## Next

That's about it for this section!

You can now head over to read about [arrays](/types/array).