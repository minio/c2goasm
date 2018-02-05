/*
 * Minio Cloud Storage, (C) 2017 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
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

	stack := NewStack(epilogue, len(arguments), 0)

	lines := writeGoasmPrologue(subroutine, stack, arguments, returnValues)
	lines = append(lines, writeGoasmEpilogue(subroutine, stack, arguments, returnValues)...)

	alignedWithTable := `TEXT ·_SimdSse2MedianFilterRhomb5x5(SB), $96-56

	MOVQ src+0(FP), DI
	MOVQ srcStride+8(FP), SI
	MOVQ width+16(FP), DX
	MOVQ height+24(FP), CX
	MOVQ channelCount+32(FP), R8
	MOVQ dst+40(FP), R9
	MOVQ dstStride+48(FP), R10
	MOVQ SP, BP
	ADDQ $16, SP
	ANDQ $-16, SP
	MOVQ BP, 72(SP)
	MOVQ R10, 64(SP)
	LEAQ LCDATA3<>(SB), BP

	MOVQ 72(SP), SP
	RET`

	for i, l := range lines {
		goldenLine := strings.Split(alignedWithTable, "\n")[i]
		if strings.TrimSpace(l) != strings.TrimSpace(goldenLine) {
			t.Errorf("TestAssemblyAlignedWithTableWithStackArgs(): \nexpected %s\ngot %s", goldenLine, l)
		}
	}

	test := "mov	r11, qword ptr [rbp + 16]"

	result := fixRbpPlusLoad(test, StackArgs{Number: 1, OffsetToFirst: 16}, stack)

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

	stack := NewStack(epilogue, len(arguments), 0)

	lines := writeGoasmPrologue(subroutine, stack, arguments, returnValues)
	lines = append(lines, writeGoasmEpilogue(subroutine, stack, arguments, returnValues)...)

	unalignedWithTable := `TEXT ·_SimdSse2MedianFilterSquare3x3(SB), $24-56

	MOVQ src+0(FP), DI
	MOVQ srcStride+8(FP), SI
	MOVQ width+16(FP), DX
	MOVQ height+24(FP), CX
	MOVQ channelCount+32(FP), R8
	MOVQ dst+40(FP), R9
	MOVQ dstStride+48(FP), R10
	ADDQ $8, SP
	MOVQ R10, 8(SP)
	LEAQ LCDATA2<>(SB), BP

	SUBQ $8, SP
	RET`

	for i, l := range lines {
		goldenLine := strings.Split(unalignedWithTable, "\n")[i]
		if strings.TrimSpace(l) != strings.TrimSpace(goldenLine) {
			t.Errorf("TestAssemblyUnalignedWithTableWithStackArgs(): \nexpected %s\ngot %s", goldenLine, l)
		}
	}

	test := "mov    rax, qword [rbp + 16]"

	result := fixRbpPlusLoad(test, StackArgs{Number: 1, OffsetToFirst: 16}, stack)

	dstStrideMov := `mov    rax, qword 8[rsp] /* [rbp + 16] */`
	if dstStrideMov != result {
		t.Errorf("TestAssemblyUnalignedWithTableWithStackArgs(): \nexpected %s\ngot %s", dstStrideMov, result)

	}
}

func TestAssemblyUnalignedWithTableWithStackArgsWithStackZeroSize(t *testing.T) {

	epilogue := Epilogue{SetRbpInstr: true, StackSize: 0, AlignedStack: false, AlignValue: 0, VZeroUpper: false}
	epilogue.Pops = append(epilogue.Pops, "rbp", "r15", "r14", "r13", "r12", "rbx")
	table := Table{Name: "LCDATA1"}

	subroutine := Subroutine{name: "SimdSse2MedianFilterRhomb3x3", epilogue: epilogue, table: table}
	arguments, returnValues := []string{}, []string{}
	arguments = append(arguments, "src", "srcStride", "width", "height", "channelCount", "dst", "dstStride")

	stack := NewStack(epilogue, len(arguments), 0)

	lines := writeGoasmPrologue(subroutine, stack, arguments, returnValues)
	lines = append(lines, writeGoasmEpilogue(subroutine, stack, arguments, returnValues)...)

	unalignedWithTableWithStackZeroSize := `TEXT ·_SimdSse2MedianFilterRhomb3x3(SB), $16-56

	MOVQ src+0(FP), DI
	MOVQ srcStride+8(FP), SI
	MOVQ width+16(FP), DX
	MOVQ height+24(FP), CX
	MOVQ channelCount+32(FP), R8
	MOVQ dst+40(FP), R9
	MOVQ dstStride+48(FP), R10
	ADDQ $8, SP
	MOVQ R10, 0(SP)
	LEAQ LCDATA1<>(SB), BP

	SUBQ $8, SP
	RET`

	for i, l := range lines {
		goldenLine := strings.Split(unalignedWithTableWithStackZeroSize, "\n")[i]
		if strings.TrimSpace(l) != strings.TrimSpace(goldenLine) {
			t.Errorf("TestAssemblyUnalignedWithTableWithStackArgsWithStackZeroSize(): \nexpected %s\ngot %s", goldenLine, l)
		}
	}

	test := "mov    rax, qword [rbp + 16]"

	result := fixRbpPlusLoad(test, StackArgs{Number: 1, OffsetToFirst: 16}, stack)

	dstStrideMov := `mov    rax, qword 0[rsp] /* [rbp + 16] */`
	if dstStrideMov != result {
		t.Errorf("TestAssemblyUnalignedWithTableWithStackArgsWithStackZeroSize(): \nexpected %s\ngot %s", dstStrideMov, result)

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
