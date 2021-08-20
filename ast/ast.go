package ast

import (
	"bytes"

	"github.com/parsaakbari1209/interpreter/token"
)

type Node interface {
	//
	// TokenLiteral will only be used for debegging and testing.
	// It returns the literal value of the token it's associated with.
	//
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node

	// statmentNode is a dummy method to differantiate it from other nodes.
	statementNode()
}

type Expression interface {
	Node

	// expressionNode is a dummy method to differantiate it from other nodes.
	expressionNode()
}

// Root node of every AST that the parser produces.
type Program struct {
	// P lang is made of a series of statements.
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct {
	// token.LET
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	// token.RETURN
	Token token.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode() {
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	// The first token of the expression.
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type Identifier struct {
	// token.IDENT
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type Integer struct {
	// token.INT
	Token token.Token
	Value int64
}

func (i *Integer) expressionNode() {
}

func (i *Integer) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Integer) String() string {
	return i.Token.Literal
}

type Prefix struct {
	// token.BANG, token.MINUS
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *Prefix) expressionNode() {
}

func (p *Prefix) TokenLiteral() string {
	return p.Token.Literal
}

func (p *Prefix) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

type Infix struct {
	// token.PLUS, token.MINUS, token.ASTERISK, token.SLASH,
	// token.LT, token.GT, token.EQ, token.NOT_EQ
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (i *Infix) expressionNode() {
}

func (i *Infix) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Infix) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" ")
	out.WriteString(i.Operator)
	out.WriteString(" ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	// token.TRUE, token.FALSE
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {
}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}
