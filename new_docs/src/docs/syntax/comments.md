---
permalink: /syntax/comments
---

# Comments

In accordance with ABS' minimalist style, there is only
one way to create comments in order to document your code,
by using the `#` character:

```bash
# This text represents a comment
x = 1
```

In ABS, when you start a comment it will run until the end
of the line, meaning there's no way to close comment blocks:

```bash
# This will not work # x = 1
x
ERROR: identifier not found: x
```

Though you can comment after a statement:

```bash
x = 1 # Now, this is a cool assignment!
```

## Next

That's about it for this section!

You can now head over to read about data types available in
ABS, starting from the [string](/types/string).
