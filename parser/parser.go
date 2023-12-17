package parser

import (
	"fmt"
	"strconv"

	"bariq/ast"
	"bariq/lexer"
	"bariq/token"
)

type Parser struct {
	l *lexer.Lexer

	errors    []string
	curToken  token.Token
	peekToken token.Token

	// used to check if the token has a a
	// prefix or infix functinon associated with it
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: make([]string, 0)}
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdent)
	p.registerPrefix(token.INT, p.parseIntLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpr)
	p.registerPrefix(token.MINUS, p.parsePrefixExpr)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpr)
	p.registerPrefix(token.ASYNC, p.parseAsyncFunctionLiteral)
	// p.registerPrefix(token.GENERATOR, p.parseGeneratorFunctionLiteral)
	p.registerPrefix(token.AWAIT, p.parseAwaitExpr)
	p.registerPrefix(token.YIELD, p.parseYieldExpr)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpr)
	p.registerInfix(token.MINUS, p.parseInfixExpr)
	p.registerInfix(token.GT, p.parseInfixExpr)
	p.registerInfix(token.LT, p.parseInfixExpr)
	p.registerInfix(token.SLASH, p.parseInfixExpr)
	p.registerInfix(token.ASTERIK, p.parseInfixExpr)
	p.registerInfix(token.EQ, p.parseInfixExpr)
	p.registerInfix(token.NEQ, p.parseInfixExpr)
	p.registerInfix(token.LPAREN, p.parseCallExpr)
	p.registerInfix(token.LBRACKET, p.parseIndexExpr)
	// read tow token so next and peek are set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) parseIndexExpr(left ast.Expr) ast.Expr {
	exp := &ast.IndexExpr{Token: p.curToken, Left: left}
	p.nextToken()
	exp.Index = p.parseCurrExpr(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

func (p *Parser) parseCallExpr(function ast.Expr) ast.Expr {
	exp := &ast.CallExpr{Token: p.curToken, Function: function}
	exp.Args = p.parseExprList(token.RPAREN)
	return exp
}

func (p *Parser) parseExprList(end token.TokenType) []ast.Expr {
	list := []ast.Expr{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseCurrExpr(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseCurrExpr(LOWEST))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return list
}

func (p *Parser) parseAsyncFunctionLiteral() ast.Expr {
	if !p.expectPeek(token.FUNCTION) {
		return nil
	}
	lit := &ast.FunctionLiteral{Async: true, Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParams()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStmt()
	return lit
}

func (p *Parser) parseFunctionLiteral() ast.Expr {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if p.peekTokenIs(token.GENERATOR) {
		lit.Gen = true
		p.nextToken()
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParams()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStmt()
	fmt.Println(lit.Gen)
	return lit
}

func (p *Parser) parseFunctionParams() []*ast.Ident {
	idents := []*ast.Ident{}
	// no params
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return idents
	}
	p.nextToken()
	ident := &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}
	idents = append(idents, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}
		idents = append(idents, ident)
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return idents
}

func (p *Parser) parseYieldExpr() ast.Expr {
	expr := &ast.YieldExpr{Token: p.curToken}
	p.nextToken()
	expr.Arg = p.parseCurrExpr(LOWEST)
	fmt.Printf("yield expr: %v\n", expr)
	return expr
}

func (p *Parser) parseAwaitExpr() ast.Expr {
	expr := &ast.AwaitExpr{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expr.Arg = p.parseCurrExpr(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseIfExpr() ast.Expr {
	expr := &ast.IfExpr{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expr.Condition = p.parseCurrExpr(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expr.Consequence = p.parseBlockStmt()
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expr.Alternative = p.parseBlockStmt()
	}
	return expr
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	block := &ast.BlockStmt{Token: p.curToken}
	block.Stmts = []ast.Stmt{}
	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Stmts = append(block.Stmts, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseGroupedExpression() ast.Expr {
	p.nextToken()
	exp := p.parseCurrExpr(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseExprStmt() *ast.ExprStmt {
	stmt := ast.ExprStmt{Token: p.curToken}
	stmt.Expr = p.parseCurrExpr(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return &stmt
}

// NOTE: take your time undestanding this
func (p *Parser) parseCurrExpr(precedence int) ast.Expr {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseInfixExpr(left ast.Expr) ast.Expr {
	exp := &ast.InfixExpr{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	// to make + right-associative. Ex: 1 + 2 + 3 -> (1 + (2 + 3))
	// if exp.Operator == "+" {
	// 	exp.Right = p.parseExpr(precedence - 1)
	// }
	exp.Right = p.parseCurrExpr(precedence)
	return exp
}

func (p *Parser) parsePrefixExpr() ast.Expr {
	exp := &ast.PrefixExpr{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	exp.Right = p.parseCurrExpr(PREFIX)
	return exp
}

func (p *Parser) parseBoolean() ast.Expr {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseIdent() ast.Expr {
	return &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseHashLiteral() ast.Expr {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expr]ast.Expr)
	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		k := p.parseCurrExpr(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()
		v := p.parseCurrExpr(LOWEST)
		hash.Pairs[k] = v
		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return hash
}

func (p *Parser) parseArrayLiteral() ast.Expr {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elmnts = p.parseExprList(token.RBRACKET)
	return array
}

func (p *Parser) parseStringLiteral() ast.Expr {
	lit := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	return lit
}

func (p *Parser) parseIntLiteral() ast.Expr {
	lit := &ast.IntLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("couldn't parse %q  as interger", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix function for %s found ", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf(
		"expected the next token to be %s, got %s",
		t,
		p.peekToken.Type,
	)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Stmts = []ast.Stmt{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Stmts = append(program.Stmts, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Stmt {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStmt()
	case token.RET:
		return p.parseRetStmt()
	default:
		// fmt.Println(p.curToken.Type, p.curToken.Literal)
		return p.parseExprStmt()
		// return nil
	}
}

// precence order
const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREETER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var precedence = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.MINUS:    SUM,
	token.PLUS:     SUM,
	token.GT:       LESSGREETER,
	token.LT:       LESSGREETER,
	token.SLASH:    PRODUCT,
	token.ASTERIK:  PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedence[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedence[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseRetStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{Token: p.curToken}
	p.nextToken()
	stmt.Value = p.parseCurrExpr(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	stmt := &ast.LetStmt{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Ident{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseCurrExpr(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

type (
	prefixParseFn func() ast.Expr
	infixParseFn  func(ast.Expr) ast.Expr
)

func (p *Parser) registerPrefix(tt token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}

func (p *Parser) registerInfix(tt token.TokenType, fn infixParseFn) {
	p.infixParseFns[tt] = fn
}
