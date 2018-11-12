package svm

/*
-----------------------
code mode data structure
      |<--        32 byte          -->|
iABC  | B:9 | C:9 |  A:8  | Opcode:6  |
iABx  |   Bx:18   |  A:8  | Opcode:6  |
iAsBx |  sBx:18   |  A:8  | Opcode:6  |
iAx   |        Ax:26      | Opcode:6  |

*/

const (
	modeIABC = iota
	modeIABx
	modeIAsBx
	modeIAx
)

const (
	maxArgBx  = 1<<18 - 1
	maxArgSBx = maxArgBx >> 1
)

/*
R(A), R(B), R(C) - the register numbered A, B or C
Kst(Bx) - the constant numbered |Bx|
KPROTO[Bx] - the function prototype numbered |Bx|
RK(B), RK(C) - the register numbered B or C; if B or C is negative, the constant numbered |B| or |C|
FPF - Fileds per flush
UpValue[B] - the upvalue numbered B
pc - program counter, whick instruction will be executed next
closure - make a closure
*/
const (
	opMove     = iota //MOVE(A, B) R(A) := R(B)
	opLoadK           //LOADK(A, Bx) R(A) := Kst(Bx)
	opLoadKx          //LOADKX(A) R(A) := Kst(extra arg)
	opLoadBool        //LOADBOOL(A, B, C) R(A) := (Bool)B; if(C) pc++
	opLoadNil         //LOADNIL(A, B) R(A),R(A+1),...,R(A+B) := nil

	opGetUpvalue //GETUPVAL(A, B) R(A) := UpValue[B]
	opGetTableUp //GETTABLUP(A, B, C) R(A) := UpValue[B][RK(C)]
	opGetTable   //GETTABLE(A, B, C) R(A) := R(B)[RK(C)]

	opSetTableUp //SETTABLEUP(A, B, C) UpValue[A][RK(B)] := RK(C)
	opSetUpvalue //SETUPVAL(A, B) UpValue[B] := R(A)
	opSetTable   //SETTABLE(A, B, C) R(A)[RK[B]] := RK(C)

	opNewTable //NEWTABLE(A, B, C) R(A) := {}(size= B, C)

	opSelf //SELF(A, B, C) R(A+1) := R(B); R(A) := R(B)[RK(C)]

	opAdd       //ADD(A, B, C) R(A) := RK(B) + RK(C)
	opSub       //SUB(A, B, C) R(A) := RK(B) - RK(C)
	opMul       //MUL(A, B, C) R(A) := RK(B) * RK(C)
	opMod       //MOD(A, B, C) R(A) := RK(B) % RK(C)
	opPow       //POW(A, B, C) R(A) := RK(B) ^ RK(C)
	opDiv       //DIV(A, B, C) R(A) := RK(B) / RK(C)
	opIDiv      //IDIV(A, B, C) R(A) := RK(B) // RK(C)
	opBinAnd    //BAND(A, B, C) R(A) := RK(B) & RK(C)
	opBinOr     //BOR(A, B, C) R(A) := RK(B) | RK(C)
	opBinXor    //BXOR(A, B, C) R(A) := RK(B) ~ RK(C)
	opBinShiftL //SHL(A, B, C) R(A) := RK(B) << RK(C)
	opBinShiftR //SHR(A, B, C) R(A) := RK(B) >> RK(C)
	opUminus    //UNM(A, B) R(A) := -R(B)
	opBinNot    //BNOT(A, B) R(A) := -R(B)
	opNot       //NOT(A, B) R(A) := not R(B)
	opLength    //LEN(A, B) R(A) := length of R(B)

	opConcat //CONCAT(A, B, C) R(A) := R(B).. .. ..(RC)

	opJump      //JMP(A, sBx) pc+=sBx, if(A) close all upvalues >= R(A) + 1
	opEqual     //EQ(A, B, C) if((RK(B) == RK(C)) ~= A) then pc++
	opLessThan  //LT(A, B, C) if((RK(B) < RK(C)) ~= A) then pc++
	opLessEqual //LE(A, B, C) if((RK(B) <= RK(C)) ~= A) then pc++

	opTest    //TEST(A, C) if not(R(A) <=> C) then pc++
	opTestSet //TESTSET(A, B, C) if(R(B) <=> C) then R(A) := R(B) else pc++

	opCall     //CALL(A, B, C) R(A), ... , R(A+C-2) := R(A)(R(A+1), ... , R(A+B-1))
	opTailCall //TAILCALL(A, B, C) return R(A)(R(A+1), ... , R(A+B-1))
	opReturn   //RETURN(A, B) return R(A), ... , R(A+B-2)

	opForLoop //FORLOOP(A, sBx) R(A) += R(A+2); if R(A) <?= R(A+1) then {pc += sBx; R(A+3) = R(A)}
	opForPrep //FORPREP(A, sBx) R(A) -= R(A+2); pc += sBx

	opTForCall //TFORCALL(A, C) R(A+3), ... , R(A+2+C) := R(A)(R(A+1), R(A+2))
	opTForLoop //TFORLOOP(A, sBx) if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx}

	opSetList //SETLIST(A, B, C) R(A)[(C-1)*FPF+i] := R(A+i), i <= i <= B

	opClosure //CLOSURE(A, Bx) R(A) := closure(KPROTO[Bx])

	opVarArg //VARARG(A, B) R(A), R(A+1), ... , R(A+B-2) = vararg

	opExtraArg //EXTRAARG(Ax) extra (larger) argument for previous opcode
)

