// Package fraction provides utilities for creating and manipulating
// rational numbers represented by a Fraction struct.
package fraction

import (
    "errors"
    "fmt"
    "math"
    "regexp"
    "strconv"
    "strings"

    "github.com/bogersw/wbmath"
)

// Fraction represents a rational number stored with non-negative numerator and
// denominator and a separate sign flag. Fields are unexported: use the package's
// constructors and methods to create and manipulate values.
//
// The fields "numerator" and "denominator" hold absolute (non-negative) integers.
// The field "sign" is -1 if the resulting value is negative: otherwise it is 1.

type Fraction struct {
    numerator   int
    denominator int
    sign        int
}

// New is a constructor function that takes two integer parameters - the
// numerator and the denominator, respectively - and returns a pointer to
// a Fraction struct and an error in case the denominator is zero.
func New(numerator, denominator int) (*Fraction, error) {
    if denominator == 0 {
        return nil, errors.New("division by zero")
    }
    sign := 1
    if numerator*denominator < 0 {
        sign = -1
    }
    return &Fraction{
            numerator:   wbmath.Abs(numerator),
            denominator: wbmath.Abs(denominator),
            sign:        sign},
        nil
}

// MustNew is a constructor identical to New but which panics if an error
// occurs.
func MustNew(numerator, denominator int) *Fraction {
    fraction, err := New(numerator, denominator)
    if err != nil {
        panic(err)
    }
    return fraction
}

// NewFromNumber is a constructor function that an integer or a floating
// point number and turns it into a rational number. It returns a pointer
// to a Fraction struct.
func NewFromNumber[T int | float64](num T) *Fraction {
    switch any(num).(type) {
    case int:
        return MustNew(int(num), 1)
    case float64:
        if float64(num) == 0 {
            return MustNew(int(num), 1)
        }
        s := strconv.FormatFloat(float64(num), 'f', -1, 64)
        var decimalPlaces uint = 0
        if i := strings.IndexByte(s, '.'); i >= 1 {
            decimalPlaces = uint(len(s) - i - 1)
        }
        numerator, _ := strconv.ParseInt(
            strings.Replace(s, ".", "", 1),
            10,
            64)
        return MustNew(int(numerator), wbmath.PowInt(10, decimalPlaces))
    default:
        // Shouldn't happen
        return nil
    }
}

// NewFromString is a constructor function that accepts strings like
// "a / b", with a and b either ints or floats (including scientific
// notation (e / E)). Optional signs can be provided. Whitespace is
// ignored. It returns a Fraction struct and an error.
func NewFromString(num string) (*Fraction, error) {

    // Numbers can be integers or floats (the last ones with or without leading digits).
    // The numbers can have an optional sign and optional scientific exponent (e / E).
    var numPart = `(?:[+\-]?(?:\d+\.?\d*|\.\d+)(?:[eE][+\-]?\d+)?)`
    // ^   — start of string anchor (match begins at string start).
    // \s* — zero or more whitespace characters (allows leading and /ot trailing spaces).
    // /   — literal slash separator.
    // $   — end of string anchor (ensures the entire string matches, no extra chars).
    var re = regexp.MustCompile(fmt.Sprintf(`^\s*(%s)\s*/\s*(%s)\s*$`, numPart, numPart))
    // Check match: FindStringSubmatch returns a slice (or nil if there was no match)
    // - index 0 is the full match,
    // - index 1 is match 1 (in our case: the numerator),
    // - index 2 is match 2 (in our case: the denominator).
    match := re.FindStringSubmatch(strings.TrimSpace(num))
    if match == nil {
        return nil, errors.New("invalid fraction format")
    }
    numeratorStr, denominatorStr := match[1], match[2]
    // If both numbers have no decimals/exponent, treat them as integers
    if !strings.ContainsAny(numeratorStr, ".eE") && !strings.ContainsAny(denominatorStr, ".eE") {
        if numerator, err := strconv.Atoi(numeratorStr); err != nil {
            return nil, err
        } else {
            if denominator, err := strconv.Atoi(denominatorStr); err != nil {
                return nil, err
            } else {
                return New(numerator, denominator)
            }
        }
    }
    // Otherwise: parse as floats.
    if numerator, err := strconv.ParseFloat(numeratorStr, 64); err != nil {
        return nil, err
    } else {
        if denominator, err := strconv.ParseFloat(denominatorStr, 64); err != nil {
            return nil, err
        } else {
            fracNumerator := NewFromNumber(numerator)
            fracDenominator := NewFromNumber(denominator)
            result := fracNumerator.Divide(fracDenominator).Simplify()
            return result, nil
        }
    }
}

