package golfcart

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func InjectRuntimeFunctions(context *Context) {
	context.stackFrame.values["assert"] = NativeFunctionValue{name: "assert"}
	context.stackFrame.values["log"] = NativeFunctionValue{name: "log"}
	context.stackFrame.values["str"] = NativeFunctionValue{name: "str"}
	context.stackFrame.values["num"] = NativeFunctionValue{name: "num"}
}

type NativeFunctionValue struct {
	name        string
	parameters  []string
	frame       *StackFrame
	expressions []*Expression
}

func (nativeFunctionValue NativeFunctionValue) String() string {
	return "NativeFunctionValue"
}

func (nativeFunctionValue NativeFunctionValue) Equals(other Value) bool {
	// TODO: should native function values be comparable?
	return false
}

func (nativeFunctionValue NativeFunctionValue) Exec(args []Value) (Value, error) {
	switch nativeFunctionValue.name {
	case "assert":
		if len(args) != 1 {
			return nil, errors.New("assert() expects 1 argument of type bool")
		}
		boolValue, okBool := args[0].(BoolValue)
		if !okBool {
			return nil, errors.New("assert() expects 1 argument of type bool")
		}
		if boolValue.val == false {
			return nil, errors.New("assert failed!")
		}
		return NilValue{}, nil
	case "str":
		if len(args) != 1 {
			return nil, errors.New("str() expects 1 argument of type num or bool")
		}
		value := args[0]
		if strValue, okStr := value.(StringValue); okStr {
			return strValue, nil
		}
		_, okNum := value.(NumberValue)
		_, okBool := value.(BoolValue)
		if !okNum && !okBool {
			return nil, errors.New("str() expects 1 argument of type num or bool")
		}
		return StringValue{val: value.String()}, nil
	case "num":
		if len(args) != 1 {
			return nil, errors.New("num() expects 1 argument")
		}
		value := args[0]
		if numValue, okNum := value.(NumberValue); okNum {
			return numValue, nil
		}
		strValue, okStr := value.(StringValue)
		if !okStr {
			return nil, errors.New("num() expects 1 argument of type str")
		}
		f, err := strconv.ParseFloat(strValue.val, 64)
		if err != nil {
			return nil, errors.New("num() couldn't convert " + strValue.val + " to num")
		}
		return NumberValue{val: f}, nil
	case "log":
		s := make([]string, len(args))
		for i := range args {
			s[i] = args[i].String()
		}
		fmt.Println(strings.Join(s, ", "))
		return NilValue{}, nil
	}
	// if len(args) != len(nativeFunctionValue.parameters) {
	// 	// TODO: improve error message (+ line number if pos?)
	// 	return nil, errors.New("function called with incorrect number of arguments")
	// }
	// for i, parameter := range nativeFunctionValue.parameters {
	// 	nativeFunctionValue.frame.Set(IdentifierValue{val: parameter}, args[i])
	// }
	// var result Value
	// var err error
	// result = NilValue{}
	// for _, expression := range nativeFunctionValue.expressions {
	// 	result, err = expression.Eval(nativeFunctionValue.frame)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// return result, nil
	panic("unimplemented NativeFunctionValue Exec")
}
