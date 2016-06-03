# bigfloat

Package `bigfloat` brings multi-precision floating-point [complex](https://en.wikipedia.org/wiki/Complex_number), [split-complex](https://en.wikipedia.org/wiki/Split-complex_number), and [dual](https://en.wikipedia.org/wiki/Dual_number) numbers to Go. It borrows heavily from the `math`, `math/cmplx`, and `math/big` packages. Indeed, it is built on top of the `big.Float` type from the `math/big` package.

[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/github.com/meirizarrygelpi/bigfloat) [![GoDoc](https://godoc.org/github.com/meirizarrygelpi/bigfloat?status.svg)](https://godoc.org/github.com/meirizarrygelpi/bigfloat)

Note that some of the tests fail because floating-point arithmetic is different from real number arithmetic.

## To Do

1. Improve documentation
1. Tests
1. Improve README