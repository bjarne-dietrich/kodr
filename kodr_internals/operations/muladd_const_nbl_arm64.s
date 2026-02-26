#include "textflag.h"

// func mulAddConstNibbleNEONASM(dstPtr, srcPtr *byte, n uintptr, loPtr, hiPtr *byte)
TEXT ·mulAddConstNibbleNEONASM(SB), NOSPLIT, $0-40
    // Load args
    MOVD dstPtr+0(FP), R0      // dst
    MOVD srcPtr+8(FP), R1      // src
    MOVD n+16(FP), R2          // n
    MOVD loPtr+24(FP), R3      // &lo[0]
    MOVD hiPtr+32(FP), R4      // &hi[0]



    // Load 16-byte tables into vector registers once:
    // (loTable in V16, hiTable in V17)
    VLD1 (R3), [V16.B16]
    VLD1 (R4), [V17.B16]

    // Build mask vector 0x0F in V18.
    VMOVI $0x0F, V18.B16

    // if n < 16 jump to tail
    CMP  $16, R2
    BLT  tail

    SUB $16, R2, R2

loop16:

    // V0 = src[0:16], src += 16
    VLD1.P 16(R1), [V0.B16]

    // V1 = vSrc & 0x0F <- lower nibble
    VAND V18.B16, V0.B16, V1.B16
    // V2 = vSrc >> 4
    VUSHR $4, V0.B16, V2.B16

   // VTBL <indexVec>, [<tableVec>], <dstVec>

    // V3 = loTable[vLoIdx]
    VTBL V1.B16, [V16.B16], V3.B16
    // V4 = hiTable[vHiIdx]
    VTBL V2.B16, [V17.B16], V4.B16

    // V3 = vLo XOR vHi
    VEOR V4.B16, V3.B16, V3.B16

    // V5 = dst[0:16]
    VLD1 (R0), [V5.B16]

    // V5 = V3 XOR V5
    VEOR V3.B16, V5.B16, V5.B16

    // dst[0:16] = V5
    VST1.P [V5.B16], 16(R0)

    SUBS  $16, R2, R2
    BGE    loop16

    // R2 is negative remainder, convert to positive tail
    ADD $16, R2, R2

tail:
    CBZ  R2, done

tail_loop:

    MOVBU (R1), R6
    AND   $0x0F, R6, R7
    LSR   $4, R6, R8
    MOVBU (R3)(R7), R9     // lo[x&0xF]
    MOVBU (R4)(R8), R10    // hi[x>>4]
    EOR   R10, R9, R9

    MOVBU (R0), R11
    EOR   R9, R11, R11
    MOVB  R11, (R0)

    ADD  $1, R0, R0
    ADD  $1, R1, R1
    SUBS  $1, R2, R2
    BNE    tail_loop

done:
    RET
