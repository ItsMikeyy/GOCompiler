package parser

import (
	"fmt"
	"mscript/ast"
	"mscript/lexer"
	"mscript/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      //=
	LESSGREATER //> or <
	SUM         //+
	PRODUCT     // *
	PREFIX      //Prefix
	CALL        //myFunction(x)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

type Parser struct {
	l      *lexer.Lexer //Copy of lexer
	errors []string     //An array of errors collected along the way

	curToken  token.Token //Current token parsing
	peekToken token.Token //Next token parsing

	//Check if cur token has a parsing function associated
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression //Argument is the left side of the infix operator
)

// Creates a new instance of Parser
// Has a copy of the lexer the current token and the next token
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	//INIT map
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)

	//Add parse identifer for IDENT
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	//Add parse integer for INT
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	//Add parse prefix for BANG and MINUS
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)

	//Add parse infix functions
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// advances both cur and peek token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{} //Construct root node
	program.Statements = []ast.Statement{}

	//Loop until EOF token
	for p.curToken.Type != token.EOF {
		//Get statement
		stmt := p.parseStatement()

		//Check if stmt != EOF
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	//Return ast root
	return program
}

// Checking the type of statement we need to parse and returning the resulting statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.ParseReturnStatement()
	default:
		return p.ParseExpressionStatement()
	}
}

func (p *Parser) parseStatements() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

// Handles parsing let statements
func (p *Parser) parseLetStatement() *ast.LetStatement {
	//Creating let statement with the current Token = to cur token (let)
	stmt := &ast.LetStatement{Token: p.curToken}

	//If what follows is not an identifier return nil
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	//Set the identifers token to cur token and the value of the current tokens literal (var name)
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	//If an equal sign does not follow the identifer return nil
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//While the current token is not a semicolon advance the tokens
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Parses ReturnStatement
// Returns ReturnStatement
func (p *Parser) ParseReturnStatement() *ast.ReturnStatement {
	//Creating ReturnStatement
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	//Loop through expression
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Parse Expression Statement
func (p *Parser) ParseExpressionStatement() *ast.ExpressionStatement {
	//Create Expression Statement
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	//TEMP
	stmt.Expression = p.parseExpression(LOWEST)

	//Advance if peek is ;
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Parse prefix left hand side of expression if prefix
func (p *Parser) parseExpression(precednce int) ast.Expression {
	//Get prefix function for cur token
	prefix := p.prefixParseFns[p.curToken.Type]

	//No function exists add error message and return nil
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	//Parse left side with prefix function
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precednce < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as interger", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// create prefix expression and call parseExpression
func (p *Parser) parsePrefixExpression() ast.Expression {
	//create PrefixExpression
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	//Advance
	p.nextToken()

	//Handle prefix parsing
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// Create Infix expression get operator precedence and call parseExpression on the right
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	//get operator precendence
	precedence := p.curPrecedence()

	//advance
	p.nextToken()

	//parse right
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// Just parse the inner expression
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	//Check for RPAREN
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// Handle if expression
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	//if ( <exp>
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()

	//parse if (<exp>)
	expression.Condition = p.parseExpression(LOWEST)

	//check for <exp> )
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	//check for (<exp>) {
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	//parse inner if code
	expression.Consequence = p.parseBlockStatement()

	//If there is else block
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		//Prase inner else statements
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// Helper for adding a function to to prefix map
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// Helper for adding a function to to infix map
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// Add error for unknown prefix function
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s", t)
	p.errors = append(p.errors, msg)
}

// Check if current token is equal to expected
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// Check what the next token is and compare what the expected token is
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// calls peekTokenIs and handles return logic
func (p *Parser) expectPeek(t token.TokenType) bool {

	if p.peekTokenIs(t) { //Next token is = to expected
		p.nextToken() //Advance
		return true
	} else { //Next token != to expected
		p.peekError(t) //Add error
		return false
	}
}

// Check what the precedence is of the next token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Check what the precedence is of the current token
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// Return the array of errors
func (p *Parser) Errors() []string {
	return p.errors
}

// Add error message to p.errors
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)

}
