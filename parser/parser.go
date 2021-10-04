package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type Parser struct {
	lexer        *lexer.Lexer // An instance of the lexer
	currentToken token.Token  // The current token
	peekToken    token.Token  // The next token

	/** PRATT **/
	_prefixParsingFunctions map[token.TokenType]prattPrefixParsingFuncntion
	_infixParsingFunctins   map[token.TokenType]prattInfixParsingFunction
}

func New(lexer *lexer.Lexer) *Parser {
	parser := Parser{lexer: lexer}
	// advance two times so as to set the value of the current token and the next token
	parser.advance()
	parser.advance()
	parser._infixParsingFunctins = make(map[token.TokenType]prattInfixParsingFunction)
	parser._prefixParsingFunctions = make(map[token.TokenType]prattPrefixParsingFuncntion)

	parser.addPrefixFn(parser.parseIdentifier, token.IDENTIFIER)
	parser.addPrefixFn(parser.parseIntegerLiteral, token.INT)
	parser.addPrefixFn(parser.parsePrefixExpression, token.BANG)
	parser.addPrefixFn(parser.parsePrefixExpression, token.MINUS)
	parser.addPrefixFn(parser.parseBoolean, token.TRUE)
	parser.addPrefixFn(parser.parseBoolean, token.FALSE)
	parser.addPrefixFn(parser.parseGroupedExpression, token.LPAREN)
	parser.addPrefixFn(parser.parseIFExpression, token.IF)
	parser.addPrefixFn(parser.parseFunctionLiteral, token.FUNCTION)

	parser.addInfixFn(parser.parseInfixExpression, token.PLUS)
	parser.addInfixFn(parser.parseInfixExpression, token.MINUS)
	parser.addInfixFn(parser.parseInfixExpression, token.SLASH)
	parser.addInfixFn(parser.parseInfixExpression, token.ASTERISK)
	parser.addInfixFn(parser.parseInfixExpression, token.EQ)
	parser.addInfixFn(parser.parseInfixExpression, token.NOT_EQ)
	parser.addInfixFn(parser.parseInfixExpression, token.LT)
	parser.addInfixFn(parser.parseInfixExpression, token.GT)
	parser.addInfixFn(parser.parseCallExpression, token.LPAREN)
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
		program.Statements = append(program.Statements, stmt)
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
		return ps.parseExpressionStatement()
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
	ps.advance()
	stmt.Value = ps._parseExpression(LOWEST)
	if ps.peekTokenIs(token.SEMICOLON) {
		ps.advance()
	}
	return &stmt
}

/** Parse RETURN Statement **/
func (ps *Parser) parseReturnStatement() *ast.RETURN_Statement {
	stmt := ast.RETURN_Statement{Token: ps.currentToken}
	ps.advance()
	stmt.ReturnValue = ps._parseExpression(LOWEST)
	if ps.peekTokenIs(token.SEMICOLON) {
		ps.advance()
	}
	return &stmt
}

/** Parse EXPRESSION Statement **/
func (ps *Parser) parseExpressionStatement() *ast.EXPRESSION_Statement {
	stmt := ast.EXPRESSION_Statement{}
	/** Magical function _parseExpression(int) **/
	stmt.Expression = ps._parseExpression(LOWEST)
	if ps.peekTokenIs(token.SEMICOLON) {
		ps.advance()
	}
	return &stmt
}

/** Parse IDENTIFIER **/
func (ps *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: ps.currentToken, Value: ps.currentToken.Literal}
}

func (ps *Parser) parseIntegerLiteral() ast.Expression {
	lit := ast.INTEGER_Literal{Token: ps.currentToken}
	val, err := strconv.ParseInt(ps.currentToken.Literal, 0, 64)
	if err != nil {
		fmt.Printf("Error in converting %s int to int64", ps.currentToken.Literal)
		return nil
	}
	lit.Value = val
	return &lit
}

/** Parse PrefixExpression **/
func (ps *Parser) parsePrefixExpression() ast.Expression {
	expr := ast.PREFIX_Expression{Token: ps.currentToken}
	ps.advance()
	expr.Right = ps._parseExpression(PREFIX)
	return &expr
}

/** Parse Infix Expression **/
func (ps *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := ast.INFIX_Expression{
		Token: ps.currentToken,
		Left:  left,
	}
	precedence := ps.currentPrecedence()
	ps.advance()
	expr.Right = ps._parseExpression(precedence)
	return &expr
}

/** Parse Boolean **/
func (ps *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: ps.currentToken, Value: ps.currentTokenIs(token.TRUE)}
}

/** Parse Grouped Expression **/
func (ps *Parser) parseGroupedExpression() ast.Expression {
	ps.advance()
	// parses the expression until it encounters a token whose precedence is == LOWEST e.g. ')'
	expr := ps._parseExpression(LOWEST)
	if !ps.peekTokenIs(token.RPAREN) {
		return nil
	}
	ps.advance()
	return expr
}

