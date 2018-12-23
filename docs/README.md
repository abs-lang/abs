<p align="center">
  <a href="https://abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/odino/abs/blob/master/bin/ABS.png?raw=true" width="310">
  </a>
</p>

ABS is a scripting language that works best when you're on
your terminal. It tries to combine the elegance of languages
such as Python, or Ruby, to the convenience of Bash.

Let's try to fetch our IP address and print the sum of its
parts, if its higher than 100. Here's how you could do it
in bash:

``` bash
# Simple program that fetches your IP and sums it up
IP=$(curl -s 'https://api.ipify.org?format=json' | jq -r ".ip")
IFS=. read first second third fourth <<EOF
${IP##*-}
EOF
total=$((first + second + third + fourth))
if [ $total -gt 100 ]
    echo "The sum of [$IP] is a large number, $total."
fi
```

And here's how you would write the same code in ABS:

``` bash
# Simple program that fetches your IP and sums it up
ip = $(curl -s 'https://api.ipify.org?format=json' | jq -rj ".ip");

total = ip.split(".").map(int).sum()
if total > 100 {
    echo("The sum of [%s] is a large number, %s.", ip, total)
}
```

Wondering how you can run this code? Simply grab the latest
[release](https://github.com/abs-lang/abs/releases) and run:

```
$ abs script.abs
```

## Table of contents

## Introduction

* [Why another scripting language?](/introduction/why-another-scripting-language)
* [How to run ABS code](/introduction/how-to-run-abs-code)

## Syntax

* [Assignments](/syntax/assignments)
* [return]((/syntax/return))
* [if]((/syntax/if))
* [for]((/syntax/for))
* [while](/syntax/while)
* [System commands](/syntax/system-commands)
* [Operators](/syntax/operators)

## Types and functions

* String
* Integer
* Array
* Hash
* Functions

## Miscellaneous

* A few technical details...
* Credits
