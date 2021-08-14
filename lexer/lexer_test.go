package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/parsaakbari1209/interpreter/token"
)

func TestNextToken(t *testing.T) {
	tt := []struct {
		name       string
		sourceCode string
		tokens     []token.Token
	}{
		{
			name:       "operators",
			sourceCode: "=+",
			tokens: []token.Token{
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:       "parentheses",
			sourceCode: "()",
			tokens: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:       "braces",
			sourceCode: "{}",
			tokens: []token.Token{
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:       "comma",
			sourceCode: ",",
			tokens: []token.Token{
				{Type: token.COMMA, Literal: ","},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:       "semicolon",
			sourceCode: ";",
			tokens: []token.Token{
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:       "sample plang 1",
			sourceCode: `let five = 5; let ten = 10; let add = fn(x, y) { x + y; }; let result = add(five, ten);`,
			tokens: []token.Token{
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "five"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.INT, Literal: "5"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "ten"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.INT, Literal: "10"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "add"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.FUNCITON, Literal: "fn"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENT, Literal: "y"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.IDENT, Literal: "y"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "result"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.IDENT, Literal: "add"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.IDENT, Literal: "five"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENT, Literal: "ten"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:       "single character symbols",
			sourceCode: "=+!-/*><",
			tokens: []token.Token{
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.BANG, Literal: "!"},
				{Type: token.MINUS, Literal: "-"},
				{Type: token.SLASH, Literal: "/"},
				{Type: token.ASTERISK, Literal: "*"},
				{Type: token.GT, Literal: ">"},
				{Type: token.LT, Literal: "<"},
			},
		},
		{
			name:       "keywords",
			sourceCode: "fn let if else return true false",
			tokens: []token.Token{
				{Type: token.FUNCITON, Literal: "fn"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IF, Literal: "if"},
				{Type: token.ELSE, Literal: "else"},
				{Type: token.RETURN, Literal: "return"},
				{Type: token.TRUE, Literal: "true"},
				{Type: token.FALSE, Literal: "false"},
			},
		},
		{
			name:       "sample_plang_2",
			sourceCode: `if (5 < 7) { return true; } else { return false; }`,
			tokens: []token.Token{
				{Type: token.IF, Literal: "if"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.INT, Literal: "5"},
				{Type: token.LT, Literal: "<"},
				{Type: token.INT, Literal: "7"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RETURN, Literal: "return"},
				{Type: token.TRUE, Literal: "true"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.ELSE, Literal: "else"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RETURN, Literal: "return"},
				{Type: token.FALSE, Literal: "false"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:       "two character operators",
			sourceCode: "== !=",
			tokens: []token.Token{
				{Type: token.EQ, Literal: "=="},
				{Type: token.NOT_EQ, Literal: "!="},
				{Type: token.EOF, Literal: ""},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			lxr := New(tc.sourceCode)

			for _, want := range tc.tokens {
				got := lxr.NextToken()

				assert.EqualValues(t, want, got)
			}
		})
	}
}
