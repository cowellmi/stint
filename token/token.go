package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Literals
	IDENT = "IDENT"
	INT   = "INT"
	RAW   = "RAW"

	// Delimiters
	TAG    = "%"
	COLON  = ":"
	LPAREN = "("
	RPAREN = ")"
	COMMA  = ","
)
