<p align="center">
  <a href="https://www.abs-lang.org/">
    <img alt="abs language logo" src="https://github.com/abs-lang/abs/blob/master/bin/abs-horizontal.png?raw=true">
  </a>
</p>

# Errors

When using ABS, you might bump into errors within your code:
when the interpreter finds an error, it will give up trying
to evaluate the script and will exit with status code `99`:

```
$ cat examples/error.abs
1 + "hello"
echo("should not reach here")

$ abs examples/error.abs
ERROR: type mismatch: NUMBER + STRING

$ echo $?
99
```

## Next

That's about it for this section!

You can now head over to read [a few more technical details about ABS ](/misc/technical-details).