package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const originalStackPointer = 8
const returnAddrOnStack = 8

var registers = [...]string{"DI", "SI", "DX", "CX", "R8", "R9"}
var regexpCall = regexp.MustCompile(`^\s*call\s*`)
var regexpPushInstr = regexp.MustCompile(`^\s*push\s*`)
var regexpPopInstr = regexp.MustCompile(`^\s*pop\s*`)
var regexpLabel = regexp.MustCompile(`^(\.?LBB.*:)`)
var regexpJumpTableRef = regexp.MustCompile(`\[rip \+ (\.?LJTI[_0-9]*)\]\s*$`)
var regexpJumpWithLabel = regexp.MustCompile(`^(\s*j\w*)\s*(\.?LBB.*)`)
var regexpRbpLoadHigher = regexp.MustCompile(`\[rbp \+ ([0-9]+)\]`)
var regexpRbpLoadLower = regexp.MustCompile(`\[rbp - ([0-9]+)\]`)
var regexpStripComments = regexp.MustCompile(`\s*#?#\s.*$`)

// Write the prologue for the subroutine
func writeGoasmPrologue(sub Subroutine, arguments, returnValues []string) []string {

	var result []string

	// Output definition of sub
	result = append(result, fmt.Sprintf("TEXT ·_%s(SB), 7, $%d-%d", sub.name,
		sub.epilogue.getTotalStackDepth(sub.table, len(arguments)),
		getTotalSizeOfArgumentsAndReturnValues(0, len(arguments)-1, returnValues)), "")

	if sub.epilogue.AlignedStack {

		offset := sub.epilogue.getTotalStackDepth(sub.table, len(arguments))
		if offset % sub.epilogue.AlignValue != 0 {
			panic(fmt.Sprintf("Offset (%d) must be a multiple of alignment value (%d)", offset,
				sub.epilogue.AlignValue))
		}

		// We can save one addq instruction by collapsing the offset into the 'offset(BP)'
		// result = append(result, fmt.Sprintf("    ADDQ $%d, BP", offset))
		destAddr := offset - originalStackPointer

		// Save original stack pointer right below newly aligned stack pointer
		result = append(result, fmt.Sprintf("    MOVQ SP, BP"))
		result = append(result, fmt.Sprintf("    ANDQ $-%d, BP", sub.epilogue.AlignValue))
		result = append(result, fmt.Sprintf("    MOVQ SP, %d(BP)", destAddr)) // Save original SP

		// In case base pointer is used for constants and the offset is non-deterministic
		if sub.table.isPresent() {
			// For the case assembly expects stack based arguments
			for arg := len(arguments) - 1; arg >= len(registers); arg-- {
				destAddr -= 8
				// Copy golang stack based arguments below saved original stack pointer
				result = append(result, fmt.Sprintf("    MOVQ %s+%d(FP), DI", arguments[arg], arg*8))
				result = append(result, fmt.Sprintf("    MOVQ DI, %d(BP)", destAddr))
			}
		}
	}

	// Load initial arguments (up to 6) in corresponding registers
	for arg, reg := range registers {

		result = append(result, fmt.Sprintf("    MOVQ %s+%d(FP), %s", arguments[arg], arg*8, reg))
		if arg+1 == len(arguments) {
			break
		}
	}

	if sub.table.isPresent() {
		// Setup base pointer for loading constants
		result = append(result, fmt.Sprintf("    LEAQ %s<>(SB), BP", sub.table.Name))
	} else if sub.epilogue.AlignedStack {
		// Setup base pointer to be able to load golang stack based arguments
		result = append(result, fmt.Sprintf("    MOVQ SP, BP"))
	}

	// Setup the stack pointer
	if sub.epilogue.AlignedStack {
		// Align stack pointer to next multiple of alignment space
		result = append(result, fmt.Sprintf("    ADDQ $%d, SP", sub.epilogue.AlignValue))
		result = append(result, fmt.Sprintf("    ANDQ $-%d, SP", sub.epilogue.AlignValue))
	} else if sub.epilogue.getStackpointerDecrement(sub.table, len(arguments)) != 0 {
		// Create stack space as needed
		result = append(result, fmt.Sprintf("    ADDQ $%d, SP", sub.epilogue.getFreeSpaceAtBottom()))
	}

	return append(result, ``)
}

