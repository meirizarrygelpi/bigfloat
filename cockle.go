// Copyright (c) 2016 Melvin Eloy Irizarry-Gelpí
// Licenced under the MIT License.

package bigfloat

import (
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"strings"
)

var symbCockle = [4]string{"", "i", "t", "u"}

// A Cockle represents a multi-precision floating-point Cockle quaternion.
type Cockle struct {
	l, r Complex
}

// Real returns the real part of z.
func (z *Cockle) Real() *big.Float {
	return (&z.l).Real()
}

// Cartesian returns the four multi-precision floating-point Cartesian
// components of z.
func (z *Cockle) Cartesian() (*big.Float, *big.Float, *big.Float, *big.Float) {
	return &z.l.l, &z.l.r, &z.r.l, &z.r.r
}

// String returns the string representation of a Cockle value.
//
// If z corresponds to a + bi + ct + du, then the string is "(a+bi+ct+du)",
// similar to complex128 values.
func (z *Cockle) String() string {
	v := make([]*big.Float, 4)
	v[0], v[1] = z.l.Cartesian()
	v[2], v[3] = z.r.Cartesian()
	a := make([]string, 9)
	a[0] = "("
	a[1] = fmt.Sprintf("%v", v[0])
	i := 1
	for j := 2; j < 8; j = j + 2 {
		if v[i].Sign() == -1 {
			a[j] = fmt.Sprintf("%v", v[i])
		} else {
			a[j] = fmt.Sprintf("+%v", v[i])
		}
		a[j+1] = symbCockle[i]
		i++
	}
	a[8] = ")"
	return strings.Join(a, "")
}