/** Parse IF Expression **/
func (ps *Parser) parseIFExpression() ast.Expression {
	expr := ast.IF_Expression{Token: ps.currentToken}
	if !ps.peekTokenIs(token.LPAREN) {
		return nil
	}
	ps.advance()
	ps.advance()
	expr.Condition = ps._parseExpression(LOWEST)
	if !ps.peekTokenIs(token.RPAREN) {
		return nil
	}
	ps.advance()
	if !ps.peekTokenIs(token.LBRACE) {
		return nil
	}
	ps.advance()
	expr.Consequence = ps.parseBlockStatement()
	if ps.peekTokenIs(token.ELSE) {
		ps.advance()
		if !ps.peekTokenIs(token.LBRACE) {
			return nil
		}
		expr.Alternative = ps.parseBlockStatement()
	}
	return &expr
}

/** Parse Block Statement **/
func (ps *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	block.Statemens = []ast.Statement{}
	ps.advance()
	for !ps.currentTokenIs(token.RBRACE) && !ps.currentTokenIs(token.EOF) {
		stmt := ps.parseStatement()
		block.Statemens = append(block.Statemens, stmt)
		ps.advance()
	}
	return block
}

func (ps *Parser) parseFunctionLiteral() ast.Expression {
	lit := ast.FunctionLiteral{Token: ps.currentToken}
	if !ps.peekTokenIs(token.LPAREN) {
		return nil
	}
	ps.advance() // current token is LPAREN
	lit.Parameters = ps.parseFunctionParameters()
	if !ps.peekTokenIs(token.LBRACE) {
		return nil
	}
	ps.advance()
	lit.Body = ps.parseBlockStatement()
	return &lit
}

func (ps *Parser) parseFunctionParameters() []*ast.Identifier {
	// called when current token is token.IDENT
	identifiers := []*ast.Identifier{}
	// checks is ps.peekToken is RPAREN -> This may be the case if there are no identifiers
	if ps.peekTokenIs(token.RPAREN) {
		ps.advance()
		return identifiers
	}
	ps.advance()
	ident := &ast.Identifier{Token: ps.currentToken, Value: ps.currentToken.Literal}
	identifiers = append(identifiers, ident)
	for ps.peekTokenIs(token.COMMA) {
		ps.advance()
		ps.advance()
		ident := &ast.Identifier{Token: ps.currentToken, Value: ps.currentToken.Literal}
		identifiers = append(identifiers, ident)
	}
	if !ps.peekTokenIs(token.RPAREN) {
		return nil
	}
	ps.advance()
	return identifiers
}

/** Parse CALL Expression **/
func (ps *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expr := &ast.CALL_Expression{Token: ps.currentToken, Function: function}
	expr.Arguments = ps.parseCallArguments()
	return expr
}

func (ps *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}
	if ps.peekTokenIs(token.RPAREN) {
		ps.advance()
		return args
	}
	ps.advance()
	args = append(args, ps._parseExpression(LOWEST))
	for ps.peekTokenIs(token.COMMA) {
		ps.advance()
		ps.advance()
		args = append(args, ps._parseExpression(LOWEST))
	}
	if !ps.peekTokenIs(token.RPAREN) {
		return nil
	}
	ps.advance()
	return args
}

/** PRATT Parser **/
type (
	prattPrefixParsingFuncntion func() ast.Expression
	prattInfixParsingFunction   func(ast.Expression) ast.Expression
)

func (ps *Parser) addPrefixFn(fn prattPrefixParsingFuncntion, tokenType token.TokenType) {
	ps._prefixParsingFunctions[tokenType] = fn
}

func (ps *Parser) addInfixFn(fn prattInfixParsingFunction, tokenType token.TokenType) {
	ps._infixParsingFunctins[tokenType] = fn
}

/** The BINDING Powers **/
const (
	_ = iota
	LOWEST
	EQUALS
	LESS_GREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESS_GREATER,
	token.GT:       LESS_GREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

func (ps *Parser) _parseExpression(bindingPower int) ast.Expression {
	prefixFn := ps._prefixParsingFunctions[ps.currentToken.Type]
	if prefixFn == nil {
		return nil
	}
	leftExpr := prefixFn()
	for bindingPower < ps.peekPrecedence() {
		infixFn := ps._infixParsingFunctins[ps.peekToken.Type]
		if infixFn == nil {
			return leftExpr
		}
		ps.advance()
		leftExpr = infixFn(leftExpr)
	}
	return leftExpr
}

func (ps *Parser) peekPrecedence() int {
	if p, ok := precedences[ps.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (ps *Parser) currentPrecedence() int {
	if p, ok := precedences[ps.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}
