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
	for {
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
	for {
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
	for {
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
	Equals(Value) (bool, error)
}

type NilValue struct{}

func (numberValue NilValue) String() string {
	return "nil"
}

func (numberValue NilValue) Equals(other Value) (bool, error) {
	if _, ok := other.(NilValue); ok {
		return true, nil
	}
	return false, nil
}

type IdentifierValue struct {
	val string
}

func (identifierValue IdentifierValue) String() string {
	return identifierValue.val
}

func (identifierValue IdentifierValue) Equals(other Value) (bool, error) {
	return false, errors.New("tried to compare with an uninitialized identifier: " +
		identifierValue.String() + " " + other.String())
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

func (numberValue NumberValue) Equals(other Value) (bool, error) {
	if other, ok := other.(NumberValue); ok {
		return numberValue.val == other.val, nil
	}
	return false, errors.New("only numbers can be compared with numbers")
}

type StringValue struct {
	val string
}

func (stringValue StringValue) String() string {
	return stringValue.val
}

func (stringValue StringValue) Equals(other Value) (bool, error) {
	if other, ok := other.(StringValue); ok {
		return stringValue.val == other.val, nil
	}
	panic("unimplemented StringValue Equals")
}

type BoolValue struct {
	val bool
}

func (boolValue BoolValue) String() string {
	return fmt.Sprintf("%t", boolValue.val)
}

func (boolValue BoolValue) Equals(other Value) (bool, error) {
	return boolValue == other, nil
}

type FunctionValue struct {
	parameters  []string
	frame       *StackFrame
	expressions []*Expression
}

func (functionValue FunctionValue) String() string {
	return "function"
}

func (functionValue FunctionValue) Equals(other Value) (bool, error) {
	return false, nil
}

func (functionValue FunctionValue) Exec(args []Value) (Value, error) {
	if len(args) != len(functionValue.parameters) {
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

type ListValue struct {
	items []Value
}

func (listValue ListValue) String() string {
	s := ""
	for _, item := range listValue.items {
		s += item.String() + " "
	}
	return "[ " + s + "]"
}

func (listValue ListValue) Equals(other Value) (bool, error) {
	return false, nil
}

type DictValue struct {
	entries map[string]Value
}

func (dictValue *DictValue) Get(key Value) (Value, error) {
	value, ok := dictValue.entries[key.String()]
	if ok {
		return value, nil
	}

	return nil, errors.New("cannot find value for '" + key.String() + "'")
}

func (dictValue *DictValue) Set(key Value, value Value) {
	dictValue.entries[key.String()] = value
}

func (dictValue DictValue) String() string {
	return "dictValue"
}

func (dictValue DictValue) Equals(other Value) (bool, error) {
	return false, nil
}

// --

func (exprList ExpressionList) String() string {
	s := make([]string, 0)
	for _, expr := range exprList.Expressions {
		s = append(s, expr.String())
	}
	return strings.Join(s, ", ")
}

func (exprList ExpressionList) Equals(_ Value) (bool, error) {
	return false, nil
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

func (expr Expression) Equals(other Value) (bool, error) {
	return expr == other, nil
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
	panic("unimplemented Expression Eval")
}

func (functionLiteral FunctionLiteral) String() string {
	return "functionLiteral"
}

func (functionLiteral FunctionLiteral) Equals(other Value) (bool, error) {
	return false, nil
}

func (functionLiteral FunctionLiteral) Eval(frame *StackFrame) (Value, error) {
	closureFrame := frame.GetChild()
	functionValue := FunctionValue{parameters: functionLiteral.Parameters, frame: closureFrame, expressions: functionLiteral.Body}
	return functionValue, nil
}

func (dictLiteral DictLiteral) String() string {
	return "dictLiteral"
}

func (dictLiteral DictLiteral) Equals(other Value) (bool, error) {
	return false, nil
}

func (dictLiteral DictLiteral) Eval(frame *StackFrame) (Value, error) {
	entries := make(map[string]Value)
	dictValue := DictValue{entries: entries}
	if dictLiteral.DictEntry != nil {
		for _, dictEntry := range *dictLiteral.DictEntry {
			var key Value
			var err error
			if dictEntry.Key != nil {
				key, err = dictEntry.Key.Eval(frame)
				if err != nil {
					return nil, err
				}
			} else if dictEntry.Ident != nil {
				key = IdentifierValue{val: *dictEntry.Ident}
			}

			value, err := dictEntry.Value.Eval(frame)
			if err != nil {
				return nil, err
			}
			dictValue.Set(key, value)
		}
	}

	return dictValue, nil
}

func (listLiteral ListLiteral) String() string {
	return "listLiteral"
}

func (listLiteral ListLiteral) Equals(other Value) (bool, error) {
	return false, nil
}

func (listLiteral ListLiteral) Eval(frame *StackFrame) (Value, error) {
	values := make([]Value, 0)
	if listLiteral.Expressions != nil {
		for _, expression := range *listLiteral.Expressions {
			result, err := expression.Eval(frame)
			if err != nil {
				return nil, err
			}
			values = append(values, result)
		}
		return ListValue{items: values}, nil
	}
	return ListValue{items: values}, nil
}

func (call Call) String() string {
	return "call"
}

func (call Call) Equals(other Value) (bool, error) {
	return false, nil
}

func (call Call) Eval(frame *StackFrame) (Value, error) {
	// TODO: pass the cursor location (call.Pos) for better errors?

	var value Value
	var err error
	if ident := call.Ident; ident != nil {
		value, err = frame.Get(IdentifierValue{val: *ident})
		if err != nil {
			return nil, err
		}
	}
	if subExpr := call.SubExpression; subExpr != nil {
		value, err = subExpr.Eval(frame)
		if err != nil {
			return nil, err
		}
	}

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

	var access Value
	if call.Access != nil {
		access = IdentifierValue{val: *call.Access}
	}
	if call.ComputedAccess != nil {
		value, err := call.ComputedAccess.Eval(frame)
		if err != nil {
			return nil, err
		}
		access = value
	}

	if functionValue, okFunc := value.(FunctionValue); okFunc {
		value, err := functionValue.Exec(args)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
	if nativeFunctionValue, okNatFunc := value.(NativeFunctionValue); okNatFunc {
		value, err := nativeFunctionValue.Exec(args)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
	if listValue, okList := value.(ListValue); okList && access != nil {
		return listAccess(listValue, access)
	}
	if dictValue, okDict := value.(DictValue); okDict {
		if access != nil {
			return dictAccess(dictValue, access)
		}
		return dictAccess(dictValue, access)
	}

	valueType, err := golfcartType([]Value{value})
	if err != nil {
		return nil, err
	}
	if call.Parameters != nil {
		s := make([]string, len(args))
		for i, arg := range args {
			s[i] = arg.String()
		}
		return nil, errors.New("there was a problem calling a value of type " + valueType.String() + " with " + strings.Join(s, ", "))
	}
	if call.Access != nil {
		return nil, errors.New("there was a problem calling a value of type " + valueType.String() + " with " + access.String())
	}
	if call.ComputedAccess != nil {
		return nil, errors.New("there was a problem calling a value of type " + valueType.String() + " with " + access.String())
	}
	panic("unreachable branch due to parse logic in Call Eval")
}

func listAccess(listValue ListValue, access Value) (Value, error) {
	if numValue, okNum := access.(NumberValue); okNum {
		index := int(numValue.val)
		if index < 0 || index > len(listValue.items)-1 {
			return nil, fmt.Errorf("list access out of bounds: %v", index)
		}
		return listValue.items[index], nil
	}

	value, err := golfcartType([]Value{access})
	if err != nil {
		return nil, err
	}
	return nil, errors.New("list access expects 1 argument of type number: not " + value.String())
}

func dictAccess(dictValue DictValue, access Value) (Value, error) {
	value, err := dictValue.Get(access)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (assignment Assignment) String() string {
	return "assignment"
}

func (assignment Assignment) Equals(other Value) (bool, error) {
	return assignment == other, nil
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
		if idValue, okId := right.(IdentifierValue); okId {
			right, err = frame.Get(idValue)
			if err != nil {
				return nil, err
			}
		}
		frame.Set(left, right)
		return NilValue{}, nil
	}
	panic("unreachable")
}

func (logicAnd LogicAnd) String() string {
	return "equality"
}

func (logicAnd LogicAnd) Equals(other Value) (bool, error) {
	return logicAnd == other, nil
}

func (logicAnd LogicAnd) Eval(frame *StackFrame) (Value, error) {
	left, err := logicAnd.LogicOr.Eval(frame)
	if err != nil {
		return nil, err
	}
	if logicAnd.Op == "" {
		return left, nil
	}
	panic("unreachable")
}

func (logicOr LogicOr) String() string {
	return "equality"
}

func (logicOr LogicOr) Equals(other Value) (bool, error) {
	return logicOr == other, nil
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

func (equality Equality) Equals(other Value) (bool, error) {
	return equality == other, nil
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

	if idValue, okId := left.(IdentifierValue); okId {
		value, err := frame.Get(idValue)
		if err != nil {
			return nil, err
		}
		left = value
	}
	if idValue, okId := right.(IdentifierValue); okId {
		value, err := frame.Get(idValue)
		if err != nil {
			return nil, err
		}
		right = value
	}

	result, err := left.Equals(right)
	if err != nil {
		return nil, err
	}
	if equality.Op == "==" {
		return BoolValue{val: result}, nil
	} else if equality.Op == "!=" {
		return BoolValue{val: !result}, nil
	}
	return nil, errors.New("unimplemented Equality Eval")
}

func (comparison Comparison) String() string {
	return "comparison"
}

func (comparison Comparison) Equals(other Value) (bool, error) {
	return comparison == other, nil
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

func (addition Addition) Equals(other Value) (bool, error) {
	return addition == other, nil
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

	panic("unreachable")
}

func (multiplication Multiplication) String() string {
	return "multiplication"
}

func (multiplication Multiplication) Equals(other Value) (bool, error) {
	return multiplication == other, nil
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
		return nil, errors.New(multiplication.Unary.Pos.String() + " '*' and '/' only supported between numbers")
	}
	rightNum, okRight := right.(NumberValue)
	if !okRight {
		return nil, errors.New(multiplication.Next.Pos.String() + " '*' and '/' only supported between numbers")
	}
	if multiplication.Op == "*" {
		return NumberValue{val: leftNum.val * rightNum.val}, nil
	}
	if multiplication.Op == "/" {
		return NumberValue{val: leftNum.val / rightNum.val}, nil
	}

	panic("unreachable")
}

func (unary Unary) String() string {
	return "unary"
}

func (unary Unary) Equals(other Value) (bool, error) {
	return unary == other, nil
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

func (primary Primary) Eval(frame *StackFrame) (Value, error) {
	if functionLiteral := primary.FunctionLiteral; functionLiteral != nil {
		return functionLiteral.Eval(frame)
	}
	if listLiteral := primary.ListLiteral; listLiteral != nil {
		return listLiteral.Eval(frame)
	}
	if dictLiteral := primary.DictLiteral; dictLiteral != nil {
		return dictLiteral.Eval(frame)
	}
	if subExpression := primary.SubExpression; subExpression != nil {
		return subExpression.Eval(frame)
	}
	if call := primary.Call; call != nil {
		return call.Eval(frame)
	}
	if ident := primary.Ident; ident != nil {
		identifierValue := IdentifierValue{val: *ident}
		return identifierValue, nil
	}
	if primary.Number != nil {
		return NumberValue{val: *primary.Number}, nil
	}
	if primary.Str != nil {
		// TODO: Parse strings without including quote `"` marks
		return StringValue{val: (*primary.Str)[1 : len((*primary.Str))-1]}, nil
	}
	if primary.True != nil {
		return BoolValue{val: true}, nil
	}
	if primary.False != nil {
		return BoolValue{val: false}, nil
	}
	if primary.Nil != nil {
		return NilValue{}, nil
	}

	panic("unimplemented Primary Eval")
}

func RunProgram(source string) (*string, error) {
	ast, err := GenerateAST(source)
	if err != nil {
		return nil, err
	}

	context := Context{}
	context.Init()
	InjectRuntime(&context)

	result, err := ast.Eval(&context)
	if err != nil {
		return nil, err
	}

	ret := result.String()
	return &ret, nil
}
