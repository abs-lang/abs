# Where should we write the GIF?
Output "/abs/docs/vhs/images/readme-sample.gif"

Set FontSize 24
Set Width 1200
Set Height 400
Set BorderRadius 10
Set Margin 10
Set MarginFill "#a1a1a1"
Set TypingSpeed 75ms

Type "cat /etc/timezone"
Enter
Sleep 0.75

Type "vi /tmp/timezone.abs"
Sleep 0.5
Enter
Sleep 0.5

Type "i"
Sleep 0.25
Type "tz = `cat /etc/timezone`"
Sleep 0.5
Enter
Type "continent, city = tz.split('/')"
Sleep 0.5
Enter
Enter
Type "echo('Best city in the world?')"
Sleep 0.5
Enter
Enter
Type "selection = stdin()"
Sleep 0.5
Enter
Enter
Type "if selection ~ city {"
Sleep 0.25
Enter
Type "  echo('You might be biased...')"
Sleep 0.25
Enter
Type "}"
Sleep 0.5

Escape

Type ":wq"
Sleep 0.25
Enter
Sleep 1
Type "abs /tmp/timezone.abs"
Sleep 1
Enter
Sleep 0.25
Type "dubai"
Sleep 0.5
Enter
Sleep 2