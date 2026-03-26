//go:build arm64 && !purego

package operations

//go:noescape
func mulAddConstNibbleNEONASM(dstPtr, srcPtr *byte, n uintptr, loPtr, hiPtr *byte)

//go:noescape
func mulConstNibbleNEONASM(dstPtr, srcPtr *byte, n uintptr, loPtr, hiPtr *byte)

//go:noescape
func mulAddConstMulTableNEONASM(dstPtr, srcPtr *byte, n uintptr, tblPtr *byte)

//go:noescape
func mulConstMulTableNEONASM(dstPtr, srcPtr *byte, n uintptr, tblPtr *byte)

//go:noescape
func xorAssignNEONASM(slicePtr, otherSlicePtr *byte, n uintptr)

var mulAddConstImpl = mulAddConstNEON
var mulAddConstNibbleImpl = mulAddConstNEON

func mulAddConstNEON(dst, src []byte, c byte) {
	if len(src) != len(dst) {
		panic("src and dst length do not match")
	}
	if len(dst) == 0 {
		return
	}

	mulAddConstNibbleNEONASM(&dst[0], &src[0], uintptr(len(src)),
		&nibbleTableLo[c][0],
		&nibbleTableHi[c][0])

}

var mulConstImpl = mulConstNEON
var mulConstNibbleImpl = mulConstNEON

func mulConstNEON(dst, src []byte, c byte) {
	if len(src) != len(dst) {
		panic("src and dst length do not match")
	}
	if len(dst) == 0 {
		return
	}

	mulConstNibbleNEONASM(&dst[0], &src[0], uintptr(len(src)),
		&nibbleTableLo[c][0],
		&nibbleTableHi[c][0])

}

var mulConstTableImpl = mulConstNEONUsingTable

func mulConstNEONUsingTable(dst, src []byte, t *MulTable) {
	if len(src) != len(dst) {
		panic("src and dst length do not match")
	}
	if len(dst) == 0 {
		return
	}

	mulConstMulTableNEONASM(&dst[0], &src[0], uintptr(len(src)), &t[0])
}

var mulAddConstTableImpl = mulAddConstNEONUsingTable

func mulAddConstNEONUsingTable(dst, src []byte, t *MulTable) {
	if len(src) != len(dst) {
		panic("src and dst length do not match")
	}
	if len(dst) == 0 {
		return
	}

	mulAddConstMulTableNEONASM(&dst[0], &src[0], uintptr(len(src)), &t[0])
}

var xorAssignSliceImpl = xorAssignNEON

func xorAssignNEON(dst, src []byte) {
	if len(src) != len(dst) {
		panic("src and dst length do not match")
	}
	if len(dst) == 0 {
		return
	}
	xorAssignNEONASM(&dst[0], &src[0], uintptr(len(src)))
}
