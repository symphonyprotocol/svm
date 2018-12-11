package svm

type luaUpvalue struct {
	val *luaValue
}

type luaClosure struct {
	proto  *LuaTrunkProto
	upvals []*luaUpvalue
	goFunc GoFunction
}

func newLuaClosure(proto *LuaTrunkProto) *luaClosure {
	c := &luaClosure{proto: proto}
	if nUpvals := len(proto.Upvalues); nUpvals > 0 {
		c.upvals = make([]*luaUpvalue, nUpvals)
	}
	return c
}

func newGoClosure(f GoFunction, nUpvals int) *luaClosure {
	c := &luaClosure{goFunc: f}
	if nUpvals > 0 {
		c.upvals = make([]*luaUpvalue, nUpvals)
	}
	return c
}
