package svm

import (
	"math"
)

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
	integerFunc func(int64, int64) int64
	floatFunc   func(float64, float64) float64
}

var operators = []operator{
	operator{func(a, b int64) int64 { return a + b }, func(a, b float64) float64 { return a + b }},
	operator{func(a, b int64) int64 { return a - b }, func(a, b float64) float64 { return a - b }},
	operator{func(a, b int64) int64 { return a * b }, func(a, b float64) float64 { return a * b }},
	operator{intMod, floatMod},
	operator{nil, math.Pow},
	operator{nil, func(a, b float64) float64 { return a / b }},
	operator{intFloorDiv, floatFloorDiv},
	operator{func(a, b int64) int64 { return a & b }, nil},
	operator{func(a, b int64) int64 { return a | b }, nil},
	operator{func(a, b int64) int64 { return a ^ b }, nil},
	operator{intShiftLeft, nil},
	operator{intShiftRight, nil},
	operator{func(a, _ int64) int64 { return -a }, func(a, _ float64) float64 { return -a }},
	operator{func(a, _ int64) int64 { return ^a }, nil},
}

//AirthOp airth operator
type AirthOp = int

//CompareOp compare operator
type CompareOp = int

//LuaState lua state object
type LuaState struct {
	stack *luaStack
	proto *LuaTrunkProto
	pc    int
}

//NewLuaState create a lua state object
func NewLuaState(proto *LuaTrunkProto) *LuaState {
	return &LuaState{
		stack: newLuaStack(),
		proto: proto,
		pc:    0,
	}
}

//PC get current pc
func (ls *LuaState) PC() int {
	return ls.pc
}

//AddPC add pc for n count
func (ls *LuaState) AddPC(n int) {
	ls.pc += n
}

//Fetch get current pc point to code
func (ls *LuaState) Fetch() uint32 {
	code := ls.proto.Code[ls.pc]
	ls.pc++
	return code
}

//GetConst get constant value at idx and push it into stack
func (ls *LuaState) GetConst(idx int) {
	c := ls.proto.Constants[idx]
	ls.stack.push(c)
}

//GetRK get constant or register address and push it into stack
func (ls *LuaState) GetRK(rk int) {
	if rk > 0xFF { //Constant
		ls.GetConst(rk & 0xFF)
	} else { //register
		ls.PushValue(rk)
	}
}
