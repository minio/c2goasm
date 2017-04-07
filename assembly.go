package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const returnAddrOnStack = 8

var registers = [...]string{"DI", "SI", "DX", "CX", "R8", "R9"}
var regexpCall = regexp.MustCompile(`^\s*call\s*`)
var regexpLabel = regexp.MustCompile(`^(\.?LBB.*:)`)
var regexpJumpTableRef = regexp.MustCompile(`\[rip \+ (\.?LJTI[_0-9]*)\]\s*$`)
var regexpJumpWithLabel = regexp.MustCompile(`^(\s*j\w*)\s*(\.?LBB.*)`)
var regexpRbpLoadHigher = regexp.MustCompile(`\[rbp \+ ([0-9]+)\]\s*$`)
var regexpRbpLoadLower = regexp.MustCompile(`\[rbp - ([0-9]+)\]`)
var regexpStripComments = regexp.MustCompile(`\s*#?#\s.*$`)

// Write the prologue for the subroutine
func writeGoasmPrologue(subroutine Subroutine, arguments int, table Table) []string {

	var result []string

	// Output definition of subroutine
	result = append(result, fmt.Sprintf("TEXT ·_%s(SB), 7, $%d-%d\n", subroutine.name,
		subroutine.epilogue.getTotalStackDepth(table, arguments), getTotalSizeOfArguments(0, arguments-1)))

	if subroutine.epilogue.AlignedStack {
		// Save original stack pointer right below newly aligned stack pointer
		result = append(result, fmt.Sprintf("    MOVQ SP, BP"))
		result = append(result, fmt.Sprintf("    ANDQ $-%d, BP", subroutine.epilogue.AlignValue))
		result = append(result, fmt.Sprintf("    MOVQ SP, -%d(BP)", returnAddrOnStack)) // Save original SP

		// In case base pointer is used for constants and the offset is non-deterministic
		if table.isPresent() {
			// For the case assembly expects stack based arguments
			for arg := arguments - 1; arg >= len(registers); arg-- {
				// Copy golang stack based arguments below saved original stack pointer
				result = append(result, fmt.Sprintf("    MOVQ arg%d+%d(FP), DI", arg+1, arg*8))
				result = append(result, fmt.Sprintf("    MOVQ DI, %d(BP)", -returnAddrOnStack+(arguments-arg)*-8))
			}
		}
	}

	// Load initial arguments (up to 6) in corresponding registers
	for arg, reg := range registers {

		result = append(result, fmt.Sprintf("    MOVQ arg%d+%d(FP), %s", arg+1, arg*8, reg))
		if arg+1 == arguments {
			break
		}
	}

	if table.isPresent() {
		// Setup base pointer for loading constants
		result = append(result, fmt.Sprintf("    LEAQ %s<>(SB), BP", table.Name))
	} else if subroutine.epilogue.AlignedStack {
		// Setup base pointer to be able to load golang stack based arguments
		result = append(result, fmt.Sprintf("    MOVQ SP, BP"))
	}

	// Setup the stack pointer
	if subroutine.epilogue.AlignedStack {
		// Aligned stack as required (zeroing out lower order bits), and create space
		result = append(result, fmt.Sprintf("    ANDQ $-%d, SP", subroutine.epilogue.AlignValue))
	}
	if subroutine.epilogue.getStackpointerDecrement(table, arguments) != 0 {
		// Create stack space as needed
		result = append(result, fmt.Sprintf("    SUBQ $%d, SP", subroutine.epilogue.getStackpointerDecrement(table, arguments)))
	}

	return append(result, ``)
}

func writeGoasmBody(lines []string, table Table, stackArgs StackArgs, epilogue Epilogue, arguments int) ([]string, error) {

	var result []string

	for iline, line := range lines {

		// If part of epilogue
		if iline >= epilogue.Start && iline < epilogue.End {

			// Instead of last line, output go assembly epilogue
			if iline == epilogue.End - 1 {
				result = append(result, writeGoasmEpilogue(epilogue, arguments, table)...)
			}
			continue
		}

		// Remove ## comments
		var skipLine bool
		line, skipLine = stripComments(line)
		if skipLine {
			continue
		}

		// Skip lines with aligns
		if strings.Contains(line, ".align") || strings.Contains(line, ".p2align") {
			continue
		}

		line, _ = fixLabels(line)
		line, _, _ = upperCaseJumps(line)
		line = upperCaseCalls(line)

		fields := strings.Fields(line)
		// Test for any non-jmp instruction (lower case mnemonic)
		if len(fields) > 0 && !strings.Contains(fields[0], ":") && isLower(fields[0]) {
			// prepend line with comment for subsequent asm2plan9s assembly
			line = "                                 // " + strings.TrimSpace(line)
		}

		line = removeUndefined(line, "ptr")
		line = removeUndefined(line, "# NOREX")

		// https://github.com/vertis/objconv/blob/master/src/disasm2.cpp
		line = replaceUndefined(line, "xmmword", "oword")
		line = replaceUndefined(line, "ymmword", "yword")

		line = fixShiftInstructions(line)
		line = fixMovabsInstructions(line)
		if table.isPresent() {
			line = fixPicLabels(line, table)
		}

		line = fixRbpPlusLoad(line, stackArgs, epilogue.getStackpointerDecrement(table, arguments)-epilogue.additionalStackSpace(table, arguments), table.isPresent(), epilogue.AlignedStack)

		detectRbpMinusMemoryAccess(line)
		detectJumpTable(line)

		result = append(result, line)
	}

	return result, nil
}

