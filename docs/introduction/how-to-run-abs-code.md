# How to run ABS code

In order to run programs written in abs, you can simply download
the latest release of ABS from Github, and dump the ABS executable
in your `PATH`. Windows, OSX and a few Linux flavors are supported.

We also provide a 1-command installer that should work across
platforms:

``` bash
bash <(curl https://www.abs-lang.org/installer.sh)
```

and will download the `abs` executable in your current
directory -- again, we recommend to move it to your `$PATH`.

Afterwards, you can run ABS scripts with:

``` bash
$ abs path/to/scripts.abs
```
You can also run an executable abs script directly from bash
using a bash shebang line at the top of the script file. 

In this example the abs executable is linked to `/usr/local/bin/abs`
and the abs script `~/bin/remote.abs` has its execute permissions set.
```bash
$ cat ~/bin/remote.abs
#! /usr/local/bin/abs
# remote paths are <target>:<path> 
from_path = arg(2) 
to_path = arg(3)
if ! (from_path && to_path) {
    if ! from_path {from_path = "<missing>"}
    if ! to_path {to_path = "<missing>"}
    echo("FROM: %s, TO: %s", from_path, to_path)
    exit(1)
}
...
# the executable abs script above is in the PATH at ~/bin/remote.abs
$ remote.abs
FROM: <missing>, TO: <missing>
```

Scripts do not have to have a specific extension,
although it's recommended to use `.abs` as a
convention.

A bit lost right now? We'd suggest to clone [ABS' main repository](https://github.com/abs-lang/abs) as you can already
start testing some code with the scripts in the
[examples](https://github.com/abs-lang/abs/tree/master/examples) directory.

## REPL

If you want to get a more *live* feeling of ABS, you can
also simply run the interpreter; without any argument. It
will launch ABS' REPL, and you will be able to test code on
the fly:

``` bash
$ abs
Hello there, welcome to the ABS programming language!
Type 'quit' when you're done, 'help' if you get lost!
⧐  ip = $(curl icanhazip.com)
⧐  ip.ok
true
⧐  ip()
ERROR: not a function: STRING
⧐  ip
94.204.178.37
```
### REPL Command History

Interactive REPL sessions can restore and save the command 
history to a history file containing a maximum number of command lines. 

The prompt live history is restored from the history file when
the REPL starts and then saved again when the REPL exits. This way you
can navigate through the command lines from your previous sessions
by using the up and down arrow keys at the prompt.

+ Note well that the live prompt history will show duplicate command
lines, but the saved history will only contain a single command
when the previous command and the current command are the same.

The history file name and the maximum number of history lines are
configurable through 
1) the ABS environment (set by the ABS init file; see below)
2) the OS environment
3) The default values are `ABS_HISTORY_FILE="~/.abs_history"` 
and `ABS_MAX_HISTORY_LINES=1000`.

+ If you wish to suppress the command line history completely, just 
set `ABS_MAX_HISTORY_LINES=0`. In this case the history file
will not be created.

For example:
```bash
$ export ABS_HISTORY_FILE="~/my_abs_hist"
$ export ABS_MAX_HISTORY_LINES=500
$ abs
Hello user, welcome to the ABS (1.1.0) programming language!
Type 'quit' when you are done, 'help' if you get lost!
⧐  pwd()
/home/user/git/abs
⧐  cd()
/home/user
⧐  echo("hello")
hello
⧐  quit
Adios!

$ cat ~/my_abs_hist`; echo
pwd()
cd()
echo("hello")
$
```

## ABS Init File

When the ABS interpreter starts running, it will load an optional
ABS script as its init file. The ABS init file path can be 
configured via the OS environment variable `ABS_INIT_FILE`. The
default value is `ABS_INIT_FILE=~/.absrc`.

If the `ABS_INIT_FILE` exists, it will be evaluated before the
interpreter begins in both interactive REPL or script modes.
The result of all expressions evaluated in the init file become
part of the ABS global environment which are available to command
line expressions or script programs.

Also, note that the `ABS_INTERACTIVE` global environment variable
is pre-set to `true` or `false` so that the init file can determine
which mode is running. This is useful if you wish to set the ABS
prompt or history configuration variables in the init file. This
will preset the prompt and history parameters for the interactive REPL. 
See [REPL Command History](#REPL_Command_History) above.

The REPL prompt may be configured using `ABS_PROMPT_LIVE_PREFIX` and
`ABS_PROMPT_PREFIX` from these ABS or OS environment variables. 
The live prompt follows the current working directory set by `cd()`
when it is enabled. The prompt prefix replaces the default prefix
which follows the live current working directory if it is enabled.

For example, you can create a bash-style prompt: 
```bash
$ cat ~/.absrc
# ABS init script ~/.absrc 
# For interactive REPL, override default prompt, history filename and size
if ABS_INTERACTIVE {
    ABS_PROMPT_LIVE_PREFIX = true
    ABS_PROMPT_PREFIX = "$ "
    ABS_HISTORY_FILE = "~/.abs_hist"
    ABS_MAX_HISTORY_LINES = 500
}
$ abs
Hello user, welcome to the ABS (1.1.0) programming language!
Type 'quit' when you are done, 'help' if you get lost!
/home/user/git/abs$ cwd = cd()
/home/user$ `ls .absrc`
.absrc
/home/user$ 
```

Also see a `template ABS Init File` at [examples](https://github.com/abs-lang/abs/tree/master/examples/absrc.abs).

## Why is abs interpreted?

ABS' goal is to be a portable, pragmatic, concise, simple language:
great performance comes second.

With this in mind, we made a deliberate choice to avoid
compiling ABS code, as it would require additional complexity
in the codebase, with very little benefits. Tell us, when
was the last time you were interested in how many milliseconds
it took to run a Bash script?

## Next

That's about it for this section!

You can now head over to read about ABS's syntax,
starting with [assignments](/syntax/assignments)!
