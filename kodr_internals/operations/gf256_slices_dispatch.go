package operations

// MulConst calculates dst[i] = src[i]*c
// In most cases you'll want to use this function instead of other implementations like MulConstNibble.
func MulConst(dst, src []byte, c byte) {
	mulConstImpl(dst, src, c)
}

// MulAddConst calculates dst[i] += src[i]*c
// In most cases you'll want to use this function instead of other implementations like MulAddConstNibble.
func MulAddConst(dst, src []byte, c byte) {
	mulAddConstImpl(dst, src, c)
}

// MulConstNibble calculates dst[i] = src[i]*c
// using a nibble table.
// In most cases you'll want to use MulConst instead of this implementation.
// The Library then decides the most efficient way to compute.
func MulConstNibble(dst, src []byte, c byte) {
	mulConstNibbleImpl(dst, src, c)
}

// MulAddConstNibble calculates dst[i] += src[i]*c
// using a nibble table.
// In most cases you'll want to use MulAddConst instead of this implementation.
// The Library then decides the most efficient way to compute.
func MulAddConstNibble(dst, src []byte, c byte) {
	mulAddConstNibbleImpl(dst, src, c)
}

// MulConstTable calculates dst[i] = src[i]*c
// using a provided multiplication table.
// In most cases you'll want to use MulConst instead of this implementation.
// The Library then decides the most efficient way to compute.
func MulConstTable(dst, src []byte, table *MulTable) {
	mulConstTableImpl(dst, src, table)
}

// MulAddConstTable calculates dst[i] += src[i]*c
// using a provided multiplication table.
// In most cases you'll want to use MulAddConst instead of this implementation.
// The Library then decides the most efficient way to compute.
func MulAddConstTable(dst, src []byte, table *MulTable) {
	mulAddConstTableImpl(dst, src, table)
}

// XorAssignSlice calculates dst[i] ^= src[i]
func XorAssignSlice(dst, src []byte) {
	xorAssignSliceImpl(dst, src)

}
