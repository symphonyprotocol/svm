package svm

const fieldsPerFlush = 50

func moveOP(i Instruction, ls *LuaState) {
	a, b, _ := i.ABC()
	ls.Copy(b, a)
}

func jmpOP(i Instruction, ls *LuaState) {
	a, sBx := i.AsBx()
	ls.AddPC(sBx)
	if a != 0 {
		ls.closeUpvalues(a)
	}
}

func loadNilOP(i Instruction, ls *LuaState) {
	a, b, _ := i.ABC()
	ls.PushNil()
	for i := a; i <= a+b; i++ {
		ls.Copy(-1, i)
	}
	ls.Pop(1)
}

func loadBoolOP(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	ls.PushBoolean(b != 0)
	ls.Replace(a)
	if c != 0 {
		ls.AddPC(1)
	}
}

func loadKOP(i Instruction, ls *LuaState) {
	a, bx := i.ABx()
	ls.GetConst(bx)
	ls.Replace(a)
}

func loadKxOP(i Instruction, ls *LuaState) {
	a, _ := i.ABx()
	ax := Instruction(ls.Fetch()).Ax()
	ls.GetConst(ax)
	ls.Replace(a)
}

func lenOP(i Instruction, ls *LuaState) {
	a, b, _ := i.ABC()
	ls.Len(b)
	ls.Replace(a)
}

func concatOP(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	n := c - b + 1
	for i := b; i <= c; i++ {
		ls.PushValue(i)
	}
	ls.Concat(n)
	ls.Replace(a)
}

func testSetOP(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	if ls.ToBoolean(b) == (c != 0) {
		ls.Copy(b, a)
	} else {
		ls.AddPC(1)
	}
}

func testOP(i Instruction, ls *LuaState) {
	a, _, c := i.ABC()
	if ls.ToBoolean(a) != (c != 0) {
		ls.AddPC(1)
	}
}

func forPrepOP(i Instruction, ls *LuaState) {
	a, sBx := i.AsBx()

	ls.PushValue(a)
	ls.PushValue(a + 2)
	ls.Arith(operatorSub)
	ls.Replace(a)
	ls.AddPC(sBx)
}

func forLoopOP(i Instruction, ls *LuaState) {
	a, sBx := i.AsBx()
	ls.PushValue(a + 2)
	ls.PushValue(a)
	ls.Arith(operatorAdd)
	ls.Replace(a)
	isPositiveStep := ls.ToNumber(a+2) >= 0
	if isPositiveStep && ls.Compare(a, a+1, operatorLessEqual) ||
		!isPositiveStep && ls.Compare(a+1, a, operatorLessEqual) {
		ls.AddPC(sBx)
		ls.Copy(a, a+3)
	}
}

func compareArith(i Instruction, ls *LuaState, op CompareOp) {
	a, b, c := i.ABC()
	ls.GetRK(b)
	ls.GetRK(c)
	if ls.Compare(-2, -1, op) != (a != 0) {
		ls.AddPC(1)
	}
	ls.Pop(2)
}

func equalOP(i Instruction, ls *LuaState) {
	compareArith(i, ls, operatorEqual)
}

func lessThanOP(i Instruction, ls *LuaState) {
	compareArith(i, ls, operatorLessThan)
}

func lessEqualOP(i Instruction, ls *LuaState) {
	compareArith(i, ls, operatorLessEqual)
}

func notOP(i Instruction, ls *LuaState) {
	a, b, _ := i.ABC()
	ls.PushBoolean(!ls.ToBoolean(b))
	ls.Replace(a)
}

func binaryArith(i Instruction, ls *LuaState, op AirthOp) {
	a, b, c := i.ABC()
	ls.GetRK(b)
	ls.GetRK(c)
	ls.Arith(op)
	ls.Replace(a)
}

func unaryArith(i Instruction, ls *LuaState, op AirthOp) {
	a, b, _ := i.ABC()
	ls.PushValue(b)
	ls.Arith(op)
	ls.Replace(a)
}

func addOP(i Instruction, ls *LuaState) {
	binaryArith(i, ls, operatorAdd)
}

func subOP(i Instruction, ls *LuaState) {
	binaryArith(i, ls, operatorSub)
}

func mulOP(i Instruction, ls *LuaState) {
	binaryArith(i, ls, operatorMul)
}

func modOP(i Instruction, ls *LuaState) {
	binaryArith(i, ls, operatorMod)
}

func powOP(i Instruction, ls *LuaState) {
	binaryArith(i, ls, operatorPow)
}

func divOP(i Instruction, ls *LuaState) {
	binaryArith(i, ls, operatorDiv)
}

func idivOP(i Instruction, ls *LuaState) {
	binaryArith(i, ls, operatorIDiv)
}

func binandOP(i Instruction, ls *LuaState) {
	binaryArith(i, ls, operatorBinAnd)
}

func binorOP(i Instruction, ls *LuaState) {
	binaryArith(i, ls, operatorBinOr)
}

