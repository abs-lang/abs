# Where should we write the GIF?
Output "/abs/docs/vhs/images/ip.gif"

# Set up a 1200x600 terminal with 46px font.
Set FontSize 24
Set Width 1200
Set Height 400
Set BorderRadius 10
Set Margin 10
Set MarginFill "#a1a1a1"
Set TypingSpeed 75ms

Type "abs"
Sleep 0.5
Enter

Wait+Screen /welcome/
Sleep 0.5

Type "ip = `curl icanhazip.com`"
Enter
Wait /.*/
Sleep 1

Type "ip"
Sleep 0.5
Enter
Wait /.*/
Sleep 0.5

Type "ip.split('.').map(int).map(f(n) {n**2})"
Sleep 0.5
Enter
Wait /.*/
Sleep 1s