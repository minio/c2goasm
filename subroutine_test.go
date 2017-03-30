package main

import (
	"strings"
	"testing"
	"github.com/cloudflare/go/src/fmt"
)

func testName(t *testing.T, fullname, expected string) {
	name := extractName(fullname)
	if name != expected {
		t.Errorf("TestNames(): \nexpected %s\ngot      %s", expected, name)
	}
}

func TestNames(t *testing.T) {

	testName(t, "_ZN4Simd4Avx213Yuv444pToBgraEPKhmS2_mS2_mmmPhmh", "SimdAvx2Yuv444pToBgra")
	testName(t, "_ZN4Simd4Avx213Yuv420pToBgraEPKhmS2_mS2_mmmPhmh", "SimdAvx2Yuv420pToBgra")
	testName(t, "_ZN4Simd4Avx213Yuv422pToBgraEPKhmS2_mS2_mmmPhmh", "SimdAvx2Yuv422pToBgra")
	testName(t, "_ZN4Simd4Avx213ReduceGray2x2EPKhmmmPhmmm", "SimdAvx2ReduceGray2x2")
	testName(t, "_ZN4Simd4Avx216AbsDifferenceSumEPKhmS2_mmmPy", "SimdAvx2AbsDifferenceSum")
}

func subroutineEqual(a, b []Subroutine) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !(a[i].name == b[i].name && equalString(a[i].body, b[i].body)) {
			return false
		}
	}

	return true
}

func testSubroutine(t *testing.T, fullsrc []string, expected []Subroutine) {
	subroutines := segmentSource(fullsrc)
	if !subroutineEqual(subroutines, expected) {
		t.Errorf("testSubroutine(): \nexpected %#v\ngot      %#v", expected, subroutines)
	}
}

func TestSubroutine(t *testing.T) {

	disabledForTesting = true

	src1 := strings.Split(`	.section	__TEXT,__text,regular,pure_instructions
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
`, "\n")

	subroutine1 := []Subroutine{}
	subroutine1 = append(subroutine1, Subroutine{name: "SimdAvx2BgraToGray", body: src1[25:98]})

	testSubroutine(t, src1, subroutine1)

	src2 := strings.Split(`	.section	__TEXT,__text,regular,pure_instructions
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

.subsections_via_symbols`, "\n")

	subroutine2 := []Subroutine{}
	subroutine2 = append(subroutine2, Subroutine{name: "SimdAvx2Yuv444pToBgra", body: src2[23:60]})
	subroutine2 = append(subroutine2, Subroutine{name: "SimdAvx2Yuv420pToBgra", body: src2[84:120]})
	subroutine2 = append(subroutine2, Subroutine{name: "SimdAvx2Yuv422pToBgra", body: src2[144:207]})

	testSubroutine(t, src2, subroutine2)

	src3 := strings.Split(`        .globl  __ZN4Simd4Avx214MultiplyAndAddEPfS1_S1_S1_
        .align  4, 0x90
__ZN4Simd4Avx214MultiplyAndAddEPfS1_S1_S1_: ## @_ZN4Simd4Avx214MultiplyAndAddEPfS1_S1_S1_
## BB#0:
        push    rbp
        mov     rbp, rsp
        vmovups ymm0, ymmword ptr [rdi]
        vmovups ymm1, ymmword ptr [rsi]
        vfmadd213ps     ymm1, ymm0, ymmword ptr [rdx]
        vmovups ymmword ptr [rcx], ymm1
        pop     rbp
        vzeroupper
        ret

.subsections_via_symbols`, "\n")

	subroutine3 := []Subroutine{}
	subroutine3 = append(subroutine3, Subroutine{name: "SimdAvx2MultiplyAndAdd", body: src3[6:13]})

	testSubroutine(t, src3, subroutine3)

	src4 := strings.Split(`        .section        __TEXT,__text,regular,pure_instructions
        .macosx_version_min 10, 11
        .intel_syntax noprefix
        .globl  __Z22MultiplyAndAddConstantPfS_S_
        .align  4, 0x90
__Z22MultiplyAndAddConstantPfS_S_:      ## @_Z22MultiplyAndAddConstantPfS_S_
## BB#0:
        push    rbp
        mov     rbp, rsp
        vmovups ymm0, ymmword ptr [rdi]
        vmovups ymm1, ymmword ptr [rsi]
        vfmadd213ps     ymm1, ymm0, ymmword ptr [rip + __ZL1a]
        vmovups ymmword ptr [rdx], ymm1
        pop     rbp
        vzeroupper
        ret

        .section        __DATA,__data
        .align  5                       ## @_ZL1a
__ZL1a:
        .long   1065353216              ## float 1.000000e+00
        .long   1073741824              ## float 2.000000e+00
        .long   1077936128              ## float 3.000000e+00
        .long   1082130432              ## float 4.000000e+00
        .long   1084227584              ## float 5.000000e+00
        .long   1086324736              ## float 6.000000e+00
        .long   1088421888              ## float 7.000000e+00
        .long   1090519040              ## float 8.000000e+00
`, "\n")

	subroutine4 := []Subroutine{}
	subroutine4 = append(subroutine4, Subroutine{name: "MultiplyAndAddConstant", body: src4[9:16]})

	testSubroutine(t, src4, subroutine4)

	subroutine5 := []Subroutine{}
	subroutine5 = append(subroutine5, Subroutine{name: "SimdSse2BgraToYuv420p", body: srcOsx[43:53]})
	subroutine5 = append(subroutine5, Subroutine{name: "SimdSse2BgraToYuv422p", body: srcOsx[94:103]})
	subroutine5 = append(subroutine5, Subroutine{name: "SimdSse2BgraToYuv444p", body: srcOsx[142:151]})

	testSubroutine(t, srcOsx, subroutine5)

	subroutine6 := []Subroutine{}
	subroutine6 = append(subroutine6, Subroutine{name: "SimdSse2BgraToYuv420p", body: srcClang[41:51]})
	subroutine6 = append(subroutine6, Subroutine{name: "SimdSse2BgraToYuv422p", body: srcClang[94:103]})
	subroutine6 = append(subroutine6, Subroutine{name: "SimdSse2BgraToYuv444p", body: srcClang[144:153]})

	testSubroutine(t, srcClang, subroutine6)
}

