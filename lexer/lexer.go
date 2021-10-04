package lexer

import "monkey/token"

// The lexer struct
type Lexer struct {
	input string // The input that is being scanned
	index int    // The index of the next character to be read
	char  byte   // The current character that is being read == input[index - 1]
}

// Helper function to create a new Lexer
func New(input string) *Lexer {
	lx := &Lexer{input: input}
	lx.readNextChar()
	return lx
}

// Helper function that returns the next char only
func (lx *Lexer) peekChar() byte {
	if lx.index == len(lx.input) {
		return 0
	} else {
		return lx.input[lx.index]
	}
}

// Helper function to read the next character and advance the pointers
func (lx *Lexer) readNextChar() {
	// If the index to be read is beyond the input then assign lx.char = '\0'
	if lx.index >= len(lx.input) {
		lx.char = 0
	} else {
		lx.char = lx.input[lx.index]
	}
	lx.index++
}

var KEYWORDS = map[string]token.TokenType{
	"fn":     token.FUNCTION,
	"false":  token.FALSE,
	"true":   token.TRUE,
	"else":   token.ELSE,
	"if":     token.IF,
	"return": token.RETURN,
	"let":    token.LET,
}

// GetNextToken returns the next token
func (lx *Lexer) GetNextToken() token.Token {
	lx.skipWhiteSpaces()
	var tok token.Token
	switch lx.char {
	case 0:
		tok = token.Token{Literal: "", Type: token.EOF}
	case '=':
		if lx.peekChar() == '=' {
			lx.readNextChar()
			tok = token.Token{Literal: "==", Type: token.EQ}
		} else {
			tok = lx.makeToken(token.ASSIGN)
		}
	case ';':
		tok = lx.makeToken(token.SEMICOLON)
	case '(':
		tok = lx.makeToken(token.LPAREN)
	case ')':
		tok = lx.makeToken(token.RPAREN)
	case ',':
		tok = lx.makeToken(token.COMMA)
	case '+':
		tok = lx.makeToken(token.PLUS)
	case '{':
		tok = lx.makeToken(token.LBRACE)
	case '}':
		tok = lx.makeToken(token.RBRACE)
	case '-':
		tok = lx.makeToken(token.MINUS)
	case '!':
		if lx.peekChar() == '=' {
			lx.readNextChar()
			tok = token.Token{Literal: "!=", Type: token.NOT_EQ}
		} else {
			tok = lx.makeToken(token.BANG)
		}
	case '*':
		tok = lx.makeToken(token.ASTERISK)
	case '/':
		tok = lx.makeToken(token.SLASH)
	case '>':
		tok = lx.makeToken(token.GT)
	case '<':
		tok = lx.makeToken(token.LT)
	default:
		if lx.isLetter() {
			lit := lx.readIdentifier()
			_type, ok := KEYWORDS[lit]
			if ok {
				return token.Token{Literal: lit, Type: _type}
			}
			return token.Token{Literal: lit, Type: token.IDENTIFIER}
		} else if lx.isDigit() {
			return token.Token{Literal: lx.readNumber(), Type: token.INT}
		}
		tok = token.Token{Literal: "", Type: token.ILLEGAL}
	}
	lx.readNextChar()
	return tok
}

// Helper function to make a token
func (lx *Lexer) makeToken(tType token.TokenType) token.Token {
	return token.Token{Literal: string(lx.char), Type: tType}
}

// Helper function to check if lx.char is a string
func (lx *Lexer) isLetter() bool {
	return lx.char >= 'a' && lx.char <= 'z' || lx.char >= 'A' && lx.char <= 'Z'
}

// Helper function to read Identifiers
func (lx *Lexer) readIdentifier() string {
	var mark int = lx.index - 1
	for lx.isLetter() {
		lx.readNextChar()
	}
	return lx.input[mark : lx.index-1]
}

// Helper function to check if a charater is a digit
func (lx *Lexer) isDigit() bool {
	return lx.char >= '0' && lx.char <= '9'
}

// Helper function to read Numbers
func (lx *Lexer) readNumber() string {
	var mark int = lx.index - 1
	for lx.isDigit() {
		lx.readNextChar()
	}
	return lx.input[mark : lx.index-1]
}

// Function to skip white spaces
func (lx *Lexer) skipWhiteSpaces() {
	for lx.char == '\t' || lx.char == '\r' || lx.char == '\n' || lx.char == ' ' {
		lx.readNextChar()
	}
}
