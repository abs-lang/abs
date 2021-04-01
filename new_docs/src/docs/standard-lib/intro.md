---
permalink: /stdlib/intro
---

# Standard library

ABS comes bundled not only with an array of builtin types and functions,
but also with a few modules to ease your development process: we refer
to this modules as the "standard library", "standard modules" or `stdlib`.

## Requiring a standard library module

Standard library modules are required the same way you'd require
any other external module, by using the `require(...)` function;
the only difference is that standard modules use the `@` character
as prefix, for example:

```py
mod = require('@module')  # Loads "module" from the standard library
mod = require('./module') # Loads "module" from the current directory
mod = require('module')   # Loads "module" that was installed through the ABS package manager
```

## Technical details

The ABS standard library is developed in ABS itself and available
for everyone to see (and poke with) at [github.com/abs-lang/abs/tree/master/stdlib](https://github.com/abs-lang/abs/tree/master/stdlib).

The `@cli` library, for example, is a simple ABS script of less
than 100 lines of code:

```bash
# CLI app
cli = {}

# Commands registered within this CLI app
cli.commands = {}

# Function used to register a command
cli.cmd = f(name, description, flags) {
    return f(fn) {
        cli.commands[name] = {};
        cli.commands[name].cmd = f() {
            # Create flags with default values
            for k, _ in flags {
                v = flag(k)

                if v {
                    flags[k] = v
                }
            }
            
            # Call the original cmd
            result = fn.call([args()[3:], flags])

            # If there's a result we print it out
            if result {
                echo(result)
            }
        }

        cli.commands[name].description = description
    }
}

# Run the CLI app
cli.run = f() {
    # ABS sees "abs script.abs xyz"
    # so the command is the 3rd argument
    cmd = arg(2)

    # Not passing a command? Let's print the help
    if !cmd {
        return cli.commands['help'].cmd()
    }

    if !cli.commands[cmd] {
        exit(99, "command '${cmd}' not found")
    }

    return cli.commands[cmd].cmd()
}

# Add a default help command that can be
# overridden by the caller
@cli.cmd("help", "print this help message", {})
f help() {
    echo("Available commands:\n")

    for cmd in cli.commands.keys().sort() {
        s = "  * ${cmd}"

        if cli.commands[cmd].description {
            s += " - " + cli.commands[cmd].description
        }

        echo(s)
    }
}

return cli
```

## Next

That's about it for this section!

You can now explore ABS' first standard library module, the [@runtime](/stdlib/runtime).