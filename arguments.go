package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func ArgumentsOnStack(lines []string) int {

	regexpRbpLoadHigher := regexp.MustCompile(`\[rbp \+ ([0-9]+)\]$`)

	offsets := make(map[int]bool)

	for _, l := range lines {

		if match := regexpRbpLoadHigher.FindStringSubmatch(l); len(match) > 1 {
			offset, _ := strconv.Atoi(match[1])
			if _, found := offsets[offset]; !found {
				offsets[offset] = true
			}
		}
	}

	fmt.Println(offsets)

	return len(offsets)
}

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