// Equals returns true if y and z are equal.
func (z *Cockle) Equals(y *Cockle) bool {
	if !z.l.Equals(&y.l) || !z.r.Equals(&y.r) {
		return false
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Cockle) Copy(y *Cockle) *Cockle {
	z.l.Copy(&y.l)
	z.r.Copy(&y.r)
	return z
}

// NewCockle returns a pointer to the Cockle value a+bi+ct+du.
func NewCockle(a, b, c, d *big.Float) *Cockle {
	z := new(Cockle)
	z.l.l.Copy(a)
	z.l.r.Copy(b)
	z.r.l.Copy(c)
	z.r.r.Copy(d)
	return z
}

// Scal sets z equal to y scaled by a, and returns z.
func (z *Cockle) Scal(y *Cockle, a *big.Float) *Cockle {
	z.l.Scal(&y.l, a)
	z.r.Scal(&y.r, a)
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Cockle) Neg(y *Cockle) *Cockle {
	z.l.Neg(&y.l)
	z.r.Neg(&y.r)
	return z
}

// Conj sets z equal to the conjugate of y, and returns z.
func (z *Cockle) Conj(y *Cockle) *Cockle {
	z.l.Conj(&y.l)
	z.r.Neg(&y.r)
	return z
}

// Add sets z equal to x+y, and returns z.
func (z *Cockle) Add(x, y *Cockle) *Cockle {
	z.l.Add(&x.l, &y.l)
	z.r.Add(&x.r, &y.r)
	return z
}

// Sub sets z equal to x-y, and returns z.
func (z *Cockle) Sub(x, y *Cockle) *Cockle {
	z.l.Sub(&x.l, &y.l)
	z.r.Sub(&x.r, &y.r)
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
//
// The multiplication rules are:
// 		Mul(i, i) = -1
// 		Mul(t, t) = Mul(u, u) = +1
// 		Mul(i, t) = -Mul(t, i) = u
// 		Mul(u, t) = -Mul(t, u) = i
// 		Mul(u, i) = -Mul(i, u) = t
// This binary operation is noncommutative but associative.
func (z *Cockle) Mul(x, y *Cockle) *Cockle {
	a := new(Complex).Copy(&x.l)
	b := new(Complex).Copy(&x.r)
	c := new(Complex).Copy(&y.l)
	d := new(Complex).Copy(&y.r)
	temp := new(Complex)
	z.l.Add(
		z.l.Mul(a, c),
		temp.Mul(temp.Conj(d), b),
	)
	z.r.Add(
		z.r.Mul(d, a),
		temp.Mul(b, temp.Conj(c)),
	)
	return z
}

// Commutator sets z equal to the commutator of x and y
// 		Mul(x, y) - Mul(y, x)
// Then it returns z.
func (z *Cockle) Commutator(x, y *Cockle) *Cockle {
	return z.Sub(
		z.Mul(x, y),
		new(Cockle).Mul(y, x),
	)
}

// Quad returns the quadrance of z. If z = a+bi+ct+du, then the quadrance is
// 		Mul(a, a) + Mul(b, b) - Mul(c, c) - Mul(d, d)
// This can be positive, negative, or zero.
func (z *Cockle) Quad() *big.Float {
	return new(big.Float).Sub(
		z.l.Quad(),
		z.r.Quad(),
	)
}

// IsZeroDiv returns true if z is a zero divisor.
func (z *Cockle) IsZeroDiv() bool {
	return z.l.Quad().Cmp(z.r.Quad()) == 0
}

// Inv sets z equal to the inverse of y, and returns z. If y is a zero divisor,
// then Inv panics.
func (z *Cockle) Inv(y *Cockle) *Cockle {
	if y.IsZeroDiv() {
		panic("inverse of zero divisor")
	}
	quad := y.Quad()
	z.Conj(y)
	z.l.l.Quo(&z.l.l, quad)
	z.l.r.Quo(&z.l.r, quad)
	z.r.l.Quo(&z.r.l, quad)
	z.r.r.Quo(&z.r.r, quad)
	return z
}

// QuoL sets z equal to the left quotient of x and y:
// 		Mul(Inv(y), x)
// Then it returns z. If y is a zero divisor, then QuoL panics.
func (z *Cockle) QuoL(x, y *Cockle) *Cockle {
	if y.IsZeroDiv() {
		panic("denominator is zero divisor")
	}
	quad := y.Quad()
	z.Conj(y)
	z.Mul(z, x)
	z.l.l.Quo(&z.l.l, quad)
	z.l.r.Quo(&z.l.r, quad)
	z.r.l.Quo(&z.r.l, quad)
	z.r.r.Quo(&z.r.r, quad)
	return z
}

// QuoR sets z equal to the right quotient of x and y:
// 		Mul(x, Inv(y))
// Then it returns z. If y is a zero divisor, then QuoR panics.
func (z *Cockle) QuoR(x, y *Cockle) *Cockle {
	if y.IsZeroDiv() {
		panic("denominator is zero divisor")
	}
	quad := y.Quad()
	z.Conj(y)
	z.Mul(x, z)
	z.l.l.Quo(&z.l.l, quad)
	z.l.r.Quo(&z.l.r, quad)
	z.r.l.Quo(&z.r.l, quad)
	z.r.r.Quo(&z.r.r, quad)
	return z
}

// CrossRatioL sets z equal to the left cross-ratio of v, w, x, and y:
// 		Inv(w - x) * (v - x) * Inv(v - y) * (w - y)
// Then it returns z.
func (z *Cockle) CrossRatioL(v, w, x, y *Cockle) *Cockle {
	temp := new(Cockle)
	z.Sub(w, x)
	z.Inv(z)
	temp.Sub(v, x)
	z.Mul(z, temp)
	temp.Sub(v, y)
	temp.Inv(temp)
	z.Mul(z, temp)
	temp.Sub(w, y)
	return z.Mul(z, temp)
}

// CrossRatioR sets z equal to the right cross-ratio of v, w, x, and y:
// 		(v - x) * Inv(w - x) * (w - y) * Inv(v - y)
// Then it returns z.
func (z *Cockle) CrossRatioR(v, w, x, y *Cockle) *Cockle {
	temp := new(Cockle)
	z.Sub(v, x)
	temp.Sub(w, x)
	temp.Inv(temp)
	z.Mul(z, temp)
	temp.Sub(w, y)
	z.Mul(z, temp)
	temp.Sub(v, y)
	temp.Inv(temp)
	return z.Mul(z, temp)
}

// MöbiusL sets z equal to the left Möbius (fractional linear) transform of y:
// 		Inv(y*c + d) * (y*a + b)
// Then it returns z.
func (z *Cockle) MöbiusL(y, a, b, c, d *Cockle) *Cockle {
	z.Mul(y, a)
	z.Add(z, b)
	temp := new(Cockle)
	temp.Mul(y, c)
	temp.Add(temp, d)
	temp.Inv(temp)
	return z.Mul(temp, z)
}

// MöbiusR sets z equal to the right Möbius (fractional linear) transform of y:
// 		(a*y + b) * Inv(c*y + d)
// Then it returns z.
func (z *Cockle) MöbiusR(y, a, b, c, d *Cockle) *Cockle {
	z.Mul(a, y)
	z.Add(z, b)
	temp := new(Cockle)
	temp.Mul(c, y)
	temp.Add(temp, d)
	temp.Inv(temp)
	return z.Mul(z, temp)
}

// IsNilpotent returns true if z raised to the n-th power vanishes.
func (z *Cockle) IsNilpotent(n int) bool {
	zero := new(Cockle)
	zeroFloat := new(big.Float)
	if z.Equals(zero) {
		return true
	}
	p := NewCockle(big.NewFloat(1), zeroFloat, zeroFloat, zeroFloat)
	for i := 0; i < n; i++ {
		p.Mul(p, z)
		if p.Equals(zero) {
			return true
		}
	}
	return false
}

// Generate returns a random Cockle value for quick.Check testing.
func (z *Cockle) Generate(rand *rand.Rand, size int) reflect.Value {
	randomCockle := &Cockle{
		*NewComplex(
			big.NewFloat(rand.Float64()),
			big.NewFloat(rand.Float64()),
		),
		*NewComplex(
			big.NewFloat(rand.Float64()),
			big.NewFloat(rand.Float64()),
		),
	}
	return reflect.ValueOf(randomCockle)
}
