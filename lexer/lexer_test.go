package lexer

import (
	"testing"

	"github.com/abs-lang/abs/token"
)

func TestNextToken(t *testing.T) {
	input := `five = 5;
ten = 10;
add = f(x, y) {
  x + y;
};
result = add(five, ten);
&&||!-/*5;
5 < 10 > 5;
1 <= 1 >= 1;
<=>
if (5 < 10) {
	return true;
} else {
	return false;
}
while (1 > 0) {
	echo("hello")
}
for x in xs {
	x
}
for x = 0; x < 10; x = x + 1 {
	x
}
10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2];
$(echo "()");
{"foo": "bar"}
$(curl icanhazip.com -X POST)
$(ls *.go);
a = [1]
a.first()
a.prop
# Comment
// Comment
hello
$(command; command)
$(command2; command2);
one | two | tree
"hel\"lo"
"hel\lo"
"hel\\\\lo"
"\"hello\""
"\"he\"\"llo\""
"hello\\"
"hello\\\\"
"\\\\hello"
**
1..10
~%
+=
-=
*=
/=
**=
%=
1.23
1.str()
null
nullo
&^>><<
$111
'123'
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "f"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.NUMBER, "5"},
		{token.LT, "<"},
		{token.NUMBER, "10"},
		{token.GT, ">"},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.NUMBER, "1"},
		{token.LT_EQ, "<="},
		{token.NUMBER, "1"},
		{token.GT_EQ, ">="},
		{token.NUMBER, "1"},
		{token.SEMICOLON, ";"},
		{token.COMBINED_COMP, "<=>"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.NUMBER, "5"},
		{token.LT, "<"},
		{token.NUMBER, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.WHILE, "while"},
		{token.LPAREN, "("},
		{token.NUMBER, "1"},
		{token.GT, ">"},
		{token.NUMBER, "0"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "echo"},
		{token.LPAREN, "("},
		{token.STRING, "hello"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.FOR, "for"},
		{token.IDENT, "x"},
		{token.IN, "in"},
		{token.IDENT, "xs"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.RBRACE, "}"},
		{token.FOR, "for"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.NUMBER, "0"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "x"},
		{token.LT, "<"},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.NUMBER, "1"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.RBRACE, "}"},
		{token.NUMBER, "10"},
		{token.EQ, "=="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},
		{token.NUMBER, "10"},
		{token.NOT_EQ, "!="},
		{token.NUMBER, "9"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.LBRACKET, "["},
		{token.NUMBER, "1"},
		{token.COMMA, ","},
		{token.NUMBER, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.COMMAND, `echo "()"`},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.COMMAND, "curl icanhazip.com -X POST"},
		{token.COMMAND, "ls *.go"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.LBRACKET, "["},
		{token.NUMBER, "1"},
		{token.RBRACKET, "]"},
		{token.IDENT, "a"},
		{token.DOT, "."},
		{token.IDENT, "first"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.IDENT, "a"},
		{token.DOT, "."},
		{token.IDENT, "prop"},
		{token.COMMENT, "# Comment"},
		{token.COMMENT, "// Comment"},
		{token.IDENT, "hello"},
		{token.COMMAND, "command; command"},
		{token.COMMAND, "command2; command2"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "one"},
		{token.PIPE, "|"},
		{token.IDENT, "two"},
		{token.PIPE, "|"},
		{token.IDENT, "tree"},
		{token.STRING, "hel\"lo"},
		{token.STRING, "hel\\lo"},
		{token.STRING, "hel\\\\lo"},
		{token.STRING, "\"hello\""},
		{token.STRING, "\"he\"\"llo\""},
		{token.STRING, "hello\\"},
		{token.STRING, "hello\\\\"},
		{token.STRING, "\\\\hello"},
		{token.EXPONENT, "**"},
		{token.NUMBER, "1"},
		{token.RANGE, ".."},
		{token.NUMBER, "10"},
		{token.TILDE, "~"},
		{token.MODULO, "%"},
		{token.COMP_PLUS, "+="},
		{token.COMP_MINUS, "-="},
		{token.COMP_ASTERISK, "*="},
		{token.COMP_SLASH, "/="},
		{token.COMP_EXPONENT, "**="},
		{token.COMP_MODULO, "%="},
		{token.NUMBER, "1.23"},
		{token.NUMBER, "1"},
		{token.DOT, "."},
		{token.IDENT, "str"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.NULL, "null"},
		{token.IDENT, "nullo"},
		{token.BIT_AND, "&"},
		{token.BIT_XOR, "^"},
		{token.BIT_RSHIFT, ">>"},
		{token.BIT_LSHIFT, "<<"},
		{token.ILLEGAL, "$111"},
		{token.STRING, "123"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (%s, %s)", i, tt.expectedType, tok.Type, tok.Literal, tt.expectedLiteral)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestRewind(t *testing.T) {
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "a"},
		{token.IDENT, "b"},
		{token.IDENT, "c"},
		{token.IDENT, "d"},
		{token.EOF, ""},
	}

	input := `a b c d`
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (%s, %s)", i, tt.expectedType, tok.Type, tok.Literal, tt.expectedLiteral)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}

	l.Rewind(0)

	tests = []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "a"},
		{token.IDENT, "b"},
		{token.IDENT, "c"},
		{token.IDENT, "d"},
		{token.EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}

	// Should skip whitespaces etc
	input = `a b c d`
	l = New(input)
	l.Rewind(3)

	tests = []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "c"},
		{token.IDENT, "d"},
		{token.EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestCurrentPosition(t *testing.T) {
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "a"},
		{token.IDENT, "b"},
		{token.IDENT, "c"},
		{token.IDENT, "d"},
		{token.EOF, ""},
	}

	input := `a b c d`
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}

	if l.CurrentPosition() != 8 {
		t.Fatalf("wrong position. expected=%d, got=%d", 6, l.CurrentPosition())
	}

	l.Rewind(0)

	if l.CurrentPosition() != 0 {
		t.Fatalf("wrong position. expected=%d, got=%d", 6, l.CurrentPosition())
	}
}
