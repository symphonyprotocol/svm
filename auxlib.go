package svm

import "fmt"

func (ls *LuaState) LoadBaseLibs() {
	for k, f := range baseFunctions {
		ls.Register(k, f)
	}
}

func (ls *LuaState) CheckString(arg int) string {
	s, ok := ls.ToStringX(arg)
	if !ok {
		ls.tagError(arg, LuaTypeString)
	}
	return s
}

func (ls *LuaState) GetMetafield(obj int, event string) LuaType {
	if !ls.GetMetaTable(obj) { /* no metatable? */
		return LuaTypeNil
	}

	ls.PushString(event)
	tt := ls.RawGet(-2)
	if tt == LuaTypeNil { /* is metafield nil? */
		ls.Pop(2) /* remove metatable and metafield */
	} else {
		ls.Remove(-2) /* remove only metatable */
	}
	return tt /* return metafield type */
}

func (ls *LuaState) CheckInteger(arg int) int64 {
	i, ok := ls.ToIntegerX(arg)
	if !ok {
		ls.intError(arg)
	}
	return i
}

func (ls *LuaState) ArgCheck(cond bool, arg int, extraMsg string) {
	if !cond {
		ls.ArgError(arg, extraMsg)
	}
}

func (ls *LuaState) CheckAny(arg int) {
	if ls.Type(arg) == LuaTypeNone {
		ls.ArgError(arg, "value expected")
	}
}

func (ls *LuaState) TypeName2(idx int) string {
	return ls.TypeName(ls.Type(idx))
}

func (ls *LuaState) PushFString(fmtStr string, a ...interface{}) {
	str := fmt.Sprintf(fmtStr, a...)
	ls.stack.push(str)
}

func (ls *LuaState) Error2(fmt string, a ...interface{}) int {
	ls.PushFString(fmt, a...) // todo
	return ls.Error()
}

func (ls *LuaState) ArgError(arg int, extraMsg string) int {
	// bad argument #arg to 'funcname' (extramsg)
	return ls.Error2("bad argument #%d (%s)", arg, extraMsg) // todo
}

func (ls *LuaState) typeError(arg int, tname string) int {
	var typeArg string /* name for the type of the actual argument */
	if ls.GetMetafield(arg, "__name") == LuaTypeString {
		typeArg = ls.ToString(-1) /* use the given type name */
	} else if ls.Type(arg) == LuaTypeLightUserData {
		typeArg = "light userdata" /* special name for messages */
	} else {
		typeArg = ls.TypeName2(arg) /* standard name */
	}
	msg := tname + " expected, got " + typeArg
	ls.PushString(msg)
	return ls.ArgError(arg, msg)
}

func (ls *LuaState) tagError(arg int, tag LuaType) {
	ls.typeError(arg, ls.TypeName(LuaType(tag)))
}

func (ls *LuaState) intError(arg int) {
	if ls.IsNumber(arg) {
		ls.ArgError(arg, "number has no integer representation")
	} else {
		ls.tagError(arg, LuaTypeNumber)
	}
}
