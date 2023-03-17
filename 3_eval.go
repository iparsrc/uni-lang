// Evaluating is the process that defines how the programming language being interpreted works.
// The statements are executed in the source language, which in this case it is Golang.

package main

import (
	"fmt"
)

// ************
// ** Scope **
// ************

type Scope struct {
	variables map[string]any
	functions map[string]any
	parent    *Scope
}

func NewScope(scope *Scope) *Scope {
	return &Scope{
		variables: make(map[string]any),
		functions: make(map[string]any),
		parent:    scope,
	}
}

func (s *Scope) GetVariable(identifier Identifier) (any, bool) {
	variable, ok := s.variables[identifier.Token.Value]
	if !ok && s.parent != nil {
		variable, ok = s.parent.GetVariable(identifier)
	}
	return variable, ok
}

func (s *Scope) SetVariable(identifier Identifier, value any) {
	s.variables[identifier.Token.Value] = value
}

func (s *Scope) GetFunction(identifier Identifier) (any, bool) {
	function, ok := s.functions[identifier.Token.Value]
	if !ok && s.parent != nil {
		function, ok = s.parent.GetFunction(identifier)
	}
	return function, ok
}

func (s *Scope) SetFunction(identifier Identifier, value any) {
	s.functions[identifier.Token.Value] = value
}

func (s *Scope) GetParent() *Scope {
	return s.parent
}

// ***************
// ** Evaluator **
// ***************

type Evaluator struct {
	parser *Parser
}

func NewEvaluator(parser *Parser) *Evaluator {
	return &Evaluator{
		parser: parser,
	}
}

func (e *Evaluator) Eval(scope *Scope) any {
	var value any
	statements := e.parser.Parse()
	for statement := range statements {
		if statement == nil {
			break
		}
		value = evalStatement(statement, scope)
	}
	return value
}

// ****************
// ** Statements **
// ****************

func evalStatement(statement Statement, scope *Scope) any {
	switch typedStatement := statement.(type) {
	case Variable:
		return evalVariable(typedStatement, scope)
	case If:
		return evalIf(typedStatement, scope)
	case While:
		return evalWhile(typedStatement, scope)
	case For:
		return evalFor(typedStatement, scope)
	case Function:
		return evalFunction(typedStatement, scope)
	case Return:
		return evalReturn(typedStatement, scope)
	case Block:
		return evalBlock(typedStatement, NewScope(scope))
	default:
		return evalExpression(typedStatement, scope)
	}
}

func evalVariable(in Variable, scope *Scope) any {
	if in.IsNew {
		scope.SetVariable(in.Name, evalExpression(in.Value, scope))
		return nil
	}
	for {
		if _, ok := scope.GetVariable(in.Name); ok {
			scope.SetVariable(in.Name, evalExpression(in.Value, scope))
		}
		scope = scope.GetParent()
		if scope == nil {
			break
		}
	}
	return nil
}

func evalIf(in If, scope *Scope) any {
	if evalExpression(in.Condition, scope).(bool) {
		return evalBlock(in.Consequence, NewScope(scope))
	}
	if in.Alternative != nil {
		return evalBlock(*in.Alternative, NewScope(scope))
	}
	return nil
}

func evalWhile(in While, scope *Scope) any {
	for evalExpression(in.Condition, scope).(bool) {
		newScope := NewScope(scope)
		if result := evalBlock(in.Consequence, newScope); result != nil {
			return result
		}
	}
	return nil
}

func evalFor(in For, scope *Scope) any {
	switch subject := evalExpression(in.Condition, scope).(type) {
	case string:
		for key, value := range subject {
			newScope := NewScope(scope)
			newScope.SetVariable(in.Key, key)
			newScope.SetVariable(in.Value, string(value))
			if result := evalBlock(in.Consequence, newScope); result != nil {
				return result
			}
		}
	case []any:
		for key, value := range subject {
			newScope := NewScope(scope)
			newScope.SetVariable(in.Key, key)
			newScope.SetVariable(in.Value, value)
			if result := evalBlock(in.Consequence, newScope); result != nil {
				return result
			}
		}
	case map[any]any:
		for key, value := range subject {
			newScope := NewScope(scope)
			newScope.SetVariable(in.Key, key)
			newScope.SetVariable(in.Value, value)
			if result := evalBlock(in.Consequence, newScope); result != nil {
				return result
			}
		}
	}
	return nil
}

