---
permalink: /syntax/return
---

# Returning values

We promise, this is going to be short!

Returning values is done through the
`return` keyword:

```bash
return "hello world"
```

Note that functions allow implicit returns,
so you don't need to explicitely use a `return`:

```bash
func = f(x) {
    x + 1
}

func(9) # 10
```

The default value of a `return` is `null`:

```bash
if x {
    return # null
}
```
