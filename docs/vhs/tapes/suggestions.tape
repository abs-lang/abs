Output "/abs/docs/vhs/images/suggestions.gif"

Set FontSize 24
Set Width 1200
Set Height 400
Set BorderRadius 10
Set Margin 10
Set MarginFill "#a1a1a1"
Set TypingSpeed 75ms

Source "/abs/docs/vhs/tapes/config.tape"

Type "ip = `curl icanhazip.com`"
Enter
Wait /.*/
Sleep 1

Type "ip.s"
Sleep 0.5
Tab@200ms 4
Enter
Sleep 0.5
Type "('.').map(int).s"
Tab@200ms 7
Enter
Sleep 0.25
Type "()"

Sleep 0.5
Enter
Wait /.*/
Sleep 1s