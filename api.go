package svm

import (
	"fmt"
)

const (
	LuaTypeNone = iota
	LuaTypeNil
	LuaTypeBoolean
	LuaTypeLightUserData
	LuaTypeNumber
	LuaTypeString
	LuaTypeTable
	LuaTypeFunction
	LuaTypeUserData
	LuaTypeThread
)

//LuaType lua type structure
type LuaType = int

//PushNil push nil value
func (ls *LuaState) PushNil() {
	ls.stack.push(nil)
}

//PushBoolean push boolean value
func (ls *LuaState) PushBoolean(b bool) {
	ls.stack.push(b)
}

//PushInteger push int64
func (ls *LuaState) PushInteger(n int64) {
	ls.stack.push(n)
}

//PushNumber push float64
func (ls *LuaState) PushNumber(n float64) {
	ls.stack.push(n)
}

//PushString push string value
func (ls *LuaState) PushString(s string) {
	ls.stack.push(s)
}

//PushValue get idx value and push it into stack
func (ls *LuaState) PushValue(idx int) {
	val := ls.stack.get(idx)
	ls.stack.push(val)
}

//GetTopIndex get stack top index
func (ls *LuaState) GetTopIndex() int {
	return ls.stack.topIndex()
}

//AbsIndex get abs index
func (ls *LuaState) AbsIndex(idx int) int {
	return ls.stack.absIndex(idx)
}

//Pop pop out n index
func (ls *LuaState) Pop(n int) {
	for i := 0; i < n; i++ {
		ls.stack.pop()
	}
}

//Copy copy fromIdx value to toIdx
func (ls *LuaState) Copy(fromIdx, toIdx int) {
	val := ls.stack.get(fromIdx)
	ls.stack.set(toIdx, val)
}

//Replace pop stack value and insert into idx
func (ls *LuaState) Replace(idx int) {
	val := ls.stack.pop()
	ls.stack.set(idx, val)
}

//Insert pop out top value in stack and insert into idx
func (ls *LuaState) Insert(idx int) {
	ls.Rotate(idx, 1)
}

//Remove delete value for idx, and move value > idx down
func (ls *LuaState) Remove(idx int) {
	ls.Rotate(idx, -1)
	ls.Pop(1)
}

//Rotate revert [idx, top] value to n position
func (ls *LuaState) Rotate(idx, n int) {
	top := ls.stack.topIndex()
	start := ls.stack.absIndex(idx)
	var m int
	if n >= 0 {
		m = top - n
	} else {
		m = start - n - 1
	}
	ls.stack.reverse(start, m)
	ls.stack.reverse(m+1, top)
	ls.stack.reverse(start, top)
}

//SetTop set idx to the new top of the stack
func (ls *LuaState) SetTop(idx int) {
	top := ls.stack.absIndex(idx)
	if top < 0 {
		panic("stack underflow!")
	}
	n := ls.stack.topIndex() - top
	if n > 0 {
		for i := 0; i < n; i++ {
			ls.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			ls.stack.push(nil)
		}
	}
}

//TypeName return lua type name
func (ls *LuaState) TypeName(tp LuaType) string {
	switch tp {
	case LuaTypeNone:
		return "no value"
	case LuaTypeNil:
		return "nil"
	case LuaTypeBoolean:
		return "boolean"
	case LuaTypeNumber:
		return "number"
	case LuaTypeString:
		return "string"
	case LuaTypeTable:
		return "table"
	case LuaTypeFunction:
		return "function"
	case LuaTypeThread:
		return "thread"
	default:
		return "userdata"
	}
}

//Type return index value's lua type
func (ls *LuaState) Type(idx int) LuaType {
	val := ls.stack.get(idx)
	return typeOf(val)
}

//IsNone is idx value == None
func (ls *LuaState) IsNone(idx int) bool {
	return ls.Type(idx) == LuaTypeNone
}

//IsNil is idx value == nil
func (ls *LuaState) IsNil(idx int) bool {
	return ls.Type(idx) == LuaTypeNil
}

//IsNoneOrNil is idx value == None or nil
func (ls *LuaState) IsNoneOrNil(idx int) bool {
	return ls.Type(idx) <= LuaTypeNil
}

//IsBoolean is idx value == bool
func (ls *LuaState) IsBoolean(idx int) bool {
	return ls.Type(idx) == LuaTypeBoolean
}

//IsString is idx value == string
func (ls *LuaState) IsString(idx int) bool {
	t := ls.Type(idx)
	return t == LuaTypeString || t == LuaTypeNumber
}

//IsNumber is idx value == number
func (ls *LuaState) IsNumber(idx int) bool {
	_, ok := ls.ToNumberX(idx)
	return ok
}

//IsInteger is idx value == integer
func (ls *LuaState) IsInteger(idx int) bool {
	val := ls.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

//ToBoolean idx value covert to boolean
func (ls *LuaState) ToBoolean(idx int) bool {
	val := ls.stack.get(idx)
	return covertToBoolean(val)
}

func covertToBoolean(val luaValue) bool {
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

//ToNumber idx value covert to float64
func (ls *LuaState) ToNumber(idx int) float64 {
	n, _ := ls.ToNumberX(idx)
	return n
}

//ToNumberX idx value try convert to float64
func (ls *LuaState) ToNumberX(idx int) (float64, bool) {
	val := ls.stack.get(idx)
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

//ToInteger covert idx value to int64
func (ls *LuaState) ToInteger(idx int) int64 {
	i, _ := ls.ToIntegerX(idx)
	return i
}

//ToIntegerX try to covert idx value to int64
func (ls *LuaState) ToIntegerX(idx int) (int64, bool) {
	val := ls.stack.get(idx)
	i, ok := val.(int64)
	return i, ok
}

//ToString covert idx value to string
func (ls *LuaState) ToString(idx int) string {
	s, _ := ls.ToStringX(idx)
	return s
}

//ToStringX  try covert idx value to string and change idx value to string if it's not
func (ls *LuaState) ToStringX(idx int) (string, bool) {
	val := ls.stack.get(idx)
	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x)
		ls.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}
