package operations

// mulConstGeneric calculates dst[i] = src[i]*c
// by looping
func mulConstGeneric(dst, src []byte, c byte) {
	if len(src) != len(dst) {
		panic("src and dst length do not match")
	}
	for i := range src {
		MulAssign2(&dst[i], src[i], c)
	}
}

// mulAddConstGeneric calculates dst[i] += src[i]*c
// by looping
func mulAddConstGeneric(dst, src []byte, c byte) {
	if len(src) != len(dst) {
		panic("src and dst length do not match")
	}
	if c == 0 {
		return
	}
	if c == 1 {
		for i := range src {
			AddAssign(&dst[i], src[i])
		}
		return
	}

	for i := range src {
		AddMulAssign(&dst[i], c, src[i])
	}

}

func mulAddConstTableGeneric(dst, src []byte, table *MulTable) {
	for i := range src {
		dst[i] ^= table[src[i]]
	}
}

func mulConstTableGeneric(dst, src []byte, table *MulTable) {
	for i := range src {
		dst[i] = table[src[i]]
	}
}

// mulConstNibbleGeneric calculates dst[i] += src[i]*c
// using a nibble table.
func mulConstNibbleGeneric(dst, src []byte, coef byte) {
	if len(src) != len(dst) {
		panic("src and dst length do not match")
	}
	if len(src) == 0 {
		return
	}

	hi := nibbleTableHi[coef]
	lo := nibbleTableLo[coef]
	for i := 0; i < len(src); i++ {
		dst[i] = hi[(src[i]>>4)] ^ lo[(src[i]&0xf)]
	}

}

// mulAddConstNibbleGeneric calculates dst[i] += src[i]*c
// using a nibble table.
func mulAddConstNibbleGeneric(dst, src []byte, coef byte) {
	if len(src) != len(dst) {
		panic("src and dst length do not match")
	}
	if len(src) == 0 {
		return
	}

	hi := nibbleTableHi[coef]
	lo := nibbleTableLo[coef]
	for i := 0; i < len(src); i++ {
		dst[i] ^= hi[(src[i]>>4)] ^ lo[(src[i]&0xf)]
	}

}

// xorAssignSliceGeneric calculates dst[i] ^= src[i]
func xorAssignSliceGeneric(dst, src []byte) {
	if len(src) != len(dst) {
		panic("src and dst length do not match")
	}
	for i := range src {
		dst[i] ^= src[i]
	}

}
