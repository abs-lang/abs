<p align="center">
  <a href="https://abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/odino/abs/blob/master/bin/abs-horizontal.png?raw=true">
  </a>
</p>

# Returning values

We promise, this is going to be short!

Returning values is done through the
`return` keyword:

``` bash
return "hello world"
```

Note that functions allow implicit returns,
so you don't need to explicitely use a `return`:

``` bash
func = f(x) {
    x + 1
}

func(9) # 10
```

## Next

That's about it for this section!

You can now head over to read about [if expressions](/syntax/if).