#include "textflag.h"

// func copy_movsb(dst, src unsafe.Pointer, n uintptr) 
TEXT ·copy_movsb(SB),NOSPLIT,$0
    MOVQ dst+0(FP), DI
    MOVQ src+8(FP), SI
    MOVQ n+16(FP), CX 
    REP; MOVSB
    RET
