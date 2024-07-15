package token

const (
	ILLEGAL   = "ILLEGAL"
	EOF       = ""
	IDENT     = "IDENT"
	INT       = "INT"
	ASSIGN    = "ASSIGN"
	PLUS      = "PLUS"
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	FUNCTION  = "FUNCTION"
	LET       = "LET"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
