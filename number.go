package svm

import (
	"math"
	"strconv"
)

func intFloorDiv(a, b int64) int64 {
	if a > 0 && b > 0 || a < 0 && b < 0 || a%b == 0 {
		return a / b
	}
	return a/b - 1
}

func floatFloorDiv(a, b float64) float64 {
	return math.Floor(a / b)
}

func intMod(a, b int64) int64 {
	return a - intFloorDiv(a, b)*b
}

func floatMod(a, b float64) float64 {
	return a - math.Floor(a/b)*b
}

func intShiftLeft(a, n int64) int64 {
	if n >= 0 {
		return a << uint64(n)
	}
	return intShiftRight(a, -n)

}

func intShiftRight(a, n int64) int64 {
	if n >= 0 {
		return int64(uint64(a) >> uint64(n))
	}
	return intShiftRight(a, -n)

}

func floatToInteger(f float64) (int64, bool) {
	i := int64(f)
	return i, float64(i) == f
}

func parseInteger(str string) (int64, bool) {
	i, err := strconv.ParseInt(str, 10, 64)
	return i, err == nil
}

func parseFloat(str string) (float64, bool) {
	f, err := strconv.ParseFloat(str, 64)
	return f, err == nil
}

func stringToInteger(s string) (int64, bool) {
	if i, ok := parseInteger(s); ok {
		return i, true
	}
	if f, ok := parseFloat(s); ok {
		return floatToInteger(f)
	}
	return 0, false
}

func convertToFloat(val luaValue) (float64, bool) {
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	case string:
		return parseFloat(x)
	default:
		return 0, false
	}
}

func converToInteger(val luaValue) (int64, bool) {
	switch x := val.(type) {
	case int64:
		return x, true
	case float64:
		return floatToInteger(x)
	case string:
		return stringToInteger(x)
	default:
		return 0, false
	}
}
