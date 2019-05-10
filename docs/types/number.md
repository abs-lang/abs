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

You can write numbers in the exponential notation:

```
1e1 # 10
1e+1 # 10
1e-1 # 0.1
```

In addition, numbers can include underscores (`_`) as visual
separators, in order to improve readability: when
ABS encounters `1_000_000` it will internally convert it
to a million. Underscore separators can be placed anywhere
on a number (`10_`, `10_00`, `10.00_00_00`) except at its start:

```
1000000 # 1M
1_000_000 # 1M, just a lot more readable
1_00_00_00 # 1M, formatted with another separator pattern
_100000000 # ERROR: identifier not found: _
```

Note there is no limit to the amount of consecutive
underscores that can be used (eg. `10__________0`).

## Supported functions

### number()

Identity:

``` bash
99.5.number() # 99.5
```

### int()

Rounds **down** the number to the closest integer:

``` bash
10.3.int() # 10
```

### round(precision?)

Rounds the number with the given precision.
The precision argument is optional, and set to `0`
by default:

``` bash
10.3.round() # 10
10.6.round() # 11
10.333.round(1) # 10.3
```

### ceil()

Rounds the number up to the closest integer:

``` bash
10.3.ceil() # 11
```

### floor()

Rounds the number down to the closest integer:

``` bash
10.9.floor() # 10
```

### str()

Returns a string containing the number:

``` bash
99.str() # "99"
```

## Next

That's about it for this section!

You can now head over to read about [arrays](/types/array).