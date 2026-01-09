package lexer

import "github.com/cowellmi/stint/token"

type Lexer struct {
	tagged       bool // state for inside %...%
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

func (l *Lexer) NextToken() token.Token {
	if l.ch == 0 {
		return token.Token{Type: token.EOF, Literal: ""}
	}

	// State control
	if l.ch == '%' {
		l.tagged = !l.tagged // toggle state
		tok := newToken(token.TAG, l.ch)
		l.readChar()
		return tok
	}

	// Tagged
	if l.tagged {
		l.skipWhitespace()
		return l.lexTagged()
	}

	// Raw
	return token.Token{
		Type:    token.RAW,
		Literal: l.readRaw(),
	}
}

func (l *Lexer) lexTagged() token.Token {
	var tok token.Token

	switch l.ch {
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	default:
		if isLetter(l.ch) {
			tok.Type = token.IDENT
			tok.Literal = l.readIdentifier()
			return tok
		}
		if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
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

func (l *Lexer) readRaw() string {
	start := l.position
	for l.ch != '%' && l.ch != 0 {
		l.readChar()
	}
	end := l.position
	return l.input[start:end]
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	end := l.position
	return l.input[start:end]
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	end := l.position
	return l.input[start:end]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
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
