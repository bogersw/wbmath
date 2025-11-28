// Package vector provides a generic, slice-backed numeric Vector[T] type.
//
// The Vector type is a wrapper around a Go slice that supports common
// numeric operations and utilities for elements constrained by
// `wbmath.SignedNumber`. Available functionality includes constructors
// (New, NewFromValue, NewFromRange), cloning (Clone), element-wise arithmetic
// with optional offsets (Add, Subtract, Multiply, Divide), scalar
// multiplication (Scale), reductions (Sum, Product), and rounding (Round).
//
// Important details:
//
// (*) Most mutating methods operate in-place and also return the modified Vector
// to allow chaining. Call Clone() first when an independent copy is needed.
//
// (*) NewFromRange always returns a Vector[float64] regardless of the input type.
//
// (*) Arithmetic with offsets only processes indices common to both vectors;
// out-of-range elements are ignored.
//
// (*) The type parameter T must satisfy wbmath.SignedNumber, so both integer and
// floating-point element types are supported.
package vector

import (
	"github.com/bogersw/wbmath"
)

type Vector[T wbmath.SignedNumber] []T

// ============================================================================
// Constructor functions
// ============================================================================

// New is a constructor function that accepts an arbitrary number of elements
// of type `SignedNumber` and returns a Vector.
func New[T wbmath.SignedNumber](elements ...T) Vector[T] {
	vec := make(Vector[T], len(elements), len(elements)*2)
	copy(vec, elements)
	return vec
}

// NewFromValue is a constructor function that takes a `value` of type
// `SignedNumber` and a count: it returns a Vector with `count` elements
// that are all equal to `value`. Can be used to create a Vector with all
// ones or zeros for example.
// Returns a Vector of type T.
func NewFromValue[T wbmath.SignedNumber](value T, count int) Vector[T] {
	vec := make(Vector[T], count, count*2)
	for i := 0; i < count; i++ {
		vec[i] = value
	}
	return vec
}

// NewFromRange is a constructor function that returns a Vector with `min`
// and `max` as the elements at the start and at the end: `steps` elements
// will be added in between (evenly distributed). If `steps` is 0 a Vector
// with only the elements `min` and `max` is returned.
// Always returns a Vector with elements of type float64.
func NewFromRange[T wbmath.SignedNumber](min T, max T, steps uint) Vector[float64] {
	fMin := float64(min)
	fMax := float64(max)
	vec := make(Vector[float64], 2+steps, (2+steps)*2)
	vec[0] = fMin
	vec[steps+1] = fMax
	for i := 1; i <= int(steps); i++ {
		vec[i] = fMin + float64(i)*((fMax-fMin)/float64(steps+1))
	}
	return vec
}

// ============================================================================
// Private methods
// ============================================================================

func (v Vector[T]) operation(other Vector[T], offset int, operation string) Vector[T] {
	if offset < len(v) && offset >= 0 {
		index := offset
		for index < len(v) {
			if index-offset >= len(other) {
				break
			}
			if operation == "add" {
				v[index] += other[index-offset]
			} else if operation == "subtract" {
				v[index] -= other[index-offset]
			} else if operation == "multiply" {
				v[index] *= other[index-offset]
			} else if operation == "divide" {
				v[index] /= other[index-offset]
			}
			index += 1
		}
	}
	return v
}

// ============================================================================
// Public methods
// ============================================================================

// Clone returns a new Vector with its own backing array (modifying the result
// will not affect the original).
func (v Vector[T]) Clone() Vector[T] {
	clone := make(Vector[T], len(v), len(v)*2)
	copy(clone, v)
	return clone
}

// Add adds the specified Vector to the current Vector (in-place, unless a
// Clone is made beforehand). Addition is element-wise, based on the index,
// but  when an offset is specified the vectors are shifted by that amount.
// When the specified Vector is shorter than the current Vector, only matching
// elements are added: when it's longer, extra elements are ignored.
func (v Vector[T]) Add(other Vector[T], offset int) Vector[T] {
	return v.operation(other, offset, "add")
}

// Subtract subtracts the specified Vector from the current Vector (in-place,
// unless a Clone is made beforehand). Subtraction is element-wise, based on the
// index, but when an offset is specified the vectors are shifted by that amount.
// When the specified Vector is shorter than the current Vector, only matching
// elements are subtracted: when it's longer, extra elements are ignored.
func (v Vector[T]) Subtract(other Vector[T], offset int) Vector[T] {
	return v.Add(other.Scale(T(-1)), offset)
}

// Multiply multiplies the specified Vector with the current Vector (in-place,
// unless a Clone is made beforehand). Multiplication is element-wise, based on
// their index, but when an offset is specified the vectors are shifted by that
// amount. When the specified Vector is shorter than the current Vector, only
// matching elements are multiplied: when it's longer, extra elements are ignored.
func (v Vector[T]) Multiply(other Vector[T], offset int) Vector[T] {
	return v.operation(other, offset, "multiply")
}

// Divide divides the current Vector by the specified Vector (in-place,
// unless a Clone is made beforehand). Division is element-wise, based on
// the index, but when an offset is specified the vectors are shifted by that
// amount. When the specified Vector is shorter than the current Vector, only
// matching elements are divided: when it's longer, extra elements are ignored.
func (v Vector[T]) Divide(other Vector[T], offset int) Vector[T] {
	return v.operation(other, offset, "divide")
}

// Scale implements scalar multiplication: every element in the current Vector
// is multiplied by `factor`. This operation is in-place, unless a Clone is made
// beforehand.
func (v Vector[T]) Scale(factor T) Vector[T] {
	for index := range v {
		v[index] *= factor
	}
	return v
}

// Sum returns the sum of the elements of a Vector.
func (v Vector[T]) Sum() T {
	var sum T = 0
	for i := 0; i < len(v); i++ {
		sum += v[i]
	}
	return sum
}

// Product returns the product of the elements of a Vector.
func (v Vector[T]) Product() T {
	var product T = 0
	for i := 0; i < len(v); i++ {
		product *= v[i]
	}
	return product
}

// Round rounds all elements of a Vector to the specified precision.
func (v Vector[T]) Round(precision uint) Vector[T] {
	if len(v) == 0 || v == nil {
		return v
	}
	switch any(v[0]).(type) {
	case float32:
		for i := range v {
			v[i] = T(wbmath.Round(float32(v[i]), precision))
		}
	case float64:
		for i := range v {
			v[i] = T(wbmath.Round(float64(v[i]), precision))
		}
	}
	return v
}
