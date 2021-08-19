package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/parsaakbari1209/interpreter/ast"
	"github.com/parsaakbari1209/interpreter/lexer"
	"github.com/parsaakbari1209/interpreter/token"
)

func TestString(t *testing.T) {
	tt := []struct {
		name    string
		program *ast.Program
		want    string
	}{
		{
			name: "one let statement",
			program: &ast.Program{
				Statements: []ast.Statement{
					&ast.LetStatement{
						Token: token.Token{Type: token.LET, Literal: "let"},
						Name: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Value: "x",
						},
						Value: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "y"},
							Value: "y",
						},
					},
				},
			},
			want: `let x = y;`,
		},
		{
			name: "one return statement",
			program: &ast.Program{
				Statements: []ast.Statement{
					&ast.ReturnStatement{
						Token: token.Token{Type: token.RETURN, Literal: "return"},
						Value: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Value: "x",
						},
					},
				},
			},
			want: `return x;`,
		},
		{
			name: "",
			program: &ast.Program{
				Statements: []ast.Statement{
					&ast.LetStatement{
						Token: token.Token{Type: token.LET, Literal: "let"},
						Name: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Value: "x",
						},
						Value: &ast.Integer{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
			want: `let x = 5;`,
		},
		{
			name: "",
			program: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{Type: token.MINUS, Literal: "-"},
						Expression: &ast.Prefix{
							Token:    token.Token{Type: token.MINUS, Literal: "-"},
							Operator: "-",
							Right: &ast.Integer{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
						},
					},
				},
			},
			want: "(-5)",
		},
		{
			name: "",
			program: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{Type: token.BANG, Literal: "!"},
						Expression: &ast.Prefix{
							Token:    token.Token{Type: token.BANG, Literal: "!"},
							Operator: "!",
							Right: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
						},
					},
				},
			},
			want: "(!x)",
		},
		{
			name: "",
			program: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Expression: &ast.Infix{
							Token:    token.Token{Type: token.INT, Literal: "5"},
							Operator: "+",
							Left: &ast.Integer{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
							Right: &ast.Integer{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
						},
					},
				},
			},
			want: "(5 + 5)",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.program.String()
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestParseLetStatement(t *testing.T) {
	tt := []struct {
		name       string
		sourceCode string
		want       *ast.Program
		wantErr    bool
	}{
		{
			name:       "one let statement",
			sourceCode: `let x = 5;`,
			want: &ast.Program{
				Statements: []ast.Statement{
					&ast.LetStatement{
						Token: token.Token{Type: token.LET, Literal: "let"},
						Name: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Value: "x",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "two let statements",
			sourceCode: `let x = 5; let y = 10;`,
			want: &ast.Program{
				Statements: []ast.Statement{
					&ast.LetStatement{
						Token: token.Token{Type: token.LET, Literal: "let"},
						Name: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Value: "x",
						},
					},
					&ast.LetStatement{
						Token: token.Token{Type: token.LET, Literal: "let"},
						Name: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "y"},
							Value: "y",
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			lxr := lexer.New(tc.sourceCode)
			psr := New(lxr)
			if tc.wantErr {
				assert.NotEmpty(t, psr.GetErrors())
			}
			if !tc.wantErr {
				assert.Empty(t, psr.GetErrors())
			}
			got := psr.ParseProgram()
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestReturnStatement(t *testing.T) {
	tt := []struct {
		name       string
		sourceCode string
		want       *ast.Program
		wantErr    bool
	}{
		{
			name:       "one return statement",
			sourceCode: `return 5;`,
			want: &ast.Program{
				Statements: []ast.Statement{
					&ast.ReturnStatement{
						Token: token.Token{Type: token.RETURN, Literal: "return"},
					},
				},
			},
		},
		{
			name:       "two return statements",
			sourceCode: `return 5; return 10;`,
			want: &ast.Program{
				Statements: []ast.Statement{
					&ast.ReturnStatement{
						Token: token.Token{Type: token.RETURN, Literal: "return"},
					},
					&ast.ReturnStatement{
						Token: token.Token{Type: token.RETURN, Literal: "return"},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			lxr := lexer.New(tc.sourceCode)
			psr := New(lxr)
			got := psr.ParseProgram()
			if tc.wantErr {
				assert.NotEmpty(t, psr.GetErrors())
			}
			if !tc.wantErr {
				assert.Empty(t, psr.GetErrors())
			}
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestIdentifierExpression(t *testing.T) {
	tt := []struct {
		name       string
		sourceCode string
		want       *ast.Program
		wantErr    bool
	}{
		{
			name:       "one identifier as expression",
			sourceCode: "x;",
			want: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Expression: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Value: "x",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "one identifier as expression without semicolon",
			sourceCode: "x",
			want: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Expression: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Value: "x",
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			lxr := lexer.New(tc.sourceCode)
			psr := New(lxr)
			got := psr.ParseProgram()
			if tc.wantErr {
				assert.NotEmpty(t, psr.GetErrors())
			}
			if !tc.wantErr {
				assert.Empty(t, psr.GetErrors())
			}
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestIntegerExpression(t *testing.T) {
	tt := []struct {
		name       string
		sourceCode string
		want       *ast.Program
		wantErr    bool
	}{
		{
			name:       "one integer expression",
			sourceCode: "5;",
			want: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Expression: &ast.Integer{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			lxr := lexer.New(tc.sourceCode)
			psr := New(lxr)
			got := psr.ParseProgram()
			if tc.wantErr {
				assert.NotEmpty(t, psr.GetErrors())
			}
			if !tc.wantErr {
				assert.Empty(t, psr.GetErrors())
			}
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestPrefixExpression(t *testing.T) {
	tt := []struct {
		name       string
		sourceCode string
		want       *ast.Program
		wantErr    bool
	}{
		{
			name:       "one prefix expression",
			sourceCode: "-5",
			want: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{Type: token.MINUS, Literal: "-"},
						Expression: &ast.Prefix{
							Token:    token.Token{Type: token.MINUS, Literal: "-"},
							Operator: "-",
							Right: &ast.Integer{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "one prefix expression",
			sourceCode: "-15",
			want: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{Type: token.MINUS, Literal: "-"},
						Expression: &ast.Prefix{
							Token:    token.Token{Type: token.MINUS, Literal: "-"},
							Operator: "-",
							Right: &ast.Integer{
								Token: token.Token{Type: token.INT, Literal: "15"},
								Value: 15,
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			lxr := lexer.New(tc.sourceCode)
			psr := New(lxr)
			got := psr.ParseProgram()
			if tc.wantErr {
				assert.NotEmpty(t, psr.GetErrors())
			}
			if !tc.wantErr {
				assert.Empty(t, psr.GetErrors())
			}
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestInfixExpression(t *testing.T) {
	tt := []struct {
		name       string
		sourceCode string
		want       *ast.Program
		wantErr    bool
	}{
		{
			name:       "oen prefix expression",
			sourceCode: "5 + 5;",
			want: &ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Expression: &ast.Infix{
							Token:    token.Token{Type: token.PLUS, Literal: "+"},
							Operator: "+",
							Left: &ast.Integer{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
							Right: &ast.Integer{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			lxr := lexer.New(tc.sourceCode)
			psr := New(lxr)
			got := psr.ParseProgram()
			if tc.wantErr {
				assert.NotEmpty(t, psr.GetErrors())
			}
			if !tc.wantErr {
				assert.Empty(t, psr.GetErrors())
			}
			assert.EqualValues(t, tc.want, got)
		})
	}
}
