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

// A Complex represents a multi-precision floating-point complex number.
type Complex struct {
	l, r big.Float
}

// Real returns the real part of z.
func (z *Complex) Real() *big.Float {
	return &z.l
}

// Cartesian returns the two cartesian components of z.
func (z *Complex) Cartesian() (*big.Float, *big.Float) {
	return &z.l, &z.r
}

// String returns the string version of a Complex value.
//
// If z corresponds to a + bi, then the string is "(a+bi)", similar to
// complex128 values.
func (z *Complex) String() string {
	a := make([]string, 5)
	a[0] = "("
	a[1] = fmt.Sprintf("%v", &z.l)
	if z.r.Signbit() {
		a[2] = fmt.Sprintf("%v", &z.r)
	} else {
		a[2] = fmt.Sprintf("+%v", &z.r)
	}
	a[3] = "i"
	a[4] = ")"
	return strings.Join(a, "")
}

// Equals returns true if y and z are equal.
func (z *Complex) Equals(y *Complex) bool {
	if z.l.Cmp(&y.l) != 0 || z.r.Cmp(&y.r) != 0 {
		return false
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Complex) Copy(y *Complex) *Complex {
	z.l.Copy(&y.l)
	z.r.Copy(&y.r)
	return z
}

// NewComplex returns a pointer to the Complex value a+bi.
func NewComplex(a, b *big.Float) *Complex {
	z := new(Complex)
	z.l.Copy(a)
	z.r.Copy(b)
	return z
}

// Scal sets z equal to y scaled by a, and returns z.
func (z *Complex) Scal(y *Complex, a *big.Float) *Complex {
	z.l.Mul(&y.l, a)
	z.r.Mul(&y.r, a)
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Complex) Neg(y *Complex) *Complex {
	z.l.Neg(&y.l)
	z.r.Neg(&y.r)
	return z
}

// Conj sets z equal to the conjugate of y, and returns z.
func (z *Complex) Conj(y *Complex) *Complex {
	z.l.Copy(&y.l)
	z.r.Neg(&y.r)
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Complex) Add(x, y *Complex) *Complex {
	z.l.Add(&x.l, &y.l)
	z.r.Add(&x.r, &y.r)
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Complex) Sub(x, y *Complex) *Complex {
	z.l.Sub(&x.l, &y.l)
	z.r.Sub(&x.r, &y.r)
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
//
// The multiplication rule is:
// 		Mul(i, i) = -1
// This binary operation is commutative and associative.
func (z *Complex) Mul(x, y *Complex) *Complex {
	a := new(big.Float).Copy(&x.l)
	b := new(big.Float).Copy(&x.r)
	c := new(big.Float).Copy(&y.l)
	d := new(big.Float).Copy(&y.r)
	temp := new(big.Float)
	z.l.Sub(
		z.l.Mul(a, c),
		temp.Mul(d, b),
	)
	z.r.Add(
		z.r.Mul(d, a),
		temp.Mul(b, c),
	)
	return z
}

// Quad returns the quadrance of z, a pointer to a big.Float value.
func (z *Complex) Quad() *big.Float {
	quad := new(big.Float)
	return quad.Add(
		quad.Mul(&z.l, &z.l),
		new(big.Float).Mul(&z.r, &z.r),
	)
}

// Inv sets z equal to the inverse of y, and returns z.
func (z *Complex) Inv(y *Complex) *Complex {
	zero := new(Complex)
	if y.Equals(zero) {
		panic("zero inverse")
	}
	quad := y.Quad()
	z.Conj(y)
	z.l.Quo(&z.l, quad)
	z.r.Quo(&z.r, quad)
	return z
}

// Quo sets z equal to the quotient of x and y, and returns z.
func (z *Complex) Quo(x, y *Complex) *Complex {
	zero := new(Complex)
	if y.Equals(zero) {
		panic("zero denominator")
	}
	quad := y.Quad()
	z.Conj(y)
	z.Mul(x, z)
	z.l.Quo(&z.l, quad)
	z.r.Quo(&z.r, quad)
	return z
}

// CrossRatio sets z equal to the cross ratio
// 		Inv(w - x) * (v - x) * Inv(v - y) * (w - y)
// Then it returns z.
func (z *Complex) CrossRatio(v, w, x, y *Complex) *Complex {
	temp := new(Complex)
	z.Sub(w, x)
	z.Inv(z)
	temp.Sub(v, x)
	z.Mul(z, temp)
	temp.Sub(v, y)
	temp.Inv(temp)
	z.Mul(z, temp)
	temp.Sub(w, y)
	z.Mul(z, temp)
	return z
}

// Möbius sets z equal to the Möbius (fractional linear) transform
// 		(a*y + b) * Inv(c*y + d)
// Then it returns z.
func (z *Complex) Möbius(y, a, b, c, d *Complex) *Complex {
	z.Mul(a, y)
	z.Add(z, b)
	temp := new(Complex)
	temp.Mul(c, y)
	temp.Add(temp, d)
	temp.Inv(temp)
	z.Mul(z, temp)
	return z
}

// Generate returns a random Complex value for quick.Check testing.
func (z *Complex) Generate(rand *rand.Rand, size int) reflect.Value {
	randomComplex := &Complex{
		*big.NewFloat(rand.Float64()),
		*big.NewFloat(rand.Float64()),
	}
	return reflect.ValueOf(randomComplex)
}
