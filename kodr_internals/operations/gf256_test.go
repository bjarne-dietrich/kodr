package operations_test

import (
	"testing"

	"github.com/itzmeanjan/kodr/kodr_internals/operations"
)

func TestInverse(t *testing.T) {
	for i := 1; i < 256; i++ {
		inv, err := operations.Inverse(byte(i))
		if err != nil {
			t.Error(err)
		}
		if operations.Mul(inv, byte(i)) != operations.MultiplicativeIdentityElement {
			t.Error("Inverse result mismatch")
		}
	}
}

func TestAssignInverse(t *testing.T) {
	for i := 1; i < 256; i++ {
		bi := byte(i)
		err := operations.InverseAssign(&bi)
		if err != nil {
			t.Error(err)
		}
		if operations.Mul(bi, byte(i)) != operations.MultiplicativeIdentityElement {
			t.Error("Inverse result mismatch")
		}
	}
}

func TestAdd(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			if operations.Add(byte(a), byte(b)) != byte(a)^byte(b) {
				t.Error("Add result mismatch")
			}
		}
	}
}

func TestAddAssign(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			ba := byte(a)
			operations.AddAssign(&ba, byte(b))
			if ba != byte(a)^byte(b) {
				t.Error("Add result mismatch")
			}
		}
	}
}

func TestSub(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			if operations.Sub(byte(a), byte(b)) != byte(a)^byte(b) {
				t.Error("Sub result mismatch")
			}
		}
	}
}

func TestSubAssign(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			ba := byte(a)
			operations.SubAssign(&ba, byte(b))
			if ba != byte(a)^byte(b) {
				t.Error("Sub result mismatch")
			}
		}
	}
}

func TestMul(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			if operations.Mul(byte(a), byte(b)) != peasantsAlgorithm(byte(a), byte(b)) {
				t.Error("Mul result mismatch")
			}
		}
	}
}

func TestMulAssign(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			ba := byte(a)
			operations.MulAssign(&ba, byte(b))
			if ba != peasantsAlgorithm(byte(a), byte(b)) {
				t.Error("Mul result mismatch")
			}
		}
	}
}

func TestMulAssign2(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			res := byte(0)
			operations.MulAssign2(&res, byte(a), byte(b))
			if res != peasantsAlgorithm(byte(a), byte(b)) {
				t.Error("Mul result mismatch")
			}
		}
	}
}

func TestDiv(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			if invb, err := operations.Inverse(byte(b)); err != nil {
				if byte(b) == operations.AdditiveIdentityElement {
					continue
				}
			} else if div, err2 := operations.Div(byte(a), byte(b)); err2 != nil {
				t.Error(err2)
			} else if div != peasantsAlgorithm(invb, byte(a)) {
				t.Error("Div result mismatch")
			}

		}
	}
}

func TestDivAssign(t *testing.T) {
	for a := 0; a < 256; a++ {
		for b := 0; b < 256; b++ {
			div := byte(a)
			if invb, err := operations.Inverse(byte(b)); err != nil {
				if byte(b) == operations.AdditiveIdentityElement {
					continue
				}
			} else if err2 := operations.DivAssign(&div, byte(b)); err2 != nil {
				t.Error(err2)
			} else if div != peasantsAlgorithm(invb, byte(a)) {
				t.Error("Div result mismatch")
			}

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

			if pow != refPow {
				t.Log("a", a, "b", b, "pow", pow)
				t.Log("refPow", refPow)
				t.Error("Pow result mismatch")
			}
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

			if pow != refPow {
				t.Log("a", a, "b", b, "pow", pow)
				t.Log("refPow", refPow)
				t.Error("Pow result mismatch")
			}
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
