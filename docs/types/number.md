# Number

Numbers are very straightforward data structures:

``` bash
123
1.23
```

Most mathematical operators are designed to work
with them

``` bash
(2 ** 5.5 / 1 / 3 + 19) % 5 # 4.08494466531301
```

Note that numbers have what we call a "zero value":
a value that evaluates to `false` when casted to boolean:

``` bash
!!0 # false
```

You can use [bitwise operators](/syntax/operators) on numbers, but bear in
mind that they will be implicitely converted to integers:

``` bash
1 ^ 1 # 0
1 ^ 0 # 1
1 ^ 0.9 # 1, as 0.9 is converted to 0
```

## Supported functions

### number()

Identity:

``` bash
99.5.number() # 99.5
```

### int()

Rounds down the number to the closest integer:

``` bash
10.3.int() # 10
```

### str()

Returns a string containing the number:

``` bash
99.str() # "99"
```

## Next

That's about it for this section!

You can now head over to read about [arrays](/types/array).