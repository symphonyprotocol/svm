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
		panic("todo!")
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

func newTable(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	ls.NewTable(floatPointByteToInt(b), floatPointByteToInt(c))
	ls.Replace(a)
}

func getTable(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	ls.GetRK(c)
	ls.GetTable(b)
	ls.Replace(a)
}

func setTable(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	ls.GetRK(b)
	ls.GetRK(c)
	ls.SetTable(a)
}

func setList(i Instruction, ls *LuaState) {
	a, b, c := i.ABC()
	if c > 0 {
		c = c - 1
	} else {
		c = Instruction(ls.Fetch()).Ax()
	}

	idx := int64(c * fieldsPerFlush)
	for j := 1; j <= b; j++ {
		idx++
		ls.PushValue(a + j)
		ls.SetI(a, idx)
	}
}
