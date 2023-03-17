// Lexical analysis, lexing, or tokenization is the process of converting a sequence of
// characters into lexical tokens. A program that performs lexical analysis may be termed
// a lexer, tokenizer, or scanner, although scanner is also a term for the first stage of
// a lexer. A lexer will take an input character stream and convert it into the smallest
// meaningful characters called tokens. The stream of lexemes can be fed to a parser
// which will convert it into a parser tree.

package main

import (
	"bufio"
	"log"
	"strings"
	"unicode"
)

// ************
// ** Token **
// ************

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Delimiters
	COMMA    TokenType = ","
	COLON    TokenType = ":"
	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"
	LCURLY   TokenType = "{"
	RCURLY   TokenType = "}"

	// Identifiers and literals
	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"
	STRING TokenType = "STRING"

	// Keywords
	TRUE    TokenType = "TRUE"
	FALSE   TokenType = "FALSE"
	VAR     TokenType = "VAR"
	IF      TokenType = "IF"
	ELSE    TokenType = "ELSE"
	WHILE   TokenType = "WHILE"
	FOR     TokenType = "FOR"
	IN      TokenType = "IN"
	FN      TokenType = "FN"
	RETURN  TokenType = "RETURN"
	LEN     TokenType = "LEN"
	PRINT   TokenType = "PRINT"
	PRINTLN TokenType = "PRINTLN"

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	NOT      TokenType = "!"
	LT       TokenType = "<"
	GT       TokenType = ">"
	LEQ      TokenType = "<="
	GEQ      TokenType = ">="
	EQ       TokenType = "=="
	NEQ      TokenType = "!="
	OR       TokenType = "OR"
	AND      TokenType = "AND"
)

type Token struct {
	Type  TokenType
	Value string
}

func NewToken(t TokenType, v string) Token {
	return Token{Type: t, Value: v}
}

// ***********
// ** Lexer **
// ***********

type Lexer struct {
	reader *bufio.Reader
}

func NewLexer(in string) *Lexer {
	return &Lexer{reader: bufio.NewReader(strings.NewReader(in))}
}

func (l *Lexer) Lex() chan Token {
	tokens := make(chan Token)
	go func() {
		defer close(tokens)
		for {
			l.lexWhitespace()
			r := l.readRune()
			switch {
			case r == 0:
				tokens <- NewToken(EOF, "")
				return
			case r == '"':
				tokens <- l.lexString(r)
			case unicode.IsDigit(r):
				tokens <- l.lexNumber(r)
			case unicode.IsLetter(r):
				tokens <- l.lexIdentifier(r)
			default:
				token := l.lexSymbol(r)
				if token.Type == ILLEGAL {
					log.Fatalf("error: invalid token found: \"%s\"", token.Value)
				}
				tokens <- token
			}
		}
	}()
	return tokens
}

func (l *Lexer) lexIdentifier(r rune) Token {
	keywords := map[string]TokenType{
		"true":    TRUE,
		"false":   FALSE,
		"var":     VAR,
		"if":      IF,
		"else":    ELSE,
		"while":   WHILE,
		"for":     FOR,
		"in":      IN,
		"fn":      FN,
		"return":  RETURN,
		"len":     LEN,
		"print":   PRINT,
		"println": PRINTLN,
		"or":      OR,
		"and":     AND,
	}
	v := string(r)
	for {
		r = l.readRune()
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			l.unreadRune()
			break
		}
		v += string(r)
	}
	if t, ok := keywords[v]; ok {
		return NewToken(t, v)
	}
	return NewToken(IDENT, v)
}

func (l *Lexer) lexNumber(r rune) Token {
	t := INT
	v := string(r)
	for {
		r = l.readRune()
		if !unicode.IsDigit(r) && r != '.' {
			l.unreadRune()
			break
		}
		v += string(r)
		if r == '.' {
			t = FLOAT
		}
	}
	return NewToken(t, v)
}

func (l *Lexer) lexString(_ rune) Token {
	str, _ := l.reader.ReadString('"')
	return NewToken(STRING, strings.TrimRight(str, `"`))
}

func (l *Lexer) lexSymbol(r rune) Token {
	symbols := map[string]TokenType{
		",":  COMMA,
		":":  COLON,
		"(":  LPAREN,
		")":  RPAREN,
		"[":  LBRACKET,
		"]":  RBRACKET,
		"{":  LCURLY,
		"}":  RCURLY,
		"=":  ASSIGN,
		"+":  PLUS,
		"-":  MINUS,
		"*":  ASTERISK,
		"/":  SLASH,
		"!":  NOT,
		"<":  LT,
		">":  GT,
		"<=": LEQ,
		">=": GEQ,
		"==": EQ,
		"!=": NEQ,
	}
	singleCharSymbol := string(r)
	doubleCharSymbol := singleCharSymbol + string(l.readRune())
	if t, ok := symbols[doubleCharSymbol]; ok {
		return NewToken(t, doubleCharSymbol)
	}
	l.unreadRune()
	if t, ok := symbols[singleCharSymbol]; ok {
		return NewToken(t, singleCharSymbol)
	}
	return NewToken(ILLEGAL, singleCharSymbol)
}

func (l *Lexer) lexWhitespace() {
	r := l.readRune()
	for r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '#' {
		if r == '#' {
			l.reader.ReadLine()
		}
		r = l.readRune()
	}
	l.unreadRune()
}

func (l *Lexer) readRune() rune {
	r, _, _ := l.reader.ReadRune()
	return r
}

func (l *Lexer) unreadRune() {
	l.reader.UnreadRune()
}
