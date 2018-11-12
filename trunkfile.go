package svm

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

/*
trunk file:
---------------
data structure:

type       |    byte

signature        4
version          1
format           1
luac_data        6
cint_size        1
sizet_size       1
instruction_size 1
lua_integer_size 1
lua_number_size  1
luac_int         8
luac_number      8
size_upvalues    1
main function    n

----------------
header detail (only for lua5.3) amd64 ***CAUSION: WE ARE ONLY SUPPORT AMD64, X32 IS IGNORED!
|1B 4C 75 61|    53   |   00   | 19 93 0D 0A 1A 0A |     04    |     08     |        04        |        08       |        08      | 78 56 00 00 00 00 00 00 | 00 00 00 00 00 28 77 40 |
| signature | version | format |     lua data      | cint size | sizet size | instruction size | lua integer size| lua number size|         luac int        |        luac number      |

function prototype

type      |      byte

source            n
line defined      4
last line defined 4
num params        1
is vararg         1
max stack size    1
code              n
constants         n
upvalues          n
protos            n
line info         n
locvars           n
upvalue names     n

-----------------
addition for string in trunk:

NULL --> 0x00
len(str) <= 0xFD (253) --> [n+1(1 byte)][string bytes]
len(str) >= 0xFE (254) --> [0xFF][n+1(8 byte)][string bytes]
*/

const (
	//trunk header defination
	trunkHeaderLuaSignature    = "\x1bLua"
	trunkHeaderLuacVersion     = 0x53
	trunkHeaderLuacFormat      = 0
	trunkHeaderLuacData        = "\x19\x93\r\n\x1a\n"
	trunkHeaderCintSize        = 4
	trunkHeaderSizetSize       = 8
	trunkHeaderInstructionSize = 4
	trunkHeaderLuaIntegerSize  = 8
	trunkHeaderLuaNumberSize   = 8
	trunkHeaderLuacInt         = 0x5678
	trunkHeaderLuacNumber      = 370.5

	trunkProtoTagNil      = 0x00
	trunkProtoTagBoolean  = 0x01
	trunkProtoTagNumber   = 0x03
	trunkProtoTagInteger  = 0x13
	trunkProtoTagShortStr = 0x04
	trunkProtoTagLongStr  = 0x14
)

//LuaTrunkHeader trunk file header struct
type LuaTrunkHeader struct {
	signature       [4]byte
	version         byte
	format          byte
	luacData        [6]byte
	cintSize        byte
	sizetSize       byte
	instructionSize byte
	luaIntegerSize  byte
	luaNumberSize   byte
	luacInt         [8]byte
	luacNumber      [8]byte
}

//LuaTrunkProto trunk proto
type LuaTrunkProto struct {
	Source          string
	LineDefined     uint32
	LastLineDefined uint32
	NumParam        byte
	IsVararg        byte
	MaxStackSize    byte
	Code            []uint32
	Constants       []interface{}
	Upvalues        []upvalue
	Protos          []*LuaTrunkProto
	LineInfo        []uint32
	LocVars         []locVar
	UpvalueNames    []string
}

//LuaTrunk trunk file struct
type LuaTrunk struct {
	LuaTrunkHeader
	sizeUpvalues byte
	mainFunc     *LuaTrunkProto
}

type upvalue struct {
	InStack byte
	Idx     byte
}

type locVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

// trunk file reader
type trunkReader struct {
	reader io.Reader
}

func (r *trunkReader) read(data interface{}) error {
	return binary.Read(r.reader, binary.LittleEndian, data)
}

func (r *trunkReader) readByte() (byte, error) {
	var bte byte
	err := r.read(&bte)
	return bte, err
}

func (r *trunkReader) readUint32() (uint32, error) {
	var i uint32
	err := r.read(&i)
	return i, err
}

func (r *trunkReader) readUint64() (uint64, error) {
	var i uint64
	err := r.read(&i)
	return i, err
}

