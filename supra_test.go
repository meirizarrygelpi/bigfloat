// Copyright (c) 2016 Melvin Eloy Irizarry-Gelpí
// Licenced under the MIT License.

package bigfloat

import (
	"math/big"
	"testing"
	"testing/quick"
)

// Commutativity

func TestSupraAddCommutative(t *testing.T) {
	f := func(x, y *Supra) bool {
		// t.Logf("x = %v, y = %v", x, y)
		l := new(Supra).Add(x, y)
		r := new(Supra).Add(y, x)
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestSupraNegConjCommutative(t *testing.T) {
	f := func(x *Supra) bool {
		// t.Logf("x = %v", x)
		l, r := new(Supra), new(Supra)
		l.Neg(l.Conj(x))
		r.Conj(r.Neg(x))
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Non-commutativity

func TestSupraMulNonCommutative(t *testing.T) {
	f := func(x, y *Supra) bool {
		// t.Logf("x = %v, y = %v", x, y)
		l := new(Supra).Commutator(x, y)
		zero := new(Supra)
		return !l.Equals(zero)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Anti-commutativity

func TestSupraSubAntiCommutative(t *testing.T) {
	f := func(x, y *Supra) bool {
		// t.Logf("x = %v, y = %v", x, y)
		l, r := new(Supra), new(Supra)
		l.Sub(x, y)
		r.Sub(y, x)
		r.Neg(r)
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Associativity

func XTestSupraAddAssociative(t *testing.T) {
	f := func(x, y, z *Supra) bool {
		// t.Logf("x = %v, y = %v, z = %v", x, y, z)
		l, r := new(Supra), new(Supra)
		l.Add(l.Add(x, y), z)
		r.Add(x, r.Add(y, z))
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func XTestSupraMulAssociative(t *testing.T) {
	f := func(x, y, z *Supra) bool {
		// t.Logf("x = %v, y = %v, z = %v", x, y, z)
		l, r := new(Supra), new(Supra)
		l.Mul(l.Mul(x, y), z)
		r.Mul(x, r.Mul(y, z))
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Identity

func TestSupraAddZero(t *testing.T) {
	zero := new(Supra)
	f := func(x *Supra) bool {
		// t.Logf("x = %v", x)
		l := new(Supra).Add(x, zero)
		return l.Equals(x)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestSupraMulOne(t *testing.T) {
	one := &Infra{
		l: *big.NewFloat(1),
	}
	zero := new(Infra)
	f := func(x *Supra) bool {
		// t.Logf("x = %v", x)
		l := new(Supra).Mul(x, &Supra{*one, *zero})
		return l.Equals(x)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func XTestSupraMulInvOne(t *testing.T) {
	one := &Infra{
		l: *big.NewFloat(1),
	}
	zero := new(Infra)
	f := func(x *Supra) bool {
		// t.Logf("x = %v", x)
		l := new(Supra)
		l.Mul(x, l.Inv(x))
		return l.Equals(&Supra{*one, *zero})
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func XTestSupraAddNegSub(t *testing.T) {
	f := func(x, y *Supra) bool {
		// t.Logf("x = %v, y = %v", x, y)
		l, r := new(Supra), new(Supra)
		l.Sub(x, y)
		r.Add(x, r.Neg(y))
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestSupraAddScalDouble(t *testing.T) {
	f := func(x *Supra) bool {
		// t.Logf("x = %v", x)
		l, r := new(Supra), new(Supra)
		l.Add(x, x)
		r.Scal(x, big.NewFloat(2))
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Involutivity

func XTestSupraInvInvolutive(t *testing.T) {
	f := func(x *Supra) bool {
		// t.Logf("x = %v", x)
		l := new(Supra)
		l.Inv(l.Inv(x))
		return l.Equals(x)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestSupraNegInvolutive(t *testing.T) {
	f := func(x *Supra) bool {
		// t.Logf("x = %v", x)
		l := new(Supra)
		l.Neg(l.Neg(x))
		return l.Equals(x)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestSupraConjInvolutive(t *testing.T) {
	f := func(x *Supra) bool {
		// t.Logf("x = %v", x)
		l := new(Supra)
		l.Conj(l.Conj(x))
		return l.Equals(x)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Anti-distributivity

func TestSupraMulConjAntiDistributive(t *testing.T) {
	f := func(x, y *Supra) bool {
		// t.Logf("x = %v, y = %v", x, y)
		l, r := new(Supra), new(Supra)
		l.Conj(l.Mul(x, y))
		r.Mul(r.Conj(y), new(Supra).Conj(x))
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func XTestSupraMulInvAntiDistributive(t *testing.T) {
	f := func(x, y *Supra) bool {
		// t.Logf("x = %v, y = %v", x, y)
		l, r := new(Supra), new(Supra)
		l.Inv(l.Mul(x, y))
		r.Mul(r.Inv(y), new(Supra).Inv(x))
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Distributivity

func TestSupraAddConjDistributive(t *testing.T) {
	f := func(x, y *Supra) bool {
		// t.Logf("x = %v, y = %v", x, y)
		l, r := new(Supra), new(Supra)
		l.Add(x, y)
		l.Conj(l)
		r.Add(r.Conj(x), new(Supra).Conj(y))
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestSupraSubConjDistributive(t *testing.T) {
	f := func(x, y *Supra) bool {
		// t.Logf("x = %v, y = %v", x, y)
		l, r := new(Supra), new(Supra)
		l.Sub(x, y)
		l.Conj(l)
		r.Sub(r.Conj(x), new(Supra).Conj(y))
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestSupraAddScalDistributive(t *testing.T) {
	f := func(x, y *Supra) bool {
		// t.Logf("x = %v, y = %v", x, y)
		a := big.NewFloat(2)
		l, r := new(Supra), new(Supra)
		l.Scal(l.Add(x, y), a)
		r.Add(r.Scal(x, a), new(Supra).Scal(y, a))
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestSupraSubScalDistributive(t *testing.T) {
	f := func(x, y *Supra) bool {
		// t.Logf("x = %v, y = %v", x, y)
		a := big.NewFloat(2)
		l, r := new(Supra), new(Supra)
		l.Scal(l.Sub(x, y), a)
		r.Sub(r.Scal(x, a), new(Supra).Scal(y, a))
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func XTestSupraAddMulDistributive(t *testing.T) {
	f := func(x, y, z *Supra) bool {
		// t.Logf("x = %v, y = %v, z = %v", x, y, z)
		l, r := new(Supra), new(Supra)
		l.Mul(l.Add(x, y), z)
		r.Add(r.Mul(x, z), new(Supra).Mul(y, z))
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func XTestSupraSubMulDistributive(t *testing.T) {
	f := func(x, y, z *Supra) bool {
		// t.Logf("x = %v, y = %v, z = %v", x, y, z)
		l, r := new(Supra), new(Supra)
		l.Mul(l.Sub(x, y), z)
		r.Sub(r.Mul(x, z), new(Supra).Mul(y, z))
		return l.Equals(r)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Positivity

func TestSupraQuadPositive(t *testing.T) {
	f := func(x *Supra) bool {
		// t.Logf("x = %v", x)
		return x.Quad().Sign() > 0
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

// Composition

func XTestSupraComposition(t *testing.T) {
	f := func(x, y *Supra) bool {
		// t.Logf("x = %v, y = %v", x, y)
		p := new(Supra)
		a, b := new(big.Float), new(big.Float)
		p.Mul(x, y)
		a.Set(p.Quad())
		b.Mul(x.Quad(), y.Quad())
		return a.Cmp(b) == 0
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
