---
permalink: /syntax/for
---

# For loops

ABS supports 2 types of `for` loops: the "standard" form and
the "in" one.

## Standard form

A standard `for` loop takes the canonical form:

```bash
for x = 0; x < 10; x = x + 1 {
    echo("Looping...")
}
```

It's important to understand scoping in this form:
if the loop declares an identifier that was already
declared earlier, it will be "temporarely" overwritten
in the loop, but will assume its original value
afterwards.

Code is better than 1000 words:

```bash
x = "hello world"

for x = 0; x < 10; x = x + 1 {
    # x is 0, 1, 2...
}

echo(x) # "hello world"
```

Similarly, a variable declared on the loop (not inside)
will cease to exist after the loop is done:

```bash
for x = 0; x < 10; x = x + 1 {
    # x is 0, 1, 2...
}

echo(x) # x is not defined here
```

Finally, variables declared inside the loop will instead
keep living afterwards:

```bash
for x = 0; x < 10; x = x + 1 {
    y = x
}

echo(y) # 9
```

## In form

The "in" form of the `for` loops allows you to iterate over
an array or an hash:

```bash
for x in [1, 2, 3] {
    # x is 1, 2, 3
}

for x in {"a": 1, "b": 2, "c": 3} {
    # x is 1, 2, 3
}
```

Both key and values are available in the loop:

```bash
for k, v in [1, 2, 3] {
    # k is 0, 1, 2
    # v is 1, 2, 3
}

for k, v in {"a": 1, "b": 2, "c": 3} {
    # k is a, b, c
    # v is 1, 2, 3
}
```

In terms of scoping, the "in" form follows the same rules
as the standard one, meaning that:

```bash
k = "hello world"

for k, v in [1, 2, 3] {
    # k is 0, 1, 2
    # v is 1, 2, 3
}

echo(k) # "hello world"
echo(v) # v is not defined
```

## break and continue

`break` and `continue` work just as you'd expect:
the former breaks out of a loop:

```bash
test = 0
for x = 0; x <= 10; x = x + 1 {
  if x < 10 {
    break
  }

  test += x
}

test # 0
```

while the later skips to the next execution of the loop:

```bash
test = 0
for x = 0; x <= 10; x = x + 1 {
  if x < 10 {
    continue
  }

  test += x
}

test # 10
```

## For ... else ...

`For` loops can also have `else` clause which executes if
the list in the `for` condition is empty.

For example, when we run a database query like the following,
if the "users" list is empty, we will only execute the statements
inside the `else` clause.

```bash
users = db.query("SELECT students WHERE age > 20")

for user in users {
  print(user)
} else {
  print("We don't have students above the age of 20")
}
```
