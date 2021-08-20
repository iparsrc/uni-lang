package parser

import (
	"fmt"
	"strconv"

	"github.com/parsaakbari1209/interpreter/ast"
	"github.com/parsaakbari1209/interpreter/lexer"
	"github.com/parsaakbari1209/interpreter/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	CALL        // myFunction(x)
)

var (
	precedences = map[token.TokenType]int{
		token.EQ:       EQUALS,
		token.NOT_EQ:   EQUALS,
		token.LT:       LESSGREATER,
		token.GT:       LESSGREATER,
		token.PLUS:     SUM,
		token.MINUS:    SUM,
		token.SLASH:    PRODUCT,
		token.ASTERISK: PRODUCT,
	}
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lxr            *lexer.Lexer
	curToken       token.Token
	peekToken      token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(lxr *lexer.Lexer) *Parser {
	p := &Parser{lxr: lxr, errors: []string{}}
	p.nextToken()
	p.nextToken()
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerPrefixParseFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixParseFn(token.INT, p.parseInteger)
	p.registerPrefixParseFn(token.TRUE, p.parseBoolean)
	p.registerPrefixParseFn(token.FALSE, p.parseBoolean)
	p.registerPrefixParseFn(token.MINUS, p.parsePrefix)
	p.registerPrefixParseFn(token.BANG, p.parsePrefix)
	p.registerInfixParseFn(token.EQ, p.parseInfix)
	p.registerInfixParseFn(token.NOT_EQ, p.parseInfix)
	p.registerInfixParseFn(token.LT, p.parseInfix)
	p.registerInfixParseFn(token.GT, p.parseInfix)
	p.registerInfixParseFn(token.PLUS, p.parseInfix)
	p.registerInfixParseFn(token.MINUS, p.parseInfix)
	p.registerInfixParseFn(token.SLASH, p.parseInfix)
	p.registerInfixParseFn(token.ASTERISK, p.parseInfix)
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}
	for !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) GetErrors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lxr.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	// TODO: Handle expressions.
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	// TODO: Handle experssions.
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixParseFn := p.prefixParseFns[p.curToken.Type]
	if prefixParseFn == nil {
		msg := fmt.Sprintf("no prefix parse function for %s found", p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	left := prefixParseFn()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixParseFn := p.infixParseFns[p.peekToken.Type]
		if infixParseFn == nil {
			return left
		}
		p.nextToken()
		left = infixParseFn(left)
	}
	return left
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseInteger() ast.Expression {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}
	return &ast.Integer{Token: p.curToken, Value: value}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.currentTokenIs(token.TRUE)}
}

func (p *Parser) parsePrefix() ast.Expression {
	expression := &ast.Prefix{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	expression := &ast.Infix{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
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

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefixParseFn(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixParseFn(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