func evalFunction(in Function, scope *Scope) any {
	scope.SetFunction(in.Name, in)
	return nil
}

func evalReturn(in Return, scope *Scope) any {
	return evalExpression(in.Value, scope)
}

func evalBlock(in Block, scope *Scope) any {
	for _, statement := range in.Statements {
		if statement == nil {
			continue
		}
		if result := evalStatement(statement, scope); result != nil {
			return result
		}
	}
	return nil
}

// *****************
// ** Expressions **
// *****************

func evalExpression(expression Expression, scope *Scope) any {
	switch typedExpression := expression.(type) {
	case Boolean:
		return evalBoolean(typedExpression, scope)
	case Integer:
		return evalInteger(typedExpression, scope)
	case Float:
		return evalFloat(typedExpression, scope)
	case String:
		return evalString(typedExpression, scope)
	case Array:
		return evalArray(typedExpression, scope)
	case Map:
		return evalMap(typedExpression, scope)
	case Index:
		return evalIndex(typedExpression, scope)
	case Call:
		return evalCall(typedExpression, scope)
	case Identifier:
		return evalIdentifier(typedExpression, scope)
	case UnaryOperation:
		return evalUnaryOperation(typedExpression, scope)
	case BinaryOperation:
		return evalBinaryOperation(typedExpression, scope)
	case Len:
		return evalLen(typedExpression, scope)
	case Print:
		return evalPrint(typedExpression, scope)
	default:
		return nil
	}
}

func evalBoolean(in Boolean, _ *Scope) any {
	return in.Value
}

func evalInteger(in Integer, _ *Scope) any {
	return in.Value
}

func evalFloat(in Float, _ *Scope) any {
	return in.Value
}

func evalString(in String, _ *Scope) any {
	return in.Value
}

func evalArray(in Array, scope *Scope) any {
	a := make([]any, len(in.Items))
	for key, value := range in.Items {
		a[key] = evalExpression(value, scope)
	}
	return a
}

func evalMap(in Map, scope *Scope) any {
	m := make(map[any]any, len(in.Items))
	for key, value := range in.Items {
		m[evalExpression(key, scope)] = evalExpression(value, scope)
	}
	return m
}

func evalIndex(in Index, scope *Scope) any {
	switch subject := evalExpression(in.Subject, scope).(type) {
	case []any:
		return subject[int(evalExpression(in.Index, scope).(int64))]
	case map[any]any:
		return subject[evalExpression(in.Index, scope)]
	default:
		return nil
	}
}

func evalCall(in Call, scope *Scope) any {
	untypedFunction, ok := scope.GetFunction(in.Identifier)
	if !ok {
		return nil
	}
	function := untypedFunction.(Function)
	if len(function.Parameters) != len(in.Arguments) {
		return nil
	}
	newScope := NewScope(scope)
	for i, argument := range in.Arguments {
		newScope.SetVariable(function.Parameters[i], evalExpression(argument, scope))
	}
	return evalBlock(function.Body, newScope)
}

func evalIdentifier(identifier Identifier, scope *Scope) any {
	if identifier.IsFunctionCall {
		function, _ := scope.GetFunction(identifier)
		return function
	}
	variable, _ := scope.GetVariable(identifier)
	return variable
}

func evalUnaryOperation(in UnaryOperation, scope *Scope) any {
	switch t := evalExpression(in.Expression, scope).(type) {
	case bool:
		if in.Token.Type == NOT {
			return !t
		}
	case int64:
		if in.Token.Type == PLUS {
			return t
		}
		if in.Token.Type == MINUS {
			return -1 * t
		}
	case float64:
		if in.Token.Type == PLUS {
			return t
		}
		if in.Token.Type == MINUS {
			return -1 * t
		}
	}
	return nil
}

