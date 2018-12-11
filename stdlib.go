package svm

import (
	"fmt"
)

type RegisterFunctions map[string]GoFunction

var baseFunctions = map[string]GoFunction{
	"print":        basePrint,
	"error":        baseErrorf,
	"pairs":        basePairs,
	"ipairs":       baseIPairs,
	"next":         baseNext,
	"pcall":        basePCall,
	"getmetatable": baseGetMetatable,
	"setmetatable": baseSetMetatable,
	"select":       baseSelect,
	"assert":       baseAssert,
	"type":         baseType,
}

func basePrint(ls *LuaState) int {
	nArgs := ls.GetTopIndex()
	for i := 0; i <= nArgs; i++ {
		if ls.IsBoolean(i) {
			fmt.Printf("%t", ls.ToBoolean(i))
		} else if ls.IsString(i) {
			fmt.Print(ls.ToString(i))
		} else {
			fmt.Print(ls.TypeName(ls.Type(i)))
		}
		if i < nArgs {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
}

func baseGetMetatable(ls *LuaState) int {
	if !ls.GetMetaTable(0) {
		ls.PushNil()
	}
	return 1
}

func baseSetMetatable(ls *LuaState) int {
	ls.SetMetatable(0)
	return 1
}

func baseNext(ls *LuaState) int {
	ls.SetTop(1)
	if ls.Next(0) {
		return 2
	}
	ls.PushNil()
	return 1
}

func basePairs(ls *LuaState) int {
	ls.PushGoFunction(baseNext)
	ls.PushValue(0)
	ls.PushNil()
	return 3
}

func iPairsAux(ls *LuaState) int {
	i := ls.ToInteger(1) + 1
	ls.PushInteger(i)
	if ls.GetI(0, i) == LuaTypeNil {
		return 1
	}
	return 2
}

func baseIPairs(ls *LuaState) int {
	ls.PushGoFunction(iPairsAux)
	ls.PushValue(0)
	ls.PushInteger(0)
	return 3
}

func baseErrorf(ls *LuaState) int {
	return ls.Error()
}

func basePCall(ls *LuaState) int {
	nArgs := ls.GetTopIndex()
	status := ls.PCall(nArgs, -1, 0)
	ls.PushBoolean(status == LuaFuncOK)
	ls.Insert(0)
	return ls.GetTopIndex() + 1
}

func baseSelect(ls *LuaState) int {
	n := int64(ls.GetTopIndex())
	if ls.Type(0) == LuaTypeString && ls.CheckString(0) == "#" {
		ls.PushInteger(n)
		return 1
	}
	i := ls.CheckInteger(0)
	if i < 0 {
		i = n + i
	} else if i > n {
		i = n
	}
	ls.ArgCheck(1 <= i, 1, "index out of range")
	return int(n - i + 1)
}

func baseAssert(ls *LuaState) int {
	if ls.ToBoolean(0) { /* condition is true? */
		return ls.GetTopIndex() /* return all arguments */
	} else { /* error */
		ls.CheckAny(0)                     /* there must be a condition */
		ls.Remove(0)                       /* remove it */
		ls.PushString("assertion failed!") /* default message */
		ls.SetTop(0)                       /* leave only message (default if no other one) */
		return baseErrorf(ls)              /* call 'error' */
	}
}

func baseType(ls *LuaState) int {
	typeStr := ls.TypeName2(0)
	ls.PushString(typeStr)
	return 1
}
