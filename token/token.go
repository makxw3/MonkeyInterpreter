package token

import (
	"fmt"
)

type TokenType string

type Token struct {
	Type    TokenType // The type of the token
	Literal string    // The literal value of the token
}

// The supported types of tokens
const (
	ILLEGAL    = "ILLEGAL"    // Any unsuported token
	EOF        = "EOF"        // The end of file token
	IDENTIFIER = "IDENTIFIER" // Represents any identifier
	INT        = "INT"        // Integers such as 1,2, ...
	ASSIGN     = "="
	PLUS       = "+"
	COMMA      = ","
	SEMICOLON  = ";"
	LPAREN     = "("
	RPAREN     = ")"
	LBRACE     = "{"
	RBRACE     = "}"
	MINUS      = "-"
	BANG       = "!"
	ASTERISK   = "*"
	SLASH      = "/"
	LT         = "<"
	GT         = ">"
	EQ         = "=="
	NOT_EQ     = "!="
	/** Keywords **/
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

func (t Token) Print() {
	fmt.Printf("<%v,%v>\n", t.Type, t.Literal)
}
