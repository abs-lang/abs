# Quickstart

In order to run programs written in ABS, you can simply [download
the latest release from Github](https://github.com/abs-lang/abs/releases)
and dump the executable in your `PATH`. Windows, OSX and a few Linux flavors are supported.

We also provide a 1-command installer that should work across
platforms:

```bash
bash <(curl https://www.abs-lang.org/installer.sh)
```

and will download the `abs` executable in your current
directory -- again, we recommend to move it to your `$PATH`.

Afterwards, you can run ABS scripts with:

```bash
$ abs path/to/script.abs
```

Scripts do not need a specific extension,
although it's recommended to use `.abs` as a
convention: we may reserve some keywords in the
future (such as `abs version` or `abs install`)
so we recommend to attach an extension to the
scripts you're trying to run.

## REPL

If you want to get a more "live" feeling of ABS, you can
also simply run the interpreter; without any argument. It
will launch ABS' REPL, and you will be able to test code on
the fly:

```
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

## ABS from bash

You can also run an executable abs script directly from bash
using a bash shebang.

In this example the abs executable is linked to `/usr/local/bin/abs`
and the script `~/bin/remote.abs` is executable (`chmod +x`):

```bash
$ cat ~/bin/hello.abs
#! /usr/local/bin/abs
echo("Hello world!")
...

# the executable abs script above is in the PATH at ~/bin/hello.abs
$ hello.abs
Hello world!
```

## Explore the docs!

A bit lost right now? Here's what we suggest you do:

* explore the [docs](/docs) to learn more about ABS' features 
* try running some ABS code in our browser-based [playground](/playground)
* check some of the [example scripts](https://github.com/abs-lang/abs/tree/master/examples) in our official repo
