package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	tt := []struct {
		name string
		in   string
		want any
	}{
		{
			name: "",
			in: `fn fibonacci(n) {
					if n == 1 {
						return 0
					}
					if n == 2 {
						return 1
					}
					return fibonacci(n-1) + fibonacci(n-2)
				}	
				fibonacci(3)		
			`,
			want: int64(1),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			lexer := NewLexer(tc.in)
			parser := NewParser(lexer)
			evaluator := NewEvaluator(parser)
			got := evaluator.Eval(NewScope(nil))
			assert.Equal(t, got, tc.want)
		})
	}
}
