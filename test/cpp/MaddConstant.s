	.section	__TEXT,__text,regular,pure_instructions
	.macosx_version_min 10, 11
	.intel_syntax noprefix
	.section	__TEXT,__const
	.align	5
LCPI0_0:
	.long	1065353216              ## float 1.000000e+00
	.long	1073741824              ## float 2.000000e+00
	.long	1077936128              ## float 3.000000e+00
	.long	1082130432              ## float 4.000000e+00
	.long	1084227584              ## float 5.000000e+00
	.long	1086324736              ## float 6.000000e+00
	.long	1088421888              ## float 7.000000e+00
	.long	1090519040              ## float 8.000000e+00
	.section	__TEXT,__text,regular,pure_instructions
	.globl	__Z12MaddConstantPfS_S_
	.align	4, 0x90
__Z12MaddConstantPfS_S_:                ## @_Z12MaddConstantPfS_S_
## BB#0:
	push	rbp
	mov	rbp, rsp
	vmovups	ymm0, ymmword ptr [rdi]
	vmovups	ymm1, ymmword ptr [rsi]
	vfmadd213ps	ymm1, ymm0, ymmword ptr [rip + LCPI0_0]
	vmovups	ymmword ptr [rdx], ymm1
	pop	rbp
	vzeroupper
	ret


.subsections_via_symbols
