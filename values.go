package main

type ValueType int

const (
	ValueTypeNull ValueType = iota
	ValueTypeNumber
	ValueBoolean
	ValueObject
	ValueNativeFunction
)

type RuntimeVal interface {
}

type NullVal struct {
	Type  ValueType
	Value string
}

type NumberVal struct {
	Type  ValueType
	Value float64
}

type BoolVal struct {
	Type  ValueType
	Value bool
}

type ObjectVal struct {
	Type       ValueType
	properties map[string]RuntimeVal
}

type NativeFnValue struct {
	Type ValueType
	call FunctionCall
}

func makeNumber(n float64) *NumberVal {
	return &NumberVal{Type: ValueTypeNumber, Value: n}
}

func makeNull() *NullVal {
	return &NullVal{Type: ValueTypeNull, Value: "null"}
}

func makeBool(v bool) *BoolVal {
	return &BoolVal{Type: ValueBoolean, Value: v}
}

func makeNativeFn(call FunctionCall) *NativeFnValue {
	return &NativeFnValue{
		Type: ValueNativeFunction,
		call: call,
	}
}
