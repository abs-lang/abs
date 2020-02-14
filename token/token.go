package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT        = "IDENT"  // add, foobar, x, y, ...
	NUMBER       = "NUMBER" // 1343456, 1.23456
	STRING       = "STRING" // "foobar"
	COMMENT      = "#"      // # Comment
	AT           = "@"      // @ At symbol
	NULL         = "NULL"   // # null
	CURRENT_ARGS = "..."    // # ... function args

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
	OR  = "||"

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
	QUESTION = "?"
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
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
)

type Token struct {
	Type     TokenType
	Position int // lexer position in file before token
	Literal  string
}

var keywords = map[string]TokenType{
	"f":        FUNCTION,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"while":    WHILE,
	"for":      FOR,
	"in":       IN,
	"null":     NULL,
	"break":    BREAK,
	"continue": CONTINUE,
}

// NumberAbbreviations is a list of abbreviations that can be used in numbers eg. 1k, 20B
var NumberAbbreviations = map[string]float64{
	"k": 1000,
	"m": 1000000,
	"b": 1000000000,
	"t": 1000000000000,
}

// NumberSeparator is a separator for numbers eg. 1_000_000
var NumberSeparator = '_'

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
