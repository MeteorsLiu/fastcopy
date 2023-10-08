#include "textflag.h"

// func copy_movsb(to, from unsafe.Pointer, n uintptr)
TEXT ·copy_movsb(SB),NOSPLIT,$0
    MOVQ to+0(FP), DI
    MOVQ from+8(FP), SI
    MOVQ n+16(FP), CX 
    REP; MOVSB
    RET
