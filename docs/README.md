<p align="center">
  <a href="https://www.abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/abs-lang/abs/blob/master/bin/abs-horizontal.png?raw=true">
  </a>
</p>

<p align="center">
  <a href="https://travis-ci.com/abs-lang/abs"><img alt="Travis Status" src="https://travis-ci.com/abs-lang/abs.svg?branch=master"></a>
  <a href="https://github.com/abs-lang/abs"><img alt="License" src="https://img.shields.io/github/license/abs-lang/abs.svg"></a>
  <a href="https://github.com/abs-lang/abs"><img alt="Version" src="https://img.shields.io/github/release-pre/abs-lang/abs.svg"></a>
</p>

ABS is a scripting language that works best when you're on
your terminal. It tries to combine the elegance of languages
such as Python, or Ruby, to the convenience of Bash.

``` bash
tz = $(cat /etc/timezone);
continent, city = tz.split("/")

echo("Best city in the world?")

selection = stdin()

if selection == city {
  echo("You might be biased...")
}
```

See it in action:

[![asciicast](https://asciinema.org/a/218909.svg)](https://asciinema.org/a/218909)

Let's now try to fetch our IP address and print the sum of its
parts, if its higher than 100. Here's how you could do it
in Bash:

``` bash
# Simple program that fetches your IP and sums it up
RES=$(curl -s 'https://api.ipify.org?format=json' || "ERR")

if [ "$RES" = "ERR" ]; then
    echo "An error occurred"
    exit 1
fi

IP=$(echo $RES | jq -r ".ip")
IFS=. read first second third fourth <<EOF
${IP##*-}
EOF

total=$((first + second + third + fourth))
if [ $total -gt 100 ]; then
    echo "The sum of [$IP] is a large number, $total."
fi
```

And here's how you could write the same code in ABS:

``` bash
# Simple program that fetches your IP and sums it up
res = $(curl -s 'https://api.ipify.org?format=json');

if !res.ok {
  echo("An error occurred: %s", res)
  exit(1)
}

ip = res.json().ip
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

You can also install ABS with the 1-command installer:

``` bash
sh <(curl https://www.abs-lang.org/installer.sh)
```

*(you might need to sudo right before that)*

## Table of contents

## Introduction

* [Why another scripting language?](/introduction/why-another-scripting-language)
* [How to run ABS code](/introduction/how-to-run-abs-code)

## Syntax

* [Assignments](/syntax/assignments)
* [return](/syntax/return)
* [if](/syntax/if)
* [for](/syntax/for)
* [while](/syntax/while)
* [System commands](/syntax/system-commands)
* [Operators](/syntax/operators)
* [Comments](/syntax/comments)

## Types and functions

* [String](/types/string)
* [Number](/types/number)
* [Array](/types/array)
* [Hash](/types/hash)
* [Functions](/types/function)
* [Builtin functions](/types/builtin-function)

## Miscellaneous

* [Errors](/misc/error)
* [A few technical details...](/misc/technical-details)
* [Credits](/misc/credits)
