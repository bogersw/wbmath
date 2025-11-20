// Package wbmath provides several general math functions.
package wbmath

import (
	"math"
)

// Number is a custom constraint that allows integers and floats.
type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 |
		uint16 | uint32 | uint64 | float32 | float64
}

// SignedNumber is a custom constraint that allows signed integers
// and floats.
type SignedNumber interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

// Abs returns the absolute value of the specified number.
func Abs[T SignedNumber](value T) T {
	if value < 0 {
		return -value
	}
	return value
}

// Gcd implements the Euclidean algorithm for computing the greatest common
// divisor (gcd). The gcd of two numbers is the largest positive integer that
// divides both numbers without leaving a remainder.
// The function is associative: for example, the gcd of three numbers
// a, b, c is equal to: gcd(a, b, c) = gcd(a, gcd(b, c). And so on.
// The function takes two integers and returns an integer (the gcd).
func Gcd(a int, b int) int {
	if a == 0 || b == 0 {
		return a + b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// IsNthRootInt checks if the specified integer value can be expressed
// in terms of a root of the specified degree.
func IsNthRootInt(value int, degree uint) bool {
	nthRoot := math.Pow(float64(value), 1.0/float64(degree))
	nthRootRounded := math.Round(nthRoot)
	return math.Abs(nthRoot-nthRootRounded) <= 0.00001
}

// PowInt returns base**exponent, with base an integer and exponent an
// unsigned integer. The returned power is an integer.
func PowInt(base int, exponent uint) int {
	// In each step the exponent is divided by two (shift bits
	// to the right by 1) and the value is squared.
	result := 1
	for exponent > 0 {
		if exponent&1 != 0 {
			// Uneven exponent
			result *= base
		}
		// Shift bits right by 1 (basically divide by 2, binary)
		exponent >>= 1
		// Square the value. To prevent doing this when the end
		// state (exp=0) is reached, we check for 0.
		if exponent != 0 {
			base *= base
		}
	}
	return result
}

// Round takes a number of type float32 / float64, rounds it to thw
// specified number of decimal places and returns the rounded number.
func Round[T float32 | float64](num T, decimalPlaces uint) T {
	scale := math.Pow(10, float64(decimalPlaces))
	return T(math.Round(float64(num)*scale) / scale)
}

// IsInteger checks if the specified number of type float32 / float64
// is an integer value.
func IsInteger[T float32 | float64](num T) bool {
	// Check for NaN's and infinity: these are never ints
	if math.IsNaN(float64(num)) || math.IsInf(float64(num), 0) {
		return false
	}
	switch any(num).(type) {
	case float32:
		return float32(num) == float32(math.Trunc(float64(num)))
	case float64:
		return float64(num) == math.Trunc(float64(num))
	default:
		return false
	}
}