var srcClang = strings.Split(`	.text
	.intel_syntax noprefix
	.section	.rodata.cst16,"aM",@progbits,16
	.align	16
.LCPI0_0:
	.byte	255                     # 0xff
.LCPI0_1:
	.byte	255                     # 0xff
	.byte	0                       # 0x0
.LCPI0_2:
	.quad	281474976776192         # 0x1000000010000
.LCPI0_3:
	.short	1606                    # 0x646
	.short	4211                    # 0x1073
.LCPI0_4:
	.short	8258                    # 0x2042
.LCPI0_5:
	.short	16                      # 0x10
.LCPI0_6:
	.short	2                       # 0x2
.LCPI0_7:
	.short	7193                    # 0x1c19
.LCPI0_8:
	.short	60768                   # 0xed60
.LCPI0_9:
	.short	128                     # 0x80
.LCPI0_10:
	.short	64373                   # 0xfb75
.LCPI0_11:
	.short	59507                   # 0xe873
.LCPI0_12:
	.zero	16
	.text
	.globl	_ZN4Simd4Sse213BgraToYuv420pEPKhmmmPhmS3_mS3_m
	.align	16, 0x90
	.type	_ZN4Simd4Sse213BgraToYuv420pEPKhmmmPhmS3_mS3_m,@function
_ZN4Simd4Sse213BgraToYuv420pEPKhmmmPhmS3_mS3_m: # @_ZN4Simd4Sse213BgraToYuv420pEPKhmmmPhmS3_mS3_m
# BB#0:
	push	rbp
	push	r15
	push	r14

.LBB0_24:                               # %_ZN4Simd4Sse213BgraToYuv420pILb1EEEvPKhmmmPhmS4_mS4_m.exit
	add	rsp, 136
	pop	rbx
	pop	r12
	pop	r13
	pop	r14
	pop	r15
	pop	rbp
	ret
.Lfunc_end0:
	.size	_ZN4Simd4Sse213BgraToYuv420pEPKhmmmPhmS3_mS3_m, .Lfunc_end0-_ZN4Simd4Sse213BgraToYuv420pEPKhmmmPhmS3_mS3_m

	.section	.rodata.cst16,"aM",@progbits,16
	.align	16
.LCPI1_0:
	.byte	255                     # 0xff
.LCPI1_1:
	.byte	0                       # 0x0
.LCPI1_2:
	.quad	281474976776192         # 0x1000000010000
.LCPI1_3:
	.short	4211                    # 0x1073
.LCPI1_4:
	.short	8192                    # 0x2000
.LCPI1_5:
	.short	16                      # 0x10
.LCPI1_6:
	.short	1                       # 0x1
.LCPI1_7:
	.short	7193                    # 0x1c19
.LCPI1_8:
	.short	60768                   # 0xed60
.LCPI1_9:
	.short	128                     # 0x80
.LCPI1_10:
	.short	64373                   # 0xfb75
.LCPI1_11:
	.short	59507                   # 0xe873
.LCPI1_12:
	.zero	16
	.text
	.globl	_ZN4Simd4Sse213BgraToYuv422pEPKhmmmPhmS3_mS3_m
	.align	16, 0x90
	.type	_ZN4Simd4Sse213BgraToYuv422pEPKhmmmPhmS3_mS3_m,@function
_ZN4Simd4Sse213BgraToYuv422pEPKhmmmPhmS3_mS3_m: # @_ZN4Simd4Sse213BgraToYuv422pEPKhmmmPhmS3_mS3_m
# BB#0:
	push	rbp
	push	r15
	push	r14
	push	r13
	push	r12
	push	rbx

.LBB1_20:                               # %_ZN4Simd4Sse213BgraToYuv422pILb1EEEvPKhmmmPhmS4_mS4_m.exit
	pop	rbx
	pop	r12
	pop	r13
	pop	r14
	pop	r15
	pop	rbp
	ret
.Lfunc_end1:
	.size	_ZN4Simd4Sse213BgraToYuv422pEPKhmmmPhmS3_mS3_m, .Lfunc_end1-_ZN4Simd4Sse213BgraToYuv422pEPKhmmmPhmS3_mS3_m

	.section	.rodata.cst16,"aM",@progbits,16
	.align	16
.LCPI2_0:
	.byte	0                       # 0x0
.LCPI2_1:
	.byte	255                     # 0xff
.LCPI2_2:
	.quad	281474976776192         # 0x1000000010000
.LCPI2_3:
	.short	1606                    # 0x646
.LCPI2_4:
	.short	8258                    # 0x2042
.LCPI2_5:
	.short	16                      # 0x10
.LCPI2_6:
	.short	7193                    # 0x1c19
.LCPI2_7:
	.short	60768                   # 0xed60
.LCPI2_8:
	.short	128                     # 0x80
.LCPI2_9:
	.short	64373                   # 0xfb75
.LCPI2_10:
	.short	59507                   # 0xe873
.LCPI2_11:
	.zero	16
	.text
	.globl	_ZN4Simd4Sse213BgraToYuv444pEPKhmmmPhmS3_mS3_m
	.align	16, 0x90
	.type	_ZN4Simd4Sse213BgraToYuv444pEPKhmmmPhmS3_mS3_m,@function
_ZN4Simd4Sse213BgraToYuv444pEPKhmmmPhmS3_mS3_m: # @_ZN4Simd4Sse213BgraToYuv444pEPKhmmmPhmS3_mS3_m
# BB#0:
	push	rbp
	push	r15
	push	r14
	push	r13
	push	r12
	push	rbx

.LBB2_20:                               # %_ZN4Simd4Sse213BgraToYuv444pILb1EEEvPKhmmmPhmS4_mS4_m.exit
	pop	rbx
	pop	r12
	pop	r13
	pop	r14
	pop	r15
	pop	rbp
	ret
.Lfunc_end2:
	.size	_ZN4Simd4Sse213BgraToYuv444pEPKhmmmPhmS3_mS3_m, .Lfunc_end2-_ZN4Simd4Sse213BgraToYuv444pEPKhmmmPhmS3_mS3_m


	.ident	"clang version 3.8.0-2ubuntu4 (tags/RELEASE_380/final)"
	.section	".note.GNU-stack","",@progbits`, "\n")

