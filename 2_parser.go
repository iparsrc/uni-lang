// Parsing, or syntax analyzing is the process of converting tokens into corresponding
// Intermediate Representation(IR) using the programming language grammar.
// Top-down parser, and bottom-up parser are two main categories of parsers.
// Top-down parser is classified into two types: Recursive decent parser, and LL(1).
// Bottom-up parser is classified into two types: LR, and Operator precedence parser.
// LR parser is classified into four types: LR(0), SLR(1), LALR(1) and CLR(1).
//
// Parser
// ├── Top-down parser
// │   ├── Recursive decent parser
// │   └── LL(1)
// └── Bottom-up parser
//     ├── LR
//     │   ├── LR(0)
//     │   ├── SLR(1)
//     │   ├── LALR(1)
//     │   └── CLR(1)
//     └── Operator precedence parser

package main

import (
	"fmt"
	"log"
	"strconv"
)

const (
	LOWEST  = iota + 1
	EQUALS  // == !=
	BOOLOP  // or and
	GREATER // < > <= >=
	SUM     // + -
	PRODUCT // * /
	PREFIX  // +x -x !x
)

// ************
// ** Parser **
// ************

type Parser struct {
	tokens       chan Token
	currentToken Token
	peekToken    Token
	errors       []error
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{
		tokens: lexer.Lex(),
		errors: make([]error, 0),
	}
}

func (p *Parser) Parse() chan Statement {
	p.next() // initialize peek token
	p.next() // initialize current token
	statements := make(chan Statement)
	go func() {
		for p.currentToken.Type != EOF {
			if statement := p.parseStatement(); statement != nil {
				statements <- statement
			}
		}
		close(statements)
	}()
	return statements
}

// ****************
// ** Statements **
// ****************

type Statement interface{}

func (p *Parser) parseStatement() Statement {
	switch p.currentToken.Type {
	case VAR:
		return p.parseVariable()
	case IF:
		return p.parseIf()
	case WHILE:
		return p.parseWhile()
	case FOR:
		return p.parseFor()
	case FN:
		return p.parseFunction()
	case RETURN:
		return p.parseReturn()
	case LCURLY:
		return p.parseBlock()
	case IDENT:
		if p.peekToken.Type == ASSIGN {
			return p.parseVariable()
		}
		fallthrough
	default:
		return p.parseExpression(LOWEST)
	}
}

type Variable struct {
	Name  Identifier
	Value Expression
	IsNew bool
}

func (p *Parser) parseVariable() Statement {
	v := Variable{}
	if p.currentToken.Type == VAR {
		p.next() // skip var keyword
		v.IsNew = true
	}
	v.Name = p.parseIdentifier().(Identifier)
	if !p.expectCurrent(ASSIGN) {
		return nil
	}
	p.next() // skip = symbol
	v.Value = p.parseExpression(LOWEST)
	return v
}

type If struct {
	Condition   Expression
	Consequence Block
	Alternative *Block
}

func (p *Parser) parseIf() Statement {
	p.next() // skip if keyword
	i := If{}
	i.Condition = p.parseExpression(LOWEST)
	i.Consequence = p.parseBlock().(Block)
	if p.currentToken.Type == ELSE {
		p.next() // skip else keyword
		alternative := p.parseBlock().(Block)
		i.Alternative = &alternative
	}
	return i
}

type While struct {
	Condition   Expression
	Consequence Block
}

func (p *Parser) parseWhile() Statement {
	p.next() // skip while keyword
	w := While{Condition: p.parseExpression(LOWEST)}
	if !p.expectCurrent(LCURLY) {
		return nil
	}
	w.Consequence = p.parseBlock().(Block)
	return w
}

type For struct {
	Key         Identifier
	Value       Identifier
	Condition   Expression
	Consequence Block
}

func (p *Parser) parseFor() Statement {
	p.next() // skip for keyword
	f := For{Key: p.parseIdentifier().(Identifier)}
	if p.currentToken.Type == COMMA {
		p.next() // skip , symbol
		f.Value = p.parseIdentifier().(Identifier)
	}
	if !p.expectCurrent(IN) {
		return nil
	}
	p.next() // read in keyword
	f.Condition = p.parseExpression(LOWEST)
	if !p.expectCurrent(LCURLY) {
		return nil
	}
	f.Consequence = p.parseBlock().(Block)
	return f
}

type Function struct {
	Name       Identifier
	Parameters []Identifier
	Body       Block
}

