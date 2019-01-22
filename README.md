<p align="center">
  <a href="https://www.abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/abs-lang/abs/blob/master/bin/ABS.png?raw=true" width="310">
  </a>
</p>

<p align="center">
  <a href="https://travis-ci.com/abs-lang/abs"><img alt="Travis Status" src="https://travis-ci.com/abs-lang/abs.svg?branch=master"></a>
  <a href="https://github.com/abs-lang/abs"><img alt="License" src="https://img.shields.io/github/license/abs-lang/abs.svg"></a>
  <a href="https://github.com/abs-lang/abs"><img alt="Version" src="https://img.shields.io/github/release-pre/abs-lang/abs.svg"></a>
</p>

# The ABS programming language

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

Let's try to fetch our IP address and print the sum of its
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

## Documentation

Visit [abs-lang.org](https://www.abs-lang.org)

## Contributing

Wanna hack on ABS locally? The recommended development
environment is inside a Docker container -- simply:

* `make build` (builds the container)
* `make run` (sends you inside the development container)
* `make test` (runs the abs tests)

After you make any change, simply run `make test` and check
if any errors pops up. If everything looks fine that means
you're ready to [open a pull request](https://github.com/abs-lang/abs/pulls)!

## Status

ABS is fresh and under active development, meaning exciting
things happen on a weekly basis.

Have a look at the roadmaps [here](https://github.com/abs-lang/abs/milestones):
to know of what version we're currently working on take a look at [abs-lang.org/misc/technical-details](https://www.abs-lang.org/misc/technical-details).
