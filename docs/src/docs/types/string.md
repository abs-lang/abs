---
permalink: /types/string
---

# String

Strings are probably the most basic data type
in all languages, yet they hold a very important
value in ABS: considering that shell scripting
is all about working around command outputs,
we assume you will likely work a lot with them.

Strings are enclosed by double or single quotes:

```bash
"hello world"
'hello world'
```

You can escape quotes with a simple backslash:

```bash
"I said: \"hello world\""
```

or use the other quote to ease escaping:

```bash
'I said: "hello world"'
```

Their individual characters can be accessed
with the index notation:

```bash
"hello world"[1] # e
```

Accessing an index that does not exist returns an empty string.

You can access the Nth last character of the string using a
negative index:

```bash
"string"[-2] # "n"
```

You can also access a range of the string with the `[start:end]` notation:

```bash
"string"[0:3] // "str"
```

where `start` is the starting position in the array, and `end` is
the ending one. If `start` is not specified, it is assumed to be 0,
and if `end` is omitted it is assumed to be the last character in the
string:

```bash
"string"[:3] // "str"
"string"[1:] // "tring"
```

If `end` is negative, it will be converted to `length of string - (-end)`:

```bash
"string"[0:-1] // "strin"
```

To concatenate strings, "sum" them:

```bash
"hello" + " " + "world" # "hello world"
```

Note that strings have what we call a "zero value":
a value that evaluates to `false` when casted to boolean:

```bash
!!"" # false
```

To test for the existence of substrings within strings use the `in` operator:

```bash
"str" in "string"   # true
"xyz" in "string"   # false
```

## Interpolation

You can also replace parts of the string with variables
declared within your program using the `$` symbol:

```bash
file = "/etc/hosts"
x = "File name is: $file"
echo(x) # "File name is: /etc/hosts"
```

