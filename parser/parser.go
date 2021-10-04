package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	lexer        *lexer.Lexer // An instance of the lexer
	currentToken token.Token  // The current token
	peekToken    token.Token  // The next token
}

func New(lexer *lexer.Lexer) *Parser {
	parser := Parser{lexer: lexer}
	// advance two times so as to set the value of the current token and the next token
	parser.advance()
	parser.advance()
	return &parser
}

func (ps *Parser) advance() {
	ps.currentToken = ps.peekToken
	ps.peekToken = ps.lexer.GetNextToken()
}

func (ps *Parser) ParseProgram() *ast.Program {
	program := ast.Program{}
	program.Statements = []ast.Statement{}

	for ps.currentToken.Type != token.EOF {
		stmt := ps.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		ps.advance()
	}
	return &program
}

func (ps *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return ps.peekToken.Type == tokenType
}

func (ps *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return ps.currentToken.Type == tokenType
}

func (ps *Parser) parseStatement() ast.Statement {
	switch ps.currentToken.Type {
	case token.LET:
		return ps.parseLetStatement()
	case token.RETURN:
		return ps.parseReturnStatement()
	default:
		return nil
	}
}

/** Parse LET Statement **/
func (ps *Parser) parseLetStatement() *ast.LET_Statement {
	stmt := ast.LET_Statement{Token: ps.currentToken}
	if !ps.peekTokenIs(token.IDENTIFIER) {
		return nil
	}
	ps.advance()
	stmt.Name = &ast.Identifier{Token: ps.currentToken, Value: ps.currentToken.Literal}
	if !ps.peekTokenIs(token.ASSIGN) {
		return nil
	}
	ps.advance()
	// TODO: For now skipping all expressions until we get to the semicolon
	for !ps.currentTokenIs(token.SEMICOLON) {
		ps.advance()
	}
	return &stmt
}

/** Parse RETURN Statement **/
func (ps *Parser) parseReturnStatement() *ast.RETURN_Statement {
	stmt := ast.RETURN_Statement{Token: ps.currentToken}
	ps.advance()
	// TODO: Implement parsing of expressions for now skip all expressons until ';'
	for !ps.currentTokenIs(token.SEMICOLON) {
		ps.advance()
	}
	return &stmt
}
