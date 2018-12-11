package svm

import (
	"fmt"
)

//GoFunction go function used in lua
type GoFunction func(*LuaState) int

//AirthOp airth operator
type AirthOp = int

//CompareOp compare operator
type CompareOp = int

//LuaState lua state object
type LuaState struct {
	registry *luaTable
	stack    *luaStack
	isDebug  bool
}

//NewLuaState create a lua state object
func NewLuaState() *LuaState {
	registery := newLuaTable(0, 0)
	registery.set(luaRidxGlobals, newLuaTable(0, 0))
	ls := &LuaState{registry: registery}
	ls.pushLuaStack(newLuaStack(ls))
	ls.isDebug = false
	return ls
}

//SetDebug set debug flag
func (ls *LuaState) SetDebug(debug bool) {
	ls.isDebug = debug
	ls.stack.setDebug(debug)
}

//RegisterCount get max stack size
func (ls *LuaState) RegisterCount() int {
	if ls.isDebug {
		fmt.Println("maxStackSize:", ls.stack.closure.proto.MaxStackSize)
	}
	return int(ls.stack.closure.proto.MaxStackSize)
}

//LoadVararg load varargs
func (ls *LuaState) LoadVararg(n int) {
	if n < 0 {
		n = len(ls.stack.varargs)
	}
	ls.stack.pushN(ls.stack.varargs, n)
}

//PC get current pc
func (ls *LuaState) PC() int {
	return ls.stack.pc
}

//AddPC add pc for n count
func (ls *LuaState) AddPC(n int) {
	ls.stack.pc += n
}

//Fetch get current pc point to code
func (ls *LuaState) Fetch() uint32 {
	code := ls.stack.closure.proto.Code[ls.stack.pc]
	ls.stack.pc++
	return code
}

//GetConst get constant value at idx and push it into stack
func (ls *LuaState) GetConst(idx int) {
	c := ls.stack.closure.proto.Constants[idx]
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