func (p *Parser) parseFunction() Statement {
	p.next() // skip fn keyword
	f := Function{Name: p.parseIdentifier().(Identifier)}
	if !p.expectCurrent(LPAREN) {
		return nil
	}
	p.next() // skip ( symbol
	for p.currentToken.Type != RPAREN {
		f.Parameters = append(f.Parameters, p.parseIdentifier().(Identifier))
		if p.currentToken.Type == COMMA {
			p.next() // skip , symbol
		}
	}
	p.next() // skip ) symbol
	if !p.expectCurrent(LCURLY) {
		return nil
	}
	f.Body = p.parseBlock().(Block)
	return f
}

type Return struct {
	Value Expression
}

func (p *Parser) parseReturn() Statement {
	p.next() // skip return keyword
	r := Return{Value: p.parseExpression(LOWEST)}
	return r
}

type Block struct {
	Statements []Statement
}

func (p *Parser) parseBlock() Statement {
	p.next() // skip { symbol
	b := Block{}
	for p.currentToken.Type != RCURLY {
		b.Statements = append(b.Statements, p.parseStatement())
	}
	p.next() // skip } symbol
	return b
}

// *****************
// ** Expressions **
// *****************

type Expression interface{}

func (p *Parser) parseExpression(precedence int) Expression {
	var left Expression
	switch p.currentToken.Type {
	case TRUE, FALSE:
		left = p.parseBoolean()
	case INT:
		left = p.parseInteger()
	case FLOAT:
		left = p.parseFloat()
	case STRING:
		left = p.parseString()
	case PLUS, MINUS, NOT:
		left = p.parseUnaryOperation()
	case IDENT:
		left = p.parseIdentifier()
		switch p.currentToken.Type {
		case LBRACKET:
			left = p.parseIndex(left)
		case LPAREN:
			left = p.parseCall(left)
		}
	case LPAREN:
		p.next() // skip ( symbol
		left = p.parseExpression(LOWEST)
		p.next() // skip ) symbol
	case LBRACKET:
		left = p.parseArray()
	case LCURLY:
		left = p.parseMap()
	case LEN:
		left = p.parseLen()
	case PRINT, PRINTLN:
		left = p.parsePrint()
	default:
		p.errors = append(p.errors, fmt.Errorf("unary parse function for %s not found", p.currentToken.Type))
		return nil
	}
	for precedence < getPrecedence(p.currentToken.Type) {
		switch p.currentToken.Type {
		case OR, AND, PLUS, MINUS, ASTERISK, SLASH, EQ, NEQ, LT, GT, LEQ, GEQ:
			left = p.parseBinaryOperation(left)
		default:
			p.errors = append(p.errors, fmt.Errorf("binary parse function for %s not found", p.currentToken.Type))
			return nil
		}
	}
	return left
}

type Boolean struct {
	Value bool
}

func (p *Parser) parseBoolean() Expression {
	if !p.expectCurrent(TRUE, FALSE) {
		return nil
	}
	b := Boolean{}
	b.Value, _ = strconv.ParseBool(p.currentToken.Value)
	p.next() // skip boolean literal
	return b
}

type Integer struct {
	Value int64
}

func (p *Parser) parseInteger() Expression {
	if !p.expectCurrent(INT) {
		return nil
	}
	i := Integer{}
	i.Value, _ = strconv.ParseInt(p.currentToken.Value, 10, 64)
	p.next() // skip integer literal
	return i
}

type Float struct {
	Value float64
}

func (p *Parser) parseFloat() Expression {
	if !p.expectCurrent(FLOAT) {
		return nil
	}
	f := Float{}
	f.Value, _ = strconv.ParseFloat(p.currentToken.Value, 64)
	p.next() // skip float literal
	return f
}

type String struct {
	Value string
}

func (p *Parser) parseString() Expression {
	if !p.expectCurrent(STRING) {
		return nil
	}
	s := String{Value: p.currentToken.Value}
	p.next() // skip string literal
	return s
}

type Array struct {
	Items []Expression
}

func (p *Parser) parseArray() Expression {
	p.next() // skip [ symbol
	a := Array{Items: make([]Expression, 0)}
	for p.currentToken.Type != RBRACKET {
		a.Items = append(a.Items, p.parseExpression(LOWEST))
		if p.currentToken.Type == COMMA {
			p.next() // skip , symbol
		}
	}
	p.next() // skip ] symbol
	return a
}

type Map struct {
	Items map[Expression]Expression
}

