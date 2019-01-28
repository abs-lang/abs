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

This is also the suggested way to push a new element into
an array:

``` bash
x = [1, 2]
x += [3]
x # [1, 2, 3]
```

It is also possible to modify an existing array element using `array[index]` assignment. This also works with compound operators such as `+=` :
```bash
a = [1, 2, 3, 4]
a # [1, 2, 3, 4]

# index assignment
a[0] = 99
a # [99, 2, 3, 4]

# compound assignment
a[0] += 1
a # [100, 2, 3, 4]
```

An array can also be extended by using an index beyond the end of the existing array. Note that intervening array elements will be set to `null`. This means that they can be set to another value later:
```bash
a = [1, 2, 3, 4]
a # [1, 2, 3, 4]

# indexes beyond end of array expand the array
a[4] = 99
a # [1, 2, 3, 4, 99]
a[6] = 66
a # [1, 2, 3, 4, 99, null, 66]

# assign to a null element
a[5] = 55
a # [1, 2, 3, 4, 99, 55, 66]
```

An array is defined as "homogeneous" when all its elements
are of a single type:

```
[1, 2, 3] # homogeneous
[null, 0, "", {}] # heterogeneous
```

This is important as some functions are only supported
on homogeneous arrays: `sum()`, for example, can only be
called on homogeneous arrays of numbers.

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

### contains(e)

> This function is deprecated and might be removed in future versions.
>
> Use the "in" operator instead: 3 in [1, 2, 3]

Checks whether `e` is present in the array. `e` can only be
a string or number and the array needs to be a homogeneous array
of strings or numbers:

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