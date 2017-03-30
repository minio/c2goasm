package main

import (
	"strings"
	"testing"
)

func equalString(a, b []string) bool {

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func testEpilogue(t *testing.T, prologue, epilogue string, expected Epilogue) {
	stack := extractEpilogueInfo(strings.Split(epilogue, "\n"))

	for _, line := range strings.Split(prologue, "\n") {
		isPrologue := stack.IsPrologueInstruction(line)
		if !isPrologue {
			t.Errorf("testStack(): \nexpected true\ngot      %#v", isPrologue)
		}
	}

	if stack.StackSize != expected.StackSize || stack.AlignedStack != expected.AlignedStack ||
		stack.VZeroUpper != expected.VZeroUpper || !equalString(stack.Pops, expected.Pops) ||
		stack.SetRbpIns != expected.SetRbpIns {
		t.Errorf("testStack(): \nexpected %#v\ngot      %#v", expected, stack)
	}
}

func TestEpilogues(t *testing.T) {

	asmPrologue1 := `	   push    rbp
	   mov     rbp, rsp`

	epilogue1 := Epilogue{SetRbpIns: true, VZeroUpper: true}
	epilogue1.Pops = append(epilogue1.Pops, "rbp")

	asmEpilogue1 := `	    pop     rbp
	    vzeroupper
	    ret`

	testEpilogue(t, asmPrologue1, asmEpilogue1, epilogue1)

	/***********************************************************************************/

	asmPrologue2 := `	   push    rbp
	   mov     rbp, rsp
	   push    r15
	   push    r14
	   push    r13
	   push    r12
	   push    rbx
	   and     rsp, -32
	   sub     rsp, 864`

	epilogue2 := Epilogue{SetRbpIns: true, StackSize: 864, AlignedStack: true, AlignValue: -32}
	epilogue2.Pops = append(epilogue2.Pops, "rbp", "r15", "r14", "r13", "r12", "rbx")

	asmEpilogue2 := `        lea     rsp, [rbp - 40]
        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp`

	testEpilogue(t, asmPrologue2, asmEpilogue2, epilogue2)

	/***********************************************************************************/

	asmPrologue3 := `push    rbp
	   mov     rbp, rsp
	   push    r15
	   push    r14
	   push    r13
	   push    r12
	   push    rbx`

	epilogue3 := Epilogue{SetRbpIns: true}
	epilogue3.Pops = append(epilogue3.Pops, "rbp", "r15", "r14", "r13", "r12", "rbx")

	asmEpilogue3 := `        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp`

	testEpilogue(t, asmPrologue3, asmEpilogue3, epilogue3)

	/***********************************************************************************/

	asmPrologue4 := `push    rbp
	   mov     rbp, rsp
	   push    r15
	   push    r14
	   push    r13
	   push    r12
	   push    rbx
	   sub     rsp, 152`

	epilogue4 := Epilogue{SetRbpIns: true, StackSize: 152}
	epilogue4.Pops = append(epilogue4.Pops, "rbp", "r15", "r14", "r13", "r12", "rbx")

	asmEpilogue4 := `        add     rsp, 152
        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp`

	testEpilogue(t, asmPrologue4, asmEpilogue4, epilogue4)

	/***********************************************************************************/

	asmPrologue5 := `push    rbp
	   mov     rbp, rsp
	   push    r15
	   push    r14
	   push    r13
	   push    r12
	   push    rbx
	   and     rsp, -32
	   sub     rsp, 192`

	epilogue5 := Epilogue{SetRbpIns: true, StackSize: 192, AlignedStack: true, AlignValue: -32, VZeroUpper: true}
	epilogue5.Pops = append(epilogue5.Pops, "rbp", "r15", "r14", "r13", "r12", "rbx")

	asmEpilogue5 := `        lea     rsp, [rbp - 40]
        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp
        vzeroupper
        ret`

	testEpilogue(t, asmPrologue5, asmEpilogue5, epilogue5)

	/***********************************************************************************/

	asmPrologue6 := `push	rbp
	mov	rbp, rsp
	push	r15
	push	r14
	push	r13
	push	r12
	push	rbx
	push	rax`

	epilogue6 := Epilogue{SetRbpIns: true, VZeroUpper: true}
	epilogue6.Pops = append(epilogue6.Pops, "rbp", "r15", "r14", "r13", "r12", "rbx")

	// `add rsp, 8` counters the additional `push rax` (there are 7 pushes and 6 pops)
	asmEpilogue6 := `	add	rsp, 8
	pop	rbx
	pop	r12
	pop	r13
	pop	r14
	pop	r15
	pop	rbp
	vzeroupper
	ret`

	testEpilogue(t, asmPrologue6, asmEpilogue6, epilogue6)

}
