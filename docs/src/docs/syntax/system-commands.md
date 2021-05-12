---
permalink: /syntax/system-commands
---

# System (shell) commands

Executing system commands is one of the most important features
of ABS, as it allows mixing the convenience of the shell with
the syntax of a modern programming language.

Commands are executed with the `` `command` `` syntax,
which resemble Bash's syntax to execute commands in a subshell:

```bash
date = `date` # "Sun Apr 1 04:30:59 +01 1995"
```

As you can see, the return value of a command is a simple
string -- the output of the program. If the program was to
encounter an error, the same string would hold the error
message:

```bash
date = `dat` # "bash: dat: command not found"
```

It would be fairly painful to have to parse strings
manually to understand if a command executed without errors;
in ABS, the returned string has a special property `ok` that
checks whether the command was successful:

```js
if `ls -la`.ok {
    echo("hello world")
}
```

## Executing commands in background

Sometimes you might want to execute a command in
background, so that the script keeps executing
while the command is running. In order to do so,
you can simply add an `&` at the end of your script:

```bash
`sleep 10 &`
echo("This will be printed right away!")
```

You might also want to check whether a command
is "done", by checking the boolean `.done` property:

```bash
cmd = `sleep 10 &`
cmd.done # false
`sleep 11`
cmd.done # true
```

If, at some point, you want to wait for the command
to finish before running additional code, you can
use the `wait` method:

```bash
cmd = `sleep 10 &`
echo("This will be printed right away!")
cmd.wait()
echo("This will be printed after 10s")
```

If you ever want to terminate a running command, you can
use the `kill` method.

```bash
cmd = `sleep 10 &`
cmd.done # false
cmd.kill()
cmd.done # true
```

Also note that when an `exec()` command string terminates with an `&`,
the `exec(command)` function will terminate immediately after launching
the command which will run independently in the background.
This means that the command must either terminate on its own or be killed
using `pkill` or similar. This way an ABS script can launch a true daemon
process that may operate on its own outside of ABS. For example you can
reboot a remote computer via ssh without interacting with it:

```bash
exec("ssh user@host.local 'sudo reboot' &")
```

## Interpolation

You can also replace parts of the command with variables
declared within your program using the `$` symbol:

```bash
file = "cpuinfo"
x = `cat /proc/$file`
echo(x) # processor: 0\nvendor_id: GenuineIntel...
```

or interpolation within an `exec(command)`

```bash
cmd = args(2)
filename = args(3)
exec("sudo $cmd $filename")
```

and if you need `$` literals in your command, you
simply need to escape them with a `\`:

```bash
`echo $PWD` # "" since the ABS variable PWD doesn't exist
`echo \$PWD` # "/go/src/github.com/abs-lang/abs"
```

## Using a different shell

By default, ABS uses `bash -c` to execute commands; on Windows
it instead uses `cmd.exe /C`.

You can specify which shell to use by setting the environment
variable `ABS_COMMAND_EXECUTOR`:

```sh
`echo \$0` # bash
env("ABS_COMMAND_EXECUTOR", "sh -c")
`echo \$0` # sh
```

## Alternative \$() syntax

Even though the use of backticks is the standard recommended
way to run system commands, for the ease of embedding ABS also
allows you to use the `$(command)` syntax:

```
$(basename $(dirname "/tmp/make/life/easy")) // "easy"
```

Commands that use the `$()` syntax need to be
on their own line, meaning that you will not
be able to have additional code on the same line.
This will throw an error:

```bash
$(sleep 10); echo("hello world")
```

## Executing commands without capturing I/O

It is also possible to execute a shell command without capturing its
input or output using the `exec(command)` function. This allows long running
or interactive programs to be run using the terminal's Standard IO
(stdin, stdout, stderr). For example:

```bash
exec("sudo visudo")
```

would open the default text editor in super user mode on the /etc/sudoers file.

Unlike the normal backtick command execution syntax above,
the `exec(command)` function call does not return a result string unless it fails.
Therefore, the `exec(command)` may be the last command executed in a script
file leaving the executed command in charge of the terminal IO until it
terminates.

For example, an ABS script might be used to marshall the command line args
for an interactive program such as the nano editor:

```bash
$ cat abs/tests/test-exec.abs
# marshall the args for the nano editor
# if the filename is not given in the args, prompt for it
# if the file is located outside the user's home dir, invoke sudo nano filename

cmd = 'nano'
filename = arg(2)
homedir = env("HOME")

while filename == '' {
    echo("Please enter file name for %s: ", cmd)
    filename = stdin()
}

if filename.prefix('~/') || filename.prefix(homedir) {
    sudo = ''
} else {
    sudo = 'sudo'
}

# execute the command with live stdIO
exec("$sudo $cmd $filename")
```