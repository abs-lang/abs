package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT   = "IDENT"  // add, foobar, x, y, ...
	NUMBER  = "NUMBER" // 1343456, 1.23456
	STRING  = "STRING" // "foobar"
	COMMENT = "#"      // # Comment
	NULL    = "NULL"   // # null

	// Operators
	TILDE         = "~"
	BANG          = "!"
	ASSIGN        = "="
	PLUS          = "+"
	MINUS         = "-"
	ASTERISK      = "*"
	SLASH         = "/"
	EXPONENT      = "**"
	MODULO        = "%"
	COMP_PLUS     = "+="
	COMP_MINUS    = "-="
	COMP_ASTERISK = "*="
	COMP_SLASH    = "/="
	COMP_EXPONENT = "**="
	COMP_MODULO   = "%="
	RANGE         = ".."

	// Logical operators
	AND = "&&"
	OR  = "OR"

	// Bitwise operators
	// It might be worth
	// to rename these
	// to AMPERSAND / CARET / etc
	BIT_AND    = "&"
	BIT_XOR    = "^"
	BIT_RSHIFT = ">>"
	BIT_LSHIFT = "<<"
	PIPE       = "|"

	LT            = "<"
	LT_EQ         = "<="
	GT            = ">"
	GT_EQ         = ">="
	COMBINED_COMP = "<=>"

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
	DOT      = "."
	COMMAND  = "$()"

	// Keywords
	FUNCTION = "F"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	WHILE    = "WHILE"
	FOR      = "FOR"
	IN       = "IN"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"f":      FUNCTION,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
	"for":    FOR,
	"in":     IN,
	"null":   NULL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