func (r *trunkReader) readLuaInteger() (int64, error) {
	i, err := r.readUint64()
	if err != nil {
		return 0, err
	}
	return int64(i), err
}

func (r *trunkReader) readLuaNumber() (float64, error) {
	f, err := r.readUint64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(f), err
}

func (r *trunkReader) readString() (string, error) {
	sizet, err := r.readByte()
	if err != nil {
		return "", err
	}
	size := uint(sizet)
	if size == 0 {
		return "", err
	}
	if size == 0xFF {
		sizer, err := r.readUint64()
		if err != nil {
			return "", err
		}
		size = uint(sizer)
	}
	strBytes := make([]byte, size-1)
	err = r.read(strBytes)
	if err != nil {
		return "", err
	}
	return string(strBytes), err
}

func (r *trunkReader) readBytes(n uint) ([]byte, error) {
	bytes := make([]byte, n)
	err := r.read(bytes)
	return bytes, err
}

func (r *trunkReader) checkHeader() {
	signature, _ := r.readBytes(4)
	fmt.Println(string(signature))
	if string(signature) != trunkHeaderLuaSignature {
		panic("trunkfile error: not a valid trunk")
	}
	version, _ := r.readByte()
	if version != trunkHeaderLuacVersion {
		panic("trunkfile error: lua version not right")
	}
	format, _ := r.readByte()
	if format != trunkHeaderLuacFormat {
		panic("trunkfile error: lua format not match")
	}
	luaData, _ := r.readBytes(6)
	if string(luaData) != trunkHeaderLuacData {
		panic("trunkfile error: lua data not match")
	}
	cintSize, _ := r.readByte()
	if cintSize != trunkHeaderCintSize {
		panic("trunkfile error: int size not match")
	}
	sizetSize, _ := r.readByte()
	if sizetSize != trunkHeaderSizetSize {
		panic("trunkfile error: size_t size not match")
	}
	instructionSize, _ := r.readByte()
	if instructionSize != trunkHeaderInstructionSize {
		panic("trunkfile error: instruction size not match")
	}
	luaIntegerSize, _ := r.readByte()
	if luaIntegerSize != trunkHeaderLuaIntegerSize {
		panic("trunkfile error: lua integer size not match")
	}
	luaNumberSize, _ := r.readByte()
	if luaNumberSize != trunkHeaderLuaNumberSize {
		panic("trunkfile error: lua number size not match")
	}
	luacInt, _ := r.readLuaInteger()
	if luacInt != trunkHeaderLuacInt {
		panic("trunkfile error: endianness not match")
	}
	luacNumber, _ := r.readLuaNumber()
	if luacNumber != trunkHeaderLuacNumber {
		panic("trunkfile error: float format not match")
	}
}

func (r *trunkReader) readCode() []uint32 {
	codeLen, err := r.readUint32()
	if err != nil {
		panic("trunkfile error: code length init error")
	}
	code := make([]uint32, codeLen)
	for i := range code {
		code[i], _ = r.readUint32()
	}
	return code
}

func (r *trunkReader) readConstants() []interface{} {
	consLen, err := r.readUint32()
	if err != nil {
		panic("trunkfile error: constant length init error")
	}
	constants := make([]interface{}, consLen)
	for i := range constants {
		constants[i] = r.readConstant()
	}
	return constants
}

func (r *trunkReader) readConstant() interface{} {
	tag, err := r.readByte()
	if err != nil {
		panic("trunkfile error: constant tag init error")
	}
	switch tag {
	case trunkProtoTagNil:
		return nil
	case trunkProtoTagBoolean:
		b, err := r.readByte()
		if err != nil {
			panic("trunkfile error: constant boolean value init error")
		}
		return b != 0
	case trunkProtoTagInteger:
		luaInt, err := r.readLuaInteger()
		if err != nil {
			panic("trunkfile error: constant lua integer init error")
		}
		return luaInt
	case trunkProtoTagNumber:
		luaNumber, err := r.readLuaNumber()
		if err != nil {
			panic("trunkfile error: constant lua number init error")
		}
		return luaNumber
	case trunkProtoTagShortStr:
		shortStr, err := r.readString()
		if err != nil {
			panic("trunkfile error: constant short string init error")
		}
		return shortStr
	case trunkProtoTagLongStr:
		longStr, err := r.readString()
		if err != nil {
			panic("trunkfile error: constant long string init error")
		}
		return longStr
	default:
		panic("trunkfile error: no constant type detected")
	}
}

