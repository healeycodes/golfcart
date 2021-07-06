package golfcart

import (
	"errors"
	"fmt"
	"strings"
)

type Context struct {
	stackFrame StackFrame
}

func (context *Context) Init() {
	context.stackFrame = StackFrame{values: make(map[string]Value)}
}

type StackFrame struct {
	values map[string]Value
	parent *StackFrame
}

func (frame *StackFrame) String() string {
	s := ""
	for true {
		s += "["
		for key, value := range frame.values {
			s += key + ": " + value.String() + ", "
		}
		s += "] --> "
		if parent := frame.parent; parent != nil {
			frame = parent
		} else {
			break
		}
	}

	return s + "*"
}

func (frame *StackFrame) GetChild() *StackFrame {
	childFrame := StackFrame{parent: frame, values: make(map[string]Value)}
	return &childFrame
}

func (frame *StackFrame) Get(key Value) (Value, error) {
	for true {
		value, ok := frame.values[key.String()]
		if ok {
			return value, nil
		}
		if parent := frame.parent; parent != nil {
			frame = parent
		} else {
			break
		}
	}

	return nil, errors.New("cannot find value for '" + key.String() + "'")
}

func (frame *StackFrame) Set(key Value, value Value) {
	currentFrame := frame
	for true {
		_, ok := frame.values[key.String()]
		if ok {
			frame.values[key.String()] = value
			return
		}
		if parent := frame.parent; parent != nil {
			frame = parent
		} else {
			break
		}
	}

	currentFrame.values[key.String()] = value
}

type Value interface {
	String() string
	Equals(Value) bool
}

type NilValue struct{}

func (numberValue NilValue) String() string {
	return "nil"
}

func (numberValue NilValue) Equals(other Value) bool {
	if _, ok := other.(NilValue); ok {
		return true
	}
	return false
}

type IdentifierValue struct {
	val string
}

func (identifierValue IdentifierValue) String() string {
	return identifierValue.val
}

func (identifierValue IdentifierValue) Equals(other Value) bool {
	panic("identifiers shouldn't be compared")
}

func (identifierValue IdentifierValue) Eval(frame *StackFrame) (Value, error) {
	value, err := frame.Get(identifierValue)
	if err != nil {
		return nil, err
	}
	return value, nil
}

type NumberValue struct {
	val float64
}

func (numberValue NumberValue) String() string {
	return fmt.Sprintf("%f", numberValue.val)
}

func (numberValue NumberValue) Equals(other Value) bool {
	if other, ok := other.(NumberValue); ok {
		return numberValue.val == other.val
	}
	panic("unimplemented NumberValue Equals")
}

type StringValue struct {
	val string
}

func (stringValue StringValue) String() string {
	return stringValue.val
}

func (stringValue StringValue) Equals(other Value) bool {
	if other, ok := other.(StringValue); ok {
		return stringValue.val == other.val
	}
	panic("unimplemented StringValue Equals")
}

type BoolValue struct {
	val bool
}

func (boolValue BoolValue) String() string {
	return fmt.Sprintf("%t", boolValue.val)
}

func (boolValue BoolValue) Equals(other Value) bool {
	return boolValue == other
}

type FunctionValue struct {
	parameters  []string
	frame       *StackFrame
	expressions []*Expression
}

func (functionValue FunctionValue) String() string {
	return "FunctionValue"
}

func (functionValue FunctionValue) Equals(other Value) bool {
	// TODO: should function values be comparable?
	return false
}