func binxorOP(i Instruction, ls *LuaState) {
	binaryArith(i, ls, operatorBinXor)
}

func shlOP(i Instruction, ls *LuaState) {
	binaryArith(i, ls, operatorShl)
}

func shrOP(i Instruction, ls *LuaState) {
	binaryArith(i, ls, operatorShr)
}

func unmOP(i Instruction, ls *LuaState) {
	unaryArith(i, ls, operatorUnm)
}

func binnotOP(i Instruction, ls *LuaState) {
	unaryArith(i, ls, operatorBinNot)
}

func newTableOP(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	ls.NewTable(floatPointByteToInt(b), floatPointByteToInt(c))
	ls.Replace(a)
}

func getTableOP(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	ls.GetRK(c)
	ls.GetTable(b)
	ls.Replace(a)
}

func setTableOP(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	ls.GetRK(b)
	ls.GetRK(c)
	ls.SetTable(a)
}

func setListOP(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	if c > 0 {
		c = c - 1
	} else {
		c = Instruction(ls.Fetch()).Ax()
	}

	bIsZero := b == 0
	if bIsZero {
		b = int(ls.ToInteger(-1)) - a - 1
		ls.Pop(1)
	}

	idx := int64(c * fieldsPerFlush)
	for j := 1; j <= b; j++ {
		idx++
		ls.PushValue(a + j)
		ls.SetI(a, idx)
	}

	if bIsZero {
		for j := ls.RegisterCount(); j <= ls.GetTopIndex(); j++ {
			idx++
			ls.PushValue(j)
			ls.SetI(a, idx)
		}

		// clear stack
		ls.SetTop(ls.RegisterCount() - 1)
	}
}

func closureOP(i Instruction, ls *LuaState) {
	a, bx := i.ABx()
	ls.LoadProto(bx)
	ls.Replace(a)
}

func callOP(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	//ls.RemoveNilTail()
	nArgs := pushFuncAndArgs(a, b, ls)
	ls.Call(nArgs, c-1)
	popResults(a, c, ls)
}

func returnOP(i Instruction, ls *LuaState) {
	a, b, _ := i.ABC()
	//ls.RemoveNilTail()
	if b == 1 {
		//no results
	} else if b > 1 {
		for i := a; i <= a+b-2; i++ {
			ls.PushValue(i)
		}
	} else {
		fixStack(a, ls)
	}
}

func varargOP(i Instruction, ls *LuaState) {
	a, b, _ := i.ABC()
	if b != 1 {
		ls.LoadVararg(b - 1)
		popResults(a, b, ls)
	}
}

func tailCallOP(i Instruction, ls *LuaState) {
	a, b, _ := i.ABC()
	c := 0
	nArgs := pushFuncAndArgs(a, b, ls)
	ls.Call(nArgs, c-1)
	popResults(a, c, ls)
}

func selfOP(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	ls.Copy(b, a+1)
	ls.GetRK(c)
	ls.GetTable(b)
	ls.Replace(a)
}

func setUpvalOP(i Instruction, ls *LuaState) {
	a, b, _ := i.ABC()
	ls.Copy(a, luaUpvalueIndex(b))
}

func getUpvalOP(i Instruction, ls *LuaState) {
	a, b, _ := i.ABC()
	ls.Copy(luaUpvalueIndex(b), a)
}

func getTabupOP(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	ls.GetRK(c)
	ls.GetTable(luaUpvalueIndex(b))
	ls.Replace(a)
}

func setTabupOP(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	ls.GetRK(b)
	ls.GetRK(c)
	ls.SetTable(luaUpvalueIndex(a))
}

func tForCallOP(i Instruction, ls *LuaState) {
	a, _, c := i.ABC()
	pushFuncAndArgs(a, 3, ls)
	ls.Call(2, c)
	popResults(a+3, c+1, ls)
}

func tForLoopOP(i Instruction, ls *LuaState) {
	a, sBx := i.AsBx()
	if !ls.IsNil(a + 1) {
		ls.Copy(a+1, a)
		ls.AddPC(sBx)
	}
}

func pushFuncAndArgs(a, b int, ls *LuaState) (nArgs int) {
	if b >= 1 {
		for i := a; i < a+b; i++ {
			ls.PushValue(i)
		}
		return b - 1
	}
	fixStack(a, ls)
	return ls.GetTopIndex() - ls.RegisterCount()
}

func popResults(a, c int, ls *LuaState) {
	if c == 1 {
		//no results
	} else if c > 1 {
		for i := a + c - 2; i >= a; i-- {
			ls.Replace(i)
		}
	} else {
		ls.PushInteger(int64(a))
	}
}

func fixStack(a int, ls *LuaState) {
	x := int(ls.ToInteger(-1))
	ls.Pop(1)

	for i := a; i < x; i++ {
		ls.PushValue(i)
	}
	ls.Rotate(ls.RegisterCount(), x-a)
}

func luaUpvalueIndex(i int) int {
	return luaRegisteryIndex - i - 1
}
