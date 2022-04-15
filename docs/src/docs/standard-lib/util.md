---
permalink: /stdlib/util
---

# @util

The `@util` module provides various utilities.

## API

```py
util = require('@util')
```

### @util.memoize(ttl)

A decorator that allows memoization of a function based on
the arguments passed to it:

```py
@util.memoize(60)
f expensive_task(x, y, z) {
    # do something very expensive here...
}
```

The first time `expensive_task` gets called with a set of arguments,
it will be execute. The next time it's called with the same set of
arguments its result will be fetched from a cache (currently
implemented in-memory). Executions are going to be cached for a
specific timeframe, `ttl` (in seconds).

Arguments are serialized using the `str()` method:

```bash
[12, {}, 0.23, "hello"].str() # "[12, {}, 0.23, \"hello\"]"
```