// Write the epilogue for the subroutine
func writeGoasmEpilogue(epilogue Epilogue, arguments int, table Table) []string {

	var result []string

	// Restore the stack pointer
	if epilogue.AlignedStack {
		// For an aligned stack, restore the stack pointer from the stack itself
		result = append(result, fmt.Sprintf("    MOVQ %d(SP), SP", epilogue.getStackpointerDecrement(table, arguments)-returnAddrOnStack))
	} else if epilogue.getStackpointerDecrement(table, arguments) != 0 {
		// For an unaligned stack, simply add the (fixed size) stack size in order restore the stack pointer
		result = append(result, fmt.Sprintf("    ADDQ $%d, SP", epilogue.getStackpointerDecrement(table, arguments)))
	}

	// Clear upper half of YMM register, if so done in the original code
	if epilogue.VZeroUpper {
		result = append(result, "    VZEROUPPER")
	}

	// Finally, return out of the subroutine
	result = append(result, "    RET")

	return result
}

// Strip comments from assembly lines
func stripComments(line string) (result string, skipLine bool) {

	if match := regexpStripComments.FindStringSubmatch(line); len(match) > 0 {
		line = line[:len(line)-len(match[0])]
		if line == "" {
			return "", true
		}
	}
	return line, false
}

// Remove leading `.` from labels
func fixLabels(line string) (string, string) {

	label := ""

	if match := regexpLabel.FindStringSubmatch(line); len(match) > 0 {
		label = strings.Replace(match[1], ".", "", 1)
		line = label
		label = strings.Replace(label, ":", "", 1)
	}

	return line, label
}

// Make jmps uppercase
func upperCaseJumps(line string) (string, string, string) {

	instruction, label := "", ""

	if match := regexpJumpWithLabel.FindStringSubmatch(line); len(match) > 1 {
		// make jmp statement uppercase
		instruction = strings.ToUpper(match[1])
		label = strings.Replace(match[2], ".", "", 1)
		line = instruction + " " + label

	}

	return line, strings.TrimSpace(instruction), label
}

// Make calls uppercase
func upperCaseCalls(line string) string {

	// Make 'call' instructions uppercase
	if match := regexpCall.FindStringSubmatch(line); len(match) > 0 {
		parts := strings.SplitN(line, `call`, 2)
		fname := strings.TrimSpace(parts[1])

		// replace c stdlib functions with equivalents
		if fname == "_memcpy" || fname == "memcpy@PLT" { // (Procedure Linkage Table)
			parts[1] = "clib·_memcpy(SB)"
		} else if fname == "_memset" || fname == "memset@PLT" { // (Procedure Linkage Table)
			parts[1] = "clib·_memset(SB)"
		} else if fname == "_floor" || fname == "floor@PLT" { // (Procedure Linkage Table)
			parts[1] = "clib·_floor(SB)"
		} else if fname == "___bzero" {
			parts[1] = "clib·_bzero(SB)"
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

func replaceUndefined(line, undef, repl string) string {

	if parts := strings.SplitN(line, undef, 2); len(parts) > 1 {
		line = parts[0] + repl + parts[1]
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
func fixRbpPlusLoad(line string, stackArgs StackArgs, stackSize uint, tableIsPresent, alignedStack bool) string {

	if match := regexpRbpLoadHigher.FindStringSubmatch(line); len(match) > 1 {
		offset, _ := strconv.Atoi(match[1])
		parts := strings.SplitN(line, "[rbp + ", 2)
		if tableIsPresent {
			// Base pointer is setup for loading constants, so cannot use
			if alignedStack {
				offset = int(stackSize) + (offset - stackArgs.OffsetToFirst)
				line = parts[0] + fmt.Sprintf("%d[rsp] /* [rbp + %s */", offset, parts[1])
			} else {
				// fixed stack size, load from stack pointer
				offset += int(stackSize)
				line = parts[0] + fmt.Sprintf("%d[rsp] /* [rbp + %s */", offset, parts[1])
			}
		} else {
			if alignedStack {
				// Base pointer equal to (initial) stack pointer, so can leave loads untouched
			} else {
				// fixed stack size, load from stack pointer
				offset = int(stackSize) + /*returnAddrOnStack*/ 8 + 8*len(registers) + (offset - stackArgs.OffsetToFirst)
				line = parts[0] + fmt.Sprintf("%d[rsp] /* [rbp + %s */", offset, parts[1])
			}
		}
	}

	return line
}

// Detect memory accesses in the form of '[rbp - constant]'
func detectRbpMinusMemoryAccess(line string) {

	if match := regexpRbpLoadLower.FindStringSubmatch(line); len(match) > 1 {

		panic(fmt.Sprintf("Not expected to find [rbp -] based loads: %s\n\nDid you specify `-mno-red-zone`?\n\n", line))
	}
}

// Detect jump tables
func detectJumpTable(line string) {

	if match := regexpJumpTableRef.FindStringSubmatch(line); len(match) > 0 {
		panic(fmt.Sprintf("Jump table detected: %s\n\nCircumvent using '-fno-jump-tables', see 'clang -cc1 -help' (version 3.9+)\n\n", match[1]))
	}
}
