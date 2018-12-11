package svm

import "math"

const (
	operatorAdd = iota
	operatorSub
	operatorMul
	operatorMod
	operatorPow
	operatorDiv
	operatorIDiv
	operatorBinAnd
	operatorBinOr
	operatorBinXor
	operatorShl
	operatorShr
	operatorUnm
	operatorBinNot
)

const (
	operatorEqual = iota
	operatorLessThan
	operatorLessEqual
)

type operator struct {
	metamethod  string
	integerFunc func(int64, int64) int64
	floatFunc   func(float64, float64) float64
}

var operators = []operator{
	operator{"__add", func(a, b int64) int64 { return a + b }, func(a, b float64) float64 { return a + b }},
	operator{"__sub", func(a, b int64) int64 { return a - b }, func(a, b float64) float64 { return a - b }},
	operator{"__mul", func(a, b int64) int64 { return a * b }, func(a, b float64) float64 { return a * b }},
	operator{"__mod", intMod, floatMod},
	operator{"__pow", nil, math.Pow},
	operator{"__div", nil, func(a, b float64) float64 { return a / b }},
	operator{"__idiv", intFloorDiv, floatFloorDiv},
	operator{"__band", func(a, b int64) int64 { return a & b }, nil},
	operator{"__bor", func(a, b int64) int64 { return a | b }, nil},
	operator{"__bxor", func(a, b int64) int64 { return a ^ b }, nil},
	operator{"__shl", intShiftLeft, nil},
	operator{"__shr", intShiftRight, nil},
	operator{"__unm", func(a, _ int64) int64 { return -a }, func(a, _ float64) float64 { return -a }},
	operator{"__bnot", func(a, _ int64) int64 { return ^a }, nil},
}
