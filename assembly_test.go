package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestAssemblyAlignedWithTableWithStackArgs(t *testing.T) {

	epilogue := Epilogue{SetRbpInstr: true, StackSize: 64, AlignedStack: true, AlignValue: 16, VZeroUpper: false}
	epilogue.Pops = append(epilogue.Pops, "rbp", "r15", "r14", "r13", "r12", "rbx")
	table := Table{Name: "LCDATA3"}

	subroutine := Subroutine{name: "SimdSse2MedianFilterRhomb5x5", epilogue: epilogue, table: table}
	arguments, returnValues := []string{}, []string{}
	arguments = append(arguments, "src", "srcStride", "width", "height", "channelCount", "dst", "dstStride")

	lines := writeGoasmPrologue(subroutine, arguments, returnValues)
	lines = append(lines, writeGoasmEpilogue(subroutine, arguments, returnValues)...)

	alignedWithTable := `TEXT ·_SimdSse2MedianFilterRhomb5x5(SB), 7, $96-56

	MOVQ SP, BP
	ANDQ $-16, BP
	MOVQ SP, 88(BP)
	MOVQ dstStride+48(FP), DI
	MOVQ DI, 80(BP)
	MOVQ src+0(FP), DI
	MOVQ srcStride+8(FP), SI
	MOVQ width+16(FP), DX
	MOVQ height+24(FP), CX
	MOVQ channelCount+32(FP), R8
	MOVQ dst+40(FP), R9
	LEAQ LCDATA3<>(SB), BP
	ADDQ $16, SP
	ANDQ $-16, SP

	MOVQ 72(SP), SP
	RET`

	for i, l := range lines {
		goldenLine := strings.Split(alignedWithTable, "\n")[i]
		if strings.TrimSpace(l) != strings.TrimSpace(goldenLine) {
			t.Errorf("TestAssemblyAlignedWithTableWithStackArgs(): \nexpected %s\ngot %s", goldenLine, l)
		}
	}

	test := "mov	r11, qword ptr [rbp + 16]"

	result := fixRbpPlusLoad(test, StackArgs{Number: 1, OffsetToFirst: 16}, epilogue.getStackpointerDecrement(table, len(arguments)), epilogue.additionalStackSpace(table, len(arguments)), table.isPresent(), epilogue.AlignedStack)

	dstStrideMov := `mov	r11, qword ptr 64[rsp] /* [rbp + 16] */`
	if dstStrideMov != result {
		t.Errorf("TestAssemblyAlignedWithTableWithStackArgs(): \nexpected %s\ngot %s", dstStrideMov, result)
	}
}

func TestAssemblyUnalignedWithTableWithStackArgs(t *testing.T) {

	epilogue := Epilogue{SetRbpInstr: true, StackSize: 8, AlignedStack: false, AlignValue: 0, VZeroUpper: false}
	epilogue.Pops = append(epilogue.Pops, "rbp", "r15", "r14", "r13", "r12", "rbx")
	table := Table{Name: "LCDATA2"}

	subroutine := Subroutine{name: "SimdSse2MedianFilterSquare3x3", epilogue: epilogue, table: table}
	arguments, returnValues := []string{}, []string{}
	arguments = append(arguments, "src", "srcStride", "width", "height", "channelCount", "dst", "dstStride")

	lines := writeGoasmPrologue(subroutine, arguments, returnValues)
	lines = append(lines, writeGoasmEpilogue(subroutine, arguments, returnValues)...)

	unalignedWithTable := `TEXT ·_SimdSse2MedianFilterSquare3x3(SB), 7, $8-56

	MOVQ src+0(FP), DI
	MOVQ srcStride+8(FP), SI
	MOVQ width+16(FP), DX
	MOVQ height+24(FP), CX
	MOVQ channelCount+32(FP), R8
	MOVQ dst+40(FP), R9
	LEAQ LCDATA2<>(SB), BP
	ADDQ $8, SP

	SUBQ $8, SP
	RET`

	for i, l := range lines {
		goldenLine := strings.Split(unalignedWithTable, "\n")[i]
		if strings.TrimSpace(l) != strings.TrimSpace(goldenLine) {
			t.Errorf("TestAssemblyUnalignedWithTableWithStackArgs(): \nexpected %s\ngot %s", goldenLine, l)
		}
	}

	test := "mov    rax, qword ptr [rbp + 16]"

	result := fixRbpPlusLoad(test, StackArgs{Number: 1, OffsetToFirst: 16}, epilogue.getStackpointerDecrement(table, len(arguments)), epilogue.additionalStackSpace(table, len(arguments)), table.isPresent(), epilogue.AlignedStack)

	dstStrideMov := `mov    rax, qword ptr 64[rsp] /* [rbp + 16] */`
	if dstStrideMov != result {
		t.Errorf("TestAssemblyUnalignedWithTableWithStackArgs(): \nexpected %s\ngot %s", dstStrideMov, result)

	}
}

func TestRbpPlusLoad(t *testing.T) {

	tests := `mov	r8, qword ptr [rbp + 24]
	mov	r11, qword ptr [rbp + 16]
	movd	xmm11, dword ptr [rbp + 32] ## xmm11 = mem[0],zero,zero,zero
	movd	xmm8, dword ptr [rbp + 40] ## xmm8 = mem[0],zero,zero,zero
	movd	xmm9, dword ptr [rbp + 48] ## xmm9 = mem[0],zero,zero,zero
	movd	xmm10, dword ptr [rbp + 56] ## xmm10 = mem[0],zero,zero,zero
	mov	rdi, qword ptr [rbp + 24]
	add	r15, qword ptr [rbp + 16]
	mov	rdi, qword ptr [rbp + 24]
	add	r15, qword ptr [rbp + 16]`

	stackArgs := StackArgs{OffsetToFirst: 256}
	for _, test := range strings.Split(tests, "\n") {
		test, _ = stripComments(test)
		result := fixRbpPlusLoad(test, stackArgs, 0, 0, true, true)

		if !(strings.Contains(result, `/*`) && strings.Contains(result, `*/`)) {
			t.Errorf("TestRbpPlusLoad(): \nexpected to find C-style comment\ngot %s", result)
		}
	}
}

func TestStripComments(t *testing.T) {

	stripped, notskip := stripComments(`	mov     qword ptr [rsp + 144], r9 ## 8-byte Spill`)
	expected := `	mov     qword ptr [rsp + 144], r9`
	if stripped != expected || notskip {
		t.Errorf("TestStripComments(): \nexpected %s\ngot %s", expected, stripped)
	}

	stripped2, notskip2 := stripComments(`	movdqa	xmm3, xmmword ptr [rip + .LCPI2_10] # xmm3 = [59507,8192,59507,8192,59507,8192,59507,8192]`)
	expected = `	movdqa	xmm3, xmmword ptr [rip + .LCPI2_10]`
	if stripped2 != expected || notskip2 {
		t.Errorf("TestStripComments(): \nexpected %s\ngot %s", expected, stripped2)
	}

	empty, skip := stripComments(`                                        ## =>This Loop Header: Depth=1`)
	expected = ``
	if empty != expected || !skip {
		t.Errorf("TestStripComments(): \nexpected %s\ngot %s", expected, empty)
	}

}
