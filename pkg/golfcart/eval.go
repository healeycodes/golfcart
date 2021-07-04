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
	parent map[string]Value
}

func (frame *StackFrame) Get(key Value) (Value, error) {
	value, ok := frame.values[key.String()]
	if ok {
		return value, nil
	}
	return nil, errors.New("cannot find value for " + key.String())
}

func (frame *StackFrame) Set(key Value, value Value) {
	frame.values[key.String()] = value
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

type VariableValue struct {
	val string
}

func (variableValue VariableValue) String() string {
	return variableValue.val
}

func (variableValue VariableValue) Equals(other Value) bool {
	panic("unimplemented VariableValue Equals")
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

	return results[0], nil
}

func (expr Expression) String() string {
	return "expr"
}

func (expr Expression) Equals(other Value) bool {
	return expr == other
}

func (expr Expression) Eval(frame *StackFrame) (Value, error) {
	if expr.Assignment != nil {
		return expr.Assignment.Eval(frame)
	}
	return nil, errors.New("Unimplemented")
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
	if assignment.Op != "" {
		return left, nil
	}
	right, err := assignment.LogicAnd.Eval(frame)
	if err != nil {
		return nil, err
	}
	frame.Set(left, right)
	return right, nil
}

func (logicAnd LogicAnd) String() string {
	return "equality"
}

func (logicAnd LogicAnd) Equals(other Value) bool {
	return logicAnd == other
}

func (logicAnd LogicAnd) Eval(frame *StackFrame) (Value, error) {
	return nil, errors.New("unimplemented LogicAnd Eval")
}

func (logicOr LogicOr) String() string {
	return "equality"
}

func (logicOr LogicOr) Equals(other Value) bool {
	return logicOr == other
}

func (logicOr LogicOr) Eval(frame *StackFrame) (Value, error) {
	return nil, errors.New("unimplemented LogicOr Eval")
}

func (equality Equality) String() string {
	return "equality"
}

func (equality Equality) Equals(other Value) bool {
	return equality == other
}

func (equality Equality) Eval(frame *StackFrame) (Value, error) {
	return nil, errors.New("unimplemented Equality Eval")
}

func (comparison Comparison) String() string {
	return "comparison"
}

func (comparison Comparison) Equals(other Value) bool {
	return comparison == other
}

func (comparison Comparison) Eval(frame *StackFrame) (Value, error) {
	return nil, errors.New("unimplemented Comparison Eval")
}

func (addition Addition) String() string {
	return "addition"
}

func (addition Addition) Equals(other Value) bool {
	return addition == other
}

func (addition Addition) Eval(frame *StackFrame) (Value, error) {
	return nil, errors.New("unimplemented Addition Eval")
}

func (multiplication Multiplication) String() string {
	return "multiplication"
}

func (multiplication Multiplication) Equals(other Value) bool {
	return multiplication == other
}

func (multiplication Multiplication) Eval(frame *StackFrame) (Value, error) {
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
	if primary.Ident != nil {
		return VariableValue{val: *primary.Ident}, nil
	}
	if primary.Number != nil {
		return NumberValue{val: *primary.Number}, nil
	}
	if primary.String != nil {
		return StringValue{val: *primary.Str}, nil
	}
	if primary.Bool != nil {
		return BoolValue{val: *primary.Bool}, nil
	}
	if primary.Nil != nil {
		return NilValue{}, nil
	}

	panic("unimplemented Primary Eval")
}
