# Array

Arrays represent a list of elements
of any other type:

``` bash
[1, 2, "hello", [1, f(x){ x + 1 }]]
```

They can be looped over:

``` bash
for x in [1, 2] {
    echo(x)
}
```

You can access elements of the array with the index
notation:

``` bash
array[3]
```

Accessing an index that does not exist returns null.

To concatenate arrays, "sum" them:

``` bash
[1, 2] + [3] # [1, 2, 3]
```

## Supported functions

### len()

Returns the length of the array:

``` bash
[1, 2].len() # 2
```

### sum()

Sums the elements of the array. Only supported on arrays of numbers:

``` bash
[1, 1, 1].sum() # 3
```

### sort()

Sorts the array. Only supported on arrays of only numbers
or only strings:

``` bash
[3, 1, 2].sort() # [1, 2, 3]
["b", "a", "c"].sort() # ["a", "b", "c"]
```

### map(f)

Modifies the array by applying the function `f` to all its elements:

``` bash
[0, 1, 2].map(f(x){x+1}) # [1, 2, 3]
```

### some(f)

Returns true when at least one of the elements in the array
returns `true` when applied to the function `f`:

``` bash
[0, 1, 2].map(f(x){x == 1}) # true
[0, 1, 2].map(f(x){x == 4}) # false
```

### every(f)

Returns true when all elements in the array
return `true` when applied to the function `f`:

``` bash
[0, 1, 2].every(f(x){type(x) == "NUMBER"}) # true
[0, 1, 2].every(f(x){x == 0}) # false
```

### find(f)

Returns the first element that returns `true` when applied to the function `f`:

``` bash
["hello", 0, 1, 2].find(f(x){type(x) == "NUMBER"}) # 0
```

### find(f)

Returns a new array with only the elements that returned
`true` when applied to the function `f`:

``` bash
["hello", 0, 1, 2].filter(f(x){type(x) == "NUMBER"}) # [0, 1, 2]
```

### json()

Parses the string as JSON, returning an [hash](/types/hash):

``` bash
"{}".json() # {}
```

Note that currently only JSON objects are supported,
and if the objects contain floats this method will
return an error. Support for floats is coming (see [#29](https://github.com/abs-lang/abs/issues/29))
as well as being able to parse all valid JSON expressions (see [#54](https://github.com/abs-lang/abs/issues/54)).

### contains(e)

Checks whether `e` is present in the array. `e` can only be
a string or number and the array needs to be a heterogeneous array
of strings or number:

``` bash
[1, 2, 3].contains(3) # true
[1, 2, 3].contains(4) # false
```

### str()

Returns the string representation of the array:

``` bash
[1, 2].str() # "[1, 2]"
```

### slice(start, end)

Returns a portion of the array, from `start` to `end`:

``` bash
(1..10).slice(0, 3) # [1, 2, 3]"
```

If `start` is negative, it slices from the end of the string,
back as many characters as the value of `start`:

``` bash
(1..10).slice(-3, 0) # [8, 9, 10]"
```

### shift(start, end)

Removes the first elements from the array, and returns it:

``` bash
a = [1, 2, 3]
a.shift() # 1
a # [2, 3]
```

### reverse()

Reverses the order of the elements in the array:

``` bash
[1, 2].reverse() # [2, 1]
```

### push()

Pushes an element at the end of the array:

``` bash
[1, 2].push(3) # [1, 2, 3]
```

### pop()

Pops the last element from the array, returning it:

``` bash
a = [1, 2, 3]
a.shift() # 3
a # [1, 2]
```

### keys()

Returns an array of the keys in the original array:

``` bash
(1..2).keys() # [0, 1]
```

### join(separator)

Joins the elements of the array by `separator`:

``` bash
[1, 2, 3].join("_") # "1_2_3"
```

## Next

That's about it for this section!

You can now head over to read about [hashes](/types/hash).