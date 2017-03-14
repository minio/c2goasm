	.section	__TEXT,__text,regular,pure_instructions
	.macosx_version_min 10, 11
	.intel_syntax noprefix
	.globl	__Z14MultiplyAndAddPfS_S_S_
	.align	4, 0x90
__Z14MultiplyAndAddPfS_S_S_:            ## @_Z14MultiplyAndAddPfS_S_S_
## BB#0:
	push	rbp
	mov	rbp, rsp
	vmovups	ymm0, ymmword ptr [rdi]
	vmovups	ymm1, ymmword ptr [rsi]
	vfmadd213ps	ymm1, ymm0, ymmword ptr [rdx]
	vmovups	ymmword ptr [rcx], ymm1
	pop	rbp
	vzeroupper
	ret


.subsections_via_symbols
