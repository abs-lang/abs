<p align="center">
  <a href="https://www.abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/abs-lang/abs/blob/master/bin/abs-horizontal.png?raw=true">
  </a>
</p>

# How to run ABS code

In order to run programs written in abs, you can simply download
the latest release of ABS from Github, and dump the ABS executable
in your `PATH`.

Afterwards, you can run ABS scripts with:

``` bash
$ abs path/to/scripts.abs
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

## Why is abs interpreted?

ABS' goal is to be a portable, pragmatic, coincise, simple language:
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