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

func testStack(t *testing.T, epilogue string, expected Stack) {
	stack := ExtractStackInfo(strings.Split(epilogue, "\n"))

	if stack.StackSize != expected.StackSize || stack.AlignedStack != expected.AlignedStack ||
	   stack.VZeroUpper != expected.VZeroUpper || !equalString(stack.Pushes, expected.Pushes) {
		t.Errorf("testStack(): \nexpected %s\ngot      %s", expected, stack)
	}
}

func TestStacks(t *testing.T) {

	//   push    rbp
	//   mov     rbp, rsp
	epilogue1 := Stack{SetRbpIns: true, VZeroUpper: true}
	epilogue1.Pushes = append(epilogue1.Pushes, "rbp")

	prologue1 := `	    pop     rbp
	    vzeroupper
	    ret`

	testStack(t, prologue1, epilogue1)

	/***********************************************************************************/

	epilogue2 := Stack{SetRbpIns: true, AlignedStack: true}
	epilogue2.Pushes = append(epilogue2.Pushes, "rbp", "r15", "r14", "r13", "r12", "rbx")

	//   push    rbp
	//   mov     rbp, rsp
	//   push    r15
	//   push    r14
	//   push    r13
	//   push    r12
	//   push    rbx
	//   and     rsp, -32
	//   sub     rsp, 864
	prologue2 := `        lea     rsp, [rbp - 40]
        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp`

	testStack(t, prologue2, epilogue2)

	/***********************************************************************************/

	//   push    rbp
	//   mov     rbp, rsp
	//   push    r15
	//   push    r14
	//   push    r13
	//   push    r12
	//   push    rbx
	epilogue3 := Stack{SetRbpIns: true}
	epilogue3.Pushes = append(epilogue3.Pushes, "rbp", "r15", "r14", "r13", "r12", "rbx")

	prologue3 := `        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp`

	testStack(t, prologue3, epilogue3)

	/***********************************************************************************/

	//   push    rbp
	//   mov     rbp, rsp
	//   push    r15
	//   push    r14
	//   push    r13
	//   push    r12
	//   push    rbx
	//   sub     rsp, 152`
	prologue4 := Stack{SetRbpIns: true, StackSize: 152}
	prologue4.Pushes = append(prologue4.Pushes, "rbp", "r15", "r14", "r13", "r12", "rbx")

	epilogue4 := `        add     rsp, 152
        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp`

	testStack(t, epilogue4, prologue4)

	/***********************************************************************************/

	//   push    rbp
	//   mov     rbp, rsp
	//   push    r15
	//   push    r14
	//   push    r13
	//   push    r12
	//   push    rbx
	//   and     rsp, -32
	//   sub     rsp, 192`
	prologue5 := Stack{SetRbpIns: true, AlignedStack: true, VZeroUpper: true}
	prologue5.Pushes = append(prologue5.Pushes, "rbp", "r15", "r14", "r13", "r12", "rbx")

	epilogue5 := `        lea     rsp, [rbp - 40]
        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp
        vzeroupper
        ret`

	testStack(t, epilogue5, prologue5)
}
