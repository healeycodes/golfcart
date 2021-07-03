package golfcart

import (
	"errors"
	"fmt"
	"strings"
)

type StackFrame struct {
	Values map[string]Value
	Parent map[string]Value
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

	panic("Unimplemented")
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

func (expr Expression) String() string {
	return "expr"
}

func (expr Expression) Equals(other Value) bool {
	return expr == other
}

func (binary Binary) String() string {
	return "binary"
}

func (binary Binary) Equals(other Value) bool {
	return binary == other
}

func (binary Binary) Eval() (Value, error) {
	left, err := binary.Arithmetic.Eval()
	if err != nil {
		return nil, err
	}

	if binary.Next == nil {
		return left, nil
	}

	right, err := binary.Next.Eval()
	if err != nil {
		return nil, err
	}

	if binary.Op == "==" {
		return BoolValue{val: left.Equals(right)}, nil
	}
	return nil, errors.New("unimplemented Binary Eval")
}

func (arithmetic Arithmetic) String() string {
	return "arithmetic"
}

func (arithmetic Arithmetic) Equals(other Value) bool {
	return arithmetic == other
}

func (arithmetic Arithmetic) Eval() (Value, error) {
	left, err := arithmetic.Unary.Eval()
	if err != nil {
		return nil, err
	}
	if arithmetic.Op == "" {
		return left, nil
	}

	right, err := arithmetic.Next.Eval()
	if err != nil {
		return nil, err
	}
	if arithmetic.Op == "+" {
		left, okLeft := left.(NumberValue)
		right, okRight := right.(NumberValue)
		if okLeft && okRight {
			return NumberValue{val: left.val + right.val}, nil
		}
		return nil, errors.New(arithmetic.EndPos.String() + " addition only supported between numbers")
	}

	return nil, errors.New("unimplemented Arithmetic Eval")
}

func (unary Unary) String() string {
	return "unary"
}

func (unary Unary) Equals(other Value) bool {
	return unary == other
}

func (unary Unary) Eval() (Value, error) {
	if unary.Op == "!" {
		if unary.Unary.Primary.Bool != nil {
			return BoolValue{val: !*unary.Unary.Primary.Bool}, nil
		}
		return nil, errors.New("Expected bool at " + unary.Unary.EndPos.String())
	}
	if unary.Op == "-" {
		if unary.Unary.Primary.Number != nil {
			return NumberValue{val: -*unary.Unary.Primary.Number}, nil
		}
		return nil, errors.New("Expected number at " + unary.EndPos.String())
	}

	if unary.Primary != nil && unary.Primary.Number != nil {
		return NumberValue{val: *unary.Primary.Number}, nil
	}
	if unary.Primary != nil && unary.Primary.Bool != nil {
		return BoolValue{val: *unary.Primary.Bool}, nil
	}

	return nil, errors.New("unimplemented Unary Eval")
}

func (expr Expression) Eval(frame *StackFrame) (Value, error) {
	if expr.Binary != nil {
		return expr.Binary.Eval()
	}
	return nil, errors.New("Unimplemented")
}

func (exprList ExpressionList) Eval(stackframe *StackFrame) (Value, error) {
	results := make([]Value, 0)
	for _, node := range exprList.Expressions {
		if node.Binary != nil {
			result, err := node.Binary.Eval()
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}

	if len(results) == 0 {
		return NilValue{}, nil
	}

	return results[0], nil
}
