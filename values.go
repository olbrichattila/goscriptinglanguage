package main

type ValueType int

const (
	ValueTypeNull ValueType = iota
	ValueTypeNumber
	ValueBoolean
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

func MK_NUMBER(n float64) *NumberVal {
	return &NumberVal{Type: ValueTypeNumber, Value: n}
}

func MK_NULL() *NullVal {
	return &NullVal{Type: ValueTypeNull, Value: "null"}
}

func MK_BOOL(v bool) *BoolVal {
	return &BoolVal{Type: ValueBoolean, Value: v}
}
