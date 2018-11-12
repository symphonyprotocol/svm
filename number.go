package svm

import "math"

//IFloorDiv mock lua integer floordiv
func IFloorDiv(a, b int64) int64 {
	if a > 0 && b > 0 || a < 0 && b < 0 || a%b == 0 {
		return a / b
	} else {
		return a/b - 1
	}
}

//FFloorDiv mock lua float64 floordiv
func FFloorDiv(a, b float64) float64 {
	return math.Floor(a / b)
}

//IMod mock mod for int64
func IMod(a, b int64) int64 {
	return a - IFloorDiv(a, b)*b
}

//FMod mock mod for float64
func FMod(a, b float64) float64 {
	return a - math.Floor(a/b)*b
}
