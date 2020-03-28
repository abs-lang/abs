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

Accessing an index that does not exist returns `null`.

You can also access the Nth last element of an array by
using a negative index:

``` bash
["a", "b", "c", "d"][-2] # "c"
```

You can also access a range of indexes with the `[start:end]` notation:

``` bash
array = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]

array[0:2] # [0, 1, 2]
```

where `start` is the starting position in the array, and `end` is
the ending one. If `start` is not specified, it is assumed to be 0,
and if `end` is omitted it is assumed to be the last index in the
array:

``` bash
array[:2] # [0, 1, 2]
array[7:] # [7, 8, 9]
```

If `end` is negative, it will be converted to `length of array - end`:

``` bash
array[:-3] # [0, 1, 2, 3, 4, 5, 6]
```

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

In a similar way, we can make a **shallow** copy of an array using the `+` operator with an empty array. Be careful, the empty array must be on the left side of the `+` operator.

```bash
a = [1, 2, 3]
a   # [1, 2, 3]

# shallow copy an array using the + operator with an empty array
# note well that the empty array must be on the left side of the +
b = [] + a
b   # [1, 2, 3]

# modify the shallow copy without changing the original
b[0] = 99
b   # [99, 2, 3]
a   # [1, 2, 3]
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

``` bash
[1, 2, 3] # homogeneous
[null, 0, "", {}] # heterogeneous
```

This is important as some functions are only supported
on homogeneous arrays: `sum()`, for example, can only be
called on homogeneous arrays of numbers.

## Supported functions

### contains(e)

> This function is deprecated and might be removed in future versions.
>
> Use the "in" operator instead: 3 in [1, 2, 3]

Checks whether `e` is present in the array. `e` can only be
a string or number and the array needs to be a homogeneous array
of strings or numbers:

``` py
[1, 2, 3].contains(3) # true
[1, 2, 3].contains(4) # false
```

### chunk(size)

Splits the array into chunks of the given size:

```py
[1, 2, 3].chunk(2) # [[1, 2], [3]]
[1, 2, 3].chunk(10) # [[1,2,3]]
[1, 2, 3].chunk(1.2) # argument to chunk must be a positive integer, got '1.2'
```

### every(f)

Returns true when all elements in the array
return `true` when applied to the function `f`:

``` py
[0, 1, 2].every(f(x){type(x) == "NUMBER"}) # true
[0, 1, 2].every(f(x){x == 0}) # false
```

### filter(f)

Returns a new array with only the elements that returned
`true` when applied to the function `f`:

``` py
["hello", 0, 1, 2].filter(f(x){type(x) == "NUMBER"}) # [0, 1, 2]
```

### find(f)

Returns the first element that returns `true` when applied to the function `f`:

``` py
["hello", 0, 1, 2].find(f(x){type(x) == "NUMBER"}) # 0
```

A shorthand syntax supports passing a hash and comparing
elements to the given hash:

```py
[null, {"key": "val", "test": 123}].find({"key": "val"}) # {"key": "val", "test": 123}
```

### len()

Returns the length of the array:

``` py
[1, 2].len() # 2
```

### join(separator)

Joins the elements of the array by `separator`, defaulting to an empty string:

``` py
[1, 2, 3].join("_") # "1_2_3"
[1, 2, 3].join()    # "123"
```

### keys()

Returns an array of the keys in the original array:

``` py
(1..2).keys() # [0, 1]
```

### map(f)

Modifies the array by applying the function `f` to all its elements:

``` py
[0, 1, 2].map(f(x){x+1}) # [1, 2, 3]
```

### pop()

Pops the last element from the array, returning it:

``` py
a = [1, 2, 3]
a.shift() # 3
a # [1, 2]
```

### push()

Pushes an element at the end of the array:

``` py
[1, 2].push(3) # [1, 2, 3]
```

This is equivalent to summing 2 arrays:

``` py
[1, 2] + [3] # [1, 2, 3]
```

### reverse()

Reverses the order of the elements in the array:

``` py
[1, 2].reverse() # [2, 1]
```

### shift(start, end)

Removes the first elements from the array, and returns it:

``` py
a = [1, 2, 3]
a.shift() # 1
a # [2, 3]
```

### slice(start, end)

Returns a portion of the array, from `start` to `end`:

``` py
(1..10).slice(0, 3) # [1, 2, 3]"
```

If `start` is negative, it slices from the end of the string,
back as many characters as the value of `start`:

``` bash
(1..10).slice(-3, 0) # [8, 9, 10]"
```

### some(f)

Returns true when at least one of the elements in the array
returns `true` when applied to the function `f`:

``` py
[0, 1, 2].map(f(x){x == 1}) # true
[0, 1, 2].map(f(x){x == 4}) # false
```

### sort()

Sorts the array. Only supported on arrays of only numbers
or only strings:

```py
[3, 1, 2].sort() # [1, 2, 3]
["b", "a", "c"].sort() # ["a", "b", "c"]
```

### str()

Returns the string representation of the array:

```py
[1, 2].str() # "[1, 2]"
```

### sum()

Sums the elements of the array. Only supported on arrays of numbers:

```py
[1, 1, 1].sum() # 3
```

### tsv([separator], [header])

Formats the array into TSV:

``` bash
[["LeBron", "James"], ["James", "Harden"]].tsv()
LeBron	James
James	Harden
```

You can also specify the separator to be used if you
prefer not to use tabs:

``` bash
[["LeBron", "James"], ["James", "Harden"]].tsv(",")
LeBron,James
James,Harden
```

The input array needs to be an array of arrays or hashes. If
you use hashes, their keys will be used as heading of the TSV:

```bash
[{"name": "Lebron", "last": "James", "jersey": 23}, {"name": "James", "last": "Harden"}].tsv()
jersey	last	name
23	James	Lebron
null	Harden	James
```

The heading will, by default, be a combination of all keys present in the hashes,
sorted alphabetically. If a key is missing in an hash, `null` will be used as value.
If you wish to specify the output format, you can pass a list of keys to be used
as header:

```bash
[{"name": "Lebron", "last": "James", "jersey": 23}, {"name": "James", "last": "Harden"}].tsv("\t", ["name", "last", "jersey", "additional_key"])
name	last	jersey	additional_key
Lebron	James	23	null
James	Harden	null	null
```

### unique()

Returns an array with unique values:

```py
[1, 1, 1, 2].unique() # [1, 2]
```

### intersect(array)

Computes the intersection between 2 arrays:

```py
[1, 2, 3].intersect([]) # []
[1, 2, 3].intersect([3]) # [3]
[1, 2, 3].intersect([3, 1]) # [1, 3]
[1, 2, 3].intersect([1, 2, 3, 4]) # [1, 2, 3]
```

### diff(array)

Computes the difference between 2 arrays,
returning elements that are only on the first array:

```py
[1, 2, 3].diff([]) # [1, 2, 3]
[1, 2, 3].diff([3]) # [1, 2]
[1, 2, 3].diff([3, 1]) # [2]
[1, 2, 3].diff([1, 2, 3, 4]) # []
```

For symmetric difference see [diff_symmetric(...)](#diff_symmetricarray)

### diff_symmetric(array)

Computes the [symmetric difference](https://en.wikipedia.org/wiki/Symmetric_difference)
between 2 arrays (elements that are only on either of the 2):

```py
[1, 2, 3].diff([]) # [1, 2, 3]
[1, 2, 3].diff([3]) # [1, 2]
[1, 2, 3].diff([3, 1]) # [2]
[1, 2, 3].diff([1, 2, 3, 4]) # [4]
```

### union(array)

Computes the [union](https://en.wikipedia.org/wiki/Union_(set_theory))
between 2 arrays:

```py
[1, 2, 3].union([1, 2, 3, 4]) # [1, 2, 3, 4]
[1, 2, 3].union([3]) # [1, 2, 3]
[].union([3, 1]) # [3, 1]
[1, 2].union([3, 4]) # [1, 2, 3, 4]
```

### flatten()

Flattens an array a single level deep:

```py
[[1, 2], 3, [4]].flatten() # [1, 2, 3, 4]
[[1, 2, 3, 4]].flatten() # [1, 2, 3, 4]
```

### flatten_deep()

Flattens an array recursively until no member is an array:

```py
[[[1, 2], [[[[3]]]], [4]]].flatten_deep() # [1, 2, 3, 4]
[[1, [2, 3], 4]].flatten_deep() # [1, 2, 3, 4]
```

## Next

That's about it for this section!

You can now head over to read about [hashes](/types/hash).