func (functionValue FunctionValue) Exec(args []Value) (Value, error) {
	if len(args) != len(functionValue.parameters) {
		// TODO: improve error message (+ line number if pos?)
		return nil, errors.New("function called with incorrect number of arguments")
	}
	for i, parameter := range functionValue.parameters {
		functionValue.frame.Set(IdentifierValue{val: parameter}, args[i])
	}
	var result Value
	var err error
	result = NilValue{}
	for _, expression := range functionValue.expressions {
		result, err = expression.Eval(functionValue.frame)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// --

func (exprList ExpressionList) String() string {
	s := make([]string, 0)
	for _, expr := range exprList.Expressions {
		s = append(s, expr.String())
	}
	return strings.Join(s, ", ")
}

func (exprList ExpressionList) Equals(_ Value) bool {
	return false
}

func (exprList ExpressionList) Eval(context *Context) (Value, error) {
	results := make([]Value, 0)
	for _, node := range exprList.Expressions {
		result, err := node.Eval(&context.stackFrame)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if len(results) == 0 {
		return NilValue{}, nil
	}

	return results[len(results)-1], nil
}

func (expr Expression) String() string {
	return "expr"
}

func (expr Expression) Equals(other Value) bool {
	return expr == other
}

func (expr Expression) Eval(frame *StackFrame) (Value, error) {
	if expr.Assignment != nil {
		result, err := expr.Assignment.Eval(frame)
		if err != nil {
			return nil, err
		}
		if identifierValue, ok := result.(IdentifierValue); ok {
			value, err := frame.Get(identifierValue)
			if err != nil {
				return nil, err
			}
			return value, nil
		}
		return result, nil
	}
	return nil, errors.New("unimplemented Expression Eval")
}

func (functionLiteral FunctionLiteral) String() string {
	return "functionLiteral"
}

func (functionLiteral FunctionLiteral) Equals(other Value) bool {
	// TODO: should function literals be comparable?
	return false
}

func (functionLiteral FunctionLiteral) Eval(frame *StackFrame) (Value, error) {
	closureFrame := frame.GetChild()
	functionValue := FunctionValue{parameters: functionLiteral.Parameters, frame: closureFrame, expressions: functionLiteral.Body}
	return functionValue, nil
}

func (call Call) String() string {
	return "call"
}

func (call Call) Equals(other Value) bool {
	// TODO: should function literals be comparable?
	return false
}

func (call Call) Eval(frame *StackFrame) (Value, error) {
	var args []Value
	if parameters := call.Parameters; parameters != nil {
		args = make([]Value, len(*parameters))
		for i, parameter := range *parameters {
			result, err := parameter.Eval(frame)
			if err != nil {
				return nil, err
			}
			args[i] = result
		}
	}
	if ident := call.Ident; ident != nil {
		value, err := frame.Get(IdentifierValue{val: *ident})
		if err != nil {
			return nil, err
		}
		if functionValue, ok := value.(FunctionValue); ok {
			// TODO: pass the cursor location (call.Pos) for better function errors?
			result, err := functionValue.Exec(args)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		// If list
		// If obj
		// var access string
		// var computedAccess Value
		// if call.Access != nil {
		// 	access = IdentifierValue{val: *call.Access}.String()
		// }
		// if call.ComputedAccess != nil {
		// 	_computedAccess, err := call.ComputedAccess.Eval(frame)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	computedAccess = _computedAccess
		// }
	}
	if call.SubExpression != nil {
		result, err := call.SubExpression.Eval(frame)
		if err != nil {
			return nil, err
		}
		if functionValue, ok := result.(FunctionValue); ok {
			result, err := functionValue.Exec(args)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	}
	panic("unimplemented Call Eval")
}

func (assignment Assignment) String() string {
	return "assignment"
}

func (assignment Assignment) Equals(other Value) bool {
	return assignment == other
}

func (assignment Assignment) Eval(frame *StackFrame) (Value, error) {
	left, err := assignment.LogicAnd.Eval(frame)
	if err != nil {
		return nil, err
	}
	if assignment.Op == "" {
		return left, nil
	}
	right, err := assignment.Next.Eval(frame)
	if err != nil {
		return nil, err
	}
	if assignment.Op == "=" {
		frame.Set(left, right)
		return NilValue{}, nil
	}
	panic("unimplemented Assignment Eval")
}

func (logicAnd LogicAnd) String() string {
	return "equality"
}

func (logicAnd LogicAnd) Equals(other Value) bool {
	return logicAnd == other
}

func (logicAnd LogicAnd) Eval(frame *StackFrame) (Value, error) {
	left, err := logicAnd.LogicOr.Eval(frame)
	if err != nil {
		return nil, err
	}
	if logicAnd.Op == "" {
		return left, nil
	}
	panic("unimplemented LogicAnd Eval")
}

func (logicOr LogicOr) String() string {
	return "equality"
}

func (logicOr LogicOr) Equals(other Value) bool {
	return logicOr == other
}

func (logicOr LogicOr) Eval(frame *StackFrame) (Value, error) {
	left, err := logicOr.Equality.Eval(frame)
	if err != nil {
		return nil, err
	}
	if logicOr.Op == "" {
		return left, nil
	}
	return nil, errors.New("unimplemented LogicOr Eval")
}

func (equality Equality) String() string {
	return "equality"
}

func (equality Equality) Equals(other Value) bool {
	return equality == other
}

func (equality Equality) Eval(frame *StackFrame) (Value, error) {
	left, err := equality.Comparison.Eval(frame)
	if err != nil {
		return nil, err
	}
	if equality.Op == "" {
		return left, nil
	}
	right, err := equality.Next.Eval(frame)
	if err != nil {
		return nil, err
	}
	if equality.Op == "==" {
		return BoolValue{val: left.Equals(right)}, nil
	} else if equality.Op == "!=" {
		return BoolValue{val: !left.Equals(right)}, nil
	}
	return nil, errors.New("unimplemented Equality Eval")
}

func (comparison Comparison) String() string {
	return "comparison"
}

func (comparison Comparison) Equals(other Value) bool {
	return comparison == other
}

func (comparison Comparison) Eval(frame *StackFrame) (Value, error) {
	left, err := comparison.Addition.Eval(frame)
	if err != nil {
		return nil, err
	}
	if comparison.Op == "" {
		return left, nil
	}
	return nil, errors.New("unimplemented Comparison Eval")
}

func (addition Addition) String() string {
	return "addition"
}

func (addition Addition) Equals(other Value) bool {
	return addition == other
}

func (addition Addition) Eval(frame *StackFrame) (Value, error) {
	left, err := addition.Multiplication.Eval(frame)
	if err != nil {
		return nil, err
	}
	if addition.Op == "" {
		return left, nil
	}
	right, err := addition.Next.Eval(frame)
	if err != nil {
		return nil, err
	}

	if leftId, okLeft := left.(IdentifierValue); okLeft {
		left, err = frame.Get(leftId)
		if err != nil {
			return nil, err
		}
	}

	if rightId, okRight := right.(IdentifierValue); okRight {
		right, err = frame.Get(rightId)
		if err != nil {
			return nil, err
		}
	}

	leftStr, okLeft := left.(StringValue)
	rightStr, okRight := right.(StringValue)
	if addition.Op == "+" && (okLeft && !okRight || okRight && !okLeft) {
		return nil, errors.New(addition.Multiplication.Pos.String() + " '+' only supported between strings")
	} else if addition.Op == "+" && (okLeft || okRight) {
		return StringValue{val: leftStr.val + rightStr.val}, nil
	}

	leftNum, okLeft := left.(NumberValue)
	if !okLeft {
		return nil, errors.New(addition.Multiplication.Pos.String() + " '+' and '-' only supported between numbers")
	}
	rightNum, okRight := right.(NumberValue)
	if !okRight {
		return nil, errors.New(addition.Next.Pos.String() + " '+' and '-' only supported between numbers")
	}
	if addition.Op == "+" {
		return NumberValue{val: leftNum.val + rightNum.val}, nil
	}
	if addition.Op == "-" {
		return NumberValue{val: leftNum.val - rightNum.val}, nil
	}

	panic("unimplemented Addition Eval")
}

func (multiplication Multiplication) String() string {
	return "multiplication"
}

func (multiplication Multiplication) Equals(other Value) bool {
	return multiplication == other
}

func (multiplication Multiplication) Eval(frame *StackFrame) (Value, error) {
	left, err := multiplication.Unary.Eval(frame)
	if err != nil {
		return nil, err
	}
	if multiplication.Op == "" {
		return left, nil
	}
	right, err := multiplication.Next.Eval(frame)
	if err != nil {
		return nil, err
	}

	if leftId, okLeft := left.(IdentifierValue); okLeft {
		left, err = frame.Get(leftId)
		if err != nil {
			return nil, err
		}
	}

	if rightId, okRight := right.(IdentifierValue); okRight {
		right, err = frame.Get(rightId)
		if err != nil {
			return nil, err
		}
	}

	leftNum, okLeft := left.(NumberValue)
	if !okLeft {
		return nil, errors.New(multiplication.Unary.Pos.String() + " '+' and '-' only supported between numbers")
	}
	rightNum, okRight := right.(NumberValue)
	if !okRight {
		return nil, errors.New(multiplication.Next.Pos.String() + " '+' and '-' only supported between numbers")
	}
	if multiplication.Op == "*" {
		return NumberValue{val: leftNum.val * rightNum.val}, nil
	}
	if multiplication.Op == "/" {
		return NumberValue{val: leftNum.val / rightNum.val}, nil
	}

	return nil, errors.New("unimplemented Multiplication Eval")
}

func (unary Unary) String() string {
	return "unary"
}

func (unary Unary) Equals(other Value) bool {
	return unary == other
}

func (unary Unary) Eval(frame *StackFrame) (Value, error) {
	if unary.Op == "!" {
		_unary, err := unary.Unary.Eval(frame)
		if err != nil {
			return nil, err
		}
		if boolValue, ok := _unary.(BoolValue); ok {
			return BoolValue{val: !boolValue.val}, nil
		}
		return nil, errors.New(unary.Pos.String() + " expected bool after '!'")
	}
	if unary.Op == "-" {
		_unary, err := unary.Unary.Eval(frame)
		if err != nil {
			return nil, err
		}
		if numberValue, ok := _unary.(NumberValue); ok {
			return NumberValue{val: -numberValue.val}, nil
		}
		return nil, errors.New(unary.Pos.String() + " expected number after '-'")
	}

	if unary.Primary != nil {
		return unary.Primary.Eval(frame)
	}

	return nil, errors.New("unimplemented Unary Eval")
}

func (primary Primary) String() string {
	return "primary"
}

func (primary Primary) Equals(other Value) bool {
	panic("unimplemented Primary Equals")
}

func (primary Primary) Eval(frame *StackFrame) (Value, error) {
	if functionLiteral := primary.FunctionLiteral; functionLiteral != nil {
		return functionLiteral.Eval(frame)
	}
	if call := primary.Call; call != nil {
		return call.Eval(frame)
	}
	if ident := primary.Ident; ident != "" {
		identifierValue := IdentifierValue{val: ident}
		return identifierValue, nil
	}
	if primary.Number != nil {
		return NumberValue{val: *primary.Number}, nil
	}
	if primary.Str != nil {
		// TODO: Parse strings without including quote `"` marks
		return StringValue{val: (*primary.Str)[1 : len((*primary.Str))-1]}, nil
	}
	if primary.Bool != nil {
		return BoolValue{val: *primary.Bool}, nil
	}
	if primary.Nil != nil {
		return NilValue{}, nil
	}

	panic("unimplemented Primary Eval")
}
