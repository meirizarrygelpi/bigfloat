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

// A Perplex represents a multi-precision floating-point perplex number.
type Perplex struct {
	l, r big.Float
}

// Cartesian returns the two cartesian components of z.
func (z *Perplex) Cartesian() (*big.Float, *big.Float) {
	return &z.l, &z.r
}

// String returns the string version of a Perplex value.
//
// If z corresponds to a + bs, then the string is "(a+bs)", similar to
// complex128 values.
func (z *Perplex) String() string {
	a := make([]string, 5)
	a[0] = "("
	a[1] = fmt.Sprintf("%v", &z.l)
	if z.r.Signbit() {
		a[2] = fmt.Sprintf("%v", &z.r)
	} else {
		a[2] = fmt.Sprintf("+%v", &z.r)
	}
	a[3] = "s"
	a[4] = ")"
	return strings.Join(a, "")
}

// Equals returns true if y and z are equal.
func (z *Perplex) Equals(y *Perplex) bool {
	if z.l.Cmp(&y.l) != 0 || z.r.Cmp(&y.r) != 0 {
		return false
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Perplex) Copy(y *Perplex) *Perplex {
	z.l.Copy(&y.l)
	z.r.Copy(&y.r)
	return z
}

// NewPerplex returns a pointer to the Perplex value a+bs.
func NewPerplex(a, b *big.Float) *Perplex {
	z := new(Perplex)
	z.l.Copy(a)
	z.r.Copy(b)
	return z
}

// Scal sets z equal to y scaled by a, and returns z.
func (z *Perplex) Scal(y *Perplex, a *big.Float) *Perplex {
	z.l.Mul(&y.l, a)
	z.r.Mul(&y.r, a)
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Perplex) Neg(y *Perplex) *Perplex {
	z.l.Neg(&y.l)
	z.r.Neg(&y.r)
	return z
}

// Conj sets z equal to the conjugate of y, and returns z.
func (z *Perplex) Conj(y *Perplex) *Perplex {
	z.l.Copy(&y.l)
	z.r.Neg(&y.r)
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Perplex) Add(x, y *Perplex) *Perplex {
	z.l.Add(&x.l, &y.l)
	z.r.Add(&x.r, &y.r)
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Perplex) Sub(x, y *Perplex) *Perplex {
	z.l.Sub(&x.l, &y.l)
	z.r.Sub(&x.r, &y.r)
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
//
// The multiplication rule is:
// 		Mul(s, s) = +1
// This binary operation is commutative and associative.
func (z *Perplex) Mul(x, y *Perplex) *Perplex {
	a := new(big.Float).Copy(&x.l)
	b := new(big.Float).Copy(&x.r)
	c := new(big.Float).Copy(&y.l)
	d := new(big.Float).Copy(&y.r)
	temp := new(big.Float)
	z.l.Add(
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
func (z *Perplex) Quad() *big.Float {
	quad := new(big.Float)
	return quad.Sub(
		quad.Mul(&z.l, &z.l),
		new(big.Float).Mul(&z.r, &z.r),
	)
}

// IsZeroDiv returns true if z is a zero divisor.
func (z *Perplex) IsZeroDiv() bool {
	if z.l.Cmp(&z.r) == 0 {
		return true
	}
	if z.l.Cmp(new(big.Float).Neg(&z.r)) == 0 {
		return true
	}
	return false
}

// Inv sets z equal to the inverse of y, and returns z.
func (z *Perplex) Inv(y *Perplex) *Perplex {
	if y.IsZeroDiv() {
		panic("zero divisor inverse")
	}
	quad := y.Quad()
	z.Conj(y)
	z.l.Quo(&z.l, quad)
	z.r.Quo(&z.r, quad)
	return z
}

// Quo sets z equal to the quotient of x and y, and returns z.
func (z *Perplex) Quo(x, y *Perplex) *Perplex {
	if y.IsZeroDiv() {
		panic("zero divisor denominator")
	}
	quad := y.Quad()
	z.Conj(y)
	z.Mul(x, z)
	z.l.Quo(&z.l, quad)
	z.r.Quo(&z.r, quad)
	return z
}

// Idempotent sets z equal to a pointer to an idempotent Perplex.
func (z *Perplex) Idempotent(sign int) *Perplex {
	z.l.SetFloat64(0.5)
	if sign < 0 {
		z.r.SetFloat64(-0.5)
		return z
	}
	z.r.SetFloat64(0.5)
	return z
}

// CrossRatio sets z equal to the cross ratio
// 		Inv(w - x) * (v - x) * Inv(v - y) * (w - y)
// Then it returns z.
func (z *Perplex) CrossRatio(v, w, x, y *Perplex) *Perplex {
	temp := new(Perplex)
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
func (z *Perplex) Möbius(y, a, b, c, d *Perplex) *Perplex {
	z.Mul(a, y)
	z.Add(z, b)
	temp := new(Perplex)
	temp.Mul(c, y)
	temp.Add(temp, d)
	temp.Inv(temp)
	z.Mul(z, temp)
	return z
}

// Generate returns a random Perplex value for quick.Check testing.
func (z *Perplex) Generate(rand *rand.Rand, size int) reflect.Value {
	randomPerplex := &Perplex{
		*big.NewFloat(rand.Float64()),
		*big.NewFloat(rand.Float64()),
	}
	return reflect.ValueOf(randomPerplex)
}
