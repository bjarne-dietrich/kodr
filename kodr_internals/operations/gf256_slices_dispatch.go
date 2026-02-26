package operations

// MulConst calculates dst[i] = src[i]*c
func MulConst(dst, src []byte, c byte) {
	mulConstImpl(dst, src, c)
}

// MulAddConst calculates dst[i] += src[i]*c
func MulAddConst(dst, src []byte, c byte) {
	mulAddConstImpl(dst, src, c)
}

// MulConstNibble calculates dst[i] = src[i]*c
// using a nibble table.
func MulConstNibble(dst, src []byte, coef byte) {
	mulAddConstNibbleImpl(dst, src, coef)
}

// MulAddConstNibble calculates dst[i] += src[i]*c
// using a nibble table.
func MulAddConstNibble(dst, src []byte, coef byte) {
	mulAddConstNibbleImpl(dst, src, coef)
}

// MulConstTable calculates dst[i] = src[i]*c
// using a provided multiplication table.
func MulConstTable(dst, src []byte, table *MulTable) {
	mulConstTableImpl(dst, src, table)
}

// MulAddConstTable calculates dst[i] += src[i]*c
// using a provided multiplication table.
func MulAddConstTable(dst, src []byte, table *MulTable) {
	mulAddConstTableImpl(dst, src, table)
}

// XorAssignSlice calculates dst[i] ^= src[i]
func XorAssignSlice(dst, src []byte) {
	xorAssignSliceImpl(dst, src)

}
