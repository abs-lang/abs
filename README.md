# The abs programming language

Turn this:

``` bash
# Simple program that fetches your IP and sums it up
IP=$(curl -s 'https://api.ipify.org?format=json' | jq -r ".ip")
IFS=. read first second third fourth <<EOF
${IP##*-}
EOF
total=$((first + second + third + fourth))
if [ $total -gt 100 ]
    echo "The sum of [$IP] is $total."
fi
```

into this:

``` bash
# Simple program that fetches your IP and sums it up
ip = $(curl -s 'https://api.ipify.org?format=json' | jq -rj ".ip");

total = ip.split(".").map(int).sum()
if total > 100 {
    echo("The sum of [%s] is %s.", ip, total)
}
```

## Why

## Status

Early stage, so it could be that the language parser / evaluator might throw a bunch
of errors if you feed it funny code.

Open an issue and let's have fun!

## Credits

* [Terence Parr (ANTLR)](https://www.antlr.org/), for introducing me to parser generators
* [Thorsten Ball (interpreter book)](https://interpreterbook.com/), for demystifying interpreters and providing the initial codebase for the abs interpreter
* [Joe Jean](https://www.joejean.net/), for suggesting the interpreter book
* [Bash](https://en.wikipedia.org/wiki/Bash_(Unix_shell)), for being terrible at control flow ;-)
