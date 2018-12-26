# Why another scripting language?

If you're a developer, chances are that you got your hands
on a server at some point during your career. If that statement
is true, chances are that server was running a flavor of Linux,
meaning you've probably encountered a terminal before, where
you most likely had to use Bash to automate some task with a
good old `.sh` script.

<p align="center">
    <img alt="terminal" src="https://github.com/abs-lang/abs/blob/master/bin/terminal.png?raw=true">
</p>

You might have felt, though, that Bash was a fairly strange
language, with a very uncommon syntax:

``` bash
if [ -z $STRING ]; then
    ...
fi
```

(*if you're wondering, the above snippet would check whether the variable
`$STRING` is an empty string*)

Far for bashing Bash (pun intended) or the generic [shell command language](http://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html),
we believe there should be a more straightforward alternative
to automating tasks, something that [Bash excels at](https://www.quora.com/What-are-the-main-advantages-of-Bash-as-a-programming-language): if there's reason why the pragmatic Python or the elegant Ruby haven't been able to overcome Bash as the *de-facto* standard
for shell scripting, that would be the inner simplicity of bash.
Running programs in parallel, interacting with the underlying system,
ease of portability...   ...these are quick and easy wins when you're
writing those `.sh` files.

As both general-purpose programming languages fans and shell lovers,
we believe there could be an alternative where a programmer would
combine the syntax and flexibility of general-purpose languages
(Python, Ruby and JS, to name a few) with the benefits of Bash.

This is why we developed the ABS programming language: a
language that is a joy to work with in the context of shell scripting:
it isn't here to replace the likes of PHP, Java or Python,
neither it wants to diminish the importance of Bash.

ABS tries to mix a more modern language with the
simplicity of Bash.

Let's take a look a look at some practical ABS code. We will now
call the API of nba.com in order to retrieve the stats for
one of last year's NBA games:

``` bash
r = $(curl "http://data.nba.net/prod/v1/20170201/0021600732_boxscore.json" -H 'DNT: 1' -H 'Accept-Encoding: gzip, deflate, sdch' -H 'Accept-Language: en' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36' -H 'Accept: */*' -H 'Referer: http://stats.nba.com/' -H 'Connection: keep-alive' --compressed);

if !r.ok {
    echo("Could not fetch game data. Bummer!")
    exit(1)
}

doc = r.json()

arena = doc.basicGameData.arena.name
city = doc.basicGameData.arena.city

echo("The game was played at the %s in %s", arena, city)

highlight = doc.basicGameData.nugget.text
if highlight.len() {
    echo("The press said: \"%s\"", highlight)
}

# The game was played at the TD Garden in Boston
# The press said: "Thomas scores 19 of 44 points in 4th quarter"
```

You will notice 3 things:

* [Isiah Thomas](https://en.wikipedia.org/wiki/Isaiah_Thomas_(basketball)) seems to be a really good player
* you should be very familiar with the above syntax
* the language is capable of seamlessly throwing shell commands into the mix

This is exactly why ABS was born: a familiar syntax, and the convenience of Bash.

A sneak-peek at some of the things ABS can elegantly do:

``` bash
# Unix pipes work
ip = $(curl icanhazip.com | tr -d '\n')

# We now have a string -> "10.10.10.12"
echo(ip)

# Let's play with it -> [10, 10, 10, 12]
parts = ip.split(".").map(int)

# 42 anyone?
echo(parts.sum())
```

and some more opinionated language features:

``` bash
# Case-insensitive string comparison
"LeBron" ~ "lebron" = true

# Array concatenation
[1] + [2] = [1, 2]
```

## Next

That's about it for the intro, we don't want to spoil the rest.
You can now head over to read [how to run ABS code!](/introduction/how-to-run-abs-code)