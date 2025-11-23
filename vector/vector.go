package vector

import "github.com/bogersw/wbmath"

type Vector[T wbmath.SignedNumber] []T

func New[T wbmath.SignedNumber](elements ...T) Vector[T] {
	vec := make(Vector[T], len(elements), len(elements)*2)
	copy(vec, elements)
	return vec
}

func (v Vector[T]) Add(other Vector[T], offset int) Vector[T] {
	result := make(Vector[T], len(v))
	copy(result, v)
	if offset < len(result) && offset >= 0 {
		index := offset
		for index < len(result) {
			if index-offset >= len(other) {
				break
			}
			result[index] += other[index-offset]
			index += 1
		}
	}
	return result
}

func (v Vector[T]) Subtract(other Vector[T], offset int) Vector[T] {
	return v.Add(other.Scale(T(-1)), offset)
}

func (v Vector[T]) Scale(factor T) Vector[T] {
	result := make(Vector[T], len(v))
	copy(result, v)
	for index, _ := range result {
		result[index] *= factor
	}
	return result
}