// MustNewFromString is a constructor identical to NewFromString but which
// panics if an error occurs.
func MustNewFromString(num string) *Fraction {
    fraction, err := NewFromString(num)
    if err != nil {
        panic(err)
    }
    return fraction
}

// Simplify determines the greatest common divisor (gcd) to make the fraction as
// simple as possible. Changes the current Fraction instance in-place and returns
// nil if the Fraction instance is nil. Note that if gcd = 0 the current Fraction
// is not changed.
func (f *Fraction) Simplify() *Fraction {
    if f == nil {
        return nil
    }
    gcd := wbmath.Gcd(f.numerator, f.denominator)
    if gcd != 0 {
        f.numerator = f.numerator / gcd
        f.denominator = f.denominator / gcd
    }
    return f
}

// Evaluate calculates and returns the fraction as a float value. Returns
// NaN if the Fraction instance is nil.
func (f *Fraction) Evaluate() float64 {
    if f == nil {
        return math.NaN()
    }
    if f.sign == -1 && f.numerator != 0 {
        return -float64(f.numerator) / float64(f.denominator)
    }
    return float64(f.numerator) / float64(f.denominator)
}

// String implements the fmt.Stringer interface and returns a string
// with a nicely formatted fraction for use by the fmt package.
func (f *Fraction) String() string {
    // Guard against funky input: String() can be called on nil pointers.
    if f == nil || f.denominator == 0 {
        return "NaN"
    }
    // Check the kind of fraction we're dealing with and if it can
    // be simplified to a whole number + a remaining fraction.
    numerator := f.numerator
    denominator := f.denominator
    wholeNumber := 0
    if numerator >= denominator {
        wholeNumber = numerator / denominator
        numerator = numerator % denominator
    }
    // Format the output string
    var result string
    if numerator == 0 {
        // Whole number only
        result = fmt.Sprintf("%d", wholeNumber)
    } else if wholeNumber == 0 {
        // Fraction only
        result = fmt.Sprintf("%d/%d", numerator, denominator)
    } else {
        // Whole number + fraction
        result = fmt.Sprintf("%d %d/%d", wholeNumber, numerator, denominator)
    }
    // Return result with the correct sign
    if f.sign == -1 {
        return fmt.Sprintf("-%s", result)
    }
    return fmt.Sprintf("%s", result)
}

// Multiply multiplies the current Fraction instance with the specified
// Fraction instance. Modifies the current Fraction instance in-place.
// Returns nil if either Fraction instance is nil.
func (f *Fraction) Multiply(other *Fraction) *Fraction {
    if f == nil || other == nil {
        return nil
    }
    f.numerator = f.numerator * other.numerator
    f.denominator = f.denominator * other.denominator
    f.sign = f.sign * other.sign
    return f
}

// MultiplyInt multiplies the current Fraction instance with the specified
// integer. Returns nil if the Fraction instance is nil.
func (f *Fraction) MultiplyInt(value int) *Fraction {
    if f == nil {
        return nil
    }
    return f.Multiply(NewFromNumber(value))
}

// Add adds the specified Fraction instance to the current Fraction instance.
// Modifies the current Fraction instance in-place. Returns nil if either
// Fraction instance is nil.
func (f *Fraction) Add(other *Fraction) *Fraction {
    if f == nil || other == nil {
        return nil
    }
    numerator := f.sign*f.numerator*other.denominator + other.sign*other.numerator*f.denominator
    f.numerator = wbmath.Abs(numerator)
    f.denominator = f.denominator * other.denominator
    if numerator >= 0 {
        f.sign = 1
    } else {
        f.sign = -1
    }
    return f
}

// AddInt adds the specified integer to the current Fraction instance.
// Returns nil if the Fraction instance is nil.
func (f *Fraction) AddInt(value int) *Fraction {
    if f == nil {
        return nil
    }
    return f.Add(NewFromNumber(value))
}

// Divide divides the current Fraction instance with the specified Fraction
// instance. Modifies the current Fraction instance in-place. Returns nil if
// either Fraction instance is nil.
func (f *Fraction) Divide(other *Fraction) *Fraction {
    if f == nil || other == nil {
        return nil
    }
    f.numerator = f.numerator * other.denominator
    f.denominator = f.denominator * other.numerator
    f.sign = f.sign / other.sign
    return f
}