func (p *Parser) parseMap() Expression {
	p.next() // skip { symbol
	m := Map{Items: make(map[Expression]Expression)}
	for p.currentToken.Type != RCURLY {
		key := p.parseExpression(LOWEST)
		if !p.expectCurrent(COLON) {
			return nil
		}
		p.next() // skip : symbol
		m.Items[key] = p.parseExpression(LOWEST)
		if p.currentToken.Type == COMMA {
			p.next() // skip , symbol
		}
	}
	p.next() // skip } symbol
	return m
}

type Index struct {
	Index   Expression
	Subject Expression
}

func (p *Parser) parseIndex(left Expression) Expression {
	p.next() // skip [ symbol
	i := Index{
		Index:   p.parseExpression(LOWEST),
		Subject: left,
	}
	p.next() // skip ] symbol
	return i
}

type Call struct {
	Identifier Identifier
	Arguments  []Expression
}

func (p *Parser) parseCall(left Expression) Expression {
	p.next() // skip ( symbol
	c := Call{Identifier: left.(Identifier), Arguments: make([]Expression, 0)}
	for p.currentToken.Type != RPAREN {
		c.Arguments = append(c.Arguments, p.parseExpression(LOWEST))
		if p.currentToken.Type == COMMA {
			p.next() // skip , symbol
		}
	}
	p.next() // skip ) symbol
	return c
}

type Identifier struct {
	Token          Token
	IsFunctionCall bool
}

func (p *Parser) parseIdentifier() Expression {
	if !p.expectCurrent(IDENT) {
		return nil
	}
	i := Identifier{Token: p.currentToken, IsFunctionCall: p.peekToken.Type == LPAREN}
	p.next() // skip identifier
	return i
}

type UnaryOperation struct {
	Token      Token
	Expression Expression
}

func (p *Parser) parseUnaryOperation() Expression {
	if !p.expectCurrent(PLUS, MINUS, NOT) {
		return nil
	}
	uo := UnaryOperation{Token: p.currentToken}
	p.next() // skip +, -, or !
	uo.Expression = p.parseExpression(PREFIX)
	return uo
}

type BinaryOperation struct {
	Token Token
	Left  Expression
	Right Expression
}

func (p *Parser) parseBinaryOperation(left Expression) Expression {
	bo := BinaryOperation{Token: p.currentToken, Left: left}
	precedence := getPrecedence(p.currentToken.Type)
	p.next() // skip operator(+, -, ...)
	bo.Right = p.parseExpression(precedence)
	return bo
}

type Len struct {
	Subject Expression
}

func (p *Parser) parseLen() Expression {
	p.next() // skip len keyword
	p.next() // skip ( symbol
	l := Len{Subject: p.parseExpression(LOWEST)}
	p.next() // skip ) symbol
	return l
}

type Print struct {
	Args      []Expression
	IsNewLine bool
}

func (p *Parser) parsePrint() Expression {
	print := Print{IsNewLine: p.currentToken.Type == PRINTLN}
	p.next() // skip print, or println keyword
	if !p.expectCurrent(LPAREN) {
		return nil
	}
	p.next() // skip ( symbol
	for p.currentToken.Type != RPAREN {
		print.Args = append(print.Args, p.parseExpression(LOWEST))
		if p.currentToken.Type == COMMA {
			p.next() // skip , symbol
		}
	}
	p.next() // skip ) symbol
	return print
}

func (p *Parser) expectCurrent(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.currentToken.Type == tokenType {
			return true
		}
	}
	if len(tokenTypes) == 1 {
		log.Fatalf("expected %s, got %s instead", tokenTypes[0], p.currentToken.Type)
		p.errors = append(p.errors, fmt.Errorf("expected %s, got %s instead", tokenTypes[0], p.currentToken.Type))
	} else {
		log.Fatalf("expected one of %s, got %s instead", tokenTypes, p.currentToken.Type)
		p.errors = append(p.errors, fmt.Errorf("expected one of %s, got %s instead", tokenTypes, p.currentToken.Type))
	}
	return false
}

func getPrecedence(in TokenType) int {
	precedences := map[TokenType]int{
		EQ:       EQUALS,
		NEQ:      EQUALS,
		OR:       BOOLOP,
		AND:      BOOLOP,
		LT:       GREATER,
		GT:       GREATER,
		LEQ:      GREATER,
		GEQ:      GREATER,
		PLUS:     SUM,
		MINUS:    SUM,
		ASTERISK: PRODUCT,
		SLASH:    PRODUCT,
	}
	if precedence, ok := precedences[in]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) next() {
	p.currentToken = p.peekToken
	p.peekToken = <-p.tokens
}
