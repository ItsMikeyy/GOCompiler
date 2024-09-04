package token

//All TokenTypes
const (
	ILLEGAL   = "ILLEGAL"
	EOF       = ""
	IDENT     = "IDENT"
	INT       = "INT"
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	BANG      = "!"
	ASTERISK  = "*"
	SLASH     = "/"
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	GT        = ">"
	LT        = "<"
	EQ        = "=="
	NOT_EQ    = "!="
	FUNCTION  = "FUNCTION"
	LET       = "LET"
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	IF        = "IF"
	ELSE      = "ELSE"
	RETURN    = "RETURN"
	STRING    = "STRING"
)

//Keywords mapped to TokenTypes
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

//Creating new type TokenType set to a string
type TokenType string

//Creating new type Token which is a struct that has a type: TokenType(string)
//and literal: string
type Token struct {
	Type    TokenType
	Literal string
}

//Used for debugging
//returns TokenType after passing in an identifier
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
