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

var symbHamilton = [4]string{"", "i", "j", "k"}

// A Hamilton represents a multi-precision floating-point Hamilton quaternion.
type Hamilton struct {
	l, r Complex
}

// Real returns the real part of z.
func (z *Hamilton) Real() *big.Float {
	return (&z.l).Real()
}

// Cartesian returns the four multi-precision floating-point Cartesian
// components of z.
func (z *Hamilton) Cartesian() (*big.Float, *big.Float, *big.Float, *big.Float) {
	return &z.l.l, &z.l.r, &z.r.l, &z.r.r
}

// String returns the string representation of a Hamilton value.
//
// If z corresponds to a + bi + cj + dk, then the string is"(a+bi+cj+dk)",
// similar to complex128 values.
func (z *Hamilton) String() string {
	v := make([]*big.Float, 4)
	v[0], v[1] = z.l.Cartesian()
	v[2], v[3] = z.r.Cartesian()
	a := make([]string, 9)
	a[0] = "("
	a[1] = fmt.Sprintf("%v", v[0])
	i := 1
	for j := 2; j < 8; j = j + 2 {
		if v[i].Sign() < 0 {
			a[j] = fmt.Sprintf("%v", v[i])
		} else {
			a[j] = fmt.Sprintf("+%v", v[i])
		}
		a[j+1] = symbHamilton[i]
		i++
	}
	a[8] = ")"
	return strings.Join(a, "")
}

// Equals returns true if y and z are equal.
func (z *Hamilton) Equals(y *Hamilton) bool {
	if !z.l.Equals(&y.l) || !z.r.Equals(&y.r) {
		return false
	}
	return true
}

// Copy copies y onto z, and returns z.
func (z *Hamilton) Copy(y *Hamilton) *Hamilton {
	z.l.Copy(&y.l)
	z.r.Copy(&y.r)
	return z
}

// NewHamilton returns a pointer to the Hamilton value a+bi+cj+dk.
func NewHamilton(a, b, c, d *big.Float) *Hamilton {
	z := new(Hamilton)
	z.l.l.Copy(a)
	z.l.r.Copy(b)
	z.r.l.Copy(c)
	z.r.r.Copy(d)
	return z
}

// Scal sets z equal to y scaled by a, and returns z.
func (z *Hamilton) Scal(y *Hamilton, a *big.Float) *Hamilton {
	z.l.Scal(&y.l, a)
	z.r.Scal(&y.r, a)
	return z
}

// Neg sets z equal to the negative of y, and returns z.
func (z *Hamilton) Neg(y *Hamilton) *Hamilton {
	z.l.Neg(&y.l)
	z.r.Neg(&y.r)
	return z
}

// Conj sets z equal to the conjugate of y, and returns z.
func (z *Hamilton) Conj(y *Hamilton) *Hamilton {
	z.l.Conj(&y.l)
	z.r.Neg(&y.r)
	return z
}

// Add sets z equal to x+y, and returns z.
func (z *Hamilton) Add(x, y *Hamilton) *Hamilton {
	z.l.Add(&x.l, &y.l)
	z.r.Add(&x.r, &y.r)
	return z
}

// Sub sets z equal to x-y, and returns z.
func (z *Hamilton) Sub(x, y *Hamilton) *Hamilton {
	z.l.Sub(&x.l, &y.l)
	z.r.Sub(&x.r, &y.r)
	return z
}

// Mul sets z equal to the product of x and y, and returns z.
//
// The multiplication rules are:
// 		Mul(i, i) = Mul(j, j) = Mul(k, k) = -1
// 		Mul(i, j) = -Mul(j, i) = k
// 		Mul(j, k) = -Mul(k, j) = i
// 		Mul(k, i) = -Mul(i, k) = j
// This binary operation is noncommutative but associative.
func (z *Hamilton) Mul(x, y *Hamilton) *Hamilton {
	a := new(Complex).Copy(&x.l)
	b := new(Complex).Copy(&x.r)
	c := new(Complex).Copy(&y.l)
	d := new(Complex).Copy(&y.r)
	temp := new(Complex)
	z.l.Sub(
		z.l.Mul(a, c),
		temp.Mul(temp.Conj(d), b),
	)
	z.r.Add(
		z.r.Mul(d, a),
		temp.Mul(b, temp.Conj(c)),
	)
	return z
}

// Commutator sets z equal to the commutator of x and y:
// 		Mul(x, y) - Mul(y, x)
// Then it returns z.
func (z *Hamilton) Commutator(x, y *Hamilton) *Hamilton {
	return z.Sub(
		z.Mul(x, y),
		new(Hamilton).Mul(y, x),
	)
}

// Quad returns the quadrance of z. If z = a+bi+cj+dk, then the quadrance is
// 		Mul(a, a) + Mul(b, b) + Mul(c, c) + Mul(d, d)
// This is always non-negative.
func (z *Hamilton) Quad() *big.Float {
	return new(big.Float).Add(
		z.l.Quad(),
		z.r.Quad(),
	)
}

// Inv sets z equal to the inverse of y, and returns z. If y is zero, then Inv
// panics.
func (z *Hamilton) Inv(y *Hamilton) *Hamilton {
	if zero := new(Hamilton); y.Equals(zero) {
		panic("inverse of zero")
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
// Then it returns z. If y is zero, then QuoL panics.
func (z *Hamilton) QuoL(x, y *Hamilton) *Hamilton {
	if zero := new(Hamilton); y.Equals(zero) {
		panic("left denominator is zero")
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
// Then it returns z. If y is zero, then QuoR panics.
func (z *Hamilton) QuoR(x, y *Hamilton) *Hamilton {
	if zero := new(Hamilton); y.Equals(zero) {
		panic("right denominator is zero")
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
func (z *Hamilton) CrossRatioL(v, w, x, y *Hamilton) *Hamilton {
	temp := new(Hamilton)
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
func (z *Hamilton) CrossRatioR(v, w, x, y *Hamilton) *Hamilton {
	temp := new(Hamilton)
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
func (z *Hamilton) MöbiusL(y, a, b, c, d *Hamilton) *Hamilton {
	z.Mul(y, a)
	z.Add(z, b)
	temp := new(Hamilton)
	temp.Mul(y, c)
	temp.Add(temp, d)
	temp.Inv(temp)
	return z.Mul(temp, z)
}

// MöbiusR sets z equal to the right Möbius (fractional linear) transform of y:
// 		(a*y + b) * Inv(c*y + d)
// Then it returns z.
func (z *Hamilton) MöbiusR(y, a, b, c, d *Hamilton) *Hamilton {
	z.Mul(a, y)
	z.Add(z, b)
	temp := new(Hamilton)
	temp.Mul(c, y)
	temp.Add(temp, d)
	temp.Inv(temp)
	return z.Mul(z, temp)
}

// Generate returns a random Hamilton value for quick.Check testing.
func (z *Hamilton) Generate(rand *rand.Rand, size int) reflect.Value {
	randomHamilton := &Hamilton{
		*NewComplex(
			big.NewFloat(rand.Float64()),
			big.NewFloat(rand.Float64()),
		),
		*NewComplex(
			big.NewFloat(rand.Float64()),
			big.NewFloat(rand.Float64()),
		),
	}
	return reflect.ValueOf(randomHamilton)
}
