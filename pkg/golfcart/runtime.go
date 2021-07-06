package golfcart

import "errors"

type NativeFunctionValue struct {
	parameters  []string
	frame       *StackFrame
	expressions []*Expression
}

func (nativeFunctionValue NativeFunctionValue) String() string {
	return "NativeFunctionValue"
}

func (nativeFunctionValue NativeFunctionValue) Equals(other Value) bool {
	// TODO: should function values be comparable?
	return false
}

func (nativeFunctionValue NativeFunctionValue) Exec(args []Value) (Value, error) {
	if len(args) != len(nativeFunctionValue.parameters) {
		// TODO: improve error message (+ line number if pos?)
		return nil, errors.New("function called with incorrect number of arguments")
	}
	for i, parameter := range nativeFunctionValue.parameters {
		nativeFunctionValue.frame.Set(IdentifierValue{val: parameter}, args[i])
	}
	var result Value
	var err error
	result = NilValue{}
	for _, expression := range nativeFunctionValue.expressions {
		result, err = expression.Eval(nativeFunctionValue.frame)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
