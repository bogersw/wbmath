package fraction

import (
    "errors"
    "fmt"
    "math"
)

type Fraction struct {
    numerator   int
    denominator int
    isNegative  bool
}

// New is a constructor function for the fraction package. It returns a Fraction
// struct and an error in case the denominator is zero.
func New(numerator int, denominator int) (Fraction, error) {
    if denominator == 0 {
        return Fraction{}, errors.New("error: integer divide by zero")
    }
    return Fraction{
            numerator:   absInt(numerator),
            denominator: absInt(denominator),
            isNegative:  numerator*denominator < 0},
        nil
}

func (f Fraction) Simplify() Fraction {

    gcd := f.Gcd()
    if gcd != 0 {
        f.numerator = f.numerator / gcd
        f.denominator = f.denominator / gcd
    }
    return f
}

// Evaluate - calculates and returns the fraction as a float value.
func (f Fraction) Evaluate() float64 {
    if f.isNegative && f.numerator != 0 {
        return -float64(f.numerator) / float64(f.denominator)
    }
    return float64(f.numerator) / float64(f.denominator)
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
    if fraction.isNegative {
        return fmt.Sprintf("-%s", result)
    }
    return fmt.Sprintf("%s", result)
}

// Multiply - multiplies two fractions and returns the result as a new fraction.
func (f Fraction) Multiply(other Fraction) Fraction {
    numerator := f.numerator * other.numerator
    denominator := f.denominator * other.denominator
    // Check sign
    isNegative := false
    if f.isNegative && other.isNegative {
        isNegative = false
    } else if f.isNegative || other.isNegative {
        isNegative = true
    }
    return Fraction{
        numerator:   numerator,
        denominator: denominator,
        isNegative:  isNegative}
}

// MultiplyInt - multiplies a fraction with an integer
func (f Fraction) MultiplyInt(value int) Fraction {

    return f.Multiply(Fraction{
        numerator:   absInt(value),
        denominator: 1,
        isNegative:  value < 0})

}

func (f Fraction) Add(other Fraction) Fraction {
    // Adds two Fractions and returns a new Fraction with the result
    signF := 1
    if f.isNegative {
        signF = -1
    }
    signOther := 1
    if other.isNegative {
        signOther = -1
    }
    numerator := signF*f.numerator*other.denominator + signOther*other.numerator*f.denominator
    denominator := f.denominator * other.denominator
    return Fraction{
        numerator:   absInt(numerator),
        denominator: denominator,
        isNegative:  numerator*denominator < 0}
}

func (f Fraction) AddInt(value int) Fraction {
    return f.Add(Fraction{
        numerator:   absInt(value),
        denominator: 1,
        isNegative:  value < 0})
}

func (f Fraction) Subtract(other Fraction) Fraction {
    // Subtracts two Fractions and returns a new Fraction with the result
    numerator := f.numerator*other.denominator - other.numerator*f.denominator
    denominator := f.denominator * other.denominator
    return Fraction{
        numerator:   numerator,
        denominator: denominator}
}

func (f Fraction) SubtractInt(value int) Fraction {
    return f.Subtract(Fraction{
        numerator:   value,
        denominator: 1})
}

// Divide divides a fraction by another fraction.
func (f Fraction) Divide(other Fraction) Fraction {
    // Divides two Fractions and returns a new Fraction with the result
    numerator := f.numerator * other.denominator
    denominator := f.denominator * other.numerator
    return Fraction{
        numerator:   numerator,
        denominator: denominator}
}

// DivideInt divides a fraction by an integer value. It returns a new Fraction
// and an error (which is nil if no error occurs).
func (f Fraction) DivideInt(value int) (Fraction, error) {
    if value == 0 {
        return Fraction{}, errors.New("error: integer divide by zero")
    }
    // var sign rune = '+'
    // if value < 0 {
    // 	sign = '-'
    // }
    return f.Divide(Fraction{
            numerator:   value,
            denominator: 1}),
        nil
}

func (f Fraction) Pow(exponent uint) Fraction {
    //numerator := f.numerator
    //if exponent == f.nthRootDegree {
    //	// Power n and an n-th root degree cancel each other
    //	return Fraction{numerator: f.numerator, denominator: f.denominator, nthRootDegree: 1}
    //}
    // Power of an n-th root equals the n-th root of the power
    // TODO: check if exponent / nthRootDegree is not a float
    numerator := powInt(f.numerator, exponent)
    denominator := powInt(f.denominator, exponent)
    return Fraction{
        numerator:   numerator,
        denominator: denominator}
}

func (f Fraction) NthRoot(degree uint) Fraction {
    //TODO: negative fractions can have uneven degrees, not even degrees
    numerator := f.numerator
    denominator := f.denominator
    if isNthRootInt(numerator, degree) && isNthRootInt(f.denominator, degree) {
        // The nth-roots of the numerator and the denominator are integers => process
        numerator = int(math.Round(math.Pow(float64(numerator), 1.0/float64(degree))))
        denominator = int(math.Round(math.Pow(float64(denominator), 1.0/float64(degree))))
        return Fraction{numerator: numerator, denominator: denominator}
    }
    return Fraction{numerator: numerator, denominator: f.denominator}
}

func (f Fraction) Gcd() int {
    // Euclidean algorithm for computing the greatest common divisor (gcd).
    // The gcd of two numbers is the largest positive integer that divides
    // both numbers without leaving a remainder.
    // The function is associative: for example, the gcd of three numbers
    // a, b, c is equal to: gcd(a, b, c) = gcd(a, gcd(b, c). And so on.
    a := absInt(f.numerator)
    b := absInt(f.denominator)

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

//func roundToDecimalPlaces(value float64, places int) float64 {
//	scale := math.Pow(10, float64(places))
//	return math.Round(value*scale) / scale
//}
