package svm

import (
	"fmt"
)

//Instruction code
type Instruction uint32

//Opcode get opcode
func (i Instruction) Opcode() int {
	return int(i & 0x3F)
}

//ABC get IABC arguments
func (i Instruction) ABC() (a, b, c int) {
	a = int(i >> 6 & 0xFF)
	c = int(i >> 14 & 0x1FF)
	b = int(i >> 23 & 0x1FF)
	return
}

//ABx get IAbx arguments
func (i Instruction) ABx() (a, bx int) {
	a = int(i >> 6 & 0xFF)
	bx = int(i >> 14)
	return
}

//AsBx get IAsBx arguments
func (i Instruction) AsBx() (a, sbx int) {
	a, bx := i.ABx()
	return a, bx - maxArgSBx
}

//Ax get IAx arguments
func (i Instruction) Ax() int {
	return int(i >> 6)
}

//OpName get op name
func (i Instruction) OpName() string {
	return opcodes[i.Opcode()].name
}

//OpMode get op mode
func (i Instruction) OpMode() byte {
	return opcodes[i.Opcode()].opMode
}

//Execute execute instruction
func (i Instruction) Execute(ls *LuaState) {
	if ls.isDebug {
		fmt.Println("OP: ", i.OpName())
	}
	switch i.Opcode() {
	case opMove:
		moveOP(i, ls)
	case opLoadK:
		loadKOP(i, ls)
	case opLoadKx:
		loadKxOP(i, ls)
	case opLoadBool:
		loadBoolOP(i, ls)
	case opLoadNil:
		loadNilOP(i, ls)
	case opGetUpvalue:
		getUpvalOP(i, ls)
	case opGetTableUp:
		getTabupOP(i, ls)
	case opGetTable:
		getTableOP(i, ls)
	case opSetTableUp:
		setTabupOP(i, ls)
	case opSetUpvalue:
		setUpvalOP(i, ls)
	case opSetTable:
		setTableOP(i, ls)
	case opNewTable:
		newTableOP(i, ls)
	case opSelf:
		selfOP(i, ls)
	case opAdd:
		addOP(i, ls)
	case opSub:
		subOP(i, ls)
	case opMul:
		mulOP(i, ls)
	case opMod:
		modOP(i, ls)
	case opPow:
		powOP(i, ls)
	case opDiv:
		divOP(i, ls)
	case opIDiv:
		idivOP(i, ls)
	case opBinAnd:
		binandOP(i, ls)
	case opBinOr:
		binorOP(i, ls)
	case opBinXor:
		binxorOP(i, ls)
	case opBinShiftL:
		shlOP(i, ls)
	case opBinShiftR:
		shrOP(i, ls)
	case opUminus:
		unmOP(i, ls)
	case opBinNot:
		binnotOP(i, ls)
	case opNot:
		notOP(i, ls)
	case opLength:
		lenOP(i, ls)
	case opConcat:
		concatOP(i, ls)
	case opJump:
		jmpOP(i, ls)
	case opEqual:
		equalOP(i, ls)
	case opLessThan:
		lessThanOP(i, ls)
	case opLessEqual:
		lessEqualOP(i, ls)
	case opTest:
		testOP(i, ls)
	case opTestSet:
		testSetOP(i, ls)
	case opCall:
		callOP(i, ls)
	case opTailCall:
		tailCallOP(i, ls)
	case opReturn:
		returnOP(i, ls)
	case opForLoop:
		forLoopOP(i, ls)
	case opForPrep:
		forPrepOP(i, ls)
	case opTForCall:
		tForCallOP(i, ls)
	case opTForLoop:
		tForLoopOP(i, ls)
	case opSetList:
		setListOP(i, ls)
	case opClosure:
		closureOP(i, ls)
	case opVarArg:
		//panic("not support vararg")
		varargOP(i, ls)
	case opExtraArg:
		//do nothing
		return
	default:
		panic("no opcode found!")
	}
}
