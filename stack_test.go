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

func testStack(t *testing.T, postamble string, expected Stack) {
	preamble := ExtractStackInfo(strings.Split(postamble, "\n"))

	if preamble.StackSize != expected.StackSize || preamble.AlignedStack != expected.AlignedStack ||
	   preamble.VZeroUpper != expected.VZeroUpper || !equalString(preamble.Pushes, expected.Pushes) {
		t.Errorf("testStack(): \nexpected %s\ngot      %s", expected, preamble)
	}
}

func TestStacks(t *testing.T) {

	//   push    rbp
	//   mov     rbp, rsp
	preamble1 := Stack{SetRbpIns: true, VZeroUpper: true}
	preamble1.Pushes = append(preamble1.Pushes, "rbp")

	postamble1 := `	    pop     rbp
	    vzeroupper
	    ret`

	testStack(t, postamble1, preamble1)

	preamble2 := Stack{SetRbpIns: true, AlignedStack: true}
	preamble2.Pushes = append(preamble2.Pushes, "rbp", "r15", "r14", "r13", "r12", "rbx")

	//   push    rbp
	//   mov     rbp, rsp
	//   push    r15
	//   push    r14
	//   push    r13
	//   push    r12
	//   push    rbx
	//   and     rsp, -32
	//   sub     rsp, 864
	postamble2 := `        lea     rsp, [rbp - 40]
        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp`

	testStack(t, postamble2, preamble2)

	//   push    rbp
	//   mov     rbp, rsp
	//   push    r15
	//   push    r14
	//   push    r13
	//   push    r12
	//   push    rbx
	preamble3 := Stack{SetRbpIns: true}
	preamble3.Pushes = append(preamble3.Pushes, "rbp", "r15", "r14", "r13", "r12", "rbx")

	postamble3 := `        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp`

	testStack(t, postamble3, preamble3)

	//   push    rbp
	//   mov     rbp, rsp
	//   push    r15
	//   push    r14
	//   push    r13
	//   push    r12
	//   push    rbx
	//   sub     rsp, 152`
	preamble4 := Stack{SetRbpIns: true, StackSize: 152}
	preamble4.Pushes = append(preamble4.Pushes, "rbp", "r15", "r14", "r13", "r12", "rbx")

	postamble4 := `        add     rsp, 152
        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp`

	testStack(t, postamble4, preamble4)

	//   push    rbp
	//   mov     rbp, rsp
	//   push    r15
	//   push    r14
	//   push    r13
	//   push    r12
	//   push    rbx
	//   and     rsp, -32
	//   sub     rsp, 192`
	preamble5 := Stack{SetRbpIns: true, AlignedStack: true, VZeroUpper: true}
	preamble5.Pushes = append(preamble5.Pushes, "rbp", "r15", "r14", "r13", "r12", "rbx")

	postamble5 := `        lea     rsp, [rbp - 40]
        pop     rbx
        pop     r12
        pop     r13
        pop     r14
        pop     r15
        pop     rbp
        vzeroupper
        ret`

	testStack(t, postamble5, preamble5)
}
