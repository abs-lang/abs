<p align="center">
  <a href="https://www.abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/abs-lang/abs/blob/master/bin/abs-horizontal.png?raw=true">
  </a>
</p>

# System (shell) commands

Executing system commands is one of the most important features
of ABS, as it allows the mixing of conveniency of the shell with
the syntax of a modern programming language.

Commands are executed with `$()`, which resembles Bash's
syntax to execute commands in a subshell:

``` bash
date = $(date) # "Sun Apr 1 04:30:59 +01 1995"
```

As you can see, the return value of a command is a simple
string -- the output of the program. If the program was to
encounter an error, the same string would hold the error
message:

``` bash
date = $(dat) # "bash: dat: command not found"
```

It would be fairly painful to have to parse strings
manually to understand if a command executed without errors;
in ABS, the returned string has a special method `ok` that
checks whether the command was successful:

``` bash
ls = $(ls -la)

if ls.ok() {
    echo("hello world")
}
```

You can also replace parts of the command with variables
declared within your program using the `$` symbol:

``` bash
file = "cpuinfo"
x = $(cat /proc/$file)
echo(x) # processor: 0\nvendor_id: GenuineIntel...
```

Currently, commands need to be on their own line, meaning
that you will not be able to have additional code
on the same line. This will throw an error:

``` bash
$(sleep 10); echo("hello world")
```

Note that this is currently a limitation that will likely
be removed in the future (see [#41](https://github.com/abs-lang/abs/issues/41)).

Commands are blocking and cannot be run in parallel, although
we're planning to support background execution in the future
(see [#70](https://github.com/abs-lang/abs/issues/70)).

Also note that, currently, the implementation of system commands
requires the `bash` executable to [be available on the system](https://github.com/abs-lang/abs/blob/5b5b0abf3115a5dd4dfe8485501f8765985ad0db/evaluator/evaluator.go#L696-L722).
Future work will make it possible to select which shell to use,
as well as bypassing the shell altogether (see [#73](https://github.com/abs-lang/abs/issues/73)).

## Next

That's about it for this section!

You can now head over to read about [operators](/syntax/operators).