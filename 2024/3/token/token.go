package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	LPAREN = "("
	RPAREN = ")"

	IDENT = "IDENT"
	INT   = "INT"
)
