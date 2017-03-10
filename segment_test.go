package main

import (
	"strings"
	"testing"
)

func testSegment(t *testing.T, fullsrc []string, expected []Segment) {
	segments := SegmentSource(fullsrc)
	if !segmentEqual(segments, expected) {
		t.Errorf("TestNames(): \nexpected %s\ngot      %s", expected, segments)
	}
}

func TestSegment(t *testing.T) {

	src1 := `	.section	__TEXT,__text,regular,pure_instructions
	.macosx_version_min 10, 11
	.intel_syntax noprefix
	.section	__TEXT,__const
	.align	5
LCPI0_0:
	.byte	255                     ## 0xff
	.byte	0                       ## 0x0
LCPI0_1:
	.short	9617                    ## 0x2591
	.short	0                       ## 0x0
LCPI0_2:
	.short	1868                    ## 0x74c
	.short	4899                    ## 0x1323
	.section	__TEXT,__literal4,4byte_literals
	.align	2
LCPI0_3:
	.long	8192                    ## 0x2000
	.section	__TEXT,__text,regular,pure_instructions
	.globl	__ZN4Simd4Avx210BgraToGrayEPKhmmmPhm
	.align	4, 0x90
__ZN4Simd4Avx210BgraToGrayEPKhmmmPhm:   ## @_ZN4Simd4Avx210BgraToGrayEPKhmmmPhm
## BB#0:
	push    rbp
	mov     rbp, rsp
	mov     rax, rdi
	and     rax, -32
	cmp     rax, rdi
	jne     LBB0_9
## BB#1:
	mov	r10, r9
	jne	LBB0_9
## BB#2:
	mov	rax, r8
	jne	LBB0_9
## BB#3:
	test	rdx, rdx
	je	LBB0_15
## BB#4:                                ## %.preheader.lr.ph.i.1
	mov	r11, rsi
	.align	4, 0x90
LBB0_5:                                 ## %.preheader.i.5
	je	LBB0_6
	.align	4, 0x90
LBB0_16:                                ## %.lr.ph.i.12
                                        ##   Parent Loop BB0_5 Depth=1
                                        ## =>  This Inner Loop Header: Depth=2
	vmovdqu	ymm4, ymmword ptr [rdi + 4*rax]
	cmp	rax, r11
	jb	LBB0_16
LBB0_6:                                 ## %._crit_edge.i.6
                                        ##   in Loop: Header=BB0_5 Depth=1
	cmp	r11, rsi
	je	LBB0_8
## BB#7:                                ##   in Loop: Header=BB0_5 Depth=1
	vmovdqu	ymm4, ymmword ptr [rdi + 4*rsi - 128]
	vmovdqu	ymmword ptr [r8 + rsi - 32], ymm4
LBB0_8:                                 ##   in Loop: Header=BB0_5 Depth=1
	add	rdi, rcx
	jne	LBB0_5
	jmp	LBB0_15
LBB0_9:
	test	rdx, rdx
	je	LBB0_15
## BB#10:                               ## %.preheader.lr.ph.i
	mov	r11, rsi
	vpbroadcastd	ymm3, dword ptr [rip + LCPI0_3]
	.align	4, 0x90
LBB0_11:                                ## %.preheader.i
                                        ## =>This Loop Header: Depth=1
                                        ##     Child Loop BB0_17 Depth 2
	mov	eax, 0
	test	r11, r11
	je	LBB0_12
	.align	4, 0x90
LBB0_17:                                ## %.lr.ph.i
                                        ##   Parent Loop BB0_11 Depth=1
                                        ## =>  This Inner Loop Header: Depth=2
	vmovdqu	ymm4, ymmword ptr [rdi + 4*rax]
	vmovdqu	ymm5, ymmword ptr [rdi + 4*rax + 32]
	jb	LBB0_17
LBB0_12:                                ## %._crit_edge.i
                                        ##   in Loop: Header=BB0_11 Depth=1
	cmp	r11, rsi
	je	LBB0_14
## BB#13:                               ##   in Loop: Header=BB0_11 Depth=1
	vmovdqu	ymm4, ymmword ptr [rdi + 4*rsi - 128]
	vmovdqu	ymmword ptr [r8 + rsi - 32], ymm4
LBB0_14:                                ##   in Loop: Header=BB0_11 Depth=1
	add	rdi, rcx
	add	r8, r9
	inc	r10
	cmp	r10, rdx
	jne	LBB0_11
LBB0_15:                                ## %_ZN4Simd4Avx210BgraToGrayILb1EEEvPKhmmmPhm.exit
	pop	rbp
	vzeroupper
	ret

.subsections_via_symbols
`

	segments1 := []Segment{}
	segments1 = append(segments1, Segment{Name: "SimdAvx2BgraToGray", Start: 22, End: 95})

	testSegment(t, strings.Split(src1, "\n"), segments1)

	src2 := `	.section	__TEXT,__text,regular,pure_instructions
	.macosx_version_min 10, 11
	.intel_syntax noprefix
	.section	__TEXT,__const
	.align	5
LCPI0_0:
	.short	16                      ## 0x10
	.short	13074                   ## 0x3312
	.short	0                       ## 0x0
	.section	__TEXT,__text,regular,pure_instructions
	.globl	__ZN4Simd4Avx213Yuv444pToBgraEPKhmS2_mS2_mmmPhmh
	.align	4, 0x90
__ZN4Simd4Avx213Yuv444pToBgraEPKhmS2_mS2_mmmPhmh: ## @_ZN4Simd4Avx213Yuv444pToBgraEPKhmS2_mS2_mmmPhmh
## BB#0:
	push    rbp
	mov     rbp, rsp
	push    r15
	push    r14
	push    r13
	push    r12
	push    rbx
	and     rsp, -32
	sub     rsp, 192
	mov     qword ptr [rsp + 56], r9 ## 8-byte Spill
	mov     r9b, byte ptr [rbp + 48]
	mov     r15, qword ptr [rbp + 40]
	mov     r13, qword ptr [rbp + 32]
	mov     r10, qword ptr [rbp + 16]
	mov     rbx, rsi
	and     rbx, -32
	cmp     rbx, rsi
	jne     LBB0_14
### BB#1:
	mov	rbx, rdi
	cmp	rbx, r13
	jne	LBB0_14
## BB#8:
	movzx	eax, r9b
	cmp	qword ptr [rbp + 24], 0
	je	LBB0_20
## BB#9:                                ## %.preheader.lr.ph.i.1
	vinserti128	ymm14, ymm0, xmm0, 1
	vmovdqu	ymmword ptr [r13 + r9 + 96], ymm0
LBB0_13:                                ##   in Loop: Header=BB0_10 Depth=1
	add	rdi, rsi
	jb	LBB0_22
LBB0_17:                                ## %._crit_edge.i
	cmp	rbx, qword ptr [rbp + 16]
	cmp	r11, qword ptr [rbp + 24]
	jne	LBB0_16
LBB0_20:                                ## %_ZN4Simd4Avx213Yuv444pToBgraILb1EEEvPKhmS3_mS3_mmmPhmh.exit
	lea	rsp, [rbp - 40]
	pop	rbx
	pop     r12
	pop     r13
	pop     r14
	pop     r15
	pop	rbp
	vzeroupper
	ret

	.section	__TEXT,__const
	.align	5
LCPI1_0:
	.byte	0                       ## 0x0
	.space	1
	.space	1
	.space	1
LCPI1_13:
	.space	32
	.section	__TEXT,__text,regular,pure_instructions
	.globl	__ZN4Simd4Avx213Yuv420pToBgraEPKhmS2_mS2_mmmPhmh
	.align	4, 0x90
__ZN4Simd4Avx213Yuv420pToBgraEPKhmS2_mS2_mmmPhmh: ## @_ZN4Simd4Avx213Yuv420pToBgraEPKhmS2_mS2_mmmPhmh
## BB#0:
	push	rbp
	mov	rbp, rsp
	push	r15
	push    r14
	push    r13
	push    r12
	push    rbx
	and     rsp, -32
	sub     rsp, 864
	mov     qword ptr [rsp + 144], r9 ## 8-byte Spill
	mov     qword ptr [rsp + 152], rcx ## 8-byte Spill
	xor	r12d, r12d
	.align	4, 0x90
LBB1_12:                                ## %.lr.ph.i.18
                                        ##   Parent Loop BB1_10 Depth=1
	cmp	r15, r11
	jb	LBB1_12
LBB1_13:                                ## %._crit_edge.i.8
                                        ##   in Loop: Header=BB1_10 Depth=1
	vmovdqa	ymm7, ymm10
	vmovdqu	ymmword ptr [rax + rsi + 224], ymm0
LBB1_15:                                ##   in Loop: Header=BB1_10 Depth=1
	add	rdi, qword ptr [rsp + 192] ## 8-byte Folded Reload
	vmovdqa	ymm7, ymmword ptr [rip + LCPI1_7] ## ymm7 = <u,u,u,u,1,1,1,1,u,u,u,u,1,1,1,1>
	.align	4, 0x90
LBB1_18:                                ## %.preheader.i
                                        ## =>This Loop Header: Depth=1
	cmp	rsi, rbx
	jb	LBB1_23
LBB1_19:                                ## %._crit_edge.i
                                        ##   in Loop: Header=BB1_18 Depth=1
	vmovdqu	ymmword ptr [rax + rsi + 224], ymm0
LBB1_21:                                ##   in Loop: Header=BB1_18 Depth=1
	add	rdi, qword ptr [rsp + 96] ## 8-byte Folded Reload
	jb	LBB1_18
LBB1_22:                                ## %_ZN4Simd4Avx213Yuv420pToBgraILb1EEEvPKhmS3_mS3_mmmPhmh.exit
	lea	rsp, [rbp - 40]
	pop	rbx
	pop     r12
	pop     r13
	pop     r14
	pop     r15
	pop	rbp
	vzeroupper
	ret

	.section	__TEXT,__const
	.align	5
LCPI2_0:
	.byte	0                       ## 0x0
	.byte	2                       ## 0x2
	.byte	15                      ## 0xf
LCPI2_12:
	.space	1
	.space	1
	.section	__TEXT,__text,regular,pure_instructions
	.globl	__ZN4Simd4Avx213Yuv422pToBgraEPKhmS2_mS2_mmmPhmh
	.align	4, 0x90
__ZN4Simd4Avx213Yuv422pToBgraEPKhmS2_mS2_mmmPhmh: ## @_ZN4Simd4Avx213Yuv422pToBgraEPKhmS2_mS2_mmmPhmh
## BB#0:
	push	rbp
	mov     rbp, rsp
	push    r15
	push    r14
	push    r13
	push    r12
	push    rbx
	and     rsp, -32
	sub     rsp, 416
	mov     qword ptr [rsp + 184], rcx ## 8-byte Spill
	mov     qword ptr [rsp + 176], rsi ## 8-byte Spill
	mov     cl, byte ptr [rbp + 48]
	mov     r12, qword ptr [rbp + 40]
	mov     rax, qword ptr [rbp + 32]
	mov     r10, qword ptr [rbp + 16]
	jne	LBB2_14
## BB#1:
	mov	rsi, rdi
	jne	LBB2_14
## BB#8:
	movzx	ecx, cl
	cmp	qword ptr [rbp + 24], 0
	mov	rcx, r9
	je	LBB2_20
## BB#9:                                ## %.preheader.lr.ph.i.1
	vinserti128	ymm12, ymm0, xmm0, 1
	.align	4, 0x90
LBB2_10:                                ## %.preheader.i.7
	.align	4, 0x90
LBB2_21:                                ## %.lr.ph.i.16
	jb	LBB2_21
LBB2_11:                                ## %._crit_edge.i.8
	je	LBB2_13
## BB#12:                               ##   in Loop: Header=BB2_10 Depth=1
	vmovdqa	ymm15, ymm9
	vmovdqu	ymmword ptr [rax + r15 + 224], ymm0
LBB2_13:                                ##   in Loop: Header=BB2_10 Depth=1
	add	rdi, qword ptr [rsp + 176] ## 8-byte Folded Reload
	jmp	LBB2_20
LBB2_14:
	mov	qword ptr [rsp + 168], r9 ## 8-byte Spill
	je	LBB2_20
## BB#15:                               ## %.preheader.lr.ph.i
	vinserti128	ymm0, ymm0, xmm0, 1
	.align	4, 0x90
LBB2_16:                                ## %.preheader.i
                                        ## =>This Loop Header: Depth=1
	je	LBB2_17
	.align	4, 0x90
LBB2_22:                                ## %.lr.ph.i
	cmp	r15, rbx
	jb	LBB2_22
LBB2_17:                                ## %._crit_edge.i
                                        ##   in Loop: Header=BB2_16 Depth=1
	cmp	rbx, qword ptr [rbp + 16]
	je	LBB2_19
## BB#18:                               ##   in Loop: Header=BB2_16 Depth=1
	vpermq	ymm1, ymmword ptr [rdx + rsi], 216 ## ymm1 = mem[0,2,1,3]
	vmovdqu	ymmword ptr [rax + r13 + 224], ymm0
LBB2_19:                                ##   in Loop: Header=BB2_16 Depth=1
	add	rdi, qword ptr [rsp + 176] ## 8-byte Folded Reload
	jne	LBB2_16
LBB2_20:                                ## %_ZN4Simd4Avx213Yuv422pToBgraILb1EEEvPKhmS3_mS3_mmmPhmh.exit
	lea	rsp, [rbp - 40]
	pop	rbx
	pop     r12
	pop     r13
	pop     r14
	pop     r15
	pop	rbp
	vzeroupper
	ret

.subsections_via_symbols`

	segments2 := []Segment{}
	segments2 = append(segments2, Segment{Name: "SimdAvx2Yuv444pToBgra", Start: 13, End: 51})
	segments2 = append(segments2, Segment{Name: "SimdAvx2Yuv420pToBgra", Start: 74, End: 111})
	segments2 = append(segments2, Segment{Name: "SimdAvx2Yuv422pToBgra", Start: 134, End: 198})

	testSegment(t, strings.Split(src2, "\n"), segments2)

}
