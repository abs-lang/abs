# String

Strings are probably the most basic data type
in all languages, yet they hold a very important
value in ABS: considering that shell scripting
is all about working around command outputs,
we assume you will likely work a lot with them.

Strings are enclosed by double or single quotes:

``` bash
"hello world"
'hello world'
```

You can escape quotes with a simple backslash:

``` bash
"I said: \"hello world\""
```

or use the other quote to ease escaping:

``` bash
'I said: "hello world"'
```

Their individual characters can be accessed
with the index notation:

``` bash
"hello world"[1] # e
```

Accessing an index that does not exist returns null.

To concatenate strings, "sum" them:

``` bash
"hello" + " " + "world" # "hello world"
```

Note that strings have what we call a "zero value":
a value that evaluates to `false` when casted to boolean:

``` bash
!!"" # false
```

## Supported functions

### len()

Returns the length of a string:

``` bash
"hello world".len() # 11
```

### fmt()

Formats a string ([sprintf convention](https://linux.die.net/man/3/sprintf)):

``` bash
"hello %s".fmt("world") # "hello world"
```

### number()

Converts a string to a number, if possible:

``` bash
"99.5".number() # 99.5
"a".number() # ERROR: int(...) can only be called on strings which represent numbers, 'a' given
```

### is_number()

Checks whether a string can be converted to a number:

``` bash
"99.5".is_number() # true
"a".is_number() # false
```

Use this function when `"...".number()` might return an error.

### int()

Converts a string to integer, if possible:

``` bash
"99.5".int() # 99
"a".int() # ERROR: int(...) can only be called on strings which represent numbers, 'a' given
```

### split(separator)

Splits a string by separator:

``` bash
"1.2.3.4".split(".") # ["1", "2", "3", "4"]
```

### lines()

Splits a string by newline:

``` bash
"first\nsecond".lines() # ["first", "second"]
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

### contains(str)

Checks whether `str` is present in the string:

``` bash
"string".contains("ing") # true
"string".contains("ong") # false
```

### str()

Identity:

``` bash
"string".str() # "string"
```

### any(str)

Checks whether any of the characters in `str` are present in the string:

``` bash
"string".any("abc") # true
"string".any("xyz") # false
```

### prefix(str)

Checks whether the string has the given prefix:

``` bash
"string".prefix("str") # true
"string".prefix("abc") # false
```

### suffix(str)

Checks whether the string has the given suffix:

``` bash
"string".suffix("ing") # true
"string".suffix("ong") # false
```

### repeat(i)

Creates a new string, repeating the original one `i` times:

``` bash
"string".repeat(2) # "stringstring"
```

### replace(x, y, n)

Replaces occurrences of `x` with `y`, `n` times.
If `n` is negative it will replace all occurrencies:

``` bash
"string".replace("i", "o", -1) # "strong"
```

### title()

Titlecases the string:

``` bash
"hello world".title() # "Hello World"
```

### lower()

Lowercases the string:

``` bash
"STRING".lower() # "string"
```

### upper()

Uppercases the string:

``` bash
"string".upper() # "STRING"
```

### trim()

Removes empty spaces from the beginning and end of the string:

``` bash
" string     ".trim() # "string"
```

### trim_by(str)

Removes `str` from the beginning and end of the string:

``` bash
"string".trim_by("g") # "strin"
```

### index(str)

Returns the index at which `str` is found:

``` bash
"string".index("t") # 1
```

### last_index(str)

Returns the last index at which `str` is found:

``` bash
"string string".last_index("g") # 13
```

### slice(start, end)

Returns a portion of the string, from `start` to `end`:

``` bash
"string".slice(0, 3) # "str"
```

If `start` is negative, it slices from the end of the string,
back as many characters as the value of `start`:

``` bash
"string".slice(-3, 0) # "ing"
```

## Next

That's about it for this section!

You can now head over to read about [numbers](/types/number).