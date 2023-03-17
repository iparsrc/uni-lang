package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	tt := []struct {
		name string
		in   string
		want []Statement
	}{
		{
			name: "expression 1",
			in:   "(1 + 2) * 3",
			want: []Statement{
				BinaryOperation{
					Token: NewToken(ASTERISK, "*"),
					Left: BinaryOperation{
						Token: NewToken(PLUS, "+"),
						Left:  Integer{Value: 1},
						Right: Integer{Value: 2},
					},
					Right: Integer{Value: 3},
				},
			},
		},
		{
			name: "expression 2",
			in:   "sum(1, 2)",
			want: []Statement{
				Call{
					Identifier: Identifier{
						Token:          NewToken(IDENT, "sum"),
						IsFunctionCall: true,
					},
					Arguments: []Expression{
						Integer{Value: 1},
						Integer{Value: 2},
					},
				},
			},
		},
		{
			name: "variable 1",
			in:   "var a = 0",
			want: []Statement{
				Variable{
					Name: Identifier{
						Token:          NewToken(IDENT, "a"),
						IsFunctionCall: false,
					},
					Value: Integer{Value: 0},
					IsNew: true,
				},
			},
		},
		{
			name: "variable 2",
			in:   "var a = 0.0",
			want: []Statement{
				Variable{
					Name: Identifier{
						Token:          NewToken(IDENT, "a"),
						IsFunctionCall: false,
					},
					Value: Float{Value: 0.0},
					IsNew: true,
				},
			},
		},
		{
			name: "variable 3",
			in:   "var a = \"Hello World!\"",
			want: []Statement{
				Variable{
					Name: Identifier{
						Token:          NewToken(IDENT, "a"),
						IsFunctionCall: false,
					},
					Value: String{Value: "Hello World!"},
					IsNew: true,
				},
			},
		},
		{
			name: "condition 1",
			in:   "if true {}",
			want: []Statement{
				If{
					Condition:   Boolean{Value: true},
					Consequence: Block{},
				},
			},
		},
		{
			name: "condition 2",
			in:   "if 1 == 2 {} else {}",
			want: []Statement{
				If{
					Condition: BinaryOperation{
						Token: NewToken(EQ, "=="),
						Left:  Integer{Value: 1},
						Right: Integer{Value: 2},
					},
					Consequence: Block{},
					Alternative: &Block{},
				},
			},
		},
		{
			name: "while 1",
			in:   "while true {}",
			want: []Statement{
				While{
					Condition:   Boolean{Value: true},
					Consequence: Block{},
				},
			},
		},
		{
			name: "while 2",
			in:   "while a == 1 {}",
			want: []Statement{
				While{
					Condition: BinaryOperation{
						Token: NewToken(EQ, "=="),
						Left: Identifier{
							Token:          NewToken(IDENT, "a"),
							IsFunctionCall: false,
						},
						Right: Integer{Value: 1},
					},
					Consequence: Block{},
				},
			},
		},
		{
			name: "for 1",
			in:   "for k, v in \"Hello World!\" {}",
			want: []Statement{
				For{
					Key: Identifier{
						Token:          NewToken(IDENT, "k"),
						IsFunctionCall: false,
					},
					Value: Identifier{
						Token:          NewToken(IDENT, "v"),
						IsFunctionCall: false,
					},
					Condition:   String{Value: "Hello World!"},
					Consequence: Block{},
				},
			},
		},
		{
			name: "for 2",
			in:   "for k, v in [0.1, 0.2] {}",
			want: []Statement{
				For{
					Key: Identifier{
						Token:          NewToken(IDENT, "k"),
						IsFunctionCall: false,
					},
					Value: Identifier{
						Token:          NewToken(IDENT, "v"),
						IsFunctionCall: false,
					},
					Condition: Array{
						Items: []Expression{
							Float{Value: 0.1},
							Float{Value: 0.2},
						},
					},
					Consequence: Block{},
				},
			},
		},
		{
			name: "for 3",
			in:   "for k, v in {\"one\": 1} {}",
			want: []Statement{
				For{
					Key: Identifier{
						Token:          NewToken(IDENT, "k"),
						IsFunctionCall: false,
					},
					Value: Identifier{
						Token:          NewToken(IDENT, "v"),
						IsFunctionCall: false,
					},
					Condition: Map{
						Items: map[Expression]Expression{
							String{Value: "one"}: Integer{Value: 1},
						},
					},
					Consequence: Block{},
				},
			},
		},
		{
			name: "function 1",
			in:   "fn sum(a, b) { return a + b }",
			want: []Statement{
				Function{
					Name: Identifier{
						Token:          NewToken(IDENT, "sum"),
						IsFunctionCall: true,
					},
					Parameters: []Identifier{
						{
							Token:          NewToken(IDENT, "a"),
							IsFunctionCall: false,
						},
						{
							Token:          NewToken(IDENT, "b"),
							IsFunctionCall: false,
						},
					},
					Body: Block{
						Statements: []Statement{
							Return{
								Value: BinaryOperation{
									Token: NewToken(PLUS, "+"),
									Left: Identifier{
										Token:          NewToken(IDENT, "a"),
										IsFunctionCall: false,
									},
									Right: Identifier{
										Token:          NewToken(IDENT, "b"),
										IsFunctionCall: false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "block 1",
			in:   "{ xy = 0 if true {}}",
			want: []Statement{
				Block{
					Statements: []Statement{
						Variable{
							Name: Identifier{
								Token:          NewToken(IDENT, "xy"),
								IsFunctionCall: false,
							},
							Value: Integer{Value: 0},
						},
						If{
							Condition:   Boolean{Value: true},
							Consequence: Block{},
						},
					},
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			lexer := NewLexer(tc.in)
			parser := NewParser(lexer)
			statements := parser.Parse()
			for _, want := range tc.want {
				got := <-statements
				assert.Equal(t, want, got)
			}
		})
	}
}
