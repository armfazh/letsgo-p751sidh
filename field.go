package cln16sidh

const Fp751NumWords = 12

type Fp751Element [Fp751NumWords]uint64

/// Compute z = x + y.
func Fp751Add(z, x, y *Fp751Element)

/// Compute z = x - y.
func Fp751Sub(z, x, y *Fp751Element)