var srcOsx = strings.Split(`	.section	__TEXT,__text,regular,pure_instructions
	.macosx_version_min 10, 11
	.intel_syntax noprefix
	.section	__TEXT,__literal16,16byte_literals
	.align	4
LCPI0_0:
	.byte	255                     ## 0xff
LCPI0_1:
	.byte	1                       ## 0x1
LCPI0_2:
	.quad	281474976776192         ## 0x1000000010000
LCPI0_3:
	.short	1606                    ## 0x646
LCPI0_4:
	.short	8258                    ## 0x2042
LCPI0_5:
	.short	16                      ## 0x10
LCPI0_6:
	.short	2                       ## 0x2
LCPI0_7:
	.short	7193                    ## 0x1c19
LCPI0_8:
	.short	60768                   ## 0xed60
LCPI0_9:
	.short	128                     ## 0x80
LCPI0_10:
	.short	64373                   ## 0xfb75
LCPI0_11:
	.short	59507                   ## 0xe873
LCPI0_12:
	.space	16
	.section	__TEXT,__text,regular,pure_instructions
	.globl	__ZN4Simd4Sse213BgraToYuv420pEPKhmmmPhmS3_mS3_m
	.align	4, 0x90
__ZN4Simd4Sse213BgraToYuv420pEPKhmmmPhmS3_mS3_m: ## @_ZN4Simd4Sse213BgraToYuv420pEPKhmmmPhmS3_mS3_m
## BB#0:
	push	rbp
	mov	rbp, rsp
	push	r15
	push	r14
	push	r13
	push	r12
	push	rbx

LBB0_24:                                ## %_ZN4Simd4Sse213BgraToYuv420pILb1EEEvPKhmmmPhmS4_mS4_m.exit
	add	rsp, 88
	pop	rbx
	pop	r12
	pop	r13
	pop	r14
	pop	r15
	pop	rbp
	ret

	.section	__TEXT,__literal16,16byte_literals
	.align	4
LCPI1_0:
	.byte	255                     ## 0xff
LCPI1_1:
	.byte	1                       ## 0x1
LCPI1_2:
	.quad	281474976776192         ## 0x1000000010000
LCPI1_3:
	.short	1606                    ## 0x646
LCPI1_4:
	.short	8258                    ## 0x2042
LCPI1_5:
	.short	16                      ## 0x10
LCPI1_6:
	.short	1                       ## 0x1
LCPI1_7:
	.short	7193                    ## 0x1c19
LCPI1_8:
	.short	60768                   ## 0xed60
LCPI1_9:
	.short	128                     ## 0x80
LCPI1_10:
	.short	64373                   ## 0xfb75
LCPI1_11:
	.short	59507                   ## 0xe873
LCPI1_12:
	.space	16
	.section	__TEXT,__text,regular,pure_instructions
	.globl	__ZN4Simd4Sse213BgraToYuv422pEPKhmmmPhmS3_mS3_m
	.align	4, 0x90
__ZN4Simd4Sse213BgraToYuv422pEPKhmmmPhmS3_mS3_m: ## @_ZN4Simd4Sse213BgraToYuv422pEPKhmmmPhmS3_mS3_m
## BB#0:
	push	rbp
	mov	rbp, rsp
	push	r15
	push	r14
	push	r13
	push	r12
	push	rbx

LBB1_20:                                ## %_ZN4Simd4Sse213BgraToYuv422pILb1EEEvPKhmmmPhmS4_mS4_m.exit
	pop	rbx
	pop	r12
	pop	r13
	pop	r14
	pop	r15
	pop	rbp
	ret

	.section	__TEXT,__literal16,16byte_literals
	.align	4
LCPI2_0:
	.byte	255                     ## 0xff
LCPI2_1:
	.byte	1                       ## 0x1
LCPI2_2:
	.quad	281474976776192         ## 0x1000000010000
LCPI2_3:
	.short	1606                    ## 0x646
LCPI2_4:
	.short	8258                    ## 0x2042
LCPI2_5:
	.short	16                      ## 0x10
LCPI2_6:
	.short	7193                    ## 0x1c19
LCPI2_7:
	.short	60768                   ## 0xed60
LCPI2_8:
	.short	128                     ## 0x80
LCPI2_9:
	.short	64373                   ## 0xfb75
LCPI2_10:
	.short	59507                   ## 0xe873
LCPI2_11:
	.space	16
	.section	__TEXT,__text,regular,pure_instructions
	.globl	__ZN4Simd4Sse213BgraToYuv444pEPKhmmmPhmS3_mS3_m
	.align	4, 0x90
__ZN4Simd4Sse213BgraToYuv444pEPKhmmmPhmS3_mS3_m: ## @_ZN4Simd4Sse213BgraToYuv444pEPKhmmmPhmS3_mS3_m
## BB#0:
	push	rbp
	mov	rbp, rsp
	push	r15
	push	r14
	push	r13
	push	r12
	push	rbx

LBB2_20:                                ## %_ZN4Simd4Sse213BgraToYuv444pILb1EEEvPKhmmmPhmS4_mS4_m.exit
	pop	rbx
	pop	r12
	pop	r13
	pop	r14
	pop	r15
	pop	rbp
	ret


.subsections_via_symbols`, "\n")
