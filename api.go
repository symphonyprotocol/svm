package svm

import (
	"bytes"
	"fmt"
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
	ls.stack.setTop(idx)
}

//RemoteNilTail remove nil tail before call
func (ls *LuaState) RemoveNilTail() {
	ls.stack.removeNilTail()
}

//TypeName return lua type name
func (ls *LuaState) TypeName(tp LuaType) string {
	return typeName(tp)
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
		return
	}
	mm := operator.metamethod
	if result, ok := callMetamethod(a, b, mm, ls); ok {
		ls.stack.push(result)
		return
	}
	panic("airthmetic error!")
}

//Compare  compare two index value in stack
func (ls *LuaState) Compare(idx1, idx2 int, op CompareOp) bool {
	a := ls.stack.get(idx1)
	b := ls.stack.get(idx2)
	switch op {
	case operatorEqual:
		return doEqual(a, b, ls)
	case operatorLessThan:
		return doLessThan(a, b, ls)
	case operatorLessEqual:
		return doLessEqual(a, b, ls)
	default:
		panic("invalid compare op!")
	}
}

//Len get length of index value and push length to the stack top
func (ls *LuaState) Len(idx int) {
	val := ls.stack.get(idx)
	if s, ok := val.(string); ok {
		ls.stack.push(int64(len(s)))
	} else if result, ok := callMetamethod(val, val, "__len", ls); ok {
		ls.stack.push(result)
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
			b := ls.stack.pop()
			a := ls.stack.pop()
			if result, ok := callMetamethod(a, b, "__concat", ls); ok {
				ls.stack.push(result)
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
	return ls.getTable(t, key, false)
}

func (ls *LuaState) getTable(t, k luaValue, raw bool) LuaType {
	if table, ok := t.(*luaTable); ok {
		val := table.get(k)
		if raw || val != nil || !table.hasMetafield("__index") {
			ls.stack.push(val)
			return typeOf(val)
		}
	}
	if !raw {
		if mf := getMetafield(t, "__index", ls); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				return ls.getTable(x, k, false)
			case *luaClosure:
				ls.stack.push(mf)
				ls.stack.push(t)
				ls.stack.push(k)
				ls.Call(2, 1)
				v := ls.stack.get(-1)
				return typeOf(v)
			}
		}
	}
	panic("not a table")
}

//SetTable set table value
func (ls *LuaState) SetTable(idx int) {
	t := ls.stack.get(idx)
	v := ls.stack.pop()
	k := ls.stack.pop()
	ls.setTable(t, k, v, false)
}

func (ls *LuaState) setTable(t, k, v luaValue, raw bool) {
	if table, ok := t.(*luaTable); ok {
		if raw || table.get(k) != nil || !table.hasMetafield("__newindex") {
			table.set(k, v)
			return
		}
	}
	if !raw {
		if mf := getMetafield(t, "__newindex", ls); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				ls.setTable(x, k, v, false)
				return
			case *luaClosure:
				ls.stack.push(mf)
				ls.stack.push(t)
				ls.stack.push(k)
				ls.stack.push(v)
				ls.Call(3, 0)
				return
			}
		}
	}
	panic("not a table")
}

//SetField set table key value
func (ls *LuaState) SetField(idx int, k string) {
	t := ls.stack.get(idx)
	v := ls.stack.pop()
	ls.setTable(t, k, v, false)
}

//SetI set key:i val:pop
func (ls *LuaState) SetI(idx int, i int64) {
	t := ls.stack.get(idx)
	v := ls.stack.pop()
	ls.setTable(t, i, v, false)
}

//GetI get key:i from table
func (ls *LuaState) GetI(idx int, i int64) LuaType {
	t := ls.stack.get(idx)
	return ls.getTable(t, i, false)
}

//Load load binary code to the stack
func (ls *LuaState) Load(chunk []byte, chunkName, mode string) int {
	reader := bytes.NewReader(chunk)
	proto := Undump(reader)
	c := newLuaClosure(proto)
	ls.stack.push(c)
	if len(proto.Upvalues) > 0 {
		env := ls.registry.get(luaRidxGlobals)
		c.upvals[0] = &luaUpvalue{&env}
	}
	return 0
}

