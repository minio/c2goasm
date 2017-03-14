	.section	__TEXT,__text,regular,pure_instructions
	.macosx_version_min 10, 11
	.intel_syntax noprefix
	.globl	__Z22MultiplyAndAddConstantPfS_S_
	.align	4, 0x90
__Z22MultiplyAndAddConstantPfS_S_:      ## @_Z22MultiplyAndAddConstantPfS_S_
## BB#0:
	push	rbp
	mov	rbp, rsp
	vmovups	ymm0, ymmword ptr [rdi]
	vmovups	ymm1, ymmword ptr [rsi]
	vfmadd213ps	ymm1, ymm0, ymmword ptr [rip + __ZL1a]
	vmovups	ymmword ptr [rdx], ymm1
	pop	rbp
	vzeroupper
	ret

	.section	__DATA,__data
	.align	5                       ## @_ZL1a
__ZL1a:
	.long	1065353216              ## float 1.000000e+00
	.long	1073741824              ## float 2.000000e+00
	.long	1077936128              ## float 3.000000e+00
	.long	1082130432              ## float 4.000000e+00
	.long	1084227584              ## float 5.000000e+00
	.long	1086324736              ## float 6.000000e+00
	.long	1088421888              ## float 7.000000e+00
	.long	1090519040              ## float 8.000000e+00


.subsections_via_symbols
