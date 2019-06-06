package lexer

import (
	"strings"
	"unicode"

	"github.com/abs-lang/abs/token"
)

type Lexer struct {
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current rune under examination
	input        []rune
	// map of input line boundaries used by linePosition() for error location
	lineMap [][2]int // array of [begin, end] pairs: [[0,12], [13,22], [23,33] ... ]
}

func New(in string) *Lexer {
	l := &Lexer{input: []rune(in)}
	// map the input line boundaries for CurrentLine()
	l.buildLineMap()
	// read the first char
	l.readChar()
	return l
}

// buildLineMap creates map of input line boundaries used by LinePosition() for error location
func (l *Lexer) buildLineMap() {
	begin := 0
	idx := 0
	for i, ch := range l.input {
		idx = i
		if ch == '\n' {
			l.lineMap = append(l.lineMap, [2]int{begin, idx})
			begin = idx + 1
		}
	}
	// last line
	l.lineMap = append(l.lineMap, [2]int{begin, idx + 1})
}

// CurrentPosition returns l.position
func (l *Lexer) CurrentPosition() int {
	return l.position
}

// linePosition (pos) returns lineNum, begin, end
func (l *Lexer) linePosition(pos int) (int, int, int) {
	idx := 0
	begin := 0
	end := 0
	for i, tuple := range l.lineMap {
		idx = i
		begin, end = tuple[0], tuple[1]
		if pos >= begin && pos <= end {
			break
		}
	}
	lineNum := idx + 1
	return lineNum, begin, end
}

// ErrorLine (pos) returns lineNum, column, errorLine
func (l *Lexer) ErrorLine(pos int) (int, int, string) {
	lineNum, begin, end := l.linePosition(pos)
	errorLine := l.input[begin:end]
	column := pos - begin + 1
	return lineNum, column, string(errorLine)
}

