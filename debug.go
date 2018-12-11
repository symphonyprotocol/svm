package svm

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func printStack(ls *luaStack) {
	fmt.Printf("{ID:%d}", ls.id)
	top := ls.topIndex()
	for i := 0; i <= top; i++ {
		val := ls.get(i)
		switch t := typeOf(val); t {
		case LuaTypeBoolean:
			fmt.Printf("[%t]", covertToBoolean(val))
		case LuaTypeNumber:
			fVal, _ := convertToFloat(val)
			fmt.Printf("[%g]", fVal)
		case LuaTypeString:
			fmt.Printf("[%q]", val)
		default:
			if c, ok := val.(*luaClosure); ok {
				if c.goFunc != nil {
					funcName := runtime.FuncForPC(reflect.ValueOf(c.goFunc).Pointer()).Name()
					names := strings.Split(funcName, "/")
					name := names[len(names)-1]
					fmt.Printf("[goFunc:%s]", name)
				} else if c.proto != nil {
					funcType := "main"
					if c.proto.LineDefined > 0 {
						funcType = "function"
					}

					fmt.Printf("[%s <%s:%d, %d>]", funcType, c.proto.Source, c.proto.LineDefined, c.proto.LastLineDefined)

				} else {
					fmt.Printf("[%s]", typeName(t))
				}
			} else if table, ok := val.(*luaTable); ok {
				fmt.Printf("[table:%d]", table.id)
			} else {
				fmt.Printf("[%s]", typeName(t))
			}
		}
	}
	fmt.Println()
}

func printState(ls *LuaState) {
	top := ls.GetTopIndex()
	for i := 0; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case LuaTypeBoolean:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case LuaTypeNumber:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case LuaTypeString:
			fmt.Printf("[%q]", ls.ToString(i))
		default:
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}

func listProto(proto *LuaTrunkProto) {
	printHeader(proto)
	printCode(proto)
	printDetail(proto)
	for _, p := range proto.Protos {
		listProto(p)
	}
}

func printHeader(proto *LuaTrunkProto) {
	funcType := "main"
	if proto.LineDefined > 0 {
		funcType = "function"
	}

	varargFlag := ""
	if proto.IsVararg > 0 {
		varargFlag = "+"
	}

	fmt.Printf("\n%s <%s:%d, %d> (%d instructions)\n", funcType, proto.Source, proto.LineDefined, proto.LastLineDefined, len(proto.Code))
	fmt.Printf("%d%s params, %d slots, %d values,", proto.NumParam, varargFlag, proto.MaxStackSize, len(proto.Upvalues))
	fmt.Printf("%d locals, %d constants, %d functions\n", len(proto.LocVars), len(proto.Constants), len(proto.Protos))
}

func printCode(proto *LuaTrunkProto) {
	for pc, c := range proto.Code {
		line := "-"
		if len(proto.LineInfo) > 0 {
			line = fmt.Sprintf("%d", proto.LineInfo[pc])
		}
		i := Instruction(c)
		fmt.Printf("\t%d\t[%s]\t0x%08X\t%s\n", pc+1, line, c, i.OpName())
	}
}

func printDetail(proto *LuaTrunkProto) {
	fmt.Printf("constants (%d):\n", len(proto.Constants))
	for i, k := range proto.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(k))
	}

	fmt.Printf("locals (%d):\n", len(proto.LocVars))
	for i, locVar := range proto.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
	}

	fmt.Printf("upvalues (%d):\n", len(proto.Upvalues))
	for i, upval := range proto.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, upvalName(proto, i), upval.InStack, upval.Idx)
	}
}

func constantToString(k interface{}) string {
	switch k.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", k)
	case float64:
		return fmt.Sprintf("%g", k)
	case int64:
		return fmt.Sprintf("%d", k)
	case string:
		return fmt.Sprintf("%q", k)
	default:
		return "?"
	}
}

func upvalName(proto *LuaTrunkProto, idx int) string {
	if len(proto.UpvalueNames) > 0 {
		return proto.UpvalueNames[idx]
	}
	return "-"
}
