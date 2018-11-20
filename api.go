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
	return convertToFloat(val)
}

//ToInteger covert idx value to int64
func (ls *LuaState) ToInteger(idx int) int64 {
	i, _ := ls.ToIntegerX(idx)
	return i
}

//ToIntegerX try to covert idx value to int64
func (ls *LuaState) ToIntegerX(idx int) (int64, bool) {
	val := ls.stack.get(idx)
	return converToInteger(val)
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

//Arith airthmetic method
func (ls *LuaState) Arith(op AirthOp) {
	var a, b luaValue
	b = ls.stack.pop()
	if op != operatorUnm && op != operatorBinNot {
		a = ls.stack.pop()
	} else {
		a = b
	}

	operator := operators[op]
	if result := doAirth(a, b, operator); result != nil {
		ls.stack.push(result)
	} else {
		panic("arithmetic error!")
	}
}

//Compare  compare two index value in stack
func (ls *LuaState) Compare(idx1, idx2 int, op CompareOp) bool {
	a := ls.stack.get(idx1)
	b := ls.stack.get(idx2)
	switch op {
	case operatorEqual:
		return doEqual(a, b)
	case operatorLessThan:
		return doLessThan(a, b)
	case operatorLessEqual:
		return doLessEqual(a, b)
	default:
		panic("invalid compare op!")
	}
}

//Len get length of index value and push length to the stack top
func (ls *LuaState) Len(idx int) {
	val := ls.stack.get(idx)
	if s, ok := val.(string); ok {
		ls.stack.push(int64(len(s)))
	} else if t, ok := val.(*luaTable); ok {
		ls.stack.push(int64(t.len()))
	} else {
		panic("length error!")
	}
}

//Concat pop top n value and concat them push back
func (ls *LuaState) Concat(n int) {
	if n == 0 {
		ls.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if ls.IsString(-1) && ls.IsString(-2) {
				strb := ls.ToString(-1)
				stra := ls.ToString(-2)
				ls.stack.pop()
				ls.stack.pop()
				ls.stack.push(stra + strb)
				continue
			}
			panic("concatenation error!")
		}
	}
}

//NewTable create luaTable
func (ls *LuaState) NewTable(arrayLen, mapCap int) {
	t := newLuaTable(arrayLen, mapCap)
	ls.stack.push(t)
}

//GetTable get luaTable from idx, push value into stack (key: top value of stack)
func (ls *LuaState) GetTable(idx int) LuaType {
	t := ls.stack.get(idx)
	key := ls.stack.pop()
	return ls.getTable(t, key)
}

func (ls *LuaState) getTable(t, k luaValue) LuaType {
	if table, ok := t.(*luaTable); ok {
		val := table.get(k)
		ls.stack.push(val)
		return typeOf(val)
	}
	panic("not a table")
}

//SetTable set table value
func (ls *LuaState) SetTable(idx int) {
	t := ls.stack.get(idx)
	v := ls.stack.pop()
	k := ls.stack.pop()
	ls.setTable(t, k, v)
}

func (ls *LuaState) setTable(t, k, v luaValue) {
	if table, ok := t.(*luaTable); ok {
		table.set(k, v)
		return
	}
	panic("not a table")
}

//SetI set key:i val:pop
func (ls *LuaState) SetI(idx int, i int64) {
	t := ls.stack.get(idx)
	v := ls.stack.pop()
	ls.setTable(t, i, v)
}

func doAirth(a, b luaValue, op operator) luaValue {
	if op.floatFunc == nil {
		if x, ok := converToInteger(a); ok {
			if y, ok := converToInteger(b); ok {
				return op.integerFunc(x, y)
			}
		}
	} else {
		if op.integerFunc != nil {
			if x, ok := a.(int64); ok {
				if y, ok := b.(int64); ok {
					return op.integerFunc(x, y)
				}
			}
		}
		if x, ok := convertToFloat(a); ok {
			if y, ok := convertToFloat(b); ok {
				return op.floatFunc(x, y)
			}
		}
	}
	return nil
}

func doEqual(a, b luaValue) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case int64:
			return x == float64(y)
		default:
			return false
		}
	default:
		return a == b
	}
}

func doLessThan(a, b luaValue) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x < y
		case int64:
			return x < float64(y)
		}
	}
	panic("comparison error!")
}

func doLessEqual(a, b luaValue) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x <= y
		case int64:
			return x <= float64(y)
		}
	}
	panic("comparison error!")
}
