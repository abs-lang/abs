<p align="center">
  <a href="https://www.abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/abs-lang/abs/blob/master/bin/abs-horizontal.png?raw=true">
  </a>
</p>

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

``` bash
1 == 1 # true
1 == "hello world" # false
```

## !=

Not equals operator, one of the few that can be used
between arguments of different type:

``` bash
1 != 1 # false
1 != "hello world" # true
```

## +

Addition, can be used to merge arrays and combine strings as well:

``` bash
1 + 1 # 2
"hello" + " " + "world" # "hello world"
[1] + [2] # [1, 2]
```

## -

Subtraction:

``` bash
0 - 1 # -1
```

## *

Multiplication:

``` bash
1 * 2 # 2
```

## /

Division:

``` bash
5 / 5 # 1
```

## >

Greater than:

``` bash
10 > 0 # true
0 > 10 # false
```

## >=

Greater or equal than:

``` bash
1 >= 1 # true
2 >= 1 # true
```

## <

Lower than:

``` bash
10 < 0 # false
0 < 10 # true
```

## <=

Lower or equal than:

``` bash
1 <= 1 # true
1 <= 2 # true
```

## &&

Logical AND, which supports [short-circuiting](https://en.wikipedia.org/wiki/Short-circuit_evaluation):

``` bash
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

``` bash
true || true # true
true || false # true
1 || 2 # 1
1 || 0 # 1
"" || "hello world" # "hello world"
"hello" || "world" # "hello"
```

## ..

Range operator, which creates an array from start to end:

``` bash
1..10 # [1, 2, 3, 4, 5, 6, 7, 8, 9]
```

## **

Mathematical exponentiation:

``` bash
2 ** 2 # 4
2 ** 0 # 1
```

## ~

The tilde (meaning "around") is used to do a case-insensitive
comparison between strings:

``` bash
"hello" == "HELLO" # false
"hello" ~ "HELLO" # true
```

## !

Negation:

``` bash
a = true
!a # false
```

## !!

Even though there is no double negation operator, using
2 bangs will result into converting the argument to boolean:

``` bash
!!1 # true
!!0 # false
!!"" # false
!!"hello" # true
```

## Next

That's about it for this section!

You can now head over to read about data types available in
ABS, starting from the [string](/types/string).