func (l *Lexer) newToken(tokenType token.TokenType) token.Token {
	return token.Token{
		Type:     tokenType,
		Position: l.position,
		Literal:  string(l.ch)}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			tok = l.newToken(token.EQ)
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok.Literal = literal
		} else {
			tok = l.newToken(token.ASSIGN)
		}
	case '+':
		if l.peekChar() == '=' {
			tok.Type = token.COMP_PLUS
			tok.Position = l.position
			tok.Literal = "+="
			l.readChar()
		} else {
			tok = l.newToken(token.PLUS)
		}
	case '-':
		if l.peekChar() == '=' {
			tok.Type = token.COMP_MINUS
			tok.Position = l.position
			tok.Literal = "-="
			l.readChar()
		} else {
			tok = l.newToken(token.MINUS)
		}
	case '%':
		if l.peekChar() == '=' {
			tok.Type = token.COMP_MODULO
			tok.Position = l.position
			tok.Literal = "%="
			l.readChar()
		} else {
			tok = l.newToken(token.MODULO)
		}
	case '!':
		if l.peekChar() == '=' {
			tok = l.newToken(token.NOT_EQ)
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok.Literal = literal
		} else {
			tok = l.newToken(token.BANG)
		}
	case '/':
		if l.peekChar() == '/' {
			tok.Type = token.COMMENT
			tok.Position = l.position
			tok.Literal = l.readLine()
		} else if l.peekChar() == '=' {
			tok.Type = token.COMP_SLASH
			tok.Position = l.position
			tok.Literal = "/="
			l.readChar()
		} else {
			tok = l.newToken(token.SLASH)
		}
	case '#':
		tok.Type = token.COMMENT
		tok.Position = l.position
		tok.Literal = l.readLine()
	case '&':
		if l.peekChar() == '&' {
			tok.Type = token.AND
			tok.Position = l.position
			tok.Literal = l.readLogicalOperator()
		} else {
			tok = l.newToken(token.BIT_AND)
		}
	case '^':
		tok = l.newToken(token.BIT_XOR)
	case '*':
		if l.peekChar() == '*' {
			l.readChar()
			if l.peekChar() == '=' {
				tok.Type = token.COMP_EXPONENT
				tok.Position = l.position
				tok.Literal = "**="
				l.readChar()
			} else {
				tok.Type = token.EXPONENT
				tok.Position = l.position
				tok.Literal = "**"
			}
		} else if l.peekChar() == '=' {
			tok.Type = token.COMP_ASTERISK
			tok.Position = l.position
			tok.Literal = "*="
			l.readChar()
		} else {
			tok = l.newToken(token.ASTERISK)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()

			if l.peekChar() == '>' {
				tok.Type = token.COMBINED_COMP
				tok.Position = l.position
				tok.Literal = "<=>"
				l.readChar()
			} else {
				tok.Type = token.LT_EQ
				tok.Position = l.position
				tok.Literal = "<="
			}
		} else if l.peekChar() == '<' {
			tok.Type = token.BIT_LSHIFT
			tok.Position = l.position
			tok.Literal = "<<"
			l.readChar()
		} else {
			tok = l.newToken(token.LT)
		}
	case '>':
		if l.peekChar() == '=' {
			tok.Type = token.GT_EQ
			tok.Position = l.position
			tok.Literal = ">="
			l.readChar()
		} else if l.peekChar() == '>' {
			tok.Type = token.BIT_RSHIFT
			tok.Position = l.position
			tok.Literal = ">>"
			l.readChar()
		} else {
			tok = l.newToken(token.GT)
		}
	case ';':
		tok = l.newToken(token.SEMICOLON)
	case ':':
		tok = l.newToken(token.COLON)
	case ',':
		tok = l.newToken(token.COMMA)
	case '.':
		if l.peekChar() == '.' {
			tok.Type = token.RANGE
			tok.Position = l.position
			tok.Literal = ".."
			l.readChar()
		} else {
			tok = l.newToken(token.DOT)
		}
	case '|':
		if l.peekChar() == '|' {
			tok.Type = token.OR
			tok.Position = l.position
			tok.Literal = l.readLogicalOperator()
		} else {
			tok = l.newToken(token.PIPE)
		}
	case '{':
		tok = l.newToken(token.LBRACE)
	case '}':
		tok = l.newToken(token.RBRACE)
	case '~':
		tok = l.newToken(token.TILDE)
	case '(':
		tok = l.newToken(token.LPAREN)
	case ')':
		tok = l.newToken(token.RPAREN)
	case '"':
		tok.Type = token.STRING
		tok.Position = l.position
		tok.Literal = l.readString('"')
	case '\'':
		tok.Type = token.STRING
		tok.Position = l.position
		tok.Literal = l.readString('\'')
	case '`':
		tok.Type = token.COMMAND
		tok.Position = l.position
		tok.Literal = l.readString('`')
	case '$':
		if l.peekChar() == '(' {
			tok.Type = token.COMMAND
			tok.Position = l.position
			tok.Literal = l.readCommand()
		} else {
			tok.Type = token.ILLEGAL
			tok.Position = l.position
			tok.Literal = l.readLine()
		}
	case '[':
		tok = l.newToken(token.LBRACKET)
	case ']':
		tok = l.newToken(token.RBRACKET)
	case 0:
		tok.Type = token.EOF
		tok.Position = l.position
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Position = l.position
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Position = l.position
			literal, kind := l.readNumber()
			tok.Type = kind
			tok.Literal = literal
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// This function will read a rune
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) Rewind(pos int) {
	l.ch = l.input[0]
	l.position = 0
	l.readPosition = l.position + 1

	for l.position < pos {
		l.NextToken()
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) prevChar(steps int) rune {
	prevPosition := l.readPosition - steps
	if prevPosition < 1 {
		return 0
	}
	return l.input[prevPosition]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

// 12
// 12.2
// 12e-1
// 12e+1
// 12e1
func (l *Lexer) readNumber() (number string, kind token.TokenType) {
	position := l.position
	hasDot := false
	kind = token.NUMBER
	hasExponent := false

	// List of character that can appear in a "number"
	for isDigit(l.ch) || l.ch == '.' || l.ch == '+' || l.ch == '-' || l.ch == 'e' || l.ch == '_' {
		// If we have a plus / minus but there was no exponent
		// in this number, it means we're at the end of the
		// number and we're at an addition / subtraction.
		if (l.ch == '+' || l.ch == '-') && !hasExponent {
			return string(l.input[position:l.position]), kind
		}

		// If the number contains as 'e',
		// we're using scientific notation
		if l.ch == 'e' {
			hasExponent = true
		}

		// If we have a dot, let's check whether this is a range
		// or maybe a method call (122.string())
		if l.ch == '.' && (l.peekChar() == '.' || !isDigit(l.peekChar())) {
			return string(l.input[position:l.position]), kind
		}

		if l.ch == '.' {
			// If we have 2 dots in a number, there's a problem
			if hasDot {
				return string(l.input[position : l.position+1]), token.ILLEGAL
			}

			hasDot = true
		}
		l.readChar()
	}

	// If the number ends with the exponent,
	// there's a problem.
	if l.input[l.position-1] == 'e' {
		return string(l.input[position:l.position]), token.ILLEGAL
	}

	return strings.ReplaceAll(string(l.input[position:l.position]), "_", ""), kind
}

// A logical operator is 2 chars, so
// we can simply read 2 chars and call
// it a day.
func (l *Lexer) readLogicalOperator() string {
	l.readChar()
	return string(l.input[l.position-1 : l.position+1])
}

// Reads a strings from the text.
// Strings can be escaped with \".
// To close a string with an escape
// character ("\"), escape the escape
// character itself ("\\").
func (l *Lexer) readString(quote byte) string {
	var chars []string
	esc := rune('\\')
	doubleEscape := false
	for {
		l.readChar()

		if l.ch == esc && l.peekChar() == esc {
			chars = append(chars, string(esc))
			l.readChar()
			// be careful here, there may be double escaped LFs in the string
			if l.peekChar() == rune(quote) {
				doubleEscape = true
			} else {
				// this is not a double escaped quote
				chars = append(chars, string(esc))
			}
			continue
		}
		// If we encounter an escape, let's check whether
		// we're trying to escape a quote. If so, let's skip
		// the escape and add the quote to the string.
		if l.ch == esc && l.peekChar() == rune(quote) {
			chars = append(chars, string(quote))
			l.readChar()
			continue
		}
		// If this is a double quoted string we need to expand embedded
		// LF, CR, and TAB to ASCII and add the ASCII code to the string
		// NB. single quoted strings don't expand special characters to ASCII
		if quote == '"' {
			if l.ch == esc && l.peekChar() == 'n' {
				chars = append(chars, "\n")
				l.readChar()
				continue
			} else if l.ch == esc && l.peekChar() == 'r' {
				chars = append(chars, "\r")
				l.readChar()
				continue
			} else if l.ch == esc && l.peekChar() == 't' {
				chars = append(chars, "\t")
				l.readChar()
				continue
			}
		}
		// The string ends when we encounter a quote
		// and the character before that was not an escape,
		// or the escape was escaped as well ("string\\").
		if (l.ch == rune(quote) && (l.prevChar(2) != esc || doubleEscape)) || l.ch == 0 {
			break
		}
		chars = append(chars, string(l.ch))
		doubleEscape = false
	}
	return strings.Join(chars, "")
}

// Go ahead until you find a new line.
// This makes it so that comments take
// a full line.
func (l *Lexer) readLine() string {
	position := l.position
	for {
		l.readChar()
		if l.ch == '\n' || l.ch == '\r' || l.ch == 0 {
			break
		}
	}
	return string(l.input[position:l.position])
}

// We want to extract the actual command
// from $(command).
// We first go ahead 2 characters (`$(`)
// and then remove either 1 or 2 of them
// at the end (`)` or `);`). This forces
// commands to be on a single line,
// which is also not a terrible thing
// (we might want to support \ at some point
// in the future).
func (l *Lexer) readCommand() string {
	position := l.position + 2
	subtract := 1
	for {
		l.readChar()

		if l.ch == '\n' || l.ch == '\r' || l.ch == 0 {
			// TODO: for compat turn this into prevCharOtherThan
			if l.prevChar(2) == ';' {
				subtract = 2
			}
			break
		}
	}
	ret := l.input[position : l.position-subtract]

	// Let's make sure the semicolo is the next token, without
	// "cutting" it out...
	if subtract == 2 {
		l.position = l.position - 1
		l.readPosition = l.position
	}

	return string(ret)
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
