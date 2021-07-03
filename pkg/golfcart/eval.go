package golfcart

import (
	"errors"
	"strings"
)

// TODO: to_json()?
// out, err := json.Marshal(expr)
// if err != nil {
// 	panic(err)
// }

// return string(out)

type StackFrame struct {
	Values map[string]Value
}

type Value interface {
	String() string
	Equals(Value) bool
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

func (binary Binary) Eval(frame *StackFrame) (Value, error) {
	if binary.Op == "==" {
		result, err := binary.Arithmetic.Eval(frame).Equals(binary.Next.Eval(frame))
		if err != nil {
			return nil, err
		}
		return result
	}
	return nil, errors.New("Unimplemented")
}

func (arithmetic Arithmetic) String() string {
	return "arithmetic"
}

func (arithmetic Arithmetic) Equals(other Value) bool {
	return arithmetic == other
}

func (arithmetic Arithmetic) Eval(frame *StackFrame) (Value, error) {
	if arithmetic.Op == "" {
		result, err := arithmetic.Unary.Eval()
		if err != nil {
			return nil, err
		}
		return result
	}
	return nil, errors.New("Unimplemented")
}

func (binary Binary) Eval(frame *StackFrame) (Value, error) {
	if binary.Op == "==" {
		// TODO: split left/right
		result, err := binary.Arithmetic.Eval(frame).Equals(binary.Next.Eval(frame))
		if err != nil {
			return nil, err
		}
		return result
	}
	return nil, errors.New("Unimplemented")
}

func (expr Expression) Eval(frame *StackFrame) (Value, error) {
	if expr.Binary != nil {
		return expr.Binary.Eval(frame)
	}
	return nil, errors.New("Unimplemented")
}

func (exprList ExpressionList) Eval(stackframe *StackFrame) {
	for _, node := range exprList.Expressions {
		println(node.String())
	}
}