If you need `$` literals in your command, you
simply need to escape them with a `\`:

```bash
"$non_existing_var" # "" since the ABS variable 'non_existing_var' doesn't exist
"\$non_existing_var" # "$non_existing_var"
```

An alternative syntax (`${...}`) is available for special
cases -- for example, when your string is embedded
within another string:

```bash
word = "word"
echo("prefix$wordsuffix") # "prefix"
echo("prefix${word}suffix") # "prefixwordsuffix"
```

## Special characters embedded in strings

Double and single quoted strings behave differently if the string contains
escaped special ASCII line control characters such as `LF "\n"`, `CR "\r"`,
and `TAB "\t"`.

If the string is double quoted these characters will be expanded to their ASCII codes.
On the other hand, if the string is single quoted, these characters will be considered
as escaped literals.

This means, for example, that double quoted LFs will cause line feeds to appear in the output:

```bash
⧐  echo("a\nb\nc")
a
b
c
⧐
```

Conversely, single quoted LFs will appear as escaped literal strings:

```bash
⧐  echo('a\nb\nc')
a\nb\nc
⧐
```

And if you need to mix escaped and unescaped special characters, then you can do this with double escapes within double quoted strings:

```bash
⧐  echo("a\\nb\nc")
a\\nb
c
⧐
```

## Unicode support

Unicode characters are supported in strings:

```bash
⧐  echo("⺐")
⺐
⧐  echo("I ❤ ABS")
I ❤ ABS
```

### Working with special characters in string functions

Special characters also work with `split()` and `join()` and other string functions as well.

1. Double quoted expanded special characters:

```bash
⧐  s = split("a\nb\nc", "\n")
⧐  echo(s)
[a, b, c]
⧐  ss = join(s, "\n")
⧐  echo(ss)
a
b
c
⧐
```

2. Single quoted literal special characters:

```bash
⧐  s = split('a\nb\nc', '\n')
⧐  echo(s)
[a, b, c]
⧐  ss = join(s, '\n')
⧐  echo(ss)
a\nb\nc
⧐
```

3. Double quoted, double escaped special characters:

```bash
⧐  s = split("a\\nb\\nc", "\\n")
⧐  echo(s)
[a, b, c]
⧐  ss = join(s, "\\n")
⧐  echo(ss)
a\\nb\\nc
⧐
```

## Supported functions

### any(str)

Checks whether any of the characters in `str` are present in the string:

```bash
"string".any("abs") # true
"string".any("xyz") # false
```

### camel()

Converts the string to camelCase:

```bash
"a short sentence".camel() # aShortSentence
```

### ceil()

Converts a string to a number, and then rounds the
number up to the closest integer.

The string must represent a number.

```bash
"10.3".ceil() # 11
"-10.3".ceil() # -10
"a".ceil() # ERROR: ceil(...) can only be called on strings which represent numbers, 'a' given
```

### floor()

Converts a string to a number, and then rounds the
number down to the closest integer.

The string must represent a number.

```bash
"10.9".floor() # 10
"-10.9".floor() # -11
"a".floor() # ERROR: floor(...) can only be called on strings which represent numbers, 'a' given
```

### fmt()

Formats a string ([sprintf convention](https://linux.die.net/man/3/sprintf)):

```bash
"hello %s".fmt("world") # "hello world"
```

In order to print a literal `%`, you can simply escape it with another `%%`:

```bash
"30%%".fmt() # 30%
"30%% %s".fmt("higher") # 30% higher
"30%".fmt() # 30%!(NOVERB)
```

### index(str)

Returns the first index at which `str` is found:

```bash
"string".index("t") # 1
"string".index("ri") # 2
```

### int()

Converts a string to a number, and then rounds it
towards zero to the closest integer.
The string must represent a number.

```bash
"99.5".int() # 99
"-99.5".int() # -99
"a".int() # ERROR: int(...) can only be called on strings which represent numbers, 'a' given
```

### is_number()

Checks whether a string can be converted to a number:

```bash
"99.5".is_number() # true
"a".is_number() # false
```

Use this function when `"...".number()` might return an error.

### json()

Parses the string as JSON, returning a [hash](/types/hash):

```bash
⧐  s = '{"a": 1, "b": "string", "c": true, "d": {"x": 10, "y": 20}}'
⧐  h = s.json()
⧐  h
{a: 1, b: string, c: true, d: {x: 10, y: 20}}
⧐  h.d
{x: 10, y: 20}
```

### kebab()

Converts the string to kebab-case:

```bash
"a short sentence".snake() # a-short-sentence
```

### last_index(str)

Returns the last index at which `str` is found:

```bash
"string string".last_index("g") # 13
"string string".last_index("ri") # 9
```

### len()

Returns the length of a string:

```bash
"hello world".len() # 11
```

### lines()

Splits a string by newline:

```bash
"first\nsecond".lines() # ["first", "second"]
```

### lower()

Lowercases the string:

```bash
"STRING".lower() # "string"
```

### number()

Converts a string to a number, if possible:

```bash
"99.5".number() # 99.5
"a".number() # ERROR: int(...) can only be called on strings which represent numbers, 'a' given
```

### prefix(str)

Checks whether the string starts with `str`:

```bash
"string".prefix("str") # true
"string".prefix("abc") # false
```

### repeat(i)

Creates a new string by repeating the original one `i` times:

```bash
"string".repeat(2) # "stringstring"
```

### replace(str1, str2 [, n])

Replaces the first `n` occurrences of `str1` in the string with `str2`.
If `n` is omitted or negative, it will replace all occurrences:

```bash
"string".replace("i", "o", -1) # "strong"
"aaaa".replace("a", "x") # "xxxx"
"aaaa".replace("a", "x", 2) # "xxaa"
"A man, a plan, a canal, Panama!".replace("a ", "ur-") # "A man, ur-plan, ur-canal, Panama!"
```

### reverse()

Returns a new string with the order of characters/glyphs reversed from the
source.

```bash
"hello world".reverse() # "dlrow olleh"
"世界".reverse() # "界世"
```

### round(precision?)

Converts a string to a number, and then rounds
the number with the given precision.

The precision argument is optional, and set to `0`
by default.

The string must represent a number.

```bash
"10.3".round() # 10
"10.6".round() # 11
"10.333".round(1) # 10.3
"a".round() # ERROR: round(...) can only be called on strings which represent numbers, 'a' given
```

You can also replace an array of strings:

```bash
"string".replace(["i", "g"], "o") # "strono"
"A man, a plan, a canal, Panama!".replace(["a ", "l"], "ur-") # "A man, ur-pur-an, ur-canaur-, Panama!"
```

### snake()

Converts the string to snake_case:

```bash
"a short sentence".snake() # a_short_sentence
```

### split(separator)

Splits a string by `separator`, defaulting to a single space:

```bash
"1.2.3.4".split(".") # ["1", "2", "3", "4"]
"1 2 3 4".split()    # ["1", "2", "3", "4"]
```

### str()

Identity:

```bash
"string".str() # "string"
```

### suffix(str)

Checks whether the string ends with `str`:

```bash
"string".suffix("ing") # true
"string".suffix("ong") # false
```

### title()

Titlecases the string:

```bash
"hello world".title() # "Hello World"
```

### trim()

Removes empty spaces from the beginning and end of the string:

```bash
" string     ".trim() # "string"
```

### trim_by(str)

Removes `str` from the beginning and end of the string:

```bash
"string".trim_by("g") # "strin"
"stringest".trim_by("st") # "ringe"
```

### upper()

Uppercases the string:

```bash
"string".upper() # "STRING"
```
