package lexer

import (
	"strings"

	"github.com/abs-lang/abs/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	// map of input line boundaries used by linePosition() for error location
	lineMap [][2]int // array of [begin, end] pairs: [[0,12], [13,22], [23,33] ... ]
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
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
	return lineNum, column, errorLine
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

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) prevChar(steps int) byte {
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
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() (number string, kind token.TokenType) {
	position := l.position
	kind = token.NUMBER
	hasDot := false

	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' && (l.peekChar() == '.' || !isDigit(l.peekChar())) {
			return l.input[position:l.position], token.NUMBER
		}

		if l.ch == '.' {
			if hasDot {
				return "", token.ILLEGAL
			}

			hasDot = true
			kind = token.NUMBER
		}
		l.readChar()
	}

	return l.input[position:l.position], kind
}

// A logical operator is 2 chars, so
// we can simply read 2 chars and call
// it a day.
func (l *Lexer) readLogicalOperator() string {
	l.readChar()
	return l.input[l.position-1 : l.position+1]
}

// Reads a strings from the text.
// Strings can be escaped with \".
// To close a string with an escape
// character ("\"), escape the escape
// character itself ("\\").
func (l *Lexer) readString(quote byte) string {
	var chars []string
	doubleEscape := false
	for {
		l.readChar()

		if l.ch == '\\' && l.peekChar() == '\\' {
			chars = append(chars, string('\\'))
			l.readChar()
			doubleEscape = true
			continue
		}

		// If we encounter a \, let's check whether
		// we're trying to escape a ". If so, let's skip
		// the / and add the " to the string.
		if l.ch == '\\' && l.peekChar() == quote {
			chars = append(chars, string(quote))
			l.readChar()
			continue
		}

		// The string ends when we encounter a "
		// and the character before that was not a \,
		// or the \ was escaped as well ("string\\").
		if (l.ch == quote && (l.prevChar(2) != '\\' || doubleEscape)) || l.ch == 0 {
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
	return l.input[position:l.position]
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

	return ret
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
