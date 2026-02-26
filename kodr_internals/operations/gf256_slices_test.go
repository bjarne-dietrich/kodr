package operations

import "testing"

func TestMulConstDispatch(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)
	mulConstRef(refA, b, c)
	MulConst(a, b, c)
	assertEqual(t, refA, a)

}

func TestMulAddConstDispatch(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)
	mulAddConstRef(refA, b, c)
	MulAddConst(a, b, c)
	assertEqual(t, refA, a)
}

func TestMulConstTableDispatch(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)
	mulConstRef(refA, b, c)
	table := BuildMulTable(c)
	MulConstTable(a, b, &table)

	assertEqual(t, refA, a)
}

func TestMulAddConstTableDispatch(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)
	mulAddConstRef(refA, b, c)
	table := BuildMulTable(c)
	MulAddConstTable(a, b, &table)
	assertEqual(t, refA, a)
}

func TestMulConstNibbleDispatch(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)
	mulConstRef(refA, b, c)
	MulConstNibble(a, b, c)
	assertEqual(t, refA, a)
}

func TestMulAddConstNibbleDispatch(t *testing.T) {
	a, refA, b, c := generateRandomData(1023)
	mulAddConstRef(refA, b, c)
	MulAddConstNibble(a, b, c)
	assertEqual(t, refA, a)
}

func TestXorAssignSliceDispatch(t *testing.T) {
	a, refA, b, _ := generateRandomData(1023)
	mulAddConstRef(refA, b, 1)
	XorAssignSlice(a, b)
	assertEqual(t, refA, a)
}
