---
permalink: /syntax/operators
---

# Operators

Operators are natural building blocks of languages, so it's
no surprise ABS has its fair share of them.

As a general rule, you should understand that ABS does not
allow to use operator on different types, with a couple exceptions.
Here is a list of operators you can use, with examples on how
to make the most out of them.

## ==

Equality operator, one of the few that can be used
between arguments of different type:

```bash
1 == 1 # true
1 == "hello world" # false
```

## !=

Not equals operator, one of the few that can be used
between arguments of different type:

```bash
1 != 1 # false
1 != "hello world" # true
```

## +

Addition, can be used to merge arrays and combine strings as well:

```bash
1 + 1 # 2
"hello" + " " + "world" # "hello world"
[1] + [2] # [1, 2]
```

## +=

Compound addition:

```bash
a = 10
a += 1 # a is now 11
```

## -

Subtraction:

```bash
0 - 1 # -1
```

## -=

Compound subtraction:

```bash
a = 10
a -= 1 # a is now 9
```

## \*

Multiplication:

```bash
1 * 2 # 2
```

## \*=

Compound multiplication:

```bash
a = 10
a *= 10 # a is now 100
```

## /

Division:

```bash
5 / 5 # 1
```

## /=

Compound division:

```bash
a = 10
a /= 2 # a is now 5
```

## in

Membership test operator (find whether a needle is in the haystack):

```bash
1 in [1, 2, 3] # true
9 in [1, 2, 3] # false
9 in 9 # unknown operator: NUMBER in NUMBER
"str" in "string" # true
"xyz" in "string" # false
"x" in {"x": 1} # true
"y" in {"x": 1} # false
```

## \*\*

Mathematical exponentiation:

```bash
2 ** 2 # 4
2 ** 0 # 1
```

## \*\*=

Compound exponentiation:

```bash
a = 10
a **= 0 # a is now 1
```

## %

Modulo:

```bash
19 % 5 # 4
```

## %=

Compound modulo:

```bash
a = 19
a %= 5 # a is now 4
```

## >

Greater than:

```bash
10 > 0 # true
0 > 10 # false
```

When used with strings, `>` redirects
the input (left argument) to a file
(right argument):

```bash
"Hi, there" > "file.txt"
`cat /etc/hosts` > "copy-of-etc-hosts.bkp"
```

## >=

Greater or equal than:

```bash
1 >= 1 # true
2 >= 1 # true
```

## <

Lower than:

```bash
10 < 0 # false
0 < 10 # true
```

## <=

Lower or equal than:

```bash
1 <= 1 # true
1 <= 2 # true
```

## <=>

The combined comparison operator allows to test whether a number
is lower, equal or higher than another one with one statement:

```bash
5 <=> 5 # 0
5 <=> 6 # -1
6 <=> 5 # 1
```

## &&

Logical AND, which supports [short-circuiting](https://en.wikipedia.org/wiki/Short-circuit_evaluation):

```bash
true && true # true
true && false # false
1 && 2 # 2
1 && 0 # 0
0 && 2 # 0
"" && "hello world" # ""
"hello" && "world" # "world"
```

## ||

Logical OR, which supports [short-circuiting](https://en.wikipedia.org/wiki/Short-circuit_evaluation):

```bash
true || true # true
true || false # true
1 || 2 # 1
1 || 0 # 1
"" || "hello world" # "hello world"
"hello" || "world" # "hello"
```

## .

Property accessor, used to access properties or methods of specific variables:

```bash
hello = {"to_who": "the world"}
hello.to_who # "the world"
```

There are some builtin functions that you can access through the property accessor:

```bash
"hello".len() # 5
```

(a comprehensive list of function is documented in the "_Types and functions_" section of the documentation)

## ?.

Optional chaining operator, used to access properties in a "safe" way.

Given the following object:

```bash
test = {"property": 1}
```

An error would be raised if you were trying to access a non-existing property
such as `test.something.something_else`:

```
ERROR: invalid property 'something_else' on type NULL
	[1:15]	test.something.something_else
```

Optional chainig prevents those errors from being raised, auto-magically
converting non-existing properties and methods to `NULL`:

```bash
test?.something?.something_else # null
test?.something?.something_else() # null
```

## ..

Range operator, which creates an array from start to end:

```bash
1..10 # [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
10..1 # [10, 9, 8, 7, 6, 5, 4, 3, 2, 1]
```

## !

Negation:

```bash
a = true
!a # false
```

## !!

Even though there is no double negation operator, using
2 bangs will result into converting the argument to boolean:

```bash
!!1 # true
!!0 # false
!!"" # false
!!"hello" # true
```

## ~

The tilde (meaning "around") is used to do similarity comparisons
(x "is almost" y).

In case of strings, it will compare them case-insensitively

```bash
"hello" == "HELLO" # false
"hello" ~ "HELLO" # true
```

In case of integers, it will compare their integer values:

```bash
1 == 1.5 # false
1 ~ 1.5 # true
```

When in front of a number, it will instead be used as a
bitwise NOT:

```bash
~0 # -1
~"hello" # ERROR: Bitwise not (~) can only be applied to numbers, got STRING (hello)
```

## &

Bitwise AND:

```bash
1 & 1 # 1
1 & "hello" # ERROR: type mismatch: NUMBER & STRING
```

## |

Bitwise OR:

```bash
1 | 1 # 1
1 | "hello" # ERROR: type mismatch: NUMBER | STRING
```

## ^

Bitwise XOR:

```bash
1 ^ 1 # 0
1 ^ "hello" # ERROR: type mismatch: NUMBER ^ STRING
```

## >>

Bitwise right shift:

```bash
1 >> 1 # 0
1 >> "hello" # ERROR: type mismatch: NUMBER >> STRING
```

When used with strings, `>>` appends
the input (left argument) to a file
(right argument):

```bash
"Hi" >> "file.txt"
" " >> "file.txt"
"there!" >> "file.txt"

`cat file.txt` # Hi there!
```

## <<

Bitwise left shift:

```bash
1 << 1 # 2
1 << "hello" # ERROR: type mismatch: NUMBER << STRING
```