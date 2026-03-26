package operations_test

import (
	"testing"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals/operations"
	"github.com/stretchr/testify/assert"
)

func TestInverse(t *testing.T) {
	for i := 1; i < 256; i++ {
		inv, err := operations.Inverse(byte(i))
		assert.Nil(t, err)
		assert.Equal(t, operations.MultiplicativeIdentityElement, operations.Mul(inv, byte(i)), "Inverse Check failed")
	}
}

func TestAssignInverse(t *testing.T) {
	for i := 1; i < 256; i++ {
		bi := byte(i)
		err := operations.InverseAssign(&bi)
		assert.Nil(t, err)
		assert.Equal(t, operations.MultiplicativeIdentityElement, operations.Mul(bi, byte(i)), "Inverse Check failed")
	}
}

func TestAdd(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			assert.Equal(t, byte(a)^byte(b), operations.Add(byte(a), byte(b)))
		}
	}
}

func TestAddAssign(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			ba := byte(a)
			operations.AddAssign(&ba, byte(b))
			assert.Equal(t, byte(a)^byte(b), ba)
		}
	}
}

func TestSub(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			assert.Equal(t, byte(a)^byte(b), operations.Sub(byte(a), byte(b)))
		}
	}
}

func TestSubAssign(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			ba := byte(a)
			operations.SubAssign(&ba, byte(b))
			assert.Equal(t, byte(a)^byte(b), ba)
		}
	}
}

func TestMul(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			assert.Equal(t, peasantsAlgorithm(byte(a), byte(b)), operations.Mul(byte(a), byte(b)))
		}
	}
}

func TestMulAssign(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			ba := byte(a)
			operations.MulAssign(&ba, byte(b))
			assert.Equal(t, peasantsAlgorithm(byte(a), byte(b)), ba)
		}
	}
}

func TestMulAssign2(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			res := byte(0)
			operations.MulAssign2(&res, byte(a), byte(b))
			assert.Equal(t, peasantsAlgorithm(byte(a), byte(b)), res)
		}
	}
}

func TestDiv(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			inverse, err := operations.Inverse(byte(b))

			if byte(b) == operations.AdditiveIdentityElement {
				assert.ErrorIs(t, err, kodr.ErrCannotInvertGf256AdditiveIdentity)
				continue
			}

			assert.Nil(t, err)

			div, err := operations.Div(byte(a), byte(b))
			assert.Nil(t, err)
			assert.Equal(t, peasantsAlgorithm(inverse, byte(a)), div)
		}
	}
}

func TestDivAssign(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			div := byte(a)
			inverse, err := operations.Inverse(byte(b))

			if byte(b) == operations.AdditiveIdentityElement {
				assert.ErrorIs(t, err, kodr.ErrCannotInvertGf256AdditiveIdentity)
				continue
			}

			assert.Nil(t, err)

			err = operations.DivAssign(&div, byte(b))
			assert.Nil(t, err)
			assert.Equal(t, peasantsAlgorithm(inverse, byte(a)), div)

		}
	}
}

func TestPow(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			pow := operations.Pow(byte(a), b)

			refPow := byte(1)
			for range b {
				operations.MulAssign(&refPow, byte(a))
			}

			assert.Equal(t, refPow, pow)
		}
	}
}

func TestPowAssign(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			pow := byte(a)
			operations.PowAssign(&pow, b)

			refPow := byte(1)
			for range b {
				operations.MulAssign(&refPow, byte(a))
			}

			assert.Equal(t, refPow, pow)
		}
	}
}

func peasantsAlgorithm(a, b byte) byte {
	var p byte
	for i := 0; i < 8; i++ {
		if (b & 1) != 0 {
			p ^= a
		}
		b >>= 1
		carry := a & 0x80
		a <<= 1
		if carry != 0 {
			a ^= operations.ReducePolynomial
		}
	}
	return p
}
