package fraction

import (
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	const eps = 1e-9
	return math.Abs(a-b) <= eps
}

func TestNewAndAccessors(t *testing.T) {
	f, err := New(3, 4)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if n, ok := f.Numerator(); !ok || n != 3 {
		t.Fatalf("Numerator() = %v, %v; want 3,true", n, ok)
	}
	if d, ok := f.Denominator(); !ok || d != 4 {
		t.Fatalf("Denominator() = %v, %v; want 4,true", d, ok)
	}

	if _, err := New(1, 0); err == nil {
		t.Fatalf("New with denominator 0 should return error")
	}
}

func TestNewFromNumberAndString(t *testing.T) {
	f1 := NewFromNumber(5)
	if v := f1.Evaluate(); !almostEqual(v, 5.0) {
		t.Fatalf("NewFromNumber(int) Evaluate = %v; want 5", v)
	}

	f2 := NewFromNumber(0.125)
	if v := f2.Evaluate(); !almostEqual(v, 0.125) {
		t.Fatalf("NewFromNumber(float) Evaluate = %v; want 0.125", v)
	}

	// From string with ints
	fs, err := NewFromString("-3/4")
	if err != nil {
		t.Fatalf("NewFromString returned error: %v", err)
	}
	if v := fs.Evaluate(); !almostEqual(v, -0.75) {
		t.Fatalf("NewFromString Evaluate = %v; want -0.75", v)
	}

	// From string with floats
	fs2, err := NewFromString("1.5 / 0.5")
	if err != nil {
		t.Fatalf("NewFromString returned error: %v", err)
	}
	if v := fs2.Evaluate(); !almostEqual(v, 3.0) {
		t.Fatalf("NewFromString Evaluate = %v; want 3", v)
	}
}

func TestSimplifyAndString(t *testing.T) {
	f, _ := New(2, 4)
	f = f.Simplify()
	if n, _ := f.Numerator(); n != 1 {
		t.Fatalf("Simplified numerator = %v; want 1", n)
	}
	if s := f.String(); s != "1/2" {
		t.Fatalf("String() = %q; want \"1/2\"", s)
	}

	// Mixed number and negative
	f2, _ := New(7, 3) // 2 1/3
	if s := f2.String(); s != "2 1/3" {
		t.Fatalf("String() = %q; want \"2 1/3\"", s)
	}
	f3, _ := New(-4, 2)
	if s := f3.String(); s != "-2" {
		t.Fatalf("String() = %q; want \"-2\"", s)
	}
}

func TestArithmeticOperations(t *testing.T) {
	// Multiply
	a, _ := New(1, 2)
	b, _ := New(3, 4)
	a.Multiply(b)
	if v := a.Evaluate(); !almostEqual(v, 0.375) {
		t.Fatalf("Multiply Evaluate = %v; want 0.375", v)
	}

	// Add (with sign handling)
	a, _ = New(1, 2)  // 1/2
	b, _ = New(-3, 4) // -3/4
	a.Add(b)          // result -1/4 (internally -2/8)
	if v := a.Evaluate(); !almostEqual(v, -0.25) {
		t.Fatalf("Add Evaluate = %v; want -0.25", v)
	}
	a = a.Simplify()
	if s := a.String(); s != "-1/4" {
		t.Fatalf("After Simplify String() = %q; want \"-1/4\"", s)
	}

	// AddInt / SubtractInt
	x, _ := New(3, 2) // 1 1/2
	x.AddInt(1)       // 2 1/2
	if !almostEqual(x.Evaluate(), 2.5) {
		t.Fatalf("AddInt Evaluate = %v; want 2.5", x.Evaluate())
	}
	x.SubtractInt(1) // back to 1.5
	if !almostEqual(x.Evaluate(), 1.5) {
		t.Fatalf("SubtractInt Evaluate = %v; want 1.5", x.Evaluate())
	}

	// Divide and DivideInt
	d1, _ := New(3, 4)
	d2, _ := New(2, 3)
	d1.Divide(d2) // (3/4) / (2/3) = 9/8 = 1.125
	if !almostEqual(d1.Evaluate(), 1.125) {
		t.Fatalf("Divide Evaluate = %v; want 1.125", d1.Evaluate())
	}
	// DivideInt error on zero
	_, err := d1.DivideInt(0)
	if err == nil {
		t.Fatalf("DivideInt by zero should return error")
	}
}

func TestNthRootAndPow(t *testing.T) {
	// Pow
	p, _ := New(2, 3)
	p.Pow(2) // (2/3)^2 = 4/9
	if !almostEqual(p.Evaluate(), 4.0/9.0) {
		t.Fatalf("Pow Evaluate = %v; want %v", p.Evaluate(), 4.0/9.0)
	}

	// NthRoot success
	r, _ := New(4, 9)
	got, err := r.NthRoot(2)
	if err != nil {
		t.Fatalf("NthRoot returned error: %v", err)
	}
	if !almostEqual(got.Evaluate(), 2.0/3.0) {
		t.Fatalf("NthRoot Evaluate = %v; want %v", got.Evaluate(), 2.0/3.0)
	}

	// NthRoot failure (even root of negative)
	neg, _ := New(-1, 4)
	if _, err := neg.NthRoot(2); err == nil {
		t.Fatalf("NthRoot of negative (even) should return error")
	}
}

func TestAsIntegerRatioAndNilReceiver(t *testing.T) {
	f, _ := New(5, 6)
	if s := f.AsIntegerRatio(); s != "5/6" {
		t.Fatalf("AsIntegerRatio = %q; want \"5/6\"", s)
	}
	// Nil receiver behavior
	var nf *Fraction
	if nf.Simplify() != nil {
		t.Fatalf("nil.Simplify() should return nil")
	}
	if !math.IsNaN(nf.Evaluate()) {
		t.Fatalf("nil.Evaluate() should return NaN")
	}
	if s := nf.String(); s != "NaN" {
		t.Fatalf("nil.String() = %q; want \"NaN\"", s)
	}
	// Operations on nil should return nil
	if nf.Multiply(f) != nil {
		t.Fatalf("nil.Multiply(...) should return nil")
	}
	if nf.Add(f) != nil {
		t.Fatalf("nil.Add(...) should return nil")
	}
	if nf.Divide(f) != nil {
		t.Fatalf("nil.Divide(...) should return nil")
	}
}

func TestMustNewFromString(t *testing.T) {
	f := MustNewFromString("4/6")
	if f == nil {
		t.Fatalf("MustNewFromString returned nil")
	}
	if v := f.Evaluate(); !almostEqual(v, 2.0/3.0) {
		t.Fatalf("MustNewFromString Evaluate = %v; want %v", v, 2.0/3.0)
	}
}
