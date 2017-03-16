
#include "textflag.h"

// void *memcpy(void *dst, const void *src, size_t n)
// DI = dst, SI = src, DX = size
TEXT clib·_memcpy(SB), NOSPLIT|NOFRAME, $16-0
	PUSHQ R8
	PUSHQ CX
	XORQ CX, CX     // clear register
MEMCPY_QUAD_LOOP:
	ADDQ $8, CX
	CMPQ CX, DX
	JA MEMCPY_QUAD_DONE
	MOVQ -8(SI)(CX*1), R8
	MOVQ R8, -8(DI)(CX*1)
	JMP MEMCPY_QUAD_LOOP
MEMCPY_QUAD_DONE:
	SUBQ $4, CX
	CMPQ CX, DX
	JA MEMCPY_LONG_DONE
	MOVL -4(SI)(CX*1), R8
	MOVL R8, -4(DI)(CX*1)
MEMCPY_LONG_DONE:
	SUBQ $2, CX
	CMPQ CX, DX
	JA MEMCPY_WORD_DONE
	MOVW -2(SI)(CX*1), R8
	MOVW R8, -2(DI)(CX*1)
MEMCPY_WORD_DONE:
	SUBQ $1, CX
	CMPQ CX, DX
	JA MEMCPY_BYTE_DONE
	MOVB -1(SI)(CX*1), R8
	MOVB R8, -1(DI)(CX*1)
MEMCPY_BYTE_DONE:
	POPQ CX
	POPQ R8
	RET

// void *memset(void *str, int c, size_t n)
TEXT clib·_memset(SB), NOSPLIT|NOFRAME, $16-0
    RET

