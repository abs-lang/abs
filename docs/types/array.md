# Array

Arrays represent lists of elements
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

You can access elements of the array with `[]` index
notation:

``` bash
array[3]
```

Accessing an array element that does not exist returns `null`.

You can also access the Nth last element of an array
with a negative index:

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

If `end` is negative, it will be converted to `length of array - (-end)`:

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

### chunk(size)

Splits the array into chunks of the given `size`:

```py
[1, 2, 3].chunk(2) # [[1, 2], [3]]
[1, 2, 3].chunk(10) # [[1,2,3]]
[1, 2, 3].chunk(1.2) # argument to chunk must be a positive integer, got '1.2'
```

### diff(array)

Computes the difference between 2 arrays,
returning elements that are only in the first array:

```py
[1, 2, 3].diff([]) # [1, 2, 3]
[1, 2, 3].diff([3]) # [1, 2]
[1, 2, 3].diff([3, 1]) # [2]
[1, 2, 3].diff([1, 2, 3, 4]) # []
```

For symmetric difference see [diff_symmetric(...)](#diff_symmetricarray)

### diff_symmetric(array)

Computes the [symmetric difference](https://en.wikipedia.org/wiki/Symmetric_difference)
between 2 arrays, returning elements that are only in one of the arrays:

```py
[1, 2, 3].diff_symmetric([]) # [1, 2, 3]
[1, 2, 3].diff_symmetric([3]) # [1, 2]
[1, 2, 3].diff_symmetric([3, 1]) # [2]
[1, 2, 3].diff_symmetric([1, 2, 3, 4]) # [4]
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

### flatten()

Concatenates the lowest "layer" of elements in a nested array:

```py
[[1, 2], 3, [4]].flatten() # [1, 2, 3, 4]
[[1, 2, 3, 4]].flatten() # [1, 2, 3, 4]
[[[1, 2], [3, 4], 5, 6], 7, 8].flatten() # [[1, 2], [3, 4], 5, 6, 7, 8]
```

### flatten_deep()

Recursively flattens an array until no element is an array:

```py
[[[1, 2], [[[[3]]]], [4]]].flatten_deep() # [1, 2, 3, 4]
[[1, [2, 3], 4]].flatten_deep() # [1, 2, 3, 4]
```

### intersect(array)

Computes the intersection between 2 arrays:

```py
[1, 2, 3].intersect([]) # []
[1, 2, 3].intersect([3]) # [3]
[1, 2, 3].intersect([3, 1]) # [1, 3]
[1, 2, 3].intersect([1, 2, 3, 4]) # [1, 2, 3]
```

### join([separator])

Joins the elements of the array with the string `separator` (default "", the empty string):

``` py
[1, 2, 3].join("_") # "1_2_3"
[1, 2, 3].join()    # "123"
```

### keys()

Returns an array of the keys in the original array:

``` py
(1..2).keys() # [0, 1]
```

### len()

Returns the length of the array:

``` py
[1, 2].len() # 2
```

### map(f)

Modifies the array by applying the function `f` to all its elements:

``` py
[0, 1, 2].map(f(x){x+1}) # [1, 2, 3]
```

### max()

Finds the highest number in an array:

```py
[].max() # NULL
[0, 5, -10, 100].max() # 100
```

### min()

Finds the lowest number in an array:

```py
[].min() # NULL
[0, 5, -10, 100].min() # -10
```

### partition(f)

Partitions the array by applying `f(element)` to all of its elements,
then grouping the elements into an array of arrays based on the results:

```py
f odd(n) {
  return !!(n % 2)
}
f div2(n) {
  return int(n / 2)
}
[0, 1, 2, 3, 4, 5].partition(odd) # [[0, 2, 4], [1, 3, 5]]
[5, 4, 3, 2, 1, 0].partition(div2) # [[5, 4], [3, 2], [1, 0]]
["1", {}, 0, "0", 1].partition(str) # [["1", 1], [{}], [0, "0"]]
```

### pop()

Removes and returns the last element from the array:

``` py
a = [1, 2, 3]
a.pop() # 3
a # [1, 2]
```

### push(x)

Inserts `x` at the end of the array:

``` py
[1, 2].push(3) # [1, 2, 3]
```

This is equivalent to summing 2 arrays:

``` py
[1, 2] + [3] # [1, 2, 3]
```

### reduce(f, accumulator)

Reduces the array to a value by iterating through its elements and applying the two-argument function `f(value, element)` to them, with `accumulator` as the initial `value`:

```py
[1, 2, 3, 4].reduce(f(value, element) { return value + element }, 0) # 10
[1, 2, 3, 4].reduce(f(value, element) { return value + element }, 10) # 20
```

### reverse()

Reverses the order of the elements in the array:

``` py
[1, 2].reverse() # [2, 1]
```

### shift()

Removes the first element from the array, and returns it:

``` py
a = [1, 2, 3]
a.shift() # 1
a # [2, 3]
```

### shuffle()

Shuffles elements in the array:

``` py
a = [1, 2, 3, 4]
a.shuffle() # [3, 1, 2, 4]
```

### some(f)

Returns true when at least one of the elements in the array
returns `true` when applied to the function `f`:

``` py
[0, 1, 2].map(f(x){x == 1}) # true
[0, 1, 2].map(f(x){x == 4}) # false
```

### sort()

Sorts the array. Only supported on homogeneous arrays of numbers
or strings:

```py
[3, 1, 2].sort() # [1, 2, 3]
["b", "a", "c"].sort() # ["a", "b", "c"]
[42, "hut", 37].sort()
ERROR: argument to 'sort' must be an homogeneous array (elements of the same type), got [42, "hut", 37]
	[1:16]	[42, "hut", 37].sort()
```

### str()

Returns the string representation of the array:

```py
[1, 2].str() # "[1, 2]"
```

### sum()

Sums the elements of the array. Only supported on homogeneous arrays of numbers:

```py
[1, 1, 1].sum() # 3
```

### tsv([separator[, header]])

Formats the array as a TSV (Tab-Separated Values):

``` bash
[["LeBron", "James"], ["James", "Harden"]].tsv()
LeBron	James
James	Harden
```

You can also specify the `separator` to be used if you
prefer not to use tabs:

``` bash
[["LeBron", "James"], ["James", "Harden"]].tsv(",")
LeBron,James
James,Harden
```

The input must be an array of arrays or hashes. If
you use hashes, their keys will be used as the first row of the TSV:

```bash
[{"name": "Lebron", "last": "James", "jersey": 23}, {"name": "James", "last": "Harden"}].tsv()
jersey	last	name
23	James	Lebron
null	Harden	James
```

The first row will, by default, be a combination of all keys present in the hashes,
sorted alphabetically. If a key is missing in a hash, `null` will be used as its value.

`header` is an optional array of output keys, whose values are output in the specified order:

```bash
[{"name": "Lebron", "last": "James", "jersey": 23}, {"name": "James", "last": "Harden"}].tsv("\t", ["name", "last", "jersey", "additional_key"])
name	last	jersey	additional_key
Lebron	James	23	null
James	Harden	null	null

[{"name": "Lebron", "last": "James", "jersey": 23}, {"name": "James", "last": "Harden"}].tsv(",", ["last", "jersey"])
last,jersey
James,23
Harden,null
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

### unique()

Returns the array with duplicate values removed. The values need not be sorted:

```py
[1, 1, 1, 2].unique() # [1, 2]
[2, 1, 2, 3].unique() # [2, 1, 3]
```

## Next

That's about it for this section!

You can now head over to read about [hashes](/types/hash).
