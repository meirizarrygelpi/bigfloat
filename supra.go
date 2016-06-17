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

var symbSupra = [4]string{"", "α", "β", "γ"}

// A Supra represents a multi-precision floating-point supra number.
type Supra struct {
	l, r Infra
}

// Real returns the real part of z.
func (z *Supra) Real() *big.Float {
	return (&z.l).Real()
}

// Cartesian returns the four multi-precision floating-point Cartesian
// components of z.
func (z *Supra) Cartesian() (*big.Float, *big.Float, *big.Float, *big.Float) {
	return &z.l.l, &z.l.r, &z.r.l, &z.r.r
}

// String returns the string representation of a Supra value.
//
// If z corresponds to a + bα + cβ + dγ, then the string is "(a+bα+cβ+dγ)",
// similar to complex128 values.
func (z *Supra) String() string {
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
		a[j+1] = symbSupra[i]
		i++
	}
	a[8] = ")"
	return strings.Join(a, "")
}

// Equals returns true if y and z are equal.
func (z *Supra) Equals(y *Supra) bool {
	if !z.l.Equals(&y.l) || !z.r.Equals(&y.r) {
		return false
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Supra) Copy(y *Supra) *Supra {
	z.l.Copy(&y.l)
	z.r.Copy(&y.r)
	return z
}

// NewSupra returns a pointer to the Supra value a+bα+cβ+dγ.
func NewSupra(a, b, c, d *big.Float) *Supra {
	z := new(Supra)
	z.l.l.Copy(a)
	z.l.r.Copy(b)
	z.r.l.Copy(c)
	z.r.r.Copy(d)
	return z
}

// Scal sets z equal to y scaled by a, and returns z.
func (z *Supra) Scal(y *Supra, a *big.Float) *Supra {
	z.l.Scal(&y.l, a)
	z.r.Scal(&y.r, a)
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Supra) Neg(y *Supra) *Supra {
	z.l.Neg(&y.l)
	z.r.Neg(&y.r)
	return z
}

// Conj sets z equal to the conjugate of y, and returns z.
func (z *Supra) Conj(y *Supra) *Supra {
	z.l.Conj(&y.l)
	z.r.Neg(&y.r)
	return z
}

// Add sets z equal to x+y, and returns z.
func (z *Supra) Add(x, y *Supra) *Supra {
	z.l.Add(&x.l, &y.l)
	z.r.Add(&x.r, &y.r)
	return z
}

// Sub sets z equal to x-y, and returns z.
func (z *Supra) Sub(x, y *Supra) *Supra {
	z.l.Sub(&x.l, &y.l)
	z.r.Sub(&x.r, &y.r)
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
//
// The multiplication rules are:
// 		Mul(α, α) = Mul(β, β) = Mul(γ, γ) = 0
// 		Mul(α, β) = -Mul(β, α) = γ
// 		Mul(β, γ) = Mul(γ, β) = 0
// 		Mul(γ, α) = Mul(α, γ) = 0
// This binary operation is noncommutative but associative.
func (z *Supra) Mul(x, y *Supra) *Supra {
	a := new(Infra).Copy(&x.l)
	b := new(Infra).Copy(&x.r)
	c := new(Infra).Copy(&y.l)
	d := new(Infra).Copy(&y.r)
	temp := new(Infra)
	z.l.Mul(a, c)
	z.r.Add(
		z.r.Mul(d, a),
		temp.Mul(b, temp.Conj(c)),
	)
	return z
}

// Commutator sets z equal to the commutator of x and y:
// 		Mul(x, y) - Mul(y, x)
// Then it returns z.
func (z *Supra) Commutator(x, y *Supra) *Supra {
	return z.Sub(
		z.Mul(x, y),
		new(Supra).Mul(y, x),
	)
}

// Quad returns the quadrance of z. If z = a+bα+cβ+dγ, then the quadrance is
// 		Mul(a, a)
// This is always non-negative.
func (z *Supra) Quad() *big.Float {
	return z.l.Quad()
}

// IsZeroDiv returns true if z is a zero divisor.
func (z *Supra) IsZeroDiv() bool {
	return z.l.IsZeroDiv()
}

// Inv sets z equal to the inverse of y, and returns z. If y is a zero divisor,
// then Inv panics.
func (z *Supra) Inv(y *Supra) *Supra {
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
func (z *Supra) QuoL(x, y *Supra) *Supra {
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
func (z *Supra) QuoR(x, y *Supra) *Supra {
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

// CrossFloatioL sets z equal to the left cross-ratio of v, w, x, and y:
// 		Inv(w - x) * (v - x) * Inv(v - y) * (w - y)
// Then it returns z.
func (z *Supra) CrossFloatioL(v, w, x, y *Supra) *Supra {
	temp := new(Supra)
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

// CrossFloatioR sets z equal to the right cross-ratio of v, w, x, and y:
// 		(v - x) * Inv(w - x) * (w - y) * Inv(v - y)
// Then it returns z.
func (z *Supra) CrossFloatioR(v, w, x, y *Supra) *Supra {
	temp := new(Supra)
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
func (z *Supra) MöbiusL(y, a, b, c, d *Supra) *Supra {
	z.Mul(y, a)
	z.Add(z, b)
	temp := new(Supra)
	temp.Mul(y, c)
	temp.Add(temp, d)
	temp.Inv(temp)
	return z.Mul(temp, z)
}

// MöbiusR sets z equal to the right Möbius (fractional linear) transform of y:
// 		(a*y + b) * Inv(c*y + d)
// Then it returns z.
func (z *Supra) MöbiusR(y, a, b, c, d *Supra) *Supra {
	z.Mul(a, y)
	z.Add(z, b)
	temp := new(Supra)
	temp.Mul(c, y)
	temp.Add(temp, d)
	temp.Inv(temp)
	return z.Mul(z, temp)
}

// Generate returns a random Supra value for quick.Check testing.
func (z *Supra) Generate(rand *rand.Rand, size int) reflect.Value {
	randomSupra := &Supra{
		*NewInfra(
			big.NewFloat(rand.Float64()),
			big.NewFloat(rand.Float64()),
		),
		*NewInfra(
			big.NewFloat(rand.Float64()),
			big.NewFloat(rand.Float64()),
		),
	}
	return reflect.ValueOf(randomSupra)
}