// DivideInt divides the current Fraction instance with the specified integer.
// Modifies the current Fraction instance in-place and returns it (or returns
// nil if the Fraction instance is nil). Also returns an error (which is nil
// if no error occurs).
func (f *Fraction) DivideInt(value int) (*Fraction, error) {
    if f == nil {
        return nil, errors.New("invalid Fraction instance")
    }
    if value == 0 {
        return nil, errors.New("division by zero")
    }
    return f.Divide(NewFromNumber(value)), nil
}

// MustDivideInt is identical to DivideInt, but it panics if an error occurs.
func (f *Fraction) MustDivideInt(value int) *Fraction {
    if _, err := f.DivideInt(value); err != nil {
        panic(err)
    }
    return f
}

// Subtract subtracts the specified Fraction instance from the current Fraction
// instance. Modifies the current Fraction instance in-place. Returns nil if
// either Fraction instance is nil.
func (f *Fraction) Subtract(other *Fraction) *Fraction {
    // Subtracts two Fractions and returns a new Fraction with the result
    if f == nil || other == nil {
        return nil
    }
    if other.sign == -1 {
        return f.Add(MustNew(other.numerator, other.denominator))
    }
    return f.Add(MustNew(-other.numerator, other.denominator))
}

// SubtractInt subtracts the specified integer from the current Fraction
// instance. Returns nil if the Fraction instance is nil.
func (f *Fraction) SubtractInt(value int) *Fraction {
    if f == nil {
        return nil
    }
    return f.Subtract(NewFromNumber(value))
}

// Pow raises the current Fraction instance to the specified power. Modifies
// the current Fraction instance in-place and returns it (or returns
// nil if the Fraction instance is nil).
func (f *Fraction) Pow(exponent uint) *Fraction {
    if f == nil {
        return nil
    }
    f.numerator = wbmath.PowInt(f.numerator, exponent)
    f.denominator = wbmath.PowInt(f.denominator, exponent)
    // For uneven powers a negative sign is preserved
    f.sign = wbmath.PowInt(f.sign, exponent)
    return f
}

// NthRoot determines the nth-root of the current Fraction instance. Modifies
// the current Fraction instance in-place and returns it (or returns
// nil if the nth-root of the Fraction instance is non-existent). Returns an
// error (which is nil if no error occurs).
func (f *Fraction) NthRoot(degree uint) (*Fraction, error) {
    if f == nil {
        return nil, errors.New("invalid Fraction instance")
    }
    if f.sign == -1 && degree%2 == 0 {
        return nil, errors.New("the even nth-root of a negative number does not exist")
    }
    if wbmath.IsNthRootInt(f.numerator, degree) && wbmath.IsNthRootInt(f.denominator, degree) {
        // The nth-roots of the numerator and the denominator are integers => process
        f.numerator = int(math.Round(math.Pow(float64(f.numerator), 1.0/float64(degree))))
        f.denominator = int(math.Round(math.Pow(float64(f.denominator), 1.0/float64(degree))))
        return f, nil
    } else {
        // The nth-roots do not yield integers => invalid Fraction
        return nil, errors.New("the nth-root of this fraction does not yield a valid fraction")
    }
}

// MustNthRoot is identical to NthRoot, but it panics if an error occurs.
func (f *Fraction) MustNthRoot(degree uint) *Fraction {
    if _, err := f.NthRoot(degree); err != nil {
        panic(err)
    }
    return f
}

// Numerator returns the numerator of the current Fraction instance. Note that
// if the fraction is negative, the returned value for the numerator will be
// negative. Returns the numerator value and a boolean value that indicates if
// the returned numerator is valid.
func (f *Fraction) Numerator() (int, bool) {
    if f == nil {
        return 0, false
    }
    return f.sign * f.numerator, true
}

// Denominator returns the denominator of the current Fraction instance.
// Returns the denominator value and a boolean value that indicates if
// the returned denominator is valid.
func (f *Fraction) Denominator() (int, bool) {
    if f == nil {
        return 0, false
    }
    return f.denominator, true
}

// AsIntegerRatio returns the string representation of the Fraction instance
// as an integer ratio [-]a/b. If the Fraction instance is nil it will return
// NaN.
func (f *Fraction) AsIntegerRatio() string {
    if f == nil {
        return "NaN"
    }
    var result string
    result = fmt.Sprintf("%d/%d", f.numerator, f.denominator)
    if f.sign == -1 {
        return fmt.Sprintf("-%s", result)
    }
    return fmt.Sprintf("%s", result)
}
