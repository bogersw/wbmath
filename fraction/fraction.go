package fraction

import (
	"errors"
	"fmt"
	"math"
)

type Fraction struct {
	numerator     int
	denominator   int
	nthRootDegree uint
	sign          rune
}

// New is a constructor function for the fraction package. It returns a Fraction
// struct and an error in case the denominator is zero.
func New(numerator int, denominator int) (Fraction, error) {
	if denominator == 0 {
		return Fraction{}, errors.New("error: integer divide by zero")
	}
	var sign rune = '+'
	if numerator*denominator < 0 {
		sign = '-'
	}
	return Fraction{
			numerator:     absInt(numerator),
			denominator:   absInt(denominator),
			nthRootDegree: 1,
			sign:          sign},
		nil
}

func (f Fraction) Simplify() Fraction {

	numerator := f.numerator
	denominator := f.denominator
	gcd := f.Gcd()
	if gcd != 0 {
		numerator = numerator / gcd
		denominator = denominator / gcd
	}
	return Fraction{
		numerator:     numerator,
		denominator:   denominator,
		nthRootDegree: f.nthRootDegree,
		sign:          f.sign}
}

func (f Fraction) Evaluate() float64 {
	//numerator := float64(f.wholeNumber*f.denominator + f.numerator)
	numerator := float64(f.numerator)
	denominator := float64(f.denominator)
	return float64(f.sign) * math.Pow(numerator/denominator, 1.0/float64(f.nthRootDegree))
}

func (f Fraction) String() string {
	var result string

	fraction := f.Simplify()
	numerator := fraction.numerator
	denominator := fraction.denominator
	wholeNumber := 0

	if numerator >= denominator {
		wholeNumber = numerator / denominator
		numerator = numerator % denominator
	}

	if numerator == 0 {
		result = fmt.Sprintf("%d", wholeNumber)
	} else if wholeNumber == 0 {
		result = fmt.Sprintf("%d/%d", numerator, denominator)
	} else {
		result = fmt.Sprintf("%d %d/%d", wholeNumber, numerator, denominator)
	}
	if fraction.sign == -1 {
		result = "-" + result
	}
	if fraction.nthRootDegree == 1 {
		return result
	}
	return fmt.Sprintf("(%s)^(1/%d)", result, fraction.nthRootDegree)
}

func (f Fraction) Multiply(other Fraction) Fraction {
	// Multiplies two Fractions and returns a new Fraction with the result
	numerator := f.numerator * other.numerator
	denominator := f.denominator * other.denominator
	var sign rune
	if (f.sign == '+' && other.sign == '-') || (f.sign == '-' && other.sign == '+') {
		sign = '-'
	}
	return Fraction{
		numerator:     numerator,
		denominator:   denominator,
		nthRootDegree: f.nthRootDegree,
		sign:          sign}
}

func (f Fraction) MultiplyInt(value int) Fraction {

	// TODO: take sign of value into account
	return Fraction{
		numerator:     f.numerator * powInt(absInt(value), f.nthRootDegree),
		denominator:   f.denominator,
		nthRootDegree: f.nthRootDegree,
		sign:          '+'}
}

func (f Fraction) Add(other Fraction) Fraction {
	// Adds two Fractions and returns a new Fraction with the result
	numerator := f.numerator*other.denominator + other.numerator*f.denominator
	denominator := f.denominator * other.denominator
	return Fraction{
		numerator:     numerator,
		denominator:   denominator,
		nthRootDegree: f.nthRootDegree,
		sign:          f.sign}
}

func (f Fraction) AddInt(value int) Fraction {
	return f.Add(Fraction{
		numerator:     0,
		denominator:   1,
		nthRootDegree: f.nthRootDegree,
		sign:          f.sign})
}

func (f Fraction) Subtract(other Fraction) Fraction {
	// Subtracts two Fractions and returns a new Fraction with the result
	numerator := f.numerator*other.denominator - other.numerator*f.denominator
	denominator := f.denominator * other.denominator
	return Fraction{
		numerator:     numerator,
		denominator:   denominator,
		nthRootDegree: f.nthRootDegree,
		sign:          f.sign}
}

func (f Fraction) SubtractInt(value int) Fraction {
	return f.Subtract(Fraction{
		numerator:     0,
		denominator:   1,
		nthRootDegree: f.nthRootDegree,
		sign:          f.sign})
}

// Divide divides a fraction by another fraction.
func (f Fraction) Divide(other Fraction) Fraction {
	// Divides two Fractions and returns a new Fraction with the result
	numerator := f.numerator * other.denominator
	denominator := f.denominator * other.numerator
	return Fraction{
		numerator:     numerator,
		denominator:   denominator,
		nthRootDegree: 1,
		sign:          '+'}
}

// DivideInt divides a fraction by an integer value. It returns a new Fraction
// and an error (which is nil if no error occurs).
func (f Fraction) DivideInt(value int) (Fraction, error) {
	if value == 0 {
		return Fraction{}, errors.New("error: integer divide by zero")
	}
	var sign rune = '+'
	if value < 0 {
		sign = '-'
	}
	return f.Divide(Fraction{
			numerator:     value,
			denominator:   1,
			nthRootDegree: 1,
			sign:          sign}),
		nil
}

func (f Fraction) Pow(exponent uint) Fraction {
	//numerator := f.numerator
	if exponent == f.nthRootDegree {
		// Power n and an n-th root degree cancel each other
		return Fraction{numerator: f.numerator, denominator: f.denominator, nthRootDegree: 1, sign: f.sign}
	}
	// Power of an n-th root equals the n-th root of the power
	// TODO: check if exponent / nthRootDegree is not a float
	numerator := powInt(f.numerator, exponent)
	denominator := powInt(f.denominator, exponent)
	return Fraction{
		numerator:     numerator,
		denominator:   denominator,
		nthRootDegree: f.nthRootDegree,
		sign:          f.sign}
}

func (f Fraction) NthRoot(degree uint) Fraction {
	//TODO: negative fractions can have uneven degrees, not even degrees
	numerator := f.numerator
	denominator := f.denominator
	if isNthRootInt(numerator, degree) && isNthRootInt(f.denominator, degree) {
		// The nth-roots of the numerator and the denominator are integers => process
		numerator = int(math.Round(math.Pow(float64(numerator), 1.0/float64(degree))))
		denominator = int(math.Round(math.Pow(float64(denominator), 1.0/float64(degree))))
		return Fraction{numerator: numerator, denominator: denominator, nthRootDegree: 1}
	}
	return Fraction{numerator: numerator, denominator: f.denominator, nthRootDegree: degree}
}

func (f Fraction) Gcd() int {
	// Euclidean algorithm for computing the greatest common divisor (gcd).
	// The gcd of two numbers is the largest positive integer that divides
	// both numbers without leaving a remainder.
	// The function is associative: for example, the gcd of three numbers
	// a, b, c is equal to: gcd(a, b, c) = gcd(a, gcd(b, c). And so on.
	a := f.numerator
	b := f.denominator

	if a == 0 || b == 0 {
		return a + b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// ============================================================================
// Private functions
// ============================================================================

func isNthRootInt(value int, degree uint) bool {
	nthRoot := math.Pow(float64(value), 1.0/float64(degree))
	nthRootRounded := math.Round(nthRoot)
	return math.Abs(nthRoot-nthRootRounded) <= 0.00001
}

// Pow returns base**exp, with base an integer and exp an unsigned
// integer.
func powInt(base int, exponent uint) int {
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

func absInt(value int) int {
	return max(-value, value)
	//if value < 0 {
	//	return -value
	//}
	//return value
}
