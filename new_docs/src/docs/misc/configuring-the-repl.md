---
permalink: /misc/configuring-the-repl
---

# Configuring the REPL

Interactive REPL sessions can restore and save the command
history to a history file containing a maximum number of command lines.

The prompt live history is restored from the history file when
the REPL starts and then saved again when the REPL exits. This way you
can navigate through the command lines from your previous sessions
by using the up and down arrow keys at the prompt.

Note well that the live prompt history will show duplicate command
lines, but the saved history will only contain a single command
when the previous command and the current command are the same.

The history file name and the maximum number of history lines are
configurable through:

- the ABS environment (set by the ABS init file; see below)
- the OS environment
- The default values are `ABS_HISTORY_FILE="~/.abs_history"` and `ABS_MAX_HISTORY_LINES=1000`.

If you wish to suppress the command line history completely, just
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

## Configuring the ABS REPL Command Line Prompt

The ABS REPL command line prompt may be configured at start up using
`ABS_PROMPT_LIVE_PREFIX` and `ABS_PROMPT_PREFIX` variables from either
the ABS or OS environments. The default values are
`ABS_PROMPT_LIVE_PREFIX=false` and `ABS_PROMPT_PREFIX="⧐ "`.

REPL "static prompt" mode will be configured if `ABS_PROMPT_PREFIX`
contains no live prompt `template string` or if
`ABS_PROMPT_LIVE_PREFIX=false`. The `static prompt` will be the
value of the `ABS_PROMPT_PREFIX` string (if present) or the default
prompt `"⧐ "`. Static prompt mode is the default for the REPL.

REPL "live prompt" mode follows the current working directory
set by `cd()` when both `ABS_PROMPT_LIVE_PREFIX=true` and the
`ABS_PROMPT_PREFIX` variable contains a live prompt `template string`.

A live prompt `template string` may contain the following
named placeholders:

- `{user}`: the current userId
- `{host}`: the local hostname
- `{dir}`: the current working directory following `cd()`

For example, you can create a `bash`-style live prompt:

```bash
$ cat ~/.absrc
# ABS init script ~/.absrc
# For interactive REPL, override default prompt, history filename and size
if ABS_INTERACTIVE {
    ABS_PROMPT_LIVE_PREFIX = true
    ABS_PROMPT_PREFIX = "{user}@{host}:{dir}$ "
    ABS_HISTORY_FILE = "~/.abs_hist"
    ABS_MAX_HISTORY_LINES = 500
}

$ abs
Hello user, welcome to the ABS (1.1.0) programming language!
Type 'quit' when you are done, 'help' if you get lost!
user@hostname:~/git/abs$ cwd = cd()
user@hostname:~$ `ls .absrc`
.absrc
user@hostname:~$
```

## Next

That's about it for this section!

You can now head over to read [about the ABS runtime](/misc/runtime).
