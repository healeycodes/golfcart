package golfcart

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

type Context struct {
	stackFrame StackFrame
}

func (context *Context) Init() {
	context.stackFrame = StackFrame{entries: make(map[string]Value)}
}

type StackFrame struct {
	entries map[string]Value
	parent  *StackFrame
}

func (frame *StackFrame) String() string {
	s := ""
	for {
		s += "{\n"
		for key, value := range frame.entries {
			s += fmt.Sprintf("\t %v: %v\n", key, value)
		}
		s += "}"
		if parent := frame.parent; parent != nil {
			frame = parent
		} else {
			break
		}
	}
	return s
}

func (frame *StackFrame) GetChild() *StackFrame {
	childFrame := StackFrame{parent: frame, entries: make(map[string]Value)}
	return &childFrame
}

func (frame *StackFrame) Get(key string) (Value, error) {
	for {
		value, ok := frame.entries[key]
		if ok {
			return value, nil
		}
		if parent := frame.parent; parent != nil {
			frame = parent
		} else {
			break
		}
	}
	return nil, fmt.Errorf("cannot find value for key '%v'", key)
}

func (frame *StackFrame) Set(key string, value Value) {
	currentFrame := frame
	for {
		_, ok := frame.entries[key]
		if ok {
			frame.entries[key] = value
			return
		}
		if parent := frame.parent; parent != nil {
			frame = parent
		} else {
			break
		}
	}
	currentFrame.entries[key] = value
}

type Value interface {
	String() string
	Equals(Value) (bool, error)
}

func formatValues(values []Value) string {
	s := make([]string, len(values))
	for i, value := range values {
		s[i] = value.String()
	}
	return "(" + strings.Join(s, ", ") + ")"
}

type ReferenceValue struct {
	val *Value
}

func (_ ReferenceValue) String() string {
	return "reference"
}

func (_ ReferenceValue) Equals(other Value) (bool, error) {
	return false, nil
}

func unref(value Value) Value {
	if refValue, okRef := value.(ReferenceValue); okRef {
		return *refValue.val
	}
	return value
}

func unwrap(value Value, frame *StackFrame) (Value, error) {
	if idValue, okId := value.(IdentifierValue); okId {
		return frame.Get(idValue.val)
	}
	value = unref(value)
	return value, nil
}

type NilValue struct{}

func (numberValue NilValue) String() string {
	return "nil"
}

