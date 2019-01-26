# If ... else

ABS supports basic `if` statements:

``` bash
if x > 0 {
    echo("hello world")
}
```

as well as `else` and `else if` alternatives:

``` bash
if x > 0 {
    echo("x is high")
} else if x < 0 {
    echo("x is low")
} else {
    echo("x is actually zero!")
}
```

You can wrap conditions in parenthesis, although we believe that,
from a readability standpoint, it's usually better to omit them:

``` bash
if (x > 0) {
    echo("hello world")
}
```

## Next

That's about it for this section!

You can now head over to read about [for loops](/syntax/for).