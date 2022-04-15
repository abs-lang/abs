---
permalink: /stdlib/runtime
---

# @runtime

The `@runtime` module provides information about the ABS
runtime that is currently executing the script.

## API

```py
runtime = require('@runtime')
```

### @runtime.version

Returns the version of the ABS interpreter:

```py
runtime.version # "x.y.z" eg. 1.0.0
```

### @runtime.name

Returns the name of the runtime:

```py
runtime.name # "abs"
```