func (numberValue NilValue) Equals(other Value) (bool, error) {
	if _, ok := unref(other).(NilValue); ok {
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
	return false, fmt.Errorf("tried to compare with an uninitialized identifier: %v %v", identifierValue, other)
}

func (idValue IdentifierValue) Eval(frame *StackFrame) (Value, error) {
	value, err := frame.Get(idValue.val)
	if err != nil {
		return nil, err
	}
	return value, nil
}

type NumberValue struct {
	val float64
}

func (numberValue NumberValue) String() string {
	return nToS(numberValue.val)
}

func (numberValue NumberValue) Equals(other Value) (bool, error) {
	if other, ok := unref(other).(NumberValue); ok {
		return numberValue.val == other.val, nil
	}
	return false, nil
}

func nvToS(numberValue NumberValue) string {
	return nToS(numberValue.val)
}

func nToS(n float64) string {
	return strconv.FormatFloat(n, 'f', -1, 64)
}

type StringValue struct {
	val []byte
}

func (stringValue StringValue) String() string {
	return string(stringValue.val)
}

func (stringValue StringValue) Equals(other Value) (bool, error) {
	if otherStr, ok := unref(other).(StringValue); ok {
		a := stringValue.val
		b := otherStr.val
		if len(a) != len(b) {
			return false, nil
		}
		for i := range a {
			if a[i] != b[i] {
				return false, nil
			}
		}
		return true, nil
	}
	return false, nil
}

type BoolValue struct {
	val bool
}

func (boolValue BoolValue) String() string {
	return fmt.Sprintf("%t", boolValue.val)
}

func (boolValue BoolValue) Equals(other Value) (bool, error) {
	if otherValue, okBool := unref(other).(BoolValue); okBool {
		return boolValue.val == otherValue.val, nil
	}
	return false, nil
}

type ReturnValue struct {
	pos lexer.Position
	val Value
}

func (returnValue ReturnValue) Error() string {
	return fmt.Sprintf("%v return expression used outside of a function", returnValue.pos)
}

type BreakValue struct {
	pos lexer.Position
}

func (breakValue BreakValue) Error() string {
	return fmt.Sprintf("%v break expression used outside of a for loop", breakValue.pos)
}

type ContinueValue struct {
	pos lexer.Position
}

func (continueValue ContinueValue) Error() string {
	return fmt.Sprintf("%v continue expression used outside of a for loop", continueValue.pos)
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
	callFrame := functionValue.frame.GetChild()
	if len(args) != len(functionValue.parameters) {
		return nil, fmt.Errorf("function called with incorrect number of arguments, wanted: %v, got: %v", len(functionValue.parameters), formatValues(args))
	}
	for i, parameter := range functionValue.parameters {
		callFrame.Set(parameter, args[i])
	}
	var result Value
	result = NilValue{}
	var err error
	for _, expression := range functionValue.expressions {
		result, err = expression.Eval(callFrame)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type ListValue struct {
	val map[int]*Value
}

func (listValue ListValue) String() string {
	formatted := make([]string, len(listValue.val))
	for i, item := range listValue.val {
		formatted[i] = (*item).String()
	}
	return "[" + strings.Join(formatted, ", ") + "]"
}

func (listValue ListValue) Equals(other Value) (bool, error) {
	return false, nil
}

func (listValue ListValue) Append(other Value) {
	listValue.val[len(listValue.val)] = &other
}

func (listValue ListValue) Prepend(other Value) {
	// Add a new zeroth item.
	// Correcting the remaining indexes costs O(N)
	for i := len(listValue.val); i > 0; i-- {
		listValue.val[i] = listValue.val[i-1]
	}
	listValue.val[0] = &other
}

func (listValue ListValue) Pop() Value {
	last := *listValue.val[len(listValue.val)-1]
	delete(listValue.val, len(listValue.val)-1)
	return last
}

func (listValue ListValue) PopLeft() Value {
	// Remove and return the zeroth item.
	// Correcting the remaining indexes costs O(N)
	first := *listValue.val[0]
	delete(listValue.val, 0)
	for i := 0; i < len(listValue.val); i++ {
		listValue.val[i] = listValue.val[i+1]
	}
	delete(listValue.val, len(listValue.val)-1)
	return first
}

type DictValue struct {
	val map[string]*Value
}

func (dictValue *DictValue) Get(key string) (*Value, error) {
	value, ok := dictValue.val[key]
	if ok {
		return value, nil
	}
	return nil, fmt.Errorf("cannot find value for key: '%v'", key)
}

func (dictValue *DictValue) GetOrSet(key string, newValue *Value) *Value {
	existingValue, ok := dictValue.val[key]
	if ok {
		return existingValue
	}
	dictValue.val[key] = newValue
	return newValue
}

func (dictValue *DictValue) Set(key string, value Value) {
	dictValue.val[key] = &value
}

func (dictValue DictValue) String() string {
	s := make([]string, 0)
	s = append(s, "{")
	for key, value := range dictValue.val {
		s = append(s, fmt.Sprintf("%v: %v", key, *value))
	}
	s = append(s, "}")
	return strings.Join(s, "")
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

func (exprList ExpressionList) Equals(other Value) (bool, error) {
	return false, nil
}

func (exprList ExpressionList) Eval(context *Context) (Value, error) {
	var result Value
	result = NilValue{}
	var err error
	for _, node := range exprList.Expressions {
		result, err = node.Eval(&context.stackFrame)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (expr Expression) String() string {
	return "expr"
}

func (expr Expression) Equals(other Value) (bool, error) {
	return false, nil
}

func (expr Expression) Eval(frame *StackFrame) (Value, error) {
	if expr.Assignment != nil {
		result, err := expr.Assignment.Eval(frame)
		if err != nil {
			return nil, err
		}
		if idValue, ok := result.(IdentifierValue); ok {
			value, err := frame.Get(idValue.val)
			if err != nil {
				return nil, err
			}
			return value, nil
		}
		return result, nil
	}
	panic("unimplemented Expression Eval")
}

func (assignment Assignment) String() string {
	return "assignment"
}

func (assignment Assignment) Equals(other Value) (bool, error) {
	return false, nil
}

func (assignment Assignment) Eval(frame *StackFrame) (Value, error) {
	left, err := assignment.LogicAnd.Eval(frame)
	if err != nil {
		return nil, err
	}
	leftRef, leftRefOk := left.(ReferenceValue)

	if assignment.Op == "" {
		if leftRefOk {
			return *leftRef.val, nil
		}
		return left, nil
	}
	right, err := assignment.Next.Eval(frame)
	if err != nil {
		return nil, err
	}
	if assignment.Op == "=" {
		if idValue, okId := right.(IdentifierValue); okId {
			right, err = frame.Get(idValue.val)
			if err != nil {
				return nil, err
			}
		}
		if leftRefOk {
			*leftRef.val = right
			return right, nil
		}
		if leftId, okId := left.(IdentifierValue); okId {
			frame.Set(leftId.val, right)
			return right, nil
		}
		return nil, fmt.Errorf("%v can't assign to non-identifier: %v", assignment.Pos, left)
	}
	panic("unreachable Assignment Eval")
}

func (logicAnd LogicAnd) String() string {
	return "equality"
}

func (logicAnd LogicAnd) Equals(other Value) (bool, error) {
	return false, nil
}

func (logicAnd LogicAnd) Eval(frame *StackFrame) (Value, error) {
	left, err := logicAnd.LogicOr.Eval(frame)
	if err != nil {
		return nil, err
	}
	if logicAnd.Op == "" {
		return left, nil
	}
	right, err := logicAnd.Next.Eval(frame)
	if err != nil {
		return nil, err
	}
	left, err = unwrap(left, frame)
	if err != nil {
		return nil, err
	}
	right, err = unwrap(right, frame)
	if err != nil {
		return nil, err
	}

	if boolValue, okBool := left.(BoolValue); okBool {
		if boolValue.val {
			if boolValue, okBool := right.(BoolValue); okBool {
				if boolValue.val {
					return boolValue, nil
				}
			} else {
				return nil, fmt.Errorf("%v only bools can be compared with 'and', not: %v", logicAnd.Pos, right)
			}
		}
	} else {
		return nil, fmt.Errorf("%v only bools can be compared with 'and', not: %v", logicAnd.Pos, left)
	}

	return BoolValue{val: false}, nil
}

func (logicOr LogicOr) String() string {
	return "equality"
}

func (logicOr LogicOr) Equals(other Value) (bool, error) {
	return false, nil
}

func (logicOr LogicOr) Eval(frame *StackFrame) (Value, error) {
	left, err := logicOr.Equality.Eval(frame)
	if err != nil {
		return nil, err
	}
	if logicOr.Op == "" {
		return left, nil
	}
	right, err := logicOr.Next.Eval(frame)
	if err != nil {
		return nil, err
	}
	left, err = unwrap(left, frame)
	if err != nil {
		return nil, err
	}
	right, err = unwrap(right, frame)
	if err != nil {
		return nil, err
	}

	if boolValue, okBool := left.(BoolValue); okBool {
		if boolValue.val {
			return boolValue, nil
		}
	} else {
		return nil, fmt.Errorf("%v only bools can be compared with 'or', not: %v", logicOr.Pos, left)
	}
	if boolValue, okBool := right.(BoolValue); okBool {
		if boolValue.val {
			return boolValue, nil
		}
	} else {
		return nil, fmt.Errorf("%v only bools can be compared with 'or', not: %v", logicOr.Pos, right)
	}
	return BoolValue{val: false}, nil
}

func (equality Equality) String() string {
	return "equality"
}

func (equality Equality) Equals(other Value) (bool, error) {
	return false, nil
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
	left = unref(left)
	right = unref(right)

	if idValue, okId := left.(IdentifierValue); okId {
		value, err := frame.Get(idValue.val)
		if err != nil {
			return nil, err
		}
		left = value
	}
	if idValue, okId := right.(IdentifierValue); okId {
		value, err := frame.Get(idValue.val)
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
	panic("unreachable Equality Eval")
}

func (comparison Comparison) String() string {
	return "comparison"
}

func (comparison Comparison) Equals(other Value) (bool, error) {
	return false, nil
}

func (comparison Comparison) Eval(frame *StackFrame) (Value, error) {
	left, err := comparison.Addition.Eval(frame)
	if err != nil {
		return nil, err
	}
	if comparison.Op == "" {
		return left, nil
	}
	right, err := comparison.Next.Eval(frame)
	if err != nil {
		return nil, err
	}

	left, err = unwrap(left, frame)
	if err != nil {
		return nil, err
	}
	right, err = unwrap(right, frame)
	if err != nil {
		return nil, err
	}

	if leftNum, okNum := left.(NumberValue); okNum {
		if rightNum, okNum := right.(NumberValue); okNum {
			return BoolValue{val: comparison.Op == "<" && leftNum.val < rightNum.val ||
				comparison.Op == "<=" && leftNum.val <= rightNum.val ||
				comparison.Op == ">" && leftNum.val > rightNum.val ||
				comparison.Op == ">=" && leftNum.val >= rightNum.val}, nil
		}
	}
	leftType, err := golfcartType([]Value{left})
	if err != nil {
		return nil, err
	}
	rightType, err := golfcartType([]Value{right})
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("%v only numbers can be compared: %v %v %v", comparison.Pos, leftType, comparison.Op, rightType)
}

func (addition Addition) String() string {
	return "addition"
}

func (addition Addition) Equals(other Value) (bool, error) {
	return false, nil
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

	left, err = unwrap(left, frame)
	if err != nil {
		return nil, err
	}
	right, err = unwrap(right, frame)
	if err != nil {
		return nil, err
	}

	err_msg := fmt.Errorf("%v '+' can only be used between [string, string], [number, number], [list, list], not: [%v, %v]",
		addition.Multiplication.Pos, left, right)

	leftStr, okLeft := left.(StringValue)
	rightStr, okRight := right.(StringValue)
	if addition.Op == "+" && (okLeft && !okRight || okRight && !okLeft) {
		return nil, err_msg
	} else if addition.Op == "+" && okLeft && okRight {
		return StringValue{val: append([]byte{}, append(leftStr.val, rightStr.val...)...)}, nil
	}

	leftList, okLeft := left.(ListValue)
	rightList, okRight := right.(ListValue)
	if addition.Op == "+" && (okLeft && !okRight || okRight && !okLeft) {
		return nil, err_msg
	} else if addition.Op == "+" && okLeft && okRight {
		newMap := ListValue{val: map[int]*Value{}}
		for i, value := range leftList.val {
			newMap.val[i] = value
		}
		len := len(newMap.val)
		for i, value := range rightList.val {
			newMap.val[i+len] = value
		}
		return newMap, nil
	}

	leftNum, okLeft := left.(NumberValue)
	rightNum, okRight := right.(NumberValue)
	if addition.Op == "+" && (okLeft && !okRight || okRight && !okLeft) {
		return nil, err_msg
	} else if addition.Op == "+" && okLeft && okRight {
		return NumberValue{val: leftNum.val + rightNum.val}, nil
	} else if addition.Op == "-" && okLeft && okRight {
		return NumberValue{val: leftNum.val - rightNum.val}, nil
	}
	panic("unreachable Addition Eval")
}

func (multiplication Multiplication) String() string {
	return "multiplication"
}

func (multiplication Multiplication) Equals(other Value) (bool, error) {
	return false, nil
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

	left, err = unwrap(left, frame)
	if err != nil {
		return nil, err
	}
	right, err = unwrap(right, frame)
	if err != nil {
		return nil, err
	}

	leftNum, okLeft := left.(NumberValue)
	if !okLeft {
		return nil, fmt.Errorf("%v '*' and '/' only supported between numbers", multiplication.Unary.Pos)
	}
	rightNum, okRight := right.(NumberValue)
	if !okRight {
		return nil, fmt.Errorf("%v '*' and '/' only supported between numbers", multiplication.Unary.Pos)
	}
	if multiplication.Op == "*" {
		return NumberValue{val: leftNum.val * rightNum.val}, nil
	}
	if multiplication.Op == "/" {
		return NumberValue{val: leftNum.val / rightNum.val}, nil
	}
	if multiplication.Op == "%" {
		return NumberValue{val: float64(int(math.Round(leftNum.val)) % int(math.Round(rightNum.val)))}, nil
	}
	panic("unreachable Multiplication Eval")
}

func (unary Unary) String() string {
	return "unary"
}

func (unary Unary) Equals(other Value) (bool, error) {
	return false, nil
}

func (unary Unary) Eval(frame *StackFrame) (Value, error) {
	if unary.Op == "!" {
		value, err := unary.Unary.Eval(frame)
		if err != nil {
			return nil, err
		}
		value, err = unwrap(value, frame)
		if err != nil {
			return nil, err
		}
		if boolValue, ok := value.(BoolValue); ok {
			return BoolValue{val: !boolValue.val}, nil
		}
		return nil, fmt.Errorf("%v expected bool after '!'", unary.Pos)
	}
	if unary.Op == "-" {
		value, err := unary.Unary.Eval(frame)
		if err != nil {
			return nil, err
		}
		value, err = unwrap(value, frame)
		if err != nil {
			return nil, err
		}
		if numberValue, ok := value.(NumberValue); ok {
			return NumberValue{val: -numberValue.val}, nil
		}
		return nil, fmt.Errorf("%v expected number after '-'", unary.Pos)
	}
	return unary.Primary.Eval(frame)
}

func (primary Primary) String() string {
	return "primary"
}

func (primary Primary) Eval(frame *StackFrame) (Value, error) {
	if ifExpression := primary.If; ifExpression != nil {
		return ifExpression.Eval(frame)
	}
	if primary.DataLiteral != nil {
		if functionLiteral := primary.DataLiteral.FunctionLiteral; functionLiteral != nil {
			return functionLiteral.Eval(frame)
		}
		if listLiteral := primary.DataLiteral.ListLiteral; listLiteral != nil {
			return listLiteral.Eval(frame)
		}
		if dictLiteral := primary.DataLiteral.DictLiteral; dictLiteral != nil {
			return dictLiteral.Eval(frame)
		}
	}
	if subExpression := primary.SubExpression; subExpression != nil {
		return subExpression.Eval(frame)
	}
	if call := primary.Call; call != nil {
		return call.Eval(frame)
	}
	if returnVal := primary.Return; returnVal != nil {
		value, err := returnVal.Expression.Eval(frame)
		if err != nil {
			return nil, err
		}
		return nil, ReturnValue{pos: returnVal.Pos, val: value}
	}
	if primary.Break != nil {
		return nil, BreakValue{pos: primary.Break.Pos}
	}
	if primary.Continue != nil {
		return nil, BreakValue{pos: primary.Continue.Pos}
	}
	if forExpression := primary.For; forExpression != nil {
		return forExpression.Eval(frame)
	}
	if forWhileExpression := primary.ForWhile; forWhileExpression != nil {
		forExpression := For{
			Condition: forWhileExpression.Condition,
			Body:      forWhileExpression.Body,
		}
		return forExpression.Eval(frame)
	}
	if forKeyExpression := primary.ForValue; forKeyExpression != nil {
		value := forKeyExpression.Value
		collection := forKeyExpression.Collection
		collectionExpression := forKeyExpression.CollectionExpression
		return evalForKeyValue(nil, value, collection, collectionExpression, forKeyExpression.Body, frame)
	}
	if forKeyExpression := primary.ForKeyValue; forKeyExpression != nil {
		key := forKeyExpression.Key
		value := forKeyExpression.Value
		collection := forKeyExpression.Collection
		collectionExpression := forKeyExpression.CollectionExpression
		return evalForKeyValue(key, value, collection, collectionExpression, forKeyExpression.Body, frame)
	}
	if primary.Number != nil {
		return NumberValue{val: *primary.Number}, nil
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
		return StringValue{val: []byte(*primary.Str)[1 : len((*primary.Str))-1]}, nil
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

func (ifExpression If) String() string {
	return "if expression"
}

func (ifExpression If) Equals(other Value) (bool, error) {
	return false, nil
}

func (ifExpression If) Eval(frame *StackFrame) (Value, error) {
	ifFrame := frame.GetChild()
	condition, err := ifExpression.Condition.Eval(ifFrame)
	if err != nil {
		return nil, err
	}
	var result Value
	result = NilValue{}
	if boolValue, okBool := condition.(BoolValue); okBool {
		var err error
		if boolValue.val {
			for _, expr := range ifExpression.IfBody {
				result, err = (*expr).Eval(ifFrame)
				if err != nil {
					return nil, err
				}
			}
			return result, nil
		}
		if ifExpression.ElseIf != nil {
			// TODO: there's some duplicated logic here
			// but it will only ever be in two places
			// If and ElseIf – not sure if it's worth refactoring?
			current := ifExpression.ElseIf
			for current != nil {
				condition, err := current.Condition.Eval(ifFrame)
				if err != nil {
					return nil, err
				}
				var result Value
				result = NilValue{}
				if boolValue, okBool := condition.(BoolValue); okBool {
					var err error
					if boolValue.val {
						for _, expr := range current.IfBody {
							result, err = (*expr).Eval(ifFrame)
							if err != nil {
								return nil, err
							}
						}
						return result, nil
					}
				} else {
					return nil, fmt.Errorf("%v if expression conditional should evaluate to true or false",
						ifExpression.Pos)
				}
				current = current.Next
			}
		}
		if ifExpression.ElseBody != nil {
			for _, expr := range ifExpression.ElseBody {
				result, err = (*expr).Eval(ifFrame)
				if err != nil {
					return nil, err
				}
			}
			return result, nil
		}
	} else {
		return nil, fmt.Errorf("%v if expression conditional should evaluate to true or false",
			ifExpression.Pos)
	}
	return result, nil
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
	dictValue := DictValue{val: make(map[string]*Value)}
	if dictLiteral.DictEntry != nil {
		for _, dictEntry := range *dictLiteral.DictEntry {
			var key string
			if dictEntry.Key != nil {
				value, err := dictEntry.Key.Eval(frame)
				if err != nil {
					return nil, err
				}
				if strValue, okStr := value.(StringValue); okStr {
					key = string(strValue.val)
				}
			} else if dictEntry.Ident != nil {
				key = *dictEntry.Ident
			}

			value, err := dictEntry.Value.Eval(frame)
			if err != nil {
				return nil, err
			}
			if key == "" {
				return nil, fmt.Errorf("%v can't set empty string as dict key – did you forget to wrap a number with \"\" quote marks?",
					dictLiteral.Pos)
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
	values := make(map[int]*Value, 0)
	if listLiteral.Expressions != nil {
		for _, expression := range *listLiteral.Expressions {
			result, err := expression.Eval(frame)
			if err != nil {
				return nil, err
			}
			values[len(values)] = &result
		}
	}
	return ListValue{val: values}, nil
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
		value, err = frame.Get(*ident)
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

	chainCall := call.CallChain
	for chainCall != nil {
		value = unref(value)
		var args []Value
		if parameters := chainCall.Parameters; parameters != nil {
			args, err = parseArgs(chainCall.Parameters, frame)
			if err != nil {
				return nil, err
			}
		}

		var access Value
		if chainCall.Access != nil {
			access = IdentifierValue{val: *chainCall.Access}
		}
		if chainCall.ComputedAccess != nil {
			access, err = chainCall.ComputedAccess.Eval(frame)
			if err != nil {
				return nil, err
			}
		}

		if listValue, okList := value.(ListValue); okList && access != nil {
			if idVal, okId := access.(IdentifierValue); okId {
				if idVal.val == "append" {
					err = listAppend(listValue, chainCall, frame)
				} else if idVal.val == "prepend" {
					err = listPrepend(listValue, chainCall, frame)
				} else if idVal.val == "pop" {
					if len(listValue.val) == 0 {
						err = fmt.Errorf("cannot pop() from an empty list")
					} else {
						return listValue.Pop(), nil
					}
				} else if idVal.val == "pop_left" {
					if len(listValue.val) == 0 {
						err = fmt.Errorf("cannot pop_left() from an empty list")
					} else {
						return listValue.PopLeft(), nil
					}
				}
				if err != nil {
					return nil, err
				}
				return NilValue{}, nil
			}
			value, err = listAccess(listValue, access)
			if err != nil {
				return nil, err
			}
		}
		if dictValue, okDict := value.(DictValue); okDict && access != nil {
			value, err = dictAccess(dictValue, access)
			if err != nil {
				return nil, err
			}
		}
		if functionValue, okFunc := value.(FunctionValue); okFunc {
			value, err = functionValue.Exec(args)
			if returnValue, okRet := err.(ReturnValue); okRet {
				value = returnValue.val
			} else if err != nil {
				return nil, err
			}
		}
		if nativeFunctionValue, okNatFunc := value.(NativeFunctionValue); okNatFunc {
			value, err = nativeFunctionValue.Exec(args)
			if err != nil {
				return nil, err
			}
		}
		if stringValue, okStr := value.(StringValue); okStr && access != nil {
			value, err = stringAccess(stringValue, access)
			if err != nil {
				return nil, err
			}
		}
		if chainCall.Next != nil {
			chainCall = chainCall.Next
		} else {
			chainCall = nil
		}
	}

	return value, nil
}

func parseArgs(expressions *[]Expression, frame *StackFrame) ([]Value, error) {
	args := make([]Value, len(*expressions))
	for i, parameter := range *expressions {
		result, err := parameter.Eval(frame)
		if err != nil {
			return nil, err
		}
		argValue := unref(result)
		args[i] = argValue
	}
	return args, nil
}

func stringAccess(stringValue StringValue, access Value) (Value, error) {
	if numValue, okNum := access.(NumberValue); okNum {
		index := int(numValue.val)
		if index < 0 || index > len(stringValue.val)-1 {
			return nil, fmt.Errorf("string access out of bounds: %v", index)
		}
		return StringValue{val: []byte{stringValue.val[index]}}, nil
	}

	value, err := golfcartType([]Value{access})
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("string access expects 1 argument of type number, not: %v", value)
}

func listAccess(listValue ListValue, access Value) (Value, error) {
	if numValue, okNum := access.(NumberValue); okNum {
		index := int(numValue.val)
		if index < 0 || index > len(listValue.val)-1 {
			return nil, fmt.Errorf("list access out of bounds: %v", index)
		}
		return ReferenceValue{val: listValue.val[index]}, nil
	}

	value, err := golfcartType([]Value{access})
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("list access expects 1 argument of type number, not %v", value)
}

func listAppend(listValue ListValue, chainCall *CallChain, frame *StackFrame) error {
	if chainCall.Next == nil || chainCall.Next.Parameters == nil || len(*chainCall.Next.Parameters) != 1 {
		return fmt.Errorf("append() expects 1 argument")
	}
	args, err := parseArgs(chainCall.Next.Parameters, frame)
	if err != nil {
		return err
	}
	listValue.Append(args[0])
	return nil
}

func listPrepend(listValue ListValue, chainCall *CallChain, frame *StackFrame) error {
	if chainCall.Next == nil || chainCall.Next.Parameters == nil || len(*chainCall.Next.Parameters) != 1 {
		return fmt.Errorf("prepend() expects 1 argument")
	}
	args, err := parseArgs(chainCall.Next.Parameters, frame)
	if err != nil {
		return err
	}
	listValue.Prepend(args[0])
	return nil
}

func dictAccess(dictValue DictValue, access Value) (Value, error) {
	var key string
	if strValue, okStr := access.(StringValue); okStr {
		key = string(strValue.val)
	} else if idValue, okId := access.(IdentifierValue); okId {
		key = string(idValue.val)
	} else if numValue, okNum := access.(NumberValue); okNum {
		key = nvToS(numValue)
	} else {
		golfType, err := golfcartType([]Value{idValue})
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Only strings are allowed as dict keys, not: %v", golfType)
	}
	var newValue Value
	newValue = NilValue{}
	value := dictValue.GetOrSet(key, &newValue)
	return ReferenceValue{val: value}, nil
}

func evalForKeyValue(keyIdent *string, valueIdent *string, collectionIdent *string, collectionExpression *Expression, expressions []*Expression, frame *StackFrame) (Value, error) {
	iterations := NumberValue{val: 0}
	forFrame := frame.GetChild()
	var values Value
	var err error
	if collectionIdent != nil {
		values, err = frame.Get(*collectionIdent)
	} else {
		values, err = collectionExpression.Eval(forFrame)
	}
	if err != nil {
		return nil, err
	}
	iterableValues := make([]Value, 0)
	iterableKeys := make([]Value, 0)
	if listValue, okList := values.(ListValue); okList {
		for k, v := range listValue.val {
			iterableKeys = append(iterableKeys, NumberValue{val: float64(k)})
			iterableValues = append(iterableValues, *v)
		}
	}
	if dictVal, okDict := values.(DictValue); okDict {
		for k, v := range dictVal.val {
			iterableKeys = append(iterableKeys, StringValue{val: []byte(k)})
			iterableValues = append(iterableValues, *v)
		}
	}
	if strVal, okStr := values.(StringValue); okStr {
		for k, v := range strVal.val {
			iterableKeys = append(iterableKeys, NumberValue{val: float64(k)})
			iterableValues = append(iterableValues, StringValue{val: []byte{v}})
		}
	}

	for i := 0; i < len(iterableValues); i++ {
		forFrame.Set(*valueIdent, iterableValues[i])
		if keyIdent != nil {
			forFrame.Set(*keyIdent, iterableKeys[i])
		}
		var err error
		iterations.val++
		for _, expr := range expressions {
			_, err = (*expr).Eval(forFrame)
			if _, okBreak := err.(BreakValue); okBreak {
				return iterations, nil
			}
			if _, okCont := err.(ContinueValue); okCont {
				continue
			}
			if err != nil {
				return nil, err
			}
		}
	}
	return iterations, nil
}

func (forExpression For) Eval(frame *StackFrame) (Value, error) {
	iterations := NumberValue{val: 0}
	forFrame := frame.GetChild()
	if forExpression.Init != nil {
		for _, assignExpr := range forExpression.Init {
			_, err := (*assignExpr).Eval(forFrame)
			if err != nil {
				return nil, err
			}
		}
	}
	for {
		var condition Value
		var err error
		if forExpression.Condition != nil {
			condition, err = forExpression.Condition.Eval(forFrame)
			if err != nil {
				return nil, err
			}
		} else {
			// Fake a condition for infinite loop variant `for {}`
			condition = BoolValue{val: true}
		}
		if boolValue, okBool := condition.(BoolValue); okBool {
			if boolValue.val {
				iterations.val++
				for _, expr := range forExpression.Body {
					_, err = (*expr).Eval(forFrame)
					if _, okBreak := err.(BreakValue); okBreak {
						return iterations, nil
					}
					if _, okCont := err.(ContinueValue); okCont {
						continue
					}
					if err != nil {
						return nil, err
					}
				}
			} else {
				break
			}
		} else {
			valueType, err := golfcartType([]Value{condition})
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("condition expression of for loop should be of type bool not: %v", valueType)
		}
		// For-variants e.g. `for true {}` `for {}` don't have a post expression
		if forExpression.Post != nil {
			_, err = forExpression.Post.Eval(forFrame)
			if err != nil {
				return nil, err
			}
		}
	}
	return iterations, nil
}
