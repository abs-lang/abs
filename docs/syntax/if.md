<p align="center">
  <a href="https://www.abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/abs-lang/abs/blob/master/bin/abs-horizontal.png?raw=true">
  </a>
</p>

# If ... else

ABS supports basic `if` statements:

``` bash
if x > 0 {
    echo("hello world")
}
```

as well as `else` alternatives:

``` bash
if x > 0 {
    echo("hello world")
} else {
    echo("hello globe")
}
```

You can wrap conditions in parenthesis, although we believe that,
from a readability standpoint, it's usually better to omit them:

``` bash
if (x > 0) {
    echo("hello world")
}
```

Note that `else if` clauses are not supported,
although they are planned (see [#27](https://github.com/abs-lang/abs/issues/27)).

## Next

That's about it for this section!

You can now head over to read about [for loops](/syntax/for).