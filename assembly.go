package main

import (
	"fmt"
	"strings"
	"unicode"
	"regexp"
	"strconv"
)

const returnAddrOnStack = 8
var registers = [...]string{"DI", "SI", "DX", "CX", "R8", "R9"}
var regexpCall = regexp.MustCompile(`^\s*call\s*`)

// Write the prologue for the subroutine
func WriteGoasmPrologue(segment Segment, arguments int, table Table) []string {

	var result []string

	// Output definition of subroutine
	result = append(result, fmt.Sprintf("TEXT ·_%s(SB), 7, $0\n", segment.Name))

	arg, reg := 0, ""

	// For the case assembly expects stack based arguments
	for arg = len(registers); arg < arguments; arg++ {
		// In case base pointer is used for constants and the offset is non-deterministic
		if table.IsPresent() && segment.stack.AlignedStack {
			// Copy golang arguments to C-style stack
			result = append(result, fmt.Sprintf("    MOVQ arg%d+%d(FP), DI", arg+1, arg*8))
			result = append(result, fmt.Sprintf("    MOVQ DI, %d(SP)", -256+(arg-6)*8))
		} else {
			// We can load the arguments from golang stack, so no need to do anything
		}
	}

	// Load initial arguments (up to 6) in corresponding registers
	for arg, reg = range registers {

		result = append(result, fmt.Sprintf("    MOVQ arg%d+%d(FP), %s", arg+1, arg*8, reg))
		if arg+1 == arguments {
			break
		}
	}

	// Setup the stack pointer
	if segment.stack.AlignedStack {
		// Aligned stack as required (zeroing out lower order bits), create space, and save original stack pointer
		result = append(result, fmt.Sprintf("    MOVQ SP, BP"))
		result = append(result, fmt.Sprintf("    ANDQ $%d, SP", 32))
		result = append(result, fmt.Sprintf("    SUBQ $%d, SP", segment.stack.StackSize))
		result = append(result, fmt.Sprintf("    MOVQ BP, -8(SP)"))
	} else if segment.stack.StackSize != 0 {
		// Unaligned stack, simply create space as required
		result = append(result, fmt.Sprintf("    SUBQ $%d, SP", segment.stack.StackSize))
	}

	if table.IsPresent() {
		// Setup base pointer for loading constants
		result = append(result, "", fmt.Sprintf("    LEAQ %s<>(SB), BP", table.Name), "")
	} else if segment.stack.AlignedStack {
		// Setup base pointer to be able to load stack based arguments
		result = append(result, "", fmt.Sprintf("    MOVQ SP, BP"), "")
	}

	return result
}

func WriteGoasmBody(lines []string, table Table, stackArgs StackArgs) ([]string, error) {

	var result []string

	for _, line := range lines {

		// Remove ## comments
		if parts := strings.SplitN(line, `##`, 2); len(parts) > 1 {
			if strings.TrimSpace(parts[0]) == "" {
				continue
			}
			line = parts[0]
		}

		// Skip lines with aligns
		if strings.Contains(line, ".align") {
			continue
		}

		line = upperCaseJumps(line)
		line = upperCaseCalls(line)

		fields := strings.Fields(line)
		// Test for any non-jmp instruction (lower case mnemonic)
		if len(fields) > 0 && !strings.Contains(fields[0], ":") && isLower(fields[0]) {
			// prepend line with comment for subsequent asm2plan9s assembly
			line = "                                 // " + strings.TrimSpace(line)
		}

		line = removeUndefined(line, "ptr")
		line = removeUndefined(line, "xmmword")
		line = removeUndefined(line, "ymmword")
		line = removeUndefined(line, "# NOREX")

		line = fixShiftInstructions(line)
		line = fixMovabsInstructions(line)
		if table.IsPresent() {
			line = fixPicLabels(line, table)
		}

		line = fixRbpPlusLoad(line, stackArgs)
		line = fixRbpMinusMemoryAccess(line)

		result = append(result, line)
	}

	return result, nil
}

// Write the epilogue for the subroutine
func WriteGoasmEpilogue(stack Stack) []string {

	var result []string

	// Restore the stack pointer
	// - for an aligned stack, restore the stack pointer from the stack itself
	// - for an unaligned stack, simply add the (fixed size) stack size in order restore the stack pointer
	if stack.AlignedStack {
		result = append(result, fmt.Sprintf("    MOVQ -8(SP), SP"))
	} else {
		if stack.StackSize != 0 {
			result = append(result, fmt.Sprintf("    ADDQ $%d, SP", stack.StackSize))
		}
	}

	// Clear upper half of YMM register, if so done in the original code
	if stack.VZeroUpper {
		result = append(result, "    VZEROUPPER")
	}

	// Finally, return out of the subroutine
	result = append(result, "    RET")

	return result
}

// Make jmps uppercase
func upperCaseJumps(line string) string {

	if parts := strings.SplitN(line, `LBB`, 2); len(parts) > 1 {
		// unless it is a label
		if !strings.Contains(parts[1], ":") {
			// make jmp statement uppercase
			line = strings.ToUpper(parts[0]) + "LBB" + parts[1]
		}
	}

	return line
}

// Make calls uppercase
func upperCaseCalls(line string) string {

	// Make 'call' instructions uppercase
	if match := regexpCall.FindStringSubmatch(line); len(match) > 0 {
		parts := strings.SplitN(line, `call`, 2)

		// replace c stdlib functions with equivalents
		if strings.TrimSpace(parts[1]) == "_memcpy" {
			parts[1] = fmt.Sprintf("clib·%s(SB)", strings.TrimSpace(parts[1]))
		}
		line = parts[0] + "CALL " + strings.TrimSpace(parts[1])
	}

	return line
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

func fixShiftInstructions(line string) string {

	line = fixShiftNoArgument(line, "shr")
	line = fixShiftNoArgument(line, "sar")

	return line
}

func fixMovabsInstructions(line string) string {

	if strings.Contains(line, "movabs") {
		parts := strings.SplitN(line, "movabs", 2)
		line = parts[0] + "mov" + parts[1]
	}

	return line
}

// Fix loads in the form of '[rbp + constant]'
// These are load instructions for stack-based arguments that occur after the first 6 arguments
// Remap to rsp/stack pointer and load from golang stack
func fixRbpPlusLoad(line string, stackArgs StackArgs) string {

	regexpRbpLoadHigher := regexp.MustCompile(`\[rbp \+ ([0-9]+)\]$`)

	if match := regexpRbpLoadHigher.FindStringSubmatch(line); len(match) > 1 {
		offset, _ := strconv.Atoi(match[1])
		offset = offset - stackArgs.OffsetToFirst + returnAddrOnStack + 8*len(registers)
		parts := strings.SplitN(line, "[rbp + ", 2)
		line = parts[0] + fmt.Sprintf("%d[rsp] /* [rbp + %s */", offset, parts[1])
	}

	return line
}

// Fix memory accesses in the form of '[rbp - constant]'
func fixRbpMinusMemoryAccess(line string) string {

	return line
}