func (r *trunkReader) readProto(parentSource string) *LuaTrunkProto {
	source, err := r.readString()
	if err != nil {
		panic("trunkfile err: cannot find source")
	}
	if source == "" {
		source = parentSource
	}
	lineNumber, err := r.readUint32()
	if err != nil {
		panic("trunkfile err: linenumber init error")
	}
	lastLineNumber, err := r.readUint32()
	if err != nil {
		panic("trunkfile err: lastLineNumber init error")
	}
	numParam, err := r.readByte()
	if err != nil {
		panic("trunkfile err: numParam init error")
	}
	isVararg, err := r.readByte()
	if err != nil {
		panic("trunkfile err: isVararg init error")
	}
	maxStackSize, err := r.readByte()
	if err != nil {
		panic("trunkfile err: maxStackSize init error")
	}
	return &LuaTrunkProto{
		Source:          source,
		LineDefined:     lineNumber,
		LastLineDefined: lastLineNumber,
		NumParam:        numParam,
		IsVararg:        isVararg,
		MaxStackSize:    maxStackSize,
		Code:            r.readCode(),
		Constants:       r.readConstants(),
		Upvalues:        r.readUpvalues(),
		Protos:          r.readProtos(source),
		LineInfo:        r.readLineInfo(),
		LocVars:         r.readLocVars(),
		UpvalueNames:    r.readUpvalueNames(),
	}
}

func (r *trunkReader) readUpvalues() []upvalue {
	upLen, err := r.readUint32()
	if err != nil {
		panic("trunkfile error: read upvalues length error")
	}
	upvalues := make([]upvalue, upLen)
	for i := range upvalues {
		instack, _ := r.readByte()
		idx, _ := r.readByte()
		upvalues[i] = upvalue{
			InStack: instack,
			Idx:     idx,
		}
	}
	return upvalues
}

func (r *trunkReader) readProtos(parentSource string) []*LuaTrunkProto {
	protosLen, err := r.readUint32()
	if err != nil {
		panic("trunkfile error: read protos length error")
	}
	protos := make([]*LuaTrunkProto, protosLen)
	for i := range protos {
		protos[i] = r.readProto(parentSource)
	}
	return protos
}

func (r *trunkReader) readLineInfo() []uint32 {
	len, err := r.readUint32()
	if err != nil {
		panic("trunkfile error: read line info length error")
	}
	lineInfo := make([]uint32, len)
	for i := range lineInfo {
		lineInfo[i], _ = r.readUint32()
	}
	return lineInfo
}

func (r *trunkReader) readLocVars() []locVar {
	len, err := r.readUint32()
	if err != nil {
		panic("trunkfile error: read locvar length error")
	}
	locVars := make([]locVar, len)
	for i := range locVars {
		name, _ := r.readString()
		start, _ := r.readUint32()
		end, _ := r.readUint32()
		locVars[i] = locVar{
			VarName: name,
			StartPC: start,
			EndPC:   end,
		}
	}
	return locVars
}

func (r trunkReader) readUpvalueNames() []string {
	len, err := r.readUint32()
	if err != nil {
		panic("trunkfile error: read upvalue length error")
	}
	names := make([]string, len)
	for i := range names {
		names[i], _ = r.readString()
	}
	return names
}

//Undump undump turnk file
func Undump(reader io.Reader) *LuaTrunkProto {
	tr := &trunkReader{reader}
	tr.checkHeader()
	tr.readByte()
	return tr.readProto("")
}
