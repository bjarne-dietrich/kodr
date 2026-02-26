#include "textflag.h"

// func mulAddConstMulTableNEONASM(dstPtr, srcPtr *byte, n uintptr, tblPtr *byte)
TEXT ·mulAddConstMulTableNEONASM(SB), NOSPLIT, $0-32
    // Load args
    MOVD dstPtr+0(FP), R0      // dst
    MOVD srcPtr+8(FP), R1      // src
    MOVD n+16(FP), R2          // n
    MOVD tblPtr+24(FP), R3      // &tblPtr[0]

    // if n < 32 jump to tail
    CMP  $32, R2
    BLT  tail

    MOVD R3, R7

    // Load 16-byte tables into vector registers once:
    VLD1.P 16(R7), [V16.B16]
    VLD1.P 16(R7), [V17.B16]
    VLD1.P 16(R7), [V18.B16]
    VLD1.P 16(R7), [V19.B16]

    VLD1.P 16(R7), [V20.B16]
    VLD1.P 16(R7), [V21.B16]
    VLD1.P 16(R7), [V22.B16]
    VLD1.P 16(R7), [V23.B16]

    VLD1.P 16(R7), [V24.B16]
    VLD1.P 16(R7), [V25.B16]
    VLD1.P 16(R7), [V26.B16]
    VLD1.P 16(R7), [V27.B16]

    VLD1.P 16(R7), [V28.B16]
    VLD1.P 16(R7), [V29.B16]
    VLD1.P 16(R7), [V30.B16]
    VLD1     (R7), [V31.B16]

    VMOVI $64, V9.B16

    SUB $16, R2, R2

loop16:

    // V0 = src[0:16], src += 16
    VLD1.P 16(R1), [V0.B16]

    // VTBL <indexVec>, [<tableVec>], <dstVec>
    VTBL V0.B16, [V16.B16, V17.B16, V18.B16, V19.B16], V1.B16
    VSUB V9.B16, V0.B16, V0.B16
    VTBL V0.B16, [V20.B16, V21.B16, V22.B16, V23.B16], V2.B16
    VSUB V9.B16, V0.B16, V0.B16
    VTBL V0.B16, [V24.B16, V25.B16, V26.B16, V27.B16], V3.B16
    VSUB V9.B16, V0.B16, V0.B16
    VTBL V0.B16, [V28.B16, V29.B16, V30.B16, V31.B16], V4.B16

    VEOR V1.B16, V2.B16, V1.B16
    VEOR V1.B16, V3.B16, V1.B16
    VEOR V1.B16, V4.B16, V1.B16

    // V5 = dst[0:16]
    VLD1 (R0), [V5.B16]

    // V5 = V1 XOR V5
    VEOR V1.B16, V5.B16, V5.B16

    // dst[0:16] = V5
    VST1.P [V5.B16], 16(R0)

    SUBS  $16, R2, R2
    BGE    loop16

    // R2 is negative remainder, convert to positive tail
    ADD $16, R2, R2

tail:
    CBZ  R2, done

tail_loop:

    MOVBU.P 1(R1), R4     // x = *src
    MOVBU (R3)(R4), R5    // y = tbl[x]
    MOVBU (R0), R6
    EOR   R5, R6, R6
    MOVB.P  R6, 1(R0)
    SUBS  $1, R2, R2
    BNE   tail_loop

done:
    RET
