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

## Example CLI app

Here is an example app with 3 commands:

* the default `help`
* `date`, to print the current date in a specific format
* `ip`, to fetch our IP address

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

## Next

That's about it for this section!

You can now head over to read a little bit about [the util module](/stdlib/util).