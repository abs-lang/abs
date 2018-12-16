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
IP=curl -s 'https://api.ipify.org?format=json' | jq -r ".ip"
total = parts | split(".") | sum

if total > 100 {
    echo "Hey that's a large number."
}
```

## Why

## Description

## Running

## TODO

* ~~remove `let`~~
* remove parens from ifs
* bash command syntax
* fix hashes
* add array std functions
* add go native functions
* add hash std functions
* add string standard functions
* pipe operator
* interpreter code `abs test.abs`

## Credits

* [Terence Parr (ANTLR), for introducing me to parser generators](https://www.antlr.org/)
* [Thorsten Ball (interpreter book), for demystifying interpreters and providing the initial codebase for the abs interpreter](https://interpreterbook.com/)
* [Joe Jean, for suggesting the interpreter book](https://www.joejean.net/)
* [Bash, for being terrible at control flow ;-)](https://en.wikipedia.org/wiki/Bash_(Unix_shell))