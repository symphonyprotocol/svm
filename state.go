package svm

//LuaState lua state object
type LuaState struct {
	stack *luaStack
	proto *LuaTrunkProto
	pc    int
}

//NewLuaState create a lua state object
func NewLuaState(proto *LuaTrunkProto) *LuaState {
	return &LuaState{
		stack: newLuaStack(),
		proto: proto,
		pc:    0,
	}
}

//PC get current pc
func (ls *LuaState) PC() int {
	return ls.pc
}

//AddPC add pc for n count
func (ls *LuaState) AddPC(n int) {
	ls.pc += n
}

//Fetch get current pc point to code
func (ls *LuaState) Fetch() uint32 {
	code := ls.proto.Code[ls.pc]
	ls.pc++
	return code
}

//GetConst get constant value at idx and push it into stack
func (ls *LuaState) GetConst(idx int) {
	c := ls.proto.Constants[idx]
	ls.stack.push(c)
}

//GetRK get constant or register address and push it into stack
func (ls *LuaState) GetRK(rk int) {

}
