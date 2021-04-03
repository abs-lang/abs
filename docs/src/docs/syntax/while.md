---
permalink: /syntax/while
---

# While

While loops are a special form of `for` loops, so much
that in some languages the canonical way to execute a
while loop is with a `for(;;)`.

ABS, though, has a dedicated construct:

```bash
x = 0

while x < 100 {
    x = x + 1
}

echo(x) # 99
```
