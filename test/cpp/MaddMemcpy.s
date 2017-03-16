	.section	__TEXT,__text,regular,pure_instructions
	.macosx_version_min 10, 11
	.intel_syntax noprefix
	.globl	__Z10MaddMemcpyPfS_S_iiS_
	.align	4, 0x90
__Z10MaddMemcpyPfS_S_iiS_:              ## @_Z10MaddMemcpyPfS_S_iiS_
## BB#0:
	push	rbp
	mov	rbp, rsp
	push	r15
	push	r14
	push	r13
	push	r12
	push	rbx
	push	rax
	mov	r14, r9
	mov	r15d, r8d
	mov	r12, rdx
	mov	r13, rsi
	mov	rbx, rdi
	movsxd	rdx, ecx
	mov	rdi, r13
	mov	rsi, rbx
	call	_memcpy
	movsxd	rdx, r15d
	mov	rdi, r12
	mov	rsi, rbx
	call	_memcpy
	vmovups	ymm0, ymmword ptr [rbx]
	vmovups	ymm1, ymmword ptr [r13]
	vfmadd213ps	ymm1, ymm0, ymmword ptr [r12]
	vmovups	ymmword ptr [r14], ymm1
	add	rsp, 8
	pop	rbx
	pop	r12
	pop	r13
	pop	r14
	pop	r15
	pop	rbp
	vzeroupper
	ret


.subsections_via_symbols
