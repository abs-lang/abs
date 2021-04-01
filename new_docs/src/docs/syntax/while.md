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

## Next

That's about it for this section!

You can now head over to read about [system (or shell) commands](/syntax/system-commands).
