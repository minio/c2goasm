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

func testStack(t *testing.T, prologue, epilogue string, expected Stack) {
	stack := ExtractStackInfo(strings.Split(epilogue, "\n"))

	if stack.StackSize != expected.StackSize || stack.AlignedStack != expected.AlignedStack ||
	   stack.VZeroUpper != expected.VZeroUpper || !equalString(stack.Pushes, expected.Pushes) ||
	   stack.SetRbpIns != expected.SetRbpIns {
		t.Errorf("testStack(): \nexpected %s\ngot      %s", expected, stack)
	}

	for _, line := range strings.Split(prologue, "\n") {
		isPrologue := stack.IsPrologueInstruction(line)
		if !isPrologue {
			t.Errorf("testStack(): \nexpected %s\ngot      %s", true, isPrologue)
		}
	}
}

func TestStacks(t *testing.T) {

	prologue1 := `	   push    rbp
	   mov     rbp, rsp`

	stack1 := Stack{SetRbpIns: true, VZeroUpper: true}
	stack1.Pushes = append(stack1.Pushes, "rbp")

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

	stack2 := Stack{SetRbpIns: true, AlignedStack: true}
	stack2.Pushes = append(stack2.Pushes, "rbp", "r15", "r14", "r13", "r12", "rbx")

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

	stack3 := Stack{SetRbpIns: true}
	stack3.Pushes = append(stack3.Pushes, "rbp", "r15", "r14", "r13", "r12", "rbx")

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

	stack4 := Stack{SetRbpIns: true, StackSize: 152}
	stack4.Pushes = append(stack4.Pushes, "rbp", "r15", "r14", "r13", "r12", "rbx")

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

	stack5 := Stack{SetRbpIns: true, AlignedStack: true, VZeroUpper: true}
	stack5.Pushes = append(stack5.Pushes, "rbp", "r15", "r14", "r13", "r12", "rbx")

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
}
