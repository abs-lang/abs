---
permalink: /introduction/how-to-run-abs-code
---

# How to run ABS code

In order to run programs written in abs, you can simply download
the latest release of ABS from Github, and dump the ABS executable
in your `PATH`. Windows, OSX and a few Linux flavors are supported.

We also provide a 1-command installer that should work across
platforms:

```bash
bash <(curl https://www.abs-lang.org/installer.sh)
```

and will download the `abs` executable in your current
directory -- again, we recommend to move it to your `$PATH`.

Afterwards, you can run ABS scripts with:

```bash
$ abs path/to/scripts.abs
```

You can also run an executable abs script directly from bash
using a bash shebang line at the top of the script file.

In this example the abs executable is linked to `/usr/local/bin/abs`
and the abs script `~/bin/remote.abs` has its execute permissions set.

```bash
$ cat ~/bin/hello.abs
#! /usr/local/bin/abs
echo("Hello world!")
...

# the executable abs script above is in the PATH at ~/bin/hello.abs
$ hello.abs
Hello world!
```

Scripts do not have to have a specific extension,
although it's recommended to use `.abs` as a
convention: we may reserve some keywords in the
future (such as `abs version` or `abs install`)
so we recommend to attach an extension to the
scripts you're trying to run.

A bit lost right now? We'd suggest to clone [ABS' main repository](https://github.com/abs-lang/abs) as you can already
start testing some code with the scripts in the
[examples](https://github.com/abs-lang/abs/tree/master/examples) directory.

## REPL

If you want to get a more _live_ feeling of ABS, you can
also simply run the interpreter; without any argument. It
will launch ABS' REPL, and you will be able to test code on
the fly:

```bash
$ abs
Hello there, welcome to the ABS programming language!
Type 'quit' when you're done, 'help' if you get lost!
⧐  ip = `curl icanhazip.com`
⧐  ip.ok
true
⧐  ip()
ERROR: not a function: STRING
⧐  ip
94.204.178.37
```

## Next

That's about it for this section!

You can now head over to try ABS directly in your
browser, on the [playground](/playground)!
