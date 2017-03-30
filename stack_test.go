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

func testStack(t *testing.T, prologue, epilogue string, expected Epilogue) {
	stack := ExtractEpilogueInfo(strings.Split(epilogue, "\n"))

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

func TestStacks(t *testing.T) {

	prologue1 := `	   push    rbp
	   mov     rbp, rsp`

	stack1 := Epilogue{SetRbpIns: true, VZeroUpper: true}
	stack1.Pops = append(stack1.Pops, "rbp")

	epilogue1 := `	    pop     rbp
	    vzeroupper
	    ret`

	testStack(t, prologue1, epilogue1, stack1)

	/***********************************************************************************/

	prologue2 := `	   push    rbp
	   mov     rbp, rsp
	   push    r15
	   push    r14
	   push    r13
	   push    r12
	   push    rbx
	   and     rsp, -32
	   sub     rsp, 864`

	stack2 := Epilogue{SetRbpIns: true, StackSize: 864, AlignedStack: true, AlignValue: -32}
	stack2.Pops = append(stack2.Pops, "rbp", "r15", "r14", "r13", "r12", "rbx")

	epilogue2 := `        lea     rsp, [rbp - 40]
        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp`

	testStack(t, prologue2, epilogue2, stack2)

	/***********************************************************************************/

	prologue3 := `push    rbp
	   mov     rbp, rsp
	   push    r15
	   push    r14
	   push    r13
	   push    r12
	   push    rbx`

	stack3 := Epilogue{SetRbpIns: true}
	stack3.Pops = append(stack3.Pops, "rbp", "r15", "r14", "r13", "r12", "rbx")

	epilogue3 := `        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp`

	testStack(t, prologue3, epilogue3, stack3)

	/***********************************************************************************/

	prologue4 := `push    rbp
	   mov     rbp, rsp
	   push    r15
	   push    r14
	   push    r13
	   push    r12
	   push    rbx
	   sub     rsp, 152`

	stack4 := Epilogue{SetRbpIns: true, StackSize: 152}
	stack4.Pops = append(stack4.Pops, "rbp", "r15", "r14", "r13", "r12", "rbx")

	epilogue4 := `        add     rsp, 152
        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp`

	testStack(t, prologue4, epilogue4, stack4)

	/***********************************************************************************/

	prologue5 := `push    rbp
	   mov     rbp, rsp
	   push    r15
	   push    r14
	   push    r13
	   push    r12
	   push    rbx
	   and     rsp, -32
	   sub     rsp, 192`

	stack5 := Epilogue{SetRbpIns: true, StackSize: 192, AlignedStack: true, AlignValue: -32, VZeroUpper: true}
	stack5.Pops = append(stack5.Pops, "rbp", "r15", "r14", "r13", "r12", "rbx")

	epilogue5 := `        lea     rsp, [rbp - 40]
        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp
        vzeroupper
        ret`

	testStack(t, prologue5, epilogue5, stack5)

	/***********************************************************************************/

	prologue6 := `push	rbp
	mov	rbp, rsp
	push	r15
	push	r14
	push	r13
	push	r12
	push	rbx
	push	rax`

	stack6 := Epilogue{SetRbpIns: true, VZeroUpper: true}
	stack6.Pops = append(stack6.Pops, "rbp", "r15", "r14", "r13", "r12", "rbx")

	// `add rsp, 8` counters the additional `push rax` (there are 7 pushes and 6 pops)
	epilogue6 := `	add	rsp, 8
	pop	rbx
	pop	r12
	pop	r13
	pop	r14
	pop	r15
	pop	rbp
	vzeroupper
	ret`

	testStack(t, prologue6, epilogue6, stack6)

}
