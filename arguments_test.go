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

func testArguments(t *testing.T, source string, expected StackArgs) {

	argsOnStack := argumentsOnStack(strings.Split(source, "\n"))

	if argsOnStack != expected {
		t.Errorf("testArguments(): \nexpected %#v\ngot      %#v", expected, argsOnStack)
	}
}

func TestArguments(t *testing.T) {

	nostackargs := `
     vmovups       ymm0, ymmword ptr [rdi]
     vmovups       ymm1, ymmword ptr [rsi]
     vfmadd213ps   ymm1, ymm0, ymmword ptr [rdx]
     vmovups       ymmword ptr [rcx], ymm1
    `

	testArguments(t, nostackargs, StackArgs{Number: 0, OffsetToFirst: 0})

	arguments7 := `
	 mov rax, qword [rbp + 16]
	 vmovups     ymm0, [rdi]
	 vmovups     ymm1, [rsi]
	 vmovups     ymm2, [r8]
	 vfmadd213ps ymm1, ymm0, [rdx]
	 vfmadd132ps ymm1, ymm2, [rcx]
	 vfmadd213ps ymm1, ymm2, [r9]
	 vmovups     [rax], ymm1
	`

	testArguments(t, arguments7, StackArgs{Number: 1, OffsetToFirst: 16})

	arguments8 := `
	 mov r10, qword [rbp + 24]
	 mov rax, qword [rbp + 16]
	 vmovups     ymm0, [rdi]
	 vmovups     ymm1, [rsi]
	 vmovups     ymm2, [rcx]
	 vmovups     ymm3, [r9]
	 vfmadd213ps ymm1, ymm0, [rdx]
	 vfmadd213ps ymm1, ymm2, [r8]
	 vfmadd213ps ymm1, ymm3, [rax]
	 vmovups     [r10], ymm1
	`

	testArguments(t, arguments8, StackArgs{Number: 2, OffsetToFirst: 16})

	arguments9 := `
	 mov r10, qword [rbp + 32]
	 mov r11, qword [rbp + 24]
	 mov rax, qword [rbp + 16]
	 vmovups     ymm0, [rdi]
	 vmovups     ymm1, [rsi]
	 vmovups     ymm2, [rcx]
	 vmovups     ymm3, [rax]
	 vfmadd213ps ymm1, ymm0, [rdx]
	 vfmadd213ps ymm1, ymm2, [r8]
	 vfmadd132ps ymm1, ymm3, [r9]
	 vfmadd213ps ymm1, ymm3, [r11]
	 vmovups     [r10], ymm1
	`

	testArguments(t, arguments9, StackArgs{Number: 3, OffsetToFirst: 16})

	arguments10 := `
	 mov r10, qword [rbp + 40]
	 mov r11, qword [rbp + 32]
	 mov rax, qword [rbp + 16]
	 mov rbx, qword [rbp + 24]
	 vmovups     ymm0, [rdi]
	 vmovups     ymm1, [rsi]
	 vmovups     ymm2, [rcx]
	 vmovups     ymm3, [r9]
	 vmovups     ymm4, [rbx]
	 vfmadd213ps ymm1, ymm0, [rdx]
	 vfmadd213ps ymm1, ymm2, [r8]
	 vfmadd213ps ymm1, ymm3, [rax]
	 vfmadd213ps ymm1, ymm4, [r11]
	 vmovups     [r10], ymm1
	`

	testArguments(t, arguments10, StackArgs{Number: 4, OffsetToFirst: 16})

	shiftbilinear := `
	 mov r15, rsi
	 mov rax, qword [rbp + 80]
	 mov r8, qword [rbp + 16]
	 mov r13, qword [rbp + 24]
	 mov rsi, qword [rbp + 32]
	 mov r12, qword [rbp + 40]
	 mov qword [rbp - 48], rdi
	 mov rdi, qword [rbp + 48]
	 mov qword [rbp - 56], rdx
	 mov rbx, qword [rbp + 56]
	 mov qword [rbp - 64], rcx
	 mov rcx, qword [rbp + 72]
	 mov qword [rbp - 72], rcx
	 mov rcx, qword [rbp + 64]
	 lea rdx, [rbp - 80]
	 mov qword [rsp + 80], rdx
	 lea rdx, [rbp - 76]
	 mov qword [rsp + 72], rdx
	 mov qword [rsp + 64], rax
	 lea rdx, [rbp - 72]
	`

	testArguments(t, shiftbilinear, StackArgs{Number: 9, OffsetToFirst: 16})

	//	shiftbilinear_rsp := `
	//	mov	rbp, rsp
	//	sub	rsp, 152
	//
	//	mov	qword ptr [rsp], r8
	//	mov	qword ptr [rsp + 8], r13
	//	mov	qword ptr [rsp + 16], rsi
	//	mov	qword ptr [rsp + 24], r12
	//	mov	qword ptr [rsp + 32], rdi
	//	mov	qword ptr [rsp + 40], rbx
	//	mov	qword ptr [rsp + 48], rcx
	//	mov	qword ptr [rsp + 56], rdx
	//	mov	qword ptr [rsp + 64], rax
	//	mov	qword ptr [rsp + 72], rdx
	//	mov	qword ptr [rsp + 80], rdx
	//
	//	mov	dword ptr [rsp], ecx
	//	mov	qword ptr [rsp + 8], r12
	//	mov	qword ptr [rsp + 16], rax
	//
	//	add	rsp, 152
	//	`
	//
	//	shiftlinear_rbp := `
	//	mov	qword ptr [rbp - 48], rdi
	//	mov	qword ptr [rbp - 56], rdx
	//	mov	qword ptr [rbp - 64], rcx
	//	mov	qword ptr [rbp - 72], rcx
	//	lea	rdx, [rbp - 80]
	//	lea	rdx, [rbp - 76]
	//	lea	rdx, [rbp - 72]
	//	lea	rdi, [rbp - 48]
	//	lea	rdx, [rbp - 56]
	//	lea	rcx, [rbp - 64]
	//	mov	rbx, qword ptr [rbp - 48]
	//	mov	rdx, qword ptr [rbp - 56]
	//	mov	r11, qword ptr [rbp - 64]
	//	mov	r9d, dword ptr [rbp - 76]
	//	mov	ecx, dword ptr [rbp - 80]
	//	mov	r12, qword ptr [rbp - 72]
	//	mov	qword ptr [rbp - 104], rax ## 8-byte Spill
	//	mov	qword ptr [rbp - 88], r11 ## 8-byte Spill
	//	mov	qword ptr [rbp - 96], rax ## 8-byte Spill
	//	cmp	qword ptr [rbp - 96], r14 ## 8-byte Folded Reload
	//	mov	qword ptr [rbp - 104], rax ## 8-byte Spill
	//	mov	qword ptr [rbp - 88], r11 ## 8-byte Spill
	//	cmp	qword ptr [rbp - 104], 7 ## 8-byte Folded Reload
	//	mov	rax, qword ptr [rbp - 88] ## 8-byte Reload
	//	mov	qword ptr [rbp - 88], rax ## 8-byte Spill
	//	mov	rax, qword ptr [rbp - 88] ## 8-byte Reload
	//	`
}

