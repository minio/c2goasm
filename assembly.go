package main

import (
	"fmt"
	"strings"
	"unicode"
)

var registers = [...]string{"DI", "SI", "DX", "CX", "R8", "R9"}

// Write the prologue for the subroutine
func WriteGoasmPrologue(segment Segment, number int, table Table) []string {

	var result []string

	// Output name of subroutine
	result = append(result, fmt.Sprintf("TEXT ·_%s(SB), 7, $0\n", segment.Name))

	arg, reg := 0, ""
	for arg, reg = range registers {

		result = append(result, fmt.Sprintf("    MOVQ arg%d+%d(FP), %s", arg+1, arg*8, reg))
		if arg+1 == number {
			break
		}
	}

	if table.IsPresent() {
		// Setup base pointer for loading constants
		result = append(result, "", fmt.Sprintf("    LEAQ LCD%s<>(SB), BP"), table.Data)
	}

	return result
}

func WriteGoasmBody(lines []string, table Table, stack Stack) ([]string, error) {

	var result []string

	eatHeader := true
	var goasm []string
	for _, line := range lines {

		// Remove ## comments
		if parts := strings.SplitN(line, `##`, 2); len(parts) > 1 {
			if strings.TrimSpace(parts[0]) == "" {
				continue
			}
			line = parts[0]
		}

		if eatHeader {

			if b, asm := stack.IsStdCallPrologue(line); b {
				fmt.Println("SKIPPING:", line)
				if asm != "" {
					goasm = append(goasm, asm)
				}
				continue
			} else {
				// Output equivalent asm instructions for golang
				result = append(result, goasm...)

				// And we are done with removing the std call header
				eatHeader = false
			}
		}

		// Skip lines with aligns
		if strings.Contains(line, ".align") {
			continue
		}

		// Make jmps uppercase
		if parts := strings.SplitN(line, `LBB`, 2); len(parts) > 1 {
			// unless it is a label
			if !strings.Contains(parts[1], ":") {
				// make jmp statement uppercase
				line = strings.ToUpper(parts[0]) + "LBB" + parts[1]
			}
		}

		fields := strings.Fields(line)
		// Test for any non-jmp instruction (lower case mnemonic)
		if len(fields) > 0 && !strings.Contains(fields[0], ":") && isLower(fields[0]) {
			// prepend line with comment for subsequent asm2plan9s assembly
			line = "                                 // " + strings.TrimSpace(line)
		}

		line = removeUndefined(line, "ptr")
		line = removeUndefined(line, "xmmword")
		line = removeUndefined(line, "ymmword")

		line = fixShifts(line)
		line = fixPicLabels(line, table)

		result = append(result, line)
	}

	return result, nil
}

// Write the epilogue for the subroutine
func (s *Stack) WriteGoasmEpilogue() []string {

	var result []string

	// Restore the stack pointer
	// - for an aligned stack, restore the stack pointer from the stack itself
	// - for an unaligned stack, simply add the (fixed size) stack size in order restore the stack pointer
	if s.AlignedStack {
		panic("TODO: Restore stack pointer from stack")
	} else {
		if s.StackSize != 0 {
			result = append(result, fmt.Sprintf("    ADD $%d, SP", s.StackSize))
		}
	}

	// Clear upper half of YMM register, if so done in the original code
	if s.VZeroUpper {
		result = append(result, "    VZEROUPPER")
	}

	// Finally, return out of the subroutine
	result = append(result, "    RET")

	return result
}


func isLower(str string) bool {

	for _, r := range str {
		return unicode.IsLower(r)
	}
	return false
}

func removeUndefined(line, undef string) string {

	if parts := strings.SplitN(line, undef, 2); len(parts) > 1 {
		line = parts[0] + strings.TrimSpace(parts[1])
	}
	return line
}

// fix Position Independent Labels
func fixPicLabels(line string, table Table) string {

	if strings.Contains(line, "[rip + ") {
		parts := strings.SplitN(line, "[rip + ", 2)
		label := parts[1][:len(parts[1])-1]

		i := -1
		var l Label
		for i, l = range table.Labels {
			if l.Name == label {
				line = parts[0] + fmt.Sprintf("%d[rbp] /* [rip + %s */", l.Offset, parts[1])
				break
			}
		}
		if i == len(table.Labels) {
			panic(fmt.Sprintf("Failed to find label to replace of position independent code: %s", label))
		}
	}

	return line
}

func fixShiftNoArgument(line, ins string) string {

	if strings.Contains(line, ins) {
		parts := strings.SplitN(line, ins, 2)
		args := strings.SplitN(parts[1], ",", 2)
		if len(args) == 1 {
			line += ", 1"
		}
	}

	return line
}

func fixShifts(line string) string {

	line = fixShiftNoArgument(line, "shr")
	line = fixShiftNoArgument(line, "sar")

	return line
}
