---
permalink: /misc/runtime
---

# Runtime

The ABS runtime lets you customize how ABS scripts are interpreted,
and exposes some useful global variables.

## ABS init file

When the ABS interpreter starts running, it will load an optional
ABS script as its init file. The ABS init file path can be
configured via the OS environment variable `ABS_INIT_FILE`. The
default value is `ABS_INIT_FILE=~/.absrc`.

If the `ABS_INIT_FILE` exists, it will be evaluated before the
interpreter begins in both interactive REPL or script modes.
The result of all expressions evaluated in the init file become
part of the ABS global environment which are available to command
line expressions or script programs.

Have a look at [an example ABS init file](https://github.com/abs-lang/abs/tree/master/examples/absrc.abs).

## ABS_INTERACTIVE

The `ABS_INTERACTIVE` global environment variable
is pre-set to `true` or `false` so that the init file can determine
which mode is running. This is useful if you wish to set the ABS REPL
command line prompt or history configuration variables in the init file.
This will preset the prompt and history parameters for the interactive
REPL (see [REPL Command History](/misc/configuring-the-repl#REPL_Command_History) above).

```
$ abs
Hello user, welcome to the ABS programming language!
Type 'quit' when you're done, 'help' if you get lost!
⧐  ABS_INTERACTIVE
true
```