//LoadProto load proto
func (ls *LuaState) LoadProto(idx int) {
	proto := ls.stack.closure.proto.Protos[idx]
	c := newLuaClosure(proto)
	ls.stack.push(c)
	for i, uvInfo := range proto.Upvalues {
		uvIdx := int(uvInfo.Idx)
		if uvInfo.InStack == 1 {
			if ls.stack.openuvs == nil {
				ls.stack.openuvs = make(map[int]*luaUpvalue)
			}
			if openuv, found := ls.stack.openuvs[uvIdx]; found {
				c.upvals[i] = openuv
			} else {
				c.upvals[i] = &luaUpvalue{&ls.stack.data[uvIdx]}
				ls.stack.openuvs[uvIdx] = c.upvals[i]
			}
		} else {
			c.upvals[i] = ls.stack.closure.upvals[uvIdx]
		}
	}
}

//Call call function in stack
func (ls *LuaState) Call(nArgs, nResults int) {
	value := ls.stack.get(-(nArgs + 1))
	c, ok := value.(*luaClosure)
	if !ok {
		if mf := getMetafield(value, "__call", ls); mf != nil {
			if c, ok = mf.(*luaClosure); ok {
				ls.stack.push(value)
				ls.Insert(-(nArgs + 2))
				nArgs++
			}
		}
	}
	if ok {
		//fmt.Printf("call %s<%d, %d>\n", c.proto.Source, c.proto.LineDefined, c.proto.LastLineDefined)
		if c.proto != nil {
			ls.callLuaClosure(nArgs, nResults, c)
		} else {
			ls.callGoClosure(nArgs, nResults, c)
		}
	} else {
		panic("not a function!")
	}
}

func (ls *LuaState) callLuaClosure(nArgs, nResults int, c *luaClosure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParam)
	isVararg := c.proto.IsVararg == 1
	newStack := newLuaStack(ls)
	newStack.closure = c
	newStack.setDebug(ls.isDebug)
	funcAndArgs := ls.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.setTop(nRegs - 1)
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}
	ls.pushLuaStack(newStack)
	ls.runLuaClosure()
	ls.popLuaStack()
	if nResults != 0 {
		results := newStack.popN(newStack.topIndex() - nRegs + 1)
		ls.stack.pushN(results, nResults)
	}
}

func (ls *LuaState) runLuaClosure() {
	for {
		inst := Instruction(ls.Fetch())
		inst.Execute(ls)
		if inst.Opcode() == opReturn {
			break
		}
	}
}

func (ls *LuaState) callGoClosure(nArgs, nResults int, c *luaClosure) {
	newStack := newLuaStack(ls)
	newStack.closure = c
	newStack.setDebug(ls.isDebug)
	args := ls.stack.popN(nArgs)
	newStack.pushN(args, nArgs)
	ls.stack.pop()
	ls.pushLuaStack(newStack)
	r := c.goFunc(ls)
	ls.popLuaStack()
	if nResults != 0 {
		results := newStack.popN(r)
		ls.stack.pushN(results, nResults)
	}
}

func (ls *LuaState) pushLuaStack(stack *luaStack) {
	stack.prev = ls.stack
	ls.stack = stack
}

func (ls *LuaState) popLuaStack() {
	stack := ls.stack
	ls.stack = stack.prev
	stack.prev = nil
}

func (ls *LuaState) pushGoFunction(f GoFunction) {
	ls.stack.push(newGoClosure(f, 0))
}

//PushGoFunction push go function to the stack
func (ls *LuaState) PushGoFunction(f GoFunction) {
	ls.stack.push(newGoClosure(f, 0))
}

func (ls *LuaState) isGoFunction(idx int) bool {
	val := ls.stack.get(idx)
	if c, ok := val.(*luaClosure); ok {
		return c.goFunc != nil
	}
	return false
}

func (ls *LuaState) toGoFunction(idx int) GoFunction {
	val := ls.stack.get(idx)
	if c, ok := val.(*luaClosure); ok {
		return c.goFunc
	}
	return nil
}

func (ls *LuaState) pushGoClosure(f GoFunction, n int) {
	c := newGoClosure(f, n)
	for i := n; i > 0; i-- {
		val := ls.stack.pop()
		c.upvals[n-1] = &luaUpvalue{&val}
	}
	ls.stack.push(c)
}

func (ls *LuaState) pushGlobalTable() {
	global := ls.registry.get(luaRidxGlobals)
	ls.stack.push(global)
}

func (ls *LuaState) getGlobal(name string) LuaType {
	t := ls.registry.get(luaRidxGlobals)
	return ls.getTable(t, name, false)
}

func (ls *LuaState) setGlobal(name string) {
	t := ls.registry.get(luaRidxGlobals)
	v := ls.stack.pop()
	ls.setTable(t, name, v, false)
}

