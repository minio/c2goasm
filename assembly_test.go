package main

import (
	"strings"
	"testing"
)

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
		result := fixRbpPlusLoad(test, stackArgs, 0, false)

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
