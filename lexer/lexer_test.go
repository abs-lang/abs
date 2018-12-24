package lexer

import (
	"testing"

	"github.com/abs-lang/abs/token"
)

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
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
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
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}

	// Should skip whitespaces etc
	input = `a   b c d`
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
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
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
