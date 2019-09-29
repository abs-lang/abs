# Installing 3rd party libraries

The ABS interpreter comes with a built-in installer for 3rd party libraries,
very similar to `npm install`, `pip install` or `go get`.

The installer, budled since the `1.8.0` release, is currently **experimental**
and a few things might change.

In order to install a package, you simply need to run `abs get`:

``` bash
$ abs get github.com/abs-lang/abs-sample-module 
🌘  - Downloading archive
Unpacking...
Creating alias...
Install Success. You can use the module with `require("abs-sample-module")`
```

Modules will be saved under the `vendor/$MODULE` directory. Each module
also gets an alias to facilitate requiring them in your code, meaning that
both of these forms are supported:

```
⧐  require("abs-sample-module/sample.abs")
{"another": f() {return hello world;}}

⧐  require("vendor/github.com/abs-lang/abs-sample-module/sample.abs")
{"another": f() {return hello world;}}
```

Module aliases are saved in the `packages.abs.json` file
which is created in the same directory where you run the
`abs get ...` command:

```
$ abs get github.com/abs-lang/abs-sample-module
🌗  - Downloading archive
Unpacking...
Creating alias...
Install Success. You can use the module with `require("abs-sample-module")`

$ cat packages.abs.json 
{
    "abs-sample-module": "./vendor/github.com/abs-lang/abs-sample-module"
}
```

If an alias is already taken, the installer will let you know that you
will need to use the full path when requiring the module:

```
$ echo '{"abs-sample-module": "xyz"}' > packages.abs.json 

$ abs get github.com/abs-lang/abs-sample-module
🌘  - Downloading archive
Unpacking...
Creating alias...This module could not be aliased because module of same name exists

Install Success. You can use the module with `require("./vendor/github.com/abs-lang/abs-sample-module")`
```

When requiring a module, ABS will try to load the `index.abs` file unless
another file is specified:

```
$ ~/projects/abs/builds/abs                                          
Hello alex, welcome to the ABS (1.9.0) programming language!
Type 'quit' when you're done, 'help' if you get lost!

⧐  require("abs-sample-module")
{"another": f() {return hello world;}}

⧐  require("abs-sample-module/index.abs")
{"another": f() {return hello world;}}

⧐  require("abs-sample-module/another.abs")
f() {return hello world;}
```

## Supported hosting platforms

Currently, the installer supports modules hosted on:

* GitHub

## Next

That's about it for this section!

You can now head over to read a little bit about [errors](/misc/error).