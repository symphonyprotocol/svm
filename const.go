package svm

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

const (
	LuaFuncOK = iota
	LuaFuncYield
	LuaFuncErrRun
	LuaFuncErrSyntax
	LuaFuncErrMem
	LuaFuncErrGcmm
	LuaFuncErrErr
	LuaFuncErrFile
)

const (
	luaMinStack             = 20
	luaMaxStack             = 1000000
	luaRegisteryIndex       = -luaMaxStack - 1000
	luaRidxGlobals    int64 = 2
)
