package svm

func intToFloatPointByte(x int) int {
	e := 0
	if x < 8 {
		return x
	}
	for x >= (8 << 4) {
		x = (x + 0xf) >> 4
		e += 4
	}
	for x >= (8 << 1) {
		x = (x + 1) >> 1
		e++
	}
	return ((e + 1) << 3) | (x - 8)
}

func floatPointByteToInt(x int) int {
	if x < 8 {
		return x
	}
	return ((x & 7) + 8) << uint((x>>3)-1)
}

func keyFloatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := floatToInteger(f); ok {
			return i
		}
	}
	return key
}
