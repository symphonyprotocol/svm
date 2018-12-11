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
	//action   func(i Instruction, ls *LuaState)
	opType int
}

var opcodes = []opcode{
	//     T   mode      name
	opcode{0, modeIABC, "MOVE", opMove},
	opcode{0, modeIABx, "LOADK", opLoadK},
	opcode{0, modeIABx, "LOADKX", opLoadKx},
	opcode{0, modeIABC, "LOADBOOL", opLoadBool},
	opcode{0, modeIABC, "LOADNIL", opLoadNil},
	opcode{0, modeIABC, "GETUPVAL", opGetUpvalue},
	opcode{0, modeIABC, "GETTABUP", opGetTableUp},
	opcode{0, modeIABC, "GETTABLE", opGetTable},
	opcode{0, modeIABC, "SETTABUP", opSetTableUp},
	opcode{0, modeIABC, "SETUPVAL", opSetUpvalue},
	opcode{0, modeIABC, "SETTABLE", opSetTable},
	opcode{0, modeIABC, "NEWTABLE", opNewTable},
	opcode{0, modeIABC, "SELF", opSelf},
	opcode{0, modeIABC, "ADD", opAdd},
	opcode{0, modeIABC, "SUB", opSub},
	opcode{0, modeIABC, "MUL", opMul},
	opcode{0, modeIABC, "MOD", opMod},
	opcode{0, modeIABC, "POW", opPow},
	opcode{0, modeIABC, "DIV", opDiv},
	opcode{0, modeIABC, "IDIV", opIDiv},
	opcode{0, modeIABC, "BAND", opBinAnd},
	opcode{0, modeIABC, "BOR", opBinOr},
	opcode{0, modeIABC, "BXOR", opBinXor},
	opcode{0, modeIABC, "SHL", opBinShiftL},
	opcode{0, modeIABC, "SHR", opBinShiftR},
	opcode{0, modeIABC, "UNM", opUminus},
	opcode{0, modeIABC, "BNOT", opBinNot},
	opcode{0, modeIABC, "NOT", opNot},
	opcode{0, modeIABC, "LEN", opLength},
	opcode{0, modeIABC, "CONCAT", opConcat},
	opcode{0, modeIAsBx, "JMP", opJump},
	opcode{0, modeIABC, "EQ", opEqual},
	opcode{0, modeIABC, "LT", opLessThan},
	opcode{0, modeIABC, "LE", opLessEqual},
	opcode{0, modeIABC, "TEST", opTest},
	opcode{0, modeIABC, "TESTSET", opTestSet},
	opcode{0, modeIABC, "CALL", opCall},
	opcode{0, modeIABC, "TAILCALL", opTailCall},
	opcode{0, modeIABC, "RETURN", opReturn},
	opcode{0, modeIAsBx, "FORLOOP", opForLoop},
	opcode{0, modeIAsBx, "FORPREP", opForPrep},
	opcode{0, modeIABC, "TFORCALL", opTForCall},
	opcode{0, modeIAsBx, "TFORLOOP", opTForLoop},
	opcode{0, modeIABC, "SETLIST", opSetList},
	opcode{0, modeIABx, "CLOSURE", opClosure},
	opcode{0, modeIABC, "VARARG", opVarArg},
	opcode{0, modeIAx, "EXTRAARG", opExtraArg},
}
