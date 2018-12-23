<p align="center">
  <a href="https://www.abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/abs-lang/abs/blob/master/bin/abs-horizontal.png?raw=true">
  </a>
</p>

# A few technical details...

The ABS interpreter is built with Golang version `1.11`, and is mostly based
on [the interpreter book](https://interpreterbook.com/) written by
[Thorsten Ball](https://twitter.com/thorstenball).

ABS is extremely different from Monkey, the "fictional" language the reader
builds throughout the book, but the base structure (lexer, parser, evaluator)
are still very much based on Thorsten's work.

## Development & contributing

The best way to start developing *for* ABS is to clone the repository
and run a `make build`: this will build a docker container with all
the necessary dependencies for developing locally (for experienced
Gophers: you might want to skip this altogether as your environment
will probably work perfectly).

With `make run` you can get inside a container built for ABS'
development, and `make test` will run all tests.

## Roadmap

We're currently working on the [preview-1 milestone](https://github.com/abs-lang/abs/milestone/1).

## Next

That's about it for this section!

You can now head over to read ABS' [credits](/misc/credits).