func evalBinaryOperation(in BinaryOperation, scope *Scope) any {
	switch left := evalExpression(in.Left, scope).(type) {
	case bool:
		switch right := evalExpression(in.Right, scope).(type) {
		case bool:
			return evalBinaryOperationBoolBool(left, right, in.Token)
		default:
			return nil
		}
	case int64:
		switch right := evalExpression(in.Right, scope).(type) {
		case int64:
			return evalBinaryOperationIntInt(left, right, in.Token)
		case float64:
			return evalBinaryOperationIntFloat(left, right, in.Token)
		default:
			return nil
		}
	case float64:
		switch right := evalExpression(in.Right, scope).(type) {
		case int64:
			return evalBinaryOperationFloatInt(left, right, in.Token)
		case float64:
			return evalBinaryOperationFloatFloat(left, right, in.Token)
		default:
			return nil
		}
	case string:
		switch right := evalExpression(in.Right, scope).(type) {
		case string:
			return evalBinaryOperationStringString(left, right, in.Token)
		default:
			return nil
		}
	default:
		return nil
	}
}

func evalBinaryOperationBoolBool(left bool, right bool, operator Token) any {
	switch operator.Type {
	case EQ:
		return left == right
	case NEQ:
		return left != right
	case OR:
		return left || right
	case AND:
		return left && right
	default:
		return nil
	}
}

func evalBinaryOperationIntInt(left int64, right int64, operator Token) any {
	switch operator.Type {
	case LT:
		return left < right
	case GT:
		return left > right
	case LEQ:
		return left <= right
	case GEQ:
		return left >= right
	case EQ:
		return left == right
	case NEQ:
		return left != right
	case PLUS:
		return left + right
	case MINUS:
		return left - right
	case ASTERISK:
		return left * right
	case SLASH:
		return left / right
	default:
		return nil
	}
}

func evalBinaryOperationIntFloat(left int64, right float64, operator Token) any {
	switch operator.Type {
	case LT:
		return float64(left) < right
	case GT:
		return float64(left) > right
	case LEQ:
		return float64(left) <= right
	case GEQ:
		return float64(left) >= right
	case EQ:
		return float64(left) == right
	case NEQ:
		return float64(left) != right
	case PLUS:
		return float64(left) + right
	case MINUS:
		return float64(left) - right
	case ASTERISK:
		return float64(left) * right
	case SLASH:
		return float64(left) / right
	default:
		return nil
	}
}

func evalBinaryOperationFloatInt(left float64, right int64, operator Token) any {
	switch operator.Type {
	case LT:
		return left < float64(right)
	case GT:
		return left > float64(right)
	case LEQ:
		return left <= float64(right)
	case GEQ:
		return left >= float64(right)
	case EQ:
		return left == float64(right)
	case NEQ:
		return left != float64(right)
	case PLUS:
		return left + float64(right)
	case MINUS:
		return left - float64(right)
	case ASTERISK:
		return left * float64(right)
	case SLASH:
		return left / float64(right)
	default:
		return nil
	}
}

func evalBinaryOperationFloatFloat(left float64, right float64, operator Token) any {
	switch operator.Type {
	case LT:
		return left < right
	case GT:
		return left > right
	case LEQ:
		return left <= right
	case GEQ:
		return left >= right
	case EQ:
		return left == right
	case NEQ:
		return left != right
	case PLUS:
		return left + right
	case MINUS:
		return left - right
	case ASTERISK:
		return left * right
	case SLASH:
		return left / right
	default:
		return nil
	}
}

func evalBinaryOperationStringString(left string, right string, operator Token) any {
	switch operator.Type {
	case PLUS:
		return left + right
	default:
		return nil
	}
}

func evalLen(in Len, scope *Scope) any {
	switch typedSubject := evalExpression(in.Subject, scope).(type) {
	case string:
		return len(typedSubject)
	case []any:
		return len(typedSubject)
	case map[any]any:
		return len(typedSubject)
	default:
		return nil
	}
}

func evalPrint(in Print, scope *Scope) any {
	var args []any
	for _, arg := range in.Args {
		args = append(args, evalExpression(arg, scope))
	}
	if in.IsNewLine {
		fmt.Println(args...)
	} else {
		fmt.Print(args...)
	}
	return nil
}
