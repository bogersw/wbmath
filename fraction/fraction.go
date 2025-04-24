package fraction

import (
	"errors"
	"fmt"
	"math"
)

type Fraction struct {
	numerator     int
	denominator   int
	wholeNumber   int
	nthRootDegree int
	sign          int
}

// New is a constructor function for the fraction package. It returns a Fraction
// struct and an error in case the denominator is zero.
func New(numerator int, denominator int) (Fraction, error) {
	if denominator == 0 {
		return Fraction{}, errors.New("error: integer divide by zero")
	}
	sign := 1
	if numerator*denominator < 0 {
		sign = -1
	}
	return Fraction{
			numerator:     int(math.Abs(float64(numerator))),
			denominator:   int(math.Abs(float64(denominator))),
			wholeNumber:   0,
			nthRootDegree: 1,
			sign:          sign},
		nil
}

func (f Fraction) Simplify() Fraction {

	numerator := f.wholeNumber*f.denominator + f.numerator
	denominator := f.denominator
	wholeNumber := 0
	gcd := f.Gcd()
	if gcd != 0 {
		numerator = numerator / gcd
		denominator = denominator / gcd
	}
	if numerator >= denominator {
		wholeNumber = numerator / denominator
		numerator = numerator % denominator
	}
	return Fraction{
		numerator:     numerator,
		denominator:   denominator,
		wholeNumber:   wholeNumber,
		nthRootDegree: f.nthRootDegree,
		sign:          f.sign}
}

func (f Fraction) Evaluate() float64 {
	numerator := float64(f.wholeNumber*f.denominator + f.numerator)
	denominator := float64(f.denominator)
	return float64(f.sign) * math.Pow(numerator/denominator, 1.0/float64(f.nthRootDegree))
}

func (f Fraction) String() string {
	var result string
	if f.numerator == 0 {
		result = fmt.Sprintf("%d", f.wholeNumber)
	} else if f.wholeNumber == 0 {
		result = fmt.Sprintf("%d/%d", f.numerator, f.denominator)
	} else {
		result = fmt.Sprintf("%d %d/%d", f.wholeNumber, f.numerator, f.denominator)
	}
	if f.sign == -1 {
		result = "-" + result
	}
	if f.nthRootDegree == 1 {
		return result
	}
	return fmt.Sprintf("(%s)^(1/%d)", result, f.nthRootDegree)
}

func (f Fraction) Multiply(other Fraction) Fraction {
	// Multiplies two Fractions and returns a new Fraction with the result
	numerator := ((f.wholeNumber * f.denominator) + f.numerator) *
		((other.wholeNumber * other.denominator) + other.numerator)
	denominator := f.denominator * other.denominator
	return Fraction{
		numerator:     numerator,
		denominator:   denominator,
		wholeNumber:   0,
		nthRootDegree: f.nthRootDegree,
		sign:          f.sign * other.sign}
}

func (f Fraction) MultiplyInt(value int) Fraction {
	return f.Multiply(Fraction{numerator: 0, denominator: 1, wholeNumber: value})
}

func (f Fraction) Add(other Fraction) Fraction {
	// Adds two Fractions and returns a new Fraction with the result
	numerator := (f.wholeNumber*f.denominator+f.numerator)*other.denominator +
		(other.wholeNumber*other.denominator+other.numerator)*f.denominator
	denominator := f.denominator * other.denominator
	return Fraction{numerator: numerator, denominator: denominator, wholeNumber: 0}
}

func (f Fraction) AddInt(value int) Fraction {
	return f.Add(Fraction{numerator: 0, denominator: 1, wholeNumber: value})
}

func (f Fraction) Subtract(other Fraction) Fraction {
	// Subtracts two Fractions and returns a new Fraction with the result
	numerator := (f.wholeNumber*f.denominator+f.numerator)*other.denominator -
		(other.wholeNumber*other.denominator+other.numerator)*f.denominator
	denominator := f.denominator * other.denominator
	return Fraction{numerator: numerator, denominator: denominator, wholeNumber: 0}
}

func (f Fraction) SubtractInt(value int) Fraction {
	return f.Subtract(Fraction{numerator: 0, denominator: 1, wholeNumber: value})
}

func (f Fraction) Divide(other Fraction) Fraction {
	// Divides two Fractions and returns a new Fraction with the result
	numerator := (f.wholeNumber*f.denominator + f.numerator) * other.denominator
	denominator := f.denominator * (other.wholeNumber*other.denominator + other.numerator)
	return Fraction{numerator: numerator, denominator: denominator, wholeNumber: 0}
}

func (f Fraction) DivideInt(value int) (Fraction, error) {
	if value == 0 {
		return Fraction{}, errors.New("error: integer divide by zero")
	}
	return f.Divide(Fraction{numerator: 0, denominator: 1, wholeNumber: value}), nil
}

func (f Fraction) Pow(exponent int) Fraction {
	numerator := f.wholeNumber*f.denominator + f.numerator
	if exponent == f.nthRootDegree {
		// Power n and an n-th root degree cancel each other
		return Fraction{numerator: f.numerator, denominator: f.denominator, wholeNumber: 0, nthRootDegree: 1}
	}
	// Power of an n-th root equals the n-th root of the power
	numerator = int(math.Round(math.Pow(float64(numerator), float64(exponent))))
	denominator := int(math.Round(math.Pow(float64(f.denominator), float64(exponent))))
	return Fraction{numerator: numerator, denominator: denominator, wholeNumber: 0, nthRootDegree: f.nthRootDegree}
}

func (f Fraction) NthRoot(degree int) Fraction {
	//TODO: negative fractions can have uneven degrees, not even degrees
	numerator := f.wholeNumber*f.denominator + f.numerator
	denominator := f.denominator
	if isNthRootInt(numerator, degree) && isNthRootInt(f.denominator, degree) {
		// The nth-roots of the numerator and the denominator are integers => process
		numerator = int(math.Round(math.Pow(float64(numerator), 1.0/float64(degree))))
		denominator = int(math.Round(math.Pow(float64(denominator), 1.0/float64(degree))))
		return Fraction{numerator: numerator, denominator: denominator, wholeNumber: 0, nthRootDegree: 1}
	}
	return Fraction{numerator: numerator, denominator: f.denominator, wholeNumber: 0, nthRootDegree: degree}
}

func (f Fraction) Gcd() int {
	// Euclidean algorithm for computing the greatest common divisor (gcd).
	// The gcd of two numbers is the largest positive integer that divides
	// both numbers without leaving a remainder.
	// The function is associative: for example, the gcd of three numbers
	// a, b, c is equal to: gcd(a, b, c) = gcd(a, gcd(b, c). And so on.
	a := f.wholeNumber*f.denominator + f.numerator
	b := f.denominator

	if a == 0 || b == 0 {
		return a + b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func isNthRootInt(value int, degree int) bool {
	nthRoot := math.Pow(float64(value), 1.0/float64(degree))
	nthRootRounded := math.Round(nthRoot)
	return math.Abs(nthRoot-nthRootRounded) <= 0.00001
}
