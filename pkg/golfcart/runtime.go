package golfcart

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

func setNativeFunc(key string, nativeFunc Value, frame *StackFrame) {
	frame.entries[key] = nativeFunc
}

func InjectRuntime(context *Context) {
	setNativeFunc("assert", NativeFunctionValue{name: "assert", Exec: golfcartAssert}, &context.stackFrame)
	setNativeFunc("in", NativeFunctionValue{name: "in", Exec: golfcartIn}, &context.stackFrame)
	setNativeFunc("log", NativeFunctionValue{name: "log", Exec: golfcartLog}, &context.stackFrame)
	setNativeFunc("type", NativeFunctionValue{name: "type", Exec: golfcartType}, &context.stackFrame)
	setNativeFunc("str", NativeFunctionValue{name: "str", Exec: golfcartStr}, &context.stackFrame)
	setNativeFunc("num", NativeFunctionValue{name: "num", Exec: golfcartNum}, &context.stackFrame)
	setNativeFunc("len", NativeFunctionValue{name: "len", Exec: golfcartLen}, &context.stackFrame)
	setNativeFunc("keys", NativeFunctionValue{name: "keys", Exec: golfcartKeys}, &context.stackFrame)
	setNativeFunc("values", NativeFunctionValue{name: "values", Exec: golfcartValues}, &context.stackFrame)
}

type NativeFunctionValue struct {
	Pos  lexer.Position
	name string
	Exec func([]Value) (Value, error)
}

func (nativeFunctionValue NativeFunctionValue) String() string {
	return nativeFunctionValue.name + " function"
}

func (nativeFunctionValue NativeFunctionValue) Equals(other Value) (bool, error) {
	if otherNatVal, okNatVal := other.(NativeFunctionValue); okNatVal {
		return nativeFunctionValue.name == otherNatVal.name, nil
	}
	return false, nil
}

func golfcartAssert(args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, errors.New("assert() expects 2 arguments")
	}
	equal, err := args[0].Equals(args[1])
	if err != nil {
		return nil, err
	}
	if !equal {
		return nil, errors.New("assert failed: " + args[0].String() + " == " + args[1].String())
	}
	return NilValue{}, nil
}

func golfcartIn(args []Value) (Value, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return StringValue{val: []byte(scanner.Text())}, nil
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
	if strValue, okStr := value.(StringValue); okStr {
		return strValue, nil
	}
	if numValue, okNum := value.(NumberValue); okNum {
		return StringValue{val: []byte(nvToS(numValue))}, nil
	}
	if boolValue, okBool := value.(BoolValue); okBool {
		if boolValue.val {
			return StringValue{val: []byte("true")}, nil
		}
		return StringValue{val: []byte("false")}, nil
	}

	return nil, errors.New("str() expects 1 argument of type string, number, or bool")
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
	f, err := strconv.ParseFloat(string(strValue.val), 64)
	if err != nil {
		return nil, errors.New("num() couldn't convert " + strValue.String() + " to num")
	}
	return NumberValue{val: f}, nil
}

func golfcartType(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, errors.New("type() expects 1 argument")
	}
	value := args[0]
	switch value.(type) {
	case IdentifierValue:
		return StringValue{val: []byte("identifier")}, nil
	case StringValue:
		return StringValue{val: []byte("string")}, nil
	case NumberValue:
		return StringValue{val: []byte("number")}, nil
	case BoolValue:
		return StringValue{val: []byte("bool")}, nil
	case FunctionValue, NativeFunctionValue:
		return StringValue{val: []byte("function")}, nil
	case ListValue:
		return StringValue{val: []byte("list")}, nil
	case DictValue:
		return StringValue{val: []byte("dict")}, nil
	case NilValue:
		return StringValue{val: []byte("nil")}, nil
	}
	return nil, errors.New("unknown type")
}

func golfcartLen(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, errors.New("len() expects 1 argument of type string, list, or dict")
	}
	value := args[0]
	if stringVal, okStr := value.(StringValue); okStr {
		return NumberValue{val: float64(len(stringVal.val))}, nil
	}
	if listVal, okList := value.(ListValue); okList {
		return NumberValue{val: float64(len(listVal.val))}, nil
	}
	if dictVal, okDict := value.(DictValue); okDict {
		return NumberValue{val: float64(len(dictVal.val))}, nil
	}
	return nil, errors.New("len() expects 1 argument of type string, list, or dict")
}

func golfcartKeys(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, errors.New("keys() expects 1 argument of type dict")
	}
	value := args[0]
	if dictVal, okDict := value.(DictValue); okDict {
		keys := make(map[int]*Value, len(dictVal.val))
		i := 0
		for k := range dictVal.val {
			var value Value
			value = StringValue{val: []byte(k)}
			keys[i] = &value
			i++
		}
		return ListValue{val: keys}, nil
	}
	return nil, errors.New("keys() expects 1 argument of type dict")
}

func golfcartValues(args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, errors.New("values() expects 1 argument of type dict")
	}
	value := args[0]
	if dictVal, okDict := value.(DictValue); okDict {
		values := make(map[int]*Value, len(dictVal.val))
		i := 0
		for _, v := range dictVal.val {
			valueCopy := *v
			values[i] = &valueCopy
			i++
		}
		return ListValue{val: values}, nil
	}
	return nil, errors.New("values() expects 1 argument of type dict")
}