//Register regisger go function to lua
func (ls *LuaState) Register(name string, f GoFunction) {
	ls.pushGoFunction(f)
	ls.setGlobal(name)
}

func (ls *LuaState) closeUpvalues(a int) {
	for i, openuv := range ls.stack.openuvs {
		if i >= a-1 {
			val := *openuv.val
			openuv.val = &val
			delete(ls.stack.openuvs, i)
		}
	}
}

//GetMetaTable stack getMetatable
func (ls *LuaState) GetMetaTable(idx int) bool {
	val := ls.stack.get(idx)
	if mt := getMetatable(val, ls); mt != nil {
		ls.stack.push(mt)
		return true
	}
	return false
}

//SetMetatable stack seMetatable
func (ls *LuaState) SetMetatable(idx int) {
	val := ls.stack.get(idx)
	mtVal := ls.stack.pop()
	if mtVal == nil {
		setMetatable(val, nil, ls)
	} else if mt, ok := mtVal.(*luaTable); ok {
		setMetatable(val, mt, ls)
	} else {
		panic("table expected")
	}
}

//RawEqual metatable equal
func (ls *LuaState) RawEqual(idx1, idx2 int) bool {
	if !ls.stack.isValid(idx1) || !ls.stack.isValid(idx2) {
		return false
	}

	a := ls.stack.get(idx1)
	b := ls.stack.get(idx2)
	return doEqual(a, b, nil)
}

//RawLen metatable len
func (ls *LuaState) RawLen(idx int) uint {
	val := ls.stack.get(idx)
	switch x := val.(type) {
	case string:
		return uint(len(x))
	case *luaTable:
		return uint(x.len())
	default:
		return 0
	}
}

//RawGet metatable get
func (ls *LuaState) RawGet(idx int) LuaType {
	t := ls.stack.get(idx)
	k := ls.stack.pop()
	return ls.getTable(t, k, true)
}

//RawGetI metatable geti
func (ls *LuaState) RawGetI(idx int, i int64) LuaType {
	t := ls.stack.get(idx)
	return ls.getTable(t, i, true)
}

//RawSet metatable set
func (ls *LuaState) RawSet(idx int) {
	t := ls.stack.get(idx)
	v := ls.stack.pop()
	k := ls.stack.pop()
	ls.setTable(t, k, v, true)
}

//RawSetI metatable seti
func (ls *LuaState) RawSetI(idx int, i int64) {
	t := ls.stack.get(idx)
	v := ls.stack.pop()
	ls.setTable(t, i, v, true)
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

//Next next iterator
func (ls *LuaState) Next(idx int) bool {
	val := ls.stack.get(idx)
	if t, ok := val.(*luaTable); ok {
		key := ls.stack.pop()
		if nexKey := t.nextKey(key); nexKey != nil {
			ls.stack.push(nexKey)
			ls.stack.push(t.get(nexKey))
			return true
		}
		return false
	}
	panic("table expected")
}

//Error raise error
func (ls *LuaState) Error() int {
	err := ls.stack.pop()
	panic(err)
}

//PCall pcall for error exception
func (ls *LuaState) PCall(nArgs, nResults, msgh int) (status int) {
	caller := ls.stack
	status = LuaFuncErrRun
	defer func() {
		if err := recover(); err != nil {
			if msgh != 0 {
				panic(err)
			}
			for ls.stack != caller {
				ls.popLuaStack()
			}
			ls.stack.push(err)
		}
	}()
	ls.Call(nArgs, nResults)
	status = LuaFuncOK
	return
}

func doEqual(a, b luaValue, ls *LuaState) bool {
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
	case *luaTable:
		if y, ok := b.(*luaTable); ok && x != y && ls != nil {
			if result, ok := callMetamethod(x, y, "__eq", ls); ok {
				return covertToBoolean(result)
			}
		}
		return a == b
	default:
		return a == b
	}
}

func doLessThan(a, b luaValue, ls *LuaState) bool {
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
	if result, ok := callMetamethod(a, b, "__lt", ls); ok {
		return covertToBoolean(result)
	}
	panic("comparison error!")
}

func doLessEqual(a, b luaValue, ls *LuaState) bool {
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
	if result, ok := callMetamethod(a, b, "__le", ls); ok {
		return covertToBoolean(result)
	} else if result, ok := callMetamethod(b, a, "__lt", ls); ok {
		return !covertToBoolean(result)
	}
	panic("comparison error!")
}
