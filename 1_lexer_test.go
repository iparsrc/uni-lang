package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	tt := []struct {
		name string
		in   string
		want []Token
	}{
		{
			name: "comment",
			in:   "# Cool!",
			want: []Token{
				{Type: EOF, Value: ""},
			},
		},
		{
			name: "delimiters",
			in:   ", : ( ) [ ] { }",
			want: []Token{
				{Type: COMMA, Value: ","},
				{Type: COLON, Value: ":"},
				{Type: LPAREN, Value: "("},
				{Type: RPAREN, Value: ")"},
				{Type: LBRACKET, Value: "["},
				{Type: RBRACKET, Value: "]"},
				{Type: LCURLY, Value: "{"},
				{Type: RCURLY, Value: "}"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name: "identifiers and literals",
			in:   `abc 0 1.5 "abc"`,
			want: []Token{
				{Type: IDENT, Value: "abc"},
				{Type: INT, Value: "0"},
				{Type: FLOAT, Value: "1.5"},
				{Type: STRING, Value: "abc"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name: "keywords",
			in:   `true false var if else while for in fn return len print println`,
			want: []Token{
				{Type: TRUE, Value: "true"},
				{Type: FALSE, Value: "false"},
				{Type: VAR, Value: "var"},
				{Type: IF, Value: "if"},
				{Type: ELSE, Value: "else"},
				{Type: WHILE, Value: "while"},
				{Type: FOR, Value: "for"},
				{Type: IN, Value: "in"},
				{Type: FN, Value: "fn"},
				{Type: RETURN, Value: "return"},
				{Type: LEN, Value: "len"},
				{Type: PRINT, Value: "print"},
				{Type: PRINTLN, Value: "println"},
				{Type: EOF, Value: ""},
			},
		},
		{
			name: "operators",
			in:   `= + - * / ! < > <= >= == != or and`,
			want: []Token{
				{Type: ASSIGN, Value: "="},
				{Type: PLUS, Value: "+"},
				{Type: MINUS, Value: "-"},
				{Type: ASTERISK, Value: "*"},
				{Type: SLASH, Value: "/"},
				{Type: NOT, Value: "!"},
				{Type: LT, Value: "<"},
				{Type: GT, Value: ">"},
				{Type: LEQ, Value: "<="},
				{Type: GEQ, Value: ">="},
				{Type: EQ, Value: "=="},
				{Type: NEQ, Value: "!="},
				{Type: OR, Value: "or"},
				{Type: AND, Value: "and"},
				{Type: EOF, Value: ""},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			lexer := NewLexer(tc.in)
			tokens := lexer.Lex()
			for _, want := range tc.want {
				got := <-tokens
				assert.Equal(t, want, got)
			}
		})
	}
}
