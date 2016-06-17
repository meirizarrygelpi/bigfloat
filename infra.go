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

// A Infra represents a multi-precision floating-point infra number.
type Infra struct {
	l, r big.Float
}

// Real returns the real part of z.
func (z *Infra) Real() *big.Float {
	return &z.l
}

// Cartesian returns the two cartesian components of z.
func (z *Infra) Cartesian() (*big.Float, *big.Float) {
	return &z.l, &z.r
}

// String returns the string version of a Infra value.
//
// If z corresponds to a + bα, then the string is "(a+bα)", similar to
// complex128 values.
func (z *Infra) String() string {
	a := make([]string, 5)
	a[0] = "("
	a[1] = fmt.Sprintf("%v", &z.l)
	if z.r.Signbit() {
		a[2] = fmt.Sprintf("%v", &z.r)
	} else {
		a[2] = fmt.Sprintf("+%v", &z.r)
	}
	a[3] = "α"
	a[4] = ")"
	return strings.Join(a, "")
}

// Equals returns true if y and z are equal.
func (z *Infra) Equals(y *Infra) bool {
	if z.l.Cmp(&y.l) != 0 || z.r.Cmp(&y.r) != 0 {
		return false
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Infra) Copy(y *Infra) *Infra {
	z.l.Copy(&y.l)
	z.r.Copy(&y.r)
	return z
}

// NewInfra returns a pointer to the Infra value a+bα.
func NewInfra(a, b *big.Float) *Infra {
	z := new(Infra)
	z.l.Copy(a)
	z.r.Copy(b)
	return z
}

// Scal sets z equal to y scaled by a, and returns z.
func (z *Infra) Scal(y *Infra, a *big.Float) *Infra {
	z.l.Mul(&y.l, a)
	z.r.Mul(&y.r, a)
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Infra) Neg(y *Infra) *Infra {
	z.l.Neg(&y.l)
	z.r.Neg(&y.r)
	return z
}

// Conj sets z equal to the conjugate of y, and returns z.
func (z *Infra) Conj(y *Infra) *Infra {
	z.l.Copy(&y.l)
	z.r.Neg(&y.r)
	return z
}

// Add sets z equal to the sum of x and y, and returns z.
func (z *Infra) Add(x, y *Infra) *Infra {
	z.l.Add(&x.l, &y.l)
	z.r.Add(&x.r, &y.r)
	return z
}

// Sub sets z equal to the difference of x and y, and returns z.
func (z *Infra) Sub(x, y *Infra) *Infra {
	z.l.Sub(&x.l, &y.l)
	z.r.Sub(&x.r, &y.r)
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
//
// The multiplication rule is:
// 		Mul(α, α) = 0
// This binary operation is commutative and associative.
func (z *Infra) Mul(x, y *Infra) *Infra {
	a := new(big.Float).Copy(&x.l)
	b := new(big.Float).Copy(&x.r)
	c := new(big.Float).Copy(&y.l)
	d := new(big.Float).Copy(&y.r)
	temp := new(big.Float)
	z.l.Mul(a, c)
	z.r.Add(
		z.r.Mul(d, a),
		temp.Mul(b, c),
	)
	return z
}

// Quad returns the quadrance of z, a pointer to a big.Float value.
func (z *Infra) Quad() *big.Float {
	return new(big.Float).Mul(&z.l, &z.l)
}

// IsZeroDiv returns true if z is a zero divisor. This is equivalent to z being
// nilpotent.
func (z *Infra) IsZeroDiv() bool {
	zero := new(big.Float)
	return z.l.Cmp(zero) == 0
}

// Inv sets z equal to the inverse of y, and returns z.
func (z *Infra) Inv(y *Infra) *Infra {
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
func (z *Infra) Quo(x, y *Infra) *Infra {
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

// CrossRatio sets z equal to the cross ratio
// 		Inv(w - x) * (v - x) * Inv(v - y) * (w - y)
// Then it returns z.
func (z *Infra) CrossRatio(v, w, x, y *Infra) *Infra {
	temp := new(Infra)
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
func (z *Infra) Möbius(y, a, b, c, d *Infra) *Infra {
	z.Mul(a, y)
	z.Add(z, b)
	temp := new(Infra)
	temp.Mul(c, y)
	temp.Add(temp, d)
	temp.Inv(temp)
	z.Mul(z, temp)
	return z
}

// Generate returns a random Infra value for quick.Check testing.
func (z *Infra) Generate(rand *rand.Rand, size int) reflect.Value {
	randomInfra := &Infra{
		*big.NewFloat(rand.Float64()),
		*big.NewFloat(rand.Float64()),
	}
	return reflect.ValueOf(randomInfra)
}
