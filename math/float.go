package math

import "math"

// Round returns round of float64
func Round(f float64, n int) float64 {
	pow10N := math.Pow10(n)
	return math.Trunc((f+0.5/pow10N)*pow10N) / pow10N
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func MinInts(a int, vals ...int) int {
	min := a
	for _, v := range vals {
		if v < min {
			min = v
		}
	}
	return min
}

func MaxInts(a int, vals ...int) int {
	max := a
	for _, v := range vals {
		if v > max {
			max = v
		}
	}
	return max
}
