# The abs programming language

Turn this:

``` bash
# Simple program that fetches your IP and sums it up
IP=$(curl -s 'https://api.ipify.org?format=json' | jq -r ".ip")
IFS=. read first second third fourth <<EOF
${IP##*-}
EOF
total=$((first + second + third + fourth))

if [ $1 -gt 100 ]
    then
    echo Hey that\'s a large number.
fi
```

into this:

``` bash
# Simple program that fetches your IP and sums it up
ip = $(curl -s 'https://api.ipify.org?format=json' | jq -r ".ip")

if ip.split(".").sum() > 100 {
    echo "Hey that's a large number."
}
```

## Why

## Description

## Running

## TODO

### 1.0

* bash command syntax
  * ~~basic command execution~~
  * ~~pipes~~
  * interpolation
  * do not require semicolon at the end of a command
  * `$(sleep1; ls -la)` fails
  * allow to access the status code of a command with `comm = $(...); comm[status]` or `comm.status`
  * ~~remove "\n" from echo output~~
* "fix" hashes
  * hash key should be string
  * allow "false" json ({k: "v"}) where k is a literal string
* ~~interpreter code `abs test.abs`~~
* builds for interpreter
* add array std functions (https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array / https://lodash.com/docs/4.17.11 -- filter by array methods)
* add string standard functions (https://golang.org/pkg/strings/)
* builtins
  * funcs
    * type
  * json
  * math (https://golang.org/pkg/math/)
  * rand (https://golang.org/pkg/math/rand/)
  * time (https://golang.org/pkg/math/)
* fix \" in strings
* pipe operator
* floats
* for
* while
* foreach
* ~~environment vars~~
* description of the language
  * assignments
  * expressions
  * functions
    * named
  * builtins
* license
* else if

### Later

* decide what to do with semicolons (either all in or ignore them)
* bash command syntax
  * special variable `$?` for BC
* add go native functions
* add hash std functions
* named functions
* until
* parsing errors with line nr etc
  * would be nice to link at doc

## Status

Early stage, so it could be that the language parser / evaluator might throw a bunch
of errors if you feed it funny code.

Open an issue and let's have fun!

## Credits

* [Terence Parr (ANTLR)](https://www.antlr.org/), for introducing me to parser generators
* [Thorsten Ball (interpreter book)](https://interpreterbook.com/), for demystifying interpreters and providing the initial codebase for the abs interpreter
* [Joe Jean](https://www.joejean.net/), for suggesting the interpreter book
* [Bash](https://en.wikipedia.org/wiki/Bash_(Unix_shell)), for being terrible at control flow ;-)