func writeGoasmBody(sub Subroutine, stackArgs StackArgs, arguments, returnValues []string) ([]string, error) {

	var result []string

	for iline, line := range sub.body {

		// If part of epilogue
		if iline >= sub.epilogue.Start && iline < sub.epilogue.End {

			// Instead of last line, output go assembly epilogue
			if iline == sub.epilogue.End - 1 {
				result = append(result, writeGoasmEpilogue(sub, arguments, returnValues)...)
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
		if sub.table.isPresent() {
			line = fixPicLabels(line, sub.table)
		}

		line = fixRbpPlusLoad(line, stackArgs, sub.epilogue.getStackpointerDecrement(sub.table, len(arguments)), sub.epilogue.additionalStackSpace(sub.table, len(arguments)), sub.table.isPresent(), sub.epilogue.AlignedStack)

		detectRbpMinusMemoryAccess(line)
		detectJumpTable(line)
		detectPushInstruction(line)
		detectPopInstruction(line)

		result = append(result, line)
	}

	return result, nil
}

// Write the epilogue for the subroutine
func writeGoasmEpilogue(sub Subroutine, arguments, returnValues []string) []string {

	var result []string

	epilogue, table := &sub.epilogue, &sub.table

	// Restore the stack pointer
	if epilogue.AlignedStack {
		// For an aligned stack, restore the stack pointer from the stack itself
		result = append(result, fmt.Sprintf("    MOVQ %d(SP), SP", epilogue.getStackpointerDecrement(*table, len(arguments))-originalStackPointer))
	} else if epilogue.getStackpointerDecrement(*table, len(arguments)) != 0 {
		// For an unaligned stack, reverse addition in order restore the stack pointer
		result = append(result, fmt.Sprintf("    SUBQ $%d, SP", epilogue.getFreeSpaceAtBottom()))
	}

	// Clear upper half of YMM register, if so done in the original code
	if epilogue.VZeroUpper {
		result = append(result, "    VZEROUPPER")
	}

	if len(returnValues) == 1 {
		// Store return value of subroutine
		result = append(result, fmt.Sprintf("    MOVQ AX, %s+%d(FP)", returnValues[0], getTotalSizeOfArgumentsAndReturnValues(0, len(arguments)-1, returnValues)- 8))
	} else if len(returnValues) > 1 {
		panic(fmt.Sprintf("Fix multiple return values: %s", returnValues))
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
// Remap to rsp/stack pointer or load from golang stack
func fixRbpPlusLoad(line string, stackArgs StackArgs, stackSize uint, tableIsPresent, alignedStack bool) string {

	if match := regexpRbpLoadHigher.FindStringSubmatch(line); len(match) > 1 {
		offset, _ := strconv.Atoi(match[1])
		parts := strings.SplitN(line, match[0], 2)
		if !alignedStack {
			//
			// There is a fixed stack size, so always load from stack pointer
			//
			// +----------------+
			// | ret1           |
			// +----------------+
			// | argN           |
			// +----------------+
			// | arg...         |
			// +----------------+
			// | arg7           |
			// +----------------+
			// | arg6           |     (register passed)
			// +----------------+
			// | arg2 ... arg5  |     (register passed)
			// +----------------+
			// | arg1           |     (register passed)
			// +----------------+
			// | return address |
			// +----------------+
			// |                |
			// |  local         |
			// |                |
			// |  stack         |
			// |                |
			// |                |
			// +----------------+ <-- bottom of stack
			//
			offset = int(stackSize) + returnAddrOnStack + 8*len(registers) + (offset - stackArgs.OffsetToFirst)
			line = parts[0] + fmt.Sprintf("%d[rsp]%s /* %s */", offset, parts[1], match[0])
		} else if tableIsPresent {
			//
			// Base pointer is setup for loading constants, so cannot use rbp
			// Non-register passed stack based arguments are moved to top of the stack
			//
			// +----------------+
			// | ret1           |
			// +----------------+
			// | argN           |
			// +----------------+
			// | arg...         |
			// +----------------+
			// | arg7           |
			// +----------------+
			// | arg6           |     (register passed)
			// +----------------+
			// | arg2 ... arg5  |     (register passed)
			// +----------------+
			// | arg1           |     (register passed)
			// +----------------+
			// | return address |
			// +----------------+
			// | align space... |     (unused)
			// +----------------+
			// |----------------|
			// || original SP  ||
			// |----------------|
			// || argN         ||     (copy)
			// |----------------|
			// || arg...       ||     (copy)
			// |----------------|
			// || arg7         ||     (copy)
			// +----------------+
			// |                |
			// |  local         |
			// |                |
			// |  stack         |
			// |                |
			// |                |
			// +----------------+ <-- aligned address (bottom of stack)
			//
			offset = int(stackSize) + (offset - stackArgs.OffsetToFirst)
			line = parts[0] + fmt.Sprintf("%d[rsp]%s /* %s */", offset, parts[1], match[0])
		} else {
			// Base pointer is equal to (initial) stack pointer, so can leave loads untouched
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

// Detect push instructions
func detectPushInstruction(line string) {

	if match := regexpPushInstr.FindStringSubmatch(line); len(match) > 0 {
		panic(fmt.Sprintf("push instruction detected: %s\n\nCannot modify `rsp` in body of assembly\n\n", match[1]))
	}
}

// Detect pop instructions
func detectPopInstruction(line string) {

	if match := regexpPopInstr.FindStringSubmatch(line); len(match) > 0 {
		panic(fmt.Sprintf("pop instruction detected: %s\n\nCannot modify `rsp` in body of assembly\n\n", match[1]))
	}
}
