//go:build arm64 && !purego

package operations

import (
	"testing"
)

func TestMulConstNeon(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)

	mulConstRef(refA, b, c)
	mulConstNEON(a, b, c)

	assertEqual(t, refA, a)
}

func TestMulAddConstNeon(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)

	mulAddConstRef(refA, b, c)
	mulAddConstNEON(a, b, c)

	assertEqual(t, refA, a)
}

func TestMulConstNeonUsingMulTable(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)

	mulConstRef(refA, b, c)
	table := BuildMulTable(c)
	mulConstNEONUsingTable(a, b, &table)

	assertEqual(t, refA, a)
}

func TestMulAddConstNeonUsingMulTable(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)

	mulAddConstRef(refA, b, c)
	table := BuildMulTable(c)
	mulAddConstNEONUsingTable(a, b, &table)

	assertEqual(t, refA, a)
}

func TestXorAssignNeon(t *testing.T) {
	a, refA, b, _ := generateRandomData(1023)

	mulAddConstRef(refA, b, 1)
	xorAssignNEON(a, b)

	assertEqual(t, refA, a)
}
