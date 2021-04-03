---
permalink: /stdlib/cli
---

# @cli

The `@cli` module provides a simple interface to
easily build CLI applications.

## API

```py
cli = require('@cli')
```

### @cli.cmd(name, description, default_flags)

A decorator that registers a command to be executed via CLI:

```py
@cli.cmd("date", "prints the current date", {format: ''})
f date(args, flags) {
    format = flags.format
    return `date ${format}`
}
```

The `name` of the command matches the command entered
by the user on the CLI (eg. `$ ./cli my_command`), while
the description is used when printing a help message,
available through `$ ./cli help`.

The user can pass however many arguments and flags they want,
and they will be passed on to the function defined as command.
For example, when the user issues the command `$ ./cli a b c --flag 25`,
args will be `["a", "b", "c", "--flag", "25"]` and flags will
be `{"flag": "25"}`. Default flags are there so that, if a flag
is not passed, it will be populated with a default value.

### @cli.run()

Runs the CLI application:

```py
cli.run()
```

The application will, by default, have an `help` command
that lists all available commands and is called if no command
is provided.

### @cli.repl()

Runs the CLI application in interactive mode:

```py
cli.repl()
```

The application will, by default, have an `help` command
that lists all available commands.

## Example CLI app

Here is an example app with 3 commands:

- the default `help`
- `date`, to print the current date in a specific format
- `ip`, to fetch our IP address

```py
#!/usr/bin/env abs
cli = require('@cli')

@cli.cmd("ip", "finds our IP address", {})
f ip_address(arguments, flags) {
    return `curl icanhazip.com`
}

@cli.cmd("date", "Is it Friday already?", {"format": ""})
f date(arguments, flags) {
    format = flags.format
    return `date ${format}`
}

cli.run()
```

You can save this script as `./cli` and make it executable
with `chmod +x ./cli`. Then you will be able to use the CLI
app:

```
$ ./cli
Available commands:

  * date - Is it Friday already?
  * help - print this help message
  * ip - finds our IP address

$ ./cli help
Available commands:

  * date - Is it Friday already?
  * help - print this help message
  * ip - finds our IP address

$ ./cli ip
87.201.252.69

$ ./cli date
Sat Apr  4 18:06:35 +04 2020

$ ./cli date --format +%s
1586009212
```

## Example REPL app

Here is an example app with 4 commands:

- the default `help`
- `count`, which prints a counter
- `incr`, which increments a counter by 1
- `incr_by`, which increments a counter by a number specified by the user

```py
#!/usr/bin/env abs
cli = require('@cli')

res = {"count": 0}

@cli.cmd("count", "prints a counter", {})
f counter(arguments, flags) {
    echo(res.count)
}

@cli.cmd("incr", "Increment our counter", {})
f incr(arguments, flags) {
    res.count += 1
    return "ok"
}

@cli.cmd("incr_by", "Increment our counter", {})
f incr_by(arguments, flags) {
    echo("Increment by how much?")
    n = stdin().number()
    res.count += n
    return "ok"
}

cli.repl()
```

You can save this script as `./cli` and make it executable
with `chmod +x ./cli`. Then you will be able to use the CLI
app:

```
$ ./cli
help
Available commands:

  * count - prints a counter
  * help - print this help message
  * incr - Increment our counter
  * incr_by - Increment our counter
count
0
incr
ok
incr
ok
count
2
incr_by
Increment by how much?
-10
ok
count
-8
```