type opcode struct {
	testFlag byte
	opMode   byte
	name     string
}

var opcodes = []opcode{
	//     T   mode      name
	opcode{0, modeIABC, "MOVE"},
	opcode{0, modeIABx, "LOADK"},
	opcode{0, modeIABx, "LOADKX"},
	opcode{0, modeIABC, "LOADBOOL"},
	opcode{0, modeIABC, "LOADNIL"},
	opcode{0, modeIABC, "GETUPVAL"},
	opcode{0, modeIABC, "GETTABUP"},
	opcode{0, modeIABC, "GETTABLE"},
	opcode{0, modeIABC, "SETTABUP"},
	opcode{0, modeIABC, "SETUPVAL"},
	opcode{0, modeIABC, "SETTABLE"},
	opcode{0, modeIABC, "NEWTABLE"},
	opcode{0, modeIABC, "SELF"},
	opcode{0, modeIABC, "ADD"},
	opcode{0, modeIABC, "SUB"},
	opcode{0, modeIABC, "MUL"},
	opcode{0, modeIABC, "MOD"},
	opcode{0, modeIABC, "POW"},
	opcode{0, modeIABC, "DIV"},
	opcode{0, modeIABC, "IDIV"},
	opcode{0, modeIABC, "BAND"},
	opcode{0, modeIABC, "BOR"},
	opcode{0, modeIABC, "BXOR"},
	opcode{0, modeIABC, "SHL"},
	opcode{0, modeIABC, "SHR"},
	opcode{0, modeIABC, "UNM"},
	opcode{0, modeIABC, "BNOT"},
	opcode{0, modeIABC, "NOT"},
	opcode{0, modeIABC, "LEN"},
	opcode{0, modeIABC, "CONCAT"},
	opcode{0, modeIAsBx, "JMP"},
	opcode{0, modeIABC, "EQ"},
	opcode{0, modeIABC, "LT"},
	opcode{0, modeIABC, "LE"},
	opcode{0, modeIABC, "TEST"},
	opcode{0, modeIABC, "TESTSET"},
	opcode{0, modeIABC, "CALL"},
	opcode{0, modeIABC, "TAILCALL"},
	opcode{0, modeIABC, "RETURN"},
	opcode{0, modeIAsBx, "FORLOOP"},
	opcode{0, modeIAsBx, "FORPREP"},
	opcode{0, modeIABC, "TFORCALL"},
	opcode{0, modeIAsBx, "TFORLOOP"},
	opcode{0, modeIABC, "SETLIST"},
	opcode{0, modeIABx, "CLOSURE"},
	opcode{0, modeIABC, "VARARG"},
	opcode{0, modeIAx, "EXTRAARG"},
}

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
