package operations

import "github.com/itzmeanjan/kodr"

const PrimitiveElement byte = 2
const MultiplicativeIdentityElement byte = 1
const One byte = 1
const Zero byte = 0
const AdditiveIdentityElement byte = 0
const ReducePolynomial byte = 0x1D

func Inverse(of byte) (byte, error) {
	if of == AdditiveIdentityElement {
		return 0, kodr.ErrCannotInvertGf256AdditiveIndentity
	}

	return gf256InverseTable[of], nil
}

func InverseAssign(a *byte) error {
	if *a == AdditiveIdentityElement {
		return kodr.ErrCannotInvertGf256AdditiveIndentity
	}

	*a = gf256InverseTable[*a]
	return nil
}

// Add Performs addition (XOR) of two GF256 Elements
func Add(a, b byte) byte {
	return a ^ b
}

// AddAssign Performs addition (XOR) of two GF256 Elements
// in place of the first element
func AddAssign(a *byte, b byte) {
	*a ^= b
}

// Sub Performs substraction (XOR) of two GF256 Elements
func Sub(a, b byte) byte {
	return a ^ b
}

// SubAssign Performs substraction (XOR) of two GF256 Elements
// in place of the first element
func SubAssign(a *byte, b byte) {
	*a ^= b
}

// Mul Performs multiplication of two GF256 Elements
func Mul(a, b byte) byte {
	if a == AdditiveIdentityElement || b == AdditiveIdentityElement {
		return 0
	}

	la := int(gf256LogTable[a])
	lb := int(gf256LogTable[b])

	return gf256ExpTable[la+lb]

}

// MulAssign Performs multiplication of two GF256 Elements
// in place of the first element a=a*b
func MulAssign(a *byte, b byte) {
	if *a == AdditiveIdentityElement {
		return
	}
	if b == AdditiveIdentityElement {
		*a = AdditiveIdentityElement
		return
	}

	la := int(gf256LogTable[*a])
	lb := int(gf256LogTable[b])

	*a = gf256ExpTable[la+lb]
}

// MulAssign2 Performs multiplication of two GF256 Elements b*c
// and stores it in a: a=b*c
func MulAssign2(a *byte, b, c byte) {
	if b == AdditiveIdentityElement {
		*a = AdditiveIdentityElement
		return
	}
	if c == AdditiveIdentityElement {
		*a = AdditiveIdentityElement
		return
	}

	lb := int(gf256LogTable[b])
	lc := int(gf256LogTable[c])

	*a = gf256ExpTable[lb+lc]
}

// Div performs division of two Gf256
// elements using multiplicative inverse
func Div(dividend, divisor byte) (byte, error) {
	if err := InverseAssign(&divisor); err != nil {
		return 0, err
	} else {
		return Mul(dividend, divisor), nil
	}
}

// DivAssign Performs division of two GF256 Elements
// in place of the first element
func DivAssign(dividend *byte, divisor byte) error {
	if err := InverseAssign(&divisor); err != nil {
		return err
	} else {
		MulAssign(dividend, divisor)
		return nil
	}
}

func Pow(a byte, b int) byte {
	if b == 0 {
		return MultiplicativeIdentityElement
	}
	if a == AdditiveIdentityElement {
		return AdditiveIdentityElement
	}

	la := int(gf256LogTable[a])
	return gf256ExpTable[(la*b)%(gf256Order-1)]
}

func PowAssign(a *byte, b int) {
	if b == 0 {
		*a = MultiplicativeIdentityElement
		return
	}
	if *a == AdditiveIdentityElement {
		return
	}

	la := int(gf256LogTable[*a])
	*a = gf256ExpTable[(la*b)%(gf256Order-1)]
}

// AddMulAssign adds b*c to a
func AddMulAssign(a *byte, b, c byte) {
	MulAssign(&b, c)
	AddAssign(a, b)
}
