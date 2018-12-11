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
	case *luaTable:
		return LuaTypeTable
	case *luaClosure:
		return LuaTypeFunction
	default:
		panic("todo!")
	}
}

func typeName(tp LuaType) string {
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

var stackSeq uint

type luaStack struct {
	id      uint
	data    []luaValue
	prev    *luaStack
	closure *luaClosure
	varargs []luaValue
	pc      int
	isDebug bool
	state   *LuaState
	openuvs map[int]*luaUpvalue
}

func newLuaStack(state *LuaState) *luaStack {
	stackSeq++
	return &luaStack{
		id:      stackSeq,
		data:    make([]luaValue, 0, 1024),
		isDebug: false,
		state:   state,
	}
}

func (ls *luaStack) setDebug(debug bool) {
	ls.isDebug = debug
}

func (ls *luaStack) topIndex() int {
	return len(ls.data) - 1
}

func (ls *luaStack) push(val luaValue) {
	ls.data = append(ls.data, val)

	if ls.isDebug {
		printStack(ls)
	}
}

func (ls *luaStack) pushN(vals []luaValue, n int) {
	nVals := len(vals)
	if n < 0 {
		n = nVals
	}

	for i := 0; i < n; i++ {
		if i < nVals {
			ls.push(vals[i])
		} else {
			ls.push(nil)
		}
	}
}

func (ls *luaStack) pop() luaValue {
	idx := ls.topIndex()
	val := ls.data[idx]
	ls.data = ls.data[:idx]

	if ls.isDebug {
		printStack(ls)
	}

	return val
}

func (ls *luaStack) popN(n int) []luaValue {
	if n < 0 {
		n = -n
	}
	vals := make([]luaValue, n)
	for i := n - 1; i >= 0; i-- {
		val := ls.pop()
		vals[i] = val
	}
	return vals
}

func (ls *luaStack) absIndex(idx int) int {
	if idx <= luaRegisteryIndex {
		return idx
	}
	if idx < 0 {
		top := ls.topIndex()
		idx = top + idx + 1
	}
	return idx
}

func (ls *luaStack) get(idx int) luaValue {
	if idx < luaRegisteryIndex {
		uvIdx := luaRegisteryIndex - idx - 1
		if ls.closure == nil || uvIdx >= len(ls.closure.upvals) {
			return nil
		}
		return *(ls.closure.upvals[uvIdx].val)
	}
	if idx == luaRegisteryIndex {
		return ls.state.registry
	}
	tmpIdx := ls.absIndex(idx)
	if tmpIdx < 0 {
		return nil
	}
	return ls.data[tmpIdx]
}

func (ls *luaStack) set(idx int, val luaValue) {
	if idx < luaRegisteryIndex {
		uvIdx := luaRegisteryIndex - idx - 1
		if ls.closure != nil && uvIdx < len(ls.closure.upvals) {
			*(ls.closure.upvals[uvIdx].val) = val
		}
		return
	}
	if idx == luaRegisteryIndex {
		ls.state.registry = val.(*luaTable)
		return
	}
	tmpIdx := ls.absIndex(idx)
	top := ls.topIndex()
	if tmpIdx >= 0 && tmpIdx <= top+1 {
		if top+1 == tmpIdx {
			ls.push(val)
		} else {
			ls.data[tmpIdx] = val
		}

		if ls.isDebug {
			printStack(ls)
		}

		return
	}
	panic("set stack out of range!")
}

func (ls *luaStack) reverse(from, to int) {
	data := ls.data
	for from < to {
		data[from], data[to] = data[to], data[from]
		from++
		to--
	}

	if ls.isDebug {
		printStack(ls)
	}
}

func (ls *luaStack) setTop(idx int) {
	newTop := ls.absIndex(idx)
	top := ls.topIndex()
	if newTop > top {
		n := newTop - top
		for i := 0; i < n; i++ {
			ls.push(nil)
		}
	}
	if newTop < top {
		n := top - newTop
		for i := 0; i < n; i++ {
			ls.pop()
		}
	}
}

func (ls *luaStack) removeNilTail() {
	top := ls.topIndex()
	for ls.get(top) == nil {
		ls.pop()
		top = ls.topIndex()
	}
}

func (ls *luaStack) isValid(idx int) bool {
	if idx < luaRegisteryIndex {
		uvIdx := luaRegisteryIndex - idx - 1
		c := ls.closure
		return c != nil && uvIdx < len(c.upvals)
	}
	if idx == luaRegisteryIndex {
		return true
	}
	tmpIdx := ls.absIndex(idx)
	if tmpIdx < len(ls.data) {
		return true
	}
	return false
}