func testProto(t *testing.T, protoName, goline string, expectedNumArgs, expectedNumRets int, expectError bool) {

	_, args, rets, err := getGolangArgs(protoName, goline)

	if len(args) != expectedNumArgs {
		t.Errorf("testProto(): \nexpected number of arguments %d\ngot      %d", expectedNumArgs, len(args))
	}
	if expectError && err == nil {
		t.Errorf("testProto(): \nexpected error %v\ngot      %v", expectError, err)
	}
	if err != nil && len(rets) != expectedNumRets {
		t.Errorf("testProto(): \nexpected number of return values %d\ngot      %d", expectedNumRets, len(rets))
	}
}

func TestPrototypes(t *testing.T) {

	proto1 := `func _SimdAvx2BgraToGray(bgra unsafe.Pointer, width uint64, height uint64, bgraStride uint64, gray unsafe.Pointer, grayStride uint64)`

	testProto(t, "SimdAvx2BgraToGray", proto1, 6, 0, false)

	proto2 := `func _SimdSsse3Reorder32bit(src unsafe.Pointer, size uint64, dst unsafe.Pointer)`

	testProto(t, "SimdSsse3Reorder32bit", proto2, 3, 0, false)

	proto3 := `func _SimdAvx2ReduceGray2x2(src unsafe.Pointer, srcWidth, srcHeight, srcStride uint64, dst unsafe.Pointer, dstWidth, dstHeight, dstStride uint64)`

	testProto(t, "SimdAvx2ReduceGray2x2", proto3, 8, 0, false)

	proto4 := `func _SimdShiftBilinear(src unsafe.Pointer, srcStride, width, height, channelCount uin64, bkg unsafe.Pointer,  bkgStride uint64, shiftX, shiftY unsafe.Pointer, cropLeft, cropTop, cropRight, cropBottom uin64, dst unsafe.Pointer, dstStride uint64)`

	testProto(t, "SimdShiftBilinear", proto4, 15, 0, false)

	proto5 := `func _SimdSse2HistogramBufAllocSize(width int) (alloc int)`

	testProto(t, "SimdSse2HistogramBufAllocSize", proto5, 1, 1, false)

	proto6 := `func _SimdSse2HistogramBufAllocSize(width int) int`

	testProto(t, "SimdSse2HistogramBufAllocSize", proto6, 1, 0, true)

}
