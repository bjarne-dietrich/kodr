#include "textflag.h"

// func xorAssignNEONASM(slicePtr, otherSlicePtr *byte, n uintptr)
TEXT ·xorAssignNEONASM(SB), NOSPLIT, $0-24
    // Load args
    MOVD slicePtr+0(FP), R0         // slice
    MOVD otherSlicePtr+8(FP), R1    // other
    MOVD n+16(FP), R2               // n

    // if n < 16 jump to tail
    CMP  $16, R2
    BLT  tail

    SUB $16, R2, R2

loop16:

    // Fetch Data
    VLD1 (R0), [V0.B16] // V0 = slice[i:16+i]
    VLD1.P 16(R1), [V1.B16] // V1 = other[i:16+i]
    VEOR V0.B16, V1.B16, V0.B16 // V0 = V0 ^ V1
    VST1.P [V0.B16], 16(R0) // slice[i:16+i] = V0

    SUBS  $16, R2, R2
    BGE    loop16

    // R2 is negative remainder, convert to positive tail
    ADD $16, R2, R2

tail:
    CBZ  R2, done

tail_loop:

    MOVBU (R0), R6
    MOVBU.P 1(R1), R7
    EOR R7, R6, R6
    MOVB.P  R6, 1(R0)

    SUBS  $1, R2, R2
    BNE    tail_loop

done:
    RET
