package golfcart

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func InjectRuntime(context *Context) {
	context.stackFrame.values["assert"] = NativeFunctionValue{name: "assert", Exec: golfcartAssert}
	context.stackFrame.values["log"] = NativeFunctionValue{name: "log", Exec: golfcartLog}
	context.stackFrame.values["type"] = NativeFunctionValue{name: "type", Exec: golfcartType}
	context.stackFrame.values["str"] = NativeFunctionValue{name: "str", Exec: golfcartStr}
	context.stackFrame.values["num"] = NativeFunctionValue{name: "num", Exec: golfcartNum}
}

type NativeFunctionValue struct {
	name string
	// TODO: improve error message (+ line number if pos?)
	Exec func([]Value) (Value, error)
}

func (nativeFunctionValue NativeFunctionValue) String() string {
	return nativeFunctionValue.name + " function"
}

func (nativeFunctionValue NativeFunctionValue) Equals(other Value) bool {
	return nativeFunctionValue.name == nativeFunctionValue.name
}

func golfcartAssert(args []Value) (Value, error) {
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
}

func golfcartLog(args []Value) (Value, error) {
	s := make([]string, len(args))
	for i := range args {
		s[i] = args[i].String()
	}
	fmt.Println(strings.Join(s, ", "))
	return NilValue{}, nil
}

func golfcartStr(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, errors.New("str() expects 1 argument of type num or bool")
	}
	value := args[0]
	switch value.(type) {
	case StringValue:
		return value, nil
	case NumberValue:
		return StringValue{val: value.String()}, nil
	case BoolValue:
		return StringValue{val: value.String()}, nil
	case FunctionValue:
		return StringValue{val: value.String()}, nil
	case NativeFunctionValue:
		return StringValue{val: value.String()}, nil
	}

	return nil, errors.New("str() expects 1 argument of type num or bool")
}

func golfcartNum(args []Value) (Value, error) {
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
}

func golfcartType(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, errors.New("type() expects 1 argument")
	}
	value := args[0]
	switch value.(type) {
	case StringValue:
		return StringValue{val: "string"}, nil
	case NumberValue:
		return StringValue{val: "number"}, nil
	case BoolValue:
		return StringValue{val: "bool"}, nil
	case FunctionValue, NativeFunctionValue:
		return StringValue{val: "function"}, nil
	case ListValue:
		return StringValue{val: "list"}, nil
	case DictValue:
		return StringValue{val: "dict"}, nil
	case NilValue:
		return StringValue{val: "nil"}, nil
	}

	panic("unimplemented golfcartType")
}
