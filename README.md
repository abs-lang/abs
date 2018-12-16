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
  * pipes
  * interpolation
  * do not require semicolon at the end of a command
  * `$(sleep1; ls -la)` fails
  * allow to access the status code of a command with `comm = $(...); comm[status]` or `comm.status`
  * remove "\n" from return output
* "fix" hashes
  * hash key should be string
  * allow "false" json ({k: "v"}) where k is a literal string
* ~~interpreter code `abs test.abs`~~
* builds for interpreter
* add array std functions
* add string standard functions
* pipe operator
* floats
* for
* while
* foreach
* description of the language
  * assignments
  * expressions
  * functions
    * named
  * math
  * builtin
* license
* else if

### Later

* bash command syntax
  * special variable `$?` for BC
* add go native functions
* add hash std functions
* named functions
* until

## Credits

* [Terence Parr (ANTLR), for introducing me to parser generators](https://www.antlr.org/)
* [Thorsten Ball (interpreter book), for demystifying interpreters and providing the initial codebase for the abs interpreter](https://interpreterbook.com/)
* [Joe Jean, for suggesting the interpreter book](https://www.joejean.net/)
* [Bash, for being terrible at control flow ;-)](https://en.wikipedia.org/wiki/Bash_(Unix_shell))