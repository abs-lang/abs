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
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) CurrentPosition() int {
	return l.position
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		if l.peekChar() == '=' {
			tok.Type = token.COMP_PLUS
			tok.Literal = "+="
			l.readChar()
		} else {
			tok = newToken(token.PLUS, l.ch)
		}
	case '-':
		if l.peekChar() == '=' {
			tok.Type = token.COMP_MINUS
			tok.Literal = "-="
			l.readChar()
		} else {
			tok = newToken(token.MINUS, l.ch)
		}
	case '%':
		if l.peekChar() == '=' {
			tok.Type = token.COMP_MODULO
			tok.Literal = "%="
			l.readChar()
		} else {
			tok = newToken(token.MODULO, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		if l.peekChar() == '/' {
			tok.Type = token.COMMENT
			tok.Literal = l.readComment()
		} else if l.peekChar() == '=' {
			tok.Type = token.COMP_SLASH
			tok.Literal = "/="
			l.readChar()
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
	case '#':
		tok.Type = token.COMMENT
		tok.Literal = l.readComment()
	case '&':
		if l.peekChar() == '&' {
			tok.Type = token.AND
			tok.Literal = l.readLogicalOperator()
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '*':
		if l.peekChar() == '*' {
			l.readChar()
			if l.peekChar() == '=' {
				tok.Type = token.COMP_EXPONENT
				tok.Literal = "**="
				l.readChar()
			} else {
				tok.Type = token.EXPONENT
				tok.Literal = "**"
			}
		} else if l.peekChar() == '=' {
			tok.Type = token.COMP_ASTERISK
			tok.Literal = "*="
			l.readChar()
		} else {
			tok = newToken(token.ASTERISK, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()

			if l.peekChar() == '>' {
				tok.Type = token.COMBINED_COMP
				tok.Literal = "<=>"
				l.readChar()
			} else {
				tok.Type = token.LT_EQ
				tok.Literal = "<="
			}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			tok.Type = token.GT_EQ
			tok.Literal = ">="
			l.readChar()
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '.':
		if l.peekChar() == '.' {
			tok.Type = token.RANGE
			tok.Literal = ".."
			l.readChar()
		} else {
			tok = newToken(token.DOT, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			tok.Type = token.OR
			tok.Literal = l.readLogicalOperator()
		} else {
			tok = newToken(token.PIPE, l.ch)
		}
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '~':
		tok = newToken(token.TILDE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '$':
		tok.Type = token.COMMAND
		tok.Literal = l.readCommand()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			literal, kind := l.readNumber()
			tok.Type = kind
			tok.Literal = literal
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
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
	l.readPosition += 1
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
func (l *Lexer) readString() string {
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
		if l.ch == '\\' && l.peekChar() == '"' {
			chars = append(chars, string('"'))
			l.readChar()
			continue
		}

		// The string ends when we encounter a "
		// and the character before that was not a \,
		// or the \ was escaped as well ("string\\").
		if (l.ch == '"' && (l.prevChar(2) != '\\' || doubleEscape)) || l.ch == 0 {
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
func (l *Lexer) readComment() string {
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

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
