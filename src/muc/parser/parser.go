package parser

import (
	"fmt"
	"muc/ast"
	"muc/lexer"
	"muc/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS			// ==
	LESSGREATER		// < or >
	SUM				// + or -
	PRODUCT			// * or /
	PREFIX			// -x or !x
	CALL			// myFunc(args)
	INDEX			// array[index]	
)

var precedences = map[token.TokenType]int{
	token.EQUAL:	 EQUALS,
	token.NOT_EQ:	 EQUALS,
	token.LESS:		 LESSGREATER,
	token.GREATER:	 LESSGREATER,
	token.PLUS:		 SUM,
	token.MINUS:	 SUM,
	token.SLASH:	 PRODUCT,
	token.ASTERISK:  PRODUCT,
	token.L_PAREN:   CALL,
	token.L_BRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer
	errors []string

	currToken token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l     : l,
		errors: []string{},

		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	p.registerPrefix(token.ID, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.L_BRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.L_BRACE, p.parseHashLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.L_PAREN, p.parseGroupExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.MACRO, p.parseMacroLiteral)
	
	for _, tok := range []token.TokenType{token.PLUS, token.MINUS, token.SLASH,
		token.ASTERISK, token.EQUAL, token.NOT_EQ, token.LESS, token.GREATER} {
		p.registerInfix(tok, p.parseInfixExpression)
	}
	p.registerInfix(token.L_PAREN, p.parseCallExpression)
	p.registerInfix(token.L_BRACKET, p.parseIndexExpression)
	
	// Next twice, set currToken and peekToken.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// foobar;
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

// 5;
func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.currToken}
	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	
	literal.Value = value
	return literal
}

// true;
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()

	expr := p.parseExpression(LOWEST)
	if !p.expectPeek(token.R_PAREN) {
		return nil
	}
	return expr
}

// "hello world"
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}
}

// [1, 2, 3]
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.currToken}
	array.Elements = p.parseExpressionList(token.R_BRACKET)
	return array
}

// myArray[2]
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{Token: p.currToken, Left: left}
	p.nextToken()
	expr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.R_BRACKET) {
		return nil
	}
	return expr
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	p.nextToken()
	if p.currTokenIs(end) {
		return list
	}
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()		// point to comma
		p.nextToken()		// skip comma
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return list
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.currToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.R_BRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value

		if !p.peekTokenIs(token.R_BRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.R_BRACE) {
		return nil
	}
	return hash
}

// !5;
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// 5 + 5;
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token: p.currToken,
		Operator: p.currToken.Literal,
		Left: left,
	}
	precedence := p.currPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		stmt := p.ParseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) ParseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.SEMICOLON:
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

// let id = expr
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(token.ID) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// return expr;
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// if <cond> <conseq> else <alternative>
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currToken}
	if !p.expectPeek(token.L_PAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.R_PAREN) {
		return nil
	}
	if !p.expectPeek(token.L_BRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.L_BRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

// expr;
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// optional semicolon
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for `%s` found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixF := p.prefixParseFns[p.currToken.Type]
	if prefixF == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}
	leftExpr := prefixF()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixF := p.infixParseFns[p.peekToken.Type]
		if infixF == nil {
			return leftExpr
		}

		p.nextToken()

		leftExpr = infixF(leftExpr)
	}
	
	return leftExpr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currToken}
	block.Statements = []ast.Statement{}
	p.nextToken()

	for !p.currTokenIs(token.R_BRACE) {
		stmt := p.ParseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: p.currToken}

	if !p.expectPeek(token.L_PAREN) {
		return nil
	}
	literal.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.L_BRACE) {
		return nil
	}
	literal.Body = p.parseBlockStatement()
	return literal
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers  := []*ast.Identifier{}
	if p.peekTokenIs(token.R_PAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	
	identifiers = append(identifiers,  
		&ast.Identifier{Token: p.currToken, Value: p.currToken.Literal})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		identifiers = append(identifiers, 
			&ast.Identifier{Token: p.currToken, Value: p.currToken.Literal})
	}
	for !p.expectPeek(token.R_PAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expr := &ast.CallExpression{Token: p.currToken, Function: function}
	expr.Arguments = p.parseExpressionList(token.R_PAREN)
	return expr
}

func (p *Parser) parseMacroLiteral() ast.Expression {
	literal := &ast.MacroLiteral{Token: p.currToken}

	if !p.expectPeek(token.L_PAREN) { return nil }
	literal.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.L_BRACE) { return nil }

	literal.Body = p.parseBlockStatement()
	return literal
}