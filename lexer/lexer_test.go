package lexer

import (
	"testing"

	"github.com/cowellmi/stint/token"
)

func TestLexer(t *testing.T) {
	input := `Hello %name%!
Welcome to area %abc:int:foo(1,   2)%.`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.RAW, "Hello "},
		{token.TAG, "%"},
		{token.IDENT, "name"},
		{token.TAG, "%"},
		{token.RAW, "!\nWelcome to area "},
		{token.TAG, "%"},
		{token.IDENT, "abc"},
		{token.COLON, ":"},
		{token.IDENT, "int"},
		{token.COLON, ":"},
		{token.IDENT, "foo"},
		{token.LPAREN, "("},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.TAG, "%"},
		{token.RAW, "."},
		{token.EOF, ""},
	}

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
}
