	.section	__TEXT,__text,regular,pure_instructions
	.macosx_version_min 10, 11
	.intel_syntax noprefix
	.globl	__Z10MaddArgs10PfS_S_S_S_S_S_S_S_S_
	.align	4, 0x90
__Z10MaddArgs10PfS_S_S_S_S_S_S_S_S_:    ## @_Z10MaddArgs10PfS_S_S_S_S_S_S_S_S_
## BB#0:
	push	rbp
	mov	rbp, rsp
	push	rbx
	mov	r10, qword ptr [rbp + 40]
	mov	r11, qword ptr [rbp + 32]
	mov	rax, qword ptr [rbp + 16]
	mov	rbx, qword ptr [rbp + 24]
	vmovups	ymm0, ymmword ptr [rdi]
	vmovups	ymm1, ymmword ptr [rsi]
	vmovups	ymm2, ymmword ptr [rcx]
	vmovups	ymm3, ymmword ptr [r9]
	vmovups	ymm4, ymmword ptr [rbx]
	vfmadd213ps	ymm1, ymm0, ymmword ptr [rdx]
	vfmadd213ps	ymm1, ymm2, ymmword ptr [r8]
	vfmadd213ps	ymm1, ymm3, ymmword ptr [rax]
	vfmadd213ps	ymm1, ymm4, ymmword ptr [r11]
	vmovups	ymmword ptr [r10], ymm1
	pop	rbx
	pop	rbp
	vzeroupper
	ret


.subsections_via_symbols
