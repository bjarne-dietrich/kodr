package operations

import (
	"crypto/rand"
	"slices"
	"testing"
)

func TestMulConstGeneric(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)
	mulConstRef(refA, b, c)
	mulConstGeneric(a, b, c)
	assertEqual(t, refA, a)

}

func TestMulAddConstGeneric(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)
	mulAddConstRef(refA, b, c)
	mulAddConstGeneric(a, b, c)
	assertEqual(t, refA, a)
}

func TestMulConstTableGeneric(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)
	mulConstRef(refA, b, c)
	table := BuildMulTable(c)
	mulConstTableGeneric(a, b, &table)

	assertEqual(t, refA, a)
}

func TestMulAddConstTableGeneric(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)
	mulAddConstRef(refA, b, c)
	table := BuildMulTable(c)
	mulAddConstTableGeneric(a, b, &table)
	assertEqual(t, refA, a)
}

func TestMulConstNibbleGeneric(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)
	mulConstRef(refA, b, c)
	mulConstNibbleGeneric(a, b, c)
	assertEqual(t, refA, a)
}

func TestMulAddConstNibbleGeneric(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)
	mulAddConstRef(refA, b, c)
	mulAddConstNibbleGeneric(a, b, c)
	assertEqual(t, refA, a)
}

func TestXorAssignSliceGeneric(t *testing.T) {
	a, refA, b, _ := generateRandomData(1023)
	mulAddConstRef(refA, b, 1)
	xorAssignSliceGeneric(a, b)
	assertEqual(t, refA, a)
}

func randBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

func mulConstRef(dst, src []byte, c byte) {
	if len(dst) != len(src) {
		panic("len mismatch")
	}
	for i := range src {
		dst[i] = peasantsAlgorithm(src[i], c)
	}
}

func mulAddConstRef(dst, src []byte, c byte) {
	if len(dst) != len(src) {
		panic("len mismatch")
	}
	for i := range src {
		dst[i] ^= peasantsAlgorithm(src[i], c)
	}
}

func generateRandomData(n int) ([]byte, []byte, []byte, byte) {
	a1 := randBytes(n)
	b := randBytes(n)
	c := randBytes(1)[0]
	a2 := make([]byte, n)
	copy(a2, a1)

	return a1, a2, b, c
}

func assertEqual(t *testing.T, reference, result []byte) {
	if slices.Compare(reference, result) != 0 {
		t.Log("Result: ", result)
		t.Log("Expected: ", reference)
		t.Error("result mismatch")
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
			a ^= ReducePolynomial
		}
	}
	return p
}
