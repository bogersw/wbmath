package fraction

import (
    "errors"
    "fmt"
    "math"
    "strconv"
    "strings"

    "github.com/bogersw/wbmath"
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
            numerator:   wbmath.Abs(numerator),
            denominator: wbmath.Abs(denominator),
            isNegative:  numerator*denominator < 0},
        nil
}

func NewFromNumber[T int | float64](num T) Fraction {
    switch any(num).(type) {
    case int:
        return Fraction{numerator: int(num), denominator: 1}
    case float64:
        s := strconv.FormatFloat(float64(num), 'f', -1, 64)
        var decimalPlaces uint = 0
        if i := strings.IndexByte(s, '.'); i >= 1 {
            decimalPlaces = uint(len(s) - i - 1)
        }
        return Fraction{
            numerator:   3,
            denominator: wbmath.PowInt(10, decimalPlaces),
        }
    default:
        return Fraction{}
    }
}

func (f Fraction) Simplify() Fraction {

    gcd := wbmath.Gcd(f.numerator, f.denominator)
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
        numerator:   wbmath.Abs(value),
        denominator: 1,
        isNegative:  value < 0})

}

func (f Fraction) Add(other Fraction) Fraction {
    // Adds two Fractions and returns a new Fraction with the result
    signSelf := 1
    if f.isNegative {
        signSelf = -1
    }
    signOther := 1
    if other.isNegative {
        signOther = -1
    }
    numerator := signSelf*f.numerator*other.denominator + signOther*other.numerator*f.denominator
    denominator := f.denominator * other.denominator
    return Fraction{
        numerator:   wbmath.Abs(numerator),
        denominator: denominator,
        isNegative:  numerator*denominator < 0}
}

func (f Fraction) AddInt(value int) Fraction {
    return f.Add(Fraction{
        numerator:   wbmath.Abs(value),
        denominator: 1,
        isNegative:  value < 0})
}

func (f Fraction) Subtract(other Fraction) Fraction {
    // Subtracts two Fractions and returns a new Fraction with the result
    signSelf := 1
    if f.isNegative {
        signSelf = -1
    }
    signOther := 1
    if other.isNegative {
        signOther = -1
    }
    numerator := signSelf*f.numerator*other.denominator - signOther*other.numerator*f.denominator
    denominator := f.denominator * other.denominator
    return Fraction{
        numerator:   wbmath.Abs(numerator),
        denominator: denominator,
        isNegative:  numerator*denominator < 0}
}

func (f Fraction) SubtractInt(value int) Fraction {
    return f.Subtract(Fraction{
        numerator:   wbmath.Abs(value),
        denominator: 1,
        isNegative:  value < 0})
}

// Divide divides a fraction by another fraction.
func (f Fraction) Divide(other Fraction) Fraction {
    // Divides two Fractions and returns a new Fraction with the result
    isNegative := false
    if f.isNegative && !other.isNegative || !f.isNegative && other.isNegative {
        isNegative = true
    }
    return Fraction{
        numerator:   f.numerator * other.denominator,
        denominator: f.denominator * other.numerator,
        isNegative:  isNegative}
}

// DivideInt divides a fraction by an integer value. It returns a new Fraction
// and an error (which is nil if no error occurs).
func (f Fraction) DivideInt(value int) (Fraction, error) {
    if value == 0 {
        return Fraction{}, errors.New("error: integer divide by zero")
    }
    return f.Divide(Fraction{
            numerator:   wbmath.Abs(value),
            denominator: 1,
            isNegative:  value < 0}),
        nil
}

func (f Fraction) Pow(exponent uint) Fraction {
    numerator := wbmath.PowInt(f.numerator, exponent)
    denominator := wbmath.PowInt(f.denominator, exponent)
    // Determine sign: for uneven powers a negative sign is preserved
    isNegative := false
    if f.isNegative && exponent%2 != 0 {
        isNegative = true
    }
    return Fraction{
        numerator:   numerator,
        denominator: denominator,
        isNegative:  isNegative}
}

func (f Fraction) NthRoot(degree uint) Fraction {
    //TODO: negative fractions can have uneven degrees, not even degrees
    numerator := f.numerator
    denominator := f.denominator
    if wbmath.IsNthRootInt(numerator, degree) && wbmath.IsNthRootInt(f.denominator, degree) {
        // The nth-roots of the numerator and the denominator are integers => process
        numerator = int(math.Round(math.Pow(float64(numerator), 1.0/float64(degree))))
        denominator = int(math.Round(math.Pow(float64(denominator), 1.0/float64(degree))))
        return Fraction{numerator: numerator, denominator: denominator}
    }
    return Fraction{numerator: numerator, denominator: f.denominator}
}
