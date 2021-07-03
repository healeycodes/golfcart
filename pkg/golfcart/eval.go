package golfcart

import (
	"errors"
	"fmt"
	"strings"
)

type StackFrame struct {
	Values map[string]Value
}

type Value interface {
	String() string
	Equals(Value) bool
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
	return nil, errors.New("Unimplemented")
}

func (arithmetic Arithmetic) String() string {
	return "arithmetic"
}

func (arithmetic Arithmetic) Equals(other Value) bool {
	return arithmetic == other
}

func (arithmetic Arithmetic) Eval() (Value, error) {
	if arithmetic.Op == "" {
		result, err := arithmetic.Unary.Eval()
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New("Unimplemented")
}

func (unary Unary) String() string {
	return "unary"
}

func (unary Unary) Equals(other Value) bool {
	return unary == other
}

func (unary Unary) Eval() (Value, error) {
	if unary.Primary != nil && unary.Primary.Number != nil {
		return NumberValue{val: *unary.Primary.Number}, nil
	}
	return nil, errors.New("Unimplemented")
}

func (expr Expression) Eval(frame *StackFrame) (Value, error) {
	if expr.Binary != nil {
		return expr.Binary.Eval()
	}
	return nil, errors.New("Unimplemented")
}

func (exprList ExpressionList) Eval(stackframe *StackFrame) (Value, error) {
	for _, node := range exprList.Expressions {
		if node.Binary != nil {
			result, err := node.Binary.Eval()
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	}
	return nil, errors.New("Unimplemented")
}
