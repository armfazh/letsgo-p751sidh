package cln16sidh

//------------------------------------------------------------------------------
// Extension Field
//------------------------------------------------------------------------------

// Represents an element of the extension field F_{p^2}.
type ExtensionFieldElement struct {
	// This field element is in Montgomery form, so that the value `a` is
	// represented by `aR mod p`.
	a fp751Element
	// This field element is in Montgomery form, so that the value `b` is
	// represented by `bR mod p`.
	b fp751Element
}

// Set dest = lhs * rhs.
//
// Allowed to overlap lhs or rhs with dest.
func (dest *ExtensionFieldElement) Mul(lhs, rhs *ExtensionFieldElement) {
	// Let (a,b,c,d) = (lhs.A,lhs.B,rhs.A,rhs.B).

	a := &lhs.a
	b := &lhs.b
	c := &rhs.a
	d := &rhs.b

	// We want to compute
	//
	// (a + bi)*(c + di) = (a*c - b*d) + (a*d + b*c)i
	//
	// Use Karatsuba's trick: note that
	//
	// (b - a)*(c - d) = (b*c + a*d) - a*c - b*d
	//
	// so (a*d + b*c) = (b-a)*(c-d) + a*c + b*d.

	var ac, bd fp751X2
	fp751Mul(&ac, a, c)				// = a*c*R*R
	fp751Mul(&bd, b, d)				// = b*d*R*R

	var b_minus_a, c_minus_d fp751Element
	fp751SubReduced(&b_minus_a, b, a)		// = (b-a)*R
	fp751SubReduced(&c_minus_d, c, d)		// = (c-d)*R

	var ad_plus_bc fp751X2
	fp751Mul(&ad_plus_bc, &b_minus_a, &c_minus_d)	// = (b-a)*(c-d)*R*R
	fp751X2AddLazy(&ad_plus_bc, &ad_plus_bc, &ac)	// = ((b-a)*(c-d) - a*c)*R*R
	fp751X2AddLazy(&ad_plus_bc, &ad_plus_bc, &bd)	// = ((b-a)*(c-d) - a*c - b*d)*R*R

	fp751MontgomeryReduce(&dest.a, &ad_plus_bc)	// = (a*d + b*c)*R mod p

	fp751X2AddLazy(&ac, &ac, &bd)			// = (a*c + b*d)*R*R
	fp751MontgomeryReduce(&dest.b, &ac)		// = (a*c + b*d)*R mod p
}

// Set dest = lhs + rhs.
//
// Allowed to overlap lhs or rhs with dest.
func (dest *ExtensionFieldElement) Add(lhs, rhs *ExtensionFieldElement) {
	fp751AddReduced(&dest.a, &lhs.a, &rhs.a)
	fp751AddReduced(&dest.b, &lhs.b, &rhs.b)
}

// Set dest = lhs - rhs.
//
// Allowed to overlap lhs or rhs with dest.
func (dest *ExtensionFieldElement) Sub(lhs, rhs *ExtensionFieldElement) {
	fp751SubReduced(&dest.a, &lhs.a, &rhs.a)
	fp751SubReduced(&dest.b, &lhs.b, &rhs.b)
}

// Returns true if lhs = rhs.  Takes variable time.
func (lhs *ExtensionFieldElement) VartimeEq(rhs *ExtensionFieldElement) bool {
	return lhs.a.vartimeEq(rhs.a) && lhs.b.vartimeEq(rhs.b)
}

//------------------------------------------------------------------------------
// Prime Field
//------------------------------------------------------------------------------

// Represents an element of the prime field F_p.
type PrimeFieldElement struct {
	// This field element is in Montgomery form, so that the value `a` is
	// represented by `aR mod p`.
	a fp751Element
}

// Set dest = lhs * rhs.
//
// Allowed to overlap lhs or rhs with dest.
func (dest *PrimeFieldElement) Mul(lhs, rhs *PrimeFieldElement) {
	a := &lhs.a				// = a*R
	b := &rhs.a				// = b*R

	var ab fp751X2
	fp751Mul(&ab, a, b)			// = a*b*R*R
	fp751MontgomeryReduce(&dest.a, &ab)	// = a*b*R mod p
}

// Set dest = lhs + rhs.
//
// Allowed to overlap lhs or rhs with dest.
func (dest *PrimeFieldElement) Add(lhs, rhs *PrimeFieldElement) {
	fp751AddReduced(&dest.a, &lhs.a, &rhs.a)
}

// Set dest = lhs - rhs.
//
// Allowed to overlap lhs or rhs with dest.
func (dest *PrimeFieldElement) Sub(lhs, rhs *PrimeFieldElement) {
	fp751SubReduced(&dest.a, &lhs.a, &rhs.a)
}

// Returns true if lhs = rhs.  Takes variable time.
func (lhs *PrimeFieldElement) VartimeEq(rhs *PrimeFieldElement) bool {
	return lhs.a.vartimeEq(rhs.a)
}

//------------------------------------------------------------------------------
// Internals
//------------------------------------------------------------------------------

const fp751NumWords = 12

// Internal representation of an element of the base field F_p.
//
// This type is distinct from PrimeFieldElement in that no particular meaning
// is assigned to the representation -- it could represent an element in
// Montgomery form, or not.  Tracking the meaning of the field element is left
// to higher types.
type fp751Element [fp751NumWords]uint64

// Represents an intermediate product of two elements of the base field F_p.
type fp751X2 [2 * fp751NumWords]uint64

// Compute z = x + y (mod p).
//go:noescape
func fp751AddReduced(z, x, y *fp751Element)

// Compute z = x - y (mod p).
//go:noescape
func fp751SubReduced(z, x, y *fp751Element)

// Compute z = x + y, without reducing mod p.
//go:noescape
func fp751AddLazy(z, x, y *fp751Element)

// Compute z = x + y, without reducing mod p.
//go:noescape
func fp751X2AddLazy(z, x, y *fp751X2)

// Compute z = x * y.
//go:noescape
func fp751Mul(z *fp751X2, x, y *fp751Element)

// Perform Montgomery reduction: set z = x R^{-1} (mod p).
// Destroys the input value.
//go:noescape
func fp751MontgomeryReduce(z *fp751Element, x *fp751X2)

// Reduce a field element in [0, 2*p) to one in [0,p).
//go:noescape
func fp751StrongReduce(x *fp751Element)

func (x fp751Element) vartimeEq(y fp751Element) bool {
	fp751StrongReduce(&x)
	fp751StrongReduce(&y)
	eq := true
	for i := 0; i < fp751NumWords; i++ {
		eq = (x[i] == y[i]) && eq
	}

	return eq
}
