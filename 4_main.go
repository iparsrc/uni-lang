package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	scope := NewScope(nil)
	if len(os.Args) < 2 {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Uni Version 0.1.0")
		for {
			fmt.Print(">> ")
			scanner.Scan()
			sourceCode := scanner.Text()
			lexer := NewLexer(sourceCode)
			parser := NewParser(lexer)
			evaluator := NewEvaluator(parser)
			if evaluated := evaluator.Eval(scope); evaluated != nil {
				fmt.Println(evaluated)
			} else {
				fmt.Println("")
			}
		}
	}
	sourceCode, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	lexer := NewLexer(string(sourceCode))
	parser := NewParser(lexer)
	evaluator := NewEvaluator(parser)
	evaluator.Eval(scope)
}
