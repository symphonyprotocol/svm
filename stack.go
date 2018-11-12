package svm

type luaValue interface{}

func typeOf(val luaValue) LuaType {
	switch val.(type) {
	case nil:
		return LuaTypeNil
	case bool:
		return LuaTypeBoolean
	case int64:
		return LuaTypeNumber
	case float64:
		return LuaTypeNumber
	case string:
		return LuaTypeString
	default:
		panic("todo!")
	}
}

type luaStack struct {
	data []luaValue
}

func newLuaStack() *luaStack {
	return &luaStack{
		data: make([]luaValue, 0, 1024),
	}
}

func (ls *luaStack) topIndex() int {
	return len(ls.data) - 1
}

func (ls *luaStack) push(val luaValue) {
	ls.data = append(ls.data, val)
}

func (ls *luaStack) pop() luaValue {
	idx := ls.topIndex()
	val := ls.data[idx]
	ls.data = ls.data[:idx]
	return val
}

func (ls *luaStack) absIndex(idx int) int {
	if idx < 0 {
		top := ls.topIndex()
		idx = top + idx + 1
	}
	return idx
}

func (ls *luaStack) get(idx int) luaValue {
	tmpIdx := ls.absIndex(idx)
	if tmpIdx < 0 {
		return nil
	}
	top := ls.topIndex()
	if tmpIdx > top {
		return nil
	}
	return ls.data[tmpIdx]
}

func (ls *luaStack) set(idx int, val luaValue) {
	tmpIdx := ls.absIndex(idx)
	top := ls.topIndex()
	if tmpIdx > top || tmpIdx < 0 {
		return
	}
	ls.data[idx] = val
}

func (ls *luaStack) reverse(from, to int) {
	data := ls.data
	for from < to {
		data[from], data[to] = data[to], data[from]
		from++
		to--
	}
}
