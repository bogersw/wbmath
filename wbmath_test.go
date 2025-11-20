package wbmath

import (
	"math"
	"testing"
)

func TestAbs(t *testing.T) {
	if got := Abs(-5); got != 5 {
		t.Fatalf("Abs(-5) = %v, want 5", got)
	}
	if got := Abs[int8](-8); got != 8 {
		t.Fatalf("Abs[int8](-8) = %v, want 8", got)
	}
	if got := Abs(-3.14); got != 3.14 {
		t.Fatalf("Abs(-3.14) = %v, want 3.14", got)
	}
}

func TestGcd(t *testing.T) {
	cases := []struct {
		a, b int
		want int
	}{
		{48, 18, 6},
		{0, 5, 5},
		{5, 0, 5},
		{0, 0, 0},
	}
	for _, c := range cases {
		if got := Gcd(c.a, c.b); got != c.want {
			t.Fatalf("Gcd(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestIsNthRootInt(t *testing.T) {
	if !IsNthRootInt(27, 3) {
		t.Fatalf("IsNthRootInt(27,3) = false, want true")
	}
	if IsNthRootInt(20, 2) {
		t.Fatalf("IsNthRootInt(20,2) = true, want false")
	}
}

func TestPowInt(t *testing.T) {
	if got := PowInt(2, 10); got != 1024 {
		t.Fatalf("PowInt(2,10) = %d, want 1024", got)
	}
	if got := PowInt(-2, 3); got != -8 {
		t.Fatalf("PowInt(-2,3) = %d, want -8", got)
	}
}

func TestRound(t *testing.T) {
	if got := Round[float64](2.3456, 2); got != 2.35 {
		t.Fatalf("Round(2.3456,2) = %v, want 2.35", got)
	}
	if got := Round[float32](1.2345, 3); got != float32(1.235) {
		t.Fatalf("Round[float32](1.2345,3) = %v, want 1.235", got)
	}
}

func TestIsInteger(t *testing.T) {
	if !IsInteger(float64(5.0)) {
		t.Fatalf("IsInteger(5.0) = false, want true")
	}
	if IsInteger(float64(5.1)) {
		t.Fatalf("IsInteger(5.1) = true, want false")
	}
	if !IsInteger(float32(3.0)) {
		t.Fatalf("IsInteger(float32(3.0)) = false, want true")
	}
	if IsInteger(math.NaN()) {
		t.Fatalf("IsInteger(NaN) = true, want false")
	}
	if IsInteger(math.Inf(1)) {
		t.Fatalf("IsInteger(Inf) = true, want false")
	}
}
